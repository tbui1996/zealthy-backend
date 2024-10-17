package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/tbui1996/zealthy-backend/internal/core/domain"
	"github.com/tbui1996/zealthy-backend/internal/core/ports"
)

type postgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) ports.UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	log.Printf("Create method called with user: %+v", user)
	if r.db == nil {
		log.Println("Database connection is nil")
		return fmt.Errorf("database connection is nil")
	}

	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`

	log.Printf("Executing query: %s", query)
	log.Printf("With values: Email=%s, Password=<redacted>", user.Email)

	err := r.db.QueryRowContext(ctx, query, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return err
	}

	log.Printf("User created successfully with ID: %s", user.ID)
	return nil
}

func (r *postgresUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	var address domain.Address
	var aboutMe, birthdate sql.NullString
	var street, city, state, zip sql.NullString

	err := r.db.QueryRowContext(ctx, `
        SELECT id, email, password, about_me, 
               address_street, address_city, address_state, address_zip, 
               birthdate 
        FROM users 
        WHERE id = $1
    `, id).Scan(
		&user.ID, &user.Email, &user.Password, &aboutMe,
		&street, &city, &state, &zip,
		&birthdate,
	)

	if err != nil {
		return nil, err
	}

	if aboutMe.Valid {
		user.AboutMe = &aboutMe.String
	}
	if birthdate.Valid {
		user.Birthdate = &birthdate.String
	}

	if street.Valid || city.Valid || state.Valid || zip.Valid {
		address = domain.Address{}
		if street.Valid {
			address.Street = &street.String
		}
		if city.Valid {
			address.City = &city.String
		}
		if state.Valid {
			address.State = &state.String
		}
		if zip.Valid {
			address.Zip = &zip.String
		}
		user.Address = &address
	}

	return &user, nil
}

func (r *postgresUserRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	query := `SELECT id, email, password, about_me, address_street, address_city, address_state, address_zip, birthdate 
              FROM users`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*domain.User
	for rows.Next() {
		var user domain.User
		var address domain.Address
		err := rows.Scan(
			&user.ID, &user.Email, &user.Password, &user.AboutMe,
			&address.Street, &address.City, &address.State, &address.Zip, &user.Birthdate)
		if err != nil {
			return nil, err
		}
		user.Address = &address
		users = append(users, &user)
	}
	return users, nil
}

func (r *postgresUserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
        UPDATE users 
        SET email = $2, 
            password = $3, 
            about_me = $4, 
            address_street = $5, 
            address_city = $6, 
            address_state = $7, 
            address_zip = $8, 
            birthdate = $9 
        WHERE id = $1
    `

	var street, city, state, zip *string
	if user.Address != nil {
		street = user.Address.Street
		city = user.Address.City
		state = user.Address.State
		zip = user.Address.Zip
	}

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.Password,
		user.AboutMe,
		street,
		city,
		state,
		zip,
		user.Birthdate)

	return err
}

func (r *postgresUserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
