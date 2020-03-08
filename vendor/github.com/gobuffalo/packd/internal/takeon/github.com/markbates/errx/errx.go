package errx

// go2 errors
type Wrapper interface {
	Unwrap() error
}

// pkg/errors
type Causer interface {
	Cause() error
}

func Unwrap(err error) error {
	switch e := err.(type) {
	case Wrapper:
		return e.Unwrap()
	case Causer:
		return e.Cause()
	}
	return err
}

var Cause = Unwrap
