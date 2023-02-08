package annotation

import (
	"fmt"
	"strconv"
	"strings"
)

func boolean(annotations map[string]string, a annotation) (b bool, err error) {
	var d bool
	if a.def.exists {
		d, err = strconv.ParseBool(a.def.value)
		if err != nil {
			err = fmt.Errorf("can't convert default value %s of '%s' to bool", a.name, a.def.value)
			return
		}
	}
	b, err = GetAsBoolOrFallback(annotations, a.name, d)
	return
}

func float(annotations map[string]string, a annotation) (f float64, err error) {
	var d float64
	if a.def.exists {
		d, err = strconv.ParseFloat(a.def.value, 64)
		if err != nil {
			err = fmt.Errorf("can't convert default value %s of '%s' to float64", a.name, a.def.value)
			return
		}
	}
	f, err = GetAsFloat64OrFallback(annotations, a.name, d)
	if err != nil {
		err = fmt.Errorf("can't read %s and parse value '%s' to float64", a.name, a.value)
	}
	return
}

func floatSlice(annotations map[string]string, a annotation) (fs []float64, err error) {
	var d []float64
	if a.def.asStringSlice() != nil {
		d = make([]float64, 0)
		for _, s := range a.def.asStringSlice() {
			var fl float64
			fl, err = strconv.ParseFloat(strings.Trim(s, " "), 64)
			if err != nil {
				err = fmt.Errorf("can't convert default %s to slice of float64", a.def.asStringSlice())
				return
			}
			d = append(d, fl)
		}
	}
	fs, err = GetAsArrayOfFloat64OrFallback(annotations, a.name, d)
	if err != nil {
		err = fmt.Errorf("can't parse %s as slice of float64 '%s'", a.name, a.value)
	}
	return
}

func boolSlice(annotations map[string]string, a annotation) (bs []bool, err error) {
	var d []bool
	if a.def.asStringSlice() != nil {
		d = make([]bool, 0)
		for _, s := range a.def.asStringSlice() {
			var b bool
			b, err = strconv.ParseBool(strings.Trim(s, " "))
			if err != nil {
				err = fmt.Errorf("can't convert default %s to slice of bool", a.def.asStringSlice())
				return
			}
			d = append(d, b)
		}
	}
	bs, err = GetAsArrayOfBoolOrFallback(annotations, a.name, d)
	if err != nil {
		err = fmt.Errorf("can't parse %s as slice of bool '%s'", a.name, a.value)
	}
	return
}
