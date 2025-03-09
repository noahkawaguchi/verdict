package models

import "fmt"

type result struct {
	poll    *Poll
	ballots []*Ballot
	// The slice at each index holds the indices of the ballots currently counting for that choice
	votes        [][]int
	winnerIdx    int
	winningRound int
}

func NewResult(poll *Poll, ballots []*Ballot) *result {
	// Validate that the ballots are for this poll
	j := 0
	for _, ballot := range ballots {
		if ballot.PollID == poll.PollID {
			ballots[j] = ballot
			j++
		}
	}
	// Initialize votes to empty slices so nil can be used for elimination
	votes := make([][]int, len(poll.Choices))
	for i := range votes {
		votes[i] = make([]int, 0)
	}
	res := &result{
		poll:         poll,
		ballots:      ballots[:j],
		votes:        votes,
		winnerIdx:    -99,
		winningRound: 0,
	}
	res.instantRunoffVoting() // Compute the result from the constructor
	return res
}

// instantRunoffVoting implements ranked-choice voting, specifically the instant runoff method, to
// calculate the winning choice amongst the submitted ballots.
func (r *result) instantRunoffVoting() {
	// Tally first-choice votes
	for i, ballot := range r.ballots {
		choice := ballot.RankOrder[0]
		r.votes[choice] = append(r.votes[choice], i)
	}
	// Majority check and elimination
	for i := range len(r.poll.Choices) { // The number of choice ranks
		// Check if any choice has a strict majority of votes
		for j, choiceBallots := range r.votes {
			if float64(len(choiceBallots))/float64(len(r.ballots)) > 0.5 {
				r.winnerIdx = j
				r.winningRound = i + 1
				return
			}
		}
		// Determine which choice is in last place
		// TODO: Handle ties for last (using votes from previous rounds)
		minVotesIdx := 0
		for i := 1; i < len(r.votes); i++ {
			if r.votes[i] != nil && // Don't consider eliminated choices
				(r.votes[minVotesIdx] == nil || // Handle the case where the first element is nil
					len(r.votes[i]) < len(r.votes[minVotesIdx])) {
				minVotesIdx = i
			}
		}
		// Redistribute the last place choice's votes to other choices
		for _, ballotIdx := range r.votes[minVotesIdx] {
			for _, choice := range r.ballots[ballotIdx].RankOrder {
				// If this choice is not the one being eliminated now and has not been eliminated
				// in a previous round, redistribute this ballot to the choice
				if choice != minVotesIdx && r.votes[choice] != nil {
					r.votes[choice] = append(r.votes[choice], ballotIdx)
					break
				}
			}
		}
		// Eliminate the last place choice
		r.votes[minVotesIdx] = nil
	}
}

func (r *result) String() string {
	if r.winnerIdx < 0 {
		return "The result was not successfully computed. Was the poll valid with at least one " +
			"corresponding ballot?"
	}
	return fmt.Sprintf("\nIn the poll %q\nThe choice %q won with %d votes in round %d\n",
		r.poll.Prompt,
		r.poll.Choices[r.winnerIdx],
		len(r.votes[r.winnerIdx]),
		r.winningRound,
	)
}
