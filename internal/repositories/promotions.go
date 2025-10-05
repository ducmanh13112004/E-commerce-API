package repositories

import (
	"context"
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

type PromotionsRepo interface {
	SPReceiveVoucher(*models.InputReceiveVoucher) (*models.RespSPReceiveVoucher, error)
	GetProgramIdFrCategoryId(categoryId int) (*models.CategoryTb, error)
	UserGetCustomerReceivedCodeTimes(string, int64) ([]models.HistoriesUsableTb, error)
	PaywifiGetListProvince(merchantId string, activate int) ([]models.MerchantLocationTb, error)
	AfiliateGetListLocationFrProgramId(programId int64, activate int) ([]models.MerchantLocationTb, error)
	AfiliateGetRemainPromotionCountFrProgramId(programId int64) (*models.TotalRemainPromotionProgramId, error)
	AfiliateLookupCustomerHasVoucherFrProgramId(customer *models.UserInfo, programId int64) (*models.InfoPromotionCodeCus, error)
	AfiliateLookupCustomerHasVoucherFrVoucherCode(customer *models.UserInfo, voucherCode string) (*models.InfoPromotionCodeCus, error)

	//Portal web
	CreateListPromotion(listPromotions []models.PromotionTb) error
	CreatePromotion(promotion *models.PromotionTb) error
	CheckCodeIsExist(listPromotion []string) ([]string, error)
	CheckPromotionCodeIsExist(code string) (int, error)
	GetPromotionByProgramId(programId int64, funcName string) ([]models.PromotionTb, error)

	// Điều chỉnh độ ưu tiên chương trình khuyến mãi
	GetPromotionsPriority([]string, int, int, string) ([]models.ProgramPromotionPriority, error)
	CheckInsertPriority(int, string) ([]models.ProgramPromotionPriority, error)
	CheckUpdatePriority(idInput int) (*models.ProgramPromotionPriority, error)
	CreatePriority(*models.ProgramPromotionPriority) (*models.ProgramPromotionPriority, error)
	UpdatePriority(*models.ProgramPromotionPriority) error
}
type promotionsRepo struct {
	db *gorm.DB
}

func NewPromotionsRepo(
	db *gorm.DB,
) PromotionsRepo {
	return &promotionsRepo{
		db: db,
	}
}
func (m *promotionsRepo) SPReceiveVoucher(input *models.InputReceiveVoucher) (*models.RespSPReceiveVoucher, error) {
	defer utils.ExecTime(utils.GetTimeUTC7(), "SPReceiveVoucher", nil)
	inputSp := &models.InputSPReceiceVoucher{
		ProgramId:      input.ProgramId,
		CustomerId:     input.CustomerId,
		CustomerPhone:  input.CustomerPhone,
		ReceiveCurrent: 1, // số lượng lấy
		StatePromotion: 1, // 1 còn voucher, 0 hết voucher
		StateUsable:    0, // 0 chưa dùng, 1, đã dùng, 2 đang chờ
	}
	voucher := &models.PromotionTb{}
	availableRemainTimes := 0
	ctxTimeout, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	transErr := m.db.Debug().WithContext(ctxTimeout).Transaction(func(tx *gorm.DB) error {
		// Kiểm tra chương trình còn khả dụng hay không - Nếu có thì trả về số lượt sd mỗi KH ?
		program, err := m.CheckProgramActive(ctxTimeout, tx, inputSp)
		if err != nil {
			return err
		}
		internal.Log.Info("1St.CheckProgramActive", zap.Any("input", input), zap.Any("program", program))
		inputSp.ProgramTb = program
		// Tính số lần KH đã lấy mã từ chương trình này, so lượt lấy voucher tối đa của khách hàng
		historiesUsable, err := m.GetCustomerReceivedCodeTimes(ctxTimeout, tx, inputSp, program.TimesEachCustomer)
		if err != nil {
			return err
		}
		internal.Log.Info("2St.GetCustomerReceivedCodeTimes", zap.Any("input", input), zap.Any("len", len(historiesUsable)))
		// Cập nhật -1 số lượng vật phẩm plays còn lại của user
		voucher, err = m.AddReceiveVoucherCustomer(ctxTimeout, tx, inputSp)
		if err != nil {
			return err
		}
		internal.Log.Info("3St.AddReceiveVoucherCustomer", zap.Any("input", input), zap.Any("voucher", voucher))
		// Tính số lần KH đã lấy mã từ chương trình này, so lượt lấy voucher tối đa của khách hàng - cả voucher vừa lấy
		historiesUsable2, err := m.GetCustomerReceivedCodeTimes(ctxTimeout, tx, inputSp, program.TimesEachCustomer)
		if err != nil {
			return err
		}
		internal.Log.Info("4St.GetCustomerReceivedCodeTimes", zap.Any("input", input), zap.Any("len", len(historiesUsable2)))
		// số lượng còn lại có thể lấy
		availableRemainTimes = program.TimesEachCustomer - len(historiesUsable2)
		return nil
	})
	if transErr != nil {
		// Check if the error is a timeout error
		if ctxTimeout.Err() == context.DeadlineExceeded {
			internal.Log.Error("SPReceiveVoucher-TimeOut", zap.Any("input", input), zap.Error(ctxTimeout.Err()))
		} else {
			internal.Log.Error("SPReceiveVoucher-Error", zap.Any("input", input), zap.Error(transErr))
		}
		return nil, transErr
	}
	respSp := &models.RespSPReceiveVoucher{
		ProgramInfo:          inputSp.ProgramTb,
		VoucherInfo:          voucher,
		AvailableRemainTimes: availableRemainTimes,
	}

	return respSp, transErr
}

func (m *promotionsRepo) CheckProgramActive(ctxTimeout context.Context, tx *gorm.DB, input *models.InputSPReceiceVoucher) (result *models.ProgramPromotionTb, err error) {
	defer utils.ExecTime(utils.GetTimeUTC7(), "CheckProgramActive", nil)
	now := utils.GetTimeUTC7()
	strDate := utils.GetStringTime(now, "Y-M-D")

	select {
	case <-ctxTimeout.Done():
		return nil, ctxTimeout.Err()
	default:
		result = &models.ProgramPromotionTb{}
		where := fmt.Sprintf("%s = ? AND %s >= ?", result.ColumnProgramId(), result.ColumnEndUsable())
		err = tx.Debug().Where(where, input.ProgramId, strDate).Preload("RedirectInfo").First(result).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			internal.Log.Error("CheckProgramActive", zap.Any("input", input), zap.Error(err))
			return nil, fmt.Errorf("Chương trình không khả dụng!")
		}
		if err != nil {
			internal.Log.Error("CheckProgramActive", zap.Any("input", input), zap.Error(err))
			return nil, fmt.Errorf("%s", internal.SysStatus.SystemError.Msg)
		}
	}
	return result, err
}

