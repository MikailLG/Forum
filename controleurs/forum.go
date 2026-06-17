package controleurs

import (
	"footconnect/SQL"
	"footconnect/services"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ForumControleur struct {
	Service *services.ForumService
}

type CommentaireAffiche struct {
	IDDiscussion int
	Auteur       string
	Contenu      string
	Date         string
}

var BoiteAuxLettresCommentaires []CommentaireAffiche

func InitialiserForumControleur(s *services.ForumService) *ForumControleur {
	return &ForumControleur{Service: s}
}

func (c *ForumControleur) AfficherAccueil(w http.ResponseWriter, r *http.Request) {
	action := r.URL.Query().Get("action")
	if action == "supprimer" {
		idStr := r.URL.Query().Get("id")
		idDiscussion, err := strconv.Atoi(strings.TrimSpace(idStr))
		if err == nil {
			// On exécute la suppression en cascade directement
			errDel := c.Service.RepoDiscu.SupprimerDiscussionEnCascade(idDiscussion)
			if errDel != nil {
				log.Println("Échec de la suppression en cascade depuis l'accueil :", errDel)
			} else {
				log.Println("Succès : Le sujet ID", idDiscussion, "a bien été supprimé via l'URL.")
				http.Redirect(w, r, "/forum/accueil", http.StatusSeeOther)
				return
			}
		}
	}
	recherche := r.URL.Query().Get("q")
	tri := r.URL.Query().Get("tri")
	if tri == "" {
		tri = "date"
	}
	limite := 10
	offset := 0
	catID := 0
	toutesLesDiscussions, err := c.Service.FiltrerLesDiscussions(tri, limite, offset, recherche, catID)
	if err != nil {
		log.Println("Erreur récupération discussions :", err)
	}
	estAdmin := false
	if cookie, err := r.Cookie("session_token"); err == nil {
		if strings.Contains(cookie.Value, "_") {
			parties := strings.Split(cookie.Value, "_")
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
		log.Println("Erreur chargement fichiers HTML de l'accueil :", errTemplate)
		http.Error(w, "Erreur interne (fichiers HTML manquants)", http.StatusInternalServerError)
		return
	}
	errExecute := tmpl.ExecuteTemplate(w, "base.html", data)
	if errExecute != nil {
		errExecute2 := tmpl.ExecuteTemplate(w, "base", data)
		if errExecute2 != nil {
			tmpl.Execute(w, data)
		}
	}
}

func (c *ForumControleur) SupprimerDiscussion(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	idDiscussion, err := strconv.Atoi(strings.TrimSpace(idStr))
	if err != nil {
		http.Error(w, "Identifiant invalide", http.StatusBadRequest)
		return
	}
	err = c.Service.RepoDiscu.SupprimerDiscussionEnCascade(idDiscussion)
	if err != nil {
		log.Println("Échec de la suppression en cascade :", err)
		http.Error(w, "Erreur lors de la suppression", http.StatusInternalServerError)
		return
	}
	log.Println("Succès : Le sujet ID", idDiscussion, "a bien été supprimé.")
	w.Header().Set("Content-Type", "text/plain")
	http.Redirect(w, r, "/forum/accueil", http.StatusSeeOther)
}

func (c *ForumControleur) AfficherFil(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	idDiscussion, err := strconv.Atoi(strings.TrimSpace(idStr))
	if err != nil {
		http.Error(w, "ID de discussion invalide", http.StatusBadRequest)
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
		discussionSelectionnee = SQL.DiscussionAffichee{
			ID:      idDiscussion,
			Titre:   "Sujet introuvable",
			Contenu: "Ce sujet n'existe pas ou a été supprimé.",
			Etat:    "fermé",
		}
	}
	var commentairesDuSujet []CommentaireAffiche
	for _, comm := range BoiteAuxLettresCommentaires {
		if comm.IDDiscussion == idDiscussion {
			commentairesDuSujet = append(commentairesDuSujet, comm)
		}
	}
	data := struct {
		Discussion   SQL.DiscussionAffichee
		Commentaires []CommentaireAffiche
	}{
		Discussion:   discussionSelectionnee,
		Commentaires: commentairesDuSujet,
	}
	tmpl, errTemplate := template.ParseFiles("html/forum/fil.html", "html/structures/base.html")
	if errTemplate != nil {
		log.Println("Erreur chargement HTML fil :", errTemplate)
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
		return
	}
	errExecute := tmpl.ExecuteTemplate(w, "base.html", data)
	if errExecute != nil {
		tmpl.Execute(w, data)
	}
}

func (c *ForumControleur) CreerDiscussion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	titre := r.FormValue("titre")
	contenu := r.FormValue("contenu")
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
		http.Error(w, "Erreur serveur lors de la création", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/forum/accueil", http.StatusSeeOther)
}

func (c *ForumControleur) Voter(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	idDiscussion, err := strconv.Atoi(idStr)
	
	typeVote := strings.ToLower(r.URL.Query().Get("type"))
	if err != nil || (typeVote != "like" && typeVote != "dislike") {
		http.Error(w, "Paramètres invalides", http.StatusBadRequest)
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
	contenu := r.FormValue("contenu")
	cookieSession, err := r.Cookie("session_token")
	if err != nil {
		log.Println("Utilisateur non connecté.")
		http.Redirect(w, r, "/connexion", http.StatusSeeOther)
		return
	}
	nomAuteur := ""
	if cookiePseudo, errP := r.Cookie("pseudo"); errP == nil && cookiePseudo.Value != "" {
		nomAuteur = cookiePseudo.Value
	} else {
		valeurCookie := cookieSession.Value
		if strings.Contains(valeurCookie, "_") {
			parties := strings.Split(valeurCookie, "_")
			if len(parties) >= 2 {
				if strings.ToUpper(parties[0]) == "TOKEN" {
					nomAuteur = parties[1] 
				} else {
					nomAuteur = parties[0]
				}
			}
		}
	}
	if nomAuteur == "" || nomAuteur == "Utilisateur" || nomAuteur == "TOKEN" || nomAuteur == "FOOTCONNECT" || nomAuteur == "Auteur_Discussion" {
		nomAuteur = "Supporter_FC" 
	}
	idDiscussion, _ := strconv.Atoi(idStr)
	if strings.TrimSpace(contenu) != "" {
		nouveauComm := CommentaireAffiche{
			IDDiscussion: idDiscussion,
			Auteur:       nomAuteur,
			Contenu:      contenu,
			Date:         "À l'instant",
		}
		BoiteAuxLettresCommentaires = append(BoiteAuxLettresCommentaires, nouveauComm)
		log.Printf("Commentaire enregistré par %s sur le sujet %d", nomAuteur, idDiscussion)
	}
	http.Redirect(w, r, "/forum/fil?id="+idStr, http.StatusSeeOther)
}
