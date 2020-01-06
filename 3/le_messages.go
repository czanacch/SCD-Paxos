package main

type Acknowledgment struct { // Messaggio usato per confermare la ricezione
	Typemessage string `json:"typemessage"`
	Sender string `json:"sender"`
	Ans string `json:"ans"`
}

type ImYourLeader struct { // Messaggio mandato periodicamente dal leader per dire che c'è
	Sender string `json:"sender"`
}

type AcceptMessage struct {
	Sender string `json:"sender"`
	Gn Number `json:"gn"`
}
var mapAcceptMessage map[string]bool = make(map[string]bool)

type InvitationMessage struct {
	Sender string `json:"sender"`
	Gn Number `json:"gn"`
}
var mapInvitationMessage map[string]bool = make(map[string]bool)

type AreYouThereMessage struct {
	Sender string `json:"sender"`
	Gn Number `json:"gn"`
}
var mapAreYouThereMessage map[string]bool = make(map[string]bool)
var mapAreYouThereMessage2 map[string]string = make(map[string]string)

type AreYouCoordinatorMessage struct {
	Sender string `json:"sender"`
}
var mapAreYouCoordinatorMessage map[string]bool = make(map[string]bool)
var mapAreYouCoordinatorMessage2 map[string]string = make(map[string]string)

type ReadyMessage struct {
	Sender string `json:"sender"`
	Gn Number `json:"gn"`
}
var mapReadyMessage map[string]bool = make(map[string]bool)