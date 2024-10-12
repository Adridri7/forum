import { fetchPosts } from "./app.js";
import { fetchAllcomments } from "./showComment.js";

export async function deletePost(post_uuid) {
    const confirmDelete = confirm("Êtes-vous sûr de vouloir supprimer ce post ?");
    if (!confirmDelete) return;

    try {
        const response = await fetch("http://localhost:8080/api/post/deletePost", {
            method: "DELETE",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ post_uuid: post_uuid }),
        });

        if (!response.ok) {
            throw new Error("Erreur lors de la suppression du post");
        }

        // Recharger les posts après suppression
        fetchPosts();
    } catch (error) {
        console.error(error);
    }
}

export async function updatePost(post_uuid, updatedContent) {
    try {
        const response = await fetch("/api/post/UpdatePost", {
            method: "PUT",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                post_uuid: post_uuid,
                content: updatedContent
            })
        });

        if (!response.ok) {
            // Handle the error
            const errorMessage = await response.text();
            alert(`Failed to update the post: ${errorMessage}`);
            return;
        }

    } catch (error) {
        alert(`An error occurred: ${error.message}`); // Handle network or other errors
    }
}

export async function updateComment(comment_id, updatedContent) {
    try {
        const response = await fetch("/api/post/UpdateComment", {
            method: "PUT",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                comment_id: comment_id,
                content: updatedContent
            })
        });

        if (!response.ok) {
            // Handle the error
            const errorMessage = await response.text();
            alert(`Failed to update the post: ${errorMessage}`);
            return;
        }

    } catch (error) {
        alert(`An error occurred: ${error.message}`); // Handle network or other errors
    }
}

export async function deleteComment(post_uuid, comment_id) {
    const confirmDelete = confirm("Êtes-vous sûr de vouloir supprimer ce post ?");
    if (!confirmDelete) return;

    try {
        const response = await fetch("/api/post/deleteComment", {
            method: "DELETE",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ comment_id: comment_id }),
        });

        if (!response.ok) {
            const error = await response.text()
            alert(`Failed to Delete the post: ${error}`);
        }

        fetchAllcomments(post_uuid);
    } catch (error) {
        console.error(error);
    }
}