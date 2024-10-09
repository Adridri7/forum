package main

import (
	"fmt"
	"forum/server/api/categories"
	comments "forum/server/api/comment"
	authentification "forum/server/api/login"
	"forum/server/api/post"
	"forum/server/api/providers"
	users "forum/server/api/user"
	"html/template"
	"net/http"
	"os"
	"time"
)

func main() {

	// if err := providers.LoadEnvVariables(); err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error %v\n", err)
	// 	return
	// }

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

	mux.HandleFunc("/authenticate", func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie("UserLogged"); err == nil {
			renderTemplate(w, "./static/homePage/index.html", nil)
		}
		renderTemplate(w, "./static/authentification/authentification.html", nil)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "./static/homePage/index.html", nil)
	})

	ourServer := http.Server{
		Addr:              ":8080",
		Handler:           mux,
		MaxHeaderBytes:    1 << 26, // 4 MB
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      45 * time.Second,
		IdleTimeout:       3 * time.Minute,
	}

	fmt.Println("Serveur démarré : http://localhost:8080/")
	fmt.Fprintln(os.Stderr, ourServer.ListenAndServe())
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
