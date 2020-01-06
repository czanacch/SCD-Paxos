package main


import (
	"math"
	"encoding/json"
	"net/http"
	"bytes"
	"time"
	_ "github.com/mattn/go-sqlite3"
)

/* Functions Leader Election */

// Reset of p_i
func Recovery() {
	stateMutex.Lock()
	s = "Election"
	counter = counter + 1
	g = Number{Id_process: my_address, Num: counter} // New group identifier
	c = my_address
	p.im_leader = true
	merge_invocated = false
	mapMutex.Lock()
	Up = make(map[string]bool)
	mapMutex.Unlock()
	s = "Normal"
	stateMutex.Unlock()
}

// Create a new group where p_i is the leader with all the members of the CoordinatorSet set inside
func Merge(CoordinatorSet map[string]bool) {	
	stateMutex.Lock()
	s = "Election"
	counter = counter + 1
	g = Number{Id_process: my_address, Num: counter} // New group identifier
	c = my_address
	p.im_leader = true
	TempSet := Up
	mapMutex.Lock()
	Up = make(map[string]bool)
	mapMutex.Unlock()
	stateMutex.Unlock()
	
	for j,_ := range CoordinatorSet {
		jsonData := InvitationMessage{
			Sender: my_address,
			Gn: g}
		jsonValue, _ := json.Marshal(jsonData)
		go http.Post("http://"+j+"/InvitationMessage", "application/json", bytes.NewBuffer(jsonValue)) // Dì esplicitamente a P_j di unirsi a g
		mapMutex.Lock()
		Up[j] = true
		mapMutex.Unlock()
	}
	for j,_ := range TempSet {
		jsonData := InvitationMessage{
			Sender: my_address,
			Gn: g}
		jsonValue, _ := json.Marshal(jsonData)
		go http.Post("http://"+j+"/InvitationMessage", "application/json", bytes.NewBuffer(jsonValue))
		mapMutex.Lock()
		Up[j] = true
		mapMutex.Unlock()
	}
	stateMutex.Lock()
	s = "Reorganization"
	stateMutex.Unlock()
	
	for j,_ := range Up {
		jsonData := ReadyMessage{
			Sender: my_address,
			Gn: g}
		jsonValue, _ := json.Marshal(jsonData)
		go http.Post("http://"+j+"/ReadyMessage", "application/json", bytes.NewBuffer(jsonValue))

		for T_readyMessage > 0 && mapReadyMessage[j] == false {
			T_readyMessage = T_readyMessage - 1
			time.Sleep(1 * time.Second)
		}
		if T_readyMessage == 0 && mapReadyMessage[j] == false {
			T_readyMessage = timer_T
			Recovery()
		} else {
			T_readyMessage = timer_T
			mapMutex.Lock()
			mapReadyMessage[j] = false
			mapMutex.Unlock()
		}
	}
	stateMutex.Lock()
	s = "Normal"
	stateMutex.Unlock()

	merge_invocated = true
}

// If p_i hasn't been receiving messages from the leader for a while...
func Timeout() {
	stateMutex.Lock()
	MyCoord := c
	MyGroup := g
	stateMutex.Unlock()
	if MyCoord != my_address {
		jsonData := AreYouThereMessage{
			Sender: my_address,
			Gn: MyGroup}
		jsonValue, _ := json.Marshal(jsonData)
		go http.Post("http://"+MyCoord+"/AreYouThereMessage", "application/json", bytes.NewBuffer(jsonValue))

		for T_areYouThereMessage > 0 && mapAreYouThereMessage[MyCoord] == false {
			T_areYouThereMessage = T_areYouThereMessage - 1
			time.Sleep(1 * time.Second)
		}
		if T_areYouThereMessage == 0 && mapAreYouThereMessage[MyCoord] == false {
			T_areYouThereMessage = timer_T
			Recovery()
		} else {
			T_areYouThereMessage = timer_T
			mapMutex.Lock()
			mapAreYouThereMessage[MyCoord] = false
			mapMutex.Unlock()
			if mapAreYouThereMessage2[MyCoord] == "No" {
				mapMutex.Lock()
				mapAreYouThereMessage2[MyCoord] = ""
				mapMutex.Unlock()
				Recovery()
			}
		}
	}
}

func Check() {
	if s == "Normal" && c == my_address {
		TempSet := make(map[string]bool)
		for _,j := range machines {
			if j != my_address {
				jsonData := AreYouCoordinatorMessage{
					Sender: my_address}
				jsonValue, _ := json.Marshal(jsonData)
				go http.Post("http://"+j+"/AreYouCoordinatorMessage", "application/json", bytes.NewBuffer(jsonValue))
			
				for T_areYouCoordinatorMessage > 0 && mapAreYouCoordinatorMessage[j] == false {
					T_areYouCoordinatorMessage = T_areYouCoordinatorMessage - 1
					time.Sleep(1 * time.Second)
				}
				if T_areYouCoordinatorMessage == 0 && mapAreYouCoordinatorMessage[j] == false {
					T_areYouCoordinatorMessage = timer_T
					continue
				} else {
					T_areYouCoordinatorMessage = timer_T
					mapMutex.Lock()
					mapAreYouCoordinatorMessage[j] = false
					mapMutex.Unlock()
					if mapAreYouCoordinatorMessage2[j] == "Yes" {
						mapMutex.Lock()
						TempSet[j] = true
						mapAreYouCoordinatorMessage2[j] = ""
						mapMutex.Unlock()
					}
				}
			}
		}
		if len(TempSet) != 0 {
			p := maxAddressNum(TempSet)
			if hash(my_address) < p {
				seconds := int(math.Log10(float64(p - hash(my_address))))
				merge_invocated = false
				callMerge(TempSet, seconds)
			}
		}
	}
}

func callMerge(TempSet map[string]bool, seconds int) {
	x := seconds
	for x > 0 {
		x = x - 1
		time.Sleep(1 * time.Second)
	}
	if c == my_address {
		Merge(TempSet)
	}
}