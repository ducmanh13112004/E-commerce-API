package repositories

import (
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/models"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type ProgramDirectRepo interface {
	CheckRedirectIdExists(funcName string, redirectIds []string) ([]string, error)
	CreateDirectCategory(funcName string, direct *models.ProgramDirectScreenMobileTb) error
	GetProgramDirectList(funcName string, filter *models.FilterProgramDirectScreenMobileTb) ([]models.ProgramDirectScreenMobileTb, error)
	GetProgramDirectListID(funcName string, redirectID string) ([]models.ProgramDirectScreenMobileTb, error)
	GetProgramDirectByID(funcName string, redirectID string) (*models.ProgramDirectScreenMobileTb, error)
	UpdateProgramDirect(funcName string, direct *models.ProgramDirectScreenMobileTb) error
	DeleteProgramDirect(funcName string, direct *models.ProgramDirectScreenMobileTb) error
	GetProgramRedirectListID(funcName string, redirectID int) ([]models.ProgramRedirectJoin, error)
	GetProgramRedirectList(funcName string, filter *models.FilterProgramRedirectTb) ([]models.ProgramRedirectTb, error)
	GetProgramRedirectByID(funcName string, redirectID int) (*models.ProgramRedirectTb, error)
	GetChannelsByRedirectID(funcName string, redirectID int) ([]models.ProgramRedirectChannelTb, error)
	DeleteProgramRedirect(funcName string, direct *models.ProgramRedirectTb) error
	UpdateProgramRedirect(funcName string, direct *models.ProgramRedirectTb) error
	Create(funcName string, direct *models.ProgramRedirectTb) error
	CreateChannel(channel *models.ProgramRedirectChannelTb) error
	GetProgramRedirectListName(funcName string) ([]models.RedirectInfo, error)
	ReplaceProgramRedirectChannels(funcName string, redirectId int, channels []models.InputProgramRedirectChannel, updateBy string, now time.Time) error
	GetProgramListNavigation(funcName string, code string) (*models.MasterData, error)
	CreateProgramNavigation(funcName string, record *models.MasterData) error
	UpdateProgramNavigation(funcName string, record *models.MasterData) error
}

type programdirectRepo struct {
	db *gorm.DB
}

func NewProgramDirectRepo(db *gorm.DB) ProgramDirectRepo {
	return &programdirectRepo{db: db}
}
func (r *programdirectRepo) CheckRedirectIdExists(funcName string, redirectIds []string) ([]string, error) {
	if len(redirectIds) == 0 {
		return nil, nil
	}

	var existingIds []string
	err := internal.DB(funcName).Model(&models.ProgramDirectScreenMobileTb{}).
		Select("redirect_id").
		Where("redirect_id IN (?)", redirectIds).
		Find(&existingIds).Error

	if err != nil {
		return nil, err
	}

	return existingIds, nil
}

func (r *programdirectRepo) GetProgramDirectByID(funcName, redirectID string) (*models.ProgramDirectScreenMobileTb, error) {
	direct := &models.ProgramDirectScreenMobileTb{}

	err := internal.DB(funcName).
		Table(models.ProgramDirectScreenMobileTb{}.TableName()).
		Where("redirect_id = ?", redirectID).
		Take(&direct).Error

	if err != nil {
		return nil, err
	}
	return direct, nil
}

func (r *programdirectRepo) GetProgramDirectList(funcName string, filter *models.FilterProgramDirectScreenMobileTb) ([]models.ProgramDirectScreenMobileTb, error) {
	var direct []models.ProgramDirectScreenMobileTb
	query := internal.DB(funcName).Model(&models.ProgramDirectScreenMobileTb{})

	if len(filter.FromDate) != 0 {
		query = query.Where("t_update >= ?", filter.FromDate+" 00:00:00")
	}
	if len(filter.ToDate) != 0 {
		query = query.Where("t_update <= ?", filter.ToDate+" 23:59:59")
	}

	err := query.Find(&direct).Error
	if err != nil {
		return nil, err
	}
	return direct, nil
}

func (r *programdirectRepo) CreateDirectCategory(funcName string, direct *models.ProgramDirectScreenMobileTb) error {
	return internal.DB(funcName).Create(direct).Error
}

func (r *programdirectRepo) GetProgramDirectListID(funcName, redirectID string) ([]models.ProgramDirectScreenMobileTb, error) {
	var direct []models.ProgramDirectScreenMobileTb
	err := internal.DB(funcName).Where("redirect_id = ?", redirectID).Find(&direct).Error
	if err != nil {
		return nil, err
	}
	return direct, nil
}

func (r *programdirectRepo) UpdateProgramDirect(funcName string, direct *models.ProgramDirectScreenMobileTb) error {
	return internal.DB(funcName).
		Model(direct).
		Select("*").
		Where("redirect_id = ?", direct.RedirectId).
		Updates(direct).Error
}

func (r *programdirectRepo) DeleteProgramDirect(funcName string, direct *models.ProgramDirectScreenMobileTb) error {
	return internal.DB(funcName).
		Model(&models.ProgramDirectScreenMobileTb{}).
		Where("redirect_id = ?", direct.RedirectId).
		Select("state", "update_by", "t_update").
		Updates(direct).Error
}

