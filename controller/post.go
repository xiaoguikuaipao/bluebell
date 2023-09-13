package controller

import (
	"database/sql"
	"strconv"

	"web_app/logic"
	"web_app/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreatePostHandler 创建帖子的Handler
func CreatePostHandler(c *gin.Context) {
	//1. validate
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("create post with invalid params", zap.Error(err))
		ResponsError(c, CodeInvalidParam)
		return
	}
	// get userid from c
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponsError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	//2. create the post
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost failed", zap.Error(err))
		ResponsError(c, CodeServerBusy)
		return
	}

	//3. return the response
	ResponsSuccess(c, nil)
}

// GetPostDetailHandler this is a handler to look into a post
func GetPostDetailHandler(c *gin.Context) {
	//1. get the parameters
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post with invalid param")
		ResponsError(c, CodeInvalidParam)
		return
	}

	//2. according to the id , get the post data
	data, err := logic.GetPostByID(pid)
	if err != nil {
		if err == sql.ErrNoRows {
			zap.L().Error("logic.getPostByID failed by invalid id", zap.Error(err))
			ResponsError(c, CodeInvalidParam)
			return
		}
		zap.L().Error("logic.getPostByID failed", zap.Error(err))
		ResponsError(c, CodeServerBusy)
		return
	}

	//3. return the response
	ResponsSuccess(c, data)
}

func GetPostListHandler(c *gin.Context) {
	page, size := GetPageInfo(c)
	//1. get the data
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList failed", zap.Error(err))
		ResponsError(c, CodeServerBusy)
		return
	}
	ResponsSuccess(c, data)
}

// GetPostListHandler2 Order the post list according to the client' needs, such as by created_time, scores
// 1. get the params from frontend
// 2. search id list
// 3. on the basis of id, get the details from db
func GetPostListHandler2(c *gin.Context) {
	// Get :   /api/v1/posts2?page=1&size=10&order=time

	// 1. bind the query data to go struct (tag is `form`)

	// default struct and parameters
	p := &models.ParamPostList{
		Page:        1,
		Size:        10,
		Order:       models.OrderTime,
		CommunityID: 0,
	}
	err := c.ShouldBindQuery(p)
	if err != nil {
		zap.L().Error("GetPostListHandler2 failed", zap.Error(err))
		ResponsError(c, CodeInvalidParam)
		return
	}

	// 2. use the p to retrieve db
	data, err := logic.GetPostListNew(p)
	if err != nil {
		zap.L().Error("logic.GetPostListNew failed", zap.Error(err))
		ResponsError(c, CodeServerBusy)
		return
	}
	ResponsSuccess(c, data)
}

//func GetCommunityPostListHandler(c *gin.Context) {
//	// Get :   /api/v1/posts2?page=1&size=10&order=time&community=2
//
//	// 1. bind the query data to go struct (tag is `form`)
//
//	// default struct and parameters
//	p := &models.ParamPostList{
//		Page:        1,
//		Size:        10,
//		Order:       models.OrderTime,
//		CommunityID: 0,
//	}
//	err := c.ShouldBindQuery(p)
//	if err != nil {
//		zap.L().Error("GetCommunityPostListHandler failed", zap.Error(err))
//		ResponsError(c, CodeInvalidParam)
//		return
//	}
//
//	// 2. use the p to retrieve db
//	data, err := logic.GetCommunityPostList(p)
//	if err != nil {
//		zap.L().Error("logic.GetCommunityPostList failed", zap.Error(err))
//		ResponsError(c, CodeServerBusy)
//		return
//	}
//	ResponsSuccess(c, data)
//}
