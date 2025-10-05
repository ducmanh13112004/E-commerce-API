package models

import (
	"time"
)

type ProgramPromotionTb struct {
	ProgramId                int64                        `gorm:"primaryKey" json:"program_id"`
	ProgramName              string                       `json:"program_name"`
	MerchantId               string                       `json:"merchant_id"`
	ProgramType              string                       `json:"program_type"`
	PromotionTitle           string                       `json:"promotion_title"`
	PromotionSubTitle        string                       `json:"promotion_sub_title"`
	Description              string                       `json:"description"`
	BeginUsable              *time.Time                   `json:"begin_usable"`
	EndUsable                *time.Time                   `json:"end_usable"`
	ReductionType            string                       `json:"reduction_type"`
	TotalPaymentMin          int                          `json:"total_payment_min"`
	PercentageReduction      *int                         `json:"percentage_reduction"`
	ReductionAmountMax       *int                         `json:"reduction_amount_max"`
	UserGuide                string                       `json:"user_guide"`
	ConditionRevCode         string                       `json:"condition_rev_code"`
	TimesEachCustomer        int                          `json:"times_each_customer"`
	State                    int                          `json:"state"`
	SrcImgIcon               string                       `json:"src_img_icon"`
	Rank                     *int                         `json:"rank"`
	TCreate                  time.Time                    `json:"t_create"`
	IssueCode                string                       `json:"issue_code"`
	PublishSource            string                       `json:"publish_source"`
	WebsiteLink              string                       `json:"website_link"`
	BypassOwner              int                          `json:"bypass_owner"`
	RedirectId               *int                         `json:"redirect_id"`
	CodeType                 string                       `json:"code_type"`
	DayNotUsed               string                       `json:"day_not_used"`
	CardType                 string                       `json:"card_type"`
	Budget                   *int                         `json:"budget"`
	Note                     string                       `json:"note"`
	UpdateBy                 string                       `json:"update_by"`
	TAction                  *time.Time                   `json:"t_action"`
	FromQuantity             int                          `json:"from_quantity"`
	TypeChoose               string                       `json:"type_choose"`
	CustomerType             string                       `json:"customer_type"`
	IsUpdateManual           int                          `json:"is_update_manual"`
	IsShare                  int                          `json:"is_share"`
	SrcImgDetail             string                       `json:"src_img_detail"`
	IsDeleted                int                          `json:"is_deleted"`
	MerchantInfo             *MerchantTb                  `gorm:"foreignKey:MerchantId;references:MerchantId" json:"merchant_info"`
	RedirectInfo             *ProgramDirectScreenMobileTb `gorm:"foreignKey:RedirectId;references:RedirectId" json:"redirect_info"`
	ProgramCategory          []ProgramCategory            `gorm:"foreignKey:ProgramId;references:ProgramId"`
	ProgramPromotionChannels []ProgramPromotionChannel    `gorm:"foreignKey:ProgramID;references:ProgramId"`
}

func (ProgramPromotionTb) TableName() string { return "program_promotion_tb" }

