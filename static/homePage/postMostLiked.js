import { DisplayMessages } from "./displayMessage.js";
import { resetUsersPost } from "./utils.js";

// Fonction pour récupérer les posts les plus likés depuis l'API
export async function FetchMostLikedPostss() {
    try {
        const response = await fetch('/api/post/fetchPostMostLiked');
        if (!response.ok) {
            throw new Error('Erreur lors de la récupération des posts les plus likés');
        }
        const mostLikedPosts = await response.json();

        console.log("Post Ranking:", mostLikedPosts)

        // Affiche les posts dans la liste
        displayMostLikedPost(mostLikedPosts);
    } catch (error) {
        console.error('Erreur lors de la récupération des posts :', error);
    }
}

export async function FetchMostLikedPosts() {
    resetUsersPost();
    const usersPost = document.getElementById('users-post');
    try {
        const response = await fetch('/api/post/fetchMostLikedPost');

        if (!response.ok) {
            const errorResponse = await response.json();
            throw new Error(`Erreur ${response.status}: ${response.statusText} - ${errorResponse.message || 'Détails supplémentaires non disponibles'}`);
        }

        const mostLikedPosts = await response.json();

        usersPost.innerHTML = '';

        // Affiche les posts dans la liste

        mostLikedPosts.forEach(posts => {
            DisplayMessages(posts);
        });
    } catch (error) {
        console.error('Erreur lors de la récupération des posts :', error.message);
    }
}
