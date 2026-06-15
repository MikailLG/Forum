package routeurs

import (
	"footconnect/controleurs"
	"net/http"
)

func ConfigurerRoutesForum(ctrl *controleurs.ForumControleur) {
	http.HandleFunc("/forum/accueil", ctrl.AfficherAccueil)
	http.HandleFunc("/forum/fil", ctrl.AfficherFil)
	http.HandleFunc("/forum/discussion/creer", ctrl.CreerDiscussion)
	http.HandleFunc("/forum/vote", ctrl.Voter)
	http.HandleFunc("/forum/commentaire/creer", ctrl.CreerCommentaire)
}
