package main

import (
	"fmt"
	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"

	_ "howlong/routers"
)

func init() {
	//设置默认数据库
	dbPath := os.Getenv("APP_DB_SQLITE_PATH")
	if dbPath == "" {
		dbPath = "./howlong.db"
	}
	fmt.Println("Load sqlite db from ", dbPath)

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
