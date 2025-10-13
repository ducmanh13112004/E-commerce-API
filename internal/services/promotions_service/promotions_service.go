package promotions

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
	"golang.org/x/mod/semver"
)

type PromotionsService interface {
	ReceiveVoucherFrProgramId(*models.InputReceiveVoucher) (*models.RespLocal, *internal.SystemStatus)

	//Portal
	GenerateListPromotionCode(quantity int, prefix string) (interface{}, *internal.SystemStatus)
	SaveListPromotion(listPromotions []models.CreatePromotion) *internal.SystemStatus
	CreatePromotion(funcName string, input *models.InputCreatePromotion) *internal.SystemStatus
	GetListPromotionByProgramId(programId int64, funcName string) (interface{}, *internal.SystemStatus)
}
type promotionsService struct {
	repo *repositories.Repositories
}

func NewPromotionsService(
	repo *repositories.Repositories,
) PromotionsService {
	return &promotionsService{
		repo: repo,
	}
}

func (s *promotionsService) ReceiveVoucherFrProgramId(input *models.InputReceiveVoucher) (*models.RespLocal, *internal.SystemStatus) {
	defer utils.ExecTime(utils.GetTimeUTC7(), "ReceiveVoucherFrProgramId", nil)
	result, err := s.repo.Promotions.SPReceiveVoucher(input)
	if err != nil {
		internal.Log.Error("SPReceiveVoucher", zap.Any("input", input), zap.Error(err))
		return nil, &internal.SystemStatus{
			Status: 0,
			Msg:    err.Error(),
		}
	}
	if result.VoucherInfo != nil {
		var redirectInfo interface{}
		promotionEndUsable := utils.GetStringTime(*result.VoucherInfo.PromotionEndUsable, "Y-M-D")
		if result.ProgramInfo != nil && result.ProgramInfo.RedirectInfo != nil {
			infoLink := result.ProgramInfo.RedirectInfo
			nameLink := infoLink.NameLink
			if semver.Compare("", "9.0") < 0 {
				nameLink = infoLink.NameLinkVersionOld
			}
			redirectInfo = &models.RedirectInfoProgram{
				RedirectId:   infoLink.RedirectId,
				RedirectName: infoLink.RedirectName,
				TypeLink:     infoLink.TypeLink,
				NameLink:     nameLink,
				State:        infoLink.State,
			}
		} else {
			redirectInfo = nil
		}
		return &models.RespLocal{
			StatusCode: 0,
			Message:    "Lấy mã khuyến mãi thành công!",
			Data: map[string]interface{}{
				"evoucher":             result.VoucherInfo.PromotionCode,
				"times_remain":         result.AvailableRemainTimes,
				"promotion_end_usable": promotionEndUsable,
				"info_link":            redirectInfo,
			},
		}, nil
	}
	return nil, internal.SysStatus.SystemBusy
}

