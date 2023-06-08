package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/go-sql-driver/mysql"
)

func main() {
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "password",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "data-access",
		AllowNativePasswords: true,
	}

	var err error

	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		fmt.Println(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		fmt.Println(pingErr)
	}

	fmt.Println("Connected!")
	User_Profil := Data_User{
		Name:     "Random",
		Password: "*******",
		Admin:    "false",
	}

	WichCategorieAreWeIn := Data_Topic{
		CategorieID: 0,
	}

	Registration := template.Must(template.ParseFiles("../html/Registration.html"))
	http.HandleFunc("/inscrire", func(w http.ResponseWriter, r *http.Request) {

		User := Data_User{
			Name:     r.FormValue("Username"),
			Email:    r.FormValue("Email"),
			Password: r.FormValue("Password"),
			Admin:    "false",
			Date:     DateObject(),
		}

		NewUser, err := AddUser(User)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(NewUser)
		Registration.Execute(w, User)
	})

	Connection := template.Must(template.ParseFiles("../html/Connection.html"))
	http.HandleFunc("/connecter", func(w http.ResponseWriter, r *http.Request) {
		User := Data_User{
			Name:     r.FormValue("Username"),
			Password: r.FormValue("Password"),
		}
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("user connecté: ", User.Name)
		fmt.Println("son mdp: ", User.Password)
		Verif := Data_Verif{
			Connect_Verif: VerifUser(User.Name, User.Password),
		}
		Verif.Connect_Verif = VerifUser(User.Name, User.Password)
		if Verif.Connect_Verif {
			User_Profil = User
			UserConnection(User_Profil)
			fmt.Println("Date profil: ", User_Profil.Date)
			fmt.Println("bio profil: ", User_Profil.Biography)
			http.Redirect(w, r, "/principal", http.StatusSeeOther)
		}
		Connection.Execute(w, User_Profil)
	})

	MainPage := template.Must(template.ParseFiles("../html/MainPage.html"))
	http.HandleFunc("/principal", func(w http.ResponseWriter, r *http.Request) {

		tab_topic1, err := PrintLastTopicCreated()
		if err != nil {
			fmt.Println(err)
		}

		Topic1 := Data_Topic{
			Topic_History: tab_topic1,
		}

		MainPage.Execute(w, Topic1)
	})

	Profil := template.Must(template.ParseFiles("../html/Profil.html"))
	http.HandleFunc("/profil", func(w http.ResponseWriter, r *http.Request) {

		Profil.Execute(w, User_Profil)
	})

	D1 := template.Must(template.ParseFiles("../html/Souls1.html"))
	http.HandleFunc("/D1", func(w http.ResponseWriter, r *http.Request) {

		WichCategorieAreWeIn.CategorieID = TakeUrlReturnIDCategorie(r.URL.String())

		Topic1 := Data_Topic{
			Trier: r.FormValue("Trier"),
		}

		tab_topic1, err := TopicPrint(WichCategorieAreWeIn.CategorieID, Topic1.Trier)
		if err != nil {
			fmt.Println(err)
		}

		Topic1.Topic_History = tab_topic1

		D1.Execute(w, Topic1)

	})

	D2 := template.Must(template.ParseFiles("../html/Souls2.html"))
	http.HandleFunc("/D2", func(w http.ResponseWriter, r *http.Request) {

		WichCategorieAreWeIn.CategorieID = TakeUrlReturnIDCategorie(r.URL.String())

		Topic2 := Data_Topic{
			Trier: r.FormValue("Trier"),
		}

		tab_topic2, err := TopicPrint(WichCategorieAreWeIn.CategorieID, Topic2.Trier)
		if err != nil {
			fmt.Println(err)
		}
		Topic2.Topic_History = tab_topic2

		D2.Execute(w, Topic2)

	})

	D3 := template.Must(template.ParseFiles("../html/Souls3.html"))
	http.HandleFunc("/D3", func(w http.ResponseWriter, r *http.Request) {

		WichCategorieAreWeIn.CategorieID = TakeUrlReturnIDCategorie(r.URL.String())

		Topic3 := Data_Topic{
			Trier: r.FormValue("Trier"),
		}

		tab_topic3, err := TopicPrint(WichCategorieAreWeIn.CategorieID, Topic3.Trier)
		if err != nil {
			fmt.Println(err)
		}
		Topic3.Topic_History = tab_topic3

		D3.Execute(w, Topic3)

	})

	CréerTopic := template.Must(template.ParseFiles("../html/TopicCreation.html"))
	http.HandleFunc("/créertopic", func(w http.ResponseWriter, r *http.Request) {

		LeTopic := Data_Topic{
			Name:        r.FormValue("Name"),
			Tags:        r.Form["Tags"],
			CategorieID: WichCategorieAreWeIn.CategorieID,
		}
		fmt.Println(LeTopic.Tags)
		if LeTopic.Tags != nil {
			if LeTopic.Name != "" {
				topID, err := AddTopic(LeTopic)
				fmt.Printf("ID of added topic: %v\n", topID)
				if err != nil {
					fmt.Println(err)
					return
				}

				LeMessage := Data_Message{
					Content:    r.FormValue("Content"),
					Author:     User_Profil.Name,
					Date:       DateObject(),
					Topic_ID:   topID,
					Topic_Name: TakeIDReturnName(topID),
				}

				if LeMessage.Content != "" {
					messID, err := AddMessage(LeMessage)
					fmt.Printf("ID of added message: %v\n", messID)
					if err != nil {
						fmt.Println(err)
						return
					}

					url := TakeIDReturnUrlCategorie(LeTopic.CategorieID)
					http.Redirect(w, r, url, http.StatusSeeOther)
				} else {
					//Erreur si message vide
				}
			} else {
				//Erreur si nom de topic vide
			}
		}
		CréerTopic.Execute(w, LeTopic)

	})

	TemplateMessage := template.Must(template.ParseFiles("../html/TemplateMessage.html"))
	http.HandleFunc("/topic/", func(w http.ResponseWriter, r *http.Request) {

		LeTopicId := TakeUrlReturnIDTopic(r.URL.String())

		message := Data_Message{
			Content:    r.FormValue("Content"),
			Author:     User_Profil.Name,
			Date:       DateObject(),
			Topic_ID:   LeTopicId,
			Topic_Name: TakeIDReturnName(LeTopicId),
		}

		tab_messages, err := MessagesPrint(message.Topic_ID)
		if err != nil {
			fmt.Println(err)
		}

		message.Message_History = tab_messages

		// TopicName, err := SelectTopicName(14)
		// if err != nil {
		// 	fmt.Println(err)
		// 	message.Topic_Name = "Error: name couldn't be found"
		// } else {
		// 	message.Topic_Name = TopicName
		// }

		if message.Content != "" {
			messID, err := AddMessage(message)
			fmt.Printf("ID of added message: %v\n", messID)
			if err != nil {
				fmt.Println(err)
				return
			}
			url := r.URL.String()
			http.Redirect(w, r, url, http.StatusSeeOther)
		}

		TemplateMessage.Execute(w, message)
	})

	css := http.FileServer(http.Dir("../css/"))
	http.Handle("/css/", http.StripPrefix("/css/", css))

	images := http.FileServer(http.Dir("../images/"))
	http.Handle("/images/", http.StripPrefix("/images/", images))

	js := http.FileServer(http.Dir("../js/"))
	http.Handle("/js/", http.StripPrefix("/js/", js))

	http.ListenAndServe(":8080", nil)

}