func (r *programdirectRepo) GetProgramRedirectByID(funcName string, redirectID int) (*models.ProgramRedirectTb, error) {
	direct := &models.ProgramRedirectTb{}
	err := internal.DB(funcName).
		Table(models.ProgramRedirectTb{}.TableName()).
		Where("redirect_id = ?", redirectID).
		Take(&direct).Error
	if err != nil {
		return nil, err
	}
	return direct, nil
}

// Lấy danh sách channel theo redirect_id
func (r *programdirectRepo) GetChannelsByRedirectID(funcName string, redirectID int) ([]models.ProgramRedirectChannelTb, error) {
	var channels []models.ProgramRedirectChannelTb
	err := internal.DB(funcName).Where("redirect_id = ?", redirectID).Find(&channels).Error
	if err != nil {
		return nil, err
	}
	return channels, nil
}

func (r *programdirectRepo) DeleteProgramRedirect(funcName string, direct *models.ProgramRedirectTb) error {
	return internal.DB(funcName).
		Model(&models.ProgramRedirectTb{}).
		Where("redirect_id = ?", direct.RedirectId).
		Select("state", "update_by", "t_update").
		Updates(direct).Error
}

func (r *programdirectRepo) GetProgramRedirectListID(funcName string, redirectID int) ([]models.ProgramRedirectJoin, error) {
	var result []models.ProgramRedirectJoin
	err := internal.DB(funcName).
		Table("program_redirect_tb pr").
		Select(`pr.redirect_id, pr.redirect_name, pr.state, pr.update_by, pr.t_create, pr.t_update,
                pc.channel, pc.min_version, pc.max_version, pc.action_type, pc.data_action, pc.data,
                pc.t_update as channel_t_update`).
		Joins("LEFT JOIN program_redirect_channel_tb pc ON pr.redirect_id = pc.redirect_id").
		Where("pr.redirect_id = ?", redirectID).
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *programdirectRepo) GetProgramRedirectList(funcName string, filter *models.FilterProgramRedirectTb) ([]models.ProgramRedirectTb, error) {
	var redirects []models.ProgramRedirectTb
	query := internal.DB(funcName).Model(&models.ProgramRedirectTb{})

	if len(filter.FromDate) != 0 {
		query = query.Where("t_update >= ?", filter.FromDate+" 00:00:00")
	}
	if len(filter.ToDate) != 0 {
		query = query.Where("t_update <= ?", filter.ToDate+" 23:59:59")
	}

	if err := query.Find(&redirects).Error; err != nil {
		return nil, err
	}
	return redirects, nil
}

// Tạo ProgramRedirect
func (r *programdirectRepo) Create(funcName string, direct *models.ProgramRedirectTb) error {
	return internal.DB(funcName).Create(direct).Error
}

// Tạo ProgramRedirectChannel
func (r *programdirectRepo) CreateChannel(channel *models.ProgramRedirectChannelTb) error {
	return internal.Db.Create(channel).Error
}

func (r *programdirectRepo) GetProgramRedirectListName(funcName string) ([]models.RedirectInfo, error) {
	var redirects []models.RedirectInfo
	query := internal.DB(funcName).Model(&models.ProgramRedirectTb{}).
		Select("redirect_id, redirect_name").
		Where("state = ?", 1)

	if err := query.Find(&redirects).Error; err != nil {
		return nil, err
	}
	return redirects, nil
}

// Update bảng chính
func (r *programdirectRepo) UpdateProgramRedirect(funcName string, direct *models.ProgramRedirectTb) error {
	return internal.DB(funcName).Model(&models.ProgramRedirectTb{}).
		Where("redirect_id = ?", direct.RedirectId).
		Updates(direct).Error
}

// Xoá toàn bộ channel cũ rồi insert lại
func (r *programdirectRepo) ReplaceProgramRedirectChannels(funcName string, redirectId int, channels []models.InputProgramRedirectChannel, updateBy string, now time.Time) error {
	tx := internal.DB(funcName).Begin()

	// Xoá cũ
	if err := tx.Where("redirect_id = ?", redirectId).Delete(&models.ProgramRedirectChannelTb{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Insert mới
	for _, ch := range channels {
		dataJson, _ := json.Marshal(ch.Data)
		str := string(dataJson)
		channelRecord := models.ProgramRedirectChannelTb{
			RedirectId: redirectId,
			Channel:    ch.Channel,
			ActionType: ch.ActionType,
			DataAction: ch.DataAction,
			Data:       &str,
			TUpdate:    &now,
		}

		if err := tx.Create(&channelRecord).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *programdirectRepo) GetProgramListNavigation(funcName string, code string) (*models.MasterData, error) {
	var data models.MasterData
	query := internal.DB(funcName).Model(&models.MasterData{}).
		Where("code = ? AND active = 1", code).
		First(&data)

	if query.Error != nil {
		return nil, query.Error
	}
	return &data, nil
}

func (r *programdirectRepo) CreateProgramNavigation(funcName string, record *models.MasterData) error {
	return internal.DB(funcName).Create(record).Error
}

func (r *programdirectRepo) UpdateProgramNavigation(funcName string, record *models.MasterData) error {
	return internal.DB(funcName).Save(record).Error
}
