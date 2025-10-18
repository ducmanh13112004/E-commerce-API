package programpromotionservice

import (
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/models"
	"ecom_promotion_v2/internal/repositories"
	"ecom_promotion_v2/internal/utils"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

type ProgramPromotionService interface {
	GetProgramPromotionList(filter *models.FilterProgramPromotion, funcName string) (interface{}, *internal.SystemStatus)
	GetProgramPromotionByID(id int64, funcName string) (interface{}, *internal.SystemStatus)
	CreateProgramPromotion(promotion *models.CreateProgramPromotion, updateBy string, funcName string) (interface{}, *internal.SystemStatus)
	UpdateProgramPromotion(promotion *models.UpdateProgramPromotion, updateBy string, funcName string) *internal.SystemStatus
	OnOffProgramPromotion(program *models.OnOffProgramPromotion, updateBy string, funcName string) *internal.SystemStatus
	DeleteProgramPromotion(id int64, updateBy string, funcName string) *internal.SystemStatus
	GetReportProgramPromotion(filter *models.FilterReportProgramPromotion) (interface{}, *internal.SystemStatus)
	GetReportProgramPromotionLocal(listProgramId []int64) (interface{}, *internal.SystemStatus)
}

type programPromotionService struct {
	repo *repositories.Repositories
}

func NewProgramPromotionService(repo *repositories.Repositories) ProgramPromotionService {
	return &programPromotionService{repo: repo}
}

// Sắp bắt đầu: 1
// Sắp kết thúc: 2
// Đang diễn ra: 3
// Đã kết thúc: 4
// Tạm ẩn: 5
func (p *programPromotionService) setStateTime(stateTime *int, program *models.ProgramPromotionTb) {
	now := utils.GetTimeUTC7()
	nowFormat := now.Format("2006-01-02")
	isDurationActive := program.BeginUsable.Before(now) && program.EndUsable.After(now)
	isTimeInDay := program.BeginUsable.Format("2006-01-02") == nowFormat || program.EndUsable.Format("2006-01-02") == nowFormat
	if isDurationActive || isTimeInDay {
		if program.State == 1 {
			// Chương trình đang diễn ra
			*stateTime = 3
			if program.EndUsable.Sub(now).Hours() <= 24 {
				//Sắp kết thúc trong 24h
				*stateTime = 2
			}
		} else if program.State == 0 {
			// Chương trình tạm ẩn
			*stateTime = 5
		}

	} else {
		if program.EndUsable.Before(now) {
			// Chương trình đã kết thúc
			*stateTime = 4
		} else if program.BeginUsable.Sub(now).Hours() <= 24 {
			// Chương trình sắp bắt đầu
			*stateTime = 1
		}
	}

}

func (p *programPromotionService) convertToMapProgramPromotion(program *models.ProgramPromotionTb, stateTime int) map[string]interface{} {
	if program == nil {
		return nil
	}
	var beginUsable string
	var endUsable string
	var percentageReduction int
	var reductionAmountMax int
	var rank int
	var tAction string
	if program.BeginUsable != nil {
		beginUsable = program.BeginUsable.Format("2006-01-02")
	}
	if program.EndUsable != nil {
		endUsable = program.EndUsable.Format("2006-01-02")
	}
	if program.PercentageReduction != nil {
		percentageReduction = *program.PercentageReduction
	}
	if program.ReductionAmountMax != nil {
		reductionAmountMax = *program.ReductionAmountMax
	}
	if program.Rank != nil {
		rank = *program.Rank
	}
	if program.TAction != nil {
		tAction = program.TAction.Format("2006-01-02 15:04:05")
	}
	var listChannel []string
	for _, channel := range program.ProgramPromotionChannels {
		listChannel = append(listChannel, channel.Channel)
	}
	return map[string]interface{}{
		"program_id":           program.ProgramId,
		"program_name":         program.ProgramName,
		"merchant_id":          program.MerchantId,
		"program_type":         program.ProgramType,
		"promotion_title":      program.PromotionTitle,
		"promotion_sub_title":  program.PromotionSubTitle,
		"description":          program.Description,
		"begin_usable":         beginUsable,
		"end_usable":           endUsable,
		"reduction_type":       program.ReductionType,
		"total_payment_min":    program.TotalPaymentMin,
		"percentage_reduction": percentageReduction,
		"reduction_amount_max": reductionAmountMax,
		"user_guide":           program.UserGuide,
		"condition_rev_code":   program.ConditionRevCode,
		"times_each_customer":  program.TimesEachCustomer,
		"state":                program.State,
		"src_img_icon":         program.SrcImgIcon,
		"rank":                 rank,
		"t_create":             program.TCreate.Format("2006-01-02 15:04:05"),
		"bypass_owner":         program.BypassOwner,
		"redirect_id":          program.RedirectId,
		"code_type":            program.CodeType,
		"update_by":            program.UpdateBy,
		"t_action":             tAction,
		"customer_type":        program.CustomerType,
		"is_update_manual":     program.IsUpdateManual,
		"is_share":             program.IsShare,
		"state_time":           stateTime,
		"src_img_detail":       program.SrcImgDetail,
		"list_channel":         listChannel,
	}
}
func (p *programPromotionService) GetProgramPromotionList(filter *models.FilterProgramPromotion, funcName string) (interface{}, *internal.SystemStatus) {
	internal.Log.Info("GetProgramPromotionList", zap.Any("funcName", funcName), zap.Any("input", filter))
	programList, err := p.repo.ProgramPromotions.GetProgramPromotionWithFilter(filter)
	if err != nil {
		internal.Log.Error("Error DB", zap.Any("funcName", funcName), zap.Any("input", filter), zap.Error(err))
		return nil, internal.SysStatus.SystemError
	}
	resp := []map[string]interface{}{}

	for _, program := range programList {
		var stateTime int
		p.setStateTime(&stateTime, &program)
		if filter.StateTime != nil && *filter.StateTime != stateTime {
			continue
		}
		resp = append(resp, p.convertToMapProgramPromotion(&program, stateTime))
	}
	return resp, nil
}

func (p *programPromotionService) GetProgramPromotionByID(id int64, funcName string) (interface{}, *internal.SystemStatus) {
	internal.Log.Info("GetProgramPromotionByIDService", zap.Any("funcName", funcName), zap.Any("input", id))
	program, err := p.repo.ProgramPromotions.GetProgramPromotionByID(id, funcName)
	if err != nil {
		internal.Log.Error("Error DB", zap.Any("funcName", funcName), zap.Any("input", id), zap.Error(err))
		return nil, internal.SysStatus.SystemError
	}
	var stateTime int
	p.setStateTime(&stateTime, program)
	return p.convertToMapProgramPromotion(program, stateTime), nil
}

func (p *programPromotionService) CreateProgramPromotion(promotion *models.CreateProgramPromotion, updateBy string, funcName string) (interface{}, *internal.SystemStatus) {
	internal.Log.Info("CreateProgramPromotion", zap.Any("funcName", funcName), zap.Any("input", promotion))
	now := utils.GetTimeUTC7()
	//validate date YYYY-MM-DD
	const layout = "2006-01-02"
	var beginUsable, endUsable time.Time
	beginUsable, err := time.Parse(layout, promotion.BeginUsable)
	if err != nil {
		internal.Log.Error("Error parse beginDate", zap.Any("funcName", funcName), zap.Any("input", promotion), zap.Error(err))
		return nil, internal.SysStatus.SystemError
	}

	endUsable, err = time.Parse(layout, promotion.EndUsable)
	if err != nil {
		internal.Log.Error("Error parse endDate", zap.Any("funcName", funcName), zap.Any("input", promotion), zap.Error(err))
		return nil, internal.SysStatus.SystemError
	}
	idString := now.Format("20060102150405")
	idInt64, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		internal.Log.Error("Error parse id", zap.Any("funcName", funcName), zap.Any("input", promotion), zap.Error(err))
		return nil, internal.SysStatus.SystemError
	}
	programPromotion := &models.ProgramPromotionTb{
		ProgramId:           idInt64,
		ProgramName:         promotion.ProgramName,
		MerchantId:          promotion.MerchantID,
		ProgramType:         promotion.ProgramType,
		PromotionTitle:      promotion.PromotionTitle,
		PromotionSubTitle:   promotion.PromotionSubTitle,
		Description:         promotion.Description,
		BeginUsable:         &beginUsable,
		EndUsable:           &endUsable,
		ReductionType:       promotion.ReductionType,
		TotalPaymentMin:     promotion.TotalPaymentMin,
		PercentageReduction: &promotion.PercentageReduction,
		ReductionAmountMax:  &promotion.ReductionAmountMax,
		UserGuide:           promotion.UserGuide,
		ConditionRevCode:    promotion.ConditionRevCode,
		TimesEachCustomer:   promotion.TimesEachCustomer,
		State:               promotion.State,
		SrcImgIcon:          promotion.SrcImgIcon,
		Rank:                &promotion.Rank,
		BypassOwner:         promotion.BypassOwner,
		RedirectId:          promotion.RedirectID,
		IsShare:             promotion.IsShare,
		CustomerType:        promotion.CustomerType,
		UpdateBy:            updateBy,
		TCreate:             now,
		TAction:             &now,
		SrcImgDetail:        promotion.SrcImgDetail,
		IsUpdateManual:      promotion.IsUpdateManual,
		TypeChoose:          "SINGLE",
	}
	resp, err := p.repo.ProgramPromotions.CreateProgramPromotion(programPromotion, promotion.ListChannel)
	if err != nil {
		// go utils.SendTelegramMessage(fmt.Sprintf("Thêm chương trình khuyến mại thất bại!\nUser: %s - ID: %s", updateBy, strconv.FormatInt(resp, 10)), 2)
		internal.Log.Error("Error create program", zap.Any("funcName", funcName), zap.Any("input", promotion), zap.Error(err))
		return nil, internal.SysStatus.SystemError
	}
	// go utils.SendTelegramMessage(fmt.Sprintf("Chương trình khuyến mại mới vừa được thêm\nUser: %s - ID: %s", updateBy, strconv.FormatInt(resp, 10)), 1)
	internal.Log.Info("Success create program", zap.Any("funcName", funcName), zap.Any("input", promotion), zap.Any("output", resp))
	return resp, nil
}

