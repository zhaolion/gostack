package cleanup

import (
	"testing"
)

type f1 struct{}

func (*f1) Close() {}

type f2 struct{}

func (*f2) Close() error { return nil }

type f3 struct{}

func (*f3) Flush() {}

type f4 struct{}

func (*f4) NoCloseOrFlush() {}

type TestStruct struct {
	F1 *f1
	F2 *f2
	F3 *f3
	F4 *f4
}

func TestRegisterStruct(t *testing.T) {
	var e Entry
	s := TestStruct{
		F1: &f1{},
		F2: &f2{},
		F3: &f3{},
		F4: &f4{},
	}
	e.Register(s)
}

func TestCleanup(t *testing.T) {
	var e Entry
	i, j := 0, 0
	func() {
		e.Register(func() { i += 1 })
		e.Register(func() { j += 2 })
		e.Run()
		e.Run() // multiple runs will be OK.
	}()
	if i != 1 || j != 2 {
		t.Errorf("Run() incorrect, want i = 1, j = 2, got i = %d, j = %d", i, j)
	}
}
