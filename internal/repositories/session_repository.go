package repositories

import (
	"cinema-system/internal/models"
	"sync"
)

type SessionRepo struct {
	mu     sync.RWMutex
	data   map[int]models.Session
	nextID int
}

func NewSessionRepo() *SessionRepo {
	return &SessionRepo{
		data:   make(map[int]models.Session),
		nextID: 1,
	}
}

func (r *SessionRepo) Create(session models.Session) models.Session {
	r.mu.Lock()
	defer r.mu.Unlock()

	session.ID = r.nextID
	r.nextID++

	r.data[session.ID] = session
	return session
}

func (r *SessionRepo) GetAll() []models.Session {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := []models.Session{}
	for _, s := range r.data {
		list = append(list, s)
	}
	return list
}
