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
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unsafe"
)

type field struct {
	a          annotation
	fieldName  string
	fieldType  *reflect.Type
	fieldValue *reflect.Value
	public     bool
}

// contains raw info about string tag field. e.g: default=hello,
type strTag struct {
	value  string
	exists bool
}

type annotation struct {
	value     string
	name      string
	tagName   string
	def       strTag
	req       strTag
	protected strTag
	present   bool
}

type meta map[string]field

// Bind binds environment variables into structure
func Bind(annotations map[string]string, s interface{}) (err error) {
	var meta meta
	if s == nil {
		return fmt.Errorf("invalid structure value (nil)")
	}
	if annotations == nil {
		return nil
	}
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s).Kind()
	if t != reflect.Ptr {
		return fmt.Errorf("argument must be pointer to structure")
	}
	if v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("argument must be pointer to structure")
	}
	meta, err = roll(annotations, v.Elem(), v.Elem().Type().Name(), "")
	if err != nil {
		return
	}
	err = bind(annotations, meta)
	return
}

// binds meta to structure pointer
func bind(annotations map[string]string, m meta) (err error) {
	for k, v := range m {
		f := reflect.NewAt(v.fieldValue.Type(), unsafe.Pointer(v.fieldValue.UnsafeAddr())).Elem()
		switch f.Interface().(type) {
		case bool:
			var b bool
			if v.a.protected.isTrue() {
				continue
			}
			b, err = boolean(annotations, v.a)
			if err != nil {
				return
			}
			f.SetBool(b)
			continue

		case int, int8, int16, int32, int64:
			if v.a.protected.isTrue() && v.fieldValue.Int() != 0 {
				continue
			}
			err = setNumeric(annotations, f, v)
			if err != nil {
				return
			}
			continue

		case float32, float64:
			if v.a.protected.isTrue() && v.fieldValue.Float() != 0 {
				continue
			}
			err = setNumeric(annotations, f, v)
			if err != nil {
				return
			}
			continue

		case uint, uint8, uint16, uint32, uint64:
			if v.a.protected.isTrue() && v.fieldValue.Uint() != 0 {
				continue
			}
			err = setNumeric(annotations, f, v)
			if err != nil {
				return
			}
			continue

		case string:
			var s string
			if v.a.protected.isTrue() && v.fieldValue.String() != "" {
				continue
			}
			s = GetAsStringOrFallback(annotations, v.a.name, v.a.def.value)
			f.SetString(s)
			continue

		case []string:
			if v.a.protected.isTrue() && !v.fieldValue.IsNil() {
				continue
			}
			var ss []string
			ss = GetAsArrayOfStringsOrFallback(annotations, v.a.name, v.a.def.asStringSlice())
			f.Set(reflect.ValueOf(ss))
			continue

		case []int, []int8, []int16, []int32, []int64, []float32, []float64, []uint, []uint8, []uint16, []uint32, []uint64:
			if v.a.protected.isTrue() && !v.fieldValue.IsNil() {
				continue
			}
			var floats []float64
			floats, err = floatSlice(annotations, v.a)
			if err != nil {
				return
			}
			setNumericSlice(f, floats)
			continue

		case []bool:
			if v.a.protected.isTrue() && !v.fieldValue.IsNil() {
				continue
			}
			var bs []bool
			bs, err = boolSlice(annotations, v.a)
			if err != nil {
				return
			}
			f.Set(reflect.ValueOf(bs))
			continue

		default:
			err = fmt.Errorf("unsupported type %s: %s", k, v.fieldValue.Type().Name())
		}
	}
	return err
}

func setNumericSlice(f reflect.Value, floats []float64) {
	switch f.Interface().(type) {
	case []uint:
		f.Set(reflect.ValueOf(convertToUInt(floats)))
		return
	case []uint8:
		f.Set(reflect.ValueOf(convertToUInt8(floats)))
		return
	case []uint16:
		f.Set(reflect.ValueOf(convertToUInt16(floats)))
		return
	case []uint32:
		f.Set(reflect.ValueOf(convertToUInt32(floats)))
		return
	case []uint64:
		f.Set(reflect.ValueOf(convertToUInt64(floats)))
		return
	case []int:
		f.Set(reflect.ValueOf(convertToInt(floats)))
		return
	case []int8:
		f.Set(reflect.ValueOf(convertToInt8(floats)))
		return
	case []int16:
		f.Set(reflect.ValueOf(convertToInt16(floats)))
		return
	case []int32:
		f.Set(reflect.ValueOf(convertToInt32(floats)))
		return
	case []int64:
		f.Set(reflect.ValueOf(convertToInt64(floats)))
		return
	case []float32:
		f.Set(reflect.ValueOf(convertToFloat32(floats)))
		return
	case []float64:
		f.Set(reflect.ValueOf(floats))
		return
	}
}

