package redis

const (
	KeyPrefix          = "bluebell:"
	KeyPostTime        = "post:time"  // Zset; post and time
	KeyPostScore       = "post:score" // Zset; post and score
	KeyPostVotedPrefix = "post:voted" // Zset; record the user and voting; the follwoing param is the post id
	KeyCommunityPf     = "community:"
)

func GetRedisKey(key string) string {
	return KeyPrefix + key
}