func (p *programPromotionService) OnOffProgramPromotion(program *models.OnOffProgramPromotion, updateBy string, funcName string) *internal.SystemStatus {
	internal.Log.Info("OnOffProgramPromotion", zap.Any("funcName", funcName), zap.Any("input", program))
	programPromotion, err := p.repo.ProgramPromotions.GetProgramPromotionByID(program.ProgramId, funcName)
	if err != nil {
		internal.Log.Error("Error DB", zap.Any("funcName", funcName), zap.Any("input", program), zap.Error(err))
		return internal.SysStatus.SystemError
	}

	now := utils.GetTimeUTC7()
	programPromotion.State = program.State
	programPromotion.UpdateBy = updateBy
	programPromotion.TAction = &now
	err = p.repo.ProgramPromotions.UpdateProgramPromotion(programPromotion)
	if err != nil {
		internal.Log.Error("Error DB", zap.Any("funcName", funcName), zap.Any("input", program), zap.Error(err))
		return internal.SysStatus.SystemError
	}
	return nil
}

func (p *programPromotionService) UpdateProgramPromotion(promotion *models.UpdateProgramPromotion, updateBy string, funcName string) *internal.SystemStatus {
	programPromotion, errDB := p.repo.ProgramPromotions.GetProgramPromotionByID(promotion.ProgramID, funcName)
	internal.Log.Info("UpdateProgramPromotion", zap.Any("funcName", funcName), zap.Any("input", promotion))
	if errDB != nil {
		internal.Log.Error("Error ProgramPromotions.GetProgramPromotionByID", zap.Any("funcName", funcName), zap.Any("input", promotion), zap.Error(errDB))
		return internal.SysStatus.SystemError
	}

	//validate date YYYY-MM-DD
	const layout = "2006-01-02"
	var beginUsable, endUsable time.Time
	beginUsable, err := time.Parse(layout, promotion.BeginUsable)
	if err != nil {
		internal.Log.Error(funcName, zap.Any("input", promotion), zap.Error(err))
		return internal.SysStatus.SystemError
	}

	endUsable, err = time.Parse(layout, promotion.EndUsable)
	if err != nil {
		internal.Log.Error(funcName, zap.Any("input", promotion), zap.Error(err))
		return internal.SysStatus.SystemError
	}

	now := utils.GetTimeUTC7()
	programPromotion.ProgramName = promotion.ProgramName
	programPromotion.MerchantId = promotion.MerchantID
	programPromotion.ProgramType = promotion.ProgramType
	programPromotion.PromotionTitle = promotion.PromotionTitle
	programPromotion.PromotionSubTitle = promotion.PromotionSubTitle
	programPromotion.Description = promotion.Description
	programPromotion.BeginUsable = &beginUsable
	programPromotion.EndUsable = &endUsable
	programPromotion.ReductionType = promotion.ReductionType
	programPromotion.TotalPaymentMin = promotion.TotalPaymentMin
	programPromotion.PercentageReduction = &promotion.PercentageReduction
	programPromotion.ReductionAmountMax = &promotion.ReductionAmountMax
	programPromotion.UserGuide = promotion.UserGuide
	programPromotion.ConditionRevCode = promotion.ConditionRevCode
	programPromotion.TimesEachCustomer = promotion.TimesEachCustomer
	programPromotion.State = promotion.State
	programPromotion.SrcImgIcon = promotion.SrcImgIcon
	programPromotion.Rank = &promotion.Rank
	programPromotion.IssueCode = promotion.IssueCode
	programPromotion.PublishSource = promotion.PublishSource
	programPromotion.WebsiteLink = promotion.WebsiteLink
	programPromotion.BypassOwner = promotion.BypassOwner
	programPromotion.RedirectId = promotion.RedirectID
	programPromotion.CodeType = promotion.CodeType
	programPromotion.DayNotUsed = promotion.DayNotUsed
	programPromotion.IsShare = promotion.IsShare
	programPromotion.CustomerType = promotion.CustomerType
	programPromotion.UpdateBy = updateBy
	programPromotion.TAction = &now
	programPromotion.SrcImgDetail = promotion.SrcImgDetail
	programPromotion.IsUpdateManual = promotion.IsUpdateManual

	err = p.repo.ProgramPromotions.UpdateProgramPromotionChannel(programPromotion, promotion.ListChannel)
	if err != nil {
		internal.Log.Error(funcName, zap.Any("input", promotion), zap.Error(err))
		return internal.SysStatus.SystemError
	}
	return nil
}

