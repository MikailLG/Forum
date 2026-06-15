package routeurs

import (
	"footconnect/controleurs"
	"net/http"
)

func ConfigurerRoutesForum(ctrl *controleurs.ForumControleur) {
	http.HandleFunc("/forum/accueil", ctrl.AfficherAccueil)
}
