package mysql

import (
	"bluebell/models"
	"strings"

	"github.com/jmoiron/sqlx"
)

func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post(post_id,title,content,author_id,community_id) values(?,?,?,?,?)`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

func GetPostById(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `select 
    post_id, title, content,author_id,community_id,create_time
    from post where post_id = ?`
	err = db.Get(post, sqlStr, pid)
	return
}

func GetPostAuthorById(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `select author_id from post where post_id = ?`
	err = db.Get(post, sqlStr, pid)
	return
}

//根据id获取用户信息
func GetUserById(id int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select
    user_id,username from user where user_id = ?`
	err = db.Get(user, sqlStr, id)
	return
}

func GetPostList(page, size int64) (posts []*models.Post, err error) {
	sqlStr := `select 
    post_id, title, content,author_id,community_id,create_time
    from post
    ORDER BY create_time DESC
    limit ?,?
    `
	posts = make([]*models.Post, 0, size)
	err = db.Select(&posts, sqlStr, (page-1)*size, size)
	return

}

//根据给定的id列表查询帖子数据
func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	sqlStr := `select post_id, title, content,author_id,community_id,create_time
    from post
    where post_id in (?)
    order by FIND_IN_SET(post_id,?)
    `
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)
	err = db.Select(&postList, query, args...)
	return
}
