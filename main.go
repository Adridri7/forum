package main

import (
	"fmt"
	"forum/server/api/categories"
	comments "forum/server/api/comment"
	authentification "forum/server/api/login"
	"forum/server/api/post"
	"forum/server/api/providers"
	users "forum/server/api/user"
	"forum/server/middleware"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// HTTPS
	/*cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatalf("Error loading certificate: %v", err)
	}

	// Configuration des cipher suites
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13, // Utilise TLS 1.3 comme version minimum
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			// Ajoute d'autres cipher suites selon tes besoins
		},
	}*/

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/api/post/createPost", post.CreatePostHandler)
	mux.HandleFunc("/api/post/fetchPost", post.FetchPostHandler)
	mux.HandleFunc("/api/post/fetchAllPost", post.FetchAllPostHandler)
	mux.HandleFunc("/api/post/deletePost", post.DeletePostHandler)
	mux.HandleFunc("/api/post/fetchPostMostLiked", post.FetchPostsMostLikedHandler)
	mux.HandleFunc("/api/post/fetchPostByUser", post.FetchUserPostHandler)

	mux.HandleFunc("/api/post/createComment", comments.CreateCommentHandler)
	mux.HandleFunc("/api/post/fetchComment", comments.FetchCommentHandler)
	mux.HandleFunc("/api/post/fetchAllComments", comments.FetchAllCommentsHandler)
	mux.HandleFunc("/api/post/deleteComment", comments.DeleteCommentHandler)
	mux.HandleFunc("/api/post/like-dislikeComment", comments.HandleLikeDislikeCommentAPI)
	mux.HandleFunc("/api/post/fetchCommentByUser", comments.FetchUserCommentsHandler)
	mux.HandleFunc("/api/post/fetchResponseUser", comments.FetchResponseUserHandler)

	mux.HandleFunc("/api/login", authentification.LoginHandler)
	mux.HandleFunc("/api/registration", authentification.RegisterHandler)
	mux.HandleFunc("/api/getSession", authentification.GetSession)

	mux.HandleFunc("/api/post/fetchAllCategories", categories.FetchAllCategoriesHandler)
	mux.HandleFunc("/api/post/fetchTendance", categories.FetchTendanceCategoriesHandler)
	mux.HandleFunc("/api/get-pp", authentification.PP_Handler)
	mux.HandleFunc("/api/google_login", providers.HandleGoogleLogin)
	mux.HandleFunc("/api/google_callback", providers.HandleGoogleCallback)

	mux.HandleFunc("/api/github_login", providers.HandleGithubLogin)
	mux.HandleFunc("/api/github_callback", providers.HandleGithubCallback)

	mux.HandleFunc("/api/discord_login", providers.HandleDiscordLogin)
	mux.HandleFunc("/api/discord_callback", providers.HandleDiscordCallback)
	mux.HandleFunc("/api/users/fetchAllUsers", users.FetchAllUsersHandler)

	mux.HandleFunc("/logout", users.LogoutHandler)
	mux.HandleFunc("/api/post/fetchPostsByCategories", categories.FetchPostByCategoriesHandler)
	mux.HandleFunc("/api/like-dislike", post.HandleLikeDislikeAPI)
	mux.HandleFunc("/authenticate", middleware.RateLimiterMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie("UserLogged"); err == nil {
			renderTemplate(w, "./static/homePage/index.html", nil)
		}
		renderTemplate(w, "./static/authentification/authentification.html", nil)
	}))

	mux.HandleFunc("/", middleware.RateLimiterMiddleware(authentification.HomeHandler))

	server := http.Server{
		Addr:              ":8080",
		Handler:           mux,
		MaxHeaderBytes:    1 << 26, // 4 MB
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      45 * time.Second,
		IdleTimeout:       3 * time.Minute,
		//TLSConfig:         tlsConfig,
	}

	// Lance une goroutine pour réinitialiser les compteurs périodiquement
	go func() {
		for {
			time.Sleep(middleware.Rl.Window)
			middleware.Rl.Cleanup()
		}
	}()

	log.Println("Server started on http://localhost:8080")

	// HTTPS
	//err = server.ListenAndServeTLS("", "") // "" car les certificats sont chargés via TLSConfig
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to start server: %s", err)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	w.Header().Set("Content-Type", "text/html")

	t, errTmpl := template.ParseFiles(tmpl)
	if errTmpl != nil {
		fmt.Fprintln(os.Stderr, errTmpl.Error())
		http.Error(w, "Error parsing template "+tmpl, http.StatusInternalServerError)
		return
	}

	if errExec := t.Execute(w, data); errExec != nil {
		fmt.Fprintln(os.Stderr, errExec.Error())
		return
	}

}
