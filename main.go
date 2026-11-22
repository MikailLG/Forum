package main

import (
	"footconnect/SQL"
	"footconnect/controleurs"
	"footconnect/database"
	"footconnect/routeurs"
	"footconnect/services"
	"log"
	"net/http"
	"os"
	"time"
)

type logResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *logResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		debut := time.Now()
		lrw := &logResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lrw, r)
		log.Printf("%s %s %d %v", r.Method, r.URL.Path, lrw.statusCode, time.Since(debut))
	})
}

func main() {
	contexteBDD := database.InitialiserBDD()
	defer contexteBDD.Close()
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/web/", http.StripPrefix("/web/", fs))
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	repoDiscu := SQL.InitialiserDiscussionRepository(contexteBDD)
	repoComm := SQL.InitialiserCommentaireRepository(contexteBDD)
	serviceForum := services.InitialiserForumService(repoDiscu, repoComm)
	ctrlForum := controleurs.InitialiserForumControleur(serviceForum)
	routeurs.ConfigurerRoutesForum(ctrlForum)
	repoUser := SQL.InitialiserUtilisateurRepository(contexteBDD)
	ctrlAuth := controleurs.InitialiserAuthentificationControleur(repoUser)
	routeurs.ConfigurerRoutesAuthentification(ctrlAuth)
	ctrlAdmin := controleurs.InitialiserAdministrationControleur(repoUser)
	http.HandleFunc("/admin/dashboard", ctrlAdmin.AfficherTableauDeBord)
	http.HandleFunc("/admin/bannir", ctrlAdmin.BannirUtilisateur)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("FootConnect est en ligne sur le port %s", port)
	if err := http.ListenAndServe(":"+port, loggerMiddleware(http.DefaultServeMux)); err != nil {
		log.Fatal("Le serveur s'est arrêté brutalement : ", err)
	}
}