func (p *programPromotionService) DeleteProgramPromotion(id int64, updateBy string, funcName string) *internal.SystemStatus {
	internal.Log.Info("DeleteProgramPromotion", zap.Any("funcName", funcName), zap.Any("input", id))
	program, err := p.repo.ProgramPromotions.GetProgramPromotionByID(id, funcName)
	if err != nil {
		internal.Log.Error("Error GetProgramPromotionId", zap.Any("funcName", funcName), zap.Any("input", id), zap.Error(err))
		return internal.SysStatus.SystemError
	}
	now := utils.GetTimeUTC7()
	if program.EndUsable != nil && (program.EndUsable.After(now) || program.EndUsable.Format("2006-01-02") == now.Format("2006-01-02")) {
		previousEndUsable := now.AddDate(0, 0, -1)
		program.EndUsable = &previousEndUsable
	}
	program.IsDeleted = 1
	program.UpdateBy = updateBy
	program.TAction = &now
	program.State = 0
	err = p.repo.ProgramPromotions.UpdateProgramPromotion(program)
	if err != nil {
		internal.Log.Error("Error DeleteProgramPromotion", zap.Any("funcName", funcName), zap.Any("input", id), zap.Error(err))
		return internal.SysStatus.SystemError
	}
	return nil
}

func (p *programPromotionService) GetReportProgramPromotion(filter *models.FilterReportProgramPromotion) (interface{}, *internal.SystemStatus) {
	funcName := "GetReportProgramPromotion"
	internal.Log.Info("Start GetReportProgramPromotion", zap.Any("funcName", funcName), zap.Any("input", filter))
	if len(filter.ListProgramId) > 0 {
		programIdNotExist, err := p.repo.ProgramPromotions.CheckProgramPromotionID(filter.ListProgramId)
		if err != nil {
			internal.Log.Error(funcName, zap.Any("input", filter), zap.Error(err))
			return nil, internal.SysStatus.SystemError
		}
		if len(programIdNotExist) >= 0 && len(programIdNotExist) < len(filter.ListProgramId) {
			internal.Log.Error(funcName, zap.Any("input", filter.ListProgramId), zap.Error(err))
			removeMap := make(map[string]bool)
			for _, str := range programIdNotExist {
				removeMap[str] = true
			}

			var result []string
			for _, str := range filter.ListProgramId {
				id := fmt.Sprintf("%d", str)
				if !removeMap[id] {
					result = append(result, id)
				}
			}
			return nil, &internal.SystemStatus{
				Status: internal.CODE_SYSTEM_ERROR,
				Msg:    "ID Chương trình khuyến mãi sau không tồn tại: " + strings.Join(result, ", "),
			}
		}
	}

	result, err := p.repo.ProgramPromotions.GetReportProgramPromotion(filter)
	if err != nil {
		internal.Log.Error(funcName, zap.Any("input", filter.ListProgramId), zap.Error(err))
		return nil, internal.SysStatus.SystemError
	}
	internal.Log.Info("End GetReportProgramPromotion", zap.Any("funcName", funcName), zap.Any("input", filter), zap.Any("rowEffect", len(result)))
	return result, nil
}

