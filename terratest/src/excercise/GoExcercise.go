package main

import (
	"fmt"
	"test"
)

var y int

const z = 20

func main() {

	var names = []string{"hello", "world", "kishor"}
	arrayPrint(names)
	/*var x int
	var x1=10
	//x2 int
	a:=20

	fmt.Println(x,x1,a)
	fmt.Printf("a is of type %T", a)
	a=21
	fmt.Println(a)
	fmt.Printf("a is of type %T", a)
	fmt.Println(y)
	fmt.Println(z)
	fmt.Printf("x is of type %T", x)*/
}

func arrayPrint(names []string) {
	for i := 0; i < len(names); i++ {
		fmt.Println(names[i])
	}

	fmt.Println(test.ProgressDeadlineSeconds)

}
