package api

import (
	db "github.com/Long4Changes/MySimpleBank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for our banking service.
// 引入 gomock 后 store 代表的不再是一个结构体类型，而是一个接口，所以不需要用指针了
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// 注册自定义验证器
	// binding.Validator 是 Gin 暴露的验证器入口
	// Engine() 返回底层引擎，但返回值是 any 类型，编译器不知道它具体有哪些方法
	// 由于我需要调用 RegisterValidation，而这个方法定义在 *validator.Validate
	// 所以需要做类型断言
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

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
	// 转账
	router.POST("/transfers", server.createTransfer)

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
