package api

import (
	"fmt"

	db "github.com/Long4Changes/MySimpleBank/db/sqlc"
	"github.com/Long4Changes/MySimpleBank/db/util"
	"github.com/Long4Changes/MySimpleBank/token"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for our banking service.
// 引入 gomock 后 store 代表的不再是一个结构体类型，而是一个接口，所以不需要用指针了
// implement login user API that returns PASETO access token
// 加入 config 和 tokenMaker 
// 这里修改之后，之前编写的 API test 文件就需要修改了
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
// 将 TokenSymmetricKey 和 AccessTokenDuration 写入 app.env 后，还加入了 config.go
// 后面从 config.go 中读取数据
// 由于 NewServer 底下有很多路由变得太长了，所以现在把路由单独拆分进 setupRouter 当中
func NewServer(config util.Config, store db.Store) (*Server, error) {
	// 只要将这里 NewPasetoMaker 改成 NewJWTMaker 就可以更换 token 类型
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	// %w 用于包装原始错误
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	// 这里实例化 server 的时候也要加入 config 和 tokenMaker
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	
	// 这一行 加入 setupRouter 了
	//router := gin.Default()

	// 注册自定义验证器
	// binding.Validator 是 Gin 暴露的验证器入口
	// Engine() 返回底层引擎，但返回值是 any 类型，编译器不知道它具体有哪些方法
	// 由于我需要调用 RegisterValidation，而这个方法定义在 *validator.Validate
	// 所以需要做类型断言
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	// 所有的 router 都加入 setupRouter

	server.setupRouter()
	// 这一行也加入 setupRouter 了
	//server.router = router
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	// add routes to router
	// 创建一个用户
	router.POST("/users", server.createUser)
	// 用户登录
	router.POST("/users/login", server.loginUser)
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
}
// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// H is the abbreviation of map[string]interface{}
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
