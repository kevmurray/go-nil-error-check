package main

import (
	"fmt"
)

// MyError is a custom error class
type MyError struct{}

func (o *MyError) Error() string { return "This is a mock error string" }

func main() {
	fmt.Printf("Error check for nil\n\n")
	var err error

	// baseline tests for basic errors

	trialReport("uninitialized err == nil", true, err == nil)

	err = fmt.Errorf("Error")
	trialReport("assigned err == nil", false, err == nil)

	err = nil
	trialReport("assigned to nil == nil", true, err == nil)

	// function tests

	// this section contains a "never true" warning because the compiler knows about this
	// problem and can tell this isn't right
	err = GetErrorPtrToError()
	trialReport("GetErrorPtrToError() == nil", false, err == nil)

	// this section contains a "never true" warning because the compiler knows about this
	// problem and can tell this isn't right
	err = GetErrorPtrToNil()
	trialReport("GetErrorPtrToNil() == nil", true, err == nil)

	err = GetErrorPtrToNilFixed1()
	trialReport("GetErrorPtrToNilFixed1() == nil", true, err == nil)

	err = GetErrorPtrToNilFixed2()
	trialReport("GetErrorPtrToNilFixed2() == nil", true, err == nil)

	// this section contains a "never true" warning because the compiler knows about this
	// problem and can tell this isn't right
	err = GetErrorPtrToNilNotFixed()
	trialReport("GetErrorPtrToNilNotFixed() == nil", true, err == nil)

	// struct tests

	// this section contains a "never true" warning because the compiler knows we are doing
	// something not right
	err = TestStruct{}.GetErrorPtrToNil()
	trialReport("struct.GetErrorPtrToNil() == nil", true, err == nil)

	// this section contains a "never true" warning because the compiler knows we are doing
	// something not right
	err = (&TestStruct{}).GetErrorPtrToNil()
	trialReport("(*struct).GetErrorPtrToNil() == nil", true, err == nil)

	err = TestStruct{}.GetErrorPtrToNilFixed()
	trialReport("struct.GetErrorPtrToNilFixed() == nil", true, err == nil)

	// interface tests

	// this section does not contain an "always true" warning because we're hiding the
	// implementation details behind the interface
	var testInterface TestInterface = TestStruct{}
	err = testInterface.GetErrorPtrToNil()
	trialReport("interface.GetErrorPtrToNil() == nil", true, err == nil)

	err = testInterface.GetErrorPtrToNilFixed()
	trialReport("interface.GetErrorPtrToNilFixed() == nil", true, err == nil)
}

func trialReport(name string, expected, actual bool) {
	var alert string
	if expected != actual {
		alert = "<-- surprising result?"
	}
	fmt.Printf("%-48s : expected %5t, actual %5t %s\n", name, expected, actual, alert)
}

/*
 * Test functions
 */

// GetErrorPtrToError returns an error result in a bad way - this works, but only because we always return
// an error
func GetErrorPtrToError() error {
	var err *MyError
	if true { // always executes
		err = &MyError{}
	}
	return err
}

// GetErrorPtrToNil returns no error in a way that makes a naive caller think an error occurred
// anyway. this is because what is returned is a variable with kind "pointer-to-MyError" and a
// value of nil, but to be "==" to nil, a variable must have a compatible kind (interface,
// struct, etc)
func GetErrorPtrToNil() error {
	var err *MyError
	if false { // never executes
		err = &MyError{}
	}
	return err
}

// GetErrorPtrToNilFixed1 fixes the problem by returning an explicit `nil`
func GetErrorPtrToNilFixed1() error {
	if false {
		return &MyError{}
	}
	return nil
}

// GetErrorPtrToNilFixed2 fixes the problem by always returning something with a "kind" of
// interface instead of a pointer to a structure
func GetErrorPtrToNilFixed2() error {
	var err error
	if false {
		err = &MyError{}
	}
	return err
}

// GetErrorPtrToNilNotFixed returns a nil error, but doesn't implement a full solution. It does
// explicitly return a `nil` value, but because the return type is `*MyError`, it gets coerced
// to a pointer-to-nil anyway
func GetErrorPtrToNilNotFixed() *MyError {
	if false { // never executes
		return &MyError{}
	}
	return nil
}

/*
 * Test interface and class
 */

type TestInterface interface {
	GetErrorPtrToNil() error
	GetErrorPtrToNilFixed() error
}

type TestStruct struct{}

func (o TestStruct) GetErrorPtrToNil() error {
	var err *MyError
	if false { // never executes
		err = &MyError{}
	}
	return err
}

func (o TestStruct) GetErrorPtrToNilFixed() error {
	var err error
	if false { // never executes
		err = &MyError{}
	}
	return err
}
