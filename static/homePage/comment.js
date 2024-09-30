import { fetchPosts } from "./app.js";

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

    const commentSection = document.createElement('div');
    commentSection.innerHTML = `
        <div class="comment-section">
            <h3>Comments</h3>
        </div>
    `;
    userPostsContainer.appendChild(commentSection);

    // Sauvegarder le postId et l'état précédent dans l'historique
    history.pushState({ postId: postId, previousHTML: previousState, title: 'Post' }, `Post ${postId}`, `#post-${postId}`);

    // Réinitialiser les écouteurs d'événements
    initEventListeners();
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
    const commentButtons = document.querySelectorAll('.comment-btn');
    commentButtons.forEach(button => {
        button.removeEventListener('click', handleCommentClick);
        button.addEventListener('click', handleCommentClick);
    });
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