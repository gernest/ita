package ita

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"
)

var hello = "hello"

type Sample struct {
	out io.Writer
}

func (s *Sample) HelloPoint() string {
	return hello
}
func (s Sample) Hello() string {
	return hello
}

func (s *Sample) HelloBlanc() {
	s.out.Write([]byte(hello))
}

func (s *Sample) HelloArgs(first, second string) []string {
	return []string{first, second}
}

func (s *Sample) HelloArgsBlank(first, second string) {
	s.out.Write([]byte(first + second))
}
func (s *Sample) HelloVariad(a ...string) []string {
	return a
}
func (s *Sample) HelloVariadBlank(a ...string) {
	for _, v := range a {
		s.out.Write([]byte(v))
	}
}
func (s *Sample) HelloMix(first string, a ...string) []string {
	a = append(a, first)
	return a
}

func (s *Sample) HelloMixBlank(first string, a ...string) {
	s.out.Write([]byte(first))
	for _, v := range a {
		s.out.Write([]byte(v))
	}
}

func TestPtrRcvr(t *testing.T) {
	buf := &bytes.Buffer{}

	obj := New(&Sample{out: buf})

	if obj.Error() != nil {
		t.Errorf("expecetd nil got %v", obj.Error())
	}

	// Check the calling to pointer receiver an normal receiver
	for _, v := range []string{"Hello", "HelloPoint"} {
		obj = obj.Call(v)
		if obj.Error() != nil {
			t.Errorf("expecetd nil got %v ", obj.Error())
		}

		rst := obj.GetResults()
		if rst.IsNil() {
			t.Errorf("expected results for %s got nil", v)
		}

		first, ok := rst.First()
		if !ok {
			t.Errorf("expeced returned value")
		}

		if !reflect.DeepEqual(hello, first) {
			t.Errorf("expecetd %s got %v", hello, first)
		}
	}

	obj = obj.Call("HelloBlanc")
	if obj.Error() != nil {
		t.Error(obj.Error())
	}
	rst := obj.GetResults()
	if !rst.IsNil() {
		t.Errorf("expected no result got %v instead", rst)
	}
	first, ok := rst.First()
	if ok {
		t.Errorf("expected zero results got %v instead", first)
	}
	if buf.String() != hello {
		t.Errorf("expected %s got %s", hello, buf.String())
	}
	buf.Reset()

	args := []struct {
		name          string
		first, second string
		blank         bool
	}{
		{"HelloArgs", "hello", "world", false},
		{"HelloArgsBlank", "hello", "world", true},
	}

	for _, v := range args {
		obj = obj.Call(v.name, v.first, v.second)
		if obj.Error() != nil {
			t.Errorf("expecetd nil got %v ", obj.Error())
		}
		if v.blank {
			rst := obj.GetResults()
			if !rst.IsNil() {
				t.Error("expected nil got results instead")
			}

			if buf.String() != v.first+v.second {
				t.Errorf("expecetd %s got %s", v.first+v.second, buf.String())
			}
			buf.Reset()
			continue

		}
		rst := obj.GetResults()
		if rst.IsNil() {
			t.Errorf("expected results for %v got nil", v)
		}

		first, ok := rst.First()
		if !ok {
			t.Errorf("expeced returned value")
		}

		if !reflect.DeepEqual([]string{v.first, v.second}, first) {
			t.Errorf("expecetd %s got %v", hello, first)
		}
	}

	variads := []struct {
		name  string
		args  []string
		blank bool
	}{
		{"HelloVariad", nil, false},
		{"HelloVariad", []string{"hello", "world"}, false},
		{"HelloVariadBlank", nil, true},
		{"HelloVariadBlank", []string{"hello", "world"}, true},
	}

	for _, v := range variads {
		var vArgs []interface{}
		if v.args != nil {
			for _, argv := range v.args {
				vArgs = append(vArgs, argv)
			}
			obj = obj.Call(v.name, vArgs...)
		} else {
			obj = obj.Call(v.name)
		}
		if obj.Error() != nil {
			t.Errorf("expecetd nil got %v ", obj.Error())
		}
		if v.blank {
			rst := obj.GetResults()
			if !rst.IsNil() {
				t.Error("expected nil got results instead")
			}

			if v.args != nil {
				if !strings.Contains(buf.String(), v.args[0]) {
					t.Errorf("expeced %s to contain %s", buf.String(), v.args[0])
				}
			}
			buf.Reset()
			continue
		}
		rst := obj.GetResults()
		if rst.IsNil() {
			t.Errorf("expected results for %v got nil", v)
		}

		first, ok := rst.First()
		if !ok {
			t.Errorf("expeced returned value")
		}
		if v.args != nil {
			if !reflect.DeepEqual(v.args, first) {
				t.Errorf("expected %v got %v", v.args, first)
			}
			if rst.Len() != 1 {
				t.Errorf("expected 1 got %d", rst.Len())
			}

			f, ok := rst.Get(0)
			if !ok {
				t.Error("expected result")
			}
			if !reflect.DeepEqual(v.args, f) {
				t.Errorf("expected %v got %v", v.args, first)
			}

			f, ok = rst.Get(1)
			if ok {
				t.Error("expected no result")
			}
			if reflect.DeepEqual(v.args, f) {
				t.Errorf("expected nil got %v", first)
			}
		}

	}
}