func (s *promotionsService) AfiliateGetFullInfoProgramFrDictData(customer *models.UserInfo, input *models.ProgramCategory, promotionCode string) (*models.OutInfoProgram, error) {
	defer utils.ExecTime(utils.GetTimeUTC7(), "AfiliateGetFullInfoProgramFrDictData", nil)

	userRcv, err := s.repo.Promotions.UserGetCustomerReceivedCodeTimes(customer.CustomerId, input.ProgramId)
	if err != nil {
		internal.Log.Error("AfiliateGetTimesRevicedOfCustomerIdFrProgramId", zap.Any("customer", customer), zap.Any("programId", input.ProgramId), zap.Error(err))
		return nil, err
	}
	timesRemain := input.ProgramPromotionTb.TimesEachCustomer - len(userRcv)

	var listAddress interface{}
	sumAddress := 0
	if input.ProgramId == 99 {
		addressList := []string{}
		resultAddress, err := s.repo.Promotions.PaywifiGetListProvince(input.ProgramPromotionTb.MerchantId, 1)
		if err != nil {
			return nil, err
		}
		for index := range resultAddress {
			if !utils.CheckItemInListContains(addressList, resultAddress[index].Province) {
				addressList = append(addressList, resultAddress[index].Province)
			}
		}
		listAddress = addressList
		sumAddress = len(resultAddress)
	} else {
		resultAddress, err := s.repo.Promotions.AfiliateGetListLocationFrProgramId(input.ProgramId, 1)
		if err != nil {
			return nil, err
		}
		locationList := []models.LocationProgram{}
		for index := range resultAddress {
			item := resultAddress[index]
			locationList = append(locationList, models.LocationProgram{
				MerchantId:   item.MerchantId,
				StoreName:    item.StoreName,
				FullAddress:  item.FullAddress,
				Coordinate:   item.Coordinate,
				LocationIcon: item.LocationIcon,
			})
		}
		listAddress = locationList
		sumAddress = len(locationList)
	}
	//  Todo lấy tình trạng số lượng voucher hiện tại
	total, err := s.repo.Promotions.AfiliateGetRemainPromotionCountFrProgramId(input.ProgramId)
	if err != nil {
		return nil, err
	}
	// thời gian sử dụng voucher
	endUsable := utils.GetStringTime(*input.ProgramPromotionTb.EndUsable, "D/M/Y")
	//  Todo bổ sung thêm msg check
	msgCheckStatus := "CHUA_CO"
	var stateUsable *int
	promotionCode = "TESTDF6ZA16C"
	if promotionCode == "" {
		// lấy thông tin code cuối cùng của user
		infoPromotionCode, err := s.repo.Promotions.AfiliateLookupCustomerHasVoucherFrProgramId(customer, input.ProgramId)
		if err != nil {
			return nil, err
		}
		internal.Log.Info("AfiliateLookupCustomerHasVoucherFrProgramId", zap.Any("customer", customer), zap.Any("ProgramId", input.ProgramId), zap.Any("result", infoPromotionCode))
		// nếu user đã lấy code
		if infoPromotionCode != nil {
			msgCheckStatus = "DA_CO"
			promotionCode = infoPromotionCode.PromotionCode
			stateUsable = &infoPromotionCode.StateUsable
			endUsable = utils.GetStringTime(*infoPromotionCode.PromotionEndUsable, "D/M/Y")
		}
	} else {
		// lấy thông tin voucher code của user
		infoPromotionCode, err := s.repo.Promotions.AfiliateLookupCustomerHasVoucherFrVoucherCode(customer, promotionCode)
		if err != nil {
			return nil, err
		}
		internal.Log.Info("AfiliateLookupCustomerHasVoucherFrVoucherCode", zap.Any("customer", customer), zap.Any("promotionCode", promotionCode), zap.Any("result", infoPromotionCode))
		// nếu user đã lấy code
		if infoPromotionCode != nil {
			msgCheckStatus = "DA_CO"
			promotionCode = infoPromotionCode.PromotionCode
			stateUsable = &infoPromotionCode.StateUsable
			endUsable = utils.GetStringTime(*infoPromotionCode.PromotionEndUsable, "D/M/Y")
		}
	}
	infoLink := input.ProgramPromotionTb.RedirectInfo
	redirectInfo := &models.RedirectInfoProgram{}
	if infoLink != nil {
		nameLink := infoLink.NameLink
		if semver.Compare(customer.AppVersion, "9.0") < 0 {
			nameLink = infoLink.NameLinkVersionOld
		}
		redirectInfo = &models.RedirectInfoProgram{
			RedirectId:   infoLink.RedirectId,
			RedirectName: infoLink.RedirectName,
			TypeLink:     infoLink.TypeLink,
			NameLink:     nameLink,
			State:        infoLink.State,
		}
	}
	resp := &models.OutInfoProgram{
		ProgramId:            input.ProgramId,
		ProgramType:          input.ProgramPromotionTb.ProgramType,
		MerchantId:           input.ProgramPromotionTb.MerchantId,
		MerchantName:         *input.ProgramPromotionTb.MerchantInfo.MerchantName,
		CategoryName:         input.CategoryId,
		MerchantSrcImg:       *input.ProgramPromotionTb.MerchantInfo.MerchantSrcImg,
		SrcImgIcon:           input.ProgramPromotionTb.SrcImgIcon,
		PromotionTitle:       input.ProgramPromotionTb.PromotionTitle,
		PromotionSubTitle:    input.ProgramPromotionTb.PromotionSubTitle,
		Description:          input.ProgramPromotionTb.Description,
		EndUsable:            endUsable,
		BeginUsable:          utils.GetStringTime(*input.ProgramPromotionTb.BeginUsable, "D/M/Y"),
		UserGuide:            input.ProgramPromotionTb.UserGuide,
		ConditionRevCode:     input.ProgramPromotionTb.ConditionRevCode,
		TimesRemain:          timesRemain,
		RemainPromotionCount: total.TotalCurrent,
		TotalUse:             total.TotalUse,
		TotalRelease:         total.TotalRelease,
		ListAddress:          listAddress,
		SumAddress:           sumAddress,
		MsgCheckStatus:       msgCheckStatus,
		MsgPromotionCode:     promotionCode,
		InfoLink:             redirectInfo,
		StateUsable:          stateUsable,
		IsUpdateManual:       input.ProgramPromotionTb.IsUpdateManual,
		IsShare:              input.ProgramPromotionTb.IsShare,
	}
	return resp, nil
}

