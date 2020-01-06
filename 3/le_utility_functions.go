package main

import (
	"encoding/json"
	"net/http"
	"bytes"
	"time"
	"hash/fnv"
	_ "github.com/mattn/go-sqlite3"
)

// Utility function to convert strings to unique 10-digit integers
func hash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32())
}

func maxAddressNum(set map[string]bool) int {
	strings := []string{}
	for key,_ := range set {
		strings = append(strings, key)
	}
	numbers := []int{}
	for _, value := range strings {
		numbers = append(numbers, hash(value))
	}
	max := numbers[0]
	for _, value := range numbers {
		if value > max {
			max = value
		}
	}
	return max
}

func invocatorCheck() {
	for {
		for timer_check > 0 {
			timer_check = timer_check - 1
			time.Sleep(1 * time.Second)
		}
		if timer_check == 0 {
			timer_check = 5
			Check()
		}
		time.Sleep(1 * time.Millisecond)
	}
}

func invocatorTimeout() {
	for {
		if my_address != c {
			for timer_leader > 0 {
				timer_leader = timer_leader - 1
				time.Sleep(1 * time.Second)
			}
			if timer_leader == 0 {
				timer_leader = 5
				Timeout()
			}
		}
		time.Sleep(1 * time.Millisecond)
	}
}

func leaderToGroup() {
	for {
		time.Sleep(1 * time.Second)
		for j,_ := range Up {
			jsonData := ImYourLeader{
				Sender: my_address}
			jsonValue, _ := json.Marshal(jsonData)
			go http.Post("http://"+j+"/ImYourLeaderMessage", "application/json", bytes.NewBuffer(jsonValue)) // Dì esplicitamente a P_j di unirsi a g
		}
		time.Sleep(1 * time.Millisecond)
	}
}