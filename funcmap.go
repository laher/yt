package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"

	yaml "gopkg.in/yaml.v2"
)

func funcMap() template.FuncMap {
	var funcMap = template.FuncMap{
		"yaml":        toYAML,
		"json":        toJSON,
		"o":           other, // TODO remove this
		"go":          toGo,
		"del":         del,
		"set":         set,
		"merge":       merge,
		"interpolate": interpolate,
		"cr":          func() string { return "\n" },
		"newdoc":      func() string { return "---\n" },
		"tableflip":   func() string { return "(╯°□°）╯︵ ┻━┻" },
		"getenv": func(key string, def ...string) string {
			val := os.Getenv(key)
			if val != "" {
				return val
			}
			if len(def) > 0 {
				return def[0]
			}
			return ""
		},
		"has": func(in interface{}, key string) bool {
			// TODO
			return false
		},
	}
	return funcMap
}

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

func merge(x, y map[interface{}]interface{}) string {
	for k, v := range y {
		x[k] = v
	}
	return "" // empty-string avoids output
}

func interpolate(input, data map[interface{}]interface{}) string {
	for k, v := range input {
		input[k] = interpolateVal(k, v, data)
	}
	return ""
}

const (
	prefixSpecialType = "__:" + string('\uFFFE') + ":" // use some unicode chars to avoid collisions with legal strings
)

func interpolateFuncMap() template.FuncMap {
	funcs := funcMap()
	funcs["num"] = func(in interface{}) string {
		return fmt.Sprintf("%snum:%v", prefixSpecialType, in)
	}
	funcs["bool"] = func(in interface{}) string {
		return fmt.Sprintf("%sbool:%v", prefixSpecialType, in)
	}
	return funcs
}

func interpolateVal(k, v interface{}, data map[interface{}]interface{}) interface{} {
	switch vs := v.(type) {
	case string:
		if strings.Contains(vs, "{{") {
			funcs := funcMap()

			canProcessSpecialTypes := !strings.Contains(vs, prefixSpecialType)
			if canProcessSpecialTypes {
				funcs = interpolateFuncMap()
			}
			tmpl, err := template.New("interpolator").Funcs(funcs).Parse(vs)
			if err != nil {
				log.Fatal(err)
			}
			output := &bytes.Buffer{}
			err = tmpl.Execute(output, data)
			if err != nil {
				log.Fatal(err)
			}
			out := output.String()
			log.Printf("Interpolate [%s]: [%s] => [%s]", k, vs, out)
			if canProcessSpecialTypes { // bools and numbers would not have any prefix
				if strings.HasPrefix(out, prefixSpecialType+"bool:") {
					switch out {
					case prefixSpecialType + "bool:true":
						return true
					case prefixSpecialType + "bool:false":
						return false
					}
				}
				if strings.HasPrefix(out, prefixSpecialType+"num:") {
					number := out[len(prefixSpecialType+"num:"):]
					s, err := strconv.ParseFloat(number, 64)
					if err != nil {
						log.Fatal(err)
					}
					return s
				}
			}
			return out
		}
		return vs
	case map[interface{}]interface{}:
		log.Printf("Deeper (map): %s", k)
		interpolate(vs, data) // returns empty string. Return the map itself
		return vs
	case []interface{}:
		log.Printf("Deeper (slice): %s", k)
		for i, item := range vs {
			vs[i] = interpolateVal(k, item, data)
		}
		return vs
	default:
		log.Printf("Straight replace: %s - %T", k, vs)
		return vs
	}
}