func (m *promotionsRepo) GetCustomerReceivedCodeTimes(ctxTimeout context.Context, tx *gorm.DB, input *models.InputSPReceiceVoucher, totalLimit int) (histories []models.HistoriesUsableTb, err error) {
	defer utils.ExecTime(utils.GetTimeUTC7(), "GetCustomerReceivedCodeTimes", nil)
	select {
	case <-ctxTimeout.Done():
		return nil, ctxTimeout.Err()
	default:
		histories = []models.HistoriesUsableTb{}
		basePro := &models.PromotionTb{}
		baseHis := &models.HistoriesUsableTb{}

		where := fmt.Sprintf("a.%s = ? AND b.%s = ?", baseHis.ColumnCustomerId(), basePro.ColumnProgramId())
		join := fmt.Sprintf(" INNER JOIN %s b ON a.%s = b.%s", basePro.TableName(), baseHis.ColumnPromotionCode(), baseHis.ColumnPromotionCode())
		err = tx.Debug().Select("a.*").Table(baseHis.TableName()+" a").Joins(join).Where(where, input.CustomerId, input.ProgramId).Find(&histories).Error
		if err != nil {
			internal.Log.Error("GetCustomerReceivedCodeTimes", zap.Any("input", input), zap.Error(err))
			return nil, fmt.Errorf("%s", internal.SysStatus.SystemError.Msg)
		}
		if len(histories) >= totalLimit {
			return nil, fmt.Errorf("Bạn đã nhận đủ số lượng mã ưu đãi cho chương trình này rồi!")
		}
	}
	return histories, err
}