func (s *promotionsService) GenerateListPromotionCode(quantity int, prefix string) (interface{}, *internal.SystemStatus) {
	funcName := "GenerateListPromotionCode"
	CHARSET := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	MAX_LENGTH_PROMOTION_CODE := 6
	const MAX_VOUCHER_CODE = 3000000

	if quantity > MAX_VOUCHER_CODE {
		return nil, internal.SysStatus.SystemError
	}

	internal.Log.Info(funcName, zap.Any("quantity", quantity), zap.Any("prefix", prefix))
	// Dùng map để đảm bảo không có mã trùng lặp
	generatedCodes := make(map[string]struct{}, quantity)
	promotionCodes := make([]string, 0, quantity)

	for len(generatedCodes) < quantity {
		code, err := utils.GenerateNanoId(CHARSET, MAX_LENGTH_PROMOTION_CODE)
		if err != nil {
			internal.Log.Error(funcName, zap.Any("quantity", quantity), zap.Any("prefix", prefix), zap.Error(err))
			return nil, internal.SysStatus.SystemError
		}

		promotionCode := fmt.Sprintf("%s%s", strings.ToUpper(prefix), code)
		// Kiểm tra trùng lặp, nếu chưa có thì thêm vào
		if _, exists := generatedCodes[promotionCode]; !exists {
			generatedCodes[promotionCode] = struct{}{}
			promotionCodes = append(promotionCodes, promotionCode)
		}
	}

	return promotionCodes, nil
}

