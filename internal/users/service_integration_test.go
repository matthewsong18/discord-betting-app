package users

import (
	"betting-discord-bot/internal/bets"
	"github.com/google/uuid"
	"testing"
)

type mockBetService struct {
	betsToReturn []bets.Bet
}

func (m *mockBetService) GetBetsFromUser(userID string) ([]bets.Bet, error) {
	// Add User ID to bet
	for i := range m.betsToReturn {
		m.betsToReturn[i].UserId = userID
	}
	return m.betsToReturn, nil
}

func (m *mockBetService) CreateBet(string, string, int) (*bets.Bet, error) {
	return nil, nil
}
func (m *mockBetService) GetBet(string, string) (*bets.Bet, error) { return nil, nil }
func (m *mockBetService) UpdateBetsByPollId(string) error          { return nil }

var _ bets.BetService = (*mockBetService)(nil)

func getTestBets(wins int, losses int, pending int) []bets.Bet {
	var betList []bets.Bet

	for i := 0; i < wins; i++ {
		betList = append(
			betList,
			bets.Bet{BetStatus: bets.Won},
		)
	}

	for i := 0; i < losses; i++ {
		betList = append(betList, bets.Bet{BetStatus: bets.Lost})
	}

	for i := 0; i < pending; i++ {
		betList = append(betList, bets.Bet{BetStatus: bets.Pending})
	}

	for i := range betList {
		betList[i].PollId = uuid.NewString()
	}

	return betList
}

func TestGetUserWinLoss(t *testing.T) {
	betsArgument := []struct {
		name    string
		betList []bets.Bet
		winLoss *WinLoss
	}{
		{
			"no bets",
			getTestBets(0, 0, 0),
			&WinLoss{Wins: 0, Losses: 0},
		},
		{
			"a winning bet",
			getTestBets(1, 0, 0),
			&WinLoss{Wins: 1, Losses: 0},
		},
		{
			"a losing bet",
			getTestBets(0, 1, 0),
			&WinLoss{Wins: 0, Losses: 1},
		},
		{
			"pending bets should not count",
			getTestBets(1, 1, 1),
			&WinLoss{1, 1},
		},
	}

	for _, tc := range betsArgument {
		t.Run(tc.name, func(t *testing.T) {
			testWinLoss(t, tc.betList, tc.winLoss)
		})
	}
}

func testWinLoss(t *testing.T, betsToReturn []bets.Bet, expectedWinLoss *WinLoss) {
	// ARRANGE
	mockBets := &mockBetService{
		betsToReturn: betsToReturn,
	}

	userRepo := NewMemoryRepository()
	userService := NewService(userRepo, mockBets)

	user, err := userService.CreateUser("test-discord-id")
	if err != nil {
		t.Fatalf("Setup failed: could not create user: %v", err)
	}

	// ACT
	actualWinLoss, err := userService.GetWinLoss(user.ID)
	if err != nil {
		t.Fatalf("GetWinLoss returned an unexpected error: %v", err)
	}

	// ASSERT
	if actualWinLoss.Wins != expectedWinLoss.Wins {
		t.Errorf("Expected Wins to be %d, but got %d", expectedWinLoss.Wins, actualWinLoss.Wins)
	}

	if actualWinLoss.Losses != expectedWinLoss.Losses {
		t.Errorf("Expected Losses to be %d, but got %d", expectedWinLoss.Losses, actualWinLoss.Losses)
	}
}
