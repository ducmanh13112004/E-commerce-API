package programdirectservice

import (
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/models"
	"ecom_promotion_v2/internal/repositories"
	"ecom_promotion_v2/internal/utils"
	"encoding/json"

	"go.uber.org/zap"
)

type ProgramRedirectService interface {
	GetProgramRedirectIDService(funcName string, redirectID int) (map[string]interface{}, *internal.SystemStatus)
	GetProgramRedirectListAllService(funcName string, filter *models.FilterProgramRedirectTb) ([]map[string]interface{}, *internal.SystemStatus)
	DeleteProgramRedirect(funcName string, Direct *models.DeleteProgramRedirectRequest, updateby string) *internal.SystemStatus
	UpdateProgramRedirect(funcName string, input *models.InputProgramRedirect, updateBy string) *internal.SystemStatus
	CreateProgramRedirect(funcName string, input *models.InputProgramRedirect, updateBy string) *internal.SystemStatus
	GetProgramRedirectListNameAllService(funcName string) ([]map[string]interface{}, *internal.SystemStatus)
	GetProgramListNavigationService(funcName, code string) (map[string]interface{}, *internal.SystemStatus)
	CreateProgramNavigationService(funcName string, bodyData *models.MasterDataInput) (map[string]interface{}, *internal.SystemStatus)
	UpdateProgramNavigationService(funcName string, bodyData *models.MasterDataInput) (map[string]interface{}, *internal.SystemStatus)
}

type programredirectService struct {
	repo *repositories.Repositories
}

func NewProgramRedirectService(
	repo *repositories.Repositories,
) ProgramRedirectService {
	return &programredirectService{
		repo: repo,
	}
}
func (s *programredirectService) GetProgramRedirectIDService(funcName string, redirectID int) (map[string]interface{}, *internal.SystemStatus) {
	internal.Log.Info("Start", zap.Any("funcName", funcName), zap.Int("redirect_id", redirectID))

	// 1. Lấy thông tin redirect
	direct, err := s.repo.ProgramDirect.GetProgramRedirectByID(funcName, redirectID)
	if err != nil {
		internal.Log.Error("ProgramDirect.GetProgramRedirectByID",
			zap.Any("funcName", funcName),
			zap.Int("redirect_id", redirectID),
			zap.Error(err),
		)
		return nil, internal.SysStatus.SystemError
	}
	if direct == nil {
		internal.Log.Error("NotFound", zap.Any("funcName", funcName), zap.Int("redirect_id", redirectID))
		return nil, &internal.SystemStatus{
			Status: 0,
			Msg:    "Không tìm thấy bản ghi",
		}
	}

	// 2. Lấy list channel
	channels, err := s.repo.ProgramDirect.GetChannelsByRedirectID(funcName, redirectID)
	if err != nil {
		internal.Log.Error("ProgramDirect.GetChannelsByRedirectID",
			zap.Any("funcName", funcName),
			zap.Int("redirect_id", redirectID),
			zap.Error(err),
		)
		return nil, internal.SysStatus.SystemError
	}

	// Convert dữ liệu channel
	listChannels := make([]map[string]interface{}, 0, len(channels))
	for _, ch := range channels {
		var data map[string]interface{}
		if ch.Data != nil {
			if err := json.Unmarshal([]byte(*ch.Data), &data); err != nil {
				internal.Log.Error("InvalidChannelData", zap.Any("funcName", funcName), zap.String("channel", ch.Channel), zap.Error(err))
				data = map[string]interface{}{}
			}
		} else {
			data = map[string]interface{}{}
		}

		listChannels = append(listChannels, map[string]interface{}{
			"channel":     ch.Channel,
			"action_type": ch.ActionType,
			"data":        data,
			"data_action": ch.DataAction,
		})
	}

	// 3. Chuẩn hóa kết quả trả về
	const defaultTime = "2025-01-02 00:00:00"
	formattedTime := defaultTime
	if direct.TUpdate != nil {
		formattedTime = direct.TUpdate.Format("2006-01-02 15:04:05")
	}

	result := map[string]interface{}{
		"redirect_id":           direct.RedirectId,
		"redirect_name":         direct.RedirectName,
		"state":                 direct.State,
		"update_by":             direct.UpdateBy,
		"t_update":              formattedTime,
		"list_redirect_channel": listChannels,
	}

	internal.Log.Info("End", zap.Any("funcName", funcName), zap.Int("redirect_id", redirectID))
	return result, nil
}

