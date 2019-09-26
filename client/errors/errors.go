package errors

const DEFAULT_LANG = 1
const DEV = 0
const EN = 1

const INPUT_VALIDATION_FAILED    = 1
const EMPTY_REQUEST_NOT_ALLOWED  = 2
const MAX_REQUESTS_EXCEEDED      = 3
const FAILED_DUE_TO_OTHER_ERRORS = 4
const INTERNAL_SERVER_ERROR      = 5

// keep em' static
var ERRORS = map[int]map[int]string{
  INPUT_VALIDATION_FAILED:
    {
      EN:  "Input validation failed",
      DEV: "Struct validations failed on tags for input",
    },
  EMPTY_REQUEST_NOT_ALLOWED:
    {
      EN:  "Empty request not allowed"
      DEV: "This endpoint does not allow the empty request - each request must be defined separately",
    },
  MAX_REQUESTS_EXCEEDED:
    {
      EN:  "Max number of requests exceeded",
      DEV: "MaxRequest parameter has been set for endpoint and is exceeded by the number of request-objects given in the input",
    },
  FAILED_DUE_TO_OTHER_ERRORS:
    {
      EN:  "Failed due to other errors",
      DEV: "Other request has already been invalidated, no reason to continue until those have been fixed",
    },
  INTERNAL_SERVER_ERROR:
    {
      EN:  "Internal server error occured. Please wait until it has been fixed, before you try again",
      DEV: "Internal server error occured. Please wait until it has been fixed, before you try again",
    },
}

// E[INPUT_VALIDATION_FAILED][EN]
// client.E[client.INPUT_VALIDATION_FAILED][client.EN]
