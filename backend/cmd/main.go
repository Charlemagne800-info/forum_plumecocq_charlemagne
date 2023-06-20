package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Category struct {
	ID          int
	Name        string
	Description string
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

type LoginPageData struct {
	ErrorMessage string
}

type RegistrationPageData struct {
	ErrorMessage string
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
	http.HandleFunc("/dashboard", DashboardHandler)
	http.HandleFunc("/category/", CategoryHandler)
	http.HandleFunc("/choiceTopic", ChoiceTopicHandler)
	http.HandleFunc("/registration", func(w http.ResponseWriter, r *http.Request) {
		RegistrationHandler(w, r, db)
	})
	fs := http.FileServer(http.Dir("../../frontend/styles"))
	http.Handle("/styles/", http.StripPrefix("/styles/", fs))

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
		tmpl, err := template.ParseFiles("../../frontend/templates/connection.html")
		if err != nil {
			log.Fatal(err)
		}

		data := LoginPageData{
			ErrorMessage: "",
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			log.Fatal(err)
		}
	} else if r.Method == "POST" {
		// Traiter les données de connexion
		username := r.FormValue("username")
		password := r.FormValue("password")

		if checkCredentials(username, password, db) {
			// Authentification réussie, redirigez vers la page de dashboard
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		} else {
			// Authentification échouée, afficher un message d'erreur
			tmpl, err := template.ParseFiles("../../frontend/templates/connection.html")
			if err != nil {
				log.Fatal(err)
			}

			data := LoginPageData{
				ErrorMessage: "Identifiants invalides",
			}

			err = tmpl.Execute(w, data)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer les données pour afficher le dashboard
	categories, err := getCategories()
	if err != nil {
		log.Fatal(err)
	}

	data := struct {
		Categories []Category
	}{
		Categories: categories,
	}

	tmpl, err := template.ParseFiles("../../frontend/templates/dashboard.html")
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
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
		found := false
		for i, t := range topics {
			if t.ID == topic.ID {
				topics[i].MessageList = append(topics[i].MessageList, message)
				found = true
				break
			}
		}
		if !found {
			topic.MessageList = append(topic.MessageList, message)
			topics = append(topics, topic)
		}
	}

	return topics, nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Extraire les valeurs du formulaire
	username := r.FormValue("username")
	password := r.FormValue("password")

	if checkCredentials(username, password, db) {
		// Authentification réussie, redirigez vers la page de dashboard
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return // Ajouter cette ligne pour arrêter l'exécution de la fonction après la redirection
	} else {
		// Authentification échouée, afficher un message d'erreur
		tmpl, err := template.ParseFiles("../../frontend/templates/connection.html")
		if err != nil {
			log.Fatal(err)
		}

		data := LoginPageData{
			ErrorMessage: "Identifiants invalides",
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			log.Fatal(err)
		}
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

func RegistrationHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == "GET" {
		// Afficher la page d'inscription
		tmpl, err := template.ParseFiles("../../frontend/templates/registration.html")
		if err != nil {
			log.Fatal(err)
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else if r.Method == "POST" {
		// Traiter les données d'inscription
		username := r.FormValue("pseudo")
		email := r.FormValue("email")
		password := r.FormValue("password")

		if !isUsernameAvailable(username, db) {
			// Nom d'utilisateur déjà pris, afficher un message d'erreur
			tmpl, err := template.ParseFiles("../../frontend/templates/registration.html")
			if err != nil {
				log.Fatal(err)
			}

			data := RegistrationPageData{
				ErrorMessage: "Nom d'utilisateur déjà pris",
			}

			err = tmpl.Execute(w, data)
			if err != nil {
				log.Fatal(err)
			}

			return
		}

		if !isEmailValid(email) {
			// Adresse e-mail invalide, afficher un message d'erreur
			tmpl, err := template.ParseFiles("../../frontend/templates/registration.html")
			if err != nil {
				log.Fatal(err)
			}

			data := RegistrationPageData{
				ErrorMessage: "Adresse e-mail invalide",
			}

			err = tmpl.Execute(w, data)
			if err != nil {
				log.Fatal(err)
			}

			return
		}

		if len(password) < 10 {
			// Mot de passe trop court, afficher un message d'erreur
			tmpl, err := template.ParseFiles("../../frontend/templates/registration.html")
			if err != nil {
				log.Fatal(err)
			}

			data := RegistrationPageData{
				ErrorMessage: "Mot de passe trop court (min 10 caractères))",
			}

			err = tmpl.Execute(w, data)
			if err != nil {
				log.Fatal(err)
			}

			return
		}

		// Créer l'utilisateur et l'ajouter à la base de données
		if createUser(username, email, password, db) {
			// Utilisateur créé avec succès, redirigez vers la page de connexion
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		} else {
			// Échec de la création de l'utilisateur, afficher un message d'erreur
			tmpl, err := template.ParseFiles("../../frontend/templates/registration.html")
			if err != nil {
				log.Fatal(err)
			}

			data := RegistrationPageData{
				ErrorMessage: "Échec de la création de l'utilisateur",
			}

			err = tmpl.Execute(w, data)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func isUsernameAvailable(username string, db *sql.DB) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	return count == 0
}

func isEmailValid(email string) bool {
	return strings.Contains(email, "@")
}

func createUser(username, email, password string, db *sql.DB) bool {
	_, err := db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, password)
	if err != nil {
		log.Fatal(err)
		return false
	}

	return true
}

func CategoryHandler(w http.ResponseWriter, r *http.Request) {
	categoryID := strings.TrimPrefix(r.URL.Path, "/category/")
	categoryID = strings.TrimSuffix(categoryID, "/")

	// Redirect to the choiceTopic.html page with the category ID as a URL parameter
	http.Redirect(w, r, "/choiceTopic?category="+categoryID, http.StatusSeeOther)
}

func getCategoryByID(categoryID string) (Category, error) {
	var category Category
	err := db.QueryRow("SELECT id_category, category_title, category_description FROM categories WHERE id_category = ?", categoryID).Scan(&category.ID, &category.Name, &category.Description)
	if err != nil {
		return Category{}, err
	}

	return category, nil
}

func getTopicsByCategoryID(categoryID string) ([]Topic, error) {
	rows, err := db.Query(`
		SELECT topics.id_topic, topics.id_category, topics.topic_title, categories.category_title, messages.id_message, messages.content, users.username
		FROM topics
		INNER JOIN categories ON topics.id_category = categories.id_category
		INNER JOIN messages ON topics.id_topic = messages.id_topic
		INNER JOIN users ON messages.id_user = users.id_user
		WHERE topics.id_category = ?
		`, categoryID)
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
		found := false
		for i, t := range topics {
			if t.ID == topic.ID {
				topics[i].MessageList = append(topics[i].MessageList, message)
				found = true
				break
			}
		}
		if !found {
			topic.MessageList = append(topic.MessageList, message)
			topics = append(topics, topic)
		}
	}

	return topics, nil
}

func ChoiceTopicHandler(w http.ResponseWriter, r *http.Request) {
	categoryID := r.URL.Query().Get("category")

	// Retrieve the category information
	category, err := getCategoryByID(categoryID)
	if err != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	// Retrieve the topics for the selected category
	topics, err := getTopicsByCategoryID(categoryID)
	if err != nil {
		http.Error(w, "Failed to get topics", http.StatusInternalServerError)
		return
	}

	// Pass the data to the choiceTopic.html template for rendering
	data := struct {
		Category Category
		Topics   []Topic
	}{
		Category: category,
		Topics:   topics,
	}

	tmpl, err := template.ParseFiles("../../frontend/templates/choiceTopic.html")
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}
