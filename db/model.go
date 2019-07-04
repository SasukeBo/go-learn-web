package db

import (
	"github.com/astaxie/beego/orm"
	"time"
)

// UserInfo model struct
type UserInfo struct {
	UID        int `orm:"PK"`
	Username   string
	Departname string
	Created    time.Time
}

// User model struct
type User struct {
	UID     int `orm:"PK"`
	Name    string
	Profile *Profile `orm:"rel(one)"`
	Post    []*Post  `orm:"reverse(many)"`
}

// Profile model struct
type Profile struct {
	ID   int
	Age  int16
	User *User `orm:"reverse(one)"`
}

// Post model struct
type Post struct {
	ID    int
	Title string
	User  *User  `orm:"rel(fk)"`
	Tags  []*Tag `orm:"rel(m2m)"`
}

// Tag model struct
type Tag struct {
	ID    int
	Name  string
	Posts []*Post `orm:"reverse(many)"`
}

func init() {
	// 注册定义的 model
	// orm.RegisterModel(new(UserInfo))
	// 也可以同时注册多个model
	// orm.RegisterModel(new(User), new(Profile), new(Post))

	// 创建 table
	// orm.RunSyncdb("default", false, true)

	// 需要在init中注册定义的model
	orm.RegisterModel(new(UserInfo), new(User), new(Profile), new(Post), new(Tag))

	// 创建 table
	orm.RunSyncdb("default", false, true)
}
