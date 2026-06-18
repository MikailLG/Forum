package SQL

import (
	"database/sql"
	"fmt"
)

type Commentaire struct {
	ID            int
	Contenu       string
	DateCreation  string
	IDUtilisateur int
	IDDiscussion  int
}

type CommentaireAffiche struct {
	IDDiscussion int
	Auteur       string
	Contenu      string
	Date         string
}

type CommentaireRepository struct {
	db *sql.DB
}

func InitialiserCommentaireRepository(db *sql.DB) *CommentaireRepository {
	return &CommentaireRepository{db: db}
}

func (r *CommentaireRepository) CreerCommentaire(c Commentaire) error {
	requete := "INSERT INTO MESSAGES (contenu, id_utilisateur, id_discussion) VALUES (?, ?, ?);"
	_, err := r.db.Exec(requete, c.Contenu, c.IDUtilisateur, c.IDDiscussion)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ajout du commentaire : %w", err)
	}
	return nil
}

func (r *CommentaireRepository) RecupererCommentairesParSujet(idDiscussion int) ([]CommentaireAffiche, error) {
	requete := `
		SELECT m.id_discussion, u.pseudo, m.contenu, m.date_creation 
		FROM MESSAGES m
		JOIN UTILISATEURS u ON m.id_utilisateur = u.id_utilisateur
		WHERE m.id_discussion = ?
		ORDER BY m.date_creation ASC;`
	lignes, err := r.db.Query(requete, idDiscussion)
	if err != nil {
		return nil, fmt.Errorf("erreur de lecture commentaires : %w", err)
	}
	defer lignes.Close()
	var commentaires []CommentaireAffiche
	for lignes.Next() {
		var comm CommentaireAffiche
		errScan := lignes.Scan(&comm.IDDiscussion, &comm.Auteur, &comm.Contenu, &comm.Date)
		if errScan != nil {
			return nil, errScan
		}
		commentaires = append(commentaires, comm)
	}
	return commentaires, nil
}
