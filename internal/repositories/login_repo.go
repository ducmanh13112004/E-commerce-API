package repositories

import (
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type LoginRepo interface {
	GetByUserName(funcName, username string) (*models.Login, error)
}

type loginrepo struct {
	db *gorm.DB
}

func NewLoginRepo(db *gorm.DB) LoginRepo {
	return &loginrepo{
		db: db,
	}
}
func (r *loginrepo) GetByUserName(funcName, username string) (*models.Login, error) {
	var user models.Login

	err := internal.DB(funcName).Where("username = ?", username).First(&user).Error
	if err != nil {
		internal.Log.Error("GetByUserName error", zap.String("funcName", funcName), zap.String("username", username), zap.Error(err))
		return nil, err
	}

	return &user, nil
}
