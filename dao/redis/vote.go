package redis

import (
	"errors"
	"math"
	"time"

	"github.com/go-redis/redis"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432
)

var (
	ErrVoteTimeExpired = errors.New("voting time is over")
	ErrVoteRepeated    = errors.New("can't vote repeatedly")
)

func VoteForPost(userID, postID string, value float64) (err error) {
	//1. estimate if this vote satisfies the vote limitation
	postTime := rdb.ZScore(GetRedisKey(KeyPostTime), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpired
	}

	// The following 2 3 should be seen as a transaction
	pipeLine := rdb.TxPipeline()
	//2. update the post scores
	// before execute the vote logic, get the user's voting record
	ov := rdb.ZScore(GetRedisKey(KeyPostVotedPrefix+postID), userID).Val()
	diff := math.Abs(value - ov)
	var op float64
	if value == ov {
		return ErrVoteRepeated
	}
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	pipeLine.ZIncrBy(GetRedisKey(KeyPostScore), op*diff*scorePerVote, postID)

	//3. record the user's voting data
	if value == 0 {
		pipeLine.ZRem(KeyPostVotedPrefix+postID, userID)

	} else {
		pipeLine.ZAdd(GetRedisKey(KeyPostVotedPrefix+postID), redis.Z{
			Score:  value,
			Member: userID,
		})
	}
	_, err = pipeLine.Exec()
	return err
}
