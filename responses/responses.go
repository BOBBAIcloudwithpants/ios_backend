package responses

type data struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type Token struct {
	Token string `json:"token"`
}

type StatusOKResponse struct {
	Code int `json:"code" example:"200"`
	data
}

type StatusBadRequestResponse struct {
	Code int `json:"code" example:"400"`
	data
}

type StatusForbiddenResponse struct {
	Code int `json:"code" example:"403"`
	data
}

type StatusNotFoundResponse struct {
	Code int `json:"code" example:"404"`
	data
}

type StatusInternalServerError struct {
	Code int `json:"code" example:"500"`
	data
}
