import { fetchPosts, fetchUserInfo, UserInfo } from "./app.js";
import { toggleMenu } from "./displayMessage.js";
import { fetchCategories } from "./fetchcategories.js";
import { fetchNotifications } from "./notifs.js";
import { FetchMostLikedPosts } from "./postMostLiked.js";
import { FetchAdminRequest, FetchHistoryRequest } from "./API_request.js";
import { fetchAllcomments } from "./showComment.js";
import { getPPFromID } from "./utils.js";

const AppState = {
    HOME: 'home',
    POST: 'post',
    SEARCH: 'search',
    TREND: "trending",
    NOTIFS: "notification",
    REQUEST: "request",
    MODERATION: "moderation"
};

// Map pour stocker les event listeners
const commentListeners = new Map();
const menuListeners = new Map();

// Fonction pour gérer le clic sur le bouton de commentaire
export function handleCommentClick(section) {
    const postId = this.closest('.message-item')?.getAttribute('post_uuid');

    if (postId) {
        // Éviter les appels redondants si on est déjà sur le même post
        const currentPostElement = document.querySelector('.message-item');
        if (currentPostElement && currentPostElement.getAttribute('post_uuid') === postId) {
            return; // Ne rien faire si on clique sur le même post
        }

        updateAppState({
            type: AppState.POST,
            data: {
                postId: postId,
                section: section
            }
        });
    }
}

export function updateAppState(newState, pushState = true) {
    const title = document.getElementById('title');
    switch (newState.type) {
        case AppState.HOME:
            title.textContent = 'General';
            fetchPosts().then(initEventListeners);
            break;
        case AppState.POST:
            title.textContent = 'Post';
            displaySinglePost(newState.data.postId, newState.data.section)
                .then(initEventListeners)
                .catch(error => {
                    console.error('Error displaying post:', error);
                });
            break;
        case AppState.SEARCH:
            title.textContent = 'Search';
            fetchCategories().then(initEventListeners);
            break;
        case AppState.TREND:
            title.textContent = 'Trending';
            FetchMostLikedPosts().then(initEventListeners);
            break;
        case AppState.NOTIFS:
            title.textContent = "Notifications";
            fetchNotifications(true).then(initEventListeners);
            break;
        case AppState.REQUEST:
            title.textContent = "Request";
            FetchHistoryRequest();
            break;
        case AppState.MODERATION:
            title.textContent = "Moderation";
            FetchAdminRequest();
            break;
        default:
            newState.type = AppState.HOME;
            title.textContent = 'General';
            fetchPosts().then(initEventListeners);
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
            case AppState.TREND:
                url = '#trending';
                break;
            case AppState.NOTIFS:
                url = "#notifications";
                break;
            case AppState.MODERATION:
                url = "#moderation";
                break;
            case AppState.REQUEST:
                url = "#request";
                break;
            default:
                url = '#home';
                break;
        }
        history.pushState(newState, '', url);
    }
}

export async function displaySinglePost(postId, section) {
    const userPostsContainer = document.querySelector(`.users-post[data-section="${section}"]`);
    const postElement = userPostsContainer?.querySelector(`.message-item[post_uuid="${postId}"]`);

    if (userPostsContainer && postElement) {
        // Nettoyez d'abord le conteneur
        userPostsContainer.innerHTML = '';

        // Créez une copie profonde du post pour éviter les problèmes de référence
        const postClone = postElement.cloneNode(true);
        userPostsContainer.appendChild(postClone);

        // Créez et ajoutez la section de commentaires
        const commentSection = createCommentInput(section);
        userPostsContainer.appendChild(commentSection);

        try {
            await fetchAllcomments(postId);
            return Promise.resolve(); // Retourne une Promise résolue
        } catch (error) {
            console.error('Error fetching comments:', error);
            return Promise.reject(error);
        };
    } else {
        console.error(`Post with id ${postId} not found in section "${section}"`);
    }
}

// Fonction pour initialiser les événements des boutons
export function initEventListeners() {
    // Nettoyer et réinitialiser les boutons de commentaire
    const commentButtons = document.querySelectorAll('.comment-btn');
    commentButtons.forEach(button => {
        // Retirer l'ancien listener s'il existe
        const oldListener = commentListeners.get(button);
        if (oldListener) {
            button.removeEventListener('click', oldListener);
            commentListeners.delete(button);
        }

        // Ajouter le nouveau listener
        const section = button.closest('.users-post')?.getAttribute('data-section');
        if (section) {
            const newListener = (e) => handleCommentClick.call(button, section);
            commentListeners.set(button, newListener);
            button.addEventListener('click', newListener);
        }
    });

    // Nettoyer et réinitialiser les boutons de menu
    const menuButtons = document.querySelectorAll('.menu-btn');
    menuButtons.forEach(button => {
        // Retirer l'ancien listener s'il existe
        const oldListener = menuListeners.get(button);
        if (oldListener) {
            button.removeEventListener('click', oldListener);
            menuListeners.delete(button);
        }

        // Ajouter le nouveau listener
        const newListener = handleMenuClick;
        menuListeners.set(button, newListener);
        button.addEventListener('click', newListener);
    });
}