func (m *promotionsRepo) AddReceiveVoucherCustomer(ctxTimeout context.Context, tx *gorm.DB, input *models.InputSPReceiceVoucher) (voucher *models.PromotionTb, err error) {
	defer utils.ExecTime(utils.GetTimeUTC7(), "AddReceiveVoucherCustomer", nil)
	select {
	case <-ctxTimeout.Done():
		return nil, ctxTimeout.Err()
	default:
		now := utils.GetTimeUTC7()
		strDate := utils.GetStringTime(now, "Y-M-D")
		// mã lấy được
		voucher := &models.PromotionTb{}
		// lấy mã hợp lệ
		promotion := &models.PromotionTb{}
		program := &models.ProgramPromotionTb{}
		join := fmt.Sprintf("INNER JOIN %s b ON a.%s = b.%s", program.TableName(), program.ColumnProgramId(), program.ColumnProgramId())
		subwhere := fmt.Sprintf("a.%s = ? AND a.%s = ? AND b.%s >= ? AND a.%s >= ? AND a.%s >= ?",
			voucher.ColumnProgramId(), voucher.ColumnStatePromotion(), program.ColumnEndUsable(), voucher.ColumnPromotionEndUsable(), voucher.ColumnTotalCurrent())
		err = tx.Debug().Table(voucher.TableName()+" a").Select("a.*").Joins(join).
			Where(subwhere, input.ProgramId, input.StatePromotion, strDate, strDate, input.ReceiveCurrent).Limit(1).
			Clauses(clause.Locking{Strength: "UPDATE"}).Scan(voucher).Error
		if err != nil {
			internal.Log.Error("1.AddReceiveVoucherCustomer", zap.Any("input", input), zap.Error(err))
			return nil, fmt.Errorf("%s", internal.SysStatus.DbFailed.Msg)
		}
		internal.Log.Info("1.1.AddReceiveVoucherCustomer", zap.Any("input", input), zap.Any("voucher", voucher))
		if voucher.PromotionCode == "" {
			return nil, fmt.Errorf("Không còn mã khuyến mãi nào khả dụng!")
		}
		statePromotionUpdate := input.StatePromotion
		if promotion.TotalCurrent == 0 {
			statePromotionUpdate = 0
		}
		// tạm ghi mã hợp lệ
		err = tx.Debug().Table(promotion.TableName()).
			Where(fmt.Sprintf("%s = (?)", promotion.ColumnPromotionCode()), voucher.PromotionCode).
			Updates(map[string]interface{}{
				promotion.ColumnTotalCurrent():   gorm.Expr(fmt.Sprintf("%s - ?", promotion.ColumnTotalCurrent()), input.ReceiveCurrent),
				promotion.ColumnTotalUse():       gorm.Expr(fmt.Sprintf("%s + ?", promotion.ColumnTotalUse()), input.ReceiveCurrent),
				promotion.ColumnStatePromotion(): statePromotionUpdate,
				// promotion.ColumnStatePromotion(): gorm.Expr(
				// 	"CASE WHEN " + promotion.ColumnTotalCurrent() + " > " + strconv.Itoa(input.ReceiveCurrent) + " THEN 1 ELSE 0 END",
				// ),
			}).Error
		if err != nil {
			internal.Log.Error("2.AddReceiveVoucherCustomer", zap.Any("input", input), zap.Error(err))
			return nil, fmt.Errorf("%s", internal.SysStatus.DbFailed.Msg)
		}
		internal.Log.Info("2.1.AddReceiveVoucherCustomer", zap.Any("input", input), zap.Any("promotion", promotion))
		usableHis := models.HistoriesUsableTb{
			CustomerID:    &input.CustomerId,
			PromotionCode: voucher.PromotionCode,
			StateUsable:   input.StateUsable,
			TCreate:       &now,
		}
		// Lưu lịch sử trừ quà
		err = tx.Debug().Model(usableHis).Create(&usableHis).Error
		if err != nil {
			internal.Log.Info("3.AddReceiveVoucherCustomer", zap.Any("input", usableHis), zap.Error(err))
			return nil, fmt.Errorf("%s", internal.SysStatus.DbFailed.Msg)
		}
		internal.Log.Info("3.1.AddReceiveVoucherCustomer", zap.Any("input", usableHis), zap.Any("usableHis", usableHis))
		return voucher, nil
	}
}

