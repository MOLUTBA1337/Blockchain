package main

import (
	"database/sql"
	"fmt"
	"log"

	"miner/firstminer/checkblock"
	config "miner/firstminer/config"
	"miner/firstminer/createblock"

	"github.com/gin-gonic/gin"
)

func init() {
	log.Println("Инициализация БД...")

	var err error

	config.DB, err = sql.Open("sqlite3", "blockchain.db?_timeout=5000")
	if err != nil {
		log.Fatal("Ошибка открытия БД:", err)
	}

	if err := config.DB.Ping(); err != nil {
		log.Fatal("Ошибка проверки соединения:", err)
	}

	log.Println("БД инициализирована успешно")

	createBlocks := `CREATE TABLE IF NOT EXISTS blocks (
		"index" INTEGER PRIMARY KEY,
		hash TEXT,
		prev_hash TEXT,
		timestamp INTEGER,
		"transaction" TEXT,
		nonce INTEGER
	);`
	createMiners := `CREATE TABLE IF NOT EXISTS miners(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	ip TEXT,
	port ITNEGER
	)`
	if _, err := config.DB.Exec(createBlocks); err != nil {
		fmt.Println("error creating blocks table: %w", err)
	}

	if _, err := config.DB.Exec(createMiners); err != nil {
		fmt.Println("ошибка создания: %w", err)
	}
	_, err = config.DB.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		log.Fatal("Не удалось включить WAL-режим")
	}

}

func main() {
	//ЗДЕСЬ СВОЕ ИМЯ НА САЙТЕ ДЛЯ ПОЛУЧЕНИЯ КОМИССИИ
	name := "Anton"
	r := gin.Default()

	main := r.Group("/")

	main.POST("/newblock", func(c *gin.Context) {
		createblock.CreateBlock(c, name)
	})
	main.POST("/checkblock", checkblock.CheckBlock)

	r.Run("0.0.0.0:8081")
}
