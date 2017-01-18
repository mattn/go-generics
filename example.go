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
