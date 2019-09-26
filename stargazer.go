package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"net/http"
)

type Repo struct {
	StargazerCount int `json:"stargazers_count"`
}

type Star struct {
	Name string
	Star int
}

func main() {
	db, err := sql.Open("sqlite3", "./asset/melpa.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("select name, 'https://api.github.com/repos/' || repo from packages where fetcher = 'github'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	q := "UPDATE packages SET star = $1 WHERE name = $2"

	var stars []Star

	for rows.Next() {
		var name string
		var url string
		err = rows.Scan(&name, &url)
		if err != nil {
			log.Printf("%v\n", err)
		}

		req, _ := http.NewRequest("GET", url, nil)
		//req.Header.Set("Authorization", "Basic `base64 user:password`")
		client := new(http.Client)
		resp, _ := client.Do(req)
		defer resp.Body.Close()
		byteArray, _ := ioutil.ReadAll(resp.Body)
		repo := new(Repo)
		json.Unmarshal(byteArray, repo)
		log.Printf("%v, %v, %v\n", url, name, repo.StargazerCount)

		star := Star{Name: name, Star: repo.StargazerCount}
		stars = append(stars, star)
		log.Printf("%v, %v\n", star.Name, star.Star)

	}
	for _, star := range stars {
		if _, err := db.Exec(q, star.Star, star.Name); err != nil {
			log.Println(err, q)
		}
	}
}
