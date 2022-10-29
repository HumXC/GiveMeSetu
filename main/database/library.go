/*
 * @Author: HumXC Hum-XC@outlook.com
 * @Date: 2022-10-29
 * @LastEditors: HumXC Hum-XC@outlook.com
 * @LastEditTime: 2022-10-29
 * @FilePath: /give-me-setu/main/database/library.go
 * @Description: 对 图库 的定义
 *
 * Copyright (c) 2022 by HumXC Hum-XC@outlook.com, All Rights Reserved.
 */
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
