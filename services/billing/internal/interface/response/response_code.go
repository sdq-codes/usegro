package response

type ResponseCode int

const (
	//SYSTEM_OPERATION_SUCCESS ResponseCode = 0
	//
	//// 7xx client errors
	BAD_REQUEST int = 9100
	//UNAUTHORIZED       ResponseCode = 701
	//FORBIDDEN          ResponseCode = 703
	//NOT_FOUND          ResponseCode = 704
	//METHOD_NOT_ALLOWED ResponseCode = 705
	VALIDATION_FAILED int = 760
	//
	//// 8xx server errors
	//INTERNAL_SERVER_ERROR ResponseCode = 800
	//NOT_IMPLEMENTED       ResponseCode = 801
	//SERVUCE_UNAVAILABLE   ResponseCode = 803
	//
	//// custom errors
	//BACKGROUND_JOB_FAILED ResponseCode = 810

	RESOURCE_CREATED      int = 7100
	RESOURCE_UPDATED      int = 7200
	RESOURCE_FETCHED      int = 7300
	RESOURCE_NOT_FOUND    int = 7500
	RESOURCE_DELETED      int = 7600
	INVALID_REQUEST       int = 9200
	INTERNAL_SERVER_ERROR int = 9300
)
