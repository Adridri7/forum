// Fonction pour récupérer les posts les plus likés depuis l'API
export async function FetchMostLikedPosts() {
    try {
        const response = await fetch('/api/post/fetchPostMostLiked');
        if (!response.ok) {
            throw new Error('Erreur lors de la récupération des posts les plus likés');
        }
        const mostLikedPosts = await response.json();

        console.log("Post Ranking:", mostLikedPosts)

        // Affiche les posts dans la liste
        // displayMostLikedPost(mostLikedPosts);
    } catch (error) {
        console.error('Erreur lors de la récupération des posts :', error);
    }
}

function displayMostLikedPost(mostLikedPosts) {
    const postContainer = document.getElementById('most-liked-posts');
    postContainer.innerHTML = '';

    if (!mostLikedPosts || Object.keys(mostLikedPosts).length === 0) {
        postContainer.innerHTML = '<p>Aucun post populaire trouvé.</p>';
        return;
    }

    // Convertir l'objet en tableau et trier par nombre de likes décroissant
    const postsArray = Object.entries(mostLikedPosts)
        .sort((a, b) => b[1] - a[1]);

    postsArray.forEach(([postUuid, likesCount], index) => {
        const postElement = document.createElement('div');
        postElement.className = 'post-item';
        postElement.innerHTML = `
            <h3>${index + 1}. Post ID: ${postUuid}</h3>
            <span class="likes-count">👍 ${likesCount} Likes</span>
        `;
        postContainer.appendChild(postElement);
    });
}
