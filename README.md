# go-nil-error-check

We recently encountered a problem where we were checking for no errors with the standard Go
check like

```golang
result, err := thing.doAThing()
if err != nil {
  panic("Doing the thing failed")
}
```

And to our surprise, the code paniced even though we were returning nil from the `doAThing`
function.

After reading [this
article](https://glucn.medium.com/golang-an-interface-holding-a-nil-value-is-not-nil-bb151f472cc7),
this made a little more sense, but I wanted to understand better how and when this problem
occurs. This project is just a toy project for me to explore when returning a `nil` error gets
mis-classified as not a `nil` error. Hopefully it can also help you understand this better too.

I suggest reading the article above, but to sum up: Go interfaces (like `error`) have a "kind"
and a "value". A `nil` interface has both a kind and a value of `nil`. In order for the equality
test to pass, both the kind and the value have to match (i.e be `nil`).

This is my understanding of how some things work:

| Value in code       | Kind               | Value               |
| ------------------- | ------------------ | ------------------- |
| nil                 | nil                | nil                 |
| var err error       | nil                | nil                 |
| var err error = nil | nil                | nil                 |
| errors.New("xx")    | ptr-to-errorString | errorString{s:"xx"} |
| MyError{}           | MyError            | MyError{}           |
| &MyError{}          | ptr-to-MyError     | MyError{}           |
| (\*MyError)(nil)    | ptr-to-MyError     | nil                 |

The last entry above is the one that causes problems. If you have a function like

```golang
func doAThing() *MyError {
  ... do some stuff that never sets a value in err
  return nil
}
```

or

```golang
func doAThing() error {
  var err *MyError
   ... do some stuff that never sets a value in err
  return err
}
```

then what is returned from the function is a value with a kind of "ptr-to-MyError" and value of
`nil`, but since `err == nil` is only true when both the kind and the value are `nil`, this will
fail the test even though it seems like it shouldn't if you just look at the value.

There are two options to fix this.

First is to fix the function so it just returns and uses `error` instead of `*MyError`. This
code works correctly:

```golang
func doAThing() error {
  if condition {
    return &MyError{}
  }
  return nil
}
```

because you are never storing the `nil` value into a `*MyError` variable type before returning
it.

Alternatively, you can explicitly check to see if the value is a pointer kind that contains a
value of `nil`, like

```golang
// check for no error
result, err := thing.doAThing()
if err == nil ||
  (reflect.ValueOf(err).Kind() == reflect.Ptr &&
  reflect.ValueOf(err).IsNil()) {
  ... no error occurred
}

// check for an error
result, err := thing.doAThing()
if err != nil &&
  (reflect.ValueOf(err).Kind() != reflect.Ptr ||
  !reflect.ValueOf(err).IsNil()) {
  ... an error occurred
}
```

Obviously the first options is better, and can be summed up by just following the Go standard
practices and not doing anything weird - if you have an error, just return an `error` type,
don't try to get all fancy with your return values.

If some of your callers need to do special things with your custom error types (like logging
different information from custom errors than default errorStrings), then make them responsible
for type checking using code like

```golang
result, err := thing.doAThing()
if err != nil {
  if myErr, ok := (*MyError)(err); ok {
    logMySpecialError(myErr)
  } else {
    logSystemError(err)
  }
}
```

(Although I'd argue strongly that if you have to do something like this then something is wrong
with the implementation of your logging or error package.)

---

You can run the code in `main.go` with `go run .` and the output would be as follows. (You'll
have to look at the code in `main.go` to see the details of what is tested for each condition,
though.)

```
Error check for nil

uninitialized err == nil                         : expected  true, actual  true
assigned err == nil                              : expected false, actual false
assigned to nil == nil                           : expected  true, actual  true
GetErrorPtrToError() == nil                      : expected false, actual false
GetErrorPtrToNil() == nil                        : expected  true, actual false <-- surprising result?
GetErrorPtrToNilFixed1() == nil                  : expected  true, actual  true
GetErrorPtrToNilFixed2() == nil                  : expected  true, actual  true
GetErrorPtrToNilNotFixed() == nil                : expected  true, actual false <-- surprising result?
struct.GetErrorPtrToNil() == nil                 : expected  true, actual false <-- surprising result?
(*struct).GetErrorPtrToNil() == nil              : expected  true, actual false <-- surprising result?
struct.GetErrorPtrToNilFixed() == nil            : expected  true, actual  true
interface.GetErrorPtrToNil() == nil              : expected  true, actual false <-- surprising result?
interface.GetErrorPtrToNilFixed() == nil         : expected  true, actual  true
```
