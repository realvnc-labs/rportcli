package models

type TwoFA struct {
	SentTo         string `json:"send_to"`
	DeliveryMethod string `json:"delivery_method"`
}

type Token struct {
	Token string `json:"token"`
	TwoFA TwoFA  `json:"two_fa,omitempty"`
}
