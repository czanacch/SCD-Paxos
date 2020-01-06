package main

import (
	"net/http"
	"fmt"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var my_address string // IP address of the process
var machines []string  // IP addresses of all processes

// Pair <id, num>
type Number struct {
	Id_process string `json:"id_process"`
	Num int `json:"num"`
}

func (a *Number) equal(b Number) bool {
	if a.Id_process == b.Id_process && a.Num == b.Num {
		return true
	} else {
		return false
	}
}

func (a *Number) less(b Number) bool {
    if a.Num < b.Num {
		return true
	} else {
		if a.Num > b.Num {
			return false
		} else {
			if a.Id_process < b.Id_process {
				return true
			} else {
				return false
			}
		}
	}
}

func (a *Number) lessEqual(b Number) bool {
	return a.equal(b) || a.less(b)
}

func initialization() {

	// Cose comuni
	my_address = "34.95.255.56:8080"
	machines = []string{"35.204.66.116:8080", "35.189.134.149:8080", "34.66.196.193:8080", "34.95.255.56:8080"}
	database = new(sql.DB)

	// Paxos
	p = new(Proposer) // Proposer initialization
	p.constructorProposer()

	a = new(Acceptor) // Acceptor initialization
	a.constructorAcceptor()

	l = new(Learner) // Learner initialization
	l.constructorLearner()
	
	// Leader election
	timer_T = 5
	timer_check = 5
	timer_leader = 5
	T_readyMessage = timer_T
	T_areYouThereMessage = timer_T
	T_areYouCoordinatorMessage = timer_T
	T_acceptMessage = timer_T

	s = "Normal"
	c = my_address
	p.im_leader = true
	counter = 0
	g = Number{Id_process: my_address, Num: 0}
	Up = make(map[string]bool)
	merge_invocated = false
}

func main() {

	initialization()
		
	fmt.Println("Starting machine", my_address)

	go leaderToGroup() // Periodic notification of the leader's work

	go invocatorCheck() // Periodic call to the Check() function

	go invocatorTimeout() // Periodic check of timer_leader

	lookDatabase()
	
	if len(l.decided_values) != 0 && len(machines) > 1 { // If I reboot after a failure
		a.mutex_PrepareRequest.Lock()
		a.mutex_AcceptRequest.Lock()
		l.mutex_AcceptResponse.Lock()
		askInformation()
	}

	go sendValues()

	http.HandleFunc("/values", handleValues)

	// LEADER ELECTION Calls
	http.HandleFunc("/AcknowledgmentMessage", handleAcknowledgmentMessage)

	http.HandleFunc("/ImYourLeaderMessage", handleImYourLeaderMessage)

	http.HandleFunc("/AcceptMessage", handleAcceptMessage)

	http.HandleFunc("/InvitationMessage", handleInvitationMessage)

	http.HandleFunc("/AreYouThereMessage", handleAreYouThereMessage)

	http.HandleFunc("/AreYouCoordinatorMessage", handleAreYouCoordinatorMessage)

	http.HandleFunc("/ReadyMessage", handleReadyMessage)
	
	// PAXOS Calls
	http.HandleFunc("/getDatabase", handleGetDatabase)

	http.HandleFunc("/setDatabase", handleSetDatabase)

	http.HandleFunc("/proposer/prepareResponse", handlePrepareResponse)
	
	http.HandleFunc("/acceptor/prepareRequest", handlePrepareRequest)

	http.HandleFunc("/acceptor/acceptRequest", handleAcceptRequest)
	
	http.HandleFunc("/learner/acceptResponse", handleAcceptResponse)

	http.ListenAndServe(":8080", nil)
}