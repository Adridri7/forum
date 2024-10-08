import { DisplayMessages } from "./displayMessage.js";
import { initEventListeners } from "./comment.js";
import { getUserInfoFromCookie, resetUsersPost } from "./utils.js";
import { NewPost } from "./newPost.js";
import { handleLogout } from "./logout.js";
import { fetchAllUsers } from "./displayUser.js";
import { toggleCommentReaction, toggleReaction } from "./reaction.js";
import { FetchMostLikedPosts } from "./postMostLiked.js";
import { FetchMostUseCategories } from "./tendance.js";
import { fetchPersonnalPosts } from "./dashboard.js";



// Sélectionnez les éléments
const toggleButton = document.getElementById('toggle-menu-btn');
const sidebar = document.getElementById('sidebar');
// const darkModeToggle = document.getElementById('dark-mode-toggle');
const addButton = document.getElementById('add-button');

const darkModeToggles = document.querySelectorAll('.dark-mode-toggle');

darkModeToggles.forEach(button => {
    button.addEventListener('click', () => {
        console.log('Dark mode toggled');
    });
});


// Fonction pour appliquer le mode
function applyMode(mode) {
    const root = document.documentElement;

    if (mode === 'dark') {
        root.style.setProperty('--background-color', '#1C1C1C');
        root.style.setProperty('--text-color', '#000000');
        root.style.setProperty('--second-text-color', '#FFFFFF');
        root.style.setProperty('--border-color', '#5E5E5F');
        root.style.setProperty('--background-message-color', '#272727');
        darkModeToggles.textContent = 'Light Mode';
    } else {
        root.style.setProperty('--background-color', '#f5f5f5');
        root.style.setProperty('--text-color', '#FFFFFF');
        root.style.setProperty('--second-text-color', '#000000');
        root.style.setProperty('--border-color', '#9C9FA8');
        root.style.setProperty('--background-message-color', '#FFFFFF');
        darkModeToggles.textContent = 'Dark Mode';
    }

    // Enregistrer la préférence dans le Local Storage
    localStorage.setItem('theme', mode);
}

// Vérifie la préférence au chargement
const userPreference = localStorage.getItem('theme');
if (userPreference) {
    applyMode(userPreference);
} else {
    // Si aucune préférence n'est trouvée, définir le mode par défaut (par exemple, light)
    applyMode('light');
}

// Écouteur d'événement pour le bouton de basculement
darkModeToggles.forEach(button => {
    button.addEventListener('click', () => {
        const currentMode = localStorage.getItem('theme') || 'light';
        const newMode = currentMode === 'dark' ? 'light' : 'dark';
        applyMode(newMode);
        // Mettre à jour le texte du bouton
        button.textContent = newMode === 'dark' ? 'Light Mode' : 'Dark Mode';
    });
});

document.addEventListener('DOMContentLoaded', () => {
    const loginButton = document.getElementById('login-btn');
    const profilMenu = document.querySelector('.profil-menu');

    // Fonction pour récupérer les informations utilisateur depuis le cookie
    const userInfo = getUserInfoFromCookie(); // Assure-toi que cette fonction existe
    console.log(userInfo)
    if (userInfo && userInfo.profileImageURL) {
        // Créer la div qui remplacera le bouton "Login"
        const profileDiv = document.createElement('div');
        profileDiv.classList.add('profile-container');

        const profileImage = document.createElement('img');
        profileImage.src = userInfo.profileImageURL;
        profileImage.alt = 'User profile';
        profileImage.classList.add('profile-image'); // Ajoute une classe pour styliser l'image

        // Créer le menu contextuel
        const menu = document.createElement('div');
        menu.classList.add('profile-menu'); // Classe pour styliser le menu
        menu.style.display = 'none'; // Masquer le menu par défaut

        // Créer le bouton "Logout"
        const logoutButton = document.createElement('button');
        logoutButton.id = "logout-btn"
        logoutButton.textContent = 'Log Out';
        logoutButton.addEventListener('click', handleLogout)

        // Ajouter le bouton "Logout" au menu
        menu.appendChild(logoutButton);

        // Ajouter l'image et le menu dans la div
        profileDiv.appendChild(profileImage);
        profileDiv.appendChild(menu);

        // Événement pour afficher/masquer le menu lorsque l'image est cliquée
        profileImage.addEventListener('click', () => {
            // Afficher ou masquer le menu
            if (menu.style.display === 'none') {
                menu.style.display = 'block';
            } else {
                menu.style.display = 'none';
            }
        });

        document.addEventListener('click', (event) => {
            if (!profileDiv.contains(event.target)) {
                menu.style.display = 'none';
            }
        });

        // Remplacer le bouton "Login" par la div
        profilMenu.replaceChild(profileDiv, loginButton);
    } else {
        // Si l'utilisateur n'est pas connecté ou a perdu le cookie, s'assurer que le bouton "Login" est visible
        if (!profilMenu.contains(loginButton)) {
            // Créer un nouveau bouton "Login" si nécessaire
            const newLoginButton = document.createElement('button');
            newLoginButton.id = 'login-btn';
            newLoginButton.textContent = 'Log in';

            // Ajouter à nouveau le bouton à la place de la div
            profilMenu.appendChild(newLoginButton);
        }

        addButton.style.display = 'none'
    }
});

