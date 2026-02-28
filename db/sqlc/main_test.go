package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/techschool/simplebank/db/util"
)

// 有了 viper 下面这段硬编码就不需要了
/*
const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:Lc@200409200034@localhost:5432/simple_bank?sslmode=disable"
)
*/

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	// "../.." 表示上一级的上一级
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cannot load config:", err)	
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
