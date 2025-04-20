package database

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	// Устанавливаем поле News как интерфейс, содержащий методы, которые должны поддерживать
	// как 'реальная' модель, так и мок-модель.
	News interface {
		Insert(news *News) error
		Get(id int64) (*News, error)
		Update(news *News) error
		Delete(id int64) error
		GetAll(title string, categories []string, status string, filters Filters) ([]*News, Metadata, error)
	}
}

// Создаем вспомогательную функцию, которая возвращает экземпляр Models, содержащий только мок-модели.
func NewMockModels() Models {
	return Models{
		News: MockNewsModel{},
	}
}

// Для удобства мы также добавляем метод New(), который возвращает структуру Models
// с инициализированным NewsModel.
func NewModels(db *sql.DB) Models {
	return Models{
		News: NewsModel{DB: db},
	}
}
