package routeurs

import (
	"footconnect/controleurs"
	"net/http"
)

func ConfigurerRoutesAuthentification(ctrl *controleurs.AuthentificationControleur) {
	http.HandleFunc("/authentification/connexion", ctrl.AfficherConnexion)
	http.HandleFunc("/authentification/inscription", ctrl.AfficherInscription)
	http.HandleFunc("/connexion", ctrl.SeConnecter)
	http.HandleFunc("/inscription", ctrl.SInscrire)
}
