package repositories

import (
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/models"
	"ecom_promotion_v2/internal/utils"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type StoreCallApiRepo interface {
	Create(input *models.StoreCallApiTb) error
}
type storeCallApiRepo struct {
	db *gorm.DB
}

func NewStoreCallApiRepo(
	db *gorm.DB,
) StoreCallApiRepo {
	return &storeCallApiRepo{
		db: db,
	}
}

func (repo *storeCallApiRepo) Create(input *models.StoreCallApiTb) error {
	ok, err := utils.CheckDBConnection(repo.db)
	if !ok {
		internal.Log.Error("Connection failed", zap.Error(err))
		return nil
	}
	input.TCreate = utils.GetTimeUTC7()
	err = repo.db.Table(input.TableName()).Create(&input).Error
	if err != nil {
		return err
	}
	return nil
}
