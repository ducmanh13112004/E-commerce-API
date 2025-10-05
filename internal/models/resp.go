package models

import "time"

type GroupData struct {
	Action string                 `json:"action" validate:"required,max=50"`
	Data   map[string]interface{} `json:"data" validate:"required"`
}

type RespWeb struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Detail interface{} `json:"detail"`
}

type RespLeb struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

type RespLocal struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

type RespLocalSnakeCase struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

type RespLocalHiInSide struct {
	StatusCode int         `json:"code"`
	Message    string      `json:"description"`
	Data       interface{} `json:"data"`
}

type RespLocalHiInSide2 struct {
	StatusCode int           `json:"code"`
	Message    string        `json:"description"`
	Data       []interface{} `json:"data"`
}

type RespLocalManAPI struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"Message"`
	Data       int    `json:"result"`
}
type InfoTransactionType struct {
	ServiceId   int    `json:"service_id"`
	ServiceName string `json:"service_name"`
	ObjId       int    `json:"ObjID"`
	LocationID  *int   `json:"locationid"`
	BrandCode   *int   `json:"branchcode"`
	BuildingId  *int   `json:"buildingid"`
	DistrictId  *int   `json:"districtid"`
	Wardid      *int   `json:"wardid"`
}
type ManAPIGetInfoTransaction struct {
	ResponseResult struct {
		ErrorCode int                   `json:"ErrorCode"`
		Message   string                `json:"Message"`
		Data      []InfoTransactionType `json:"Result"`
	} `json:"ResponseResult"`
}

type RespSalePlatformAuthentication struct {
	StatusCode int                    `json:"Code"`
	Message    map[string]interface{} `json:"Message"`
	Data       struct {
		AccessToken string `json:"AccessToken"`
	} `json:"Data"`
}

type SalePlatformCategory struct {
	ComboType            int      `json:"ComboType"`
	ServiceIds           string   `json:"ServiceIds"`
	SubServiceTypeIds    string   `json:"SubServiceTypeIds"`
	SubServiceIds        string   `json:"SubServiceIds"`
	SubServiceNames      []string `json:"SubServiceNames"`
	SysName              string   `json:"SysName"`
	FullName             string   `json:"FullName"`
	Description          string   `json:"Description"`
	ShortDescription     string   `json:"ShortDescription"`
	DisplayOrder         int      `json:"DisplayOrder"`
	RealMonth            float64  `json:"RealMonth"`
	RealPromotionMonth   float64  `json:"RealPromotionMonth"`
	PercentPromotion     int      `json:"PercentPromotion"`
	UnitDesc             string   `json:"UnitDesc"`
	FromQty              int      `json:"FromQty"`
	ToQty                int      `json:"ToQty"`
	PolicyCode           string   `json:"PolicyCode"`
	ParentCategoryType   int      `json:"ParentCategoryType"`
	ShowOnMostInterested bool     `json:"ShowOnMostInterested"`
	ShowOnBestOffer      bool     `json:"ShowOnBestOffer"`
}

type RespSalePlatformCategory struct {
	Code    int `json:"Code"`
	Message struct {
		Message   string      `json:"Message"`
		ExMessage interface{} `json:"ExMessage"`
	} `json:"Message"`
	Data []SalePlatformCategory `json:"Data"`
}
type SPCategoryPackages struct {
	SalesPolicyPackages []struct {
		PrePaid            int         `json:"PrePaid"`
		MonthUsed          int         `json:"MonthUsed"`
		RealMonth          float64     `json:"RealMonth"`
		RealPromotionMonth float64     `json:"RealPromotionMonth"`
		MonthUsedPriceVAT  float64     `json:"MonthUsedPriceVAT"`
		PrePaidPriceVAT    float64     `json:"PrePaidPriceVAT"`
		ApplyProduct       interface{} `json:"ApplyProduct"`
		DecideName         string      `json:"DecideName"`
		PolicyId           int         `json:"PolicyId"`
		PolicyCode         string      `json:"PolicyCode"`
		ProgramID          int         `json:"ProgramID"`
		PolicyName         string      `json:"PolicyName"`
		ComboType          int         `json:"ComboType"`
		PolicyType         int         `json:"PolicyType"`
		PsId               int         `json:"PsId"`
		GroupId            int         `json:"GroupId"`
		DeployTypeId       int         `json:"DeployTypeId"`
		DeployTypeName     string      `json:"DeployTypeName"`
		FromQty            int         `json:"FromQty"`
		ToQty              int         `json:"ToQty"`
		ServiceCode        int         `json:"ServiceCode"`
		Price              float64     `json:"Price"`
		PriceVAT           float64     `json:"PriceVAT"`
		PromotionId        int         `json:"PromotionId"`
		PromotionName      string      `json:"PromotionName"`
		IsBuilding         int         `json:"IsBuilding"`
	} `json:"SalesPolicyPackages"`
	ServiceId          int    `json:"ServiceId"`
	ServiceName        string `json:"ServiceName"`
	SubServiceTypeId   int    `json:"SubServiceTypeId"`
	SubServiceTypeName string `json:"SubServiceTypeName"`
	SubServiceId       int    `json:"SubServiceId"`
	SubServiceName     string `json:"SubServiceName"`
	DisplayName        string `json:"DisplayName"`
	Description        string `json:"Description"`
}
type SalePlatformCategoryDetail struct {
	Packages []SPCategoryPackages `json:"Packages"`
	Products []struct {
		IsDefault           bool `json:"IsDefault"`
		SalesPolicyProducts []struct {
			RevokeId      int         `json:"RevokeId"`
			RevokeName    string      `json:"RevokeName"`
			StatusId      int         `json:"StatusId"`
			StatusName    string      `json:"StatusName"`
			ParentName    string      `json:"ParentName"`
			ParentNameVN  string      `json:"ParentNameVN"`
			LanWire       interface{} `json:"LanWire"`
			UsesId        int         `json:"UsesId"`
			UsesName      interface{} `json:"UsesName"`
			ApplyServices []struct {
				ApplySubServiceId int `json:"ApplySubServiceId"`
				ApplyPrepaid      int `json:"ApplyPrepaid"`
			} `json:"ApplyServices"`
			IsDefault      int         `json:"IsDefault"`
			GroupNo        int         `json:"GroupNo"`
			DecideName     string      `json:"DecideName"`
			PolicyId       int         `json:"PolicyId"`
			PolicyCode     string      `json:"PolicyCode"`
			ProgramID      int         `json:"ProgramID"`
			PolicyName     string      `json:"PolicyName"`
			ComboType      int         `json:"ComboType"`
			PolicyType     int         `json:"PolicyType"`
			PsId           int         `json:"PsId"`
			GroupId        int         `json:"GroupId"`
			DeployTypeId   int         `json:"DeployTypeId"`
			DeployTypeName string      `json:"DeployTypeName"`
			FromQty        int         `json:"FromQty"`
			ToQty          int         `json:"ToQty"`
			ServiceCode    int         `json:"ServiceCode"`
			Price          float64     `json:"Price"`
			PriceVAT       float64     `json:"PriceVAT"`
			PromotionId    int         `json:"PromotionId"`
			PromotionName  interface{} `json:"PromotionName"`
			IsBuilding     int         `json:"IsBuilding"`
		} `json:"SalesPolicyProducts"`
		ServiceId          int    `json:"ServiceId"`
		ServiceName        string `json:"ServiceName"`
		SubServiceTypeId   int    `json:"SubServiceTypeId"`
		SubServiceTypeName string `json:"SubServiceTypeName"`
		SubServiceId       int    `json:"SubServiceId"`
		SubServiceName     string `json:"SubServiceName"`
		DisplayName        string `json:"DisplayName"`
		Description        string `json:"Description"`
	} `json:"Products"`
}

type RespSalePlatformCategoryDetail struct {
	Code    int `json:"Code"`
	Message struct {
		Message   string      `json:"Message"`
		ExMessage interface{} `json:"ExMessage"`
	} `json:"Message"`
	Data    SalePlatformCategoryDetail `json:"Data"`
	ErrorId string                     `json:"ErrorId"`
}

type ResultMyTickets struct {
	Id                int       `json:"event_id"`
	EventName         string    `json:"event_name"`
	EventImg          string    `json:"event_img"`
	EventLocation     string    `json:"event_location"`
	EventAddress      string    `json:"event_address"`
	EventStatus       string    `json:"event_status"`
	EventStartDate    time.Time `json:"event_start_date"`
	EventEndDate      time.Time `json:"event_end_date"`
	EventStartDateStr string    `json:"event_start_date_str"`
	State             int       `json:"state"`
	TransId           int       `json:"trans_id"`
	MyTicketType      string    `json:"my_ticket_type"`
	MyTicketStatus    string    `json:"my_ticket_status"`
	GroupId           int       `json:"group_id"`
	PaymentStatus     string    `json:"payment_status"`
	TCreate           time.Time `json:"t_create"`
	ProductId         int       `json:"product_id"`
}
type ResultMappingTransIdPaymentStatus struct {
	TransId       int
	PaymentStatus string
}
type ResultGetDetailMyTicketBelongtoOtherGroup struct {
	GroupName     string
	LeaderName    string
	LeaderPhone   string
	LeaderEmail   string
	FullName      string
	BibName       string
	Gender        int
	Birthday      string
	IdentityCards string
	Nationality   string
	City          string
	Email         string
	Phone         string
	TShirtSize    string
	RangeName     string
}

type RespListDataLocal struct {
	StatusCode int                 `json:"statusCode"`
	Message    string              `json:"message"`
	Data       []map[string]string `json:"data"`
}

type Icons struct {
	Id            int     `json:"id"`
	Title         string  `json:"title"`
	ActionType    string  `json:"actionType"`
	DataAction    string  `json:"dataAction"`
	IconUrl       string  `json:"iconUrl"`
	IsNew         int     `json:"isNew"`
	IsLocked      int     `json:"isLocked"`
	IsForceUpdate int     `json:"isForceUpdate"`
	Data          *string `json:"data"`
}
type RespLocalIconsWebkit struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Data       struct {
		Icons []Icons `json:"icons"`
	} `json:"data"`
}
type OutputMenuIcon struct {
	Id            int     `json:"id"`
	Title         string  `json:"title"`
	ActionType    string  `json:"actionType"`
	DataAction    string  `json:"dataAction"`
	IconUrl       string  `json:"iconUrl"`
	IsNew         int     `json:"isNew"`
	IsLocked      int     `json:"isLocked"`
	IsForceUpdate int     `json:"isForceUpdate"`
	Data          *string `json:"data"`
	Priority      int     `json:"Priority"`
}

type ItemLocalEcomCreateOrder struct {
	Title  string `json:"title"`
	Amount int    `json:"amount"`
}
type CustomerInfoOrder struct {
	Address  string `json:"address"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	ID       string `json:"id"`
	PhoneNb  string `json:"phone_nb"`
	Type     string `json:"type"`
}
type RuleInfoOrder struct {
	Replace string `json:"replace"`
	Text    string `json:"text"`
	URL     string `json:"url"`
}
type VatInfoOrder struct {
	AddressVat string `json:"address_vat"`
	ID         string `json:"id"`
	Note       string `json:"note"`
}
type SkuOrder struct {
	Amount   int    `json:"amount"`
	Quantity int    `json:"quantity"`
	Sku      string `json:"sku"`
	Title    string `json:"title"`
}
type DataLocalEcomCreateOrder struct {
	CustomerInfo      *CustomerInfoOrder         `json:"customer_info"`
	IsChangeAmount    int                        `json:"is_change_amount"`
	TransID           int                        `json:"trans_id"`
	TransType         string                     `json:"trans_type"`
	ServiceType       string                     `json:"service_type"`
	TotalPayment      int                        `json:"total_payment"`
	Vouchers          []string                   `json:"vouchers"`
	Items             []ItemLocalEcomCreateOrder `json:"items"`
	MetaData          string                     `json:"meta_data"`
	OrderID           string                     `json:"order_id"`
	PreDiscountAmount int                        `json:"pre_discount_amount"`
	PromoCode         string                     `json:"promo_code"`
	ReferralCode      *string                    `json:"referral_code"`
	RuleInfo          RuleInfoOrder              `json:"rule_info"`
	Skus              []SkuOrder                 `json:"skus"`
	StateType         string                     `json:"state_type"`
	TotalAmount       int                        `json:"total_amount"`
	VatInfo           *VatInfoOrder              `json:"vat_info"`
	VoucherDiscount   int                        `json:"voucher_discount"`
}
type RespLocalEcomCreateOrder struct {
	StatusCode int                       `json:"statusCode"`
	Message    string                    `json:"message"`
	Data       *DataLocalEcomCreateOrder `json:"data"`
}

