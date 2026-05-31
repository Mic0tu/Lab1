package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"notebook/internal/domain"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateUser(ctx context.Context, u domain.User) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO users (id, username, email) VALUES ($1, $2, $3)`,
		u.ID, u.Username, u.Email,
	)
	return err
}

func (r *Repository) GetUser(ctx context.Context, id uuid.UUID) (domain.User, error) {
	var u domain.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, username, email FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Username, &u.Email)
	return u, err
}

func (r *Repository) ListUsers(ctx context.Context) ([]domain.User, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, username, email FROM users ORDER BY username`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *Repository) UpdateUser(ctx context.Context, u domain.User) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE users SET username = $2, email = $3 WHERE id = $1`,
		u.ID, u.Username, u.Email,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) CreateProfile(ctx context.Context, p domain.Profile) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO profiles (id, name, lname) VALUES ($1, $2, $3)`,
		p.ID, p.Name, p.Lname,
	)
	return err
}

func (r *Repository) GetProfile(ctx context.Context, id uuid.UUID) (domain.Profile, error) {
	var p domain.Profile
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, lname FROM profiles WHERE id = $1`, id,
	).Scan(&p.ID, &p.Name, &p.Lname)
	return p, err
}

func (r *Repository) ListProfiles(ctx context.Context) ([]domain.Profile, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, lname FROM profiles ORDER BY lname, name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []domain.Profile
	for rows.Next() {
		var p domain.Profile
		if err := rows.Scan(&p.ID, &p.Name, &p.Lname); err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}
	return profiles, rows.Err()
}

func (r *Repository) UpdateProfile(ctx context.Context, p domain.Profile) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE profiles SET name = $2, lname = $3 WHERE id = $1`,
		p.ID, p.Name, p.Lname,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) DeleteProfile(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM profiles WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
