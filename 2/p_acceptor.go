package main

import (
	"sync"
	"encoding/json"
	"net/http"
	"bytes"
	_ "github.com/mattn/go-sqlite3"
)

/************ PAXOS ACCEPTOR *************/
type Acceptor struct {
	v int
	n Number
	max_n Number

	mutex_PrepareRequest sync.Mutex
	channel_PrepareRequest chan PrepareRequest
	mutex_AcceptRequest sync.Mutex
	channel_AcceptRequest chan AcceptRequest
}

func (a *Acceptor) constructorAcceptor() {
	a.v = 0
	a.n = Number{Id_process: "", Num: 0}
	a.max_n = Number{Id_process: "", Num: 0}

	a.mutex_PrepareRequest = sync.Mutex{}
	a.channel_PrepareRequest = make(chan PrepareRequest)
	a.mutex_AcceptRequest = sync.Mutex{}
	a.channel_AcceptRequest = make(chan AcceptRequest)
}

func handlePrepareRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var prepareRequest PrepareRequest
	_ = json.NewDecoder(r.Body).Decode(&prepareRequest)
	
	go managePrepareRequest()
	a.channel_PrepareRequest <- prepareRequest
}

func managePrepareRequest() {
	a.mutex_PrepareRequest.Lock()
	prepareRequest := <-a.channel_PrepareRequest
	if a.max_n.less(prepareRequest.N) { // If the number is greater than anything promised

		a.max_n = prepareRequest.N // New maximum "proposal number"
		number_instance = prepareRequest.Number_instance
		
		jsonData := PrepareResponse{
			Id: my_address,
			N: a.n,
			V: a.v,
			Number_instance: number_instance}
		jsonValue, _ := json.Marshal(jsonData)
		
		go http.Post("http://"+prepareRequest.Id+"/proposer/prepareResponse", "application/json", bytes.NewBuffer(jsonValue))
	}
	a.mutex_PrepareRequest.Unlock()
}

func handleAcceptRequest(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var acceptRequest AcceptRequest
	_ = json.NewDecoder(r.Body).Decode(&acceptRequest)

	go manageAcceptRequest()
	a.channel_AcceptRequest <- acceptRequest
}

func manageAcceptRequest() {
	
	a.mutex_AcceptRequest.Lock()
	acceptRequest := <-a.channel_AcceptRequest
	
	if a.max_n.lessEqual(acceptRequest.N) && acceptRequest.Number_instance == number_instance {
		
		a.max_n = acceptRequest.N
		a.n = acceptRequest.N
		a.v = acceptRequest.V

		for _,addr := range machines {
			jsonData := AcceptResponse{
				Id: my_address,
				V: a.v,
				Number_instance: number_instance }
			jsonValue, _ := json.Marshal(jsonData)

			go http.Post("http://"+addr+"/learner/acceptResponse", "application/json", bytes.NewBuffer(jsonValue))
		}
	}
	a.mutex_AcceptRequest.Unlock()
}