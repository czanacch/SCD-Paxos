package main

import (
	"sync"
	"encoding/json"
	"net/http"
	"bytes"
	_ "github.com/mattn/go-sqlite3"
)

/*********** PAXOS PROPOSER **************/

type Proposer struct {
	n Number
	v int
	max_number_received Number
	responded map[string]bool
	im_leader bool

	mutex_PrepareResponse sync.Mutex
	channel_PrepareResponse chan PrepareResponse

}

func (p *Proposer) constructorProposer() {
	p.n = Number{Id_process: my_address, Num: 0}
	p.v = 0
	p.max_number_received = Number{Id_process: "", Num: 0}
	p.responded = make(map[string]bool)
	p.im_leader = false

	p.mutex_PrepareResponse = sync.Mutex{}
	p.channel_PrepareResponse = make(chan PrepareResponse)
}

func (p *Proposer) prepareRequest() {
	number_instance = number_instance + 1
	p.n.Num = p.v + len(machines)

	for _,addr := range machines {
		jsonData := PrepareRequest{
			Id: my_address,
			N: p.n,
			Number_instance : number_instance }
		jsonValue, _ := json.Marshal(jsonData)

		go http.Post("http://"+addr+"/acceptor/prepareRequest", "application/json", bytes.NewBuffer(jsonValue))
	}
}

func handlePrepareResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var prepareResponse PrepareResponse
	_ = json.NewDecoder(r.Body).Decode(&prepareResponse)

	go managePrepareResponse()
	p.channel_PrepareResponse <- prepareResponse
}

func managePrepareResponse() {
	p.mutex_PrepareResponse.Lock()
	prepareResponse := <-p.channel_PrepareResponse
	if prepareResponse.Number_instance == number_instance {
		if len(p.responded) < len(machines)/2+1 {
			
			mapMutex.Lock()
			p.responded[prepareResponse.Id] = true // indicates that the prepareResponse.Id process responded with a promise
			mapMutex.Unlock()

			if p.max_number_received.less(prepareResponse.N) {
				p.max_number_received = prepareResponse.N
				p.v = prepareResponse.V
			}

		}
		if len(p.responded) >= len(machines)/2+1 {
			p.acceptRequest()
		}
	}
	p.mutex_PrepareResponse.Unlock()
}


func (p *Proposer) acceptRequest() {
	for _,addr := range machines {
		jsonData := AcceptRequest{
			Id: my_address,
			N: p.n,
			V: p.v,
			Number_instance : number_instance}
		jsonValue, _ := json.Marshal(jsonData)
		
		go http.Post("http://"+addr+"/acceptor/acceptRequest", "application/json", bytes.NewBuffer(jsonValue))
	}
}