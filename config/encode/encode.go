package encode

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/zhaolion/gostack/util/xerr"
	"gopkg.in/yaml.v2"
)

// GetValueFn ...
type GetValueFn func(string) (string, error)

// Encode ...
func Encode(i interface{}, tag string, fn GetValueFn) error {
	val := reflect.ValueOf(i)
	typ := val.Type()
	if kind := typ.Kind(); kind != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return errors.New("pointer to struct type required")
	}
	return walkStructValue(val.Elem(), tag, fn)
}

func walkStructValue(v reflect.Value, tag string, fn GetValueFn) (err error) {
	typ := v.Type()
	for i := 0; i < typ.NumField(); i++ {
		if err != nil {
			return err
		}
		tagVal := typ.Field(i).Tag.Get(tag)
		if tagVal == "-" {
			continue
		}
		fVal := reflect.Indirect(v.Field(i))

		// 获取 tag, 类似 `dva:"foo.yaml"`
		name, _ := parseTag(tagVal)

		switch fVal.Kind() {
		case reflect.Struct:
			if isJson(name) {
				fieldObj := reflect.New(typ.Field(i).Type).Interface()
				value, err := fn(name)
				if err != nil {
					return xerr.WithStack(err)
				}
				if err := json.Unmarshal([]byte(value), fieldObj); err != nil {
					return xerr.WithStack(err)
				}

				fVal.Set(reflect.Indirect(reflect.ValueOf(fieldObj)))

				continue
			} else if isYaml(name) {
				value, err := fn(name)
				if err != nil {
					return xerr.WithStack(err)
				}
				fieldObj := reflect.New(typ.Field(i).Type).Interface()
				if err := yaml.Unmarshal([]byte(value), fieldObj); err != nil {
					return xerr.WithStack(err)
				}
				fVal.Set(reflect.Indirect(reflect.ValueOf(fieldObj)))

				continue
			} else {
				err = walkStructValue(fVal, tag, fn)
			}
		case reflect.Slice, reflect.Array:
			if isCSV(name) {
				value, err := fn(name)
				if err != nil {
					return xerr.WithStack(err)
				}
				fieldObj := reflect.New(typ.Field(i).Type).Interface()
				if err := gocsv.UnmarshalBytes([]byte(value), fieldObj); err != nil {
					return xerr.WithStack(err)
				}
				fVal.Set(reflect.Indirect(reflect.ValueOf(fieldObj)))

				continue
			} else {
				if name == "" || name == "-" {
					continue
				}

				var value string
				value, err = fn(name)
				if value != "" {
					err = SetValue(fVal, value)
				}
			}
		default:
			if name == "" || name == "-" {
				continue
			}

			var value string
			value, err = fn(name)
			if value != "" {
				err = SetValue(fVal, value)
			}
		}
	}
	return err
}

// parseTag ...
func parseTag(tag string) (string, []string) {
	s := strings.Split(tag, ",")
	return s[0], s[1:]
}

// SetValue ...
func SetValue(v reflect.Value, value string) (err error) {
	switch v.Kind() {
	case reflect.Bool:
		err = setBoolVal(v, value)
	case reflect.String:
		v.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		err = setIntVal(v, value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		err = setUintVal(v, value)
	case reflect.Float32, reflect.Float64:
		err = setFloatVal(v, value)
	case reflect.Array, reflect.Slice:
		err = setSliceVal(v, value)
	default:
		err = fmt.Errorf("unknown supported type: %s, %s", v.Kind(), v.Type().Name())
	}
	return err
}

func setBoolVal(v reflect.Value, value string) error {
	if value == "" {
		return nil
	} else if value == "true" {
		v.SetBool(true)
	} else if value == "false" {
		v.SetBool(false)
	} else {
		return fmt.Errorf("invalid bool value: %s", value)
	}
	return nil
}

func setIntVal(v reflect.Value, value string) error {
	x, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	v.SetInt(x)
	return nil
}

func setUintVal(v reflect.Value, value string) error {
	x, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return err
	}
	v.SetUint(x)
	return nil
}

func setFloatVal(v reflect.Value, value string) error {
	x, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	v.SetFloat(x)
	return nil
}

func setSliceVal(v reflect.Value, value string) error {
	vv := strings.Split(value, ",")
	var nv reflect.Value
	if v.Kind() == reflect.Array {
		if v.Len() != len(vv) {
			return errors.New("unmatched length of array")
		}
		nv = v
	} else {
		nv = reflect.MakeSlice(v.Type(), len(vv), len(vv))
	}
	for i := 0; i < len(vv); i++ {
		if err := SetValue(nv.Index(i), vv[i]); err != nil {
			return err
		}
	}
	v.Set(nv)
	return nil
}

func isJson(key string) bool {
	return strings.HasSuffix(key, ".json") || strings.HasSuffix(key, ".json.secret")
}

func isYaml(key string) bool {
	return strings.HasSuffix(key, ".yml") || strings.HasSuffix(key, ".yaml") ||
		strings.HasSuffix(key, ".yml.secret") || strings.HasSuffix(key, ".yaml.secret")
}

func isCSV(key string) bool {
	return strings.HasSuffix(key, ".csv") || strings.HasSuffix(key, ".csv.secret")
}
