package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/snowflake"
)

//存放业务逻辑的代码
func SignUp(p *models.ParamSignUp) {
	//判断用户存不存在
	mysql.QueryUserByUsername()
	// 生成UID
	snowflake.GenID()
	//保存进数据库
	mysql.InsertUser()
}
