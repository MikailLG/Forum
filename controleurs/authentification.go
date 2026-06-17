package controleurs

import (
	"footconnect/SQL"
	"html/template"
	"log"
	"net/http"
)

type AuthentificationControleur struct {
	RepoUser *SQL.UtilisateurRepository
}

func InitialiserAuthentificationControleur(ru *SQL.UtilisateurRepository) *AuthentificationControleur {
	return &AuthentificationControleur{RepoUser: ru}
}

func (c *AuthentificationControleur) AfficherConnexion(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("html/authentification/connexion.html")
	tmpl.Execute(w, nil)
}

func (c *AuthentificationControleur) AfficherInscription(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("html/authentification/inscription.html")
	tmpl.Execute(w, nil)
}

func (c *AuthentificationControleur) SeConnecter(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email") 
	motDePasse := r.FormValue("not_de_passe")
	u, err := c.RepoUser.VerifierConnexion(email, motDePasse)
	if err != nil {
		log.Println("Échec de connexion pour l'email :", email)
		data := struct{ Erreur string }{Erreur: err.Error()}
		tmpl, _ := template.ParseFiles("html/authentification/connexion.html")
		tmpl.Execute(w, data)
		return
	}
	token := SQL.GenererTokenSession(u)
	cookie := &http.Cookie{
		Name:  "session_token",
		Value: token,
		Path:  "/",
	}
	http.SetCookie(w, cookie)
		cookiePseudo := &http.Cookie{
		Name:  "pseudo",
		Value: u.Pseudo,
		Path:  "/",
	}
	http.SetCookie(w, cookiePseudo)
	log.Println("Session et cookie pseudo créés pour :", u.Pseudo)
	http.Redirect(w, r, "/forum/accueil", http.StatusSeeOther)
}

func (c *AuthentificationControleur) SInscrire(w http.ResponseWriter, r *http.Request) {
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
	_, err := c.RepoUser.CreerUtilisateur(nouvelUtilisateur, motDePasse)
	if err != nil {
		log.Println("Inscription refusée")
		data := struct{ Erreur string }{Erreur: err.Error()}
		tmpl, _ := template.ParseFiles("html/authentification/inscription.html")
		tmpl.Execute(w, data)
		return
	}
	log.Println("nouveau inscrit :", pseudo, email)
	http.Redirect(w, r, "/forum/accueil", http.StatusSeeOther)
}
