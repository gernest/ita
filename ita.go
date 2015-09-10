// Package ita provides a clean API for dynamically calling structs methods
package ita

import (
	"errors"
	"reflect"
)

var (
	//ErrMethodNotFound is returned when there is no method found
	ErrMethodNotFound = errors.New("ita: method not found")

	//ErrZeroValue is returend when the vales is not valid
	ErrZeroValue = errors.New("ita: zero value argument")

	//ErrTooManyArgs is returned when arguments required are more than those passed to the method
	ErrTooManyArgs = errors.New("ita: the method has more arguments than the ones provided")

	//ErrTooFewArgs is returned when arguments required are less than those supplied to the method
	ErrTooFewArgs = errors.New("ita: you have provided more arguments than required by the method")
)

// Result stores the returned values from dynamically callling struct methods
type Result struct {
	rst []reflect.Value
}

// First returns the first result returned by caling a method
//
//For instance if you have method Foo
//	func (s *MyStruct)Foo()( a string, err error){...}
//
// First will return a but in interface form, meaning you will have to manualy
// cast the result to string. For automatic casting, please see Map, and MapTo
func (r *Result) First() (interface{}, bool) {
	if r.IsNil() {
		return nil, false
	}
	return r.Get(0)
}

// Get retrieves the result by Index, It will not panic if the index is out of range.
// instead it returns a nil value and ok is set to false.
func (r *Result) Get(index int) (val interface{}, ok bool) {
	if r.HasIndex(index) {
		return r.rst[index].Interface(), true
	}
	return nil, false
}

// IsNil checks if there is no result, or the result was nil
func (r *Result) IsNil() bool {
	return r.rst == nil
}

// Len returns the number of results. If a method returns tow values then Len()
// will return 2.
func (r *Result) Len() int {
	return len(r.rst)
}

// Last returns the last value returned by calling a method.
func (r *Result) Last() (interface{}, bool) {
	if r.IsNil() {
		return nil, false
	}
	return r.Get(r.Len() - 1)
}

// HasIndex returns true if the result has the given index.
func (r *Result) HasIndex(index int) bool {
	if r.IsNil() {
		return false
	}
	if index >= r.Len() {
		return false
	}
	return true
}

// Set stores rst in the Result object
func (r *Result) Set(rst []reflect.Value) {
	r.rst = rst
}

// Struct wraps the struct that we want to call its methods
type Struct struct {
	value  reflect.Value
	data   reflect.Value
	result *Result
	err    error
}

// Error returns the most recent error encountered while calling Call(). Be wise to
// remember to check this just incase to be sure  errors are under control.
func (s *Struct) Error() error {
	return s.err
}

// New wraps v into a Struct and returns the struct ready for plumbing.
func New(v interface{}) *Struct {
	value := reflect.ValueOf(v)
	var data reflect.Value
	switch value.Kind() {
	case reflect.Ptr:
		data = value
	default:
		data = reflect.New(reflect.TypeOf(v)).Elem()
	}
	return &Struct{
		value:  value,
		data:   data,
		result: &Result{},
		err:    nil,
	}
}

// Call calls method by nameand passing any arguments to the unserlying struct.
// This supports variadic methods and also methods that requires no arguments.
func (s *Struct) Call(method string, a ...interface{}) *Struct {
	clone := s.clone()
	m := clone.data.MethodByName(method)
	if !m.IsValid() {
		clone.err = ErrMethodNotFound
		return clone
	}
	typ := m.Type()
	args, err := getlArgs(typ, a)
	if err != nil {
		clone.err = err
		return clone
	}
	if typ.NumOut() > 0 {
		clone.callWithResults(m, args)
	} else {
		clone.call(m, args)
	}
	return clone
}

// callWIth results calls metho v with arguments args and stores the result
func (s *Struct) callWithResults(v reflect.Value, args []reflect.Value) {
	s.result.Set(v.Call(args))
}

// call calls method v with arguments args and discards the results
func (s *Struct) call(v reflect.Value, args []reflect.Value) {
	v.Call(args)
}

// getArgs retrieves argument list from args.
// This tries to avoid common panics, and returns the corresponding errors instead.
// panics can still occur but at least the ones related with arguments are mitigated.
func getlArgs(typ reflect.Type, args []interface{}) ([]reflect.Value, error) {
	var rst []reflect.Value
	for _, v := range args {
		val := reflect.ValueOf(v)
		if !val.IsValid() { // eliminate zero value panic
			return nil, ErrZeroValue
		}
		rst = append(rst, val)
	}
	switch {
	case typ.IsVariadic():
		if typ.NumIn() < len(args)-1 { // eliminate few argumenst panic
			return nil, ErrTooFewArgs
		}
	case !typ.IsVariadic() && typ.NumIn() > len(args): // eliminate too many arguments panic
		return nil, ErrTooManyArgs
	}
	return rst, nil
}

// clone creates a copy of the Struct. This way, references to the underlying object
// are kept intact just i case the resuts are still in use.
func (s *Struct) clone() *Struct {
	return &Struct{
		value:  s.value,
		data:   s.data,
		result: &Result{},
		err:    nil,
	}
}

// GetResults returns the results of the last Call.
func (s *Struct) GetResults() *Result {
	return s.result
}
