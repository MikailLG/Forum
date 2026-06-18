package routeurs

import (
	"footconnect/controleurs"
	"net/http"
)

func ConfigurerRoutesAuthentification(ctrl *controleurs.AuthentificationControleur) {
	http.HandleFunc("/authentification/connexion", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			ctrl.SeConnecter(w, r)
		} else {
			ctrl.AfficherConnexion(w, r)
		}
	})
	http.HandleFunc("/authentification/inscription", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			ctrl.SInscrire(w, r)
		} else {
			ctrl.AfficherInscription(w, r)
		}
	})
}
