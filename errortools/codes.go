package errortools

import "net/http"

const (
	InternalServerErrCode = "ECM500"
)

const (
	ValidationFailure = "ECM400"
	BindingError      = "ECM401"
	MissingQueryParam = "ECM402"
	MinLengthError    = "ECM403"
	MissingField      = "ECM404"
	BadFormat         = "ECM405"
	DataFetchFail     = "ECM406"
	JSONSyntaxError   = "ECM407"
	NoRecord          = "ECM408"
	InvalidFieldType  = "ECM409"
	NotFound          = "ECM410"
	Invalid           = "ECM411"
	InvalidUUID       = "ECM412"
	Inactive          = "ECM413"
	HasActiveProduct  = "ECM414"
	HasVariants       = "ECM415"
	Required          = "ECM416"
)

const (
	UnauthorizedAccess = "ECM600"
)

var errorCodeNames = map[string]string{
	InternalServerErrCode: "internal server error",
	MissingQueryParam:     "mandatory to pass the query param",
	ValidationFailure:     "Validation error",
	MinLengthError:        "minimum length should be %v characters",
	BindingError: `The request parameter was not bound correctly.
	please check the input params and try again`,
	MissingField:       "%s is mandatory.",
	BadFormat:          "Please check the format of the field",
	DataFetchFail:      "Couldn't get the  requested details. Please check the inputs",
	JSONSyntaxError:    "JSON syntax error at byte offset",
	NoRecord:           "No records found",
	UnauthorizedAccess: "Unauthorized access",
	NotFound:           "Project not found",
	Invalid:            "Invalid %s",
	InvalidUUID:        "Invalid value for uuid type for %s",
	Inactive:           "%s is already inactive",
	HasActiveProduct:   "Invalid operation, has active products associated with it.",
	HasVariants:        "Invalid operation, has active variants associated with it.",
	Required:           "Require the %s",
}

var errorHTTPCodes = map[string]int{
	InternalServerErrCode: http.StatusInternalServerError,
	ValidationFailure:     http.StatusBadRequest,
	UnauthorizedAccess:    http.StatusUnauthorized,
	NotFound:              http.StatusNotFound,
}
