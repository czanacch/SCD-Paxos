package main

import (
	"encoding/json"
	"net/http"
	"bytes"
	_ "github.com/mattn/go-sqlite3"
)

/****** Leader Election Events *********/

func handleAcknowledgmentMessage(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Content-Type", "application/json")
	var data Acknowledgment
	_ = json.NewDecoder(r.Body).Decode(&data)

	switch data.Typemessage {
		case "AcceptMessage": {
			mapMutex.Lock()
			mapAcceptMessage[data.Sender] = true
			mapMutex.Unlock()
		}
		case "InvitationMessage": {
			mapMutex.Lock()
			mapInvitationMessage[data.Sender] = true
			mapMutex.Unlock()
		}
		case "AreYouThereMessage": {
			mapMutex.Lock()
			mapAreYouThereMessage[data.Sender] = true
			mapAreYouThereMessage2[data.Sender] = data.Ans
			mapMutex.Unlock()
		}
		case "AreYouCoordinatorMessage": {
			mapMutex.Lock()
			mapAreYouCoordinatorMessage[data.Sender] = true
			mapAreYouCoordinatorMessage2[data.Sender] = data.Ans
			mapMutex.Unlock()
		}
		case "ReadyMessage": {
			mapMutex.Lock()
			mapReadyMessage[data.Sender] = true
			mapMutex.Unlock()
		}
	}
}

func handleImYourLeaderMessage(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var data Acknowledgment
	_ = json.NewDecoder(r.Body).Decode(&data)

	if data.Sender == c {
		timer_leader = 5
	}
}

// p_j (which triggers this event) has accepted p_i's invitation to join the Gn group
func handleAcceptMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data AcceptMessage
	_ = json.NewDecoder(r.Body).Decode(&data)

	jsonData := Acknowledgment{
		Typemessage: "AcceptMessage",
		Sender: my_address }
	jsonValue, _ := json.Marshal(jsonData)
	go http.Post("http://"+data.Sender+"/AcknowledgmentMessage", "application/json", bytes.NewBuffer(jsonValue))

	stateMutex.Lock()
	if s == "Election" && g == data.Gn && c == my_address {
		j := data.Sender
		mapMutex.Lock()
		Up[j] = true
		mapMutex.Unlock()
	}
	stateMutex.Unlock()
}

// p_i is invited by a process p_j to join in the Gn group
func handleInvitationMessage(w http.ResponseWriter, r *http.Request) {
	var Temp string
	var TempSet map[string]bool
	
	w.Header().Set("Content-Type", "application/json")
	var data InvitationMessage
	_ = json.NewDecoder(r.Body).Decode(&data)

	jsonData := Acknowledgment{
		Typemessage: "InvitationMessage",
		Sender: my_address }
	jsonValue, _ := json.Marshal(jsonData)
	go http.Post("http://"+data.Sender+"/AcknowledgmentMessage", "application/json", bytes.NewBuffer(jsonValue))

	stateMutex.Lock()
	if s == "Normal" {
		Temp = c
		TempSet = Up
		s = "Election"
		c = data.Sender
		p.im_leader = false
		g = data.Gn
	}
	stateMutex.Unlock()
	
	if Temp == my_address {
		for address,_ := range TempSet {
			jsonData := InvitationMessage{
				Sender: data.Sender,
				Gn: data.Gn }
			jsonValue, _ := json.Marshal(jsonData)
			go http.Post("http://"+address+"/InvitationMessage", "application/json", bytes.NewBuffer(jsonValue))
		}
	}

	jsonData2 := AcceptMessage{
		Sender: my_address,
		Gn: data.Gn }
	jsonValue, _ = json.Marshal(jsonData2)
	go http.Post("http://"+data.Sender+"/AcceptMessage", "application/json", bytes.NewBuffer(jsonValue))

	mapMutex.Lock()
	mapAcceptMessage = make(map[string]bool)
	mapMutex.Unlock()

	stateMutex.Lock()
	s = "Reorganization"
	stateMutex.Unlock()
}

// The process p_j that sends this message to p_i wants to know if p_i is the leader of the Gn group and
// if it considers p_j as a member of that group
func handleAreYouThereMessage(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var data AreYouThereMessage
	_ = json.NewDecoder(r.Body).Decode(&data)
	
	// In this case the Acknowledgment brings with it something additional (Ans)
	stateMutex.Lock()
	if g == data.Gn && c == my_address && Up[data.Sender] == true {
		jsonData := Acknowledgment{
			Typemessage: "AreYouThereMessage",
			Sender: my_address,
			Ans: "Yes" }
		jsonValue, _ := json.Marshal(jsonData)
		go http.Post("http://"+data.Sender+"/AcknowledgmentMessage", "application/json", bytes.NewBuffer(jsonValue))
	} else {
		jsonData := Acknowledgment{
			Typemessage: "AreYouThereMessage",
			Sender: my_address,
			Ans: "No" }
		jsonValue, _ := json.Marshal(jsonData)
		go http.Post("http://"+data.Sender+"/AcknowledgmentMessage", "application/json", bytes.NewBuffer(jsonValue))
	}
	stateMutex.Unlock()
}

// The sending process wants to know if p_i is a leader in the "Normal" state. p_i can answer YES or NO
func handleAreYouCoordinatorMessage(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var data AreYouCoordinatorMessage
	_ = json.NewDecoder(r.Body).Decode(&data)

	// In this case the Acknowledgment brings with it something additional (Ans)
	stateMutex.Lock()
	if s == "Normal" && c == my_address {
		jsonData := Acknowledgment{
			Typemessage: "AreYouCoordinatorMessage",
			Sender: my_address,
			Ans: "Yes" }
		jsonValue, _ := json.Marshal(jsonData)
		go http.Post("http://"+data.Sender+"/AcknowledgmentMessage", "application/json", bytes.NewBuffer(jsonValue))
	} else {
		jsonData := Acknowledgment{
			Typemessage: "AreYouCoordinatorMessage",
			Sender: my_address,
			Ans: "No" }
		jsonValue, _ := json.Marshal(jsonData)
		go http.Post("http://"+data.Sender+"/AcknowledgmentMessage", "application/json", bytes.NewBuffer(jsonValue))
	}
	stateMutex.Unlock()

}

func handleReadyMessage(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var data ReadyMessage
	_ = json.NewDecoder(r.Body).Decode(&data)

	jsonData := Acknowledgment{
		Typemessage: "ReadyMessage",
		Sender: my_address }
	jsonValue, _ := json.Marshal(jsonData)
	go http.Post("http://"+data.Sender+"/AcknowledgmentMessage", "application/json", bytes.NewBuffer(jsonValue))

	stateMutex.Lock()
	if s == "Reorganization"  && g == data.Gn {
		s = "Normal"
	}
	stateMutex.Unlock()
}