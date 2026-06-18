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
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/web/", http.StripPrefix("/web/", fs))
	repoDiscu := SQL.InitialiserDiscussionRepository(contexteBDD)
	repoComm := SQL.InitialiserCommentaireRepository(contexteBDD)
	serviceForum := services.InitialiserForumService(repoDiscu, repoComm)
	ctrlForum := controleurs.InitialiserForumControleur(serviceForum)
	routeurs.ConfigurerRoutesForum(ctrlForum)
	repoUser := SQL.InitialiserUtilisateurRepository(contexteBDD) 
	ctrlAuth := controleurs.InitialiserAuthentificationControleur(repoUser)
	routeurs.ConfigurerRoutesAuthentification(ctrlAuth)
	ctrlAdmin := controleurs.InitialiserAdministrationControleur(repoUser)
	http.HandleFunc("/admin/dashboard", ctrlAdmin.AfficherTableauDeBord)
	http.HandleFunc("/admin/bannir", ctrlAdmin.BannirUtilisateur)
	log.Println("FootConnect est en ligne : http://localhost:8080/forum/accueil")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Le serveur s'est arrêté brutalement : ", err)
	}
}
