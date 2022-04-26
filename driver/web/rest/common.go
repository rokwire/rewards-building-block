package rest

import (
	"net/http"
	"strconv"
)

func getStringQueryParam(r *http.Request, paramName string) *string {
	params, ok := r.URL.Query()[paramName]
	if ok && len(params[0]) > 0 {
		value := params[0]
		return &value
	}
	return nil
}

func getInt64QueryParam(r *http.Request, paramName string) *int64 {
	params, ok := r.URL.Query()[paramName]
	if ok && len(params[0]) > 0 {
		val, err := strconv.ParseInt(params[0], 0, 64)
		if err == nil {
			return &val
		}
	}
	return nil
}

func getIntQueryParam(r *http.Request, paramName string, defaultValue int) int {
	params, ok := r.URL.Query()[paramName]
	if ok && len(params[0]) > 0 {
		val, err := strconv.Atoi(params[0])
		if err == nil {
			return val
		}
	}

	return defaultValue
}

func getBoolQueryParam(r *http.Request, paramName string, defaultValue *bool) *bool {
	params, ok := r.URL.Query()[paramName]
	if ok && len(params[0]) > 0 {
		val, err := strconv.Atoi(params[0])
		if err == nil {
			bValue := val != 0
			return &bValue
		}
	}
	return defaultValue
}
