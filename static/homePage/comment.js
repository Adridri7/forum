import { fetchPosts } from "./app.js";
import { toggleMenu } from "./displayMessage.js";
import { fetchAllcomments } from "./showComment.js";
import { getUserInfoFromCookie } from "./utils.js";

let previousState = null;

// Fonction pour gérer le clic sur le bouton de commentaire
function handleCommentClick() {

    const postId = this.closest('.message-item')?.getAttribute('post_uuid');
    const postElement = this.closest('.message-item');

    previousState = document.getElementById('users-post').innerHTML;

    const title = document.getElementById('title');
    title.textContent = 'Post';

    const userPostsContainer = document.getElementById('users-post');

    removeAllMessagesExceptCreatedPost();
    userPostsContainer.appendChild(postElement);

    const commentSection = createCommentInput();

    userPostsContainer.appendChild(commentSection);
    fetchAllcomments();

    // Sauvegarder le postId et l'état précédent dans l'historique
    history.pushState({ postId: postId, previousHTML: previousState, title: 'Post' }, `Post ${postId}`, `#post-${postId}`);

    // Réinitialiser les écouteurs d'événements
    initEventListeners();
}

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
function removeAllMessagesExceptCreatedPost() {
    const messageItems = document.querySelectorAll('.message-item');
    messageItems.forEach(item => {
        item.remove();
    });
}

// Fonction pour initialiser les événements des boutons
export function initEventListeners() {
    // Sélectionne tous les boutons de commentaire
    const commentButtons = document.querySelectorAll('.comment-btn');

    // Réinitialise les événements pour chaque bouton de commentaire
    commentButtons.forEach(button => {
        button.removeEventListener('click', handleCommentClick);
        button.addEventListener('click', handleCommentClick);
    });

    // Sélectionne tous les boutons de menu (dans chaque message-item)
    const menuButtons = document.querySelectorAll('.menu-btn');

    // Réinitialise les événements pour chaque bouton de menu
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

// Gestion du lien #home pour revenir à la page principale avec tous les posts
document.getElementById('home-link').addEventListener('click', function (e) {
    e.preventDefault();

    previousState = document.getElementById('users-post').innerHTML;

    // Mettre à jour le titre à "General"
    const title = document.getElementById('title');
    title.textContent = 'General';

    // Recharger les posts
    fetchPosts();

    // Sauvegarder l'état de la page d'accueil dans l'historique
    history.pushState({ title: 'General', previousHTML: previousState }, 'Home', '#home');
});

// Écouter les événements de navigation (popstate)
window.addEventListener('popstate', function (event) {

    if (event.state) {

        // Récupérer l'état précédent
        const previousHTML = event.state.previousHTML;
        const restoredTitle = event.state.title;

        // Rétablir l'état précédent
        const userPostsContainer = document.getElementById('users-post');
        userPostsContainer.innerHTML = previousHTML || '';

        // Mettre à jour le titre avec le titre restauré
        const title = document.getElementById('title');
        title.textContent = restoredTitle;

        // Mettre à jour l'URL sans ajouter une nouvelle entrée dans l'historique
        const newUrl = event.state.postId ? `#post-${event.state.postId}` : '#home';
        history.replaceState(event.state, '', newUrl);

        initEventListeners();
    } else {
        fetchPosts();

        // Mettre à jour le titre pour refléter l'état "General"
        const title = document.getElementById('title');
        title.textContent = 'General';

        // Mettre à jour l'URL pour refléter l'état général
        history.replaceState(null, '', '#home');
    }
});

// Initialisation des événements des boutons au chargement de la page
document.addEventListener("DOMContentLoaded", () => {
    fetchPosts();
    initEventListeners();
});