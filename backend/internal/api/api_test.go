package api_test

import "github.com/noahkawaguchi/verdict/backend/internal/models"

// mockDatastore implements the datastore interface for testing purposes.
type mockDatastore struct {
	PutPollMock    func(poll *models.Poll) error
	GetPollMock    func(pollID string) (*models.Poll, error)
	PutBallotMock  func(ballot *models.Ballot) error
	GetBallotsMock func(pollID string) ([]*models.Ballot, error)
}

func (m *mockDatastore) PutPoll(poll *models.Poll) error {
	if m.PutPollMock != nil {
		return m.PutPollMock(poll)
	}
	return nil
}

func (m *mockDatastore) GetPoll(pollID string) (*models.Poll, error) {
	if m.GetPollMock != nil {
		return m.GetPollMock(pollID)
	}
	return nil, nil
}

func (m *mockDatastore) PutBallot(ballot *models.Ballot) error {
	if m.PutBallotMock != nil {
		return m.PutBallotMock(ballot)
	}
	return nil
}

func (m *mockDatastore) GetBallots(pollID string) ([]*models.Ballot, error) {
	if m.GetBallotsMock != nil {
		return m.GetBallotsMock(pollID)
	}
	return nil, nil
}
