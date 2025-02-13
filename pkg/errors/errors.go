package errors

const (
	// ===================-  COMMON  -===================
	ErrUnmarshalResponse = "ERR_FAILED_TO_DECODE_REQ"
	ErrMarshalResponse   = "ERR_FAILED_TO_ENCODE_JSON_RESP"
	// ===================-  AUTH  -===================
	ErrInvalidAuthReqParams = "ERR_INVALID_AUTH_REQ_PARAMS"
	ErrWrongPassword        = "ERR_WRONG_PASSWORD"
	ErrWrongPasswordFormat  = "ERR_WRONG_PASSWORD_FORMAT"
	ErrWrongUsernameFormat  = "ERR_WRONG_PASSWORD_FORMAT"
	ErrAuthHeader           = "ERR_AUTH_HEADER_IS_MISSING"
	ErrInvalidToken         = "ERR_INVALID_AUTH_TOKEN"
	ErrInvalidClaims        = "ERR_CANNOT_PARSE_CLAIMS"
	ErrLogin                = "ERR_FAILED_TO_LOGIN"
	// ===================-  INFO  -===================
	ErrGetInfo = "ERR_GET_INFO"
	// ===================-  BUY ITEM  -===================
	ErrInvalidGetBuyItemReqParams = "ERR_INVALID_GET_BUY_REQ_PARAMS"
	ErrGetBuyItem                 = "ERR_GET_BUY_ITEM"
	ErrItemDoesntExist            = "ERR_ITEM_DOESNT_EXIST"
	// ===================-  COINS  -===================
	ErrInvalidSendCoinsReqParams = "ERR_INVALID_SEND_COINS_REQ_PARAMS"
	ErrInvalidRecipient          = "ERR_RECIPIENT_DOESNT_EXIST"
	ErrInvalidRecipientYourself  = "ERR_RECIPIENT_IS_YOURSELF"
	ErrNotEnoughCoins            = "ERR_NOT_ENOUGH_COINS"
)
