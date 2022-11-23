package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"give-me-setu/main/conf"
	"give-me-setu/util"
	"path"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

// setu 保存每个色图
// tags 保存标签
type TagDBs interface {
	AllTags() map[string]string
}
type SetuDBs interface {
	Add(setu Setu) error
	Del(id string) error
	GetByID(id string) (Setu, error)
	GetByIDs(ids []string) ([]Setu, error)
	Mod(setu Setu) error
}

type ServeDB struct {
	database *sql.DB
	SetuDB   SetuDBs
}

func (d *ServeDB) Close() {
	d.database.Close()
	d.SetuDB = nil
	db = nil
}

const (
	CREATE_DATABASE_MYSQL   = "CREATE DATABASE ? DEFAULT CHARACTER SET = 'utf8mb4'"
	CREATE_TABLE_SETU_MYSQL = `
		CREATE TABLE IF NOT EXISTS setu(  
			id VARCHAR(255) NOT NULL PRIMARY KEY COMMENT 'Primary Key',
			create_time DATETIME NOT NULL COMMENT 'Create Time',
			modify_time DATETIME NOT NULL COMMENT 'Modify Time',
			title VARCHAR(255) NOT NULL COMMENT 'File Name On Storage',
			origin VARCHAR(255) NOT NULL COMMENT 'Where Is It From',
			tags VARCHAR(1023)
		)`
	CREATE_TABLE_SETU_SQLITE = `
		CREATE TABLE IF NOT EXISTS setu (  
			id TEXT NOT NULL PRIMARY KEY,
			create_time DATETIME NOT NULL,
			modify_time DATETIME NOT NULL,
			title TEXT NOT NULL,
			origin TEXT NOT NULL,
			ext TEXT NOT NULL
		)`
	ADD_SETU_SQLITE = `
		INSERT INTO setu (
			id, create_time, modify_time, title, origin, ext
		) VALUES (?,?,?,?,?,?)
	`
	DEL_SETU_SQLITE = `
		DELETE FROM setu WHERE id=?
	`
	MOD_SETU_SQLITE = `
		UPDATE setu SET modify_time=?, title=?, origin=? WHERE id=?
	`
	GET_SETU_BY_ID = `
		SELECT id, title, origin, ext, create_time, modify_time FROM setu WHERE id=?
	`
	GET_SETU_BY_IDS = `
		SELECT id, title, origin, ext, create_time, modify_time FROM setu WHERE id IN (%s)
	`
)

// setu 表的数据库操作
type SetuDB struct {
	*sql.DB
}

// 修改这个结构体需要同时检查 sql 语句
type Setu struct {
	ID         string
	Title      string
	Origin     string
	Ext        string
	CreateTime int64
	ModTime    int64
}

// 添加一个条目到 setu 表
// 参数为要添加到表里的实际内容，其中 id 为文件的 md5 校验和，也是文件存储在磁盘里的文件名
func (l *SetuDB) Add(setu Setu) error {
	createTime := time.Now().Unix()
	_, err := l.Exec(ADD_SETU_SQLITE, setu.ID, createTime, createTime, setu.Title, setu.Origin, setu.Ext)
	if err != nil {
		return err
	}
	return nil
}

func (l *SetuDB) GetByID(id string) (Setu, error) {
	s := Setu{}
	rows, err := l.Query(GET_SETU_BY_ID, id)
	if err != nil {
		return s, err
	}
	defer rows.Close()

	ok := rows.Next()
	if !ok {
		return s, errors.New("Cam not find this id: " + id)
	}
	var cTime time.Time
	var mTime time.Time
	err = rows.Scan(&s.ID, &s.Title, &s.Origin, &s.Ext, &cTime, &mTime)
	if err != nil {
		return s, err
	}
	s.CreateTime = cTime.Unix()
	s.ModTime = mTime.Unix()
	return s, nil
}
func (l *SetuDB) GetByIDs(ids []string) ([]Setu, error) {
	result := make([]Setu, 0)
	leng := len(ids)
	var idss = make([]any, leng)
	for i := 0; i < leng; i++ {
		idss[i] = ids[i]
	}
	sql := fmt.Sprintf(GET_SETU_BY_IDS, strings.Repeat("?,", leng-1)+"?")
	rows, err := l.Query(sql, idss...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		i++
		s := Setu{}
		var cTime time.Time
		var mTime time.Time
		err := rows.Scan(&s.ID, &s.Title, &s.Origin, &s.Ext, &cTime, &mTime)
		if err != nil {
			return nil, err
		}
		s.CreateTime = cTime.Unix()
		s.ModTime = mTime.Unix()
		result = append(result, s)
	}

	if i != len(ids) {
		return make([]Setu, 0), fmt.Errorf("want %d items, but got %d", len(ids), i)
	}
	return result, nil
}

// 修改数据
func (l *SetuDB) Mod(setu Setu) error {
	modTime := time.Now().Unix()
	_, err := l.Exec(MOD_SETU_SQLITE, modTime, setu.Title, setu.Origin, setu.Ext, setu.ID)
	if err != nil {
		return err
	}
	return nil
}

// 从 setu 表删除一个条目
func (l *SetuDB) Del(id string) error {
	_, err := l.Exec(DEL_SETU_SQLITE, id)
	if err != nil {
		return err
	}
	return nil
}

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
				s := util.Replace(CREATE_DATABASE_MYSQL, "?", cfg.Database.Name)
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
		db = db_
	}
	return &ServeDB{
		database: db,
		SetuDB:   &SetuDB{db},
	}, nil
}
