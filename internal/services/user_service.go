package services

import (
	"finanapp/internal/db"
	"finanapp/internal/models"
	"log"
)

func CreateUser(user *models.User) error {
	if err := db.DB.Create(user).Error; err != nil {
		log.Printf("Erro ao criar usuário: %v", err)
		return err
	}
	return nil
}