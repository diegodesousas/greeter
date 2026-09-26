package http

var ErrorHandler = NewErrorHandler(
	ErrorHandlerNotFound{},
	ErrorHandlerBadRequest{},
	ErrorHandlerValidation{},
)
