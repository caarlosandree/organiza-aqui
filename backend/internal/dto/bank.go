package dto

// BankDTO representa um banco na camada de apresentação
type BankDTO struct {
	ID       string `json:"id"`
	ISPB     string `json:"ispb"`
	Code     int    `json:"code"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
}

// BrasilAPIBankResponse representa a resposta da API BrasilAPI
type BrasilAPIBankResponse struct {
	ISPB     string `json:"ispb"`
	Name     string `json:"name"`
	Code     int    `json:"code"`
	FullName string `json:"fullName"` // A API retorna "fullName" em camelCase
}
