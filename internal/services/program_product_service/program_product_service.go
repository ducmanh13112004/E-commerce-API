package programproductservice

import (
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/models"
	"ecom_promotion_v2/internal/repositories"

	"go.uber.org/zap"
)

type ProgramProductService interface {
	CreateListSku(input *models.InsertListSku) *internal.SystemStatus
}

type programProductService struct {
	repo *repositories.Repositories
}

func NewProgramProductService(repo *repositories.Repositories) ProgramProductService {
	return &programProductService{repo: repo}
}

func (p *programProductService) CreateListSku(input *models.InsertListSku) *internal.SystemStatus {
	funcName := "CreateListSku"
	internal.Log.Info("CreateListSkuService", zap.Any("funcName", funcName), zap.Any("input", input))
	var listProgramProduct []models.ProgramProductTb
	for _, programProduct := range input.ListSku {
		listProgramProduct = append(listProgramProduct, models.ProgramProductTb{
			ProgramID: input.ProgramID,
			SKU:       programProduct.Sku,
			Discount:  programProduct.DiscountVoucher,
		})
	}
	err := p.repo.ProgramProducts.CreateListProgramProduct(listProgramProduct)
	if err != nil {
		internal.Log.Error(funcName, zap.Any("input", input), zap.Error(err))
		return internal.SysStatus.SystemError
	}
	return nil
}
