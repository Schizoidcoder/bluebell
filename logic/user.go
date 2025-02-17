package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
)

//存放业务逻辑的代码
func SignUp(p *models.ParamSignUp) (err error) {
	//判断用户存不存在
	if err := mysql.CheckUserExistByName(p.Username); err != nil {
		//数据库查询出错
		return err
	}
	// 生成UID
	userID := snowflake.GenID()
	//构造一个User实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	//保存进数据库
	err = mysql.InsertUser(user)
	return
}

func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	//传的是指针，就能拿到user.UserID
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	//生成JWT
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return
	}
	user.Token = token
	return
}
