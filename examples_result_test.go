package ita

import (
	"fmt"
	"reflect"
)

func ExampleResult_First() {
	var rst []reflect.Value
	out := []int{1, 2, 3}
	for _, o := range out {
		rst = append(rst, reflect.ValueOf(o))
	}
	result := &Result{}
	result.Set(rst)

	fmt.Println(result.First())

	//Output:
	//1 true
}

func ExampleResult_Get() {
	var rst []reflect.Value
	out := []int{1, 2, 3}
	for _, o := range out {
		rst = append(rst, reflect.ValueOf(o))
	}
	result := &Result{}
	result.Set(rst)

	fmt.Println(result.Get(0))
	fmt.Println(result.Get(5))

	//Output:
	//1 true
	//<nil> false
}
