package main

import (
	"fmt"
	"footconnect/SQL"
	"footconnect/database"
	"log"
)

func main() {
	contexteBDD := database.InitialiserBDD()
	defer contexteBDD.Close()
	repoDiscussion := SQL.InitialiserDiscussionRepository(contexteBDD)
	repoCommentaire := SQL.InitialiserCommentaireRepository(contexteBDD)
	fmt.Println("--- 1. Creation d'un Sujet de Discussion ---")
	nouveauSujet := SQL.Discussion{
		Titre:         "Pronostics Ligue des Champions 2026",
		Contenu:       "Qui va soulever la coupe aux grandes oreilles cette année d'après vous ?",
		IDUtilisateur: 1,
		IDCategorie:   1,
	}
	idSujet, err := repoDiscussion.CreerDiscussion(nouveauSujet)
	if err != nil {
		log.Fatalf("Echec creation sujet : %v", err)
	}
	fmt.Printf("Sujet cree avec l'ID : %d\n", idSujet)
	fmt.Println("\n--- 2. Gestion de l'Etat du Sujet ---")
	errEtat := repoDiscussion.ModifierEtat(idSujet, "archivé")
	if errEtat != nil {
		log.Fatalf("Echec modification etat : %v", errEtat)
	}
	fmt.Println("Etat du sujet passe a : 'archivé' avec succes")
	fmt.Println("\n--- 3. Ajout d'un Commentaire lie au Sujet ---")
	reponse := SQL.Commentaire{
		Contenu:       "Moi je mise tout sur le Real Madrid comme d'habitude !",
		IDUtilisateur: 1,
		IDDiscussion:  idSujet,
	}
	idCommentaire, errComm := repoCommentaire.CreerCommentaire(reponse)
	if errComm != nil {
		log.Fatalf("Echec ajout commentaire : %v", errComm)
	}
	fmt.Printf("Commentaire ajoute avec l'ID : %d\n", idCommentaire)
}
