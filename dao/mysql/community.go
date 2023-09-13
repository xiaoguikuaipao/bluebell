package mysql

import (
	"database/sql"

	"web_app/models"

	"go.uber.org/zap"
)

func GetCommunityList() ([]*models.Community, error) {
	communityList := new([]*models.Community)
	sqlStr := `select community_id community_name from community`
	err := db.Select(communityList, sqlStr)
	if err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in db")
			err = nil
		}
	}
	return *communityList, err
}

func GetCommunityDetailByID(id int64) (*models.CommunityDetail, error) {
	community := new(models.CommunityDetail)
	sqlStr := `select 
				    community_id, community_name, introduction, create_time 
				from community 
				where community_id = ?
				`
	err := db.Get(community, sqlStr, id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvalidID
		}
	}
	return community, err
}
