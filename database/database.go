package database

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
)

var (
	DBConn *gorm.DB
	err    error
)

func MySQLConnect(username, password, host, databaseName string) {
	if host != "" {
		host = fmt.Sprintf("tcp(%s)", host)
	}

	dns := fmt.Sprintf("%s:%s@%s/%s?parseTime=true", username, password, host, databaseName)
	mysqlDB, err := sql.Open("mysql", dns)

	if err != nil {
	}

	DBConn, err = gorm.Open(mysql.New(mysql.Config{
		Conn: mysqlDB,
	}), &gorm.Config{})
}
