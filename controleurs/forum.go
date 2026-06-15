package controleurs

import (
	"footconnect/SQL"
	"footconnect/services"
	"html/template"
	"net/http"
	"strconv"
)

type ForumControleur struct {
	Service *services.ForumService
}

func InitialiserForumControleur(s *services.ForumService) *ForumControleur {
	return &ForumControleur{Service: s}
}

func (c *ForumControleur) AfficherAccueil(w http.ResponseWriter, r *http.Request) {
	tri := r.URL.Query().Get("tri")
	recherche := r.URL.Query().Get("q")
	idCategorie, _ := strconv.Atoi(r.URL.Query().Get("cat"))
	limite, _ := strconv.Atoi(r.URL.Query().Get("limite"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limite
	discussions, err := c.Service.FiltrerLesDiscussions(tri, limite, offset, recherche, idCategorie)
	if err != nil {
		http.Error(w, "Erreur de chargement des données", http.StatusInternalServerError)
		return
	}
	data := struct {
		Discussions []SQL.DiscussionAffichee
	}{
		Discussions: discussions,
	}
	tmpl, _ := template.ParseFiles("html/structures/base.html", "html/forum/accueil.html")
	tmpl.ExecuteTemplate(w, "base.html", data)
}
