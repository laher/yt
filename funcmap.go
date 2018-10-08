package main

import (
	"encoding/json"
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

func toYAML(x interface{}) (string, error) {
	bts, err := yaml.Marshal(x)
	return string(bts), err
}

func toJSON(x interface{}) (string, error) {
	switch xT := x.(type) {
	case map[interface{}]interface{}:
		tryMap := coerceKeysForJSON(xT)
		bts, err := json.Marshal(tryMap)
		return string(bts), err
	}
	bts, err := json.Marshal(x)
	return string(bts), err
}

func coerceKeysForJSON(in map[interface{}]interface{}) map[string]interface{} {

	tryMap := map[string]interface{}{}
	for k, v := range in {
		ks := fmt.Sprintf("%v", k)
		switch vT := v.(type) {
		case map[interface{}]interface{}:
			tryMap[ks] = coerceKeysForJSON(vT)
		default:
			tryMap[ks] = v
		}
	}
	return tryMap
}

func other(x interface{}) string {
	return fmt.Sprintf("%+v", x)
}

func toGo(x interface{}) string {
	return fmt.Sprintf("%#v", x)
}

func del(x interface{}, y string) interface{} {
	switch t := x.(type) {
	case map[interface{}]interface{}:
		delete(t, y)
	default:
		panic(fmt.Sprintf("x (%t) is not a map[interface{}]interface{}", x))
	}
	return x
}

func set(x map[interface{}]interface{}, k string, y interface{}) string {
	x[k] = y
	return "" // empty-string avoids output
}
