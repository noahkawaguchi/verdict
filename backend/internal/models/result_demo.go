package models

import (
	"fmt"
	"strconv"
)

func ResultDemo() {
	numUsers := 10
	dummyUserIDs := make([]string, numUsers)
	for i := range numUsers {
		dummyUserIDs[i] = strconv.Itoa(i)
	}
	poll, _, ballots := DummyData(dummyUserIDs)
	result := NewResult(poll, ballots)
	result.InstantRunoffVoting()
	fmt.Println(result)
}
