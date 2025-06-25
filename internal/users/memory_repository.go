package users

import "errors"

type memoryRepository struct {
	users map[string]*User
}

func NewMemoryRepository() UserRepository {
	return &memoryRepository{
		users: make(map[string]*User),
	}
}

func (repo memoryRepository) Save(user *User) error {
	if user == nil {
		return errors.New("user is nil")
	}

	repo.users[user.ID] = user
	return nil
}

func (repo memoryRepository) GetByID(id string) (*User, error) {
	user, exists := repo.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (repo memoryRepository) GetByDiscordID(discordID string) (*User, error) {
	for i, user := range repo.users {
		if user.DiscordID == discordID {
			return repo.users[i], nil
		}
	}
	return nil, errors.New("user not found")
}

func (repo memoryRepository) Delete(discordID string) error {
	for id, user := range repo.users {
		if user.DiscordID == discordID {
			delete(repo.users, id)
			return nil
		}
	}
	return errors.New("user not found")
}

var _ UserRepository = (*memoryRepository)(nil)
