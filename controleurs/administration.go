package controleurs

import (
	"footconnect/SQL"
	"html/template"
	"log"
	"net/http"
	"strconv"
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
		http.Error(w, "Accès interdit : connectez-vous.", http.StatusForbidden)
		return
	}
	parties := strings.Split(cookie.Value, "_")
	if len(parties) < 2 {
		http.Redirect(w, r, "/authentification/connexion", http.StatusSeeOther)
		return
	}
	idRole := parties[len(parties)-1]
	if idRole != "2" { 
		http.Error(w, "Accès interdit : réservé au Staff.", http.StatusForbidden)
		return
	}
	tmpl, errTemplate := template.ParseFiles("html/administration/tableau_de_bord.html")
	if errTemplate != nil {
		log.Println("Erreur template admin:", errTemplate)
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func (c *AdministrationControleur) BannirUtilisateur(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.FormValue("id_a_supprimer") 
	idDuBanni, err := strconv.Atoi(strings.TrimSpace(idStr))
	if err != nil {
		log.Println("Erreur conversion ID :", err)
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}
	errSQL := c.RepoUser.SupprimerUtilisateur(idDuBanni)
	if errSQL != nil {
		log.Println("Erreur lors de la suppression de l'utilisateur", idDuBanni, ":", errSQL)
	} else {
		log.Println("Utilisateur supprimé avec succès. ID :", idDuBanni)
	}
	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}
