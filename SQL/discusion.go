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

func (r *DiscussionRepository) Voter(idUtilisateur int, idDiscussion int, typeReaction string) error {
    if typeReaction != "like" && typeReaction != "dislike" {
        return fmt.Errorf("type de réaction invalide")
    }
    requete := `INSERT INTO REACTIONS (id_utilisateur, id_message, type_reaction) 
                VALUES (?, ?, ?) 
                ON DUPLICATE KEY UPDATE type_reaction = ?;`
    _, err := r.db.Exec(requete, idUtilisateur, idDiscussion, typeReaction, typeReaction)
    if err != nil {
        return fmt.Errorf("erreur lors de l'enregistrement du vote : %w", err)
    }
    return nil
}

func (r *DiscussionRepository) RecupererDiscussions(triCritere string) ([]DiscussionAffichee, error) {
    var requete string
    if triCritere == "popularite" {
        requete = `SELECT d.id_discussion, d.titre, d.contenu, d.statut, d.id_utilisateur, d.id_categorie,
                   (COUNT(CASE WHEN re.type_reaction = 'like' THEN 1 END) - COUNT(CASE WHEN re.type_reaction = 'dislike' THEN 1 END)) AS score
                   FROM DISCUSSIONS d
                   LEFT JOIN REACTIONS re ON d.id_discussion = re.id_message
                   WHERE d.statut != 'archivé'
                   GROUP BY d.id_discussion
                   ORDER BY score DESC;`
    } else {
        requete = `SELECT d.id_discussion, d.titre, d.contenu, d.statut, d.id_utilisateur, d.id_categorie,
                   (COUNT(CASE WHEN re.type_reaction = 'like' THEN 1 END) - COUNT(CASE WHEN re.type_reaction = 'dislike' THEN 1 END)) AS score
                   FROM DISCUSSIONS d
                   LEFT JOIN REACTIONS re ON d.id_discussion = re.id_message
                   WHERE d.statut != 'archivé'
                   GROUP BY d.id_discussion
                   ORDER BY d.id_discussion DESC;`
    }
    lignes, err := r.db.Query(requete)
    if err != nil {
        return nil, fmt.Errorf("erreur recuperation des discussions : %w", err)
    }
    defer lignes.Close()
    var liste []DiscussionAffichee
    for lignes.Next() {
        var d DiscussionAffichee
        errScan := lignes.Scan(&d.ID, &d.Titre, &d.Contenu, &d.Etat, &d.IDUtilisateur, &d.IDCategorie, &d.ScorePopularite)
        if errScan != nil {
			return nil, fmt.Errorf("erreur scan discussion : %w", errScan)
		}
		liste = append(liste, d)
	}
	return liste, nil
}
