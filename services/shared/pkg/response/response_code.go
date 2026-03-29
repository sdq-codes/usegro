package response

type ResponseCode int

const (
	BAD_REQUEST       int = 9100
	VALIDATION_FAILED int = 760

	RESOURCE_CREATED   int = 7100
	RESOURCE_UPDATED   int = 7200
	RESOURCE_FETCHED   int = 7300
	RESOURCE_NOT_FOUND int = 7500
)
