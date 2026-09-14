package api

import (
	db "simple_bank/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for our banking service
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server with the given store and sets up routing.
func NewServer(store db.Store) *Server {
	server := &Server{
		store:  store,
		router: gin.Default(),
	}

	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.router.POST("/accounts", server.createAccount)
	server.router.GET("/accounts/:id", server.getAccount)
	server.router.GET("/accounts", server.listAccounts)
	server.router.PUT("/accounts/:id", server.updateAccount)
	server.router.DELETE("/accounts/:id", server.deleteAccount)

	server.router.POST("/transfers", server.createTransfer)

	return server
}

// Start runs the HTTP server on the specified address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// errorResponse formats an error into a JSON response.
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
