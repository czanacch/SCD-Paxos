package main

import (
	"sync"
	"encoding/json"
	"net/http"
	"fmt"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

/*********** PAXOS LEARNER **************/

type Learner struct {
	responded_values map[int]int // mapping from values to number of times the value arrived

	decided bool
	decided_values []int // slice containing all the values decided so far (it's a database copy)

	mutex_AcceptResponse sync.Mutex
	channel_AcceptResponse chan AcceptResponse
}

func (l *Learner) constructorLearner() {
	l.responded_values = make(map[int]int)
	l.decided = false
	l.decided_values = []int{}

	l.mutex_AcceptResponse = sync.Mutex{}
	l.channel_AcceptResponse = make(chan AcceptResponse)

}

func handleAcceptResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var acceptResponse AcceptResponse
	_ = json.NewDecoder(r.Body).Decode(&acceptResponse)

	go manageAcceptResponse()
	l.channel_AcceptResponse <- acceptResponse
}

func manageAcceptResponse() {
	l.mutex_AcceptResponse.Lock()
	acceptResponse := <-l.channel_AcceptResponse
	mapMutex.Lock()
	l.responded_values[acceptResponse.V] = l.responded_values[acceptResponse.V] + 1
	
	for v,num := range l.responded_values {
		
		if num >= len(machines)/2+1 && l.decided == false && contains(l.decided_values,v) == false && (len(l.decided_values) == 0 || v >= l.decided_values[len(l.decided_values)-1]) {


			fmt.Println("Learner ",my_address,": Chosen ",v)
			l.decided_values = append(l.decided_values, v)
			database, _ = sql.Open("sqlite3", "./database.db")
			statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS DecidedValues (value INTEGER)")
			statement.Exec()
			statement, _ = database.Prepare("INSERT INTO DecidedValues (value) VALUES (?)")
			statement.Exec(v)
			l.decided = true
			break

		}
	}
	mapMutex.Unlock()
	
	if l.decided == true {
		resetAfterDecision()
	}
	l.mutex_AcceptResponse.Unlock()
}