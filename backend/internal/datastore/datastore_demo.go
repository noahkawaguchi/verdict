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
	dummyUserID3 := "44"
	poll, ballot1, ballot2, ballot3 := models.DummyData(
		dummyPollID, dummyUserID1, dummyUserID2, dummyUserID3)

	// Try to put the poll
	if err := PutPoll(ctx, poll); err != nil {
		log.Println("Failed to put poll:", err)
	} else {
		fmt.Println("Successfully put poll")
	}
	// Try to put the ballots
	for _, ballot := range []*models.Ballot{ballot1, ballot2, ballot3} {
		if err := PutBallot(ctx, ballot); err != nil {
			log.Println("Failed to put ballot:", err)
		} else {
			fmt.Println("Successfully put ballot")
		}
	}
	// Try to get the poll
	gotPoll, err := getPoll(ctx, dummyPollID)
	if err != nil {
		log.Println("Failed to get poll:", err)
	} else {
		fmt.Println("Successfully got poll:", gotPoll)
	}
	// Try to get the ballots
	for _, dummyUserID := range []string{dummyUserID1, dummyUserID2, dummyUserID3} {
		gotBallot, err := getBallot(ctx, dummyPollID, dummyUserID)
		if err != nil {
			log.Println("Failed to get ballot:", err)
		} else {
			fmt.Println("Successfully got ballot:", gotBallot)
		}
	}
	// Get ballots only using poll ID
	pollBallots, err := getPollBallots(ctx, dummyPollID)
	if  err != nil {
		log.Println("Failed to get pollBallots:", err)
	} else {
		fmt.Println("Successfully got pollBallots:", pollBallots)
	}
	// Tally the results
	fmt.Println(gotPoll.TallyVotes(pollBallots...))
}
