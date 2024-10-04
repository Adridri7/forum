import { fetchPosts } from './app.js';  // Assurez-vous que le chemin est correct

export async function toggleReaction(event, postUuid) {
    console.log('toggleReaction called', postUuid);
    const button = event.target.closest('.like-btn, .dislike-btn');
    const action = button.classList.contains('like-btn') ? 'like' : 'dislike';
    console.log('Action:', action);

    try {
        const response = await fetch('/api/like-dislike', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ postId: postUuid, action: action }),
            credentials: 'same-origin'
        });

        if (!response.ok) {
            throw new Error('Erreur lors de la mise à jour de la réaction');
        }

        const result = await response.json();
        updateReactionUI(postUuid, result.likes, result.dislikes);

        // Optionnel : Rafraîchir tous les posts si nécessaire
        // fetchPosts();
    } catch (error) {
        console.error('Erreur :', error);
    }
}

function updateReactionUI(postUuid, likes, dislikes) {
    const messageItem = document.querySelector(`.message-item[post_uuid="${postUuid}"]`);
    if (!messageItem) return;

    const likeBtn = messageItem.querySelector('.like-btn');
    const dislikeBtn = messageItem.querySelector('.dislike-btn');
    const likeCount = messageItem.querySelector('.like-count');
    const dislikeCount = messageItem.querySelector('.dislike-count');

    if (likeCount) likeCount.textContent = likes;
    if (dislikeCount) dislikeCount.textContent = dislikes;

    // Mettre à jour l'apparence des boutons si nécessaire
    likeBtn.classList.toggle('active', likes > 0);
    dislikeBtn.classList.toggle('active', dislikes > 0);
}
