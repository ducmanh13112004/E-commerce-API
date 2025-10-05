package models

type RespLocalAuth struct {
	StatusCode int        `json:"statusCode"`
	Message    string     `json:"message"`
	Data       RedisModel `json:"data"`
}

type RedisContract struct {
	Address                 string `json:"address"`
	BranchCode              string `json:"branchCode"`
	BranchName              string `json:"branchName"`
	CheckEContractConfirmed int    `json:"checkEContractConfirmed"`
	ContractId              string `json:"contractId"`
	ContractNo              string `json:"contractNo"`
	ContractPhone           string `json:"contractPhone"`
	CustomerIsActive        int    `json:"customerIsActive"`
	CustomerSetContractName string `json:"customerSetContractName"`
	Email                   string `json:"email"`
	FullName                string `json:"fullName"`
	IsBill                  string `json:"isBill"`
	IsCamera                int    `json:"isCamera"`
	IsFptEmployee           int    `json:"isFptEmployee"`
	IsFptPlay               int    `json:"isFptPlay"`
	IsFsafe                 int    `json:"isFsafe"`
	IsInternet              int    `json:"isInternet"`
	IsNetTV                 string `json:"isNetTV"`
	IsOwner                 int    `json:"isOwner"`
	LastTimeServiceFeedback string `json:"lastTimeServiceFeedback"`
	LocalType               string `json:"localType"`
	LocationCode            string `json:"locationCode"`
	LocationId              string `json:"locationId"`
	LocationZone            string `json:"locationZone"`
	Order                   int    `json:"order"`
	Passport                string `json:"passport"`
	PayTv                   string `json:"payTv"`
	Status                  string `json:"status"`
	Traffics                struct {
		Download string `json:"download"`
		Upload   string `json:"upload"`
	} `json:"traffics"`
	Type    string `json:"type"`
	TypeVip string `json:"typeVip"`
}

type RedisModel struct {
	ClientId              string          `json:"clientId"`
	CodeVerifier          string          `json:"codeVerifier"`
	CreateByGrant         int             `json:"createByGrant"`
	CustomerId            int             `json:"customerId"`
	DeviceId              string          `json:"deviceId"`
	DeviceInfo            string          `json:"deviceInfo"`
	DeviceName            string          `json:"deviceName"`
	DevicePlatform        string          `json:"devicePlatform"`
	DeviceToken           string          `json:"deviceToken"`
	IsActivePhone         int             `json:"isActivePhone"`
	IsCanhTo              int             `json:"isCanhTo"`
	IsCustomerFpt         int             `json:"isCustomerFpt"`
	IsPassword            int             `json:"isPassword"`
	JwtSecret             string          `json:"jwtSecret"`
	ListContract          []RedisContract `json:"listContract"`
	Phone                 string          `json:"phone"`
	Provider              string          `json:"provider"`
	ProviderAccessToken   string          `json:"providerAccessToken"`
	ProviderId            string          `json:"providerId"`
	AppVersion            string          `json:"appVersion"`
	IsNeedActiveEContract int             `json:"isNeedActiveEContract"`
	EContractContractNo   *string         `json:"eContractContractNo"`
}

type PortalPayload struct {
	UserID   string `json:"user_id"`
	UserCode string `json:"user_code"`
	PhoneNb  string `json:"phone_nb"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	// UpdateBy       string `json:"update_by"`
	Company        string `json:"company"`
	CompanyName    string `json:"company_name"`
	Department     string `json:"department"`
	DepartmentName string `json:"department_name"`
	EmpCode        string `json:"emp_code"`
	EmpID          string `json:"emp_id"`
	GroupID        string `json:"group_id"`
	Exp            int64  `json:"exp"`
	Expired        string `json:"expired"`
}
