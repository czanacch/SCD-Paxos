package main

import (
	"sort"
	"sync"
	"encoding/json"
	"net/http"
	"time"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

/********************************************* PAXOS *******************************************/

var p *Proposer
var a *Acceptor
var l *Learner
var database *sql.DB // database in which the decided values are saved

var timer_databaseArrived int = 5
var timer_clockSignal int = 5 // How often the leader proposes value
var restored bool = false

var mutex_SetDatabase = sync.Mutex{}

var number_instance int = 0

// Quando un processo ha deciso un valore, resetta varie informazioni
func resetAfterDecision() {
	p.max_number_received = Number{Id_process: "", Num: 0}
	
	mapMutex.Lock()
	p.responded = make(map[string]bool)
	mapMutex.Unlock()

	a.v = 0
	a.n = Number{Id_process: "", Num: 0}
	a.max_n = Number{Id_process: "", Num: 0}

	l.decided = false
	
	mapMutex.Lock()		
	l.responded_values = make(map[int]int)
	mapMutex.Unlock()	
}

// Event that occurs when a user sends a sequence of commands to the process: it inserts the values sent in the UserValues database
func handleValues(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")	
	var data ActivationMessage
	_ = json.NewDecoder(r.Body).Decode(&data)

	for _,value := range data.Values {
			database, _ = sql.Open("sqlite3", "./database.db")
			statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS UserValues (value INTEGER)")
			statement.Exec()
			statement, _ = database.Prepare("INSERT INTO UserValues (value) VALUES (?)")
			statement.Exec(value)
	}
}

func sendValues() {

	for {
		value := 0
		copy := []int{}
		rows, _ := database.Query("SELECT * FROM UserValues")
		if rows != nil {
			for rows.Next() {
				rows.Scan(&value)
				copy = append(copy, value)
			}
		}
		sort.Ints(copy) // now the values contained in date.Values are in ascending order

		max_value := returnMax(l.decided_values) // takes the maximum of the values contained in l.decided_values
		for _, val := range copy {
			if val <= max_value {
				copy = copy[1:] // updates the list until it reaches a value greater than an already decided value
			}
		}

		for len(copy) > 0 && p.im_leader == true && (merge_invocated == true || len(machines) == 1) {	
			time.Sleep(time.Duration(timer_clockSignal) * time.Second)
			// Also in this case it is necessary to propose the maximum value [optimization]
			max_value = returnMax(l.decided_values) // takes the maximum of the values contained in l.decided_values
			for _, val := range copy {
				if val <= max_value {
					copy = copy[1:] // updates the list until it reaches a value greater than an already decided value
				}
			}
			if len(copy) > 0 {
				p.v, copy = copy[0], copy[1:] // Pop from the slice
			}
			p.prepareRequest() // do the Prepare Request

		}
		time.Sleep(1 * time.Millisecond)
	}
}