func (m *promotionsRepo) GetProgramIdFrCategoryId(categoryId int) (*models.CategoryTb, error) {
	defer utils.ExecTime(utils.GetTimeUTC7(), "GetProgramIdFrCategoryId", nil)
	now := utils.GetTimeUTC7()
	result := &models.CategoryTb{}
	baseProg := &models.ProgramPromotionTb{}
	// Add(-7 * time.Hour): múi giờ từ getime ra là +7, khu
	timeFormat := now.Truncate(24 * time.Hour).Add(-7 * time.Hour)
	where := fmt.Sprintf("%s = ? and %s = ?", result.ColumnCategoryId(), result.ColumnState())
	err := m.db.Debug().Where(where, categoryId, 1).
		Preload("ListProgram.ProgPromotionTb", func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s= ? and %s >= ?", baseProg.ColumnState(), baseProg.ColumnEndUsable()), 1, timeFormat).Order("`" + baseProg.ColumnRank() + "`" + " ASC").
				Preload("RedirectInfo").
				Preload("MerchantInfo")
		}).Order(result.ColumnCategoryPriority() + " ASC").First(result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return result, err
}

func (m *promotionsRepo) UserGetCustomerReceivedCodeTimes(customerId string, programId int64) ([]models.HistoriesUsableTb, error) {
	defer utils.ExecTime(utils.GetTimeUTC7(), "UserGetCustomerReceivedCodeTimes", nil)

	histories := []models.HistoriesUsableTb{}
	basePro := &models.PromotionTb{}
	baseHis := &models.HistoriesUsableTb{}

	where := fmt.Sprintf("a.%s = ? AND b.%s = ?", baseHis.ColumnCustomerId(), basePro.ColumnProgramId())
	join := fmt.Sprintf(" INNER JOIN %s b ON a.%s = b.%s", basePro.TableName(), baseHis.ColumnPromotionCode(), baseHis.ColumnPromotionCode())
	if err := m.db.Debug().Select("a.*").Table(baseHis.TableName()+" a").Joins(join).Where(where, customerId, programId).Find(&histories).Error; err != nil {
		internal.Log.Error("UserGetCustomerReceivedCodeTimes", zap.Any("customerId", customerId), zap.Any("programId", programId), zap.Error(err))
		return nil, fmt.Errorf("%s", internal.SysStatus.SystemError.Msg)
	}
	return histories, nil
}

func (m *promotionsRepo) PaywifiGetListProvince(merchantId string, activate int) ([]models.MerchantLocationTb, error) {
	defer utils.ExecTime(utils.GetTimeUTC7(), "PaywifiGetListProvince", nil)
	result := []models.MerchantLocationTb{}
	base := &models.MerchantLocationTb{}
	where := fmt.Sprintf("%s = ? and %s = ?", base.ColumnMerchantId(), base.ColumnActivate())
	err := m.db.Debug().Where(where, merchantId, activate).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, err
}

