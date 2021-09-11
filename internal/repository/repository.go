package repository

import "github.com/jackc/pgx/v4/pgxpool"

type Repository struct {
	Database *pgxpool.Pool
}

func NewRepository(database *pgxpool.Pool) *Repository {
	return &Repository{database}
}
