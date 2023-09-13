package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/blessedmadukoma/go-simple-bank/api"
	db "github.com/blessedmadukoma/go-simple-bank/db/sqlc"
	"github.com/blessedmadukoma/go-simple-bank/gapi"
	"github.com/blessedmadukoma/go-simple-bank/pb"
	"github.com/blessedmadukoma/go-simple-bank/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	_ "github.com/blessedmadukoma/go-simple-bank/docs/statik"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	// connect to database
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	// runGinServer(config, store)
	go runGatewayServer(config, store) // run in a separate goroutine
	runGrpcServer(config, store)

}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.StartServer(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server!")
	}
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)

	// optional but HIGHLY RECOMMENDED
	reflection.Register(grpcServer) // explore what rpc methods are available and easily call them

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}

	log.Printf("Starting gRPC server at %s\n", config.GRPCServerAddress)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server:", err)
	}
}

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	jsonOptions := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOptions)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register gateway server:", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// Swagger UI
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal("cannot create statik file system: ", err)
	}

	swaggerHandler := http.StripPrefix("/swagger", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}

	log.Printf("Starting HTTP gateway server at %s\n", config.HTTPServerAddress)
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start HTTP gateway server:", err)
	}
}
