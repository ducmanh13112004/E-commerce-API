package programdirectservice

import (
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/models"
	"ecom_promotion_v2/internal/repositories"
	"ecom_promotion_v2/internal/utils"
	"fmt"
	// "github.com/google/uuid"
	"go.uber.org/zap"
)

type ProgramDirectService interface {
	CreateProgramDirect(funcName string, direct *models.ProgramDirectScreenMobileTb) *internal.SystemStatus
	GetProgramDirectListAllService(filter *models.FilterProgramDirectScreenMobileTb) ([]map[string]interface{}, *internal.SystemStatus)
	GetProgramDirectIDService(direct string) ([]map[string]interface{}, *internal.SystemStatus)
	UpdateProgramDirect(direct *models.ProgramDirectScreenMobileTb) *internal.SystemStatus
	DeleteProgramDirect(Direct *models.DeleteProgramDirectRequest, updateby string) *internal.SystemStatus
}

type programdirectService struct {
	repo *repositories.Repositories
}

func NewProgramDirectService(
	repo *repositories.Repositories,
) ProgramDirectService {
	return &programdirectService{
		repo: repo,
	}
}

func (s *programdirectService) CreateProgramDirect(funcName string, direct *models.ProgramDirectScreenMobileTb) *internal.SystemStatus {
	internal.Log.Info("Start CreateProgramDirect", zap.Any("funcName", funcName), zap.Any("input", direct))
	// Định nghĩa bộ ký tự và độ dài redirect_id
	CHARSET := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	MAX_LENGTH_REDIRECT_ID := 6
	MAX_COUNT_GENERATE := 5
	// Cập nhật thời gian hiện tại
	now := utils.GetTimeUTC7()
	direct.TUpdate = &now

	var redirectID string
	var err error

	// Thử tạo mã redirect_id tối đa MAX_COUNT_GENERATE lần nếu có va chạm
	for i := 0; i < MAX_COUNT_GENERATE; i++ {
		redirectID, err = utils.GenerateNanoId(CHARSET, MAX_LENGTH_REDIRECT_ID)
		if err != nil {
			internal.Log.Error("Error GenerateNanoId",
				zap.Any("funcName", funcName),
				zap.Error(err),
			)
			return internal.SysStatus.SystemError
		}

		// Kiểm tra xem redirect_id
		exists, err := s.repo.ProgramDirect.CheckRedirectIdExists(funcName, []string{redirectID})
		if err != nil {
			internal.Log.Error("Error CheckRedirectIdExists",
				zap.Any("funcName", funcName),
				zap.Error(err),
			)
			return internal.SysStatus.SystemError
		}

		if len(exists) == 0 {
			direct.RedirectId = redirectID
			break
		}
	}

	// Nếu thử MAX_COUNT_GENERATE lần mà vẫn bị trùng, báo lỗi
	if direct.RedirectId == "" {
		internal.Log.Error("Error MaximumLoopToGenerateRedirectId",
			zap.Any("funcName", funcName),
			zap.Error(fmt.Errorf("không thể tạo redirect_id do va chạm quá nhiều")),
		)
		return &internal.SystemStatus{
			Status: internal.CODE_SYSTEM_ERROR,
			Msg:    "Không thể tạo redirect_id vì va chạm quá nhiều",
		}
	}

	// Lưu vào DB
	err = s.repo.ProgramDirect.CreateDirectCategory(funcName, direct)
	if err != nil {
		internal.Log.Error("CreateDirectCategory failed",
			zap.Any("funcName", funcName),
			zap.Error(err),
		)
		return internal.SysStatus.SystemError
	}

	// Ghi log khi tạo thành công
	internal.Log.Info("Category created successfully",
		zap.Any("funcName", funcName),
		zap.Any("redirect_id", direct.RedirectId),
	)

	return nil
}

func (p *programdirectService) GetProgramDirectListAllService(filter *models.FilterProgramDirectScreenMobileTb) ([]map[string]interface{}, *internal.SystemStatus) {
	funcName := "GetProgramDirectListAllService"

	directList, err := p.repo.ProgramDirect.GetProgramDirectList(funcName, filter)
	if err != nil {
		internal.Log.Error(funcName, zap.Error(err))
		return nil, internal.SysStatus.SystemError
	}

	listCategories := make([]map[string]interface{}, 0, len(directList))

	// Giá trị mặc định
	defaultTime := "2025-01-02 00:00:00"

	for _, direct := range directList {
		// Kiểm tra nếu `TUpdate` có giá trị, nếu không thì dùng giá trị mặc định
		var formattedTime string
		if direct.TUpdate != nil {
			formattedTime = direct.TUpdate.Format("2006-01-02 15:04:05")
		} else {
			formattedTime = defaultTime
		}

		// Rút gọn `redirect_id` thành 6 ký tự
		// shortRedirectID := ""
		// if len(direct.RedirectId) >= 6 {
		// 	shortRedirectID = direct.RedirectId[:6]
		// } else {
		// 	shortRedirectID = direct.RedirectId
		// }

		// Kiểm tra nếu `NameLink`, `TypeLink`, `NameLinkVersionOld` rỗng hoặc nil, trả về chuỗi `""`

		listCategories = append(listCategories, map[string]interface{}{
			"redirect_id":           direct.RedirectId, // Sử dụng ID rút gọn
			"redirect_name":         direct.RedirectName,
			"type_link":             direct.TypeLink,
			"action_type":           direct.Actiontype,
			"data_action":           direct.DataAction,
			"name_link":             direct.NameLink,
			"name_link_version_old": direct.NameLinkVersionOld,
			"note":                  direct.Note,
			"state":                 direct.State,
			"update_by":             direct.UpdateBy,
			"t_update":              formattedTime,
		})
	}

	return listCategories, nil
}