func (m *promotionsRepo) AfiliateGetListLocationFrProgramId(programId int64, activate int) ([]models.MerchantLocationTb, error) {
	defer utils.ExecTime(utils.GetTimeUTC7(), "AfiliateGetListLocationFrProgramId", nil)
	result := []models.MerchantLocationTb{}
	baseA := &models.ProgramLocationTb{}
	baseB := &models.ProgramPromotionTb{}
	baseC := &models.MerchantLocationTb{}

	join := fmt.Sprintf("INNER JOIN %s b ON a.%s = b.%s INNER JOIN %s c ON a.%s = c.%s", baseB.TableName(), baseA.ColumnProgramId(), baseB.ColumnProgramId(), baseC.TableName(), baseA.ColumnLocationId(), baseC.ColumnLocationId())
	where := fmt.Sprintf("a.%s = ? and c.%s = ?", baseA.ColumnProgramId(), baseC.ColumnActivate())
	err := m.db.Debug().Select("c.*").Table(baseA.TableName()+" a").Joins(join).Where(where, programId, activate).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, err
}

func (m *promotionsRepo) AfiliateGetRemainPromotionCountFrProgramId(programId int64) (*models.TotalRemainPromotionProgramId, error) {
	defer utils.ExecTime(utils.GetTimeUTC7(), "AfiliateGetListLocationFrProgramId", nil)
	result := &models.TotalRemainPromotionProgramId{}
	base := models.PromotionTb{}
	where := fmt.Sprintf("%s = ?", base.ColumnProgramId())
	err := m.db.Debug().Select("SUM(total_current) as sum_total_current, SUM(total_use) as sum_total_use, SUM(total_release) as sum_total_release").
		Table(base.TableName()).Where(where, programId).Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, err
}

func (m *promotionsRepo) AfiliateLookupCustomerHasVoucherFrProgramId(customer *models.UserInfo, programId int64) (*models.InfoPromotionCodeCus, error) {
	defer utils.ExecTime(utils.GetTimeUTC7(), "AfiliateLookupCustomerHasVoucherFrProgramId", nil)
	// info := &models.InfoPromotionCodeCus{}
	result := []models.InfoPromotionCodeCus{}
	baseHis := &models.HistoriesUsableTb{}
	basePro := &models.PromotionTb{}
	join := fmt.Sprintf("INNER JOIN %s b ON a.%s = b.%s ", basePro.TableName(), baseHis.ColumnPromotionCode(), basePro.ColumnPromotionCode())
	where := fmt.Sprintf("a.%s = ? AND b.%s = ?", baseHis.ColumnCustomerId(), basePro.ColumnProgramId())
	err := m.db.Debug().Select("a.*, b.promotion_end_usable").Table(baseHis.TableName()+" a").Joins(join).Where(where, customer.CustomerId, programId).Find(&result).Error
	if err != nil {
		internal.Log.Error("AfiliateLookupCustomerHasVoucherFrProgramId", zap.Any("customer", customer), zap.Any("programId", programId), zap.Error(err))
		return nil, err
	}
	if len(result) > 0 {
		return &result[len(result)-1], nil
	} else {
		return nil, nil
	}
}

func (m *promotionsRepo) AfiliateLookupCustomerHasVoucherFrVoucherCode(customer *models.UserInfo, voucherCode string) (*models.InfoPromotionCodeCus, error) {
	defer utils.ExecTime(utils.GetTimeUTC7(), "AfiliateLookupCustomerHasVoucherFrVoucherCode", nil)
	result := []models.InfoPromotionCodeCus{}
	baseHis := &models.HistoriesUsableTb{}
	basePro := &models.PromotionTb{}
	join := fmt.Sprintf("INNER JOIN %s b ON a.%s = b.%s ", basePro.TableName(), baseHis.ColumnPromotionCode(), basePro.ColumnPromotionCode())
	where := fmt.Sprintf("a.%s = ? AND b.%s = ?", baseHis.ColumnCustomerId(), basePro.ColumnPromotionCode())
	err := m.db.Debug().Select("a.*, b.promotion_end_usable").Table(baseHis.TableName()+" a").Joins(join).Where(where, customer.CustomerId, voucherCode).Find(&result).Error
	if err != nil {
		internal.Log.Error("AfiliateLookupCustomerHasVoucherFrProgramId", zap.Any("customer", customer), zap.Any("voucherCode", voucherCode), zap.Error(err))
		return nil, err
	}
	if len(result) > 0 {
		return &result[len(result)-1], nil
	} else {
		return nil, nil
	}
}

