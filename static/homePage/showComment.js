import { DisplayMessages } from "./displayMessage.js";

export async function fetchAllcomments() {

    const firstPostItem = document.querySelector('[post_uuid]');
    const postUuid = firstPostItem.getAttribute('post_uuid');

    console.log("fetch comment: ", postUuid)

    try {
        const response = await fetch(`/api/post/fetchComment`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            // Pour récupérer un éléments précis
            body: JSON.stringify({ post_uuid: postUuid })
        });


        if (!response.ok) {
            throw new Error('Erreur lors de la récupération des commentaires');
        }


        const comment = await response.json();

        console.log(comment[0].content)

        comment.forEach(comment => {
            DisplayMessages(comment);
        });

    } catch (error) {
        console.error(error);
    }
}