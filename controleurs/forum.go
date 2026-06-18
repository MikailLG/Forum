package controleurs

import (
	"footconnect/SQL"
	"footconnect/services"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ForumControleur struct {
	Service *services.ForumService
}

func InitialiserForumControleur(s *services.ForumService) *ForumControleur {
	return &ForumControleur{Service: s}
}

func (c *ForumControleur) AfficherAccueil(w http.ResponseWriter, r *http.Request) {
	recherche := r.URL.Query().Get("q")
	tri := r.URL.Query().Get("tri")
	if tri == "" {
		tri = "date"
	}
	limite, offset, catID := 10, 0, 0
	toutesLesDiscussions, err := c.Service.FiltrerLesDiscussions(tri, limite, offset, recherche, catID)
	if err != nil {
		log.Println("Erreur récupération discussions :", err)
	}
	estAdmin := false
	if cookie, err := r.Cookie("session_token"); err == nil {
		parties := strings.Split(cookie.Value, "_")
		if len(parties) > 0 {
			idRole := parties[len(parties)-1]
			if idRole == "2" {
				estAdmin = true
			}
		}
	}
	data := struct {
		EstAdmin    bool
		Discussions []SQL.DiscussionAffichee
	}{
		EstAdmin:    estAdmin,
		Discussions: toutesLesDiscussions,
	}
	tmpl, errTemplate := template.ParseFiles("html/forum/accueil.html", "html/structures/base.html")
	if errTemplate != nil {
		log.Println("Erreur chargement templates accueil :", errTemplate)
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "base.html", data)
}

func (c *ForumControleur) SupprimerDiscussion(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	idDiscussion, err := strconv.Atoi(strings.TrimSpace(idStr))
	if err != nil {
		http.Error(w, "Identifiant invalide", http.StatusBadRequest)
		return
	}
	errCascade := c.Service.RepoDiscu.SupprimerDiscussionEnCascade(idDiscussion)
	if errCascade != nil {
		log.Println("Échec suppression ID", idDiscussion, ":", errCascade)
		http.Error(w, "Erreur lors de la suppression", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/forum/accueil", http.StatusSeeOther)
}

func (c *ForumControleur) AfficherFil(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	idDiscussion, err := strconv.Atoi(strings.TrimSpace(idStr))
	if err != nil {
		http.Error(w, "Identifiant invalide", http.StatusBadRequest)
		return
	}
	limite, offset, catID := 100, 0, 0
	discussions, err := c.Service.FiltrerLesDiscussions("date", limite, offset, "", catID)
	var discussionSelectionnee SQL.DiscussionAffichee
	trouve := false
	if err == nil {
		for _, d := range discussions {
			if d.ID == idDiscussion {
				discussionSelectionnee = d
				trouve = true
				break
			}
		}
	}
	if !trouve {
		http.Error(w, "Sujet introuvable", http.StatusNotFound)
		return
	}
	commentairesDuSujet, errComm := c.Service.ObtenirCommentairesDunSujet(idDiscussion)
	if errComm != nil {
		log.Println("Erreur récupération commentaires BDD :", errComm)
	}
	estAdmin := false
	if cookie, err := r.Cookie("session_token"); err == nil {
		parties := strings.Split(cookie.Value, "_")
		if len(parties) > 0 {
			idRole := parties[len(parties)-1]
			if idRole == "2" {
				estAdmin = true
			}
		}
	}
	data := struct {
		EstAdmin     bool
		Discussion   SQL.DiscussionAffichee
		Commentaires []SQL.CommentaireAffiche
	}{
		EstAdmin:     estAdmin,
		Discussion:   discussionSelectionnee,
		Commentaires: commentairesDuSujet,
	}
	tmpl, errTemplate := template.ParseFiles("html/forum/fil.html", "html/structures/base.html")
	if errTemplate != nil {
		log.Println("Erreur chargement templates fil :", errTemplate)
		http.Error(w, "Fichier HTML introuvable", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "base.html", data)
}

func (c *ForumControleur) CreerDiscussion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	titre := strings.TrimSpace(r.FormValue("titre"))
	contenu := strings.TrimSpace(r.FormValue("contenu"))
	idUtilisateur := 1
	if cookie, err := r.Cookie("session_token"); err == nil {
		parties := strings.Split(cookie.Value, "_")
		if len(parties) >= 3 {
			if id, errConv := strconv.Atoi(parties[2]); errConv == nil {
				idUtilisateur = id
			}
		}
	}
	nouvelleDiscu := SQL.Discussion{
		Titre:         titre,
		Contenu:       contenu,
		IDUtilisateur: idUtilisateur,
		IDCategorie:   1,
	}
	_, err := c.Service.RepoDiscu.CreerDiscussion(nouvelleDiscu)
	if err != nil {
		log.Println("Erreur création discussion :", err)
		http.Error(w, "Erreur lors de la publication", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/forum/accueil", http.StatusSeeOther)
}

func (c *ForumControleur) Voter(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	idDiscussion, err := strconv.Atoi(idStr)
	typeVote := strings.ToLower(r.URL.Query().Get("type"))
	if err != nil || (typeVote != "like" && typeVote != "dislike") {
		http.Error(w, "Paramètres incorrects", http.StatusBadRequest)
		return
	}
	idUtilisateur := 1
	if cookie, err := r.Cookie("session_token"); err == nil {
		parties := strings.Split(cookie.Value, "_")
		if len(parties) >= 3 {
			if id, errConv := strconv.Atoi(parties[2]); errConv == nil {
				idUtilisateur = id
			}
		}
	}
	_ = c.Service.RepoDiscu.Voter(idUtilisateur, idDiscussion, typeVote)
	http.Redirect(w, r, "/forum/accueil", http.StatusSeeOther)
}

func (c *ForumControleur) CreerCommentaire(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.FormValue("id_discussion")
	contenu := strings.TrimSpace(r.FormValue("contenu"))
	idDiscussion, _ := strconv.Atoi(idStr)
	idUtilisateur := 1
	if cookieSession, err := r.Cookie("session_token"); err == nil {
		parties := strings.Split(cookieSession.Value, "_")
		if len(parties) >= 3 {
			if id, errConv := strconv.Atoi(parties[2]); errConv == nil {
				idUtilisateur = id
			}
		}
	}
	if contenu != "" && idDiscussion > 0 {
		nouveauComm := SQL.Commentaire{
			Contenu:       contenu,
			IDDiscussion: idDiscussion,
			IDUtilisateur: idUtilisateur,
			DateCreation:  time.Now().Format("2006-01-02 15:04:05"),
		}
		errInsert := c.Service.RepoComm.CreerCommentaire(nouveauComm)
		if errInsert != nil {
			log.Println("Erreur insertion commentaire BDD :", errInsert)
		}
	}
	if idStr == "" {
		idStr = strconv.Itoa(idDiscussion)
	}
	http.Redirect(w, r, "/forum/fil?id="+idStr, http.StatusSeeOther)
}