type RespLocalAddress struct {
	StatusCode int               `json:"statusCode"`
	Message    string            `json:"message"`
	Data       *DataLocalAddress `json:"data"`
}
type DataLocalAddress struct {
	ID                     int    `json:"id"`
	FullName               string `json:"fullName"`
	Phone                  string `json:"phone"`
	Email                  string `json:"email"`
	Type                   string `json:"type"`
	IsDefault              int    `json:"isDefault"`
	LocationProvinceID     int    `json:"locationProvinceId"`
	LocationProvinceName   string `json:"locationProvinceName"`
	LocationProvinceNameEn string `json:"locationProvinceNameEn"`
	LocationDistrictID     int    `json:"locationDistrictId"`
	LocationDistrictName   string `json:"locationDistrictName"`
	LocationDistrictNameEn string `json:"locationDistrictNameEn"`
	LocationWardID         int    `json:"locationWardId"`
	LocationWardName       string `json:"locationWardName"`
	LocationWardNameEn     string `json:"locationWardNameEn"`
	LocationStreetID       int    `json:"locationStreetId"`
	StreetName             string `json:"streetName"`
	StreetNameEn           string `json:"streetNameEn"`
	ApartmentNumber        string `json:"apartmentNumber"`
}

type RespGetContractsBOT struct {
	StatusCode int        `json:"statusCode"`
	Message    string     `json:"message"`
	Data       []Contract `json:"data"`
}

type Contract struct {
	ContractNo       string `json:"contractNo"`
	ContractId       string `json:"contractId"`
	CustomerIsActive int    `json:"customerIsActive"`
	Address          string `json:"address"`
	LocationId       string `json:"LocationId"`
	IsInternet       int    `json:"isInternet"`
	IsFptPlay        int    `json:"isFptPlay"`
	IsCamera         int    `json:"isCamera"`
}
