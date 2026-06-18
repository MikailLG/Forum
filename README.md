# FootConnect - Plateforme Communautaire de Football

Bienvenue sur **FootConnect**, une plateforme web conçue pour rassembler tous les passionnés du ballon. Développée dans le cadre de notre projet de fin de première année, cette application offre un espace d'échange complet aux supporters pour débattre sur l'actualité du football.

## Notre Binôme
* **Liazid Gracia Mikaïl**
* **Silue Mitaumikan Emmanuel**

---

## Fonctionnalités Principales

### Visiteurs (Utilisateur non connecté)
* **Lecture des Fils de Discussion :** Accès libre aux débats de la communauté classés par thématiques (Ligue 1, Ligue des Champions, Mercato...).
* **Consultation des Commentaires :** Lecture en direct des avis et réactions des autres supporters.
* **Authentification :** Accès aux formulaires de création de compte et de connexion sécurisée.

### Supporters (Utilisateurs connectés)
* **Création de Sujets :** Possibilité de lancer un nouveau débat ou fil de discussion sur le forum.
* **Participation active :** Ajout de messages et commentaires sous les sujets existants.
* **Système de Réactions :** Possibilité de voter (Like/Dislike) sur les sujets pour évaluer leur pertinence et modifier leur score de popularité.
* **Profil Personnalisé :** Sélection et affichage de son club de cœur lors de l'inscription.

### Espace Modération (Staff)
* **Tableau de Bord Admin :** Espace sécurisé (accessible uniquement via le rôle Administrateur).
* **Modération stricte :** Possibilité de bannir un utilisateur et de nettoyer le forum avec une suppression en cascade automatique (sujets, messages et votes associés).

---

## Architecture & Technologies

Pour ce projet, nous avons fait le choix d'un développement robuste sans frameworks lourds externes, en s'appuyant sur la puissance de la bibliothèque standard de Go :

* **Backend :** Go (Golang) natif (Routage HTTP, moteur de templates, validation de formulaires).
* **Frontend :** HTML / CSS.
* **Base de données :** MySQL.
* **Sécurité :** Validation stricte des mots de passe à l'inscription et hachage de sécurité.

---

## Structure de la Base de Données

L'application repose sur un modèle relationnel strict comprenant les entités suivantes :
* `ROLES` : id_role, libelle.
* `UTILISATEURS` : id_utilisateur, pseudo, email, mot_de_passe_hache, club_de_coeur, est_banni, id_role.
* `CATEGORIES` : id_categorie, nom_categorie.
* `DISCUSSIONS` : id_discussion, titre, contenu, statut (ouvert/fermé/archivé), date_creation, id_utilisateur, id_categorie.
* `MESSAGES` : id_message, contenu, date_creation, id_utilisateur, id_discussion.
* `REACTIONS` : id_reaction, type_reaction (like/dislike), id_utilisateur, id_message.

---

## Gestion de Projet & Organisation

L'organisation du travail de notre binôme a été entièrement structurée, découpée et suivie à l'aide d'un tableau Trello. 
* **Lien vers notre Trello :** [Consulter notre tableau Trello](https://trello.com/invite/b/6a100bebf894fb2031ae050d/ATTIe8fe819ad9587f41a9f7bd12af7d9b174D73F01E/mon-tableau-trello)

---

## Lancement

### 1. Configuration de la Base de Données
1. go run main.go
2. http://localhost:8080/forum/accueil
