package datastore

import (
	"context"
	"fmt"
	"log"

	"github.com/noahkawaguchi/verdict/backend/internal/models"
)

func DatastoreDemo() {
	ctx := context.TODO() // This will eventually come from the Lambda handler
	dummyPollID := "133"
	dummyUserID1 := "42"
	dummyUserID2 := "43"
	poll, ballot1, ballot2 := models.DummyData(dummyPollID, dummyUserID1, dummyUserID2)

	// Try to put the poll
	if err := PutPoll(ctx, poll); err != nil {
		log.Println("Failed to put poll:", err)
	} else {
		fmt.Println("Successfully put poll")
	}
	// Try to put the ballots
	if err := PutBallot(ctx, ballot1); err != nil {
		log.Println("Failed to put ballot1:", err)
	} else {
		fmt.Println("Successfully put ballot1")
	}
	if err := PutBallot(ctx, ballot2); err != nil {
		log.Println("Failed to put ballot2:", err)
	} else {
		fmt.Println("Successfully put ballot2")
	}
	// Try to get the poll
	gotPoll, err := getPoll(ctx, dummyPollID)
	if err != nil {
		log.Println("Failed to get poll:", err)
	} else {
		fmt.Println("Successfully got poll:", gotPoll)
	}
	// Try to get the ballots
	gotBallot1, err := getBallot(ctx, dummyPollID, dummyUserID1)
	if  err != nil {
		log.Println("Failed to get ballot1:", err)
	} else {
		fmt.Println("Successfully got ballot1:", gotBallot1)
	}
	gotBallot2, err := getBallot(ctx, dummyPollID, dummyUserID2)
	if  err != nil {
		log.Println("Failed to get ballot2:", err)
	} else {
		fmt.Println("Successfully got ballot2:", gotBallot2)
	}
	// Tally the results
	fmt.Println(gotPoll.TallyVotes(gotBallot1, gotBallot2))
}
