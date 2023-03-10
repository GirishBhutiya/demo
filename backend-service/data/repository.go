package data

type Repository interface {
	GetAll() ([]*User, error)
	GetByEmail(email string) (*User, error)
	GetOne(id int) (*User, error)
	Update(user User) error
	DeleteByID(id int) error
	Insert(user User) (User, error)
	ResetPassword(password string, user User) error
	PasswordMatches(plainText string, user User) (bool, error)
}
