package cron

import (
	"log"
	"notif-me/env"
)

func GetAll(env *env.Env) ([]MangaUpdate, error) {
	log.Println("Get all manga update...")

	var records []MangaUpdate
	res := env.Db.Find(&records)

	if res.Error != nil {
		return nil, res.Error
	}

	return records, nil
}

func Create(env *env.Env, payload MangaUpdate) error {
	res := env.Db.Create(payload)

	return res.Error
}

func UpdateById(env *env.Env, id string, payload map[string]interface{}) error {
	res := env.Db.Model(&MangaUpdate{}).Where("id = ?", id).Updates(payload)

	return res.Error
}