func (p *promotionsRepo) CreateListPromotion(listPromotions []models.PromotionTb) error {
	return p.db.CreateInBatches(listPromotions, 10000).Error
}

func (p *promotionsRepo) CheckCodeIsExist(listPromotion []string) ([]string, error) {
	var listCode []string
	err := p.db.Select("promotion_code").Model(&models.PromotionTb{}).Where("promotion_code IN (?)", listPromotion).Find(&listCode).Error
	if err != nil {
		return nil, err
	}
	return listCode, nil
}
func (p *promotionsRepo) CreatePromotion(promotion *models.PromotionTb) error {
	return p.db.Create(promotion).Error
}

func (p *promotionsRepo) CheckPromotionCodeIsExist(code string) (int, error) {
	var count int64
	err := p.db.Model(&models.PromotionTb{}).Where("promotion_code = ?", code).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
func (p *promotionsRepo) GetPromotionByProgramId(programId int64, funcName string) ([]models.PromotionTb, error) {
	ctx := context.WithValue(context.Background(), internal.FuncNameKey, funcName)
	var listPromotion []models.PromotionTb
	err := p.db.WithContext(ctx).Where("program_id = ? AND state_promotion = 1", programId).Find(&listPromotion).Error
	if err != nil {
		return nil, err
	}
	return listPromotion, nil
}

func (p *promotionsRepo) GetPromotionsPriority(screenList []string, limit, offset int, now string) ([]models.ProgramPromotionPriority, error) {
	result := []models.ProgramPromotionPriority{}
	base := models.ProgramPromotionPriority{}
	where := fmt.Sprintf("%s = ? AND (%s = ? OR (%s = ? AND '%s' between %s and %s))", base.ColumnState(), base.ColumnActiveTime(), base.ColumnActiveTime(), now, base.ColumnStartTime(), base.ColumnEndTime())
	args := []interface{}{1, 0, 1}
	if len(screenList) > 0 {
		where += fmt.Sprintf(" AND %s IN (?)", base.ColumnPriorityCategory())
		args = append(args, screenList)
	}
	err := p.db.Debug().Where(where, args...).Order(base.ColumnPriority()).Limit(limit).Offset(offset).Find(&result).Error
	return result, err
}
func (p *promotionsRepo) CheckInsertPriority(priority int, now string) ([]models.ProgramPromotionPriority, error) {
	result := []models.ProgramPromotionPriority{}
	base := models.ProgramPromotionPriority{}
	where := fmt.Sprintf("%s = ? AND %s = ? AND (%s = ? OR (%s = ? AND '%s' between %s and %s))", base.ColumnPriority(), base.ColumnState(), base.ColumnActiveTime(), base.ColumnActiveTime(), now, base.ColumnStartTime(), base.ColumnEndTime())
	err := p.db.Debug().Where(where, priority, 1, 0, 1).Find(&result).Error
	return result, err
}
func (p *promotionsRepo) CheckUpdatePriority(idInput int) (*models.ProgramPromotionPriority, error) {
	result := &models.ProgramPromotionPriority{}
	where := fmt.Sprintf("%s = ?", result.ColumnId())
	err := p.db.Debug().Where(where, idInput).First(result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return result, err
}
func (p *promotionsRepo) CreatePriority(input *models.ProgramPromotionPriority) (*models.ProgramPromotionPriority, error) {
	err := p.db.Debug().Create(input).Error
	return input, err
}

func (p *promotionsRepo) UpdatePriority(input *models.ProgramPromotionPriority) error {
	return p.db.Debug().Model(&models.ProgramPromotionPriority{}).Where(fmt.Sprintf("%s = ?", input.ColumnId()), input.Id).Updates(map[string]interface{}{
		input.ColumnPriority():   input.Priority,
		input.ColumnActiveTime(): input.ActiveTime,
		input.ColumnStartTime():  input.StartTime,
		input.ColumnEndTime():    input.EndTime,
		input.ColumnState():      input.State,
		input.ColumnMetaData():   input.MetaData,
		input.ColumnUpdatedAt():  input.UpdatedAt,
	}).Error
}
