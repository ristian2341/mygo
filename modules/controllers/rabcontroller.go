package controllers

import (
	"mygo/config"
	"mygo/modules/models"
)

func getDataRab(q string, code string) ([]models.Rab, error) {
	var dataRabs []models.Rab

	query := config.DB.Model(&dataRabs)

	// Kita buat wildcard string sekali saja
	searchTerm := "%" + q + "%"

	if q != "" {
		query = query.Where("CodeProject LIKE ? or email like ? or nama like ? ", searchTerm, searchTerm, searchTerm)
	}

	if code != "" {
		query = query.Where("code = ?", code)
	}

	err := query.Find(&dataRabs).Error

	return dataRabs, err

}
