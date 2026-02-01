package repositories

import (
	"cinema-system/internal/models"
	"sync"
	"time"
)

type TicketRepo struct {
	mu     sync.RWMutex
	data   map[int]models.Ticket
	nextID int
}

func NewTicketRepo() *TicketRepo {
	return &TicketRepo{
		data:   make(map[int]models.Ticket),
		nextID: 1,
	}
}

func (r *TicketRepo) Book(ticket models.Ticket) models.Ticket {
	r.mu.Lock()
	defer r.mu.Unlock()

	ticket.ID = r.nextID
	ticket.Status = "booked"
	ticket.CreatedAt = time.Now()
	r.nextID++

	r.data[ticket.ID] = ticket
	return ticket
}

func (r *TicketRepo) GetAll() []models.Ticket {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := []models.Ticket{}
	for _, t := range r.data {
		list = append(list, t)
	}
	return list
}
