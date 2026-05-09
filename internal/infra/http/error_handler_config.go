package http

var ErrorHandler = NewErrorHandler(
	ErrorHandlerNotFound{},
	ErrorHandlerValidation{},
)
