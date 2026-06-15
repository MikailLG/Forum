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

func (c *ForumControleur) CreerDiscussion(w http.ResponseWriter, r *http.Request) {
	titre := r.FormValue("titre")
	contenu := r.FormValue("contenu")
	nouvelle := SQL.Discussion{
		Titre:         titre,
		Contenu:       contenu,
		IDUtilisateur: 1,
		IDCategorie:   1,
	}
	c.Service.RepoDiscu.CreerDiscussion(nouvelle)
	http.Redirect(w, r, "/forum/accueil", http.StatusSeeOther)
}

func (c *ForumControleur) AfficherFil(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	idDiscussion, _ := strconv.Atoi(idStr)
	toutes, _ := c.Service.RepoDiscu.RecupererDiscussions("chrono", 100, 0, "", 0)
	var discussionSelectionnee SQL.DiscussionAffichee
	for _, d := range toutes {
		if d.ID == idDiscussion {
			discussionSelectionnee = d
			break
		}
	}
	data := struct {
		Discussion   SQL.DiscussionAffichee
		Commentaires []SQL.Commentaire
	}{
		Discussion:   discussionSelectionnee,
		Commentaires: []SQL.Commentaire{},
	}
	tmpl, _ := template.ParseFiles("html/structures/base.html", "html/forum/fil.html")
	tmpl.ExecuteTemplate(w, "base.html", data)
}

func (c *ForumControleur) Voter(w http.ResponseWriter, r *http.Request) {
	idMessage, _ := strconv.Atoi(r.FormValue("id_message"))
	typeReaction := r.FormValue("type_reaction")
	c.Service.RepoDiscu.Voter(1, idMessage, typeReaction)
	http.Redirect(w, r, "/forum/accueil", http.StatusSeeOther)
}

func (c *ForumControleur) CreerCommentaire(w http.ResponseWriter, r *http.Request) {
	idDiscussionStr := r.FormValue("id_discussion")
	idDiscussion, _ := strconv.Atoi(idDiscussionStr)
	contenu := r.FormValue("contenu")
	nouveauComm := SQL.Commentaire{
		Contenu:       contenu,
		IDUtilisateur: 1,
		IDDiscussion:  idDiscussion,
	}
	c.Service.RepoComm.CreerCommentaire(nouveauComm)
	http.Redirect(w, r, "/forum/fil?id="+idDiscussionStr, http.StatusSeeOther)
}
