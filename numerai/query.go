package numerai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/xerrors"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

const numeraiUrl = "https://api-tournament.numer.ai"

type QueryArgs map[string]interface{}
type QueryResult reflect.Value

func RawQuery(query string, args QueryArgs) (r QueryResult, err error) {
	requestBytes := bytes.Buffer{}
	request := struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables"`
	}{
		Query:     query,
		Variables: args,
	}
	if err = json.NewEncoder(&requestBytes).Encode(request); err != nil {
		return
	}
	rqst, err := http.NewRequest(http.MethodPost, numeraiUrl, &requestBytes)
	if err != nil {
		return
	}
	rqst.Header.Set("Content-Type", "application/json; charset=utf-8")
	rqst.Header.Set("Accept", "application/json; charset=utf-8")
	res, err := http.DefaultClient.Do(rqst)
	if err != nil {
		return
	}
	defer res.Body.Close()
	bf := bytes.Buffer{}
	if _, err = io.Copy(&bf, res.Body); err != nil {
		return
	}
	x := map[string]interface{}{}
	if err = json.Unmarshal(bf.Bytes(), &x); err != nil {
		return
	}
	if l, ok := x["errors"]; ok {
		fmt.Println(l)
		q := l.([]interface{})[0].(map[string]interface{})
		err = xerrors.Errorf("GraphQL error:", q["message"].(string))
		return
	}
	r = QueryResult(reflect.ValueOf(x))
	return
}

func (q QueryResult) Q(a interface{}) QueryResult {
	v := reflect.Value(q)
	for v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if v.Kind() == reflect.Map {
		return QueryResult(v.MapIndex(reflect.ValueOf(a)))
	} else if v.Kind() == reflect.Slice {
		return QueryResult(v.Index(int(reflect.ValueOf(a).Int())))
	}
	return QueryResult{}
}

func (q QueryResult) V(a interface{}) (r reflect.Value) {
	v := reflect.Value(q)
	for v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if v.Kind() == reflect.Map {
		r = v.MapIndex(reflect.ValueOf(a))
	} else if v.Kind() == reflect.Slice {
		r = v.Index(int(reflect.ValueOf(a).Int()))
	}
	for r.Kind() == reflect.Interface {
		r = r.Elem()
	}
	return
}

func (q QueryResult) List() []QueryResult {
	v := (reflect.Value)(q)
	for v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if v.Kind() == reflect.Slice {
		r := make([]QueryResult, v.Len(), v.Len())
		for i := range r {
			r[i] = QueryResult(v.Index(i))
		}
		return r
	}
	panic("is not list")
}

func (q QueryResult) Chan() chan QueryResult {
	v := (reflect.Value)(q)
	for v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if v.Kind() == reflect.Slice {
		c := make(chan QueryResult)
		go func() {
			defer close(c)
			for i := 0; i < v.Len(); i++ {
				c <- QueryResult(v.Index(i))
			}
		}()
		return c
	}
	panic("is not list")
}

func (q QueryResult) String(a interface{}) string {
	v := q.V(a)
	if !v.IsValid() {
		return ""
	}
	return v.String()
}

func (q QueryResult) Real(a interface{}) float32 {
	v := q.V(a)
	if !v.IsValid() {
		return 0
	}
	switch v.Kind() {
	case reflect.String:
		f, err := strconv.ParseFloat(v.String(), 32)
		if err != nil {
			panic(err.Error())
		}
		return float32(f)
	case reflect.Float64, reflect.Float32:
		return float32(v.Float())
	default:
		return float32(v.Int())
	}
}

func (q QueryResult) Float(a interface{}) float64 {
	v := q.V(a)
	if !v.IsValid() {
		return 0
	}
	switch v.Kind() {
	case reflect.String:
		f, err := strconv.ParseFloat(v.String(), 32)
		if err != nil {
			panic(err.Error())
		}
		return f
	case reflect.Float64, reflect.Float32:
		return v.Float()
	default:
		return float64(v.Int())
	}
}

func (q QueryResult) Int(a interface{}) int {
	v := q.V(a)
	if !v.IsValid() {
		return 0
	}
	switch v.Kind() {
	case reflect.String:
		i, err := strconv.ParseInt(v.String(), 10, 32)
		if err != nil {
			panic(err.Error())
		}
		return int(i)
	case reflect.Float32, reflect.Float64:
		return int(v.Float())
	}
	return int(v.Int())
}

func (q QueryResult) Time(a interface{}) time.Time {
	v := q.V(a)
	if !v.IsValid() {
		return time.Time{}
	}
	if v.Kind() == reflect.String {
		t, err := time.Parse(time.RFC3339, v.String())
		if err != nil {
			panic(err.Error())
		}
		return t
	}
	panic("not string")
}

func (q QueryResult) Fill(a interface{}) {
	av := reflect.ValueOf(a).Elem() // a is a ptr to struct
	at := av.Type()
	for i := 0; i < at.NumField(); i++ {
		f := av.Field(i)
		ft := at.Field(i)
		switch f.Kind() {
		case reflect.String:
			f.Set(reflect.ValueOf(q.String(ft.Name)))
		case reflect.Float32:
			f.Set(reflect.ValueOf(q.Real(ft.Name)))
		case reflect.Float64:
			f.Set(reflect.ValueOf(q.Float(ft.Name)))
		case reflect.Int, reflect.Int32, reflect.Int64:
			f.Set(reflect.ValueOf(q.Int(ft.Name)))
		default:
			if f.Type() == reflect.TypeOf(time.Time{}) {
				f.Set(reflect.ValueOf(q.Time(ft.Name)))
			} else {
				panic("unsupported type of filed " + ft.Name)
			}
		}
	}
}