func (s *programredirectService) GetProgramRedirectListAllService(funcName string, filter *models.FilterProgramRedirectTb) ([]map[string]interface{}, *internal.SystemStatus) {

	internal.Log.Info("Start", zap.Any("funcName", funcName), zap.Any("filter", filter))

	// 1. Lấy dữ liệu từ repo
	redirects, err := s.repo.ProgramDirect.GetProgramRedirectList(funcName, filter)
	if err != nil {
		internal.Log.Error("ProgramDirect.GetProgramRedirectList",
			zap.Any("funcName", funcName),
			zap.Any("filter", filter),
			zap.Error(err),
		)
		return nil, internal.SysStatus.SystemError
	}

	// 2. Chuẩn hóa dữ liệu trả về
	const defaultTime = "2025-01-02 00:00:00"
	result := make([]map[string]interface{}, 0, len(redirects))

	for _, item := range redirects {
		tUpdate := defaultTime
		if item.TUpdate != nil {
			tUpdate = item.TUpdate.Format("2006-01-02 15:04:05")
		}

		record := map[string]interface{}{
			"redirect_id":   item.RedirectId,
			"redirect_name": item.RedirectName,
			"state":         item.State,
			"update_by":     item.UpdateBy,
			"t_update":      tUpdate,
		}

		result = append(result, record)
	}

	internal.Log.Info("End", zap.Any("funcName", funcName), zap.Int("total_records", len(result)))
	return result, nil
}

func (s *programredirectService) DeleteProgramRedirect(funcName string, input *models.DeleteProgramRedirectRequest, updateBy string) *internal.SystemStatus {

	internal.Log.Info("Start", zap.String("funcName", funcName), zap.Any("input", input), zap.String("update_by", updateBy))

	// 1. Lấy thông tin từ DB
	existingDirect, err := s.repo.ProgramDirect.GetProgramRedirectByID(funcName, input.RedirectId)
	if err != nil {
		internal.Log.Error("ProgramDirect.GetProgramRedirectByID",
			zap.String("funcName", funcName),
			zap.Int("redirect_id", input.RedirectId),
			zap.Error(err))
		return internal.SysStatus.SystemError
	}

	// 2. Cập nhật thông tin cần thiết
	now := utils.GetTimeUTC7()
	existingDirect.State = input.State
	existingDirect.UpdateBy = updateBy
	existingDirect.TUpdate = &now

	// 3. Ghi log trước khi thực hiện xóa
	internal.Log.Info("ProgramDirect.DeleteProgramRedirect - Before Delete",
		zap.String("funcName", funcName),
		zap.Any("record", existingDirect))

	// 4. Gọi repo để xóa (thường là soft delete = update state)
	if err := s.repo.ProgramDirect.DeleteProgramRedirect(funcName, existingDirect); err != nil {
		internal.Log.Error("ProgramDirect.DeleteProgramRedirect",
			zap.String("funcName", funcName),
			zap.Any("record", existingDirect),
			zap.Error(err))
		return internal.SysStatus.SystemError
	}

	internal.Log.Info("End", zap.String("funcName", funcName), zap.Int("redirect_id", existingDirect.RedirectId))
	return nil
}

func (s *programredirectService) UpdateProgramRedirect(funcName string, input *models.InputProgramRedirect, updateBy string) *internal.SystemStatus {
	internal.Log.Info("Start", zap.String("funcName", funcName), zap.Any("input", input), zap.String("update_by", updateBy))

	// 1. Kiểm tra record chính có tồn tại
	existingDirect, err := s.repo.ProgramDirect.GetProgramRedirectByID(funcName, input.RedirectId)
	if err != nil {
		internal.Log.Error("ProgramDirect.GetProgramRedirectByID",
			zap.String("funcName", funcName),
			zap.Int("redirect_id", input.RedirectId),
			zap.Error(err))
		return internal.SysStatus.SystemError
	}

	// 2. Cập nhật record chính
	now := utils.GetTimeUTC7()
	existingDirect.RedirectName = input.RedirectName
	existingDirect.State = input.State
	existingDirect.UpdateBy = updateBy
	existingDirect.TUpdate = &now

	if err := s.repo.ProgramDirect.UpdateProgramRedirect(funcName, existingDirect); err != nil {
		internal.Log.Error("ProgramDirect.UpdateProgramRedirect",
			zap.String("funcName", funcName),
			zap.Any("record", existingDirect),
			zap.Error(err))
		return internal.SysStatus.SystemError
	}

	// 3. Update channels
	if err := s.repo.ProgramDirect.ReplaceProgramRedirectChannels(funcName, existingDirect.RedirectId, input.ListRedirectChannel, updateBy, now); err != nil {
		internal.Log.Error("ProgramDirect.ReplaceProgramRedirectChannels",
			zap.String("funcName", funcName),
			zap.Int("redirect_id", existingDirect.RedirectId),
			zap.Any("channels", input.ListRedirectChannel),
			zap.Error(err))
		return internal.SysStatus.SystemError
	}

	internal.Log.Info("End", zap.String("funcName", funcName), zap.Int("redirect_id", existingDirect.RedirectId))
	return nil
}

