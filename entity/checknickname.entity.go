package entity

type CheckNicknameResponse struct {
	IsSuccess bool   `json:"isSuccess"`
	Username  string `json:"username"`
	UserId    string `json:"userId"`
	ZoneId    string `json:"zoneId"`
}

type Response struct {
	InitCallBackendAPI   bool               `json:"initCallBackendAPI"`
	ErrorCode            string             `json:"errorCode"`
	Confirmation         bool               `json:"confirmation"`
	IsUserConfirmation   bool               `json:"isUserConfirmation"`
	ErrorMsg             string             `json:"errorMsg"`
	PaymentChannel       string             `json:"paymentChannel"`
	Result               string             `json:"result"`
	ChannelPrice         string             `json:"channelPrice"`
	ConfirmationFields   ConfirmationFields `json:"confirmationFields"`
	Success              bool               `json:"success"`
	Denom                string             `json:"denom"`
	User                 User               `json:"user"`
	IsThirdPartyMerchant bool               `json:"isThirdPartyMerchant"`
	TxnId                string             `json:"txnId"`
}

type Role struct {
	Client_type    string `json:"client_type"`
	Packed_role_id string `json:"packed_role_id"`
	Role           string `json:"role"`
	Role_id        string `json:"role_id"`
	Server         string `json:"server"`
	Server_id      string `json:"server_id"`
}

type ConfirmationFields struct {
	Roles               []Role `json:"roles,omitempty"`
	ZipCode             string `json:"zipCode"`
	Country             string `json:"country"`
	TotalPrice          string `json:"totalPrice"`
	Create_role_country string `json:"create_role_country"`
	UserIdAndZoneId     string `json:"userIdAndZoneId"`
	UserId              string `json:"userId"`
	ProductName         string `json:"productName"`
	PaymentChannel      string `json:"paymentChannel"`
	This_login_country  string `json:"this_login_country"`
	ChannelPrice        string `json:"channelPrice"`
	ZoneId              string `json:"zoneId"`
	TaxAmount           string `json:"taxAmount"`
	Email               string `json:"email"`
	InputRoleId         string `json:"inputRoleId"`
	Username            string `json:"username"`
}

type User struct {
	Target string `json:"target"`
	Secret string `json:"secret"`
}

type MobapayResponse struct {
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Data    MobapayDataResponse `json:"data"`
}

type SepulsaResponse struct {
	Status  string              `json:"status"`
	Rescode string              `json:"rescode"`
	Data    SepulsaDataResponse `json:"data"`
}

type SepulsaDataResponse struct {
	Id                  int                       `json:"id"`
	Status              string                    `json:"status"`
	Url                 string                    `json:"url"`
	Lines               []SepulsaLineDataResponse `json:"lines"`
	Total_excl_fee      string                    `json:"total_excl_fee"`
	Total_incl_fee      string                    `json:"total_incl_fee"`
	Contains_a_voucher  bool                      `json:"contains_a_voucher"`
	Contains_a_merchant bool                      `json:"contains_a_merchant"`
	Merchant_label      string                    `json:"merchant_label"`
	Merchant_point      int                       `json:"merchant_point"`
	Meowth_error        bool                      `json:"meowth_error"`
}

type SepulsaLineDataResponse struct {
	Attributes []struct {
		Option struct {
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"option"`
		Value string `json:"value"`
	} `json:"attributes"`
}

type MobapayDataResponse struct {
	User_info MobapayDataUserInfoResponse `json:"user_info"`
}

type MobapayDataUserInfoResponse struct {
	Code      int    `json:"code"`
	User_name string `json:"user_name"`
}

type CodashopResponse struct {
	InitCallBackendAPI   bool                               `json:"initCallBackendAPI"`
	Confirmation         bool                               `json:"confirmation"`
	IsUserConfirmation   bool                               `json:"isUserConfirmation"`
	ErrorMsg             string                             `json:"errorMsg"`
	PaymentChannel       string                             `json:"paymentChannel"`
	Result               string                             `json:"result"`
	ChannelPrice         string                             `json:"channelPrice"`
	ConfirmationFields   CodashopConfirmationFieldsResponse `json:"confirmationFields"`
	Success              bool                               `json:"success"`
	Denom                string                             `json:"denom"`
	IsThirdPartyMerchant bool                               `json:"isThirdPartyMerchant"`
	TxnId                string                             `json:"txnId"`
}

type CodashopConfirmationFieldsResponse struct {
	Roles []struct {
		Client_type    string `json:"client_type"`
		Packed_role_id string `json:"packed_role_id"`
		Role           string `json:"role"`
		Role_id        string `json:"role_id"`
		Server         string `json:"server"`
		Server_id      string `json:"server_id"`
	} `json:"code"`
	ZipCode             string `json:"zipCode"`
	Country             string `json:"country"`
	TotalPrice          string `json:"totalPrice"`
	Create_role_country string `json:"create_role_country"`
	UserIdAndZoneId     string `json:"userIdAndZoneId"`
	ProductName         string `json:"productName"`
	PaymentChannel      string `json:"paymentChannel"`
	This_login_country  string `json:"this_login_country"`
	ChannelPrice        string `json:"channelPrice"`
	ZoneId              string `json:"zoneId"`
	TaxAmount           string `json:"taxAmount"`
	Email               string `json:"email"`
	InputRoleId         string `json:"inputRoleId"`
	Username            string `json:"username"`
}
