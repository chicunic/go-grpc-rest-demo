package service

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go-grpc-rest-demo/internal/server/errors"
	"go-grpc-rest-demo/internal/server/model"
)

type UserService struct {
	users  map[string]*model.User
	nextID int64
	mu     sync.RWMutex
}

func NewUserService() *UserService {
	return &UserService{
		users:  make(map[string]*model.User),
		nextID: 1,
	}
}

func (s *UserService) generateID() string {
	return strconv.FormatInt(atomic.AddInt64(&s.nextID, 1)-1, 10)
}

func (s *UserService) CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	if req.Username == "" || req.Email == "" || req.FullName == "" {
		return nil, errors.NewValidationError("fields", "username, email, and full_name are required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, user := range s.users {
		if user.Username == req.Username {
			return nil, errors.NewAlreadyExistsError("user", "username", req.Username)
		}
		if user.Email == req.Email {
			return nil, errors.NewAlreadyExistsError("user", "email", req.Email)
		}
	}

	now := time.Now()
	user := &model.User{
		ID:        s.generateID(),
		Username:  req.Username,
		Email:     req.Email,
		FullName:  req.FullName,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.users[user.ID] = user
	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, id string) (*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return nil, errors.NewNotFoundError("user", id)
	}
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *model.UpdateUserRequest) (*model.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[req.ID]
	if !exists {
		return nil, errors.NewNotFoundError("user", req.ID)
	}

	if req.Username != nil {
		if err := s.checkUniqueField(req.ID, "username", *req.Username); err != nil {
			return nil, err
		}
		user.Username = *req.Username
	}

	if req.Email != nil {
		if err := s.checkUniqueField(req.ID, "email", *req.Email); err != nil {
			return nil, err
		}
		user.Email = *req.Email
	}

	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	user.UpdatedAt = time.Now()

	return user, nil
}

func (s *UserService) checkUniqueField(excludeID, field, value string) error {
	for id, u := range s.users {
		if id == excludeID {
			continue
		}
		var existing string
		if field == "username" {
			existing = u.Username
		} else {
			existing = u.Email
		}
		if existing == value {
			return errors.NewAlreadyExistsError("user", field, value)
		}
	}
	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[id]; !exists {
		return errors.NewNotFoundError("user", id)
	}

	delete(s.users, id)
	return nil
}

func (s *UserService) ListUsers(ctx context.Context, req *model.ListUsersRequest) ([]model.User, int32, int32, int32, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := s.filterUsers(req.Filter)
	s.sortUsers(users, req.SortBy)

	paged, total, page, pageSize := paginate(users, req.Page, req.PageSize)
	return paged, total, page, pageSize, nil
}

func (s *UserService) filterUsers(filter *string) []model.User {
	var users []model.User
	var filterLower string
	if filter != nil {
		filterLower = strings.ToLower(*filter)
	}

	for _, user := range s.users {
		if filterLower != "" && !s.matchesFilter(user, filterLower) {
			continue
		}
		users = append(users, *user)
	}
	return users
}

func (s *UserService) matchesFilter(user *model.User, filter string) bool {
	return strings.Contains(strings.ToLower(user.Username), filter) ||
		strings.Contains(strings.ToLower(user.Email), filter) ||
		strings.Contains(strings.ToLower(user.FullName), filter)
}

func (s *UserService) sortUsers(users []model.User, sortBy *string) {
	field := "id"
	if sortBy != nil && *sortBy != "" {
		field = *sortBy
	}

	sort.Slice(users, func(i, j int) bool {
		switch field {
		case "username":
			return users[i].Username < users[j].Username
		case "email":
			return users[i].Email < users[j].Email
		case "full_name":
			return users[i].FullName < users[j].FullName
		case "created_at":
			return users[i].CreatedAt.Before(users[j].CreatedAt)
		default:
			return users[i].ID < users[j].ID
		}
	})
}