func (s *promotionsService) SaveListPromotion(listPromotions []models.CreatePromotion) *internal.SystemStatus {
	funcName := "CreateListPromotion"
	if len(listPromotions) == 0 {
		return internal.SysStatus.SystemError
	}

	var listPromotionCode []string
	for _, promotion := range listPromotions {
		listPromotionCode = append(listPromotionCode, promotion.PromotionCode)
	}
	checkCode, err := s.repo.Promotions.CheckCodeIsExist(listPromotionCode)
	if err != nil {
		internal.Log.Error(funcName, zap.Any("input", listPromotions), zap.Error(err))
		return internal.SysStatus.SystemError
	}

	if len(checkCode) > 0 {
		internal.Log.Error(funcName, zap.Any("input", listPromotions), zap.Any("checkCode", checkCode), zap.Error(err))
		return &internal.SystemStatus{
			Status: internal.CODE_SYSTEM_ERROR,
			Msg:    internal.MSG_SYSTEM_ERROR,
			Detail: map[string]interface{}{
				"list_promotion_code": checkCode,
			},
		}
	}
	var listCreatePromotion []models.PromotionTb
	layout := "2006-01-02"
	now := time.Now()

	for _, promotion := range listPromotions {
		endUsable, err := time.Parse(layout, promotion.PromotionEndUsable)
		if err != nil {
			internal.Log.Error(funcName, zap.Any("input", listPromotions), zap.Error(err))
			return internal.SysStatus.SystemError
		}
		listCreatePromotion = append(listCreatePromotion, models.PromotionTb{
			ProgramId:          promotion.ProgramId,
			PromotionCode:      promotion.PromotionCode,
			TotalRelease:       promotion.TotalRelease,
			TotalCurrent:       promotion.TotalRelease,
			TotalUse:           0,
			StatePromotion:     1,
			TCreate:            &now,
			PromotionEndUsable: &endUsable,
		})
	}
	err = s.repo.Promotions.CreateListPromotion(listCreatePromotion)
	if err != nil {
		internal.Log.Error(funcName, zap.Any("input", listPromotions), zap.Error(err))
		return internal.SysStatus.SystemError
	}

	return nil
}

func (s *promotionsService) createPromotionWithAutoGen(funcName string, input *models.InputCreatePromotion, limitPromotion int) *internal.SystemStatus {
	internal.Log.Info("Start createPromotionWithAutoGen", zap.Any("funcName", funcName), zap.Any("input", input))
	CHARSET := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	MAX_LENGTH_PROMOTION_CODE := 6
	MAX_COUNT_GENERATE := 5

	quantity := *input.Data.Quantity
	prefix := *input.Data.Prefix
	if quantity > limitPromotion {
		return &internal.SystemStatus{
			Status: internal.CODE_SYSTEM_ERROR,
			Msg:    "Số lượng tạo mã tối đa là " + strconv.Itoa(limitPromotion),
		}
	}

	internal.Log.Info(funcName, zap.Any("quantity", quantity), zap.Any("prefix", prefix))
	// Dùng map để đảm bảo không có mã trùng lặp
	generatedCodes := map[string]struct{}{}

	for len(generatedCodes) < quantity {

		for i := 0; len(generatedCodes) < quantity && i < limitPromotion; i++ {
			code, err := utils.GenerateNanoId(CHARSET, MAX_LENGTH_PROMOTION_CODE)
			if err != nil {
				internal.Log.Error("Error GenerateNanoId", zap.Any("funcName", funcName), zap.Any("quantity", quantity), zap.Any("prefix", prefix), zap.Error(err))
				return internal.SysStatus.SystemError
			}
			promotionCode := fmt.Sprintf("%s%s", strings.ToUpper(prefix), code)
			generatedCodes[promotionCode] = struct{}{}
		}

		if len(generatedCodes) < quantity {
			// if MAX_COUNT_GENERATE == 0 {
			// 	internal.Log.Error(funcName, zap.Any("quantity", quantity), zap.Any("prefix", prefix))
			// 	return &internal.SystemStatus{
			// 		Status: internal.CODE_SYSTEM_ERROR,
			// 		Msg:    "Không thể tạo mã code vì va chạm quá nhiều",
			// 	}
			// }
			// MAX_COUNT_GENERATE--
			internal.Log.Error("Unable to generate enough promotion code", zap.Any("funcName", funcName), zap.Any("quantity", quantity), zap.Any("prefix", prefix))
			return &internal.SystemStatus{
				Status: internal.CODE_SYSTEM_ERROR,
				Msg:    "Không thể tạo mã code vì va chạm quá nhiều",
			}
		}

		// Kiểm tra xem có mã nào đã tồn tại trong DB không
		var checkListCode []string
		for code := range generatedCodes {
			checkListCode = append(checkListCode, code)
		}
		existingCodes, err := s.repo.Promotions.CheckCodeIsExist(checkListCode)
		if err != nil {
			internal.Log.Error("Error CheckCodeIsExist", zap.Any("funcName", funcName), zap.Any("quantity", quantity), zap.Any("prefix", prefix), zap.Error(err))
			return internal.SysStatus.SystemError
		}

		if len(existingCodes) == 0 {
			listProgram := make([]models.PromotionTb, 0, quantity)
			now := utils.GetTimeUTC7()
			layout := "2006-01-02"
			endUsable, err := time.Parse(layout, *input.Data.PromotionEndUsable)
			if err != nil {
				internal.Log.Error("Error ParseTimePromotionEndUsable", zap.Any("funcName", funcName), zap.Any("quantity", quantity), zap.Any("prefix", prefix), zap.Error(err))
				return internal.SysStatus.SystemError
			}

			for code := range generatedCodes {
				listProgram = append(listProgram, models.PromotionTb{
					ProgramId:          input.ProgramId,
					PromotionCode:      code,
					TotalRelease:       1,
					TotalCurrent:       1,
					TotalUse:           0,
					StatePromotion:     1,
					TCreate:            &now,
					PromotionEndUsable: &endUsable,
				})
			}
			err = s.repo.Promotions.CreateListPromotion(listProgram)
			if err != nil {
				internal.Log.Error("Error ParseTimePromotionEndUsable", zap.Any("funcName", funcName), zap.Any("quantity", quantity), zap.Any("prefix", prefix), zap.Error(err))
				return internal.SysStatus.SystemError
			}
			return nil
		}

		if len(existingCodes) > 0 {
			// Xóa code trên map
			for _, existingCode := range existingCodes {
				delete(generatedCodes, existingCode)
			}

		}
		// Khi thử GEN 5 lần mà vẫn bị va chạm với promotion_code trong db, thì dừng lại và trả lỗi
		if MAX_COUNT_GENERATE == 0 {
			internal.Log.Error("Error MaximumLoopToGeneratePromotionCode", zap.Any("funcName", funcName), zap.Any("quantity", quantity), zap.Any("prefix", prefix))
			return &internal.SystemStatus{
				Status: internal.CODE_SYSTEM_ERROR,
				Msg:    "Không thể tạo mã code vì va chạm quá nhiều",
			}
		}
		MAX_COUNT_GENERATE--

	}
	internal.Log.Error("Error", zap.Any("funcName", funcName), zap.Any("quantity", quantity), zap.Any("prefix", prefix))
	return internal.SysStatus.SystemError
}

