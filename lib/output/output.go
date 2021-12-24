package output

type Output interface {
	Write(map[string]interface{}, map[string]error) error
}

type OutputKey int

const (
	Console OutputKey = iota + 1
)

func GetOutput(o OutputKey) Output {
	switch o {
	case Console:
		return NewConsoleOutput()
	}
	return nil
}
