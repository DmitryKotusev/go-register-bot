package errors

type InvalidCredentailsError struct {
	Message string
}

func (e InvalidCredentailsError) Error() string {
	return e.Message
}

type UnauthorizedError struct {
	Message string
}

func (e UnauthorizedError) Error() string {
	return e.Message
}

type ProceedingsCountError struct {
	Message string
}

func (e ProceedingsCountError) Error() string {
	return e.Message
}
