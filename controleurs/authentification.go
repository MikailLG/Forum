package controleurs

import (
	"html/template"
	"net/http"
)

type AuthentificationControleur struct{}

func InitialiserAuthentificationControleur() *AuthentificationControleur {
	return &AuthentificationControleur{}
}

func (c *AuthentificationControleur) AfficherConnexion(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("html/authentification/connexion.html")
	tmpl.Execute(w, nil)
}

func (c *AuthentificationControleur) AfficherInscription(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("html/authentification/inscription.html")
	tmpl.Execute(w, nil)
}
