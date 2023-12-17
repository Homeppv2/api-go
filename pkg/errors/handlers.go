package errors

func HandleServiceError(err error) error {
	switch err {
	case ErrUserNotFound:
		return NewError(ErrUserNotFound.Error(), ErrCodeNotFound)
	case ErrUserEmailAlreadyExist:
		return NewError(ErrUserEmailAlreadyExist.Error(), ErrCodeInvalidArgument)
	case ErrSessionNotFound:
		return NewError(ErrSessionNotFound.Error(), ErrCodeNotFound)
	case ErrControllerNotFound:
		return NewError(ErrControllerNotFound.Error(), ErrCodeNotFound)
	case ErrControllerAlreadyIsUsed:
		return NewError(ErrControllerAlreadyIsUsed.Error(), ErrCodeInvalidArgument)
	case ErrControllerAlreadyExist:
		return NewError(ErrControllerAlreadyExist.Error(), ErrCodeInvalidArgument)
	case ErrInvalidLogin:
		return NewError(ErrInvalidLogin.Error(), ErrCodeInvalidArgument)
	case ErrValidationError:
		return NewError(ErrValidationError.Error(), ErrCodeInvalidArgument)
	default:
		return NewError(ErrUnknown.Error(), ErrCodeUnknown)
	}
}
