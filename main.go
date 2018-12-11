package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"

	_ "howlong/routers"
)

func init() {
	//设置默认数据库
	orm.RegisterDataBase("default", "sqlite3", "./howlong.db")
	// 创建table
	orm.RunSyncdb("default", false, true)
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	} else {
		beego.SetLevel(beego.LevelInformational)
	}
	beego.Run()
}
