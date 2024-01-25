package database

import "github.com/AvivKermann/Chirpy/models"

func (db *DB) GetChirps() ([]models.Chirp, error) {
	dbContent, err := db.loadDB()

	db.mu.RLock()
	defer db.mu.RUnlock()
	chirps := []models.Chirp{}

	if err != nil {
		return chirps, err
	}

	for _, chirp := range dbContent.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil

}

func (db *DB) GetSingleChirp(chirpId int) (models.Chirp, bool) {
	dbContent, err := db.loadDB()

	db.mu.RLock()
	defer db.mu.RUnlock()
	if err != nil {
		return models.Chirp{}, false
	}

	chirp, exist := dbContent.Chirps[chirpId]
	if !exist {
		return models.Chirp{}, exist
	}

	return chirp, exist
}
func (db *DB) CreateChirp(content string) (models.Chirp, error) {
	dbContent, err := db.loadDB()

	if err != nil {
		return models.Chirp{}, err
	}

	index := len(dbContent.Chirps) + 1
	newChirp := models.Chirp{
		ID:   index,
		Body: content,
	}

	dbContent.Chirps[index] = newChirp
	err = db.writeDB(dbContent)
	if err != nil {
		return models.Chirp{}, err
	}

	return newChirp, nil

}
