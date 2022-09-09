package services

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type databaseJokeService struct {
	db *sql.DB
}

var _ JokeService = &databaseJokeService{}

func NewDatabaseJokeService(db *sql.DB) JokeService {
	return &databaseJokeService{db: db}
}

func (s *databaseJokeService) Get(uuid string) (*Joke, error) {
	query := "SELECT uuid, joke FROM jokes WHERE uuid = $1 LIMIT 1"
	row := s.db.QueryRow(query, uuid)

	var joke Joke
	err := row.Scan(&joke.UUID, &joke.Joke)
	if err != nil {
		return nil, fmt.Errorf("could not find joke. Reason: %v", err)
	}
	return &joke, nil
}

func (s *databaseJokeService) Add(joke string) (*Joke, error) {
	id := uuid.New().String()
	query := "INSERT INTO jokes (uuid, joke) VALUES ($1, $2)"
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("could not create database transaction. Reason: %v", err)
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("invalid SQL statement. Reason: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, joke)
	if err != nil {
		return nil, fmt.Errorf("could not add joke to database. Reason: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("could not add joke to database. Reason: %v", err)
	}

	stmt.QueryRow(id, joke)
	return s.Get(id)
}

func (s *databaseJokeService) Random() (*Joke, error) {
	query := "SELECT uuid, joke FROM jokes ORDER BY RANDOM() LIMIT 1"
	row := s.db.QueryRow(query)

	var joke Joke
	err := row.Scan(&joke.UUID, &joke.Joke)
	if err != nil {
		return nil, fmt.Errorf("could not map database to joke. Reason: %v", err)
	}
	return &joke, nil
}