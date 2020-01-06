package main

/************ Paxos Utility functions *********/

func contains(s []int, e int) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func returnMax(s []int) int {
	if len(s) == 0 {
		return 0
	}
	max := s[0]
	for _, v := range s {
			if (v > max) {
				max = v
			}
	}
	return max
}