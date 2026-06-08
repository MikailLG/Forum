CREATE DATABASE IF NOT EXISTS footconnect_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE footconnect_db;

CREATE TABLE IF NOT EXISTS ROLES (
    id_role INT AUTO_INCREMENT PRIMARY KEY,
    libelle VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS UTILISATEURS (
    id_utilisateur INT AUTO_INCREMENT PRIMARY KEY,
    pseudo VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    mot_de_passe_hache VARCHAR(255) NOT NULL,
    club_de_coeur VARCHAR(100),
    est_banni BOOLEAN DEFAULT FALSE,
    id_role INT NOT NULL,
    FOREIGN KEY (id_role) REFERENCES ROLES(id_role)
);

CREATE TABLE IF NOT EXISTS CATEGORIES (
    id_categorie INT AUTO_INCREMENT PRIMARY KEY,
    nom_categorie VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS DISCUSSIONS (
    id_discussion INT AUTO_INCREMENT PRIMARY KEY,
    titre VARCHAR(150) NOT NULL,
    contenu TEXT NOT NULL,
    statut ENUM('ouvert', 'ferme', 'archive') DEFAULT 'ouvert',
    date_creation TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    id_utilisateur INT NOT NULL,
    id_categorie INT NOT NULL,
    FOREIGN KEY (id_utilisateur) REFERENCES UTILISATEURS(id_utilisateur) ON DELETE CASCADE,
    FOREIGN KEY (id_categorie) REFERENCES CATEGORIES(id_categorie)
);

CREATE TABLE IF NOT EXISTS MESSAGES (
    id_message INT AUTO_INCREMENT PRIMARY KEY,
    contenu TEXT NOT NULL,
    date_creation TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    id_utilisateur INT NOT NULL,
    id_discussion INT NOT NULL,
    FOREIGN KEY (id_utilisateur) REFERENCES UTILISATEURS(id_utilisateur) ON DELETE CASCADE,
    FOREIGN KEY (id_discussion) REFERENCES DISCUSSIONS(id_discussion) ON DELETE CASCADE 
);

CREATE TABLE IF NOT EXISTS REACTIONS (
    id_reaction INT AUTO_INCREMENT PRIMARY KEY,
    type_reaction ENUM('like', 'dislike') NOT NULL, 
    id_utilisateur INT NOT NULL,
    id_message INT NOT NULL,
    UNIQUE KEY unique_reaction_utilisateur (id_utilisateur, id_message), 
    FOREIGN KEY (id_utilisateur) REFERENCES UTILISATEURS(id_utilisateur) ON DELETE CASCADE,
    FOREIGN KEY (id_message) REFERENCES MESSAGES(id_message) ON DELETE CASCADE 
);

INSERT INTO ROLES (libelle) VALUES ('Supporter'), ('Administrateur');
INSERT INTO CATEGORIES (nom_categorie) VALUES ('Ligue 1'), ('Ligue des Champions'), ('Mercato'), ('Coupe du Monde');