// Événement pour le menu
toggleButton.addEventListener('click', () => {
    sidebar.classList.toggle('close');
});

let currentUser = ''
export async function fetchPosts() {
    resetUsersPost();
    const messagesList = document.getElementById('users-post');
    messagesList.innerHTML = '<p>Loading...</p>';
    try {
        const response = await fetch("http://localhost:8080/api/post/fetchAllPost");
        if (!response.ok) {
            throw new Error("Error retrieving posts");
        }

        const posts = await response.json();
        messagesList.innerHTML = '';

        if (posts.length === 0) {
            messagesList.innerHTML = '<p>No posts available.</p>';
        } else {
            currentUser = posts?.username || 'Anonymous';
            posts.sort((b, a) => new Date(b.created_at) - new Date(a.created_at));
            posts.forEach(post => {
                DisplayMessages(post);
            });
        }
    } catch (error) {
        messagesList.innerHTML = '<p>Error loading posts. Please try again.</p>';
        console.error(error);
    }
    initEventListeners();
}

export async function deletePost(post_uuid) {
    const confirmDelete = confirm("Êtes-vous sûr de vouloir supprimer ce post ?");
    if (!confirmDelete) return;

    try {
        const response = await fetch("http://localhost:8080/api/post/deletePost", {
            method: "DELETE",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ post_uuid: post_uuid }),
        });

        if (!response.ok) {
            throw new Error("Erreur lors de la suppression du post");
        }

        // Recharger les posts après suppression
        fetchPosts();
    } catch (error) {
        console.error(error);
    }
}

export function Reaction(event) {
    console.log("ça passe dans react ?");
    const likeButton = event.target.closest('.like-btn');
    const dislikeButton = event.target.closest('.dislike-btn');

    if (likeButton || dislikeButton) {
        const messageItem = (likeButton || dislikeButton).closest('.message-item');
        if (messageItem) {
            const postUuid = messageItem.getAttribute('post_uuid');
            toggleReaction(event, postUuid);
        }
    }
}

export function CommentReaction(event) {
    console.log("ça passe dans la fonction commentReact");
    const likeButton = event.target.closest('.like-comment-btn');
    const dislikeButton = event.target.closest('.dislike-comment-btn');

    console.log("ça trouve like button :", likeButton);
    console.log("ça trouve dislike button :", dislikeButton);

    if (likeButton || dislikeButton) {
        const messageItem = (likeButton || dislikeButton).closest('.message-item');
        if (messageItem) {
            const postUuid = messageItem.getAttribute('post_uuid');
            console.log("post_uuid :", postUuid);
            toggleCommentReaction(event, postUuid);
        }
    }
}

document.addEventListener('DOMContentLoaded', () => {
    fetchPosts();
    addButton.addEventListener('click', NewPost);
    fetchAllUsers();
    FetchMostUseCategories();
    FetchMostLikedPosts();

    // Gérer uniquement les événements sur les boutons de post
    document.body.addEventListener('click', (event) => {
        const isPostReaction = event.target.closest('.like-btn') || event.target.closest('.dislike-btn');
        const isCommentReaction = event.target.closest('.like-comment-btn') || event.target.closest('.dislike-comment-btn');

        if (isPostReaction) {
            Reaction(event); // Appeler la fonction pour les réactions sur les posts
        } else if (isCommentReaction) {
            CommentReaction(event); // Appeler la fonction pour les réactions sur les commentaires
        }
    });
});

const homeLink = document.getElementById('home-link');
const dashboardLink = document.getElementById('dashboard-link');

// Sélectionne les sections
const homeSection = document.getElementById('home-section');
const dashboardSection = document.getElementById('dashboard-section');

// Fonction pour masquer toutes les sections
function hideAllSections() {
    homeSection.style.display = 'none';
    dashboardSection.style.display = 'none';
}

// Ajoute des événements pour chaque lien
homeLink.addEventListener('click', () => {
    hideAllSections();
    homeSection.style.display = 'block';
});

dashboardLink.addEventListener('click', () => {
    const userInfo = getUserInfoFromCookie(); // Récupère les informations utilisateur
    const personnalPost = document.getElementById('personnal-post');

    if (!userInfo) {
        alert("You must be logged in to access the dashboard.");
        return;
    }

    hideAllSections();
    dashboardSection.style.display = 'flex';
    const currentActiveItem = document.querySelector('#nav-bar li.active');
    // Vérifier quel élément est actuellement actif et appeler la fonction correspondante
    if (currentActiveItem) {
        const activeId = currentActiveItem.id;

        // Fetch basé sur l'élément actif
        switch (activeId) {
            case 'personnal-post':
                fetchPersonnalPosts(); // Ou toute autre fonction que vous avez pour les posts
                break;
            case 'personnal-comment': // Si vous avez un élément avec id 'liked-posts'
                fetchPersonnalComment(); // Appelez la fonction appropriée ici
                break;
            case 'personnal-reaction':
                fetchPersonnalResponse();
            // Ajoutez d'autres cas selon vos besoins
            default:
                break;
        }
    }

    const navItems = document.querySelectorAll('#nav-bar li');

    navItems.forEach(item => {
        item.addEventListener('click', () => {
            // Retirer la classe active de tous les éléments
            navItems.forEach(nav => nav.classList.remove('active'));
            // Ajouter la classe active à l'élément cliqué
            item.classList.add('active');
        });
    });
});