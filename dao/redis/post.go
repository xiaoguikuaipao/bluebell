package redis

import (
	"strconv"
	"time"

	"web_app/models"

	"github.com/go-redis/redis"
)

func CreatePost(postID int64, communityID int64) (err error) {
	// the following operation is supposed to be dealt with as a transaction
	pipeLine := rdb.TxPipeline()
	//this key is used to record the created time
	pipeLine.ZAdd(GetRedisKey(KeyPostTime), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// when a post created, its initial score was the timestamp
	pipeLine.ZAdd(GetRedisKey(KeyPostScore), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// add the post ID to the community set
	// 用来和post分数或post时间的zset做interstore
	cKey := GetRedisKey(KeyCommunityPf + strconv.FormatInt(communityID, 10))
	pipeLine.SAdd(cKey, postID)

	_, err = pipeLine.Exec()
	return err
}

func getIDsFromKey(key string, page, size int64) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1
	return rdb.ZRevRange(key, start, end).Result()

}

func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	key := GetRedisKey(KeyPostTime)
	if p.Order == models.OrderScore {
		key = GetRedisKey(KeyPostScore)
	}
	// retrieve the redis
	return getIDsFromKey(key, p.Page, p.Size)
}

func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// Use the zinterstore, combine the community set and order zset
	cKey := GetRedisKey(KeyCommunityPf + strconv.FormatInt(p.CommunityID, 10))
	orderKey := GetRedisKey(KeyPostTime)
	if p.Order == models.OrderScore {
		orderKey = GetRedisKey(KeyPostScore)
	}

	//The Key is a temp set, which is a zinterstore
	Key := orderKey + strconv.FormatInt(p.CommunityID, 10)
	if rdb.Exists(Key).Val() < 1 {
		pipeline := rdb.Pipeline()
		rdb.ZInterStore(Key, redis.ZStore{
			Aggregate: "MAX",
		}, cKey, orderKey)
		// Use the cache key to reduce the overhead of zinterstore
		pipeline.Expire(Key, 60*time.Second)
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	return getIDsFromKey(Key, p.Page, p.Size)

}

func GetPostVoteData(ids []string) (data []int64, err error) {
	// Use pipeline to reduce the RTT
	// distinguish the TxPipeline and Pipeline
	data = make([]int64, 0, len(ids))
	pipeline := rdb.Pipeline()
	for _, id := range ids {
		key := GetRedisKey(KeyPostVotedPrefix + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	//for _, id := range ids {
	//	key := GetRedisKey(KeyPostVotedPrefix + id)
	//	v := rdb.ZCount(key, "1", "1").Val()
	//	data = append(data, v)
	//}
	return
}
