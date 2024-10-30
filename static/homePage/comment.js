import { fetchPosts, fetchUserInfo, headerBar, UserInfo } from "./app.js";
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


// Fonction pour gérer le clic sur le bouton de commentaire
export function handleCommentClick(section) {
    const postId = this.closest('.message-item')?.getAttribute('post_uuid');

    if (postId) {
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
            fetchPosts().then(initEventListeners); // Appelle initEventListeners après fetchPosts
            break;
        case AppState.POST:
            title.textContent = 'Post';
            displaySinglePost(newState.data.postId, newState.data.section);
            break;
        case AppState.SEARCH:
            title.textContent = 'Search';
            fetchCategories().then(initEventListeners); // Appelle initEventListeners après fetchCategories
            break;
        case AppState.TREND:
            title.textContent = 'Trending';
            FetchMostLikedPosts().then(initEventListeners); // Init si FetchMostLikedPosts retourne une Promise
            break;
        case AppState.NOTIFS:
            title.textContent = "Notifications";
            fetchNotifications(true).then(initEventListeners);
            break;
        case AppState.REQUEST:
            title.textContent = "Request";
            FetchHistoryRequest();
            break
        case AppState.MODERATION:
            title.textContent = "Moderation";
            FetchAdminRequest();
        default:
            newState.type = AppState.HOME;
            title.textContent = 'General';
            fetchPosts().then(initEventListeners);
            headerBar();
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
            default:
                url = '#home';
        }
        history.pushState(newState, '', url);
    }

}


export function displaySinglePost(postId, section) {
    const userPostsContainer = document.querySelector(`.users-post[data-section="${section}"]`);
    const postElement = userPostsContainer?.querySelector(`.message-item[post_uuid="${postId}"]`);

    if (userPostsContainer && postElement) {
        userPostsContainer.innerHTML = '';
        userPostsContainer.appendChild(postElement);

        const commentSection = createCommentInput(section);
        userPostsContainer.appendChild(commentSection);

        fetchAllcomments(postId);
    } else {
        console.error(`Post with id ${postId} not found in section "${section}"`);
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
        } else if (hash === '#trend') {
            updateAppState({ type: AppState.TREND }, false)
        } else if (hash === "#notifications") {
            updateAppState({ type: AppState.NOTIFS }, false)
        } else {
            updateAppState({ type: AppState.HOME }, false);
        }
    }
});

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
    if (!UserInfo || !UserInfo.user_uuid) {
        // Désactivez le champ d'entrée et le bouton si l'utilisateur n'est pas connecté
        commentInput.disabled = true;
        submitButton.disabled = true;

        // Ajoutez un message d'information
        const infoMessage = document.createElement('div');
        infoMessage.classList.add('info-message');
        commentInput.placeholder = 'Please login to comment';
        commentInputContainer.appendChild(infoMessage);
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
    const userPost = document.getElementById('users-post');
    const firstMessageItem = document.querySelector('.user-post .message-item');
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
            // fetchAllcomments(post_uuid);
        } else {
            const error = await response.json();
            alert("Erreur lors de l'ajout du commentaire: " + error.message);
        }
    } catch (error) {
        console.error("Erreur lors de l'envoi du commentaire:", error);
        alert("Une erreur s'est produite lors de l'envoi du commentaire. Veuillez réessayer.");
    }
}

// Fonction pour initialiser les événements des boutons
export function initEventListeners() {
    const commentButtons = document.querySelectorAll('.comment-btn');

    commentButtons.forEach(button => {
        const section = button.closest('.users-post')?.getAttribute('data-section');
        if (section) {
            button.removeEventListener('click', handleCommentClick.bind(button, section));
            button.addEventListener('click', handleCommentClick.bind(button, section));
        }
    });

    // Initialisation des boutons de menu
    const menuButtons = document.querySelectorAll('.menu-btn');
    menuButtons.forEach(button => {
        button.removeEventListener('click', handleMenuClick);
        button.addEventListener('click', handleMenuClick);
    });
}

export function handleMenuClick(event) {
    event.stopPropagation();
    const postElement = event.currentTarget.closest('.message-item');
    const postId = postElement.getAttribute('post_uuid');

    toggleMenu(event, postId);
}


document.addEventListener("DOMContentLoaded", () => {
    updateAppState({ type: AppState.HOME });
});

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
    updateAppState({ type: AppState.TREND })
});

document.getElementById('notifications-link').addEventListener('click', function (e) {
    e.preventDefault();
    if (!UserInfo) {
        return;
    }
    updateAppState({ type: AppState.NOTIFS })
});

document.getElementById('request-link').addEventListener('click', function (e) {
    e.preventDefault();
    updateAppState({ type: AppState.REQUEST })
});

document.getElementById('moderation-link').addEventListener('click', function (e) {
    e.preventDefault();
    updateAppState({ type: AppState.MODERATION })
});