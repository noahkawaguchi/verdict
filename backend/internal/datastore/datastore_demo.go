package datastore

import (
	"context"
	"fmt"
	"log"

	"github.com/noahkawaguchi/verdict/backend/internal/models"
)

func DatastoreDemo() {
	ctx := context.TODO() // This will eventually come from the Lambda handler
	dummyUserIDs := []string{"41", "42", "43", "44", "45"}
	poll, pollID, ballots := models.DummyData(dummyUserIDs)

	// Try to put the poll
	if err := PutPoll(ctx, poll); err != nil {
		log.Println("Failed to put poll:", err)
	} else {
		fmt.Printf("Successfully put poll %q\n", poll.GetPrompt())
	}
	// Try to put the ballots
	for i, ballot := range ballots {
		if err := PutBallot(ctx, ballot); err != nil {
			log.Println("Failed to put ballot:", err)
		} else {
			fmt.Printf("Successfully put the ballot of user %v\n", dummyUserIDs[i])
		}
	}
	// Try to get the poll
	gotPoll, err := GetPoll(ctx, pollID)
	if err != nil {
		log.Println("Failed to get poll:", err)
	} else {
		fmt.Println("Successfully got poll:", gotPoll)
	}
	// Try to get the ballots
	for _, dummyUserID := range dummyUserIDs {
		gotBallot, err := getBallot(ctx, pollID, dummyUserID)
		if err != nil {
			log.Println("Failed to get ballot:", err)
		} else {
			fmt.Println("Successfully got ballot:", gotBallot)
		}
	}
	// Get ballots only using poll ID
	pollBallots, err := GetPollBallots(ctx, pollID)
	if err != nil {
		log.Println("Failed to get pollBallots:", err)
	} else {
		fmt.Println("Successfully got pollBallots:")
		for _, pb := range pollBallots {
			fmt.Println(" ", pb)
		}
	}
	// Compute the results
	fmt.Println(models.NewResult(poll, ballots))
}
