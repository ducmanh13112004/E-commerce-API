package repositories

import (
	"database/sql"
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/models"
	"ecom_promotion_v2/internal/utils"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ControlApiRepo interface {
	GetInfoFrApiName(string) (*models.ControlApiTb, error)
	UpdateActiveFromAPIName(apiname string, active int) error
	CheckRun(apiname string) (bool, error)
	CheckRunV2(apiname string, second int) (bool, error)
	Update(input models.ControlApiTb) error
	Create(input *models.ControlApiTb) error
}
type controlApiRepo struct {
	db *gorm.DB
}

func NewControlApiRepo(
	db *gorm.DB,
) ControlApiRepo {
	return &controlApiRepo{
		db: db,
	}
}
func (repo *controlApiRepo) GetInfoFrApiName(apiname string) (*models.ControlApiTb, error) {
	result := &models.ControlApiTb{}
	ok, err := utils.CheckDBConnection(repo.db)
	if !ok {
		internal.Log.Error("Connection failed", zap.Error(err))
		return nil, err
	}
	where := fmt.Sprintf("%s = ?", result.ColumnApiName())
	err = repo.db.Debug().Table(result.TableName()).Where(where, apiname).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *controlApiRepo) UpdateActiveFromAPIName(apiname string, active int) error {
	ok, err := utils.CheckDBConnection(repo.db)
	if !ok {
		internal.Log.Error("Connection failed", zap.Error(err))
		return err
	}
	base := models.ControlApiTb{}
	where := fmt.Sprintf("%s = ?", base.ColumnApiName())
	err = repo.db.Debug().Table(base.TableName()).Where(where, apiname).UpdateColumn(base.ColumnActive(), active).Error
	return err
}

func (repo *controlApiRepo) Update(input models.ControlApiTb) error {
	ok, err := utils.CheckDBConnection(repo.db)
	if !ok {
		internal.Log.Error("Connection failed", zap.Error(err))
		return err
	}
	where := fmt.Sprintf("%s = ?", input.ColumnApiName())
	err = repo.db.Debug().Select("*").Table(input.TableName()).Where(where, input.ApiName).Updates(&input).Error
	if err != nil {
		return err
	}
	return nil
}
func (repo *controlApiRepo) CheckRun(apiname string) (bool, error) {
	err := repo.db.Debug().Transaction(func(tx *gorm.DB) error {
		result := &models.ControlApiTb{}
		where := fmt.Sprintf("%s = ?", result.ColumnApiName())
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Table(result.TableName()).Where(where, apiname).Find(&result).Error
		if err != nil {
			return err
		}
		if result.Active == 2 {
			result.Active = 1
			where := fmt.Sprintf("%s = ?", result.ColumnApiName())
			err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Table(result.TableName()).Where(where, apiname).Updates(&result).Error
			if err != nil {
				return err
			}
			return nil
		} else {
			return errors.New("False")
		}

	}, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil && err.Error() == "False" {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
func (repo *controlApiRepo) CheckRunV2(apiname string, second int) (bool, error) {
	err := repo.db.Debug().Transaction(func(tx *gorm.DB) error {
		result := &models.ControlApiTb{}
		where := fmt.Sprintf("%s = ?", result.ColumnApiName())
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Table(result.TableName()).Where(where, apiname).Find(&result).Error
		if err != nil {
			return err
		}
		now := utils.GetTimeUTC7()
		timeSchedule, err := utils.ParseTimeFrString("D/M/Y H:M:S", result.Conditions)
		if err != nil {
			return err
		}
		fmt.Println(apiname, second, now, result.Conditions, now.Sub(timeSchedule))
		if result.Active == 1 && now.Sub(timeSchedule) >= (time.Second*time.Duration(second-1)) {
			result.Conditions = utils.GetStringTime(now, "D/M/Y H:M:S")
			where := fmt.Sprintf("%s = ?", result.ColumnApiName())
			err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Table(result.TableName()).Where(where, apiname).Updates(&result).Error
			if err != nil {
				return err
			}
			return nil
		} else {
			return errors.New("False")
		}

	}, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil && err.Error() == "False" {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (repo *controlApiRepo) Create(input *models.ControlApiTb) error {
	ok, err := utils.CheckDBConnection(repo.db)
	if !ok {
		internal.Log.Error("Connection failed", zap.Error(err))
		return err
	}
	base := models.ControlApiTb{}
	err = repo.db.Debug().Table(base.TableName()).Create(&input).Error
	return err
}
