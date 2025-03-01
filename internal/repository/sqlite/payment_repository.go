package sqlite

import (
	"database/sql"
	"time"

	"HelpBot/internal/domain"
)

// PaymentRepository реализует интерфейс domain.PaymentRepository для SQLite
type PaymentRepository struct {
	db *sql.DB
}

// NewPaymentRepository создает новый экземпляр PaymentRepository
func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}

// SaveRequest сохраняет запрос на пополнение баланса
func (r *PaymentRepository) SaveRequest(chatID int64, amount float64) error {
	_, err := r.db.Exec(`
		INSERT INTO payment_requests (chat_id, amount, status, created_at)
		VALUES (?, ?, ?, ?)
	`, chatID, amount, "pending", time.Now())
	return err
}

// GetRequests возвращает все запросы на пополнение баланса
func (r *PaymentRepository) GetRequests() ([]*domain.PaymentRequest, error) {
	rows, err := r.db.Query(`
		SELECT id, chat_id, amount, status, created_at
		FROM payment_requests
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*domain.PaymentRequest
	for rows.Next() {
		var request domain.PaymentRequest
		if err := rows.Scan(
			&request.ID,
			&request.ChatID,
			&request.Amount,
			&request.Status,
			&request.CreatedAt,
		); err != nil {
			return nil, err
		}
		requests = append(requests, &request)
	}
	return requests, nil
}

// UpdateRequestStatus обновляет статус запроса
func (r *PaymentRepository) UpdateRequestStatus(id int64, status string) error {
	_, err := r.db.Exec(`
		UPDATE payment_requests SET status = ? WHERE id = ?
	`, status, id)
	return err
}
