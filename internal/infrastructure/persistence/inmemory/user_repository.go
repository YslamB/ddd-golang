package inmemory

import (
	"context"
	"fmt"
	"sync"

	domain_common "gddd/internal/domain/common"
	domain_user "gddd/internal/domain/user"
)

type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[domain_user.UserID]*domain_user.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[domain_user.UserID]*domain_user.User),
	}
}

func (r *InMemoryUserRepository) Save(ctx context.Context, u *domain_user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[u.ID] = u
	return nil
}

func (r *InMemoryUserRepository) FindByID(ctx context.Context, id domain_user.UserID) (*domain_user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.users[id]
	if !ok {
		return nil, domain_common.NewNotFoundError(fmt.Sprintf("user with ID %s not found", id))
	}

	copiedUser := *user
	return &copiedUser, nil
}

func (r *InMemoryUserRepository) FindByEmail(ctx context.Context, email string) (*domain_user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, user := range r.users {
		if user.Email == email {

			copiedUser := *user
			return &copiedUser, nil
		}
	}
	return nil, domain_common.NewNotFoundError(fmt.Sprintf("user with email %s not found", email))
}

func (r *InMemoryUserRepository) Delete(ctx context.Context, id domain_user.UserID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[id]; !ok {
		return domain_common.NewNotFoundError(fmt.Sprintf("user with ID %s not found for deletion", id))
	}
	delete(r.users, id)
	return nil
}
