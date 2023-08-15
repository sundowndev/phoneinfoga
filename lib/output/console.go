package output

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"io"
	"reflect"
	"sort"
	"strings"
)

type ConsoleOutput struct {
	w io.Writer
}

func NewConsoleOutput(w io.Writer) *ConsoleOutput {
	return &ConsoleOutput{w: w}
}

func (o *ConsoleOutput) Write(result map[string]interface{}, errs map[string]error) error {
	succeeded := 0
	for _, name := range getSortedResultKeys(result) {
		res := result[name]
		if res == nil {
			logrus.WithField("name", name).Debug("Scanner returned result <nil>")
			continue
		}
		_, _ = fmt.Fprintf(o.w, color.WhiteString("Results for %s\n"), name)
		o.displayResult(res, "")
		_, _ = fmt.Fprintf(o.w, "\n")
		succeeded++
	}

	if len(errs) > 0 {
		_, _ = fmt.Fprintln(o.w, "The following scanners returned errors:")
		for _, name := range getSortedErrorKeys(errs) {
			_, _ = fmt.Fprintf(o.w, "%s: %s\n", name, errs[name])
		}
		_, _ = fmt.Fprintf(o.w, "\n")
	}

	_, _ = fmt.Fprintf(o.w, "%d scanner(s) succeeded\n", succeeded)

	return nil
}

func (o *ConsoleOutput) displayResult(val interface{}, prefix string) {
	reflectType := reflect.TypeOf(val)
	reflectValue := reflect.ValueOf(val)

	if reflectValue.Kind() == reflect.Slice {
		for i := 0; i < reflectValue.Len(); i++ {
			item := reflectValue.Index(i)
			if item.Kind() == reflect.Ptr {
				item = reflectValue.Index(i).Elem()
			}
			o.displayResult(item.Interface(), prefix)

			// If it's the latest element, add a newline
			if i < reflectValue.Len()-1 {
				_, _ = fmt.Fprintf(o.w, "\n")
			}
		}
		return
	}

	for i := 0; i < reflectType.NumField(); i++ {
		valueValue := reflectValue.Field(i).Interface()

		field, ok := reflectType.Field(i).Tag.Lookup("console")
		if !ok || field == "-" {
			continue
		}

		if strings.Contains(field, "omitempty") && reflectValue.Field(i).IsZero() {
			continue
		}
		fieldTitle := strings.Split(field, ",")[0]

		switch reflectValue.Field(i).Kind() {
		case reflect.String:
			_, _ = fmt.Fprintf(o.w, "%s%s: ", prefix, fieldTitle)
			_, _ = fmt.Fprintf(o.w, color.YellowString("%s\n"), valueValue)
		case reflect.Bool:
			_, _ = fmt.Fprintf(o.w, "%s%s: ", prefix, fieldTitle)
			_, _ = fmt.Fprintf(o.w, color.YellowString("%v\n"), valueValue)
		case reflect.Int:
			_, _ = fmt.Fprintf(o.w, "%s%s: ", prefix, fieldTitle)
			_, _ = fmt.Fprintf(o.w, color.YellowString("%d\n"), valueValue)
		case reflect.Struct:
			_, _ = fmt.Fprintf(o.w, "%s%s:\n", prefix, fieldTitle)
			o.displayResult(valueValue, prefix+"\t")
		case reflect.Slice:
			_, _ = fmt.Fprintf(o.w, color.WhiteString("%s:\n"), fieldTitle)
			o.displayResult(valueValue, prefix+"\t")
		}
	}
}

func getSortedResultKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func getSortedErrorKeys(m map[string]error) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
