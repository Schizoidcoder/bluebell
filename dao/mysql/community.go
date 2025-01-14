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