func (p *programdirectService) GetProgramDirectIDService(redirectID string) ([]map[string]interface{}, *internal.SystemStatus) {
	funcName := "GetProgramDirectIDService"
	directListID, err := p.repo.ProgramDirect.GetProgramDirectListID(funcName, redirectID)
	if err != nil {
		internal.Log.Error(funcName, zap.Error(err))
		return nil, internal.SysStatus.SystemError
	}

	ListdirectID := make([]map[string]interface{}, 0, len(directListID))

	for _, direct := range directListID {

		var formattedTime string
		if direct.TUpdate != nil {
			formattedTime = direct.TUpdate.Format("2006-01-02 15:04:05")
		} else {
			formattedTime = "2025-01-02 00:00:00"
		}

		ListdirectID = append(ListdirectID, map[string]interface{}{
			"redirect_name":         direct.RedirectName,
			"type_link":             direct.TypeLink,
			"action_type":           direct.Actiontype,
			"data_action":           direct.DataAction,
			"name_link":             direct.NameLink, // Đảm bảo giá trị không bị nil
			"name_link_version_old": direct.NameLinkVersionOld,
			"note":                  direct.Note,
			"state":                 direct.State,
			"update_by":             direct.UpdateBy,
			"t_update":              formattedTime,
		})
	}
	return ListdirectID, nil
}

func (s *programdirectService) UpdateProgramDirect(direct *models.ProgramDirectScreenMobileTb) *internal.SystemStatus {
	funcName := "UpdateProgramDirect"
	// Kiểm tra xem danh mục có tồn tại không
	existingDirect, err := s.repo.ProgramDirect.GetProgramDirectByID(direct.RedirectId, funcName)
	if err != nil {
		internal.Log.Error("GetProgramDirectByID: ProgramDirectID not found",
			zap.String("function", funcName),
			zap.String("redirect_id", direct.RedirectId), //lưu ý kiểu dữ liệu
			zap.Error(err),
		)
		return internal.SysStatus.SystemError

	}

	// Cập nhật thông tin danh mục
	now := utils.GetTimeUTC7()
	direct.TUpdate = &now
	existingDirect.RedirectName = direct.RedirectName
	existingDirect.TypeLink = direct.TypeLink
	existingDirect.Actiontype = direct.Actiontype
	existingDirect.DataAction = direct.DataAction
	existingDirect.NameLink = direct.NameLink
	existingDirect.NameLinkVersionOld = direct.NameLinkVersionOld
	existingDirect.Note = direct.Note
	existingDirect.State = direct.State
	existingDirect.UpdateBy = direct.UpdateBy
	existingDirect.TUpdate = direct.TUpdate
	// Gọi hàm cập nhật danh mục
	err = s.repo.ProgramDirect.UpdateProgramDirect(funcName, existingDirect)
	if err != nil {
		internal.Log.Error("UpdateProgramDirect failed",
			zap.String("function", funcName),
			zap.Any("input", direct),
			zap.Error(err),
		)
		return internal.SysStatus.SystemError
	}

	return nil
}

func (s *programdirectService) DeleteProgramDirect(Direct *models.DeleteProgramDirectRequest, updateby string) *internal.SystemStatus {
	funcName := "DeleteProgramDirect"

	// Lấy thông tin từ DB
	Directdb, err := s.repo.ProgramDirect.GetProgramDirectByID(Direct.RedirectId, funcName)
	if err != nil {
		internal.Log.Error("Database error when fetching Direct",
			zap.String("function", funcName),
			zap.String("redirect_id", Direct.RedirectId),
			zap.Error(err),
		)
		return internal.SysStatus.SystemError
	}

	// Cập nhật thông tin cần thiết
	now := utils.GetTimeUTC7()
	Directdb.State = Direct.State
	Directdb.UpdateBy = updateby
	Directdb.TUpdate = &now

	// Ghi log trước khi thực hiện xóa
	internal.Log.Info("Deleting programdirect",
		zap.String("function", funcName),
		zap.Any("input", Directdb),
	)

	// Gọi hàm xóa
	err = s.repo.ProgramDirect.DeleteProgramDirect(funcName, Directdb)
	if err != nil {
		internal.Log.Error("DeleteProgramDirect failed",
			zap.String("function", funcName),
			zap.Any("input", Directdb),
			zap.Error(err),
		)
		return internal.SysStatus.SystemError
	}

	return nil
}
