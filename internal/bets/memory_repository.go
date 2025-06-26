package bets

import "errors"

type memoryRepository struct {
	betList map[BetKey]Bet
}

func NewMemoryRepository() BetRepository {
	return &memoryRepository{
		betList: make(map[BetKey]Bet),
	}
}

func (repo memoryRepository) Save(bet *Bet) error {
	key := BetKey{bet.PollId, bet.UserId}
	if _, exists := repo.betList[key]; exists {
		return errors.New("user already placed a bet on this poll")
	}

	repo.betList[key] = *bet
	return nil
}

func (repo memoryRepository) GetByPollIdAndUserId(pollID string, userID string) (*Bet, error) {
	key := BetKey{PollId: pollID, UserId: userID}
	if bet, exists := repo.betList[key]; exists {
		return &bet, nil
	}
	return nil, errors.New("bet not found for the given poll and user")
}

func (repo memoryRepository) GetBetsFromUser(userID string) ([]Bet, error) {
	var bets []Bet
	for key, bet := range repo.betList {
		if key.UserId == userID {
			bets = append(bets, bet)
		}
	}
	return bets, nil
}

func (repo memoryRepository) GetBetsByPollId(pollID string) ([]Bet, error) {
	var bets []Bet
	for key, bet := range repo.betList {
		if key.PollId == pollID {
			bets = append(bets, bet)
		}
	}
	return bets, nil
}

func (repo memoryRepository) UpdateBet(bet *Bet) error {
	key := BetKey{bet.PollId, bet.UserId}
	if _, exists := repo.betList[key]; !exists {
		return errors.New("bet not found for the given poll and user")
	}

	repo.betList[key] = *bet
	return nil
}

var _ BetRepository = (*memoryRepository)(nil)