func (s *programredirectService) CreateProgramRedirect(funcName string, input *models.InputProgramRedirect, updateBy string) *internal.SystemStatus {
	internal.Log.Info("Start", zap.Any("funcName", funcName), zap.Any("input", input))
	now := utils.GetTimeUTC7()

	// 1. Tạo record chính (bảng program_redirect_tb)
	redirect := models.ProgramRedirectTb{
		RedirectName: input.RedirectName,
		State:        input.State,
		UpdateBy:     updateBy,
		TCreate:      &now,
		TUpdate:      &now,
	}

	// Insert bảng cha
	if err := s.repo.ProgramDirect.Create(funcName, &redirect); err != nil {
		internal.Log.Error("ProgramDirect.Create",
			zap.Any("funcName", funcName),
			zap.Any("input", redirect),
			zap.Error(err),
		)
		return internal.SysStatus.SystemError
	}

	// 2. Tạo record bảng con (program_redirect_channel_tb)
	for _, ch := range input.ListRedirectChannel {
		var dataStr *string
		if len(ch.Data) > 0 {
			jsonBytes, _ := json.Marshal(ch.Data)
			str := string(jsonBytes)
			dataStr = &str
		} else {
			// Nếu Data rỗng → NULL trong DB
			dataStr = nil
		}

		channel := models.ProgramRedirectChannelTb{
			RedirectId: redirect.RedirectId,
			Channel:    ch.Channel,
			ActionType: ch.ActionType,
			DataAction: ch.DataAction,
			Data:       dataStr, // NULL nếu nil
			TUpdate:    &now,
		}

		if err := s.repo.ProgramDirect.CreateChannel(&channel); err != nil {
			internal.Log.Error("ProgramDirect.CreateChannel", zap.Any("funcName", funcName), zap.Any("input", channel), zap.Error(err))
			return internal.SysStatus.SystemError
		}
	}

	internal.Log.Info("End", zap.Any("funcName", funcName), zap.Int("redirect_id", redirect.RedirectId))
	return nil
}

func (s *programredirectService) GetProgramRedirectListNameAllService(funcName string) ([]map[string]interface{}, *internal.SystemStatus) {
	internal.Log.Info("Start", zap.String("funcName", funcName))

	// 1. Lấy dữ liệu từ repo
	redirectInfos, err := s.repo.ProgramDirect.GetProgramRedirectListName(funcName)
	if err != nil {
		internal.Log.Error("ProgramDirect.GetProgramRedirectListName", zap.String("funcName", funcName), zap.Error(err))
		return nil, internal.SysStatus.SystemError
	}

	// 2. Build dữ liệu trả về
	result := make([]map[string]interface{}, 0, len(redirectInfos))
	for _, info := range redirectInfos {
		record := map[string]interface{}{
			"redirect_id":   info.RedirectId,
			"redirect_name": info.RedirectName,
		}
		result = append(result, record)
	}

	// 3. Log End
	internal.Log.Info("End", zap.String("funcName", funcName), zap.Int("count", len(result)))

	return result, nil
}

func (s *programredirectService) GetProgramListNavigationService(funcName, code string) (map[string]interface{}, *internal.SystemStatus) {
	internal.Log.Info("Start", zap.String("funcName", funcName), zap.String("code", code))

	// 1. Lấy dữ liệu từ repo
	record, err := s.repo.ProgramDirect.GetProgramListNavigation(funcName, code)
	if err != nil {
		internal.Log.Error("ProgramDirect.GetProgramListNavigation",
			zap.String("funcName", funcName),
			zap.String("code", code),
			zap.Error(err))
		return nil, internal.SysStatus.SystemError
	}

	// 2. Parse JSON
	var arr []map[string]interface{}
	if err := json.Unmarshal([]byte(record.DataKeyValue), &arr); err != nil {
		internal.Log.Error("json.Unmarshal",
			zap.String("funcName", funcName),
			zap.String("raw_data", record.DataKeyValue),
			zap.Error(err))
		return nil, internal.SysStatus.SystemError
	}

	// 3. Chuẩn bị dữ liệu trả về
	result := map[string]interface{}{
		"code":           record.Code,
		"data_key_value": arr,
		"update_by":      record.UpdateBy,
		"update_at":      record.UpdateAt,
	}

	internal.Log.Info("End", zap.String("funcName", funcName), zap.String("code", code))
	return result, nil
}

