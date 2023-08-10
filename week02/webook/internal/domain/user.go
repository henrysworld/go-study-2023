package domain

import (
	"time"
)

// User 领域对象，是 DDD 中的 entity
// BO(business object)
type User struct {
	Id       int64
	NickName string
	Birthday string
	Bio      string
	Email    string
	Password string
	Ctime    time.Time
}

//type Address struct {
//}
