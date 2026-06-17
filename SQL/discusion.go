package SQL

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Discussion struct {
	ID            int
	Titre         string
	Contenu       string
	DateCreation  time.Time
	Etat          string
	IDUtilisateur int
	IDCategorie   int
}

type DiscussionRepository struct {
	db *sql.DB
}

type DiscussionAffichee struct {
	ID              int
	Titre           string
	Contenu         string
	Etat            string
	IDUtilisateur   int
	IDCategorie     int
	ScorePopularite int
}

type Reaction struct {
	IDUtilisateur int
	IDDiscussion  int
	TypeReaction  string
}

func InitialiserDiscussionRepository(db *sql.DB) *DiscussionRepository {
	return &DiscussionRepository{db: db}
}

func (r *DiscussionRepository) CreerDiscussion(d Discussion) (int, error) {
	if d.Etat == "" {
		d.Etat = "ouvert"
	}
	requete := "INSERT INTO DISCUSSIONS (titre, contenu, statut, id_utilisateur, id_categorie) VALUES (?, ?, ?, ?, ?);"
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
	requete := "UPDATE DISCUSSIONS SET statut = ? WHERE id_discussion = ?;"
	_, err := r.db.Exec(requete, nouvelEtat, idDiscussion)
	if err != nil {
		return fmt.Errorf("erreur modification etat sujet : %w", err)
	}
	return nil
}

func (r *DiscussionRepository) Voter(idUtilisateur int, idDiscussion int, typeVote string) error {
	_, _ = r.db.Exec("SET FOREIGN_KEY_CHECKS = 0;")
	requete := `
		INSERT INTO REACTIONS (id_utilisateur, id_message, type_reaction) 
		VALUES (?, ?, ?) 
		ON DUPLICATE KEY UPDATE type_reaction = ?;`
	
	_, err := r.db.Exec(requete, idUtilisateur, idDiscussion, typeVote, typeVote)
	_, _ = r.db.Exec("SET FOREIGN_KEY_CHECKS = 1;")

	if err != nil {
		log.Println("Erreur lors de l'enregistrement du vote :", err)
		return err
	}
	return nil
}

func (r *DiscussionRepository) RecupererDiscussions(triCritere string, limite int, offset int, recherche string, idCategorie int) ([]DiscussionAffichee, error) {
	requete := `SELECT d.id_discussion, d.titre, d.contenu, d.statut, d.id_utilisateur, d.id_categorie,
               GREATEST((COUNT(CASE WHEN re.type_reaction = 'like' THEN 1 END) - COUNT(CASE WHEN re.type_reaction = 'dislike' THEN 1 END)), 0) AS score
               FROM DISCUSSIONS d
               LEFT JOIN REACTIONS re ON d.id_discussion = re.id_message
               WHERE d.statut != 'archivé'`
	var args []interface{}
	if recherche != "" {
		requete += " AND d.titre LIKE ?"
		args = append(args, "%"+recherche+"%")
	}
	if idCategorie > 0 {
		requete += " AND d.id_categorie = ?"
		args = append(args, idCategorie)
	}
	requete += " GROUP BY d.id_discussion"
	if triCritere == "popularite" {
		requete += " ORDER BY score DESC"
	} else {
		requete += " ORDER BY d.id_discussion DESC"
	}
	requete += " LIMIT ? OFFSET ?;"
	args = append(args, limite, offset)
	lignes, err := r.db.Query(requete, args...)
	if err != nil {
		return nil, fmt.Errorf("erreur filtrage des discussions : %w", err)
	}
	defer lignes.Close()
	var liste []DiscussionAffichee
	for lignes.Next() {
		var d DiscussionAffichee
		errScan := lignes.Scan(&d.ID, &d.Titre, &d.Contenu, &d.Etat, &d.IDUtilisateur, &d.IDCategorie, &d.ScorePopularite)
		if errScan != nil {
			return nil, fmt.Errorf("erreur scan : %w", errScan)
		}
		liste = append(liste, d)
	}
	return liste, nil
}

func (r *DiscussionRepository) SupprimerDiscussionEnCascade(idDiscussion int) error {
	_, _ = r.db.Exec("SET FOREIGN_KEY_CHECKS = 0;")
	_, _ = r.db.Exec("DELETE FROM REACTIONS WHERE id_message = ?;", idDiscussion)
	_, _ = r.db.Exec("DELETE FROM COMMENTAIRES WHERE id_discussion = ?;", idDiscussion)
	_, _ = r.db.Exec("DELETE FROM COMMENTAIRE WHERE id_discussion = ?;", idDiscussion)
	_, _ = r.db.Exec("DELETE FROM MESSAGES WHERE id_discussion = ?;", idDiscussion)
	_, errFatal := r.db.Exec("DELETE FROM DISCUSSIONS WHERE id_discussion = ?;", idDiscussion)
	_, _ = r.db.Exec("SET FOREIGN_KEY_CHECKS = 1;")
	return errFatal
}
