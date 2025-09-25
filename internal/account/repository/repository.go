package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hampgoodwin/GoLuca/internal/meta"
	ierrors "github.com/hampgoodwin/GoLuca/pkg/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Repository struct {
	database *pgxpool.Pool
}

func NewRepository(database *pgxpool.Pool) *Repository {
	return &Repository{database}
}

// GetAccount gets an account from the database
func (r *Repository) GetAccount(ctx context.Context, accountID string) (Account, error) {
	ctx, span := otel.Tracer(meta.ServiceName).Start(ctx, "repositry.GetAccount", trace.WithAttributes(
		attribute.String("account_id", accountID),
	))
	defer span.End()

	acct := Account{}
	if err := r.database.QueryRow(ctx,
		`SELECT id, parent_id, name, type, basis, created_at
		FROM account
		WHERE id=$1
		;`,
		accountID).Scan(
		&acct.ID, &acct.ParentID, &acct.Name, &acct.Type, &acct.Basis, &acct.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return acct, errors.Join(fmt.Errorf("account %q not found: %w", accountID, err), ierrors.NotFoundErr{Type: "account", ID: accountID})
		}
		return acct, errors.Join(fmt.Errorf("scanning account result row: %w", err), ierrors.ErrNotKnown)
	}
	return acct, nil
}

// ListAccounts get accounts paginated based on a cursor and limit
func (r *Repository) ListAccounts(ctx context.Context, accountID string, createdAt time.Time, limit uint64) ([]Account, error) {
	query := `SELECT id, parent_id, name, type, basis, created_at
		FROM account
		WHERE 1=1`
	var params []any
	if accountID != "" && !createdAt.IsZero() {
		params = append(params, accountID)
		query += fmt.Sprintf(" AND account.id <= $%d", len(params))
		params = append(params, createdAt)
		query += fmt.Sprintf(" AND account.created_at <= $%d", len(params))
	}
	query += " ORDER BY created_at DESC"
	if limit != 0 {
		params = append(params, limit)
		query += fmt.Sprintf(" LIMIT $%d", len(params))
	}
	query += ";"
	rows, err := r.database.Query(ctx, query, params...)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("fetching accounts from data store: %w", err), ierrors.ErrNotKnown)
	}
	defer rows.Close()
	accounts := []Account{}
	for rows.Next() {
		acct := Account{}
		if err := rows.Scan(&acct.ID, &acct.ParentID, &acct.Name, &acct.Type, &acct.Basis, &acct.CreatedAt); err != nil {
			return nil, errors.Join(fmt.Errorf("scanning account result row: %w", err), ierrors.ErrNotKnown)
		}
		accounts = append(accounts, acct)
	}
	return accounts, nil
}

// CreateAccount creates an account record in the database and returns the created record
func (r *Repository) CreateAccount(ctx context.Context, create Account) (Account, error) {
	// get a db-transaction
	tx, err := r.database.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return Account{}, errors.Join(fmt.Errorf("beginning create account db transaction: %w", err), ierrors.ErrNotKnown)
	}

	returning := Account{}
	if err := tx.QueryRow(ctx, `
		INSERT INTO account(id, parent_id, name, type, basis, created_at)
		VALUES($1, $2, $3, $4, $5, $6)
		RETURNING id, parent_id, name, type, basis, created_at
		;`,
		create.ID, create.ParentID, create.Name, create.Type, create.Basis, create.CreatedAt).Scan(
		&returning.ID,
		&returning.ParentID,
		&returning.Name,
		&returning.Type,
		&returning.Basis,
		&returning.CreatedAt,
	); err != nil {
		return returning, errors.Join(fmt.Errorf("scanning account returned from insert: %w", err), ierrors.ErrNotKnown)
	}
	if err := tx.Commit(ctx); err != nil {
		return returning, errors.Join(fmt.Errorf("committing account insert transaction: %w", err), ierrors.ErrNotKnown)
	}
	return returning, nil
}
