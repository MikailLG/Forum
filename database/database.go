package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

func InitialiserBDD() *sql.DB {
	utilisateur := getEnv("DB_USER", "root")
	motDePasse := getEnv("DB_PWD", "")
	hote := getEnv("DB_HOST", "127.0.0.1")
	port := getEnv("DB_PORT", "3306")
	nomBase := getEnv("DB_NAME", "footconnect_db")
	chaineConnexion := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", utilisateur, motDePasse, hote, port, nomBase)
	contexteBDD, erreurOuverture := sql.Open("mysql", chaineConnexion)
	if erreurOuverture != nil {
		log.Fatalf("Erreur ouverture base de données : \n\t%s", erreurOuverture.Error())
	}
	erreurPing := contexteBDD.Ping()
	if erreurPing != nil {
		contexteBDD.Close()
		log.Fatalf("Erreur ping base de données : \n\t%s", erreurPing.Error())
	}
	contexteBDD.SetMaxOpenConns(10)
	contexteBDD.SetMaxIdleConns(5)
	log.Printf("BDD : Connexion réseau établie avec succès.")
	return contexteBDD
}

func getEnv(cle, valeurParDefaut string) string {
	valeur := os.Getenv(cle)
	if valeur == "" {
		return valeurParDefaut
	}
	return valeur
}
