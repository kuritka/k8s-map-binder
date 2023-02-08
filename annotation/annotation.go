/*
Copyright 2021 The k8gb Contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Generated by GoLic, for more details see: https://github.com/AbsaOSS/golic
*/

package annotation

import (
	"strconv"
	"strings"
)

// GetAsStringOrFallback returns the annotation variable for the given key
// and falls back to the given defaultValue if not set
func GetAsStringOrFallback(annotations map[string]string, key, defaultValue string) string {
	if v, found := annotations[key]; found {
		return v
	}
	return defaultValue
}

// GetAsArrayOfStringsOrFallback returns the annotation variable for the given key
// and falls back to the given defaultValue if not set
// GetAsArrayOfStringsOrFallback trims all whitespaces from input i.e. "us, fr, au" -> {"us","fr","au"}
func GetAsArrayOfStringsOrFallback(annotations map[string]string, key string, defaultValue []string) []string {
	if v, ex := annotations[key]; ex {
		if v == "" {
			return []string{}
		}
		arr := strings.Split(strings.ReplaceAll(v, " ", ""), ",")
		if len(arr) != 0 {
			return arr
		}
	}
	return defaultValue
}

// GetAsArrayOfIntsOrFallback returns the annotation variable for the given key
// and falls back to the given defaultValue if not set
func GetAsArrayOfIntsOrFallback(annotations map[string]string, key string, defaultValue []int) (ints []int, err error) {
	if v, ex := annotations[key]; ex {
		if v == "" {
			return []int{}, nil
		}
		slice := strings.Split(strings.ReplaceAll(v, " ", ""), ",")
		for _, s := range slice {
			var i int
			i, err = strconv.Atoi(s)
			if err != nil {
				return defaultValue, err
			}
			ints = append(ints, i)
		}
		return ints, nil
	}
	return defaultValue, nil
}

// GetAsArrayOfFloat64OrFallback returns the annotation variable for the given key
// and falls back to the given defaultValue if not set
func GetAsArrayOfFloat64OrFallback(annotations map[string]string, key string, defaultValue []float64) (floats []float64, err error) {
	if v, ex := annotations[key]; ex {
		if v == "" {
			return []float64{}, nil
		}
		slice := strings.Split(strings.ReplaceAll(v, " ", ""), ",")
		for _, s := range slice {
			var f float64
			f, err = strconv.ParseFloat(s, 64)
			if err != nil {
				return defaultValue, err
			}
			floats = append(floats, f)
		}
		return floats, nil
	}
	return defaultValue, nil
}

// GetAsArrayOfBoolOrFallback returns the annotation variable for the given key
// and falls back to the given defaultValue if not set
func GetAsArrayOfBoolOrFallback(annotations map[string]string, key string, defaultValue []bool) (bools []bool, err error) {
	if v, ex := annotations[key]; ex {
		if v == "" {
			return []bool{}, nil
		}
		slice := strings.Split(strings.ReplaceAll(v, " ", ""), ",")
		for _, s := range slice {
			var b bool
			b, err = strconv.ParseBool(s)
			if err != nil {
				return defaultValue, err
			}
			bools = append(bools, b)
		}
		return bools, nil
	}
	return defaultValue, nil
}

// GetAsIntOrFallback returns the annotation variable (parsed as integer) for
// the given key and falls back to the given defaultValue if not set
func GetAsIntOrFallback(annotations map[string]string, key string, defaultValue int) (int, error) {
	if v, ex := annotations[key]; ex {
		value, err := strconv.Atoi(v)
		if err != nil {
			return defaultValue, err
		}
		return value, nil
	}
	return defaultValue, nil
}

// GetAsFloat64OrFallback returns the annotation variable (parsed as float64) for
// the given key and falls back to the given defaultValue if not set
func GetAsFloat64OrFallback(annotations map[string]string, key string, defaultValue float64) (float64, error) {
	if v, ex := annotations[key]; ex {
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return defaultValue, err
		}
		return value, nil
	}
	return defaultValue, nil
}

// GetAsBoolOrFallback returns the annotation variable for the given key,
// parses it as boolean and falls back to the given defaultValue if not set
func GetAsBoolOrFallback(annotations map[string]string, key string, defaultValue bool) (val bool, err error) {
	if v, ex := annotations[key]; ex {
		val, err = strconv.ParseBool(v)
		if err != nil {
			return
		}
		return
	}
	return defaultValue, nil
}
