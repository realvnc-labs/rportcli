package models

type TwoFA struct {
	SentTo         string `json:"send_to"`
	DeliveryMethod string `json:"delivery_method"`
	TotPKeyStatus  string `json:"totp_key_status"`
}

type Token struct {
	Token string `json:"token"`
	TwoFA TwoFA  `json:"two_fa,omitempty"`
}
