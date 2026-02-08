package repositories

import (
	"cinema-system/internal/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TicketRepository struct {
	collection *mongo.Collection
}

func NewTicketRepository(db *mongo.Database) *TicketRepository {
	return &TicketRepository{
		collection: db.Collection("tickets"),
	}
}

func (r *TicketRepository) Create(ctx context.Context, ticket *models.Ticket) error {
	result, err := r.collection.InsertOne(ctx, ticket)
	if err != nil {
		return err
	}
	ticket.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *TicketRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Ticket, error) {
	var ticket models.Ticket
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&ticket)
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *TicketRepository) GetByUser(ctx context.Context, userID primitive.ObjectID) ([]models.Ticket, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tickets []models.Ticket
	if err = cursor.All(ctx, &tickets); err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *TicketRepository) GetBySession(ctx context.Context, sessionID primitive.ObjectID) ([]models.Ticket, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"session_id": sessionID,
		"status":     bson.M{"$ne": models.TicketCancelled},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tickets []models.Ticket
	if err = cursor.All(ctx, &tickets); err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *TicketRepository) UpdateStatus(ctx context.Context, ticketID primitive.ObjectID, status models.TicketStatus) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": ticketID},
		bson.M{"$set": bson.M{"status": status}},
	)
	return err
}

func (r *TicketRepository) GetAll(ctx context.Context) ([]models.Ticket, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tickets []models.Ticket
	if err = cursor.All(ctx, &tickets); err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *TicketRepository) CheckSeatAvailability(ctx context.Context, sessionID primitive.ObjectID, row, seat int) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{
		"session_id":  sessionID,
		"row_number":  row,
		"seat_number": seat,
		"status":      bson.M{"$ne": models.TicketCancelled},
	})
	return count == 0, err
}
