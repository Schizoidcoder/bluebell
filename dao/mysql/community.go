package mysql

import (
	"bluebell/models"
	"database/sql"

	"go.uber.org/zap"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := "select community_id,community_name from community"
	if err = db.Select(&communityList, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Error("there is no community in db")
			err = nil //返回空列表
		}
		zap.L().Error("GetCommunityList failed", zap.Error(err))
		return nil, err
	}
	return
}

// GetCommunityDetailByID 根据ID查询社区详情
func GetCommunityDetailByID(id int64) (community *models.CommunityDetail, err error) {
	community = new(models.CommunityDetail)
	sqlStr := "select community_id,community_name, introduction,create_time " +
		"from community where community_id = ?"
	err = db.Get(community, sqlStr, id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvalidID
		}
	}
	return
}
