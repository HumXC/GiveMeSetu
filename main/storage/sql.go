package storage

import (
	"database/sql"
	"give-me-setu/main/conf"
	"give-me-setu/util"
	"path"

	"github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

// setu 保存每个色图
// tags 保存标签
type Tags interface {
	AllTags() map[string]string
}
type Libs interface {
}

type ServeDB struct {
	database *sql.DB
	Lib      Libs
	Tag      Tags
}

func (d *ServeDB) Close() {
	d.database.Close()
}

const (
	CREATE_DATABASE         = "CREATE DATABASE ? DEFAULT CHARACTER SET = 'utf8mb4';"
	CREATE_TABLE_SETU_MYSQL = `
		CREATE TABLE IF NOT EXISTS setu(  
			id INT NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
			create_time DATETIME NOT NULL COMMENT 'Create Time',
			name VARCHAR(255) NOT NULL COMMENT 'File Name On Storage',
			origin VARCHAR(255) NOT NULL COMMENT 'Where Is It From',
			tags VARCHAR(1023)
		)`
	CREATE_TABLE_SETU_SQLITE = `
		CREATE TABLE IF NOT EXISTS setu(  
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			create_time DATETIME NOT NULL,
			name TEXT NOT NULL,
			origin TEXT NOT NULL,
			tags TEXT
		)`
	CREATE_TABLE_TAG_SQLITE = `
		CREATE TABLE IF NOT EXISTS tag(  
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			create_time DATETIME NOT NULL,
			name TEXT NOT NULL,
			isDisable BOOLEAN NOT NULL
		)`
)

var db *sql.DB

// 初始化一个数据库
func GetDB(cfg conf.Config) (*ServeDB, error) {
	if db != nil {
		db.Close()
		panic("Allready have one db!")
	}
	switch cfg.Database.Driver {
	case "mysql":
		url := cfg.Database.User + ":" + cfg.Database.Password + "@tcp(" + cfg.Database.Host + ")/" + cfg.Database.Name
		db_, err := sql.Open("mysql", url)
		if err != nil {
			return nil, err
		}
		err = db_.Ping()
		if err != nil {
			sqlErr := err.(*mysql.MySQLError)
			// 数据库不存在
			if sqlErr.Number == 1049 {
				url := cfg.Database.User + ":" + cfg.Database.Password + "@tcp(" + cfg.Database.Host + ")/"
				mysql, err := sql.Open("mysql", url)
				if err != nil {
					return nil, err
				}
				err = mysql.Ping()
				if err != nil {
					return nil, err
				}
				s := util.Replace(CREATE_DATABASE, "?", cfg.Database.Name)
				_, err = mysql.Exec(s)
				if err != nil {
					return nil, err
				}
				mysql.Close()
			}
		}
		err = db_.Ping()
		if err != nil {
			return nil, err
		}
		_, err = db_.Exec(CREATE_TABLE_SETU_MYSQL)
		if err != nil {
			return nil, err
		}
		db = db_

	default:
		db_, err := sql.Open("sqlite3", path.Join(cfg.DataDir, cfg.Database.Name+".db"))
		if err != nil {
			return nil, err
		}
		err = db_.Ping()
		if err != nil {
			return nil, err
		}
		_, err = db_.Exec(CREATE_TABLE_SETU_SQLITE)
		if err != nil {
			return nil, err
		}
		_, err = db_.Exec(CREATE_TABLE_TAG_SQLITE)
		if err != nil {
			return nil, err
		}
		db = db_
	}
	return &ServeDB{
		database: db,
	}, nil
}
