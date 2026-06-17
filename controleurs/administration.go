package controleurs

import (
	"footconnect/SQL"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type AdministrationControleur struct {
	RepoUser *SQL.UtilisateurRepository
}

func InitialiserAdministrationControleur(ru *SQL.UtilisateurRepository) *AdministrationControleur {
	return &AdministrationControleur{RepoUser: ru}
}

func (c *AdministrationControleur) AfficherTableauDeBord(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Println("Accès refusé : Aucun cookie trouvé.")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("<p>Accès refusé : Vous devez être connecté et Administrateur pour voir cette page.</p><p><a href='/forum/accueil'>Retourner à l'accueil</a></p>"))
		return
	}
	parties := strings.Split(cookie.Value, "_")
	if len(parties) < 5 {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("<p><a href='/authentification/connexion'>Veuillez vous reconnecter</a></p>"))
		return
	}
	idRole := parties[len(parties)-1]
	if idRole != "2" {
		log.Println("Accès refusé : L'utilisateur n'est pas Administrateur.")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("<p>Accès refusé : Vous devez être Administrateur du site FootConnect pour accéder au panneau de modération.</p><p><a href='/forum/accueil'>Retourner à l'accueil</a></p>"))
		return
	}
	tmpl, errTemplate := template.ParseFiles("html/administration/tableau_de_bord.html")
	if errTemplate != nil {
		log.Println("Erreur : Impossible de charger le fichier HTML :", errTemplate)
		http.Error(w, "Fichier HTML du tableau de bord introuvable. Vérifie le nom du fichier !", http.StatusInternalServerError)
		return
	}
	errExecute := tmpl.Execute(w, nil)
	if errExecute != nil {
		log.Println("Erreur lors du rendu du template :", errExecute)
	}
}

func (c *AdministrationControleur) BannirUtilisateur(w http.ResponseWriter, r *http.Request) {
	pseudoDuBanni := r.FormValue("pseudo_a_bannir")
	log.Println("L'administrateur a banni le joueur :", pseudoDuBanni)
	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}