func (ProgramPromotionTb) ColumnProgramId() string           { return "program_id" }
func (ProgramPromotionTb) ColumnProgramName() string         { return "program_name" }
func (ProgramPromotionTb) ColumnMerchantId() string          { return "merchant_id" }
func (ProgramPromotionTb) ColumnProgramType() string         { return "program_type" }
func (ProgramPromotionTb) ColumnPromotionTitle() string      { return "promotion_title" }
func (ProgramPromotionTb) ColumnPromotionSubTitle() string   { return "promotion_sub_title" }
func (ProgramPromotionTb) ColumnDescription() string         { return "description" }
func (ProgramPromotionTb) ColumnBeginUsable() string         { return "begin_usable" }
func (ProgramPromotionTb) ColumnEndUsable() string           { return "end_usable" }
func (ProgramPromotionTb) ColumnReductionType() string       { return "reduction_type" }
func (ProgramPromotionTb) ColumnTotalPaymentMin() string     { return "total_payment_min" }
func (ProgramPromotionTb) ColumnPercentageReduction() string { return "percentage_reduction" }
func (ProgramPromotionTb) ColumnReductionAmountMax() string  { return "reduction_amount_max" }
func (ProgramPromotionTb) ColumnUserGuide() string           { return "user_guide" }
func (ProgramPromotionTb) ColumnConditionRevCode() string    { return "condition_rev_code" }
func (ProgramPromotionTb) ColumnTimesEachCustomer() string   { return "times_each_customer" }
func (ProgramPromotionTb) ColumnState() string               { return "state" }
func (ProgramPromotionTb) ColumnSrcImgIcon() string          { return "src_img_icon" }
func (ProgramPromotionTb) ColumnRank() string                { return "rank" }
func (ProgramPromotionTb) ColumnTCreate() string             { return "t_create" }
func (ProgramPromotionTb) ColumnIssueCode() string           { return "issue_code" }
func (ProgramPromotionTb) ColumnPublishSource() string       { return "publish_source" }
func (ProgramPromotionTb) ColumnWebsiteLink() string         { return "website_link" }
func (ProgramPromotionTb) ColumnBypassOwner() string         { return "bypass_owner" }
func (ProgramPromotionTb) ColumnRedirectId() string          { return "redirect_id" }
func (ProgramPromotionTb) ColumnCodeType() string            { return "code_type" }
func (ProgramPromotionTb) ColumnDayNotUsed() string          { return "day_not_used" }
func (ProgramPromotionTb) ColumnCardType() string            { return "card_type" }
func (ProgramPromotionTb) ColumnBudget() string              { return "budget" }
func (ProgramPromotionTb) ColumnNote() string                { return "note" }
func (ProgramPromotionTb) ColumnUpdateBy() string            { return "update_by" }
func (ProgramPromotionTb) ColumnTAction() string             { return "t_action" }
func (ProgramPromotionTb) ColumnFromQuantity() string        { return "from_quantity" }
func (ProgramPromotionTb) ColumnTypeChoose() string          { return "type_choose" }
func (ProgramPromotionTb) ColumnCustomerType() string        { return "customer_type" }
func (ProgramPromotionTb) ColumnIsUpdateManual() string      { return "is_update_manual" }
func (ProgramPromotionTb) ColumnIsShare() string             { return "is_share" }
func (ProgramPromotionTb) ColumnSrcImgDetail() string        { return "src_img_detail" }

type ProgramPromotionChannel struct {
	ID               int64              `gorm:"primaryKey" json:"id"`
	ProgramID        int64              `gorm:"column:program_id;not null"`
	Channel          string             `gorm:"column:channel;size:45;not null"`
	ProgramPromotion ProgramPromotionTb `gorm:"foreignKey:ProgramID;references:ProgramId"`
}

func (ProgramPromotionChannel) TableName() string { return "program_promotion_channel_tb" }

func (ProgramPromotionChannel) ColumnProgramID() string { return "program_id" }
func (ProgramPromotionChannel) ColumnChannel() string   { return "channel" }

type OutInfoProgram struct {
	ProgramId            int64                `json:"program_id"`
	ProgramType          string               `json:"program_type"`
	MerchantId           string               `json:"merchant_id"`
	MerchantName         string               `json:"merchant_name"`
	CategoryName         int                  `json:"category_name"`
	MerchantSrcImg       string               `json:"merchant_src_img"`
	SrcImgIcon           string               `json:"src_img_icon"`
	PromotionTitle       string               `json:"promotion_title"`
	PromotionSubTitle    string               `json:"promotion_sub_title"`
	Description          string               `json:"description"`
	EndUsable            string               `json:"end_usable"`
	BeginUsable          string               `json:"begin_usable"`
	UserGuide            string               `json:"user_guide"`
	ConditionRevCode     string               `json:"condition_rev_code"`
	TimesRemain          int                  `json:"times_remain"`
	RemainPromotionCount int                  `json:"remain_promotion_count"`
	TotalUse             int                  `json:"total_use"`
	TotalRelease         int                  `json:"total_release"`
	ListAddress          interface{}          `json:"list_address"`
	SumAddress           int                  `json:"sum_address"`
	MsgCheckStatus       string               `json:"msg_check_status"`
	MsgPromotionCode     string               `json:"msg_promotion_code"`
	InfoLink             *RedirectInfoProgram `json:"info_link"`
	StateUsable          *int                 `json:"state_usable"`
	IsUpdateManual       int                  `json:"is_update_manual"`
	IsShare              int                  `json:"is_share"`
}

type FilterProgramPromotion struct {
	FromDate   string `json:"from_date"`
	ToDate     string `json:"to_date"`
	MerchantID string `json:"merchant_id"`
	ProgramId  *int64 `json:"program_id,omitempty"`
	State      *int   `json:"state,omitempty"`
	StateTime  *int   `json:"state_time,omitempty" validate:"omitempty,oneof=1 2 3 4 5"`
}

