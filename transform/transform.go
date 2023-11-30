package transform

import (
	"encoding/json"

	"github.com/itchyny/gojq"
)

func ConvertAndTransform(object any, xform string) (string, error) {
	jsonBytes, err := json.Marshal(object)
	if err != nil {
		return "", err
	}
	jsonString := string(jsonBytes)
	return Transform(jsonString, xform)
}

func Transform(jsonStr, xform string) (string, error) {
	query, err := gojq.Parse(xform)
	if err != nil {
		return "", err
	}

	var result map[string]any

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return "", err
	}

	iter := query.Run(result) // or query.RunWithContext
	retval := make([]byte, 0)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return "", err
		}
		retval, err = json.MarshalIndent(v, "", "  ")
		if err != nil {
			return "", err
		}
	}
	return string(retval), nil
}
