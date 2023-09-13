package logic

import (
	"strconv"

	"web_app/dao/redis"
	"web_app/models"

	"go.uber.org/zap"
)

// The vote foundation
//1. the data of user's voting

/*
There are some conditions of voting:
when direction = 1, there are 2 situations,
when direction = 0(which means cancelling the past voting), there are 2 situations,
when direction = -1, there are 2 situations.

The limitations of voting:
	1 weeks after the post created, user lose the right of voting:
		After the timeout, the redis Zset store into mysql consistently,
		And delete the redis Zset.
*/

// VoteForPost Vote for post
/* This project use the simple Vote algorithms,
One vote = 432 scores, which means 86400/200, the post need 200 votes every day in order to catch up the timestamp scores
*/
func VoteForPost(userID int64, p *models.ParamVoteData) (err error) {
	zap.L().Debug(
		"VoteForPost",
		zap.Int64("userid", userID),
		zap.String("postID", p.PostID),
		zap.Int8("direction", p.Direction),
	)
	return redis.VoteForPost(strconv.FormatInt(userID, 10), p.PostID, float64(p.Direction))
}
