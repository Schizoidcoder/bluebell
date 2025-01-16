package redis

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

const (
	oneWeekInSeconds = 7 * 24 * 60 * 60
	scorePerVote     = 432 //每一票值多少分
)

var (
	ErrVotesTimeExpire = errors.New("投票时间已过")
)

func CreatePost(postID int64) error {
	pipeline := rdb.TxPipeline()
	//帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	//帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	_, err := pipeline.Exec()
	return err
}

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
func VoteForPost(userID, postID string, value float64) error {
	// 1.判断投票的限制
	//去redis取帖子发布时间
	postTime := rdb.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVotesTimeExpire
	}
	//2和3需要放到一个pipeline事务中操作
	//2.更新帖子的分数
	//先查当前用户给当前帖子的投票记录
	ov := rdb.ZScore(getRedisKey(KeyPostVotedZSetPF+postID), userID).Val()
	//计算两次投票的差值
	diff := value - ov
	fmt.Println(value, ov, diff)
	pipeline := rdb.TxPipeline()
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), diff*scorePerVote, postID).Result()

	// 3.记录用户为该帖子投票的数据
	if value == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPF+postID), userID).Result()
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPF+postID), redis.Z{
			Score:  value,
			Member: userID,
		}).Result()
	}
	_, err := pipeline.Exec()
	return err
}
