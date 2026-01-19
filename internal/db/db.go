package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func Init(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25) // 最大打开连接数
	db.SetMaxIdleConns(25) // 最大空闲连接数
	db.SetConnMaxLifetime(5 * time.Minute) // 连接最大存活时间

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("数据库连接失败：%w", err)
	}

	log.Println("数据库连接成功")
	return db, nil
}