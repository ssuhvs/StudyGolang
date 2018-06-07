package dao

import (
	"github.com/jinzhu/gorm"
	"../models"
	"../tools"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"../global"
)

var (
	shdbname = global.Config.DbInfo.SeekHelpDb     //用户数据库
	shdbtye  = global.Config.DbInfo.SeekHelpDbType //数据库类型
	shdb     *gorm.DB                              //数据库连接
)

func init() {
	global.WgDb.Add(1)
	go SeekHelpDbInit()
}

//初始化
func SeekHelpDbInit() {
	shdbname = global.CurrPath + shdbname
	//fmt.Println("用户数据库地址:",logindbname)
	tdb, err := gorm.Open(shdbtye, shdbname)
	tools.PanicErr(err, "帮助数据库初始化")
	shdb = tdb
	if !shdb.HasTable(&models.SeekHelp{}) {
		shdb.CreateTable(&models.SeekHelp{})
	}
	//fmt.Println("帮助数据库初始化完成")
	global.WgDb.Done()
}

//发布求助
func PublishSeekHelp()(sid int,err error){
	return
}