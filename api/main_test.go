package api

import (
	"os"
	"testing"
	"time"

	db "github.com/Long4Changes/MySimpleBank/db/sqlc"
	"github.com/Long4Changes/MySimpleBank/db/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

// TestMain 是 api 包的测试入口
// TestMain 的作用域是包级别
// 每次 go test 运行 api 包测试的时候，都会执行一次 TestMain
// m.Run() 才会真正跑这个包里的所有 Test
// gin.SetMode(gin.DebugMode) --> gin.SetMode(gin.TestMode) 关闭/减少 Gin 的调试输出，使得测试日志更干净、可读
