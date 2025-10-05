package repositories

import (
	"ecom_promotion_v2/internal/models"
	// "errors"
	"gorm.io/gorm"
	"time"
)

type CategoryRepo interface {
	CreateCategory(*models.CategoryTb) error
	GetCategoryByID(id int) (*models.CategoryTb, error)
	UpdateCategory(category *models.CategoryTb) error
	DeleteCategory(category *models.CategoryTb) error
	GetCategoryList(filter *models.FilterCategory) ([]models.CategoryTb, error)
	GetCategoryListName() ([]models.CategoryTb, error)
}
type categoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepo(db *gorm.DB) CategoryRepo {
	return &categoryRepo{
		db: db,
	}
}

func (r *categoryRepo) CreateCategory(category *models.CategoryTb) error {
	// Kiểm tra trùng tên category
	// var count int64
	// r.db.Table(models.CategoryTb{}.TableName()).Where("category_name = ?", category.CategoryName).Count(&count)
	// if count > 0 {
	// 	return errors.New("category_name already exists")
	// }
	now := time.Now()
	category.TUpdate = &now
	return r.db.Table(models.CategoryTb{}.TableName()).Create(category).Error
}

func (r *categoryRepo) UpdateCategory(category *models.CategoryTb) error {
	return r.db.Model(category).Select("*").Where("category_id = ?", category.CategoryId).Updates(category).Error

}

func (r *categoryRepo) GetCategoryByID(id int) (*models.CategoryTb, error) {
	category := &models.CategoryTb{}

	err := r.db.Table(models.CategoryTb{}.TableName()).Where("category_id = ?", id).Take(&category).Error
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (r *categoryRepo) DeleteCategory(category *models.CategoryTb) error {

	return r.db.Model(&models.CategoryTb{}).
		Where("category_id = ?", category.CategoryId).
		Select("state", "update_by", "t_update").
		Updates(category).Error

}

func (r *categoryRepo) GetCategoryList(filter *models.FilterCategory) ([]models.CategoryTb, error) {
	query := r.db
	if (filter.CategoryId) != 0 {
		query = query.Where("category_id= ? ", filter.CategoryId)
	}
	if filter.State != nil {
		query = query.Where("state= ? ", filter.State)
	}
	var category []models.CategoryTb
	err := query.Find(&category).Error
	if err != nil {
		return nil, err
	}
	return category, nil
}
func (r *categoryRepo) GetCategoryListName() ([]models.CategoryTb, error) {
	var category []models.CategoryTb
	err := r.db.Select("category_id, category_name").Find(&category).Error
	if err != nil {
		return nil, err
	}
	return category, nil
}
