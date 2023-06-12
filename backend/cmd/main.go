package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Category struct {
	ID   int
	Name string
}

func main() {
	// Informations de connexion à la base de données
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/forum")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Vérifier la connexion à la base de données
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	// Récupérer toutes les catégories
	rows, err := db.Query("SELECT id_category, category_title FROM categories")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var categories []Category

	// Parcourir les résultats de la requête
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			panic(err.Error())
		}
		categories = append(categories, category)
	}

	// Vérifier les éventuelles erreurs lors du parcours des résultats
	err = rows.Err()
	if err != nil {
		panic(err.Error())
	}

	// Afficher les catégories
	for _, category := range categories {
		fmt.Printf("ID: %d, Name: %s\n", category.ID, category.Name)
	}
}