func (s *programredirectService) CreateProgramNavigationService(funcName string, bodyData *models.MasterDataInput) (map[string]interface{}, *internal.SystemStatus) {
	internal.Log.Info("Start", zap.String("funcName", funcName), zap.Any("input", bodyData))

	// 1. Chuyển DataKeyValue sang JSON string
	dataBytes, err := json.Marshal(bodyData.DataKeyValue)
	if err != nil {
		internal.Log.Error("json.Marshal",
			zap.String("funcName", funcName),
			zap.Any("data_key_value", bodyData.DataKeyValue),
			zap.Error(err))
		return nil, internal.SysStatus.SystemError
	}

	// 2. Tạo record
	now := utils.GetTimeUTC7()
	record := &models.MasterData{
		Code:         bodyData.Code,
		DataKeyValue: string(dataBytes),
		UpdateBy:     &bodyData.UpdateBy,
		Active:       1,
		CreateAt:     &now,
		UpdateAt:     &now,
	}

	// 3. Insert DB
	if err := s.repo.ProgramDirect.CreateProgramNavigation(funcName, record); err != nil {
		internal.Log.Error("ProgramDirect.CreateProgramNavigation",
			zap.String("funcName", funcName),
			zap.Any("record", record),
			zap.Error(err))
		return nil, internal.SysStatus.SystemError
	}

	// 4. Chuẩn bị dữ liệu trả về
	result := map[string]interface{}{
		"id":             record.ID,
		"code":           record.Code,
		"data_key_value": bodyData.DataKeyValue,
		"update_by":      record.UpdateBy,
		"update_at":      record.UpdateAt,
	}

	internal.Log.Info("End", zap.String("funcName", funcName), zap.Int("id", record.ID))
	return result, nil
}
func (s *programredirectService) UpdateProgramNavigationService(funcName string, bodyData *models.MasterDataInput) (map[string]interface{}, *internal.SystemStatus) {
	internal.Log.Info("Start", zap.String("funcName", funcName), zap.Any("input", bodyData))

	// 1. Lấy record hiện tại
	record, err := s.repo.ProgramDirect.GetProgramListNavigation(bodyData.Code, funcName)
	if err != nil {
		internal.Log.Error("ProgramDirect.GetProgramListNavigation",
			zap.String("funcName", funcName),
			zap.String("code", bodyData.Code),
			zap.Error(err))
		return nil, internal.SysStatus.NotFound
	}

	// 2. Chuyển DataKeyValue sang JSON string
	dataBytes, err := json.Marshal(bodyData.DataKeyValue)
	if err != nil {
		internal.Log.Error("json.Marshal",
			zap.String("funcName", funcName),
			zap.Any("data_key_value", bodyData.DataKeyValue),
			zap.Error(err))
		return nil, internal.SysStatus.SystemError
	}

	// 3. Cập nhật các trường
	record.DataKeyValue = string(dataBytes)
	record.UpdateBy = &bodyData.UpdateBy
	now := utils.GetTimeUTC7()
	record.UpdateAt = &now

	// 4. Gọi repo để update
	if err := s.repo.ProgramDirect.UpdateProgramNavigation(funcName, record); err != nil {
		internal.Log.Error("ProgramDirect.UpdateProgramNavigation",
			zap.String("funcName", funcName),
			zap.Any("record", record),
			zap.Error(err))
		return nil, internal.SysStatus.SystemError
	}

	// 5. Chuẩn bị dữ liệu trả về
	result := map[string]interface{}{
		"id":             record.ID,
		"code":           record.Code,
		"data_key_value": bodyData.DataKeyValue,
		"update_by":      record.UpdateBy,
		"update_at":      record.UpdateAt,
	}

	internal.Log.Info("End", zap.String("funcName", funcName), zap.Int("id", record.ID))
	return result, nil
}
