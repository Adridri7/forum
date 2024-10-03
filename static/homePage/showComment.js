import { DisplayMessages } from "./displayMessage.js";

export async function fetchAllcomments() {

    const firstPostItem = document.querySelector('[post_uuid]');

    // Extrait la valeur de `post_uuid` de l'élément HTML
    const postUuid = firstPostItem.getAttribute('post_uuid');

    console.log("fetch comment: ", postUuid)

    try {
        const response = await fetch(`/api/post/fetchComment`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            // Envoi le post_uuid au backend sous form de JSON pour traitement
            body: JSON.stringify({ post_uuid: postUuid })
        });


        if (!response.ok) {
            throw new Error('Erreur lors de la récupération des commentaires');
        }

        // Récupère et parse les commentaires retournés par le backend à partir de post_uuid
        const comment = await response.json();
        comment.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));

        comment.forEach(comment => {
            DisplayMessages(comment);
        });

    } catch (error) {
        console.error(error);
    }
}