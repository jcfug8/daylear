package fileretriever

// ErrInvalidArgument is an error that occurs when the fileretriever is given an invalid argument.
type ErrInvalidArgument struct {
	Msg string
}

func (e ErrInvalidArgument) Error() string {
	return e.Msg
}
