package config

type Body struct {
	ResponseCode int
	Message      string
}

type LDMResponse struct {
	InfaToken string `json:"infaToken"`
	Error     string `json:"error"`
}

type LDMHeader struct {
	InfaToken  string
	JsessionID string
}
