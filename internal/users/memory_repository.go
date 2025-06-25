package users

type memoryRepository struct {
}

func NewMemoryRepository() UserRepository {
	return memoryRepository{}
}

func (repo memoryRepository) Save(user *User) error {
	//TODO implement me
	panic("implement me")
}

func (repo memoryRepository) GetByID(id string) (*User, error) {
	//TODO implement me
	panic("implement me")
}

func (repo memoryRepository) GetByDiscordID(discordID string) (*User, error) {
	//TODO implement me
	panic("implement me")
}

func (repo memoryRepository) Delete(discordID string) error {
	//TODO implement me
	panic("implement me")
}

var _ UserRepository = (*memoryRepository)(nil)
