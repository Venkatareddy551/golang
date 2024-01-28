package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"regexp"
)

func isEnclosedPlaceholder(input string) bool {
	prefix := "{{"
	suffix := "}}"
	return strings.HasPrefix(input, prefix) && strings.HasSuffix(input, suffix)
}

var placeholderRegex = regexp.MustCompile(`{{([^}]+)}}`)


// ResolveExpression resolves the expression by traversing the JSON object
func ResolveExpression(jsonData map[string]interface{}, expression string) (interface{}, error) {
	if isEnclosedPlaceholder(expression) {
		// Split the expression into parts
		expressions := placeholderRegex.ReplaceAllString(expression, "$1")
		parts := strings.Split(expressions, ".")

		// Traverse the JSON structure based on the parts of the expressions
		current := jsonData
		for _, part := range parts {
			val, ok := current[part]
			if !ok {
				return nil, fmt.Errorf("key %s not found", part)
			}

			// Check if the value is a nested JSON object
			if nested, nestedOk := val.(map[string]interface{}); nestedOk {
				current = nested
			} else {
				// If the value is not a nested object, it is the final value
				return val, nil
			}
		}

		return nil, fmt.Errorf("expressions could not be fully resolved")
	} else {
		return expression, nil
	}
}

func addAToValues(dataa interface{}) interface{} {
    // Example JSON data
	jsonData := `{
		"RequestBody": {
			"data": "resolved_value"
		},
		"ResponseBody": {
		    "message":"other_value"
		}
	}`

	// Unmarshal JSON data into a map
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		fmt.Println("Error unmarshalling JSON data:", err)
		//return
	}
	switch v := dataa.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, value := range v {
			result[key] = addAToValues(value)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, value := range v {
			result[i] = addAToValues(value)
		}
		return result
	case string:
		fmt.Println("--->>>:", isEnclosedPlaceholder(v))
		//return ResolveExpression(data, v)
		//return placeholderRegex.ReplaceAllString(v, "$1")
		// Resolve the expression
		results, err := ResolveExpression(data, v)
		if err != nil {
			fmt.Println("Error resolving expression:", err)
			return nil
		}

		// Print the resolved value
		fmt.Printf("Resolved value for expression '%s': %v\n", v, results)
		return results
	default:
		return v
	}
}

func main() {
	// Read JSON file
	fileData, err := ioutil.ReadFile("updateDataMutation.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	// Unmarshal JSON data into a map
	var jsonData interface{}
	err = json.Unmarshal(fileData, &jsonData)
	if err != nil {
		fmt.Println("Error unmarshalling JSON data:", err)
		return
	}

	// Add "a" to all values in the JSON
	modifiedJSON := addAToValues(jsonData)

	// Print the modified JSON
	modifiedJSONBytes, err := json.MarshalIndent(modifiedJSON, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling modified JSON:", err)
		return
	}

	fmt.Println(string(modifiedJSONBytes))
}
