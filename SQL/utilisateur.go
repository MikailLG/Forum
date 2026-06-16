package SQL

import (
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"unicode"
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

func ValiderMotDePasse(mdp string) error {
	if len(mdp) < 12 {
		return fmt.Errorf("le mot de passe doit contenir au moins 12 caractères")
	}
	aUneMajuscule := false
	aUnCarSpecial := false
	for _, caractere := range mdp {
		if unicode.IsUpper(caractere) {
			aUneMajuscule = true
		}
		if unicode.IsPunct(caractere) || unicode.IsSymbol(caractere) {
			aUnCarSpecial = true
		}
	}
	if !aUneMajuscule {
		return fmt.Errorf("le mot de passe doit contenir au moins une lettre majuscule")
	}
	if !aUnCarSpecial {
		return fmt.Errorf("le mot de passe doit contenir au moins un caractère spécial")
	}
	return nil
}

func HacherEnSHA512(mdpBrut string) string {
	algorithme := sha512.New()
	algorithme.Write([]byte(mdpBrut))
	octetsDuHash := algorithme.Sum(nil)
	return hex.EncodeToString(octetsDuHash)
}

func (r *UtilisateurRepository) CreerUtilisateur(u Utilisateur, mdpBrut string) (int, error) {
	if erreurValidation := ValiderMotDePasse(mdpBrut); erreurValidation != nil {
		return -1, erreurValidation
	}
	u.MotDePasseHache = HacherEnSHA512(mdpBrut)
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

func (r *UtilisateurRepository) VerifierConnexion(email string, mdpSoumis string) (*Utilisateur, error) {
    requete := "SELECT id_utilisateur, pseudo, email, mot_de_passe_hache, club_de_coeur, est_banni, id_role FROM UTILISATEURS WHERE email = ?;"
    ligne := r.db.QueryRow(requete, email)
    var u Utilisateur
    erreurScan := ligne.Scan(&u.ID, &u.Pseudo, &u.Email, &u.MotDePasseHache, &u.ClubDeCoeur, &u.EstBanni, &u.IDRole)
    if erreurScan != nil {
        if errors.Is(erreurScan, sql.ErrNoRows) {
            return nil, fmt.Errorf("email ou mot de passe incorrect")
        }
        return nil, fmt.Errorf("erreur lors de la requete de connexion : %w", erreurScan)
    }
    if u.EstBanni {
        return nil, fmt.Errorf("ce compte a été banni de la plateforme")
    }
    hashSoumis := HacherEnSHA512(mdpSoumis)
    if hashSoumis != u.MotDePasseHache {
        return nil, fmt.Errorf("email ou mot de passe incorrect")
    }
    return &u, nil
}

func GenererTokenSession(u *Utilisateur) string {
	tokenUnique := fmt.Sprintf("TOKEN_FOOTCONNECT_%d_%s_%d", u.ID, u.Pseudo, u.IDRole)
	return tokenUnique
}