func setNumeric(annotations map[string]string, f reflect.Value, v field) (err error) {
	var fl float64
	fl, err = float(annotations, v.a)
	if err != nil {
		return
	}
	switch f.Interface().(type) {
	case int, int8, int16, int32, int64:
		f.SetInt(int64(fl))
		return
	case float32, float64:
		f.SetFloat(fl)
		return
	case uint, uint8, uint16, uint32, uint64:
		f.SetUint(uint64(fl))
		return
	}
	return
}

// recoursive function builds meta structure
func roll(annotations map[string]string, value reflect.Value, n, prefix string) (m meta, err error) {
	const tagAnnotation = "annotation"

	m = meta{}
	for i := 0; i < value.NumField(); i++ {
		var a annotation
		vf := value.Field(i)
		tf := value.Type().Field(i)
		key := fmt.Sprintf("%s.%s", n, tf.Name)
		tag := tf.Tag.Get(tagAnnotation)
		if vf.Kind() == reflect.Struct {
			var sm meta
			prefix := strings.TrimPrefix(fmt.Sprintf("%s_%s", prefix, getTagName(tag)), "_")
			sm, err = roll(annotations, vf, key, prefix)
			if err != nil {
				return
			}
			for k, v := range sm {
				m[k] = v
			}
			continue
		}
		if tag == "" {
			continue
		}
		if a, err = parseTag(annotations, tag, prefix); err != nil {
			return
		}
		if !a.present && a.req.value == "true" {
			err = fmt.Errorf("%s is required", a.name)
			return
		}
		m[key] = field{
			a:          a,
			fieldName:  tf.Name,
			fieldType:  &tf.Type,
			fieldValue: &vf,
			public:     tf.PkgPath == "",
		}
	}
	return m, err
}

// parseTag, retrieves annotation info and metadata
func parseTag(annotations map[string]string, tag, prefix string) (e annotation, err error) {
	var def, req, protected strTag
	var tagName = getTagName(tag)
	req, err = getTagProperty(tag, "require")
	if err != nil {
		return
	}
	def, err = getTagProperty(tag, "default")
	if err != nil {
		return
	}
	protected, err = getTagProperty(tag, "protected")
	if err != nil {
		return
	}
	name := getName(tagName, prefix)
	value, exists := annotations[name]
	e = annotation{
		name:      name,
		tagName:   tagName,
		value:     value,
		req:       req,
		def:       def,
		protected: protected,
		present:   exists,
	}
	return
}

func getName(envName, prefix string) string {
	if prefix != "" {
		return fmt.Sprintf("%s_%s", prefix, envName)
	}
	return envName
}

func getTagName(tag string) string {
	//implement regex if needed `return e.g: regexp.MustCompile("[a-zA-Z_]+[a-zA-Z0-9_]*").FindString(tag)`
	parts := strings.Split(tag, ",")
	if len(parts) > 0 {
		return strings.TrimLeft(strings.TrimRight(parts[0], " "), " ")
	}
	return strings.TrimLeft(strings.TrimRight(tag, " "), " ")
}

// parses value from annotation tag and returns <tag value, tag value exists, error>
func getTagProperty(tag, t string) (r strTag, err error) {
	const arr = `\[\w*\s*\!*\@*\#*\$*\%*\^*\&*\**\(*\)*\_*\-*\+*\<*\>*\?*\~*\=*\,*\.*\/*\{*\}*\|*\;*\:*\/*\'*\"*\/*\\*`
	const scalar = `\[*\]*\w*\s*\!*\@*\#*\$*\%*\^*\&*\**\(*\)*\_*\-*\+*\<*\>*\?*\~*\=*\.*\/*\{*\}*\|*\;*\:*\/*\'*\"*\/*\\*`
	r = strTag{}
	var findRegex, removeRegex *regexp.Regexp
	//	findRegex, err = regexp.Compile(",\\s*" + t + "\\s*=((\\s*([\\[\\w*\\,*\\.*\\s*\\-*])*\\])|(\\s*\\w*\\.*\\-*)*)")
	findRegex, err = regexp.Compile(",\\s*" + t + "\\s*=((\\s*\\[[" + arr + "]*\\])|(" + scalar + ")*)")
	if err != nil {
		err = fmt.Errorf("ivalid %s", t)
		return
	}
	removeRegex, err = regexp.Compile(",\\s*" + t + "\\s*=\\s*")
	if err != nil {
		err = fmt.Errorf("ivalid %s", t)
		return
	}
	match := findRegex.FindString(tag)
	if match == "" {
		return
	}
	remove := removeRegex.FindString(strings.ToLower(tag))
	r.value = strings.ReplaceAll(match, remove, "")
	r.exists = true
	return
}

func (t strTag) asStringSlice() (s []string) {
	if !t.exists {
		return
	}
	envdef := strings.TrimSuffix(strings.TrimPrefix(t.value, " "), " ")
	envdef = strings.TrimSuffix(strings.TrimPrefix(envdef, "["), "]")
	s = strings.Split(envdef, ",")
	if s[0] == "" {
		s = []string{}
	}
	return
}

func (t strTag) isTrue() bool {
	if !t.exists {
		return false
	}
	b, _ := strconv.ParseBool(t.value)
	return b
}
