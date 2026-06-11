package SQL

import (
	"database/sql"
	"fmt"
	"time"
)

type Commentaire struct {
	ID             int
	Contenu        string
	DateCreation   time.Time
	IDUtilisateur  int
	IDDiscussion   int
}

type CommentaireRepository struct {
	db *sql.DB
}

func InitialiserCommentaireRepository(db *sql.DB) *CommentaireRepository {
	return &CommentaireRepository{db: db}
}

func (r *CommentaireRepository) CreerCommentaire(c Commentaire) (int, error) {
	requete := "INSERT INTO COMMENTAIRES (contenu, id_utilisateur, id_discussion) VALUES (?, ?, ?);"
	resultatSQL, erreurSQL := r.db.Exec(requete, c.Contenu, c.IDUtilisateur, c.IDDiscussion)
	if erreurSQL != nil {
		return -1, fmt.Errorf("erreur lors de l'ajout du commentaire : %w", erreurSQL)
	}
	idGenere, erreurID := resultatSQL.LastInsertId()
	if erreurID != nil {
		return -1, fmt.Errorf("erreur recuperation ID commentaire : %w", erreurID)
	}
	return int(idGenere), nil
}