func (p *programPromotionService) GetReportProgramPromotionLocal(listProgramId []int64) (interface{}, *internal.SystemStatus) {
	funcName := "GetReportProgramPromotion"
	internal.Log.Info("GetReportProgramPromotionLocal", zap.Any("funcName", funcName), zap.Any("input", listProgramId))

	programIdNotExist, err := p.repo.ProgramPromotions.CheckProgramPromotionID(listProgramId)
	if err != nil {
		internal.Log.Error(funcName, zap.Any("input", listProgramId), zap.Error(err))
		return nil, internal.SysStatus.SystemError
	}
	if len(programIdNotExist) >= 0 && len(programIdNotExist) < len(listProgramId) {
		internal.Log.Error(funcName, zap.Any("input", listProgramId), zap.Error(err))
		return nil, &internal.SystemStatus{
			Status: internal.CODE_SYSTEM_ERROR,
			Msg:    internal.MSG_SYSTEM_ERROR,
		}
	}

	result, err := p.repo.ProgramPromotions.GetReportProgramPromotion(&models.FilterReportProgramPromotion{ListProgramId: listProgramId})
	if err != nil {
		internal.Log.Error(funcName, zap.Any("input", listProgramId), zap.Error(err))
		return nil, internal.SysStatus.SystemError
	}
	internal.Log.Info("End GetReportProgramPromotion", zap.Any("funcName", funcName), zap.Any("input", listProgramId), zap.Any("rowEffect", len(result)))
	return result, nil
}
