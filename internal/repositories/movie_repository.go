package repositories

import (
	"context"
	"film-downloader/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MovieRepository struct {
	db *pgxpool.Pool
}

func NewMovieRepository(db *pgxpool.Pool) *MovieRepository {
	return &MovieRepository{db: db}
}

func (r *MovieRepository) CheckMovieExists(movieID string) (bool, error) {
	rows, err := r.db.Query(context.Background(), "SELECT id FROM movies WHERE id = $1", movieID)

	if err != nil {
		return false, err
	}

	defer rows.Close()
	return rows.Next(), nil
}

func (r *MovieRepository) CreateMovie(movie models.MovieResponse) error {
	return nil
}