export function createCommentInput(section) {
    fetchUserInfo();

    const commentInputContainer = document.createElement('div');
    commentInputContainer.classList.add('comment-input-container');

    const profileImageContainer = document.createElement('div');
    profileImageContainer.classList.add('profile-image-container');

    const profileImage = document.createElement('img');
    if (UserInfo) {
        getPPFromID(UserInfo.user_uuid).then(img => { profileImage.src = img });
    } else {
        profileImage.src = "https://c.clc2l.com/t/d/i/discord-4OXyS2.png";
    }
    profileImage.alt = 'Profil-picture';
    profileImage.classList.add('profile-image');

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
    submitButton.classList.add('button-disabled');

    commentInput.addEventListener('input', () => {
        if (commentInput.value.trim() === '') {
            submitButton.disabled = true;
            submitButton.classList.add('button-disabled');
        } else {
            submitButton.disabled = false;
            submitButton.classList.remove('button-disabled');
        }
    });

    // Vérifier si l'utilisateur est connecté
    if (!UserInfo || !UserInfo.user_uuid) {
        commentInput.disabled = true;
        submitButton.disabled = true;
        commentInput.placeholder = 'Please login to comment';
    }

    form.appendChild(profileImageContainer);
    form.appendChild(commentInput);
    form.appendChild(submitButton);

    const usersPostContainer = document.querySelector(`.users-post[data-section="${section}"]`);
    const firstPostItem = usersPostContainer?.querySelector('.message-item');

    let postUuid = null;
    if (firstPostItem) {
        postUuid = firstPostItem.getAttribute('post_uuid');
    }

    form.addEventListener('submit', async (event) => {
        event.preventDefault();

        if (postUuid) {
            await createComment(postUuid, UserInfo.user_uuid);
            commentInput.value = '';
        } else {
            console.error("post_uuid est introuvable, le commentaire ne peut pas être créé.");
            alert("Impossible de soumettre le commentaire. Veuillez réessayer.");
        }
    });

    commentInputContainer.appendChild(form);
    return commentInputContainer;
}

export async function createComment(post_uuid, user_uuid) {
    const form = document.getElementById("create-comment-form");
    const formData = new FormData(form);

    const data = {};
    formData.forEach((value, key) => (data[key] = value));

    data.post_uuid = post_uuid;
    data.user_uuid = user_uuid;

    try {
        const response = await fetch("/api/post/createComment", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(data)
        });

        if (response.ok) {
            alert("Commentaire ajouté avec succès!");
        } else {
            const error = await response.json();
            alert("Erreur lors de l'ajout du commentaire: " + error.message);
        }
    } catch (error) {
        console.error("Erreur lors de l'envoi du commentaire:", error);
        alert("Une erreur s'est produite lors de l'envoi du commentaire. Veuillez réessayer.");
    }
}

export function handleMenuClick(event) {
    event.stopPropagation();
    const postElement = event.currentTarget.closest('.message-item');
    const postId = postElement.getAttribute('post_uuid');
    toggleMenu(event, postId);
}

// Event Listeners pour la navigation
window.addEventListener('popstate', function (event) {
    if (event.state) {
        updateAppState(event.state, false);
    } else {
        const hash = window.location.hash;
        if (hash.startsWith('#post-')) {
            const postId = hash.slice(6);
            updateAppState({ type: AppState.POST, data: { postId } }, false);
        } else if (hash === '#search') {
            updateAppState({ type: AppState.SEARCH }, false);
        } else if (hash === '#trend') {
            updateAppState({ type: AppState.TREND }, false);
        } else if (hash === "#notifications") {
            updateAppState({ type: AppState.NOTIFS }, false);
        } else {
            updateAppState({ type: AppState.HOME }, false);
        }
    }
});

// Initialisation des event listeners au chargement de la page
document.addEventListener("DOMContentLoaded", () => {
    updateAppState({ type: AppState.HOME });
});

// Event listeners pour la navigation
document.getElementById('home-link').addEventListener('click', function (e) {
    e.preventDefault();
    updateAppState({ type: AppState.HOME });
});

document.getElementById('search-link').addEventListener('click', function (e) {
    e.preventDefault();
    updateAppState({ type: AppState.SEARCH });
});

document.getElementById('trend-link').addEventListener('click', function (e) {
    e.preventDefault();
    updateAppState({ type: AppState.TREND });
});

document.getElementById('notifications-link').addEventListener('click', function (e) {
    e.preventDefault();
    if (!UserInfo) {
        return;
    }
    updateAppState({ type: AppState.NOTIFS });
});

document.getElementById('request-link').addEventListener('click', function (e) {
    e.preventDefault();
    updateAppState({ type: AppState.REQUEST });
});

document.getElementById('moderation-link').addEventListener('click', function (e) {
    e.preventDefault();
    updateAppState({ type: AppState.MODERATION });
});