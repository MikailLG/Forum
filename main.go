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
	repoUtilisateur := SQL.InitialiserUtilisateurRepository(contexteBDD)
	nouveauSupporter := SQL.Utilisateur{
		Pseudo:          "Zizou10",
		Email:           "zidane@foot.fr",
		MotDePasseHache: "un_super_hash_sha512_ici",
		ClubDeCoeur:     "Real Madrid",
		IDRole:          1,
	}
	id, err := repoUtilisateur.CreerUtilisateur(nouveauSupporter)
	if err != nil {
		log.Printf("Inscription Échoué : %v", err)
		return
	}
	fmt.Printf("Inscription Réussi ! Utilisateur créé avec l'ID : %d\n", id)
}
