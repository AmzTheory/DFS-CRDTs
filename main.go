package main

import (
	"fmt"

	lls "github.com/emirpasic/gods/stacks/linkedliststack"
)

func main() {
	//actual DFS running
	st := lls.New()
	st.Push("Ahmed")
	st.Push("Saeed")
	s, _ := st.Pop()
	fmt.Println(s)
}
