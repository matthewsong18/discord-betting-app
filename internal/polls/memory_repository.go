package polls

import "errors"

type memoryRepository struct {
	polls map[string]*Poll
}

func NewMemoryRepository() PollRepository {
	return &memoryRepository{
		polls: make(map[string]*Poll),
	}
}

var ErrPollNotFound = errors.New("poll not found")

func (m memoryRepository) Save(poll *Poll) error {
	if _, exists := m.polls[poll.ID]; exists {
		return errors.New("poll already exists")
	}
	m.polls[poll.ID] = poll
	return nil
}

func (m memoryRepository) GetById(id string) (*Poll, error) {
	if poll, exists := m.polls[id]; exists {
		return poll, nil
	}
	return nil, ErrPollNotFound
}

func (m memoryRepository) Update(poll *Poll) error {
	if _, exists := m.polls[poll.ID]; !exists {
		return ErrPollNotFound
	}
	m.polls[poll.ID] = poll
	return nil
}

func (m memoryRepository) Delete(pollID string) error {
	if _, exists := m.polls[pollID]; !exists {
		return ErrPollNotFound
	}
	delete(m.polls, pollID)
	return nil
}

func (m memoryRepository) GetOpenPolls() ([]*Poll, error) {
	var openPolls []*Poll
	for _, poll := range m.polls {
		if poll.Status == Open {
			openPolls = append(openPolls, poll)
		}
	}

	return openPolls, nil
}

var _ PollRepository = (*memoryRepository)(nil)
