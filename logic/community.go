package logic

import (
	"web_app/dao/mysql"
	"web_app/models"
)

// GetCommunityList 在mysql中获得社区分类列表
func GetCommunityList() ([]*models.Community, error) {
	//get the data from db
	return mysql.GetCommunityList()
}

// GetCommunityDetailByID 根据社区ID获得社区具体信息
func GetCommunityDetailByID(id int64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetailByID(id)
}
