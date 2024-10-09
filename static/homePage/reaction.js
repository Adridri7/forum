
export async function toggleReaction(event, postUuid) {
    const button = event.target.closest('.like-btn, .dislike-btn');
    const action = button.classList.contains('like-btn') ? 'like' : 'dislike';

    try {
        const response = await fetch('/api/like-dislike', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ postId: postUuid, action: action }),
            credentials: 'same-origin'
        });

        // Vérifie si la réponse est un succès
        if (!response.ok) {
            // Tente de récupérer l'erreur sous forme de texte ou JSON
            const errorMessage = await response.text();
            throw new Error(`Erreur lors de la mise à jour de la réaction: ${errorMessage}`);
        }

        // Si tout est ok, traite la réponse
        const result = await response.json();
        updateReactionUI(postUuid, result.likes, result.dislikes, result.userReaction);

    } catch (error) {
        // Affiche l'erreur complète ici
        console.error(error);
    }
}

function updateReactionUI(postUuid, likes, dislikes, userReaction) {
    const messageItem = document.querySelector(`.message-item[post_uuid="${postUuid}"]`);
    if (!messageItem) return;

    const likeBtn = messageItem.querySelector('.like-btn');
    const dislikeBtn = messageItem.querySelector('.dislike-btn');
    const likeCount = messageItem.querySelector('.like-count');
    const dislikeCount = messageItem.querySelector('.dislike-count');

    // Mettre à jour les compteurs
    if (likeCount) likeCount.textContent = likes;
    if (dislikeCount) dislikeCount.textContent = dislikes;

    // Mettre à jour l'apparence des boutons
    likeBtn.classList.toggle('active', userReaction.hasLiked);
    dislikeBtn.classList.toggle('active', userReaction.hasDisliked);
}

/*
function displayPosts(posts) {
    posts.forEach(post => {
        const postElement = createPostElement(post);

        const likeBtn = postElement.querySelector('.like-btn');
        const dislikeBtn = postElement.querySelector('.dislike-btn');

        // Initialiser l'état des boutons de réaction
        if (post.userReaction.hasLiked) {
            likeBtn.classList.add('active');
        }
        if (post.userReaction.hasDisliked) {
            dislikeBtn.classList.add('active');
        }

        // Ajouter l'élément au DOM
        document.querySelector('.posts-container').appendChild(postElement);
    });
}*/


export async function toggleCommentReaction(event, commentID,) {
    const button = event.target.closest('.like-comment-btn, .dislike-comment-btn');
    if (!button) {
        console.error('Button not found');
        return;
    }
    const action = button.classList.contains('like-comment-btn') ? 'like' : 'dislike';

    try {
        const response = await fetch('/api/post/like-dislikeComment', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ commentId: commentID, action: action }),
            credentials: 'same-origin'
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`Erreur lors de la mise à jour de la réaction: ${errorText}`);
        }

        const result = await response.json();
        updateReactionCommentUI(commentID, result.likes, result.dislikes, result.userReaction);

    } catch (error) {
        console.error('Erreur :', error);
    }
}

function updateReactionCommentUI(postUuid, likes, dislikes, userReaction) {
    const messageItem = document.querySelector(`.message-item[post_uuid="${postUuid}"]`);
    if (!messageItem) return;

    const likeBtn = messageItem.querySelector('.like-comment-btn');
    const dislikeBtn = messageItem.querySelector('.dislike-comment-btn');
    const likeCount = messageItem.querySelector('.like-count');
    const dislikeCount = messageItem.querySelector('.dislike-count');

    // Mettre à jour les compteurs
    if (likeCount) likeCount.textContent = likes;
    if (dislikeCount) dislikeCount.textContent = dislikes;

    // Mettre à jour l'apparence des boutons
    likeBtn.classList.toggle('active', userReaction.hasLiked);
    dislikeBtn.classList.toggle('active', userReaction.hasDisliked);
}