package SQL

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
)

type Utilisateur struct {
	ID              int
	Pseudo          string
	Email           string
	MotDePasseHache string
	ClubDeCoeur     string
	EstBanni        bool
	IDRole          int
}

type UtilisateurRepository struct {
	db *sql.DB
}

func InitialiserUtilisateurRepository(db *sql.DB) *UtilisateurRepository {
	return &UtilisateurRepository{db: db}
}

func (r *UtilisateurRepository) CreerUtilisateur(u Utilisateur) (int, error) {
	requete := "INSERT INTO UTILISATEURS (pseudo, email, mot_de_passe_hache, club_de_coeur, id_role) VALUES (?, ?, ?, ?, ?);"
	resultatSQL, erreurSQL := r.db.Exec(requete, u.Pseudo, u.Email, u.MotDePasseHache, u.ClubDeCoeur, u.IDRole)
	if erreurSQL != nil {
		var erreurMySQL *mysql.MySQLError
		if errors.As(erreurSQL, &erreurMySQL) {
			if erreurMySQL.Number == 1062 {
				return -1, fmt.Errorf("le pseudo ou l'adresse email est déjà utilisé")
			}
		}
		return -1, fmt.Errorf("erreur lors de l'inscription : %w", erreurSQL)
	}
	idGenere, erreurID := resultatSQL.LastInsertId()
	if erreurID != nil {
		return -1, fmt.Errorf("erreur récupération de l'ID utilisateur : %w", erreurID)
	}

	return int(idGenere), nil
}
