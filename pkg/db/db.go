package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var (
	db_host     = os.Getenv("DB_HOST")
	db_port     = os.Getenv("DB_PORT")
	db_user     = os.Getenv("DB_USER")
	db_password = os.Getenv("DB_PASSWORD")
	db_name     = os.Getenv("DB_NAME")
)

const (
	//QUERIES
	ADD_QUERY_STRING    = `INSERT INTO public.users (id) VALUES ($1)`
	DELETE_QUERY_STRING = `DELETE FROM public.users WHERE id=($1)`

	GET_QUERY_STRING     = `SELECT * FROM public.users WHERE id=($1)`
	GET_ALL_QUERY_STRING = `SELECT id FROM public.users`
)

var connection = fmt.Sprintf(
	"host=%s port=%s user=%s password=%s dbname=%s",
	db_host, db_port, db_user, db_password, db_name)

const driver = "postgres"

func openConnection() (*sql.DB, error) {
	db, err := sql.Open(driver, connection)
	if err != nil {
		log.Printf("Unable to connect to database: %s", err.Error())
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Printf("Unable to connect to database. Ping failed: %s", err.Error())
		return nil, err
	}
	log.Println("Successfully connected to Database")
	return db, nil
}

func CheckIfUserExists(chatID int64) (bool, error) {
	db, err := openConnection()
	defer db.Close()
	if err != nil {
		return false, err
	}
	result := db.QueryRow(GET_QUERY_STRING, chatID)
	err = result.Scan(&chatID)
	switch err {
	case sql.ErrNoRows:
		log.Printf("User %d doesn't exists", chatID)
		return false, nil
	case nil:
		log.Printf("User %d exists", chatID)
		return true, nil
	default:
		return false, err
	}
}

func AddUserToDB(chatID int64) error {
	db, err := openConnection()
	defer db.Close()
	if err != nil {
		return err
	}
	_, err = db.Exec(ADD_QUERY_STRING, chatID)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("ERROR Can't ADD user to DB")
	}
	log.Printf("User %d added", chatID)
	return nil
}

func DeleteUserFromDB(chatID int64) error {
	db, err := openConnection()
	defer db.Close()
	if err != nil {
		return err
	}
	_, err = db.Exec(DELETE_QUERY_STRING, chatID)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("ERROR Can't ADD user to DB")
	}
	log.Printf("User %d deleted", chatID)
	return nil
}

func GetUsersList(psqlInfo string) ([]int64, error) {
	db, err := openConnection()
	defer db.Close()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(GET_ALL_QUERY_STRING)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("ERROR Can't get users from db")
	}
	userList := []int64{}
	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			log.Println(err)
			return nil, fmt.Errorf("ERROR Can't get users from db")
		}
		userList = append(userList, id)
	}
	return userList, nil
}
