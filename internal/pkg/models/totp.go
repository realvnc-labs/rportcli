package models

type TotPSecretOutput struct {
	Secret  string `json:"secret"`
	Comment string `json:"comment"`
	File    string `json:"file"`
}

type TotPSecretResp struct {
	Secret        string `json:"secret"`
	QRImageBase64 string `json:"qr"`
}

func (tps *TotPSecretOutput) Headers() []string {
	return []string{
		"SECRET",
		"QR FILE",
	}
}

func (tps *TotPSecretOutput) Row() []string {
	return []string{
		tps.Secret,
		tps.File,
	}
}
