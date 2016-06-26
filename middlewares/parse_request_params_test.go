package middlewares

import (
	"bytes"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInt64SetReflectValueHappyPath(t *testing.T) {
	testData := []struct {
		name   string
		input  []string
		Set    int64
		wanted int64
	}{
		{
			name:   "happy path 1",
			input:  []string{"3"},
			wanted: 3,
		},
		{
			name:   "happy path 2",
			input:  []string{"-5"},
			wanted: -5,
		},
	}
	for _, test := range testData {
		var in interface{}
		in = &test
		value := reflect.ValueOf(in).Elem()
		err := setReflectValue(value.Field(2), test.input)
		assert.Nil(t, err, test.name)
		assert.Equal(t, test.wanted, test.Set, test.name)
	}
}

func TestInt64SetReflectValueError(t *testing.T) {
	// test in value
	testData := []struct {
		name  string
		input []string
		Set   int64
	}{
		{
			name:  "overflow",
			input: []string{"10000000000000000000000000000000000000"},
		},
		{
			name:  "invalid string",
			input: []string{"abc"},
		},
		{
			name:  " empty string",
			input: []string{""},
		},
	}
	for _, test := range testData {
		var in interface{}
		in = &test
		value := reflect.ValueOf(in).Elem()
		err := setReflectValue(value.Field(2), test.input)
		assert.NotNil(t, err, test.name)
	}
}

func TestFloat64SetReflectValueHappyPath(t *testing.T) {
	// test in value
	testData := []struct {
		name   string
		input  []string
		Set    float64
		wanted float64
	}{
		{
			name:   "happy path 1",
			input:  []string{"3.14"},
			wanted: 3.14,
		},
		{
			name:   "happy path 2",
			input:  []string{"-5.88"},
			wanted: -5.88,
		},
	}
	for _, test := range testData {
		var in interface{}
		in = &test
		value := reflect.ValueOf(in).Elem()
		err := setReflectValue(value.Field(2), test.input)
		assert.Nil(t, err, test.name)
		assert.Equal(t, test.wanted, test.Set, test.name)
	}
}

func TestFloat64SetReflectValueError(t *testing.T) {
	// test in value
	testData := []struct {
		name  string
		input []string
		Set   float64
	}{
		//{
		//name:  "overflow",
		//input: []string{"10000000000000000000333333330000333322222435550.323433333333333333333333333333333333333"},
		//},
		{
			name:  "invalid string",
			input: []string{"abc"},
		},
		{
			name:  "empty string",
			input: []string{""},
		},
		{
			name:  "empty value",
			input: []string{},
		},
	}
	for _, test := range testData {
		var in interface{}
		in = &test
		value := reflect.ValueOf(in).Elem()
		err := setReflectValue(value.Field(2), test.input)
		assert.NotNil(t, err, test.name)
	}
}

func TestBoolSetReflectValueHappyPath(t *testing.T) {
	// test in value
	testData := []struct {
		name   string
		input  []string
		Set    bool
		wanted bool
	}{
		{
			name:   "happy path 1",
			input:  []string{"true"},
			wanted: true,
		},
		{
			name:   "happy path 2",
			input:  []string{"false"},
			wanted: false,
		},
	}
	for _, test := range testData {
		var in interface{}
		in = &test
		value := reflect.ValueOf(in).Elem()
		err := setReflectValue(value.Field(2), test.input)
		assert.Nil(t, err, test.name)
		assert.Equal(t, test.wanted, test.Set, test.name)
	}
}

func TestBoolSetReflectValueError(t *testing.T) {
	// test in value
	testData := []struct {
		name  string
		input []string
		Set   bool
	}{
		{
			name:  "invalid string",
			input: []string{"abc"},
		},
		{
			name:  "empty string",
			input: []string{""},
		},
		{
			name:  "empty value",
			input: []string{},
		},
	}
	for _, test := range testData {
		var in interface{}
		in = &test
		value := reflect.ValueOf(in).Elem()
		err := setReflectValue(value.Field(2), test.input)
		assert.NotNil(t, err, test.name)
	}
}

func TestArraySetReflectValueHappyPath(t *testing.T) {
	// test in value
	testData := []struct {
		name   string
		input  []string
		Set    [3]int64
		wanted [3]int64
	}{
		{
			name:   "happy path 1",
			input:  []string{"123", "456", "789"},
			wanted: [3]int64{123, 456, 789},
		},
		{
			name:   "happy path 2",
			input:  []string{"111", "222", "333", "444"},
			wanted: [3]int64{111, 222, 333},
		},
	}
	for _, test := range testData {
		var in interface{}
		in = &test
		value := reflect.ValueOf(in).Elem()
		err := setReflectValue(value.Field(2), test.input)
		assert.Nil(t, err, test.name)
		for i, _ := range test.Set {
			assert.Equal(t, test.wanted[i], test.Set[i], test.name)
		}
	}
}

func TestArraySetReflectValueError(t *testing.T) {
	// test in value
	testData := []struct {
		name  string
		input []string
		Set   [3]int64
	}{
		{
			name:  "invalid string",
			input: []string{"abc"},
		},
		{
			name:  "empty string",
			input: []string{""},
		},
		{
			name:  "empty input string array",
			input: []string{},
		},
		{
			name:  "not enough values",
			input: []string{"111", "222"},
		},
	}
	for _, test := range testData {
		var in interface{}
		in = &test
		value := reflect.ValueOf(in).Elem()
		err := setReflectValue(value.Field(2), test.input)
		assert.NotNil(t, err, test.name)
	}
}

func TestSliceSetReflectValueHappyPath(t *testing.T) {
	// test in value
	testData := []struct {
		name   string
		input  []string
		Set    []int64
		wanted []int64
	}{
		{
			name:   "happy path 1",
			input:  []string{"123", "456", "789"},
			wanted: []int64{123, 456, 789},
		},
		{
			name:   "happy path 2",
			input:  []string{"111"},
			wanted: []int64{111},
		},
	}
	for _, test := range testData {
		var in interface{}
		in = &test
		value := reflect.ValueOf(in).Elem()
		err := setReflectValue(value.Field(2), test.input)
		assert.Nil(t, err, test.name)
		assert.Equal(t, len(test.wanted), len(test.Set), "happy path same size")
		for i, _ := range test.Set {
			assert.Equal(t, test.wanted[i], test.Set[i], test.name)
		}
	}
}

func TestSliceSetReflectValueError(t *testing.T) {
	// test in value
	testData := []struct {
		name  string
		input []string
		Set   []int64
	}{
		{
			name:  "invalid string",
			input: []string{"abc", "cdef"},
		},
		{
			name:  "empty string",
			input: []string{""},
		},
		{
			name:  "empty value",
			input: []string{},
		},
	}
	for _, test := range testData {
		var in interface{}
		in = &test
		value := reflect.ValueOf(in).Elem()
		err := setReflectValue(value.Field(2), test.input)
		assert.NotNil(t, err, test.name)
	}
}

func TestGetJsonValuesHappyPath(t *testing.T) {
	testData := []struct {
		name   string
		input  string
		wanted map[string][]string
	}{
		{
			name:   "happy path 1",
			input:  `{}`,
			wanted: map[string][]string{},
		},
		{
			name:   "happy path 2",
			input:  `{"key1": "value1","key2": ["value2", "value3", "value4"]}`,
			wanted: map[string][]string{"key1": []string{"value1"}, "key2": []string{"value2", "value3", "value4"}},
		},
	}

	for _, test := range testData {
		r := strings.NewReader(test.input)
		res, err := getValuesFromJSON(r)

		assert.Nil(t, err, test.name)

		eq := reflect.DeepEqual(res, test.wanted)

		assert.Equal(t, true, eq, test.name)
	}
}

func TestGetJsonValuesError(t *testing.T) {
	testData := []struct {
		name  string
		input string
	}{
		{
			name:  "invalid json format 2",
			input: `{"key1", "value", "key2": "value2}`,
		},
		{
			name:  "invalid data (not string)",
			input: `{"key": 1}`,
		},
		{
			name:  "too complex format",
			input: `{"key1":"value1", "key2": {"skey1":"v1", "skey2":"v2"}}`,
		},
	}

	for _, test := range testData {
		r := strings.NewReader(test.input)
		_, err := getValuesFromJSON(r)

		assert.NotNil(t, err, test.name)
	}
}

type user1 struct {
	Name     string   `dto:"name" required:"false"`
	Age      int      `dto:"age" rquired:"true"`
	Hobby    []string `dto:"hobby" required:"false"`
	setItems map[string]bool
}

func (u *user1) AlreadySet(dtoName string) bool {
	_, ok := u.setItems[dtoName]
	return ok
}

func (u *user1) MarkSet(dtoName string) {
	if u.setItems == nil {
		u.setItems = make(map[string]bool)
	}
	u.setItems[dtoName] = true
}

func (u *user1) MarkAllUnset() {
	u.setItems = make(map[string]bool)
}

func TestParseGETRequestParamsHappyPath(t *testing.T) {
	testData := []struct {
		name   string
		url    string
		wanted *user1
	}{
		{
			name: "happy path 1",
			url:  "http://localhost:8080/users?age=28&hobby=sport&hobby=reading",
			wanted: &user1{
				Name:     "",
				Age:      28,
				Hobby:    []string{"sport", "reading"},
				setItems: map[string]bool{"age": true, "hobby": true},
			},
		},
	}

	for _, test := range testData {
		r, _ := http.NewRequest("GET", test.url, strings.NewReader(""))
		u := &user1{}
		ok := ParseRequestParams(r, u)
		assert.True(t, ok, test.name)
		assert.Equal(t, test.wanted.Age, u.Age, test.name)
		assert.Equal(t, test.wanted.Hobby, u.Hobby, test.name)
		assert.Equal(t, test.wanted.setItems, u.setItems, test.name)
	}
}

func TestParseGETRequestParamsError(t *testing.T) {
	testData := []struct {
		name string
		url  string
	}{
		{
			name: "not enough required params",
			url:  "http://localhost:8080/users?hobby=photography",
		},
	}

	for _, test := range testData {
		r, _ := http.NewRequest("GET", test.url, strings.NewReader(""))
		u := &user1{}
		ok := ParseRequestParams(r, u)
		assert.False(t, ok, test.name)
	}
}

func TestParsePostFormRequestParamsHappyPath(t *testing.T) {
	testData := []struct {
		name   string
		url    string
		data   url.Values
		wanted *user1
	}{
		{
			name: "happy path 1",
			url:  "http://localhost:8080/users",
			data: url.Values{"age": []string{"28"}, "hobby": []string{"sport", "photography"}},
			wanted: &user1{
				Name:     "",
				Age:      28,
				Hobby:    []string{"sport", "photography"},
				setItems: map[string]bool{"age": true, "hobby": true},
			},
		},
		{
			name: "happy path 2",
			url:  "http://localhost:8080/users?age=28",
			data: url.Values{"hobby": []string{"sport", "photography"}},
			wanted: &user1{
				Name:     "",
				Age:      28,
				Hobby:    []string{"sport", "photography"},
				setItems: map[string]bool{"age": true, "hobby": true},
			},
		},
	}
	for _, test := range testData {
		r, _ := http.NewRequest("POST", test.url, bytes.NewBufferString(test.data.Encode()))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		u := &user1{}
		ok := ParseRequestParams(r, u)
		assert.True(t, ok, test.name)
		assert.Equal(t, test.wanted.Age, u.Age, test.name)
		assert.Equal(t, test.wanted.Hobby, u.Hobby, test.name)
		assert.Equal(t, test.wanted.setItems, u.setItems, test.name)
	}
}

func TestParsePostFormRequestParamsError(t *testing.T) {
	testData := []struct {
		name string
		url  string
		data url.Values
	}{
		{
			name: "not enough request params",
			url:  "http://localhost:8080/users",
			data: url.Values{"hobby": []string{"sport", "photography"}},
		},
	}
	for _, test := range testData {
		r, _ := http.NewRequest("POST", test.url, bytes.NewBufferString(test.data.Encode()))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		u := &user1{}
		ok := ParseRequestParams(r, u)
		assert.False(t, ok, test.name)
	}
}

func TestParsePostJsonRequestParamsHappyPath(t *testing.T) {
	testData := []struct {
		name     string
		url      string
		jsonBody string
		wanted   *user1
	}{
		{
			name:     "happy path 1",
			url:      "http://localhost:8080/users",
			jsonBody: `{"age":"28","hobby":["sport","photography"]}`,
			wanted: &user1{
				Name:     "",
				Age:      28,
				Hobby:    []string{"sport", "photography"},
				setItems: map[string]bool{"age": true, "hobby": true},
			},
		},
		{
			name:     "happy path 2",
			url:      "http://localhost:8080/users?age=28",
			jsonBody: `{"hobby":["sport","photography"]}`,
			wanted: &user1{
				Name:     "",
				Age:      28,
				Hobby:    []string{"sport", "photography"},
				setItems: map[string]bool{"age": true, "hobby": true},
			},
		},
	}
	for _, test := range testData {
		r, _ := http.NewRequest("POST", test.url, bytes.NewBufferString(test.jsonBody))
		r.Header.Add("Content-Type", "application/json")
		u := &user1{}
		ok := ParseRequestParams(r, u)
		assert.True(t, ok, test.name)
		assert.Equal(t, test.wanted.Age, u.Age, test.name)
		assert.Equal(t, test.wanted.Hobby, u.Hobby, test.name)
		assert.Equal(t, test.wanted.setItems, u.setItems, test.name)
	}
}

func TestParsePostJsonRequestParamsError(t *testing.T) {
	testData := []struct {
		name     string
		url      string
		jsonBody string
	}{
		{
			name:     "not enough request params",
			url:      "http://localhost:8080/users",
			jsonBody: `{"hobby":["sport","photography"]}`,
		},
	}
	for _, test := range testData {
		r, _ := http.NewRequest("POST", test.url, bytes.NewBufferString(test.jsonBody))
		r.Header.Add("Content-Type", "application/json")
		u := &user1{}
		ok := ParseRequestParams(r, u)
		assert.False(t, ok, test.name)
	}
}
