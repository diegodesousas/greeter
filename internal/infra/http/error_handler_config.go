package http

var ErrorHandler = NewErrorHandler(
	ErrorHandlerNotFound{},
	ErrorHandlerBadRequest{},
	ErrorHandlerServiceUnavailable{},
	ErrorHandlerValidation{},
)
