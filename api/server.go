package api

import (
	db "github.com/blessedmadukoma/go-simple-bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server struct serves HTTP requests for our banking service
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	// do not trust all proxies
	// router.SetTrustedProxies([]string{"192.168.1.2"})
	router.SetTrustedProxies(nil)
	router.TrustedPlatform = gin.PlatformCloudflare

	// create routes
	router.POST("/api/accounts", server.createAccount)
	router.GET("/api/accounts/:id", server.getAccount)
	router.GET("/api/accounts", server.listAccounts)
	router.PUT("/api/accounts/:id", server.updateAccount)
	// router.DELETE("/accounts/:id", server.deleteAccount)

	server.router = router

	return server
}

// StartServer runs the HTTP server on a specific address
func (server *Server) StartServer(address string) error {
	return server.router.Run(address)
}

// errorResponse handles the error response by using map[string]interface{} to return the error and it's message
func errorResponse(s string, err error) gin.H {
	return gin.H{"error": s + " -> " + err.Error()}
}
