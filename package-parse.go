package main

import _ "github.com/mattn/go-sqlite3"
import "io/ioutil"

import "encoding/json"
import "log"
import "database/sql"

func main() {
	db, err := sql.Open("sqlite3", "./asset/melpa.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("DB Open.")

	log.Println("START parse archive.json")
	jsonBlob, err := ioutil.ReadFile("./asset/archive.json")
	if err != nil {
		log.Fatal(err)
	}

	var packages interface{}
	err = json.Unmarshal(jsonBlob, &packages)
	if err != nil {
		log.Fatal(err)
	}

	for name, body := range packages.(map[string]interface{}) {
		storeArchiveJSON(db, name, body.(map[string]interface{}))
	}
	log.Println("END parse archive.json")

	log.Println("START parse recipes.json")
	jsonBlob, err = ioutil.ReadFile("./asset/recipes.json")
	if err != nil {
		log.Fatal(err)
	}

	var recipes interface{}
	err = json.Unmarshal(jsonBlob, &recipes)
	if err != nil {
		log.Fatal(err)
	}

	for name, body := range recipes.(map[string]interface{}) {
		storeRecipeJSON(db, name, body.(map[string]interface{}))
	}
	log.Println("END parse recipes.json")

	log.Println("START parse download_counts.json")
	jsonBlob, err = ioutil.ReadFile("./asset/download_counts.json")
	if err != nil {
		log.Fatal(err)
	}

	var dlcnt interface{}
	err = json.Unmarshal(jsonBlob, &dlcnt)
	if err != nil {
		log.Fatal(err)
	}

	for name, cnt := range dlcnt.(map[string]interface{}) {
		storeCountJSON(db, name, int(cnt.(float64)))
	}
	log.Println("END parse download_count.json")
}

func storeArchiveJSON(db *sql.DB, name string, body map[string]interface{}) {
	var q string
	descStr := body["desc"].(string)
	typeStr := body["type"].(string)

	urlStr := ""
	maintainerStr := ""

	if body["props"] != nil {
		urlIF := body["props"].(map[string]interface{})["url"]
		if urlIF != nil {
			urlStr = urlIF.(string)
		}

		maintainerIF := body["props"].(map[string]interface{})["maintainer"]
		if maintainerIF != nil {
			maintainerStr = maintainerIF.(string)
		}
	}
	q = "REPLACE INTO packages (name, desc, type, url, maintainer) VALUES ($1, $2, $3, $4, $5)"

	if _, err := db.Exec(q, name, descStr, typeStr, urlStr, maintainerStr); err != nil {
		log.Println(err, q)
	}

	if body["deps"] != nil {
		for deps := range body["deps"].(map[string]interface{}) {
			q = "REPLACE INTO dependencies (name, depends) VALUES ($1, $2)"
			if _, err := db.Exec(q, name, deps); err != nil {
				log.Println(err, q)
			}
		}
	}

	if body["props"] != nil {
		if body["props"].(map[string]interface{})["keywords"] != nil {
			for _, keyword := range body["props"].(map[string]interface{})["keywords"].([]interface{}) {
				q = "REPLACE INTO keywords (name, keyword) VALUES ($1, $2)"
				if _, err := db.Exec(q, name, keyword.(string)); err != nil {
					log.Println(err, q)
				}
			}
		}

	}

	if body["props"] != nil {
		if body["props"].(map[string]interface{})["authors"] != nil {
			for _, author := range body["props"].(map[string]interface{})["authors"].([]interface{}) {
				q = "REPLACE INTO authors (name, author) VALUES ($1, $2)"
				if _, err := db.Exec(q, name, author.(string)); err != nil {
					log.Println(err, q)
				}
			}
		}

	}
}

func storeRecipeJSON(db *sql.DB, name string, body map[string]interface{}) {
	fetcherStr := body["fetcher"].(string)
	repoStr := ""
	if body["repo"] != nil {
		repoStr = body["repo"].(string)
	} else if body["url"] != nil {
		repoStr = body["url"].(string)
	}

	q := "UPDATE packages SET fetcher = $1, repo = $2 WHERE name = $3"
	if _, err := db.Exec(q, fetcherStr, repoStr, name); err != nil {
		log.Println(err, q)
	}
}

func storeCountJSON(db *sql.DB, name string, cnt int) {
	q := "UPDATE packages SET dl = $1 WHERE name = $2"
	if _, err := db.Exec(q, cnt, name); err != nil {
		log.Println(err, q)
	}
}
