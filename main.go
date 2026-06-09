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
	fmt.Println("--- Tentative de Connexion de l'Utilisateur ---")
	pseudoConnexion := "Mbappe2026"
	motDePasseSaisi := "AllezLeFoot2026!"
	utilisateurConnecte, err := repoUtilisateur.VerifierConnexion(pseudoConnexion, motDePasseSaisi)
	if err != nil {
		log.Fatalf("Échec de la connexion : %v", err)
	}
	fmt.Printf("Authentification réussie pour : %s (Rôle ID : %d)\n", utilisateurConnecte.Pseudo, utilisateurConnecte.IDRole)
	token := SQL.GenererTokenSession(utilisateurConnecte)	
	fmt.Println("\n--- Clé de Session Générée ---")
	fmt.Printf("Token de session : %s\n", token)
}
