package main

import (
	crdt "CRDTsGO"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	_ "github.com/mattn/go-sqlite3"
)

var dbPath string
var data string

//Structs and type
type RepElem struct {
	Name        string
	ElemType       string
}

type contentMap map[RepElem]string

type RepLayer struct {
	dfs    *Dfs
	or     *crdt.ORSet
	cmap   contentMap
}

//initalisation
func newRepLayer(id int) *RepLayer {
	dbPath = "./src/DFS/data.db"
	data = "data"
	dic,or :=readDB(id)



	l := RepLayer{
		dfs:  new(Dfs),
		or:   or,
		cmap: dic,
	}

	return &l
}
func (l *RepLayer) setDfs(dfs *Dfs) {
	l.dfs = dfs
}

//run locally & remotely comunicate theire messages to pushupState, which in turn get excuted by pushupstate, then passed to upper layers

func (l *RepLayer) runLocally(send chan RemoteMsg, recieve chan HierToRep) {
	for {
		msg := <-recieve
		var el, u interface{}
		el = RepElem{Name: msg.path, ElemType: msg.fileType}
		if msg.op == "add" {
			u = l.or.SrcAdd(el)
		} else if msg.op == "rm" {
			u = l.or.SrcRemove(el)
		}
		// rmsg:=RemoteMsg{SenderID:-1,Op:msg.op,Params:[]string{msg.path,msg.fileType},}
		rmsg := RemoteMsg{SenderID: -1, Op: msg.op,P1:el,P2:u,cancel:msg.cancel,}
		// fmt.Println(u)
		send <- rmsg
		go l.dfs.sendRemote(rmsg) //TODO: might miss local execution
		// send <- l.returnCurrentSet() //send the updated set to hier
	}
}

//execute operation and update hier
func (l *RepLayer) executeOp(send chan map[*RepElem]string, recieve chan RemoteMsg) {
	send <- l.returnCurrentSet() //send the initial state
	var opMsg RemoteMsg
	var el, u interface{}
	var r []interface{}
	for { //wait for operation to be executed local/remotely
		opMsg = <-recieve
		// fmt.Println(opMsg)
		if opMsg.Op == "add" {
			// el = opMsg.Params[0]
			// u = opMsg.Params[1]
			el=opMsg.P1
			u=opMsg.P2
			l.add(el.(RepElem), u.(string))
		} else if opMsg.Op == "rm" {
			// el = opMsg.Params[0]
			// u = opMsg.Params[1]
			el=opMsg.P1
			u=opMsg.P2
			r = u.([]interface{})
			l.remove(r, el)
		}

		if(opMsg.cancel!=nil){opMsg.cancel()}
		send <- l.returnCurrentSet()
	}
}

//listen remotely
func (l *RepLayer) runRemotely(send chan RemoteMsg, recieve chan RemoteMsg) {
	var rmsg RemoteMsg
	for {
		rmsg = <-recieve

		send <- rmsg
	}
}

//update inteface

func (l *RepLayer) add(el RepElem, u string) {
	l.or.Add(u, el)
	if _, ok := l.cmap[el]; !ok { //start with an empty content in the case it hasn't be creatd
		l.cmap[el] = ""
	}

}

func (l *RepLayer) remove(r []interface{}, el interface{}) {
	l.or.Remove(r, el)
}


func (l *RepLayer) returnCurrentSet() map[*RepElem]string {
	temp := make(map[*RepElem]string)
	for _, k := range l.or.Values() {
		kk := (k.(RepElem))
		temp[&kk] = l.cmap[kk]
	}
	return temp
}

//read the databse
func readDB(id int) (contentMap, *crdt.ORSet) {
	or := crdt.NewORSet()
	contentMap := make(map[RepElem]string)

	database, err := sql.Open("sqlite3", dbPath)
	checkErr(err)

	rows, err := database.Query("SELECT path,type,content,used,token from " + data + " where dfsId=" + strconv.Itoa(id))
	checkErr(err)
	var path, content, ElemType, token string
	var used int
	for rows.Next() {
		rows.Scan(&path, &ElemType, &content, &used, &token)
		el := RepElem{Name: path, ElemType: ElemType}
		contentMap[el] = content
		if used == 1 {

			///Add for each token into OR
			for _, v := range strings.Split(token, ",") {
				or.Add(v, el)
			}
		}

	}
	rows.Close()
	return contentMap, or
}

func (l *RepLayer) writeDB() {
	database, _ := sql.Open("sqlite3", dbPath)
	//delete dfs records
	id := (*l.dfs).id
	statement, err := database.Prepare("delete from " + data + " where dfsID=" + strconv.Itoa(id))
	statement.Exec()
	checkErr(err)
	// //creating data table
	// statement, err = database.Prepare("create table " + data + " (id INTEGER PRIMARY KEY, path TEXT ,type TEXT ,content TEXT,used INTEGER,dfsID Integer)")
	// statement.Exec()

	statement, err = database.Prepare("INSERT INTO " + data + " (path, type,content,used,dfsID,token) VALUES (?,?,?,?,?,?)")

	var used int
	var token string
	//insert data points
	for k, v := range l.cmap {
		used = 0
		token = ""
		if (*l.or).Contains(k) {
			used = 1
			token =getTokens(l.or,k)
		}

		checkErr(err)
		statement.Exec(k.Name, k.ElemType, v, used, l.dfs.id, token)

	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

//GetTokens get set of tokens associated with specific Element in the OR.set
func getTokens(or *crdt.ORSet, el interface{}) string {
	items := []string{}

	it := or.Get(el)
	for _,v:=range it.Values() {
		items = append(items, fmt.Sprintf("%v", v))
	}

	return strings.Join(items, ", ")
}

