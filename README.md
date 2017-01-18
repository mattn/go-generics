# go-generics 

This is very very experimental

## Usage

```
$ gen example.go > out.go
```

### Input
```go
package main

type T generics

func foo(a, b T) T {
	if a < b {
		return b
	}
	return a
}

type F struct {
}

func main() {
	println(foo[T](1, 2))
	println(foo(1.2, 2))
	println(foo("foo", "bar"))
	println(foo("foo", F{}))
	println(foo("foo", &F{}))
}
```

### Output

```go
package main

type T interface{}

func foo(a, b T) T

type F struct {
}

func main() {
	println(foo[T](1, 2))
	println(foo_of_float64_int64(1.2, 2))
	println(foo_of_string_string("foo", "bar"))
	println(foo_of_string_F("foo", F{}))
	println(foo_of_string__F("foo", &F{}))
}
func foo_of_float64_int64(a float64, b int64) T {
	if a < b {
		return b
	}
	return a
}
func foo_of_string_string(a string, b string) T {
	if a < b {
		return b
	}
	return a
}
func foo_of_string_F(a string, b F) T {
	if a < b {
		return b
	}
	return a
}
func foo_of_string__F(a string, b *F) T {
	if a < b {
		return b
	}
	return a
}
```
