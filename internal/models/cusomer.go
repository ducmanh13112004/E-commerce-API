package models

type CustomerSopRequest struct {
	Name          string  `json:"name"`
	Phone         string  `json:"phone"`
	CustomerEmail *string `json:"customer_email,omitempty"`
	Dob           string  `json:"dob"`
	Address       *string `json:"address,omitempty"`
	StaffEmail    string  `json:"staff_email"`
	ContractCode  *string `json:"contract_code,omitempty"`
	Gender        *string `json:"gender"`
	OrderPrice    *int    `json:"order_price,omitempty"`
}

type CustomerSopResponse struct {
	Status    string      `json:"status"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`       // chỉ có khi thành công
	ErrorCode interface{} `json:"error_code,omitempty"` // chỉ có khi thất bại
	Detail    string      `json:"detail,omitempty"`     // chỉ có khi thất bại
}

type CustomerLocationResponse struct {
	Data struct {
		BookAddressID interface{} `json:"book_address_id"`
		ContractInfo  struct {
			ContractID int    `json:"ContractId"`
			ContractNo string `json:"ContractNo"`
			IsInternet int    `json:"IsInternet"`
		} `json:"contract_info"`
		LocationInfo struct {
			AddressNo    string `json:"address_no"`
			BuildingID   int    `json:"building_id"`
			CusTypeID    int    `json:"cus_type_id"`
			CusTypeL2ID  int    `json:"cus_type_l2_id"`
			DistrictID   int    `json:"district_id"`
			DistrictName string `json:"district_name"`
			FullAddress  string `json:"full_address"`
			LocationID   int    `json:"location_id"`
			LocationName string `json:"location_name"`
			ObjectTypeID int    `json:"object_type_id"`
			StreetID     int    `json:"street_id"`
			StreetName   string `json:"street_name"`
			WardID       int    `json:"ward_id"`
			WardName     string `json:"ward_name"`
		} `json:"location_info"`
	} `json:"data"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}
