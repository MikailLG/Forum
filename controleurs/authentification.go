package controleurs

import (
	"footconnect/SQL"
	"html/template"
	"log"
	"net/http"
	"os"
)

type AuthentificationControleur struct {
	RepoUser *SQL.UtilisateurRepository
}

func InitialiserAuthentificationControleur(ru *SQL.UtilisateurRepository) *AuthentificationControleur {
	return &AuthentificationControleur{RepoUser: ru}
}

func (c *AuthentificationControleur) AfficherConnexion(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("html/authentification/connexion.html")
	if err != nil {
		log.Println("Erreur template connexion :", err)
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func (c *AuthentificationControleur) AfficherInscription(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("html/authentification/inscription.html")
	if err != nil {
		log.Println("Erreur template inscription :", err)
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func (c *AuthentificationControleur) SeConnecter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	email := r.FormValue("email")
	motDePasse := r.FormValue("not_de_passe")
	u, err := c.RepoUser.VerifierConnexion(email, motDePasse)
	if err != nil {
		log.Println("Échec de connexion pour :", email)
		data := struct{ Erreur string }{Erreur: "Identifiants incorrects."}
		tmpl, _ := template.ParseFiles("html/authentification/connexion.html")
		tmpl.Execute(w, data)
		return
	}
	estEnProd := os.Getenv("PORT") != ""
	token := SQL.GenererTokenSession(u)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   estEnProd,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   3600 * 24,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "pseudo",
		Value:    u.Pseudo,
		Path:     "/",
		HttpOnly: true,
		Secure:   estEnProd,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   3600 * 24,
	})
	log.Println("Session et cookie créés pour l'utilisateur :", u.Pseudo)
	http.Redirect(w, r, "/forum/accueil", http.StatusSeeOther)
}

func (c *AuthentificationControleur) SInscrire(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	pseudo := r.FormValue("pseudo")
	email := r.FormValue("email")
	motDePasse := r.FormValue("not_de_passe")
	club := r.FormValue("club")
	nouvelUtilisateur := SQL.Utilisateur{
		Pseudo:      pseudo,
		Email:       email,
		ClubDeCoeur: club,
		IDRole:      1,
	}
	err := c.RepoUser.CreerUtilisateur(nouvelUtilisateur, motDePasse)
	if err != nil {
		log.Println("Inscription rejetée :", err)
		data := struct{ Erreur string }{Erreur: err.Error()}
		tmpl, _ := template.ParseFiles("html/authentification/inscription.html")
		tmpl.Execute(w, data)
		return
	}
	log.Println("Nouveau utilisateur inscrit sur FootConnect :", pseudo)
	http.Redirect(w, r, "/authentification/connexion", http.StatusSeeOther)
}
