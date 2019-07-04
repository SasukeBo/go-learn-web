package db

import (
	// "database/sql"
	"github.com/astaxie/beego/orm"
	// 导入驱动
	_ "github.com/lib/pq"
)

func init() {
	// 注册驱动
	orm.RegisterDriver("postgres", orm.DRPostgres)

	// 设置默认数据库
	orm.RegisterDataBase("default", "postgres",
		`
		user=sasuke
		password=Wb922149@...S
		dbname=sasuke
		host=127.0.0.1
		port=5432
		sslmode=disable
	`)
}
