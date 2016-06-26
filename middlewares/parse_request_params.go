package middlewares

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yikailee/golang/dto"
)

var (
	ErrCantSetValue      = errors.New("can not set value")
	ErrUnhandleType      = errors.New("can not handle type")
	ErrEmptyValue        = errors.New("empty value")
	ErrNotEnoughValue    = errors.New("not enough values")
	ErrOverflowInt       = errors.New("reflect set int overflow")
	ErrOverflowFloat     = errors.New("reflect set float overflow")
	ErrInvalidReflectVal = errors.New("invalid reflect value")
	ErrErrorType         = errors.New("error type")
	ErrIllegalJSON       = errors.New("ilegal json format")
)

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

func appendValues(dst, src map[string][]string) {
	for k, vs := range src {
		dst[k] = append(dst[k], vs...)
	}
}

func getValuesFromJSON(r io.Reader) (values map[string][]string, err error) {
	content := json.NewDecoder(r)
	values = make(map[string][]string)
	for {
		k, e := content.Token()
		if e == io.EOF {
			break
		}
		if e != nil {
			err = ErrIllegalJSON
			break
		}

		kn := reflect.TypeOf(k).Name()
		if kn == "Delim" {
			continue
		}
		if kn != "string" {
			err = ErrIllegalJSON
			break
		}
		kv := reflect.ValueOf(k).String()

		v, e := content.Token()
		if e != nil {
			err = ErrIllegalJSON
			break
		}
		vn := reflect.TypeOf(v).Name()
		if vn != "Delim" {
			if vn != "string" {
				err = ErrIllegalJSON
				break
			}

			vv := reflect.ValueOf(v).String()
			values[kv] = append(values[kv], vv)
			continue

		} else if reflect.ValueOf(v).Int() != 91 { // '[' == 91
			err = ErrIllegalJSON
			break
		}

		for content.More() {
			v, err := content.Token()
			if err != nil || reflect.TypeOf(v).Name() != "string" {
				err = ErrIllegalJSON
				break
			}
			vv := reflect.ValueOf(v).String()
			values[kv] = append(values[kv], vv)
		}
		if err != nil {
			break
		}
	}
	if err != nil {
		values = make(map[string][]string)
	}
	return
}

func parseJSONBody(r *http.Request) (err error) {
	if r.Method != "POST" && r.Method != "PUT" && r.Method != "PATCH" {
		return
	}
	if r.Body == nil {
		err = errors.New("missing request json body")
		return
	}
	ct := r.Header.Get("Content-Type")
	if ct == "" {
		ct = "application/octet-stream"
	}
	ct, _, err = mime.ParseMediaType(ct)
	switch {
	case ct == "application/json":
		var reader io.Reader = r.Body
		maxContentSize := int64(10 << 20) // 10 MB is a lot of text.
		reader = io.LimitReader(r.Body, maxContentSize+1)
		vs, e := getValuesFromJSON(reader)
		if err == nil {
			err = e
		}
		if err == nil {
			appendValues(r.Form, vs)
		}
	}

	return
}

func parseUrlEmbededParams(r *http.Request) (values map[string][]string, err error) {
	// check whether params if embeded in url
	values = make(map[string][]string)
	vars := mux.Vars(r)
	if nil != vars {
		for k, v := range vars {
			values[k] = append(values[k], v)
		}
	}
	return
}

func ParseRequestParams(r *http.Request, in dto.ValidRequestDTO) bool {
	values, _ := parseUrlEmbededParams(r)
	r.ParseMultipartForm(defaultMaxMemory)
	parseJSONBody(r)
	appendValues(values, r.Form)

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		wanted := t.Field(i).Tag.Get("dto")
		if wanted == "" {
			continue
		}
		required := !(t.Field(i).Tag.Get("required") == "false")
		if _, ok := values[wanted]; !ok {
			if required {
				return false
			}
			continue
		}

		err := setReflectValue(v.Field(i), values[wanted])
		if err != nil {
			if required {
				return false
			}
			continue
		}
		in.MarkSet(wanted)
	}

	return true
}

func setReflectSingleValue(rVal reflect.Value, str string) error {
	if !rVal.CanSet() {
		return ErrCantSetValue
	}
	if !rVal.IsValid() {
		return ErrInvalidReflectVal
	}

	// handle int type
	intType := 0
	switch rVal.Kind() {
	case reflect.Int8:
		intType = 8
	case reflect.Int16:
		intType = 16
	case reflect.Int32:
	case reflect.Int:
		intType = 32
	case reflect.Int64:
		intType = 64
	}
	if 0 != intType {
		p, err := strconv.ParseInt(str, 10, intType)
		if err != nil {
			return err
		}
		if rVal.OverflowInt(p) {
			return ErrOverflowInt
		}
		rVal.SetInt(p)
		return nil
	}

	// handle float type
	fType := 0
	switch rVal.Kind() {
	case reflect.Float32:
		fType = 32
	case reflect.Float64:
		fType = 64
	}
	if 0 != fType {
		p, err := strconv.ParseFloat(str, fType)
		if err != nil {
			return err
		}
		if rVal.OverflowFloat(p) {
			return ErrOverflowFloat
		}
		rVal.SetFloat(p)
		return nil
	}

	// handle other type
	switch rVal.Kind() {
	case reflect.Bool:
		p, err := strconv.ParseBool(str)
		if err != nil {
			return err
		}
		rVal.SetBool(p)
	case reflect.String:
		rVal.SetString(str)
	default:
		return ErrUnhandleType
	}

	return nil
}

func setReflectValue(rVal reflect.Value, strs []string) error {
	if len(strs) == 0 {
		return ErrEmptyValue
	}

	switch rVal.Kind() {
	case reflect.Array:
		var n = rVal.Len()
		if n > len(strs) {
			return ErrNotEnoughValue
		}
		for i := 0; i < n; i++ {
			err := setReflectSingleValue(rVal.Index(i), strs[i])
			if err != nil {
				return err
			}
		}
	case reflect.Slice:
		t := rVal.Type()
		n := len(strs)
		rVal.Set(reflect.MakeSlice(t, n, n))
		for i := 0; i < n; i++ {
			err := setReflectSingleValue(rVal.Index(i), strs[i])
			if err != nil {
				return err
			}
		}
	default:
		return setReflectSingleValue(rVal, strs[0])
	}

	return nil
}
