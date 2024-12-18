import { DisplayMessages } from "./displayMessage.js";

export async function fetchAllcomments(postUuid) {
    try {
        const response = await fetch(`/api/post/fetchComment`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ post_uuid: postUuid })
        });

        if (!response.ok) {
            throw new Error('Erreur lors de la récupération des commentaires');
        }

        const comments = await response.json();
        comments.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));

        comments.forEach(comment => {
            DisplayMessages(comment, "home", true);
        });

        return Promise.resolve();

    } catch (error) {
        console.error(error);
        return Promise.reject(error);
    }
}