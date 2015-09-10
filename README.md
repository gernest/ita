## ita [![GoDoc](https://godoc.org/github.com/gernest/ita?status.svg)](https://godoc.org/github.com/gernest/ita)[![Coverage Status](https://coveralls.io/repos/gernest/ita/badge.svg?branch=master&service=github)](https://coveralls.io/github/gernest/ita?branch=master)[![Build Status](https://travis-ci.org/gernest/ita.svg)](https://travis-ci.org/gernest/ita)
ita is a Go(Golang) library that provides a clean API for dynamically calling structs methods.

The name ita is a swahili word for call.

## Motivation

I needed a clean API for calling struct methods, that will need to avoid unnecessary panics. I identified some of the the panics I was to avoid, and return errors instead.

* Call with too few input arguments. This happens when the method arguments are more than what are provided to the method.

* Call with too many input arguments. This happens when the method arguments are less than the values passed to the method

* using zero Value argument. This happens when one of the arguments has zero value or Invalid.

* call of nil function. This happens when the method is nil

So, for my use case I had no reason to panic on any of those scenarions. Instead I wanted to have an error raised so I can make other choices rather than recover from the panic( Did I say I almost fainted after a good dose of panics!)

So, Ita works with structs methods only. And its my intetion to try to give meaningful error messages rather than panicking.

But be warned, IT IS YOUR JOB TO CHECK FOR ERRORS

The latest error encountered is always available via the `Error()` method of `Struct` instance.


## features

* [x] call methods by name
* [X] easy access to results
* [x] support varidaic functions
* [x] reduced unnecessary panics, returns errors instead( Please see Motivation section)
* [ ] Method chaining (Coming soon)

## Installation

	go get github.com/gernest/ita

## How to use

```go
package main

import (
	"fmt"

	"github.com/gernest/ita"
)

type Hello struct{}

func (h *Hello) Say() string {
	return "hello"
}

func (h *Hello) Who(name string) string {
	return name
}
func (h *Hello) Double(word string) (string, string) {
	return word, word
}

func main() {
	hell := ita.New(&Hello{})

	// call Who
	who := hell.Call("Who", "gernest")
	
	// If you want to check errors
	if who.Error()!=nil{
		// do something
	}

	// call Say
	say := hell.Call("Say")

	// call double
	double := hell.Call("Double", "!")

	// get the returned results for calling Who
	name, _ := who.GetResults().First()

	// get the result returned by Say
	message, _ := say.GetResults().First()

	// get results returned by double
	doubleResult := double.GetResults()

	firstDouble, _ := doubleResult.First()
	lastDouble, _ := doubleResult.Last()

	fmt.Printf("%s %s %s%s", message, name, firstDouble, lastDouble)

}
```


You can see tests for more examples.




## Contributing

Start with clicking the star button to make the author and his neighbors happy. Then fork it and submit a pull request for whatever change you want to be added to this project.

Or Open an issue for any questions.

## Author
Geofrey Ernest <geofreyernest@live.com>

twitter  : [@gernesti](https://twitter.com/gernesti)

Facebook : [Geofrey Ernest](https://www.facebook.com/geofrey.ernest.35)


## Roadmap

*  Add support for assiging results to values provided by the user.

## Licence
This project is released under MIT licence see [LICENCE](LICENCE) for more details.