package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/techschool/simplebank/db/sqlc"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// add routes to router
	// 创建一个账户
	router.POST("/accounts", server.createAccount)
	// 根据 id 查询一个账户
	router.GET("/accounts/:id", server.getAccount)
	// 分页展示账户
	router.GET("/accounts", server.listAccount)
	// 更新账户
	router.POST("/accounts/update", server.updateAccount)
	// 删除账户
	router.POST("/accounts/delete/:id", server.deleteAccount)

	server.router = router
	return server
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
// H is the abbreviation of map[string]interface{}
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
