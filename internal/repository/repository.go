package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Repository struct {
	database *pgxpool.Pool
}

func NewRepository(database *pgxpool.Pool) *Repository {
	return &Repository{database}
}
