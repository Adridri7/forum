import { fetchPosts, Reaction } from "./app.js";
import { toggleMenu } from "./displayMessage.js";
import { fetchCategories } from "./fetchcategories.js";
import { toggleCommentReaction } from "./reaction.js";
import { fetchAllcomments } from "./showComment.js";
import { getUserInfoFromCookie } from "./utils.js";

const AppState = {
    HOME: 'home',
    POST: 'post',
    SEARCH: 'search'
};

let currentState = {
    type: AppState.HOME,
    data: null
};

let previousState = null;

// Fonction pour gérer le clic sur le bouton de commentaire
function handleCommentClick() {
    const postId = this.closest('.message-item')?.getAttribute('post_uuid');

    if (postId) {
        updateAppState({
            type: AppState.POST,
            data: { postId: postId, previousHTML: document.getElementById('users-post').innerHTML }
        });
    }
}


function updateAppState(newState, pushState = true) {
    const title = document.getElementById('title');
    const userPostsContainer = document.getElementById('users-post');

    switch (newState.type) {
        case AppState.HOME:
            title.textContent = 'General';
            fetchPosts();
            break;
        case AppState.POST:
            title.textContent = 'Post';
            displaySinglePost(newState.data.postId);
            break;
        case AppState.SEARCH:
            title.textContent = 'Search';
            fetchCategories();
            break;
        default:
            // Si l'état n'est pas reconnu, revenez à l'accueil
            newState.type = AppState.HOME;
            title.textContent = 'General';
            fetchPosts();
            break;
    }

    if (pushState) {
        let url;
        switch (newState.type) {
            case AppState.POST:
                url = `#post-${newState.data.postId}`;
                break;
            case AppState.SEARCH:
                url = '#search';
                break;
            default:
                url = '#home';
        }
        history.pushState(newState, '', url);
    }

    currentState = newState;
    initEventListeners();
}


function displaySinglePost(postId) {
    const userPostsContainer = document.getElementById('users-post');
    const postElement = userPostsContainer.querySelector(`.message-item[post_uuid="${postId}"]`);

    if (postElement) {
        userPostsContainer.innerHTML = '';
        userPostsContainer.appendChild(postElement);

        const commentSection = createCommentInput();
        userPostsContainer.appendChild(commentSection);

        fetchAllcomments(postId);
    } else {
        console.error(`Post with id ${postId} not found`);
    }
}


window.addEventListener('popstate', function (event) {
    if (event.state) {
        updateAppState(event.state, false);
    } else {
        // Détermine l'état en fonction de l'URL actuelle
        const hash = window.location.hash;
        if (hash.startsWith('#post-')) {
            const postId = hash.slice(6); // Retire '#post-' du début
            updateAppState({ type: AppState.POST, data: { postId } }, false);
        } else if (hash === '#search') {
            updateAppState({ type: AppState.SEARCH }, false);
        } else {
            updateAppState({ type: AppState.HOME }, false);
        }
    }
});

