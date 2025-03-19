package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
)

type result struct {
	poll    *Poll
	ballots []*Ballot
	// The slice at each index holds the indices of the ballots currently counting for that choice
	votes        [][]int
	winnerIdx    int
	winningRound int
}

// CalculateResultData performs instant runoff voting using the provided poll and ballots,
// returning the computed data in JSON string format.
func CalculateResultData(poll *Poll, ballots []*Ballot) (string, error) {
	// Initialize votes to empty slices so nil can be used for elimination
	votes := make([][]int, len(poll.choices))
	for i := range votes {
		votes[i] = make([]int, 0)
	}
	res := &result{
		poll:         poll,
		ballots:      ballots,
		votes:        votes,
		winnerIdx:    -99,
		winningRound: 0,
	}
	res.instantRunoffVoting() // Compute the result
	if res.winnerIdx < 0 {    // Something went wrong
		return "", errors.New("the result was not successfully computed")
	}
	// Marshal into JSON format
	if body, err := json.Marshal(&struct {
		Prompt        string `json:"prompt"`
		TotalVotes    int    `json:"totalVotes"`
		WinningVotes  int    `json:"winningVotes"`
		WinningChoice string `json:"winningChoice"`
		WinningRound  int    `json:"winningRound"`
	}{
		Prompt:        res.poll.prompt,
		TotalVotes:    len(res.ballots),
		WinningVotes:  len(res.votes[res.winnerIdx]),
		WinningChoice: res.poll.choices[res.winnerIdx],
		WinningRound:  res.winningRound,
	}); err != nil {
		return "", errors.New("failed to marshal result")
	} else {
		return string(body), nil
	}
}

// instantRunoffVoting implements ranked choice voting, specifically the instant runoff method, to
// calculate the winning choice amongst the submitted ballots.
func (r *result) instantRunoffVoting() {
	// Tally first-choice votes
	for i, ballot := range r.ballots {
		firstChoiceIdx := ballot.rankOrder[0]
		r.votes[firstChoiceIdx] = append(r.votes[firstChoiceIdx], i)
	}
	// Majority check and elimination
	for i := range len(r.poll.choices) { // The number of choice ranks
		// Check if any choice has a strict majority of votes
		for j, choiceBallots := range r.votes {
			if float64(len(choiceBallots))/float64(len(r.ballots)) > 0.5 {
				r.winnerIdx = j
				r.winningRound = i + 1
				return
			}
		}
		// Find the choice(s) in last place
		minVotes := math.MaxInt
		var minIndices []int
		for j, choiceBallots := range r.votes {
			if choiceBallots != nil { // Don't consider eliminated choices
				if len(choiceBallots) < minVotes { // New last place found
					minVotes = len(choiceBallots)
					minIndices = []int{j}
				} else if len(choiceBallots) == minVotes { // Tie for last
					minIndices = append(minIndices, j)
				}
			}
		}
		// Break ties for last if necessary
		var loserIdx int
		if len(minIndices) > 1 {
			loserIdx = r.breakTiesForLast(minIndices)
		} else {
			loserIdx = minIndices[0]
		}
		// Redistribute the losing choice's votes to other choices
		for _, ballotIdx := range r.votes[loserIdx] {
			for _, choice := range r.ballots[ballotIdx].rankOrder {
				// If this choice is not the one being eliminated now and has not been eliminated
				// in a previous round, redistribute this ballot to the choice
				if choice != loserIdx && r.votes[choice] != nil {
					r.votes[choice] = append(r.votes[choice], ballotIdx)
					break
				}
			}
		}
		// Eliminate the losing choice
		r.votes[loserIdx] = nil
	}
}

// breakTiesForLast handles cases in instant runoff voting where multiple choices are tied for
// last place.
func (r *result) breakTiesForLast(tiedIndices []int) int {
	tieBreakVotes := make([]int, len(r.votes))
	// Tally votes using the highest rank that is one of the tied candidates
	for _, ballot := range r.ballots {
		for _, choiceIdx := range ballot.rankOrder {
			if slices.Contains(tiedIndices, choiceIdx) {
				tieBreakVotes[choiceIdx]++
				break
			}
		}
	}
	// Find the choice(s) in last place
	minVotes := math.MaxInt
	var minIndices []int
	for _, tiedIdx := range tiedIndices {
		if tieBreakVotes[tiedIdx] < minVotes { // New last place found
			minVotes = tieBreakVotes[tiedIdx]
			minIndices = []int{tiedIdx}
		} else if tieBreakVotes[tiedIdx] == minVotes { // Tie for last again
			minIndices = append(minIndices, tiedIdx)
		}
	}
	switch len(minIndices) {
	case 1: // Single minimum found
		return minIndices[0]
	case len(tiedIndices): // No choices were eliminated
		return minIndices[rand.IntN(len(minIndices))] // Choose randomly to avoid infinite recursion
	default:
		return r.breakTiesForLast(minIndices)
	}
}

func (r *result) String() string {
	if r.winnerIdx < 0 {
		return "The result was not successfully computed. Was the poll valid with at least one " +
			"corresponding ballot?"
	}
	return fmt.Sprintf("\nIn the poll \"%s,\" the choice %q won with "+
		"%d out of %d votes in round %d.\n",
		r.poll.prompt,
		r.poll.choices[r.winnerIdx],
		len(r.votes[r.winnerIdx]),
		len(r.ballots),
		r.winningRound,
	)
}
