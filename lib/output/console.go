package output

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
)

type ConsoleOutput struct{}

func NewConsoleOutput() *ConsoleOutput {
	return &ConsoleOutput{}
}

func (o *ConsoleOutput) Write(result map[string]interface{}, errs map[string]error) error {
	for name, res := range result {
		if res == nil {
			logrus.WithField("name", name).Debug("Scanner returned result <nil>")
			continue
		}
		fmt.Printf("\nResults for %s\n", name)
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
			fmt.Printf("%s: %v\n", field, fmt.Sprintf("%v", v))
		}
	}

	if len(errs) > 0 {
		fmt.Println("\nThe following scanners returned errors:")
	}
	for name, err := range errs {
		fmt.Printf("%s: %s\n", name, err)
	}

	fmt.Printf("\n%d scanner(s) succeeded\n", len(result))

	return nil
}

func (o *ConsoleOutput) unmarshal(res interface{}) string {
	output := ""
	typeOf := reflect.TypeOf(res)
	for i := 0; i < typeOf.NumField(); i++ {
		field, ok := typeOf.Field(i).Tag.Lookup("console")
		if !ok || field == "-" {
			logrus.WithFields(map[string]interface{}{
				"found": ok,
				"value": field,
			}).Debug("Console field was ignored")
			continue
		}
		v := reflect.ValueOf(res).FieldByName(typeOf.Field(i).Name)
		output += fmt.Sprintf("%s: %v\n", field, fmt.Sprintf("%v", v))
	}
	return output
}
