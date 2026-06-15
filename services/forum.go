package services

import "footconnect/SQL"

type ForumService struct {
	RepoDiscu *SQL.DiscussionRepository
	RepoComm  *SQL.CommentaireRepository
}

func InitialiserForumService(rd *SQL.DiscussionRepository, rc *SQL.CommentaireRepository) *ForumService {
	return &ForumService{RepoDiscu: rd, RepoComm: rc}
}

func (s *ForumService) FiltrerLesDiscussions(tri string, limite, offset int, recherche string, catID int) ([]SQL.DiscussionAffichee, error) {
	if limite != 10 && limite != 20 && limite != 30 {
		limite = 10
	}
	return s.RepoDiscu.RecupererDiscussions(tri, limite, offset, recherche, catID)
}
