package logic

import (
	"fmt"
	"github.com/SasukeBo/learn-web/db"
	"github.com/astaxie/beego/orm"
	"time"
)

// TestDB do some db test
func TestDB() {
	o := orm.NewOrm()
	userInfo := db.UserInfo{
		Username:   "sasuke",
		Departname: "自动化产品处",
		UID:        123,
		Created:    time.Now(),
	}

	// 插入表
	id, err := o.Insert(&userInfo)
	fmt.Printf("ID: %d, ERR: %v\n", id, err)

	// 更新表
	userInfo.Username = "SasukeBo"
	num, err := o.Update(&userInfo)
	fmt.Printf("NUM: %d, ERR: %v\n", num, err)

	// 读取 one
	u := db.UserInfo{UID: userInfo.UID}
	err = o.Read(&u)
	fmt.Printf("ERR: %v\n", err)

	// 删除行
	num, err = o.Delete(&u)
	fmt.Printf("NUM: %d, ERR: %v\n", num, err)
}
