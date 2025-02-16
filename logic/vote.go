package logic

import (
	"bluebell/Kafka"
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"strconv"

	"go.uber.org/zap"
)

//简化版的投票分数
// 投一票就加432 86400/200 -> 需要200张赞成票可以给你的帖子续一天

// VoteForPost 为帖子投票的函数
/*
direction=1时，有两种情况：
     1.之前没有投过票，现在投赞成票
     2. 之前投反对票，现在改投赞成票
direction=0时，有两种情况:
     1.之前投过赞成票，现在要取消投票
     2.之前投过反对票，现在要取消投票
direction=-1时，有两种情况：
     1。之前没有投过票，现在投反对票
     2. 之前投赞成票，现在改投反对票
更新分数和投票记录
投票的限制：
每个帖子自发表之日起一个星期之内允许用户投票，超过一个星期就不允许再投票了
      1.到期之后将redis中保存的赞成票数及反对票数存储到mysql表中
      2。到期之后删除那个KeyPostVotedZSetPF
*/
func VoteForPost(userID int64, p *models.ParamVoteData) error {
	// 1.判断投票的限制
	//2.更新帖子的分数

	// 3.记录用户为该帖子投票的数据
	zap.L().Debug("VoteForPost",
		zap.Int64("user_id", userID),
		zap.String("postID", p.PostID),
		zap.Int8("direction", p.Direction))
	preValue, err := redis.VoteForPostCheck(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
	if err != nil {
		return err
	}
	var user *models.User
	user, err = mysql.GetUserById(userID)
	if err != nil {
		zap.L().Error("投票：查询用户失败mysql.GetUserById", zap.Error(err))
		return err
	}
	var Post *models.Post
	var postId int64
	postId, err = strconv.ParseInt(p.PostID, 10, 64)
	if err != nil {
		zap.L().Error("投票：ParsePostID失败", zap.Error(err))
		return err
	}
	Post, err = mysql.GetPostAuthorById(postId)
	if err != nil {
		return err
	}
	if p.Direction == 1 {
		err = Kafka.KafkaSendMessage(userID, Post.AuthorID, p.PostID, user.Username, "like_event", "点赞")
	} else if p.Direction == 0 {
		err = Kafka.KafkaSendMessage(userID, Post.AuthorID, p.PostID, user.Username, "like_event", "取消")
	} else if p.Direction == -1 {
		err = Kafka.KafkaSendMessage(userID, Post.AuthorID, p.PostID, user.Username, "like_event", "点踩")
	}
	if err != nil {
		zap.L().Error("投票失败", zap.Error(err))
		return err
	}
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction), preValue)
}
