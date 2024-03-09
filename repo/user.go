package repo

import "database/sql"

type User struct {
	ID       int
	Username string
	Email    string
	Name     string
}
type UserRepo interface {
	Migrate() error
	FindAll() ([]User, error)
	FindById(id int) (User, error)
	Create(user User) error
	Update(user User) error
	Delete(id int) error
}

type SQLiteUserRepo struct {
	Db *sql.DB
}

func NewSQLiteUserRepo(db *sql.DB) *SQLiteUserRepo {
	return &SQLiteUserRepo{Db: db}
}

func (r *SQLiteUserRepo) FindAll() ([]User, error) {
	rows, err := r.Db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Name)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *SQLiteUserRepo) Migrate() error {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT,
			email TEXT,
			name TEXT
		);
	`
	_, err := r.Db.Exec(query)
	return err
}

func (r *SQLiteUserRepo) FindById(id int) (User, error) {
	user := User{}
	err := r.Db.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(&user.ID, &user.Username, &user.Email, &user.Name)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *SQLiteUserRepo) Create(user User) error {
	_, err := r.Db.Exec("INSERT INTO users (username, email, name) VALUES (?, ?, ?)", user.Username, user.Email, user.Name)
	if err != nil {
		return err
	}
	return nil
}

func (r *SQLiteUserRepo) Update(user User) error {
	_, err := r.Db.Exec("UPDATE users SET username = ?, email = ?, name = ? WHERE id = ?", user.Username, user.Email, user.Name, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *SQLiteUserRepo) Delete(id int) error {
	_, err := r.Db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

type UserService struct {
	Repo UserRepo
}

func NewUserService(repo UserRepo) *UserService {
	return &UserService{repo}
}
