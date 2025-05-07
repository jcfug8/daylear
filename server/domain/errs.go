package domain

type ErrInvalidArgument struct {
	Msg string
}

func (e ErrInvalidArgument) Error() string {
	return e.Msg
}

type ErrPermissionDenied struct {
	Msg string
}

func (e ErrPermissionDenied) Error() string {
	return e.Msg
}
