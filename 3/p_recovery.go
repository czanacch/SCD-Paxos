package main

import (
	"encoding/json"
	"net/http"
	"bytes"
	"time"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

/**** RECOVERY ****/

// Insert the values of the database DecidedValues into l.decided_values
func lookDatabase() {
	value := 0
	database, _ = sql.Open("sqlite3", "./database.db")
	rows, _ := database.Query("SELECT * FROM DecidedValues")
	if rows != nil {
		for rows.Next() {
			rows.Scan(&value)
			l.decided_values = append(l.decided_values, value)
		}
	}
}

// Send a database request message to all machines
func askInformation() {
	for _,addr := range machines {
		if addr != my_address {
			jsonData := Notification{
				Address: my_address }
			jsonValue, _ := json.Marshal(jsonData)
			go http.Post("http://"+addr+"/getDatabase", "application/json", bytes.NewBuffer(jsonValue))
		}
	}
	time.Sleep(5 * time.Second)
	if restored == false {
		timer_databaseArrived = 5
		go decrement_timer()
	}
}

// Keep asking for databases 
func decrement_timer() {
    for timer_databaseArrived > 0 && restored == false {
            timer_databaseArrived = timer_databaseArrived - 1
            time.Sleep(1 * time.Second) // Aspetta un secondo
    }
	if timer_databaseArrived == 0 {
		askInformation()
	}
}

// Event in response to the "database request" message
func handleGetDatabase(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var data Notification
	_ = json.NewDecoder(r.Body).Decode(&data)

	jsonData := RecoveryMessage{
		Decided_values: l.decided_values}
	jsonValue, _ := json.Marshal(jsonData)
	go http.Post("http://"+data.Address+"/setDatabase", "application/json", bytes.NewBuffer(jsonValue))
}

// Event in response to the "database arrival" message
func handleSetDatabase(w http.ResponseWriter, r *http.Request) {
	mutex_SetDatabase.Lock()
	if restored == false {
		w.Header().Set("Content-Type", "application/json")
		var data RecoveryMessage
		_ = json.NewDecoder(r.Body).Decode(&data)
			
		for _,value := range data.Decided_values {
			if !contains(l.decided_values, value) { // se il valore non è nel mio database
				l.decided_values = append(l.decided_values, value)
				database, _ = sql.Open("sqlite3", "./database.db")
				statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS DecidedValues (value INTEGER)")
				statement.Exec()
				statement, _ = database.Prepare("INSERT INTO DecidedValues (value) VALUES (?)")
				statement.Exec(value)
			}
		}
		restored = true
		// End of "Recovery state"
		a.mutex_PrepareRequest.Unlock()
		a.mutex_AcceptRequest.Unlock()
		l.mutex_AcceptResponse.Unlock()
	}
	mutex_SetDatabase.Unlock()
}