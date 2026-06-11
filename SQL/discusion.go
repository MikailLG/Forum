package SQL

import (
	"database/sql"
	"fmt"
	"time"
)

type Discussion struct {
	ID             int
	Titre          string
	Contenu        string
	DateCreation   time.Time
	Etat           string
	IDUtilisateur  int
	IDCategorie    int
}

type DiscussionRepository struct {
	db *sql.DB
}

func InitialiserDiscussionRepository(db *sql.DB) *DiscussionRepository {
	return &DiscussionRepository{db: db}
}

func (r *DiscussionRepository) CreerDiscussion(d Discussion) (int, error) {
	if d.Etat == "" {
		d.Etat = "ouvert"
	}
	requete := "INSERT INTO DISCUSSIONS (titre, contenu, etat, id_utilisateur, id_categorie) VALUES (?, ?, ?, ?, ?);"
	resultatSQL, erreurSQL := r.db.Exec(requete, d.Titre, d.Contenu, d.Etat, d.IDUtilisateur, d.IDCategorie)
	if erreurSQL != nil {
		return -1, fmt.Errorf("erreur lors de la creation du sujet : %w", erreurSQL)
	}
	idGenere, erreurID := resultatSQL.LastInsertId()
	if erreurID != nil {
		return -1, fmt.Errorf("erreur recuperation ID sujet : %w", erreurID)
	}
	return int(idGenere), nil
}

func (r *DiscussionRepository) ModifierEtat(idDiscussion int, nouvelEtat string) error {
	if nouvelEtat != "ouvert" && nouvelEtat != "fermé" && nouvelEtat != "archivé" {
		return fmt.Errorf("etat invalide : %s", nouvelEtat)
	}
	requete := "UPDATE DISCUSSIONS SET etat = ? WHERE id_discussion = ?;"
	_, err := r.db.Exec(requete, nouvelEtat, idDiscussion)
	if err != nil {
		return fmt.Errorf("erreur modification etat sujet : %w", err)
	}
	return nil
}
