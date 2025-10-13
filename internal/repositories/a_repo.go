package repositories

import (
	"gorm.io/gorm"
)

type Repositories struct {
	ControlApi        ControlApiRepo
	StoreCallApi      StoreCallApiRepo
	Promotions        PromotionsRepo
	Categories        CategoryRepo
	ProgramPromotions ProgramPromotionRepo
	ProgramCategories ProgramCategoryRepo
	ProgramProducts   ProgramProductRepo
	ProgramDirect     ProgramDirectRepo
	LoginApi          LoginRepo
}

func NewRepositories(
	db *gorm.DB,
) *Repositories {
	return &Repositories{
		LoginApi:          NewLoginRepo(db),
		ControlApi:        NewControlApiRepo(db),
		StoreCallApi:      NewStoreCallApiRepo(db),
		Promotions:        NewPromotionsRepo(db),
		Categories:        NewCategoryRepo(db),
		ProgramPromotions: NewProgramPromotionsRepo(db.Debug()),
		ProgramCategories: NewProgramCategoryRepo(db),
		ProgramProducts:   NewProgramProductRepo(db),
		ProgramDirect:     NewProgramDirectRepo(db),
	}
}
