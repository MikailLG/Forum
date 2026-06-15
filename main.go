package main

import (
	"footconnect/SQL"
	"footconnect/controleurs"
	"footconnect/database"
	"footconnect/routeurs"
	"footconnect/services"
	"log"
	"net/http"
)

func main() {
	contexteBDD := database.InitialiserBDD()
	defer contexteBDD.Close()
	repoDiscu := SQL.InitialiserDiscussionRepository(contexteBDD)
	repoComm := SQL.InitialiserCommentaireRepository(contexteBDD)
	serviceForum := services.InitialiserForumService(repoDiscu, repoComm)
	ctrlForum := controleurs.InitialiserForumControleur(serviceForum)
	routeurs.ConfigurerRoutesForum(ctrlForum)
	ctrlAuth := controleurs.InitialiserAuthentificationControleur()
	routeurs.ConfigurerRoutesAuthentification(ctrlAuth)
	log.Println("Serveur en ligne : http://localhost:8080/forum/accueil")
	http.ListenAndServe(":8080", nil)
}
