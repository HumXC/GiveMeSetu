package database

import "database/sql"

// 保存在 library 表里的结构
type Setu struct {
	Name    string
	Origin  string
	Library string
	Tags    []string
}

type ImgLibDB struct {
	database *sql.DB
}

func (d *ImgLibDB) Close() {
	d.database.Close()
}

func (d *ImgLibDB) GetASetu(libName string, setuName string) {

}
