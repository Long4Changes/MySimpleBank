package main

// _"github.com/lib/pq" 是空白导入
// Go的 database/sql 是通用接口层，本身不包含具体数据库的驱动
// pq 包再 init() 里会把 PostgresSQL 驱动注册到 database/sql
// 之后再调用 sql.Open("postgres", ....) 时，"postgres"这个驱动名才能被识别
import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/techschool/simplebank/api"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/db/util"
)

// 有了 viper 这段配置的硬编码就不需要了
/*
const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:Lc@200409200034@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)
*/

func main() {
	// "." 表示当前文件夹
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	// 使用 viper 下面的字段需要换掉
	// dbDriver --> config.DBDriver
	// dbSource --> config.DBSource
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	// serverAddress --> config.ServerAddress
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot connect to server:", err)
	}

}
