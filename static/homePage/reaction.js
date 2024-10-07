
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
        console.log(result);
        updateReactionUI(postUuid, result.likes, result.dislikes, result.userReaction);

        // Optionnel : Rafraîchir tous les posts si nécessaire
        // fetchPosts();
    } catch (error) {
        console.error('Erreur :', error);
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

    console.log('Message item:', messageItem);
    console.log('Like button:', likeBtn);
    console.log('Dislike button:', dislikeBtn);
    console.log('User has liked:', userReaction.hasLiked);
    console.log('User has disliked:', userReaction.hasDisliked);

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
    console.log('toggleReactionComment:', commentID);
    const button = event.target.closest('.like-comment-btn, .dislike-comment-btn');
    if (!button) {
        console.error('Button not found');
        return;
    }
    const action = button.classList.contains('like-comment-btn') ? 'like' : 'dislike';
    console.log('Action:', action);

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
        console.log("Résultat:", result);
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

    console.log('Message item:', messageItem);
    console.log('Like button:', likeBtn);
    console.log('Dislike button:', dislikeBtn);
    console.log('User has liked:', userReaction.hasLiked);
    console.log('User has disliked:', userReaction.hasDisliked);

    // Mettre à jour l'apparence des boutons
    likeBtn.classList.toggle('active', userReaction.hasLiked);
    dislikeBtn.classList.toggle('active', userReaction.hasDisliked);
}