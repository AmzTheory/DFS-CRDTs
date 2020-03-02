package main

// "fmt"
// lls "github.com/emirpasic/gods/stacks/linkedliststack"

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


	dfs.updateAddHier("/", "1st", "txt")

	// a:=new(*int)
	// // b:=new(*int)
	// st:=lls.New()

	// fmt.Println(a)
	// st.Push(a)
	// fmt.Println(st.Pop())

	dfs.UpdateRemoveReplication("/folder", "dir")
	// dfs.rep.printCurrentState()
	dfs.hier.printCurrentState()
}
func (p *person) start(age *int) {
	p.age = age
}
func (p *person) changeAge(age int) {
	(*p.age) = age
}
func doSomething() {

}
