package store

import "2022_1_CJ/internal/app/model"

// UserRepository ...
type UserRepository struct {
	store *Store
}

// Create ...
func (r *UserRepository) Create(u *model.User) (*model.User, error) {
	if err := u.Validate(); err != nil {
		return nil, err
	}

	if err := u.BeforeCreate(); err != nil {
		return nil, err
	}

	if err := r.store.db.QueryRow(
		"INSERT INTO users (name, email, phone, encrypted_password) VALUES ($1, $2, $3, $4) RETURNING id",
		u.Name,
		u.Email,
		u.Phone,
		u.EncryptedPassword,
	).Scan(&u.ID); err != nil {
		return nil, err
	}

	return u, nil
}

// FindByEmail ...
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, name, email, phone, encrypted_password FROM users WHERE email = $1",
		email,
	).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Phone,
		&u.EncryptedPassword,
	); err != nil {
		return nil, err
	}

	return u, nil
}
