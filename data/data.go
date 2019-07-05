package data

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type User struct {
	Id       uint
	Uuid     string
	First    string
	Last     string
	Email    string
	Password string
	Birthday string
	Salt     string
}

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("postgres", "user=mark dbname=facebook password=zuckerberg sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}

func (u *User) Register() (err error) {
	statement := "INSERT INTO users (uuid, first_name, last_name, email, password, salt, birthday) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Query(u.Uuid, u.First, u.Last, u.Email, u.Password, u.Salt, u.Birthday)
	return
}

func User_by_email(email string) (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, first_name, last_name, email, birthday, password, salt FROM users WHERE email = $1", email).Scan(&user.Id, &user.Uuid, &user.First, &user.Last, &user.Email, &user.Birthday, &user.Password, &user.Salt)
	return
}
