package shared

type HttpHeaders struct {
	corsOrigins string
	corsMethods string
}

func NewHttpHeaders(corsOrigins string, corsMethods string) *HttpHeaders {
	return &HttpHeaders{corsOrigins: corsOrigins, corsMethods: corsMethods}
}

func (h *HttpHeaders) CreateHeaders() map[string]string {
	headers := make(map[string]string)
	headers["Access-Control-Allow-Origin"] = h.corsOrigins
	headers["Access-Control-Allow-Methods"] = h.corsMethods
	headers["Access-Control-Allow-Headers"] = "*"
	headers["Content-Type"] = "application/json"
	return headers
}
