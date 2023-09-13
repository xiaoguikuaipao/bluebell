package logic

import (
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/models"
	"web_app/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) error {
	//1. generate the post id
	p.ID = snowflake.GenID()

	//2. stores into db
	err := mysql.CreatePost(p)
	if err != nil {
		return err
	}

	err = redis.CreatePost(p.ID, p.CommunityID)
	if err != nil {
		return err
	}
	//3. return
	return err
}

//func GetPostByID(pid int64) (*models.Post, error) {
//	return mysql.GetPostByID(pid)
//}

func GetPostByID(pid int64) (*models.ApiPostDetail, error) {
	data := new(models.ApiPostDetail)
	// search data and combine the data
	post, err := mysql.GetPostByID(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostByID failed",
			zap.Int64("pid", pid),
			zap.Error(err))
		return nil, err
	}

	//get and combine the author information according to the post.AuthorID
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserByID failed",
			zap.Int64("author_id", post.AuthorID),
			zap.Error(err))
		data.AuthorName = "This user don't exist"
	}

	//get and combine the communityDetail information according to the post.CommunityID
	communityDetail, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID failed",
			zap.Int64("community_id", post.CommunityID),
			zap.Error(err))
		data.Community = &models.Community{
			ID:   0,
			Name: "Not find",
		}
	}
	data.AuthorName = user.Username
	data.Community = &communityDetail.Community
	data.Post = post
	return data, err
}

func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		postDetail, err := GetPostByID(post.ID)
		if err != nil {
			continue
		}
		data = append(data, postDetail)
	}
	return
}

func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("GetPostIDsInOrder returns none id", zap.Error(err))
		return
	}
	//get the voting info related to each post
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		zap.L().Warn("getting voteData fails", zap.Error(err))
	}

	//combine the posts and the extra info related to the post
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for idx, post := range posts {
		postDetail, err := GetPostByID(post.ID)
		postDetail.VoteNum = voteData[idx]
		if err != nil {
			continue
		}
		data = append(data, postDetail)
	}
	return
}

func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("GetPostIDsInOrder returns none id", zap.Error(err))
		return
	}
	//get the voting info related to each post
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		zap.L().Warn("getting voteData fails", zap.Error(err))
	}

	//combine the posts and the extra info related to the post
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for idx, post := range posts {
		postDetail, err := GetPostByID(post.ID)
		postDetail.VoteNum = voteData[idx]
		if err != nil {
			continue
		}
		data = append(data, postDetail)
	}
	return
}

func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	if p.CommunityID == 0 {
		data, err = GetPostList2(p)
	} else {
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		return nil, err
	}
	return
}
