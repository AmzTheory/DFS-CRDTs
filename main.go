package main

// import (
// 		"os"
// 		"fmt"
// 		"strings"
// 		"bufio"
// 	)


type person struct {
	name string
	age  *int
}

func main() {
	//actual DFS running

	/*
		implement remove properly
			removed looks to be working properly with policies yet to be implemented
			used set implementation instead of slice

		implement interface layer
		test visualising
		implement OR set & map
		implement tree policies

			done work offline(locally)
		-------------------

		mulitple replicas
		communication

		Threading
			done with phase-1
		-----------------




	*/

	
	dfs := newDfs()
	dfs.start()
	dfs.updateAddHier("/", "1st", "txt")
	dfs.updateAddHier("/", "2nd", "txt")
	dfs.updateAddHier("/", "3rd", "txt")
	dfs.updateAddHier("/", "folder", "dir")
	dfs.updateAddHier("/folder/", "one", "txt")

	dfs.updateAddHier("/folder/", "folder2", "dir")
	dfs.updateAddHier("/folder/folder2/", "rand", "txt")

	dfs.updateRemoveHier("/1st","txt")



	dfs.ui.wait()
	// dfs.updateAddHier("/", "1st", "txt")

	// a:=new(*int)
	// // b:=new(*int)
	// st:=lls.New()

	// fmt.Println(a)
	// st.Push(a)
	// fmt.Println(st.Pop())

	// reader := bufio.NewReader(os.Stdin)
	// fmt.Println("Simple Shell")
	// fmt.Println("---------------------")

	// for {
	// 	fmt.Print("-> ")
	// 	text, _ := reader.ReadString('\n')
	// 	// convert CRLF to LF
	// 	text = strings.Replace(text, "\n", "", -1)

	// 	if strings.Compare("hi", text) == 0 {
	// 	fmt.Println("hello, Yourself")
	// 	}
	// }

}
func (p *person) start(age *int) {
	p.age = age
}
func (p *person) changeAge(age int) {
	(*p.age) = age
}
func doSomething() {

}
