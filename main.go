package main

// 		"os"

// 		"strings"
// 		"bufio"

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
)

type person struct {
	name string
	age  *int
}

func main() {

	dfs := newDfs()
	dfs.start()
	// dfs.updateAddHier("/", "1st", "txt")
	// dfs.updateAddHier("/", "2nd", "txt")
	// dfs.updateAddHier("/", "3rd", "txt")
	// dfs.updateAddHier("/", "folder", "dir")
	// dfs.updateAddHier("/folder/", "one", "txt")
	// dfs.updateAddHier("/folder/", "folder2", "dir")
	// dfs.updateAddHier("/folder/folder2/", "rand", "txt")
	// // dfs.updateRemoveHier("/1st", "txt")

	// //start the DFS with concurrency
	dfs.runAll()


}
func writeDB(){
	database, err := sql.Open("sqlite3", "./src/Dfs/data.db")
	checkErr(err)
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
	statement.Exec()
	checkErr(err)
    statement, err = database.Prepare("INSERT INTO people (firstname, lastname) VALUES (?, ?)")
	checkErr(err)
	
	statement.Exec("Nic", "Raboy")
	rows, err := database.Query("SELECT id, firstname, lastname FROM people")
	checkErr(err)
	
	var id int
    var firstname string
    var lastname string
    for rows.Next() {
        rows.Scan(&id, &firstname, &lastname)
        fmt.Println(firstname + " " + lastname)
	}
	rows.Close()
	
}