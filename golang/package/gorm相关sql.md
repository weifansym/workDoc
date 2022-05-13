## 常用gorm相关sql
### 执行原生sql查询，下面两个查询条件是等价的。
```
func (f FollowerRepository) GetFollowerForPush(ctx context.Context, liveUid uint) (*entity.FollowQueryForPush, error) {
	queryStr := "SELECT MIN(id) as minId,COUNT(*) as totalCount FROM `tbl_bin_follower` " +
		"WHERE user_id =? and status=? and push_switch=?;"
	sqlParams := []interface{}{liveUid, entity.FollowerStatusNormal, entity.FollowerPushSwitchOpen}
	var pushData *entity.FollowQueryForPush
	db := f.GetDB(ctx)
	db = db.Raw(queryStr, sqlParams...)
	err := db.Scan(&pushData).Error
	if err != nil {
		return nil, error_code.DBError.WithDetails(err.Error())
	}
	return pushData, nil
}

func (f FollowerRepository) GetFollowerForPush(ctx context.Context, liveUid uint) (*entity.FollowQueryForPush, error) {
	var pushData *entity.FollowQueryForPush
	result := f.GetDB(ctx).
		Model(&entity.Follower{}).
		Select("MIN(id) as minId, COUNT(*) as totalCount").
		Where("user_id = ? and status=? and push_switch=? ", liveUid, entity.FollowerStatusNormal, entity.FollowerPushSwitchOpen).
		Scan(&pushData)
	if result.Error != nil {
		return nil, error_code.DBError.WithDetails(result.Error.Error())
	}
	return pushData, nil
}

```
