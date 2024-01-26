package database

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"

	"github.com/AvivKermann/Chirpy/models"
)

type DBStructure struct {
	Chirps        map[int]models.Chirp           `json:"chirps"`
	Users         map[int]models.User            `json:"user"`
	RefreshTokens map[string]models.RefreshToken `json:"refresh_tokens"`
}
type DB struct {
	path string
	mu   sync.RWMutex
}

func NewDB(path string) (*DB, error) {

	dbContent := DBStructure{
		Chirps:        map[int]models.Chirp{},
		Users:         map[int]models.User{},
		RefreshTokens: map[string]models.RefreshToken{},
	}
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		_, err := os.Create(path)
		if err != nil {
			log.Fatal(err)
		}
		dbPointer := &DB{
			path: path,
		}
		dbPointer.writeDB(dbContent)
		return dbPointer, nil
	}
	err := os.Remove(path)
	if err != nil {
		log.Fatal(err)
	}
	_, err = os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	dbPointer := &DB{
		path: path,
	}
	dbPointer.writeDB(dbContent)
	return dbPointer, nil
}

func (db *DB) writeDB(dbContent DBStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	jsonDB, err := json.Marshal(dbContent)

	if err != nil {
		return err
	}
	err = os.WriteFile(db.path, jsonDB, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	byteData, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	dbContent := DBStructure{}

	err = json.Unmarshal(byteData, &dbContent)
	if err != nil {
		return DBStructure{}, err
	}

	return dbContent, nil
}
