package api

import (
	"fmt"

	db "github.com/blessedmadukoma/go-simple-bank/db/sqlc"
	"github.com/blessedmadukoma/go-simple-bank/token"
	"github.com/blessedmadukoma/go-simple-bank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server struct serves HTTP requests for our banking service
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	// tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	// register a custom validation method with gin
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	// do not trust all proxies
	// router.SetTrustedProxies([]string{"192.168.1.2"})
	router.SetTrustedProxies(nil)
	router.TrustedPlatform = gin.PlatformCloudflare

	// routes
	routes := router.Group("/api")
	{

		// users routes
		userRoute := routes.Group("/users")
		{
			userRoute.GET("/:username", server.getUserByUsername)

			// users: auth routes
			userRoute.POST("/register", server.createUser)
			userRoute.POST("/login", server.loginUser)
		}

		// authenticated routes
		authRoute := routes.Group("/").Use(authMiddleware(server.tokenMaker))

		// accounts routes
		authRoute.GET("/accounts/:id", server.getAccount)
		authRoute.GET("/accounts", server.listAccounts)
		authRoute.PUT("/accounts/:id", server.updateAccount)
		authRoute.POST("/accounts", server.createAccount)
		// authRoute.DELETE("/accounts/:id", server.deleteAccount)

		// transfers routes
		authRoute.POST("/transfers", server.createTransfer)

	}

	server.router = router
}

// StartServer runs the HTTP server on a specific address
func (server *Server) StartServer(address string) error {
	return server.router.Run(address)
}

// errorResponse handles the error response by using map[string]interface{} to return the error and it's message
func errorResponse(s string, err error) gin.H {
	return gin.H{"error": s + " -> " + err.Error()}
}
