package main

import (
	"fmt"
	"os"
	filepath "path"
	"reflect"
	"strings"
)

type any interface{}

/**
 * Строка из переменной или из настроек
 * @type {[type]}
 */
func EnvOrConfOrDefault(env string, obj any, def string) string {
	if str := os.Getenv("GATE_HTTP_HOST"); len(strings.TrimSpace(str)) != 0 {
		return str
	}
	if obj != nil {
		return fmt.Sprint(obj)
	}
	return def
}

//function types
type mapFn func(interface{}) interface{}
type mapFnErr func(interface{}) (interface{}, error)

// func(value, memo) interface
type reduceFn func(interface{}, interface{}) interface{}
type filterFn func(interface{}) bool
type reduceFnErr func(interface{}, interface{}) (interface{}, error)
type filterFnErr func(interface{}) (bool, error)

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

//Map(slice, func)
func Map(in interface{}, fn mapFn) interface{} {
	val := reflect.ValueOf(in)
	out := make([]interface{}, 0)

	for i := 0; i < val.Len(); i++ {
		out = append(out, fn(val.Index(i).Interface()))
	}

	return out
}

//MapErr(slice, func)
func MapErr(in interface{}, fn mapFnErr) (interface{}, error) {
	val := reflect.ValueOf(in)
	out := make([]interface{}, 0)

	for i := 0; i < val.Len(); i++ {
		res, err := fn(val.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		out = append(out, res)
	}

	return out, nil
}

//ReduceErr(slice, starting value, func)
func ReduceErr(in interface{}, memo interface{}, fn reduceFnErr) (interface{}, error) {
	var err error
	result := memo
	val := reflect.ValueOf(in)

	for i := 0; i < val.Len(); i++ {
		result, err = fn(result, val.Index(i).Interface())
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

//Reduce(slice, starting value, func)
func Reduce(in interface{}, memo interface{}, fn reduceFn) interface{} {
	result := memo
	val := reflect.ValueOf(in)

	for i := 0; i < val.Len(); i++ {
		result = fn(result, val.Index(i).Interface())
	}

	return result
}

//HasElem(slice, elem value)
func HasElem(in interface{}, elem interface{}) bool {
	val := reflect.ValueOf(in)

	for i := 0; i < val.Len(); i++ {
		current := val.Index(i).Interface()
		if current == elem {
			return true
		}
	}

	return false
}

//HasElemFn(slice, elem value)
func HasElemFn(in interface{}, fn filterFn) bool {
	val := reflect.ValueOf(in)

	for i := 0; i < val.Len(); i++ {
		current := val.Index(i).Interface()
		if fn(current) {
			return true
		}
	}

	return false
}

//HasElemFnErr(slice, elem value)
func HasElemFnErr(in interface{}, fn filterFnErr) (bool, error) {
	val := reflect.ValueOf(in)

	for i := 0; i < val.Len(); i++ {
		current := val.Index(i).Interface()
		res, err := fn(current)
		if err != nil {
			return false, err
		}
		if res {
			return res, nil
		}
	}

	return false, nil
}

//FilterErr(slice, predicate func)
func FilterErr(in interface{}, fn filterFnErr) (interface{}, error) {
	val := reflect.ValueOf(in)
	out := make([]interface{}, 0, val.Len())

	for i := 0; i < val.Len(); i++ {
		current := val.Index(i).Interface()
		res, err := fn(current)
		if err != nil {
			return nil, err
		}
		if res {
			out = append(out, current)
		}
	}

	return out, nil
}

//Filter(slice, predicate func)
func Filter(in interface{}, fn filterFn) interface{} {
	val := reflect.ValueOf(in)
	out := make([]interface{}, 0, val.Len())

	for i := 0; i < val.Len(); i++ {
		current := val.Index(i).Interface()
		if fn(current) {
			out = append(out, current)
		}
	}

	return out
}
