package middleware

import (
	"net/url"

	"github.com/spf13/cast"
	"github.com/xeipuuv/gojsonschema"
)

type JsonSchemaCastingFunc func(val string) interface{}

var (
	JsonSchemaCastingInt JsonSchemaCastingFunc = func(val string) interface{} {
		return cast.ToInt(val)
	}
	JsonSchemaCastingBool JsonSchemaCastingFunc = func(val string) interface{} {
		return cast.ToBool(val)
	}
)

/*
if not empty string will be assign in map
casting Type can using only exactly one value

	Example JsonSchemaQueryParams(c.Request().URL.Query(), map[string]middleware.JsonSchemaCastingFunc{
		"page":     middleware.JsonSchemaCastingInt,
		"per_page": middleware.JsonSchemaCastingInt,
	})
*/
func JsonSchemaQueryParams(values url.Values, castingTypeValueFuncs map[string]JsonSchemaCastingFunc) map[string]interface{} {
	var m = make(map[string]interface{})
	var casting = func(key string, value string) interface{} {
		if len(castingTypeValueFuncs) > 0 {
			if rf, ok := castingTypeValueFuncs[key]; ok {
				return rf(value)
			}
		}
		return value
	}
	if len(values) > 0 {
		for k, v := range values {
			switch {
			case len(v) == 1:
				var val = v[0]
				if val != "" {
					m[k] = casting(k, val)
				}
			case len(v) > 1:
				var vals = make([]interface{}, 0, len(v))
				for _, val := range v {
					if val != "" {
						vals = append(vals, casting(k, val))
					}
				}
				if len(vals) > 0 {
					m[k] = vals
				}
			}
		}
	}

	return m
}

func JsonSchemaFormat(results []gojsonschema.ResultError) map[string]interface{} {
	m := map[string]interface{}{"errors": []interface{}{}}
	for _, item := range results {
		field := item.Field()
		if field == "(root)" {
			field = ""
			field = cast.ToString(item.Details()["property"])
		}

		m["message"] = "invalid arguments"
		m["errors"] = append(m["errors"].([]interface{}), []interface{}{
			map[string]interface{}{
				"field":       field,
				"description": item.Description(),
				"details":     item.Details(),
				"error_from":  item.Type(),
			},
		})
	}
	return m
}
