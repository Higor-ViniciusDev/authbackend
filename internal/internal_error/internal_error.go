package internal_error

type InternalError struct {
	Message string
	Err     string
}

func (ie *InternalError) Error() string {
	return ie.Message + ": " + ie.Err
}

func NewNotFoundError(message string) *InternalError {
	return &InternalError{
		Message: message,
		Err:     "not_found",
	}
}

func NewInternalServerError(message string) *InternalError {
	return &InternalError{
		Message: message,
		Err:     "internal_server_error",
	}
}

func NewBadRequestError(message string) *InternalError {
	return &InternalError{
		Message: message,
		Err:     "bad_request",
	}
}

func NewManyRequestError(message string) *InternalError {
	return &InternalError{
		Message: message,
		Err:     "many_request",
	}
}

func NewUnauthorizedAccess(menssage string) *InternalError {
	return &InternalError{
		Message: menssage,
		Err:     "unauthorized",
	}
}

func NewUnauthorizedEmailAlreadyExists(message string) *InternalError {
	return &InternalError{
		Message: message,
		Err:     "unauthorized_email_already_exists",
	}
}

func NewUnauthorizedEmailNotVerified(message string) *InternalError {
	return &InternalError{
		Message: message,
		Err:     "unauthorized_email_not_verified",
	}
}
