package alchem

import (
	"encoding/json"
	"strings"

	"github.com/itchyny/gojq"
)

func ConvertAndTransform[T any](object any, xform string) (string, error) {
	jsonBytes, err := json.Marshal(object)
	if err != nil {
		return "", err
	}
	jsonString := strings.Replace(string(jsonBytes), "\\\"", "", -1)
	return Transform[T](jsonString, xform)
}

func Transform[T any](jsonStr, xform string) (string, error) {
	query, err := gojq.Parse(xform)
	if err != nil {
		return "", err
	}

	//var result T

	//if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
	//	return "", err
	//}

	iter := query.Run(jsonStr) // or query.RunWithContext
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
