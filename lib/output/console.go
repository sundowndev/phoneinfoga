package output

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"io"
	"reflect"
)

type ConsoleOutput struct {
	w io.Writer
}

func NewConsoleOutput(w io.Writer) *ConsoleOutput {
	return &ConsoleOutput{w: w}
}

func (o *ConsoleOutput) Write(result map[string]interface{}, errs map[string]error) error {
	for name, res := range result {
		if res == nil {
			logrus.WithField("name", name).Debug("Scanner returned result <nil>")
			continue
		}
		_, _ = fmt.Fprintf(o.w, color.WhiteString("Results for %s\n"), name)
		typeOf := reflect.TypeOf(res)
		for i := 0; i < typeOf.NumField(); i++ {
			v := reflect.ValueOf(res).FieldByName(typeOf.Field(i).Name)
			field, ok := typeOf.Field(i).Tag.Lookup("console")
			if !ok || field == "-" {
				logrus.WithFields(map[string]interface{}{
					"found": ok,
					"value": field,
				}).Debug("Console field was ignored")
				continue
			}
			_, _ = fmt.Fprintf(o.w, "%s: %v\n", field, fmt.Sprintf("%v", v))
		}
		_, _ = fmt.Fprintf(o.w, "\n")
	}

	if len(errs) > 0 {
		_, _ = fmt.Fprintln(o.w, "The following scanners returned errors:")
		for name, err := range errs {
			_, _ = fmt.Fprintf(o.w, "%s: %s\n", name, err)
		}
		_, _ = fmt.Fprintf(o.w, "\n")
	}

	_, _ = fmt.Fprintf(o.w, "%d scanner(s) succeeded\n", len(result))

	return nil
}
