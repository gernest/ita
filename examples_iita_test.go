package ita

import (
	"fmt"
)

type Hello struct{}

func (h *Hello) Say() string {
	return "hello"
}

func (h *Hello) Who(say string) string {
	return say
}

func ExampleStruct() {
	hell := New(&Hello{})

	who := hell.Call("Who", "gernest")

	say := hell.Call("Say")

	first, _ := say.GetResults().First()
	second, _ := who.GetResults().First()
	fmt.Println(first, second)

	//Output:
	//hello gernest

}
