package api

import (
	"database/sql"
	"net/http"

	db "github.com/Long4Changes/MySimpleBank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// 看到就复习一下
// request就是传入， response就是传出
type createAccountRequest struct {
	// binding 用于表示这个字段是必须的
	Owner string `json:"owner" binding:"required"`
	// 这里的 oneof 用于表示 Currency 只能是USD和EUR其中之一
	// oneof = USD EUR CAD
	// 引入 validator 之后就不再需要 oneof 了
	Currency string `json:"currency" binding:"required,currency"`
}

// 从前端传过来的数据有三种形式
// 1. JSON 数据
// 2. Uri Parameters
// 3. Query Parameters
func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	// 系统尝试把前端传过来的JSON数据，强行塞进你定义的 createAccountRequest结构体当中
	// 其实是做一个数据的转换，根据前段传过来参数的类型，后端定义能接住这种类型的结构
	// 不要忘记 ShouldBindJSON() 里面传的是 地址 &，
	// ShouldBindJSON() 要把解析出的 JSON 字段写进去，所以要传指针
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	// 这里实际调用的是 server.store.Queries.CreateAccount(...)
	account, err := server.store.CreateAccount(ctx, arg)
	// Handlel DB errors in Go
	// 把 err 从 error 类型转换为 PostgreSQL 驱动的具体错误类型 *pq.Error
	// ok == true 说明这次错误的确来自 PostgreSQL，于是打印 pqErr.Code.Name() 用于调试日志
	// log.Println(pqErr.Code.Name())
	// 比如这里打印的 foreign_key_violation
	// ok == false 说明不是 *pq.Error 就不打印这行日志
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
		// foreign_key_violation: 给一个没有在 user 表有 username 的 owner 创建 account
		// unique_violation: 重复给同一个 owner 创建相同 currency 的 account
		// 如果不按照下面的形式，上面的两种错误都会返回 InternalServerError 500
		// 500 指服务端错误，但实际是客户端的问题，所以处理一下返回 403 更合适
		// 通过 *pq.Error 判断是否是 PostgresSQL 的错误，如果是则对数据库的错误进行特殊处理 
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	// ID 属于 Uri Params 所以标签里要用 uri
	// 这里的限制是让 ID 的最小值为1，也就是避免负数
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	// 把前端传过来的Uri参数，塞进 getAccountRequest 结构体当中
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		// 如果传进来的 ID 没有找到记录
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

// 分页
type listAccountRequest struct {
	// PageID和PageSize都属于 Query Params 所以标签里要用 form
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccount(ctx *gin.Context) {
	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListAccountsParams{
		// Limit: 每一页最多返回多少数据
		Limit: req.PageSize,
		// Offset: 跳过前面多少条数据再开始取
		Offset: (req.PageID - 1) * req.PageSize,
	}
	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

// 更新账户
type UpdateAccountRequest struct {
	ID      int64 `json:"id" binding:"required,min=1"`
	Balance int64 `json:"balance" binding:"required"`
}

func (server *Server) updateAccount(ctx *gin.Context) {
	var req UpdateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateAccountParams{
		ID:      req.ID,
		Balance: req.Balance,
	}

	account, err := server.store.UpdateAccount(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

// 删除账户
type deleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req deleteAccountRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := server.store.DeleteAccount(ctx, req.ID); err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// gin.H{...} 用于快速构造 JSON 响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "account deleted successfully",
		"id":      req.ID,
	})
}
