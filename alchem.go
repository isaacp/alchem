package alchem

import (
	"encoding/json"

	"github.com/itchyny/gojq"
)

func TransformObject[T any](object any, xform string) (T, error) {
	var retVal T
	jsonBytes, err := json.Marshal(object)
	if err != nil {
		return retVal, err
	}
	jsonString := string(jsonBytes)
	jsn, err := TransformJson(jsonString, xform)
	if err != nil {
		return retVal, err
	}

	err = json.Unmarshal([]byte(jsn), &retVal)
	if err != nil {
		return retVal, err
	}

	return retVal, nil
}

func TransformJson(jsonStr, xform string) (string, error) {
	query, err := gojq.Parse(xform)
	if err != nil {
		return "", err
	}

	var result []any

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return "", err
	}

	iter := query.Run(result)
	mutant := make([]byte, 0)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return "", err
		}
		mutant, err = json.Marshal(v)
		if err != nil {
			return "", err
		}
	}

	return string(mutant), nil
}
