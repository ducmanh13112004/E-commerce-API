package repositories

import (
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProgramCategoryRepo interface {
	GetProgramCategoryList(filter *models.FilterProgramcategory) ([]models.ProgramCategory, error)
	GetProgramCategory() ([]models.CategoryTb, error)
	CreateProgramCategories([]models.ProgramCategory) error
	GetProgramCategoryByID(id int) (*models.ProgramCategory, error)
	UpdateProgramCategories(categories []models.ProgramCategory) error
	DeleteProgramCategory(category models.ProgramCategory) error
	IsProgramIDExists(programID int) (bool, error)
	IsCategoryIDExists(categoryID int) (bool, error)
	IsProgramidCategoryidExists(programID int64, categoryID int) (bool, error)
	UpdateOnOffProgramCategory(category *models.ProgramCategory) error
}

type programcategoryRepo struct {
	db *gorm.DB
}

func NewProgramCategoryRepo(db *gorm.DB) ProgramCategoryRepo {
	return &programcategoryRepo{
		db: db,
	}
}

func (r *programcategoryRepo) IsProgramIDExists(programID int) (bool, error) {
	var count int64
	err := r.db.Model(&models.ProgramPromotionTb{}).
		Where("program_id = ?", programID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (r *programcategoryRepo) IsCategoryIDExists(categoryID int) (bool, error) {
	var count int64
	err := r.db.Model(&models.CategoryTb{}).
		Where("category_id = ?", categoryID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (r *programcategoryRepo) GetProgramCategoryByID(id int) (*models.ProgramCategory, error) {
	category := &models.ProgramCategory{}

	err := r.db.Table(models.ProgramCategory{}.TableName()).Where("id = ?", id).Take(&category).Error
	if err != nil {
		return nil, err
	}
	return category, nil
}

//	func (r *programcategoryRepo) GetProgramCategoryList() ([]models.ProgramCategory, error) {
//		var category []models.ProgramCategory
//		err := r.db.Preload("CategoryTb").Find(&category).Error
//		if err != nil {
//			return nil, err
//		}
//		return category, nil
//	}
func (r *programcategoryRepo) GetProgramCategoryList(filter *models.FilterProgramcategory) ([]models.ProgramCategory, error) {
	query := r.db.Preload("CategoryTb").Preload("ProgramPromotionTb")
	if len(filter.FromDate) != 0 {
		query = query.Where("t_update >= ?", filter.FromDate+" 00:00:00")
	}
	if len(filter.ToDate) != 0 {
		query = query.Where("t_update <= ?", filter.ToDate+" 23:59:59")
	}
	if filter.ProgramId != nil && *filter.ProgramId != 0 {
		query = query.Where("program_id = ?", *filter.ProgramId)
	}
	if filter.CategoryId != nil && *filter.CategoryId != 0 {
		query = query.Where("category_id = ?", *filter.CategoryId)
	}

	var categories []models.ProgramCategory
	err := query.Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *programcategoryRepo) GetProgramCategory() ([]models.CategoryTb, error) {
	var categories []models.CategoryTb
	err := r.db.
		Where("state = ?", 1).
		Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

//	func (r *programcategoryRepo) CreateProgramCategory(category *models.ProgramCategory) (int, error) {
//		err := r.db.Create(&category).Error
//		if err != nil {
//			return 0, err
//		}
//		return category.Id, nil
//	}
func (r *programcategoryRepo) CreateProgramCategories(listcategories []models.ProgramCategory) error {
	return r.db.CreateInBatches(listcategories, 10000).Error
}

func (r *programcategoryRepo) IsProgramidCategoryidExists(programID int64, categoryID int) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM program_category WHERE program_id = ? AND category_id = ?`

	err := r.db.Raw(query, programID, categoryID).Scan(&count).Error
	if err != nil {
		internal.Log.Error("Error checking if program_id and category_id exist",
			zap.String("funcname", "IsProgramidCategoryidExists"),
			zap.Int64("program_id", programID),
			zap.Int("category_id", categoryID),
			zap.Error(err),
		)
		return false, err
	}

	return count > 0, nil
}

func (r *programcategoryRepo) UpdateProgramCategories(categories []models.ProgramCategory) error {
	tx := r.db.Begin()
	for _, category := range categories {
		err := tx.Model(&models.ProgramCategory{}).Select("*").Where("id = ?", category.Id).Updates(category).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (p *programcategoryRepo) UpdateOnOffProgramCategory(category *models.ProgramCategory) error {
	return p.db.Model(category).Select("*").Where("id = ?", category.Id).Updates(category).Error
}

//	func (r *programcategoryRepo) DeleteProgramCategory(category models.ProgramCategory) error {
//		return r.db.Model(&models.ProgramCategory{}).Where("id = ?", category.Id).Select("active", "update_by", "t_update").Updates(category).Error
//	}
func (r *programcategoryRepo) DeleteProgramCategory(category models.ProgramCategory) error {
	return r.db.Where("id = ?", category.Id).Delete(&models.ProgramCategory{}).Error
}