function createCommentInput() {
    const userInfo = getUserInfoFromCookie();

    const commentInputContainer = document.createElement('div');
    commentInputContainer.classList.add('comment-input-container');

    // Créer un conteneur pour l'image de profil
    const profileImageContainer = document.createElement('div');
    profileImageContainer.classList.add('profile-image-container'); // Ajoutez une classe pour le style si nécessaire

    const profileImage = document.createElement('img');
    profileImage.src = userInfo.profileImageURL;
    profileImage.alt = 'Profil-picture';
    profileImage.classList.add('profile-image');

    // Ajouter l'image de profil au conteneur
    profileImageContainer.appendChild(profileImage);

    const form = document.createElement('form');
    form.id = 'create-comment-form';

    const commentInput = document.createElement('input');
    commentInput.type = 'text';
    commentInput.name = 'content';
    commentInput.placeholder = `Post your reply`;
    commentInput.classList.add('comment-input');

    const submitButton = document.createElement('button');
    submitButton.type = 'submit';
    submitButton.textContent = 'Post';
    submitButton.disabled = true;
    submitButton.classList.add('button-disabled')

    commentInput.addEventListener('input', () => {
        if (commentInput.value.trim() === '') {
            submitButton.disabled = true;
            submitButton.classList.add('button-disabled');
        } else {
            submitButton.disabled = false;
            submitButton.classList.remove('button-disabled');
        }
    });

    // Vérifiez si l'utilisateur est connecté
    if (!userInfo || !userInfo.uuid) {
        // Désactivez le champ d'entrée et le bouton si l'utilisateur n'est pas connecté
        commentInput.disabled = true;
        submitButton.disabled = true;

        // Ajoutez un message d'information
        const infoMessage = document.createElement('div');
        infoMessage.classList.add('info-message');
        commentInput.placeholder = 'Please login to comment';
        commentInputContainer.appendChild(infoMessage);
    }


    // Ajouter le conteneur de l'image et les autres éléments au formulaire
    form.appendChild(profileImageContainer);
    form.appendChild(commentInput);
    form.appendChild(submitButton);

    const usersPostContainer = document.getElementById('users-post');
    const firstPostItem = usersPostContainer.querySelector('.message-item');

    let postUuid = null;
    if (firstPostItem) {
        postUuid = firstPostItem.getAttribute('post_uuid');
        console.log("post_uuid:", postUuid);
    } else {
        console.log("Aucun post trouvé.");
    }

    form.addEventListener('submit', async (event) => {
        event.preventDefault();

        console.log(userInfo.uuid);

        if (postUuid) {
            await createComment(postUuid, userInfo.uuid);
            commentInput.value = '';
        } else {
            console.error("post_uuid est introuvable, le commentaire ne peut pas être créé.");
            alert("Impossible de soumettre le commentaire. Veuillez réessayer.");
        }
    });

    commentInputContainer.appendChild(form);
    return commentInputContainer;
}

async function createComment(post_uuid, user_uuid) {
    const form = document.getElementById("create-comment-form");
    const formData = new FormData(form);

    const data = {};
    formData.forEach((value, key) => (data[key] = value));

    data.post_uuid = post_uuid;
    data.user_uuid = user_uuid;

    console.log("Données envoyées:", data);

    try {
        const response = await fetch("/api/post/createComment", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(data)
        });

        if (response.ok) {
            console.log(data)
            alert("Commentaire ajouté avec succès!");
            // Vous pourriez vouloir actualiser la liste des commentaires ici
        } else {
            const error = await response.json();
            alert("Erreur lors de l'ajout du commentaire: " + error.message);
        }
    } catch (error) {
        console.error("Erreur lors de l'envoi du commentaire:", error);
        alert("Une erreur s'est produite lors de l'envoi du commentaire. Veuillez réessayer.");
    }
}

// Fonction pour supprimer tous les messages sauf le post créé
// function removeAllMessagesExceptCreatedPost() {
//     const messageItems = document.querySelectorAll('.message-item');
//     messageItems.forEach(item => {
//         item.remove();
//     });
// }

// Fonction pour initialiser les événements des boutons
export function initEventListeners() {
    // Sélectionne tous les boutons de commentaire
    const commentButtons = document.querySelectorAll('.comment-btn');

    // Réinitialise les événements pour chaque bouton de commentaire
    commentButtons.forEach(button => {
        button.removeEventListener('click', handleCommentClick);
        button.addEventListener('click', handleCommentClick);
    });

    // const reactionCommentButton = document.getElementById('reaction-comment-btn');
    // reactionCommentButton.addEventListener('click', CommentReaction);

    // Sélectionne tous les boutons de menu (dans chaque message-item)
    const menuButtons = document.querySelectorAll('.menu-btn');

    // Réinitialise les événements pour chaque bouton de menuÒ
    menuButtons.forEach(button => {
        button.removeEventListener('click', handleMenuClick);
        button.addEventListener('click', handleMenuClick);
    });
}

function handleMenuClick(event) {
    event.stopPropagation();
    const postElement = event.currentTarget.closest('.message-item');
    const postId = postElement.getAttribute('post_uuid');

    toggleMenu(event, postId);
}


// Initialisation des événements des boutons au chargement de la page
document.addEventListener("DOMContentLoaded", () => {
    // Toujours commencer par la page d'accueil
    updateAppState({ type: AppState.HOME });
});

// Gestion du lien #home pour revenir à la page principale avec tous les posts
document.getElementById('home-link').addEventListener('click', function (e) {
    e.preventDefault();
    updateAppState({ type: AppState.HOME });
});

document.getElementById('search-link').addEventListener('click', function (e) {
    e.preventDefault();
    updateAppState({ type: AppState.SEARCH });
});