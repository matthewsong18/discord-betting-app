package polls

import (
	"errors"
	"github.com/google/uuid"
)

type service struct {
	pollList []Poll
}

func NewService() PollService {
	return &service{
		pollList: make([]Poll, 0),
	}
}

func (s *service) CreatePoll(title string, options []string) (*Poll, error) {
	if notExactlyTwo(options) {
		return nil, errors.New("poll must have exactly two options")
	}

	poll := &Poll{
		ID:      uuid.New().String(),
		Title:   title,
		Options: options,
		IsOpen:  true,
		Outcome: -1, // -1 indicates no outcome selected yet
	}

	s.pollList = append(s.pollList, *poll)

	return poll, nil
}

func notExactlyTwo(options []string) bool {
	return len(options) != 2
}

func (s *service) ClosePoll(poll *Poll) {
	for i, storedPoll := range s.pollList {
		if storedPoll.ID == poll.ID {
			s.pollList[i].IsOpen = false
		}
	}
}

func (s *service) SelectOutcome(poll *Poll, outcomeIndex int) error {
	if outcomeIndex < 0 || outcomeIndex >= len(poll.Options) {
		return errors.New("invalid outcome index")
	}
	poll.Outcome = outcomeIndex

	return nil
}

func (s *service) GetPollById(id string) (Poll, error) {
	for _, poll := range s.pollList {
		// Skip if the poll ID does not match
		if poll.ID != id {
			continue
		}

		return poll, nil
	}
	return Poll{}, errors.New("poll not found")
}

var _ PollService = (*service)(nil)
