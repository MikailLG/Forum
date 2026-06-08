package database

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

func InitialiserBDD() *sql.DB {

	utilisateur := "root"
	motDePasse := ""
	hote := "127.0.0.1"
	port := "3306"
	nomBase := "footconnect_db"
	chaineConnexion := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", utilisateur, motDePasse, hote, port, nomBase)
	contexteBDD, erreurOuverture := sql.Open("mysql", chaineConnexion)
	if erreurOuverture != nil {
		log.Fatalf("Erreur connexion base de donnees Erreur: \n\t%s", erreurOuverture.Error())
	}
	erreurPing := contexteBDD.Ping()
	if erreurPing != nil {
		contexteBDD.Close()
		log.Fatalf("Erreur ping base de donnees Erreur: \n\t%s", erreurPing.Error())
	}
	log.Printf("BDD : Connexion réussie")
	return contexteBDD
}
