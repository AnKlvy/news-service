package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"news_service.andreyklimov.net/internal/validator"
)

type News struct {
	ID         int64     `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Categories []string  `json:"categories"`
	Status     string    `json:"status"`
	ImageURL   *string   `json:"image_url,omitempty"`
	Version    int32     `json:"version"`
}

// ValidateNews выполняет валидацию данных новости.
func ValidateNews(v *validator.Validator, news *News) {
	v.Check(news.Title != "", "title", "must be provided")
	v.Check(len(news.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(news.Content != "", "content", "must be provided")

	v.Check(news.Categories != nil, "categories", "must be provided")
	v.Check(len(news.Categories) >= 1, "categories", "must contain at least 1 categories")
	v.Check(len(news.Categories) <= 10, "categories", "must not contain more than 10 categories")
	v.Check(validator.Unique(news.Categories), "categories", "must not contain duplicate values")

	v.Check(news.Status != "", "status", "must be provided")
	v.Check(validator.PermittedValue(news.Status, "DRAFT", "PUBLISHED", "ARCHIVED"), "status", "must be a valid status")

	if news.ImageURL != nil {
		v.Check(len(*news.ImageURL) <= 1000, "image_url", "must not be more than 1000 bytes long")
	}
}

// Определяем структуру NewsModel, которая содержит пул соединений с базой данных.
type NewsModel struct {
	DB *sql.DB
}

func (m NewsModel) Insert(news *News) error {
	query := `
    INSERT INTO news (title, content, categories, status, image_url)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id, created_at, version`
	args := []any{news.Title, news.Content, pq.Array(news.Categories), news.Status, news.ImageURL}

	// Создаём контекст с тайм-аутом 3 секунды.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Используем QueryRowContext() и передаём контекст в качестве первого аргумента.
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&news.ID, &news.CreatedAt, &news.Version)
}

func (m NewsModel) Get(id int64) (*News, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
    SELECT id, created_at, updated_at, title, content, categories, status, image_url, version
    FROM news
    WHERE id = $1`

	var news News
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&news.ID,
		&news.CreatedAt,
		&news.UpdatedAt,
		&news.Title,
		&news.Content,
		pq.Array(&news.Categories),
		&news.Status,
		&news.ImageURL,
		&news.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &news, nil
}

func (m NewsModel) Update(news *News) error {
	query := `
    UPDATE news
    SET title = $1, content = $2, categories = $3, status = $4, image_url = $5, updated_at = now(), version = version + 1
    WHERE id = $6 AND version = $7
    RETURNING version`
	args := []any{
		news.Title,
		news.Content,
		pq.Array(news.Categories),
		news.Status,
		news.ImageURL,
		news.ID,
		news.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&news.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m NewsModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
    DELETE FROM news
    WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (m NewsModel) GetAll(title string, categories []string, status string, filters Filters) ([]*News, Metadata, error) {
	query := fmt.Sprintf(
		`SELECT count(*) OVER(), id, created_at, updated_at, title, content, categories, status, image_url, version
		 FROM news
		 WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		 AND (categories @> $2 OR $2 = '{}')
		 AND (status = $3 OR $3 = '')
		 ORDER BY %s %s, id ASC
		 LIMIT $4 OFFSET $5`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{title, pq.Array(categories), status, filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	news := []*News{}

	for rows.Next() {
		var new News
		err := rows.Scan(
			&totalRecords,
			&new.ID,
			&new.CreatedAt,
			&new.UpdatedAt,
			&new.Title,
			&new.Content,
			pq.Array(&new.Categories),
			&new.Status,
			&new.ImageURL,
			&new.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		news = append(news, &new)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return news, metadata, nil
}

type MockNewsModel struct{}

func (m MockNewsModel) Insert(news *News) error {
	return nil
}

func (m MockNewsModel) Get(id int64) (*News, error) {
	return nil, nil
}

func (m MockNewsModel) Update(news *News) error {
	return nil
}

func (m MockNewsModel) Delete(id int64) error {
	return nil
}

func (m MockNewsModel) GetAll(title string, categories []string, status string, filters Filters) ([]*News, Metadata, error) {
	return nil, Metadata{}, nil
}
