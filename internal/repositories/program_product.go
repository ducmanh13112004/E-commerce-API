package repositories

import (
	"ecom_promotion_v2/internal/models"

	"gorm.io/gorm"
)

type ProgramProductRepo interface {
	CreateListProgramProduct(listProgramProduct []models.ProgramProductTb) error
}

type programProductRepo struct {
	db *gorm.DB
}

func NewProgramProductRepo(db *gorm.DB) ProgramProductRepo {
	return &programProductRepo{db: db}
}

func (p *programProductRepo) CreateListProgramProduct(listProgramProduct []models.ProgramProductTb) error {
	return p.db.CreateInBatches(listProgramProduct, 10000).Error
}
