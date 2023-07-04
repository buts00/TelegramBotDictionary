package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"math/rand"
	"strings"
)

func Database() *sql.DB {
	db, err := sql.Open("postgres", ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func AllWords(db *sql.DB, chatId int) []Word {
	rows, err := db.Query("SELECT  * FROM dictionary WHERE id = ($1)", chatId)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	defer rows.Close()
	words := []Word{}

	for rows.Next() {
		var word Word
		err = rows.Scan(&word.Id, &word.word, &word.translation)
		if err != nil {
			log.Fatal(err)
		}
		words = append(words, Word{word.Id, word.word, word.translation})
	}
	return words
}

func RandomWord(chatId int) string {
	allWords := AllWords(Database(), chatId)
	randIndx := rand.Intn(len(allWords))
	randWord := allWords[randIndx].word + " / " + allWords[randIndx].translation
	return randWord
}

func InsertWord(newWord string, chatId int) {
	newWordArr := strings.Split(newWord, " ")
	db := Database()
	defer db.Close()
	query, err := db.Prepare("INSERT INTO dictionary (id,word,translation) VALUES ($1,$2,$3)")
	defer query.Close()
	if err != nil {
		log.Fatal(err)
	}
	_, err = query.Exec(chatId, newWordArr[0], newWordArr[1])
	if err != nil {
		log.Fatal(err)
	}
}
