package main

import (
	"github.com/synkube/app/core/common"
)

func main() {
	jsonString := `{
        "name": "John Doe",
        "age": 30,
        "address": {
            "city": "San Francisco",
            "state": "CA",
            "country": "USA"
        },
        "hobbies": ["hiking", "swimming", "reading"]
    }`

	yamlString := `
name: Jane Doe
age: 25
address:
  city: Los Angeles
  state: CA
  country: USA
hobbies:
  - dancing
  - cooking
  - travelling
`

	jsonData := map[string]interface{}{
		"name": "Jack Doe",
		"age":  35,
		"address": map[string]string{
			"city":    "New York",
			"state":   "NY",
			"country": "USA",
		},
		"hobbies": []string{"writing", "jogging", "photography"},
	}

	// Pretty print from JSON string
	common.PrettyPrintYAML(jsonString)

	// Pretty print from YAML string
	common.PrettyPrintYAML(yamlString)

	// Pretty print from struct
	common.PrettyPrintYAML(jsonData)
}
