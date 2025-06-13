package bets

type service struct{}

func NewService() BetService {
	return &service{}
}

func (s service) CreateBet(pollId string, userId int, selectedOptionIndex int) (*Bet, error) {
	bet := &Bet{
		PollId:              pollId,
		UserId:              userId,
		SelectedOptionIndex: selectedOptionIndex,
	}

	return bet, nil
}
