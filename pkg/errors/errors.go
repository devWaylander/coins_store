package errors

const (
	// ===================-   COMMON   -===================
	ErrUnmarshalResponse = "ERR_FAILED_TO_DECODE_REQ"
	ErrMarshalResponse   = "ERR_FAILED_TO_ENCODE_JSON_RESP"
	// ===================-    AUTH    -===================
	ErrInvalidAuthReqParams = "ERR_INVALID_AUTH_REQ_PARAMS"
	ErrAuthHeader           = "ERR_AUTH_HEADER_IS_MISSING"
	ErrInvalidToken         = "ERR_INVALID_AUTH_TOKEN"
	ErrInvalidClaims        = "ERR_CANNOT_PARSE_CLAIMS"
)