func (s *promotionsService) createPromotionWithShareCode(funcName string, input *models.InputCreatePromotion, limitPromotion int) *internal.SystemStatus {
	if *input.Data.TotalRelease > limitPromotion {
		return &internal.SystemStatus{
			Status: internal.CODE_SYSTEM_ERROR,
			Msg:    "Số lượng tạo mã tối đa là " + strconv.Itoa(limitPromotion),
		}
	}
	count, err := s.repo.Promotions.CheckPromotionCodeIsExist(*input.Data.PromotionCode)
	if err != nil {
		internal.Log.Error("Error CheckPromotionCodeIsExist", zap.Any("funcName", funcName), zap.Any("input", input), zap.Error(err))
		return internal.SysStatus.SystemError
	}
	if count > 0 {
		return &internal.SystemStatus{
			Status: internal.CODE_SYSTEM_ERROR,
			Msg:    "Mã khuyến mãi đã tồn tại",
		}
	}
	now := utils.GetTimeUTC7()
	layout := "2006-01-02"
	endUsable, err := time.Parse(layout, *input.Data.PromotionEndUsable)
	if err != nil {
		internal.Log.Error("Error ParseTimePromotionEndUsable", zap.Any("funcName", funcName), zap.Any("input", input), zap.Error(err))
		return internal.SysStatus.SystemError
	}
	promotion := &models.PromotionTb{
		ProgramId:          input.ProgramId,
		PromotionCode:      strings.ToUpper(*input.Data.PromotionCode),
		TotalRelease:       *input.Data.TotalRelease,
		TotalCurrent:       *input.Data.TotalRelease,
		TotalUse:           0,
		StatePromotion:     1,
		TCreate:            &now,
		PromotionEndUsable: &endUsable,
	}

	err = s.repo.Promotions.CreatePromotion(promotion)
	if err != nil {
		internal.Log.Error("Error CreatePromotion", zap.Any("funcName", funcName), zap.Any("input", input), zap.Error(err))
		return internal.SysStatus.SystemError
	}

	return nil
}

