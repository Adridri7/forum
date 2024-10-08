package main

import (
	"fmt"
	comments "forum/server/api/comment"
	"forum/server/api/post"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/api/post/createPost", post.CreatePostHandler)
	http.HandleFunc("/api/post/fetchPost", post.FetchPostHandler)
	http.HandleFunc("/api/post/fetchAllPost", post.FetchAllPostHandler)
	http.HandleFunc("/api/post/deletePost", post.DeletePostHandler)

	http.HandleFunc("/api/createComment", comments.CreateCommentHandler)

	// A faire pour tester : ajouter une route pour la page createPost.html
	http.Handle("/", http.FileServer(http.Dir("./static")))

	fmt.Println("Serveur démarré sur le port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
