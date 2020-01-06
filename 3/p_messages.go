package main

/************ Paxos Messages *********/

// Message sent by the user
type ActivationMessage struct {
    Values []int `json:"values"`
}

// PrepareRequest (Proposer -> Acceptor)
type PrepareRequest struct {
	Id string `json:"id"`
	N Number `json:"n"`
	Number_instance int `json:"number_instance"`
}

// PrepareResponse (Acceptor -> Proposer)
type PrepareResponse struct {
	Id string `json:"id"`
	N Number `json:"n"`
	V int `json:"v"`
	Number_instance int `json:"number_instance"`
}

// AcceptRequest (Acceptor -> Proposer)
type AcceptRequest struct {
	Id string `json:"id"`
	N Number `json:"n"`
	V int `json:"v"`
	Number_instance int `json:"number_instance"`
}

// AcceptResponse (Acceptor -> Learner)
type AcceptResponse struct {
	Id string `json:"id"`
	V int `json:"v"`
	Number_instance int `json:"number_instance"`
}

// ---------------------------------------------------------------------------------------------

type Notification struct {
    Address string `json:"address"`
}

// Message from a Learner to a process that sends a "database request"
type RecoveryMessage struct {
    Decided_values []int `json:"decided_values"`
}