func (s *promotionsService) createPromotionWithExcelData(funcName string, input *models.InputCreatePromotion, limitPromotion int) *internal.SystemStatus {
	if len(input.Data.ListPromotion) > limitPromotion {
		return &internal.SystemStatus{
			Status: internal.CODE_SYSTEM_ERROR,
			Msg:    "Số lượng tạo mã tối đa là " + strconv.Itoa(limitPromotion),
		}
	}
	var listPromotion []string
	for _, promotion := range input.Data.ListPromotion {
		listPromotion = append(listPromotion, promotion.PromotionCode)
	}
	// dùng hash map kiểm tra phần tử trùng
	freqPromotion := make(map[string]int)
	duplicates := []string{}
	for _, promotion := range listPromotion {
		freqPromotion[promotion]++
	}
	for key, count := range freqPromotion {
		if count > 1 {
			duplicates = append(duplicates, key)
		}
	}
	if len(duplicates) > 0 {
		internal.Log.Error("Error CheckCodeIsExist", zap.Any("funcName", funcName), zap.Any("input", input))
		return &internal.SystemStatus{
			Status: internal.CODE_SYSTEM_ERROR,
			Msg:    "Mã khuyến mại sau trong excel bị trùng: " + strings.Join(duplicates, ", "),
		}
	}
	listExistCode, err := s.repo.Promotions.CheckCodeIsExist(listPromotion)
	if err != nil {
		internal.Log.Error("Error CheckCodeIsExist", zap.Any("funcName", funcName), zap.Any("input", input), zap.Error(err))
		return internal.SysStatus.SystemError
	}

	if len(listExistCode) > 0 {
		internal.Log.Error("Error CodeIsExist", zap.Any("funcName", funcName), zap.Any("input", input), zap.Error(err))
		var message strings.Builder
		message.WriteString("Mã khuyến mãi sau đã tồn tại: ")
		message.WriteString(strings.Join(listExistCode, ", "))
		return &internal.SystemStatus{
			Status: internal.CODE_SYSTEM_ERROR,
			Msg:    message.String(),
		}
	}

	now := utils.GetTimeUTC7()
	layout := "2006-01-02"
	listCreatePromotion := make([]models.PromotionTb, 0, len(input.Data.ListPromotion))
	for _, promotion := range input.Data.ListPromotion {
		endUsable, err := time.Parse(layout, promotion.PromotionEndUsable)
		if err != nil {
			internal.Log.Error("Error ParseTimePromotionEndUsable", zap.Any("funcName", funcName), zap.Any("input", input), zap.Error(err))
			return internal.SysStatus.WrongParams
		}
		listCreatePromotion = append(listCreatePromotion, models.PromotionTb{
			ProgramId:          input.ProgramId,
			PromotionCode:      strings.ToUpper(promotion.PromotionCode),
			TotalRelease:       1,
			TotalCurrent:       1,
			TotalUse:           0,
			StatePromotion:     1,
			TCreate:            &now,
			PromotionEndUsable: &endUsable,
		})
	}
	err = s.repo.Promotions.CreateListPromotion(listCreatePromotion)
	if err != nil {
		internal.Log.Error("Error CreateListPromotion", zap.Any("input", input), zap.Error(err))
		return internal.SysStatus.SystemError
	}
	internal.Log.Info("Success CreatePromotionService", zap.Any("funcName", funcName), zap.Any("input", input))
	return nil
}

