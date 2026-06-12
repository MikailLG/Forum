package main

import (
	"fmt"
	"footconnect/SQL"
	"footconnect/database"
	"log"
	"net/http"
)

func main() {
	contexteBDD := database.InitialiserBDD()
	defer contexteBDD.Close()
	repoUser := SQL.InitialiserUtilisateurRepository(contexteBDD)
	repoDiscu := SQL.InitialiserDiscussionRepository(contexteBDD)
	repoComm := SQL.InitialiserCommentaireRepository(contexteBDD)
	_, _, _ = repoUser, repoDiscu, repoComm
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, "<h2>Serveur FootConnect opérationnel !")
	})
	port := ":8080"
	log.Printf("Le serveur FootConnect a démarré sur l'adresse http://localhost%s\n", port)
	erreurServeur := http.ListenAndServe(port, nil)
	if erreurServeur != nil {
		log.Fatalf("Échec du lancement du serveur HTTP : %v", erreurServeur)
	}
}
