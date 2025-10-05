package repositories

import (
	"context"
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/models"

	"gorm.io/gorm"
)

type ProgramPromotionRepo interface {
	GetProgramPromotionWithFilter(filter *models.FilterProgramPromotion) ([]models.ProgramPromotionTb, error)
	GetProgramPromotionByID(id int64, funcName string) (*models.ProgramPromotionTb, error)
	CreateProgramPromotion(program *models.ProgramPromotionTb, listChannel []string) (int64, error)
	UpdateProgramPromotionChannel(program *models.ProgramPromotionTb, listChannel []string) error
	UpdateProgramPromotion(program *models.ProgramPromotionTb) error
	GetReportProgramPromotion(filter *models.FilterReportProgramPromotion) ([]models.ReportProgramPromotion, error)
	CheckProgramPromotionID(listIdProgram []int64) ([]string, error)
	CheckProgramIsExist(programId int64, funcName string) (bool, error)
}

type programPromotionsRepo struct {
	db *gorm.DB
}

func NewProgramPromotionsRepo(db *gorm.DB) ProgramPromotionRepo {
	return &programPromotionsRepo{db: db}
}

func (p *programPromotionsRepo) GetProgramPromotionWithFilter(filter *models.FilterProgramPromotion) ([]models.ProgramPromotionTb, error) {
	var result []models.ProgramPromotionTb
	query := p.db.Preload("ProgramPromotionChannels")
	if filter.ProgramId != nil && *filter.ProgramId != 0 {
		err := query.Where("is_deleted = ? AND program_id = ?", 0, filter.ProgramId).Find(&result).Error
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	if len(filter.MerchantID) != 0 {
		query = query.Where("merchant_id = ? ", filter.MerchantID)
	}
	if len(filter.FromDate) != 0 {
		query = query.Where("t_create >= ?", filter.FromDate+" 00:00:00")
	}
	if len(filter.ToDate) != 0 {
		query = query.Where("t_create <= ?", filter.ToDate+" 23:59:59")
	}
	if filter.State != nil {
		query = query.Where("state = ?", filter.State)
	}

	err := query.Where("is_deleted = ?", 0).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *programPromotionsRepo) GetProgramPromotionByID(id int64, funcName string) (*models.ProgramPromotionTb, error) {
	ctx := context.WithValue(context.Background(), internal.FuncNameKey, funcName)
	var result = &models.ProgramPromotionTb{}
	err := p.db.Debug().WithContext(ctx).Preload("ProgramPromotionChannels").Where("program_id = ? and is_deleted = ?", id, 0).Take(result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *programPromotionsRepo) CreateProgramPromotion(program *models.ProgramPromotionTb, listChannel []string) (int64, error) {
	tx := p.db.Begin()
	err := tx.Create(&program).Error
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	var programChannels []models.ProgramPromotionChannel
	for _, channel := range listChannel {
		programChannel := models.ProgramPromotionChannel{
			ProgramID: program.ProgramId,
			Channel:   channel,
		}
		programChannels = append(programChannels, programChannel)
	}
	err = tx.CreateInBatches(programChannels, 100).Error
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()
	return program.ProgramId, nil
}

func (p *programPromotionsRepo) UpdateProgramPromotionChannel(program *models.ProgramPromotionTb, listChannel []string) error {
	tx := p.db.Begin()
	// Bước 1: update program
	err := tx.Model(program).Select("*").Where("program_id = ?", program.ProgramId).Updates(program).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Bước 2: delete program promotion channel theo program_id
	err = tx.Model(&models.ProgramPromotionChannel{}).Where("program_id = ?", program.ProgramId).Delete(&models.ProgramPromotionChannel{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Bước 3: thêm lại program promotion channel mới cập nhật
	var programChannels []models.ProgramPromotionChannel
	for _, channel := range listChannel {
		programChannel := models.ProgramPromotionChannel{
			ProgramID: program.ProgramId,
			Channel:   channel,
		}
		programChannels = append(programChannels, programChannel)
	}
	err = tx.CreateInBatches(programChannels, 100).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (p *programPromotionsRepo) UpdateProgramPromotion(program *models.ProgramPromotionTb) error {
	return p.db.Model(program).Select("*").Where("program_id = ?", program.ProgramId).Updates(program).Error
}

// SELECT
//
//	pp.program_id,
//	pp.program_name,
//	(select sum(total_release) from promotion_tb p1 where p1.program_id = pp.program_id) AS total_release,
//	COUNT(h.id) AS total_collect,
//	COUNT(
//	  CASE WHEN h.state_usable = 1 THEN h.id END
//	) AS total_use
//
// FROM
//
//	(SELECT * FROM program_promotion_tb pp WHERE
//	pp.t_create >= '2020-01-01 00:00:00'
//	AND pp.t_create <= '2025-02-01 23:59:59' ) AS pp
//	LEFT JOIN promotion_tb p ON pp.program_id = p.program_id
//	LEFT JOIN histories_usable_tb h ON p.promotion_code = h.promotion_code
//
// GROUP BY
//
//	pp.program_id
func (p *programPromotionsRepo) GetReportProgramPromotion(filter *models.FilterReportProgramPromotion) ([]models.ReportProgramPromotion, error) {
	var result []models.ReportProgramPromotion
	query := p.db.Table("program_promotion_tb AS pp").
		Select("pp.program_id, pp.program_name, " +
			"(SELECT SUM(total_release) from promotion_tb p1 where p1.program_id = pp.program_id) AS total_release, " +
			"COUNT( h.id) AS total_collect, " +
			"COUNT( CASE WHEN h.state_usable = 1 THEN h.id END) AS total_use").
		Joins("LEFT JOIN promotion_tb p ON pp.program_id = p.program_id").
		Joins("LEFT JOIN histories_usable_tb h ON p.promotion_code = h.promotion_code")

	if len(filter.ListProgramId) > 0 {
		query = query.Where("pp.program_id IN ?", filter.ListProgramId)
	} else {
		query = query.Where("pp.t_create >= ? AND pp.t_create <= ?", filter.FromDate+" 00:00:00", filter.ToDate+" 23:59:59")
	}
	err := query.Where("is_deleted = ?", 0).Group("pp.program_id").Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (p *programPromotionsRepo) CheckProgramPromotionID(listIdProgram []int64) ([]string, error) {
	var result []string
	err := p.db.Select("program_id").Model(&models.ProgramPromotionTb{}).Where("program_id IN ? AND is_deleted = ?", listIdProgram, 0).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *programPromotionsRepo) CheckProgramIsExist(programId int64, funcName string) (bool, error) {
	ctx := context.WithValue(context.Background(), internal.FuncNameKey, funcName)
	var count int64
	err := p.db.WithContext(ctx).Select("count(1)").Model(&models.ProgramPromotionTb{}).Where("program_id = ? AND is_deleted = ?", programId, 0).Find(&count).Error
	if err != nil || count == 0 {
		return false, err
	}
	return true, nil
}
