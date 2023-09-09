package bot

var (
	UNoUser int = 0
	UAdmin  int = 1
	UUser   int = 2
)
var (
	NoRows = "sql: no rows in result set"
)

const (
	E_STATE_NOTHING             int = 0
	E_STATE_ADDCARD_ARTICLE         = 100
	E_STATE_ADDCARD_NAME            = 101
	E_STATE_ADDCARD_DESCRIPTION     = 102
	E_STATE_ADDCARD_PRICE           = 103
	E_STATE_ADDCARD_PHOTO           = 104
	E_STATE_ADDCARD_CATEGORY        = 105
	E_STATE_ADDCARD_SHOW            = 106
	E_STATE_ADMIN_MAX               = 499
)
