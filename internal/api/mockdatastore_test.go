package api_test

import "github.com/noahkawaguchi/verdict/internal/voting"

// mockDatastore implements the datastore interface for testing purposes.
type mockDatastore struct {
	PutPollMock    func(poll *voting.Poll) error
	GetPollMock    func(pollID string) (*voting.Poll, error)
	PutBallotMock  func(ballot *voting.Ballot) error
	GetBallotsMock func(pollID string) ([]*voting.Ballot, error)
}

func (m *mockDatastore) PutPoll(poll *voting.Poll) error {
	if m.PutPollMock != nil {
		return m.PutPollMock(poll)
	}
	return nil
}

func (m *mockDatastore) GetPoll(pollID string) (*voting.Poll, error) {
	if m.GetPollMock != nil {
		return m.GetPollMock(pollID)
	}
	return nil, nil
}

func (m *mockDatastore) PutBallot(ballot *voting.Ballot) error {
	if m.PutBallotMock != nil {
		return m.PutBallotMock(ballot)
	}
	return nil
}

func (m *mockDatastore) GetBallots(pollID string) ([]*voting.Ballot, error) {
	if m.GetBallotsMock != nil {
		return m.GetBallotsMock(pollID)
	}
	return nil, nil
}
