package common

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"
)

func BuildInfo() string {
	return "Build info: Version 1.0.0"
}

func AddIndent(input, indent string) string {
	lines := strings.Split(input, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = indent + line
		}
	}
	return strings.Join(lines, "\n")
}

// PrettyPrintYAML prints the data in pretty, indented YAML format with a message about the data type.
func PrettyPrintYAML(data interface{}) {
	var originalData interface{}

	// Check if the input is a string
	if reflect.TypeOf(data).Kind() == reflect.String {
		var jsonData interface{}
		jsonErr := json.Unmarshal([]byte(data.(string)), &jsonData)
		if jsonErr == nil {
			originalData = jsonData
		} else {
			var yamlData interface{}
			yamlErr := yaml.Unmarshal([]byte(data.(string)), &yamlData)
			if yamlErr == nil {
				originalData = yamlData
			} else {
				log.Printf("Error unmarshaling data as JSON or YAML: %s, %s", jsonErr, yamlErr)
				return
			}
		}
	} else {
		originalData = data
	}
	// Print the message about the data type and YAML format
	fmt.Printf("Printing data of type %s in YAML format:\n", reflect.TypeOf(originalData).Kind())

	// Pretty print the data as YAML
	yamlBytes, err := yaml.Marshal(originalData)
	if err != nil {
		log.Printf("Error pretty-printing YAML: %s", err)
		return
	}
	indentedYAML := AddIndent(string(yamlBytes), "  ")
	fmt.Printf("%s\n", indentedYAML)
}

func PrettyPrint(data interface{}) {
	PrettyPrintIndent(data, 0)
}

func PrettyPrintIndent(data interface{}, indentLevel int) {
	indent := strings.Repeat("  ", indentLevel)
	v := reflect.ValueOf(data)

	switch v.Kind() {
	case reflect.Struct:
		fmt.Printf("%sStruct:\n", indent)
		for i := 0; i < v.NumField(); i++ {
			fieldName := v.Type().Field(i).Name
			fieldValue := v.Field(i).Interface()
			fmt.Printf("%s  %s: ", indent, fieldName)
			PrettyPrintIndent(fieldValue, indentLevel+1)
		}
	case reflect.Map:
		fmt.Printf("%sMap:\n", indent)
		for _, key := range v.MapKeys() {
			fmt.Printf("%s  %v: ", indent, key.Interface())
			PrettyPrintIndent(v.MapIndex(key).Interface(), indentLevel+1)
		}
	case reflect.Slice, reflect.Array:
		fmt.Printf("%sArray/Slice:\n", indent)
		for i := 0; i < v.Len(); i++ {
			fmt.Printf("%s  [%d]: ", indent, i)
			PrettyPrintIndent(v.Index(i).Interface(), indentLevel+1)
		}
	default:
		fmt.Printf("%s%v\n", indent, v.Interface())
	}
}
