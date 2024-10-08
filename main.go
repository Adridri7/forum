package main

import (
	"fmt"
	"forum/server/api/categories"
	comments "forum/server/api/comment"
	authentification "forum/server/api/login"
	"forum/server/api/post"
	users "forum/server/api/user"
	"html/template"
	"net/http"
	"os"
)

func main() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/api/post/createPost", post.CreatePostHandler)
	http.HandleFunc("/api/post/fetchPost", post.FetchPostHandler)
	http.HandleFunc("/api/post/fetchAllPost", post.FetchAllPostHandler)
	http.HandleFunc("/api/post/deletePost", post.DeletePostHandler)
	http.HandleFunc("/api/post/fetchPostMostLiked", post.FetchPostsMostLikedHandler)

	http.HandleFunc("/api/post/createComment", comments.CreateCommentHandler)
	http.HandleFunc("/api/post/fetchComment", comments.FetchCommentHandler)
	http.HandleFunc("/api/post/fetchAllComments", comments.FetchAllCommentsHandler)
	http.HandleFunc("/api/post/deleteComment", comments.DeleteCommentHandler)
	http.HandleFunc("/api/post/like-dislikeComment", comments.HandleLikeDislikeCommentAPI)

	http.HandleFunc("/api/login", authentification.LoginHandler)
	http.HandleFunc("/api/registration", authentification.RegisterHandler)

	http.HandleFunc("/api/post/fetchAllCategories", categories.FetchAllCategoriesHandler)
	http.HandleFunc("/api/post/fetchTendance", categories.FetchTendanceCategoriesHandler)
	http.HandleFunc("/api/get-pp", authentification.PP_Handler)
	http.HandleFunc("/api/users/fetchAllUsers", users.FetchAllUsersHandler)

	http.HandleFunc("/logout", users.LogoutHandler)

	http.HandleFunc("/api/post/fetchPostsByCategories", categories.FetchPostByCategoriesHandler)

	http.HandleFunc("/api/like-dislike", post.HandleLikeDislikeAPI)

	http.HandleFunc("/authenticate", func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie("UserLogged"); err == nil {
			renderTemplate(w, "./static/homePage/index.html", nil)
		}
		renderTemplate(w, "./static/authentification/authentification.html", nil)
	})

	// A faire pour tester : ajouter une route pour la page createPost.html
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "./static/homePage/index.html", nil)
	})

	fmt.Println("Serveur démarré : http://localhost:8080/")
	fmt.Fprintln(os.Stderr, http.ListenAndServe(":8080", nil))
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
