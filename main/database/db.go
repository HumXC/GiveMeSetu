/*
 * @Author: HumXC Hum-XC@outlook.com
 * @Date: 2022-10-27
 * @LastEditors: HumXC Hum-XC@outlook.com
 * @LastEditTime: 2022-10-29
 * @FilePath: /give-me-setu/main/database/db.go
 * @Description: 初始化数据库
 *
 * Copyright (c) 2022 by HumXC Hum-XC@outlook.com, All Rights Reserved.
 */
package database

import (
	"database/sql"
	"give-me-setu/main/conf"
	"give-me-setu/util"
	"log"
	"path"

	"github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

const (
	CREATE_DATABASE      = "CREATE DATABASE ? DEFAULT CHARACTER SET = 'utf8mb4';"
	CREATE_TABLE_LIBRARY = `
		CREATE TABLE IF NOT EXISTS library(  
			id INT NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
			create_time DATETIME NOT NULL COMMENT 'Create Time',
			name VARCHAR(255) NOT NULL COMMENT 'File Name On Storage',
			origin VARCHAR(255) NOT NULL COMMENT 'Where Is It From',
			tags VARCHAR(1023)
		)`
)

type DB struct {
	database *sql.DB
}

func (d *DB) Close() {
	d.database.Close()
}
func Get(cfg conf.Config) *DB {
	var db *sql.DB
	switch cfg.Database.Driver {
	case "mysql":
		url := cfg.Database.User + ":" + cfg.Database.Password + "@tcp(" + cfg.Database.Host + ")/" + cfg.Database.Name
		db_, err := sql.Open("mysql", url)
		if err != nil {
			log.Fatal(err)
		}
		err = db_.Ping()
		if err != nil {
			sqlErr := err.(*mysql.MySQLError)
			// 数据库不存在
			if sqlErr.Number == 1049 {
				url := cfg.Database.User + ":" + cfg.Database.Password + "@tcp(" + cfg.Database.Host + ")/"
				mysql, err := sql.Open("mysql", url)
				if err != nil {
					log.Fatal(err)
				}
				err = mysql.Ping()
				if err != nil {
					log.Fatal(err)
				}
				s := util.Replace(CREATE_DATABASE, "?", cfg.Database.Name)
				_, err = mysql.Exec(s)
				if err != nil {
					log.Fatal(err)
				}
				mysql.Close()
			}
		}
		err = db_.Ping()
		if err != nil {
			log.Fatal(err)
		}
		db = db_

	default:
		db_, err := sql.Open("sqlite3", path.Join(cfg.DataDir, cfg.Database.Name+".db"))
		if err != nil {
			log.Fatal(err)
		}
		db = db_
	}
	return &DB{
		database: db,
	}
}
func InitDB(db *DB) error {
	err := db.database.Ping()
	if err != nil {
		return err
	}
	_, err = db.database.Exec(CREATE_TABLE_LIBRARY)
	return err
}
