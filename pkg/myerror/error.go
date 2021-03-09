package myerror

type ApplicationError struct {
	Message       string
	OriginalError error
	Code          int
}

func (e ApplicationError) Error() string {
	if e.OriginalError != nil {
		return e.Message + ": " + e.OriginalError.Error()
	} else {
		return e.Message
	}
}
