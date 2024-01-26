package database

import "github.com/AvivKermann/Chirpy/models"

func (db *DB) CreateNewRefreshToken(token string) error {
	dbContent, err := db.loadDB()
	if err != nil {
		return err
	}

	newToken := models.RefreshToken{
		Token:    token,
		IsActive: true,
	}
	dbContent.RefreshTokens[token] = newToken

	err = db.writeDB(dbContent)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) RevokeRefreshToken(token string) bool {
	dbContent, err := db.loadDB()

	if err != nil {
		return false
	}

	revokedToken := dbContent.RefreshTokens[token]
	revokedToken.IsActive = false
	dbContent.RefreshTokens[token] = revokedToken
	db.writeDB(dbContent)
	return true
}

func (db *DB) RefreshTokenExistAndActive(token string) bool {
	dbContent, err := db.loadDB()
	if err != nil {
		return false
	}
	dbToken, exist := dbContent.RefreshTokens[token]

	return exist && dbToken.IsActive
}
