package data

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const dbTimeout = time.Second * 3

var db *sql.DB

type postgresRepository struct {
	Conn *sql.DB
}

func NewPostgresRepository(pool *sql.DB) *postgresRepository {
	db = pool
	return &postgresRepository{
		Conn: pool,
	}
}

// New is the function used to create an instance of the data package. It returns the type
// Model, which embeds all the types we want to be available to our application.
// func New(dbPool *sql.DB) Models {
// 	db = dbPool

// 	return Models{
// 		User: User{},
// 	}
// }

// Models is the type for this package. Note that any model that is included as a member
// in this type is available to us throughout the application, anywhere that the
// app variable is used, provided that the model is also added in the New function.
// type Models struct {
// 	User User
// }

// User is the structure which holds one user from the database.
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name,omitempty"`
	Password  string    `json:"-"`
	Active    int       `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Post struct {
	ID        int       `json:"post_id"`
	Title     string    `json:"post_title"`
	Content   string    `json:"post_content"`
	UserId    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetAllPost returns a slice of all posts, sorted by userid
func (u *postgresRepository) GetAllPost() ([]*Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, title, content, user_id,created_at, updated_at
	from posts order by user_id`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post

	for rows.Next() {
		var pst Post
		err := rows.Scan(
			&pst.ID,
			&pst.Title,
			&pst.Content,
			&pst.UserId,
			&pst.CreatedAt,
			&pst.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		posts = append(posts, &pst)
	}

	return posts, nil
}

// GetOnePost returns one post by id
func (u *postgresRepository) GetOnePost(id int) (*Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, title, content, user_id,created_at, updated_at from posts where id = $1`

	var pst Post
	row := db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&pst.ID,
		&pst.Title,
		&pst.Content,
		&pst.UserId,
		&pst.CreatedAt,
		&pst.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &pst, nil
}

// DeletePost deletes one post from the database, by Post.ID
func (u *postgresRepository) DeletePost(pst Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from posts where id = $1`

	_, err := db.ExecContext(ctx, stmt, pst.ID)
	if err != nil {
		return err
	}

	return nil
}

// InsertPost inserts a new post into the database, and returns the ID of the newly inserted row
func (u *postgresRepository) InsertPost(pst Post) (Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `insert into posts (title, content, user_id, created_at, updated_at)
		values ($1, $2, $3, $4, $5) returning id`

	var post Post
	row := db.QueryRowContext(ctx, stmt,
		pst.Title,
		pst.Content,
		pst.UserId,
		time.Now(),
		time.Now(),
	)

	err := row.Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.UserId,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return Post{}, err
	}

	return post, nil

}

// GetAll returns a slice of all users, sorted by last name
func (u *postgresRepository) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at
	from users order by last_name`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FullName,
			&user.Password,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

// GetByEmail returns one user by email
func (u *postgresRepository) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at from users where email = $1`

	var user User
	row := db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FullName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetOne returns one user by id
func (u *postgresRepository) GetOne(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at from users where id = $1`

	var user User
	row := db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FullName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update updates one user in the database, using the information
// stored in the receiver u
func (u *postgresRepository) Update(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update users set
		email = $1,
		full_name = $2,
		user_active = $3,
		updated_at = $4
		where id = $5
	`

	_, err := db.ExecContext(ctx, stmt,
		user.Email,
		user.FullName,
		user.Active,
		time.Now(),
		user.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// Delete deletes one user from the database, by User.ID
func (u *postgresRepository) Delete(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from users where id = $1`

	_, err := db.ExecContext(ctx, stmt, user.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteByID deletes one user from the database, by ID
func (u *postgresRepository) DeleteByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from users where id = $1`

	_, err := db.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

// Insert inserts a new user into the database, and returns the ID of the newly inserted row
func (u *postgresRepository) Insert(user User) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return User{}, err
	}

	stmt := `insert into users (email, full_name, password, user_active, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6) returning id`

	var usr User
	row := db.QueryRowContext(ctx, stmt,
		user.Email,
		user.FullName,
		hashedPassword,
		user.Active,
		time.Now(),
		time.Now(),
	)

	err = row.Scan(
		&usr.ID,
		&usr.Email,
		&usr.FullName,
		&usr.Password,
		&usr.Active,
		&usr.CreatedAt,
		&usr.UpdatedAt,
	)

	if err != nil {
		return User{}, err
	}

	return user, nil

}

// ResetPassword is the method we will use to change a user's password.
func (u *postgresRepository) ResetPassword(password string, user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `update users set password = $1 where id = $2`
	_, err = db.ExecContext(ctx, stmt, hashedPassword, user.ID)
	if err != nil {
		return err
	}

	return nil
}

// PasswordMatches uses Go's bcrypt package to compare a user supplied password
// with the hash we have stored for a given user in the database. If the password
// and hash match, we return true; otherwise, we return false.
func (u *postgresRepository) PasswordMatches(plainText string, user User) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
