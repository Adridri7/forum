import { DisplayMessages } from "./displayMessage.js";

export async function fetchAllcomments() {
    try {
        const response = await fetch('/api/post/fetchAllComments'); // Remplace par l'URL de ton API
        if (!response.ok) {
            throw new Error('Erreur lors de la récupération des commentaires');
        }


        const comment = await response.json();
        console.log(comment)
        console.log(comment[0].content)

        comment.forEach(comment => {
            DisplayMessages(comment);
        });

    } catch (error) {
        console.error(error);
    }
}