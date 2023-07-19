package alchem

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/itchyny/gojq"
)

func ConvertAndTransform(object any, xform string) (string, error) {
	jsonBytes, err := json.Marshal(object)
	if err != nil {
		return "", err
	}
	jsonString := string(jsonBytes)
	fmt.Println("TEST: " + jsonString)
	return Transform(jsonString, xform)
}

func Transform(jsonStr, xform string) (string, error) {
	query, err := gojq.Parse(xform)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}

	var result map[string]any

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		fmt.Println(err)
		return "", nil
	}

	iter := query.Run(result) // or query.RunWithContext
	retval := make([]byte, 0)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			log.Fatalln(err)
		}
		retval, err = json.MarshalIndent(v, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
	}
	return string(retval), nil
}
