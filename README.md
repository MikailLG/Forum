# FootConnect - Plateforme Communautaire de Football

Bienvenue sur **FootConnect**, une plateforme web conçue pour rassembler tous les passionnés du ballon rond. Développée dans le cadre de notre projet de fin de première année, cette application permet aux supporters de débattre des matchs, de suivre les actualités du mercato, et de consulter les derniers scores en direct.


## Fonctionnalités Principales

### Fonctionnalités Visiteurs (Utilisateur non connecté)
* **Lecture des Fils de Discussion :** Accès libre aux débats de la communauté (Ligue 1, Ligue des Champions, Mercato...).
* **Consultation des Commentaires :** Lecture des avis et réactions des autres supporters.
* **Authentification :** Création de compte et connexion sécurisée.
* **Live Scores (Bonus) :** Consultation des scores et résultats des derniers matchs en temps réel.

### Fonctionnalités Utilisateurs (Supporters connectés)
* **Création de Sujets :** Possibilité de lancer un nouveau débat ou fil de discussion.
* **Participation :** Ajout de messages et commentaires sous les sujets existants.
* **Interactions :** Système de réactions (Upvote/Downvote ou émojis) pour évaluer la pertinence d'un message.
* **Profil Personnalisé :** Gestion des informations personnelles et sélection de son club de cœur.

## Technologie

* **Backend :** Go (Golang)
* **Frontend :** HTML / CSS / JavaScript
* **Base de données :** MySQL
* **API Live Scores :** Intégration d'une API externe pour récupérer les scores en direct.

## Structure Prévisionnelle de la Base de Données

L'application repose sur un modèle relationnel comprenant les entités suivantes :
* `User` : ID, Pseudo, Email, Mot de passe (haché avec `bcrypt`), Club de cœur, ID_Rôle.
* `Post` : ID, Titre, Contenu, Date de création, ID_Utilisateur, ID_Catégorie.
* `Comment` : ID, Contenu, Date de création, ID_Utilisateur, ID_Post.
* `Category` : ID, Nom (ex: Ligue 1, Mercato, International).
* `Reaction` : ID, Type (Like/Dislike), ID_Utilisateur, ID_Commentaire.

## Gestion de Projet & Organisation

L'organisation du travail de notre binôme est entièrement structurée et suivie à l'aide d'un tableau Trello. 
* **Lien vers notre Trello :** https://trello.com/invite/b/6a100bebf894fb2031ae050d/ATTIe8fe819ad9587f41a9f7bd12af7d9b174D73F01E/mon-tableau-trello


## 👥 Notre Binôme

* ** Liazid Gracia Mikaïl :** 
* ** Silue Mitaumikan Emmanuel:** 

## Installation et Lancement

(Prochainement)