func (s *promotionsService) CreatePromotion(funcName string, input *models.InputCreatePromotion) *internal.SystemStatus {
	internal.Log.Info("CreatePromotion", zap.Any("funcName", funcName), zap.Any("input", input))
	const config_key = "LIMIT_PROMOTION_CODE"
	config, err := s.repo.ControlApi.GetInfoFrApiName(config_key)
	if err != nil || config == nil {
		internal.Log.Error("Error GetInfoFrApiName", zap.Any("funcName", funcName), zap.Any("input", input), zap.Error(err))
		return internal.SysStatus.SystemError
	}
	if config == nil {
		internal.Log.Error("Error GetInfoFrApiName", zap.Any("funcName", funcName), zap.Any("input", input))
		return internal.SysStatus.SystemError
	}

	limitCode, err := strconv.Atoi(config.Conditions)
	if err != nil {
		internal.Log.Error("Error ConvertIntLimitCode", zap.Any("funcName", funcName), zap.Any("input", input), zap.Error(err))
		return internal.SysStatus.SystemError
	}
	internal.Log.Info("Start CreatePromotionService", zap.Any("funcName", funcName), zap.Any("input", input))
	var errService *internal.SystemStatus = nil
	if input.Type == 0 {
		errService = s.createPromotionWithAutoGen(funcName, input, limitCode)
	} else if input.Type == 1 {
		errService = s.createPromotionWithShareCode(funcName, input, limitCode)
	} else if input.Type == 2 {
		errService = s.createPromotionWithExcelData(funcName, input, limitCode)
	} else {
		errService = internal.SysStatus.WrongParams
	}
	internal.Log.Info("End CreatePromotionService", zap.Any("funcName", funcName), zap.Any("input", input))
	return errService
}

func (s *promotionsService) GetListPromotionByProgramId(programId int64, funcName string) (interface{}, *internal.SystemStatus) {
	internal.Log.Info("GetListPromotionByProgramId", zap.Any("funcName", funcName), zap.Any("input", programId))
	isProgramExist, err := s.repo.ProgramPromotions.CheckProgramIsExist(programId, funcName)
	if err != nil {
		internal.Log.Error("Error CheckProgramIsExistDB", zap.Any("funcName", funcName), zap.Any("input", programId), zap.Error(err))
		return nil, internal.SysStatus.SystemError
	}

	if !isProgramExist {
		internal.Log.Error("Error CheckProgramIsExistDB", zap.Any("funcName", funcName), zap.Any("input", programId), zap.Error(err))
		return nil, &internal.SystemStatus{
			Status: internal.CODE_SYSTEM_ERROR,
			Msg:    "Chương trình khuyến mãi không tồn tại",
		}
	}

	listPromotions, err := s.repo.Promotions.GetPromotionByProgramId(programId, funcName)
	if err != nil {
		internal.Log.Error("Error GetPromotionByProgramIdDB", zap.Any("funcName", funcName), zap.Any("input", programId), zap.Error(err))
		return nil, internal.SysStatus.SystemError
	}
	resultPromotion := make([]map[string]interface{}, len(listPromotions))
	for index, promotion := range listPromotions {
		var endUsable string
		if promotion.PromotionEndUsable != nil {
			endUsable = promotion.PromotionEndUsable.Format("02/01/2006")
		}
		var promotionCode string
		lengthCode := len(promotion.PromotionCode)
		if lengthCode <= 5 {
			promotionCode = "*****"
		}
		code := []rune(promotion.PromotionCode)
		middle := lengthCode / 2
		start := middle - 2 // Bắt đầu từ 2 ký tự trước giữa
		if start < 0 {
			start = 0
		}
		end := start + 5
		if end > lengthCode {
			end = lengthCode
		}
		// Thay thế 5 ký tự ở giữa bằng *****
		for i := start; i < end && i < lengthCode; i++ {
			code[i] = '*'
		}
		promotionCode = string(code)

		resultPromotion[index] = map[string]interface{}{
			"program_id":           promotion.ProgramId,
			"promotion_code":       promotionCode,
			"promotion_end_usable": endUsable,
			"total_release":        promotion.TotalRelease,
		}

	}
	return resultPromotion, nil
}