type CreateProgramPromotion struct {
	ProgramName         string   `json:"program_name"`
	MerchantID          string   `json:"merchant_id"`
	ProgramType         string   `json:"program_type"`
	PromotionTitle      string   `json:"promotion_title" validate:"max=100"`
	PromotionSubTitle   string   `json:"promotion_sub_title" validate:"max=255"`
	Description         string   `json:"description"`
	BeginUsable         string   `json:"begin_usable"`
	EndUsable           string   `json:"end_usable"`
	ReductionType       string   `json:"reduction_type"`
	TotalPaymentMin     int      `json:"total_payment_min"`
	PercentageReduction int      `json:"percentage_reduction"`
	ReductionAmountMax  int      `json:"reduction_amount_max"`
	UserGuide           string   `json:"user_guide"`
	ConditionRevCode    string   `json:"condition_rev_code"`
	TimesEachCustomer   int      `json:"times_each_customer"`
	State               int      `json:"state"`
	SrcImgIcon          string   `json:"src_img_icon"`
	Rank                int      `json:"rank"`
	BypassOwner         int      `json:"bypass_owner"`
	RedirectID          *int     `json:"redirect_id"`
	IsShare             int      `json:"is_share"`
	CustomerType        string   `json:"customer_type"`
	SrcImgDetail        string   `json:"src_img_detail"`
	IsUpdateManual      int      `json:"is_update_manual"`
	ListChannel         []string `json:"list_channel" validate:"required"`
}

type UpdateProgramPromotion struct {
	ProgramID           int64    `json:"program_id"`
	ProgramName         string   `json:"program_name"`
	MerchantID          string   `json:"merchant_id"`
	ProgramType         string   `json:"program_type"`
	PromotionTitle      string   `json:"promotion_title" validate:"max=100"`
	PromotionSubTitle   string   `json:"promotion_sub_title" validate:"max=255"`
	Description         string   `json:"description"`
	BeginUsable         string   `json:"begin_usable"`
	EndUsable           string   `json:"end_usable"`
	ReductionType       string   `json:"reduction_type"`
	TotalPaymentMin     int      `json:"total_payment_min"`
	PercentageReduction int      `json:"percentage_reduction"`
	ReductionAmountMax  int      `json:"reduction_amount_max"`
	UserGuide           string   `json:"user_guide"`
	ConditionRevCode    string   `json:"condition_rev_code"`
	TimesEachCustomer   int      `json:"times_each_customer"`
	State               int      `json:"state"`
	SrcImgIcon          string   `json:"src_img_icon"`
	Rank                int      `json:"rank"`
	IssueCode           string   `json:"issue_code"`
	PublishSource       string   `json:"publish_source"`
	WebsiteLink         string   `json:"website_link"`
	BypassOwner         int      `json:"bypass_owner"`
	RedirectID          *int     `json:"redirect_id"`
	CodeType            string   `json:"code_type"`
	DayNotUsed          string   `json:"day_not_used"`
	IsShare             int      `json:"is_share"`
	CustomerType        string   `json:"customer_type"`
	SrcImgDetail        string   `json:"src_img_detail"`
	IsUpdateManual      int      `json:"is_update_manual"`
	ListChannel         []string `json:"list_channel"`
}

type OnOffProgramPromotion struct {
	ProgramId int64 `json:"program_id" validate:"required"`
	State     int   `json:"state" validate:"oneof=0 1"`
}

type ReportProgramPromotion struct {
	ProgramId    int64  `gorm:"column:program_id" json:"program_id" validate:"required"`
	ProgramName  string `gorm:"column:program_name" json:"program_name" validate:"required"`
	TotalRelease int    `gorm:"column:total_release" json:"total_release" validate:"required"`
	TotalCollect int    `gorm:"column:total_collect" json:"total_collect" validate:"required"`
	TotalUse     int    `gorm:"column:total_use" json:"total_use" validate:"required"`
}

type FilterReportProgramPromotion struct {
	ListProgramId []int64 `json:"list_program_id" validate:"omitempty,required"`
	FromDate      string  `json:"from_date" validate:"omitempty,required"`
	ToDate        string  `json:"to_date" validate:"omitempty,required"`
}

type ReportProgramPromotionLocal struct {
	ListProgramId []int64 `json:"list_program_id" validate:"required"`
}
