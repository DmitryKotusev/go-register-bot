package errors

type InvalidCredentailsError struct {
	Message string
}

func (e InvalidCredentailsError) Error() string {
	return e.Message
}
