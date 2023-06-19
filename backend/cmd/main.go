package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Category struct {
	ID   int
	Name string
}

type Topic struct {
	ID          int
	CategoryID  int
	Title       string
	Category    string
	MessageList []Message
}

type Message struct {
	ID       int
	Content  string
	Username string
}

var db *sql.DB

func main() {
	var err error
	db, err = OpenDBConnection("root:@tcp(127.0.0.1:3306)/forum")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		LoginHandler(w, r, db)
	})

	fmt.Println("Server listening on http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func OpenDBConnection(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Afficher la page de connexion
		tmpl, err := template.ParseFiles("templates/connexion.html")
		if err != nil {
			log.Fatal(err)
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else if r.Method == "POST" {
		// Traiter les données de connexion
		username := r.FormValue("username")
		password := r.FormValue("password")

		if checkCredentials(username, password, db) {
			// Authentification réussie, redirigez vers une autre page
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		} else {
			// Authentification échouée, afficher un message d'erreur
			fmt.Fprintln(w, "Identifiants invalides")
		}
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func getCategories() ([]Category, error) {
	rows, err := db.Query("SELECT id_category, category_title FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func getTopics() ([]Topic, error) {
	rows, err := db.Query(`
		SELECT topics.id_topic, topics.id_category, topics.topic_title, categories.category_title, messages.id_message, messages.content, users.username
		FROM topics
		INNER JOIN categories ON topics.id_category = categories.id_category
		INNER JOIN messages ON topics.id_topic = messages.id_topic
		INNER JOIN users ON messages.id_user = users.id_user
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []Topic
	for rows.Next() {
		var topic Topic
		var message Message
		err := rows.Scan(&topic.ID, &topic.CategoryID, &topic.Title, &topic.Category, &message.ID, &message.Content, &message.Username)
		if err != nil {
			return nil, err
		}

		// Add the message to the corresponding topic
		for i, t := range topics {
			if t.ID == topic.ID {
				topics[i].MessageList = append(topics[i].MessageList, message)
				break
			}
		}

		// If the topic does not exist in the list, add it
		if len(topics) == 0 || topics[len(topics)-1].ID != topic.ID {
			topic.MessageList = append(topic.MessageList, message)
			topics = append(topics, topic)
		}
	}

	return topics, nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if checkCredentials(username, password, db) {
		// Authentification réussie, redirigez vers une autre page
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	} else {
		// Authentification échouée, redirigez vers une autre page d'erreur
		http.Redirect(w, r, "/login-error", http.StatusSeeOther)
	}
}

func checkCredentials(username, password string, db *sql.DB) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ? AND password = ?", username, password).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	return count > 0
}
