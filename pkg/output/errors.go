package output

type OpenFileError struct {
	err error
}

func (e OpenFileError) Error() string {
	return e.err.Error()
}

func NewOpenFileError(err error) OpenFileError {
	return OpenFileError{err: err}
}
