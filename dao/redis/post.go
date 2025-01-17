package redis

import (
	"bluebell/models"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func getIDsFromKey(key string, page, size int64) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1
	//Zrevrange 按分数从大到小的顺序查询指定数量的元素
	return rdb.ZRevRange(key, start, end).Result()
}
func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	//从redis 获取id
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	//确定查询的索引起始点
	return getIDsFromKey(key, p.Page, p.Size)
}

// GetPostVoteData 根据ids查询每篇帖子的投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	//data = make([]int64, 0, len(ids))
	//for _, id := range ids {
	//	key := getRedisKey(KeyPostVotedZSetPF + id)
	//	// 查找key中分数是1的元素的数量
	//	v1 := rdb.ZCount(key, "1", "1").Val()
	//	data = append(data, v1)
	//}

	// 使用pipeline一次发送多条命令减少RTT
	pipeline := rdb.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPF + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(ids))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// GetCommunityPostIDsInOrder 按社区根据ids查询每篇帖子的投赞成票的数据
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 使用 zinterstore 把分区的帖子set与帖子分数的zset 生成一个新的zset
	// 针对新的zset 按之前的逻辑取数据
	//社区的key
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityID)))
	//利用缓存key减少zinterstore执行的次数
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	if rdb.Exists(key).Val() < 1 {
		//不存在 需要计算
		pipeline := rdb.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, cKey, orderKey)                   // zinterstore 计算
		pipeline.Expire(key, 60*time.Second) //设置超时时间
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	return getIDsFromKey(key, p.Page, p.Size)

}
