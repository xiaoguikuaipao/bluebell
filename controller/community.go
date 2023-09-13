package controller

import (
	"strconv"

	"web_app/logic"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CommunityHandler 列出社区列表的Handler
func CommunityHandler(c *gin.Context) {
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList", zap.Error(err))
		ResponsError(c, CodeServerBusy) //you'd better not expose the error to the Frontend
		return
	}
	ResponsSuccess(c, data)
}

// CommunityDetailHandler 展示社区详细信息的Handler
func CommunityDetailHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponsError(c, CodeInvalidParam)
		return
	}
	data, err := logic.GetCommunityDetailByID(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityList", zap.Error(err))
		ResponsError(c, CodeServerBusy) //you'd better not expose the error to the Frontend
		return
	}
	ResponsSuccess(c, data)
}
