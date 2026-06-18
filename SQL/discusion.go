package SQL

import (
	"database/sql"
	"fmt"
)

type Discussion struct {
	ID            int
	Titre         string
	Contenu       string
	DateCreation  string
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

func InitialiserDiscussionRepository(db *sql.DB) *DiscussionRepository {
	return &DiscussionRepository{db: db}
}

func (r *DiscussionRepository) CreerDiscussion(d Discussion) (int, error) {
	if d.Etat == "" {
		d.Etat = "ouvert"
	}
	requete := "INSERT INTO DISCUSSIONS (titre, contenu, statut, id_utilisateur, id_categorie) VALUES (?, ?, ?, ?, ?);"
	resultatSQL, err := r.db.Exec(requete, d.Titre, d.Contenu, d.Etat, d.IDUtilisateur, d.IDCategorie)
	if err != nil {
		return -1, fmt.Errorf("erreur lors de la creation du sujet : %w", err)
	}
	idGenere, _ := resultatSQL.LastInsertId()
	return int(idGenere), nil
}

func (r *DiscussionRepository) ModifierEtat(idDiscussion int, nouvelEtat string) error {
	if nouvelEtat != "ouvert" && nouvelEtat != "fermé" && nouvelEtat != "archivé" {
		return fmt.Errorf("etat invalide : %s", nouvelEtat)
	}
	requete := "UPDATE DISCUSSIONS SET statut = ? WHERE id_discussion = ?;"
	_, err := r.db.Exec(requete, nouvelEtat, idDiscussion)
	return err
}

func (r *DiscussionRepository) Voter(idUtilisateur int, idDiscussion int, typeVote string) error {
	requete := `
		INSERT INTO REACTIONS (id_utilisateur, id_message, type_reaction) 
		VALUES (?, ?, ?) 
		ON DUPLICATE KEY UPDATE type_reaction = ?;`
	_, err := r.db.Exec(requete, idUtilisateur, idDiscussion, typeVote, typeVote)
	return err
}

func (r *DiscussionRepository) RecupererDiscussions(triCritere string, limite int, offset int, recherche string, idCategorie int) ([]DiscussionAffichee, error) {
	requete := `SELECT d.id_discussion, d.titre, d.contenu, d.statut, d.id_utilisateur, d.id_categorie,
               COUNT(CASE WHEN re.type_reaction = 'like' THEN 1 END) AS score
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
	_, _ = r.db.Exec("DELETE FROM REACTIONS WHERE id_message = ?;", idDiscussion)
	_, _ = r.db.Exec("DELETE FROM MESSAGES WHERE id_discussion = ?;", idDiscussion)	
	_, errFatal := r.db.Exec("DELETE FROM DISCUSSIONS WHERE id_discussion = ?;", idDiscussion)
	return errFatal
}
