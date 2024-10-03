import { DisplayMessages } from "./displayMessage.js";
import { initEventListeners } from "./comment.js";
import { fetchCategories } from "./fetchcategories.js";
import { getUserInfoFromCookie } from "./utils.js";
import { fetchPostsByCategory } from './fetchcategories.js';
import { NewPost } from "./newPost.js";
import { handleLogout } from "./logout.js";
import { fetchAllUsers } from "./displayUser.js";



// Sélectionnez les éléments
const toggleButton = document.getElementById('toggle-menu-btn');
const sidebar = document.getElementById('sidebar');
const darkModeToggle = document.getElementById('dark-mode-toggle');
const addButton = document.getElementById('add-button');


// Fonction pour appliquer le mode
function applyMode(mode) {
    const root = document.documentElement;

    if (mode === 'dark') {
        root.style.setProperty('--background-color', '#1C1C1C');
        root.style.setProperty('--text-color', '#000000');
        root.style.setProperty('--second-text-color', '#FFFFFF');
        root.style.setProperty('--border-color', '#5E5E5F');
        root.style.setProperty('--background-message-color', '#272727');
        darkModeToggle.textContent = 'Light Mode';
    } else {
        root.style.setProperty('--background-color', '#f5f5f5');
        root.style.setProperty('--text-color', '#FFFFFF');
        root.style.setProperty('--second-text-color', '#000000');
        root.style.setProperty('--border-color', '#9C9FA8');
        root.style.setProperty('--background-message-color', '#FFFFFF');
        darkModeToggle.textContent = 'Dark Mode';
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
darkModeToggle.addEventListener('click', () => {
    const currentMode = localStorage.getItem('theme') || 'light';
    const newMode = currentMode === 'dark' ? 'light' : 'dark';
    applyMode(newMode);
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
    }
});

// Événement pour le menu
toggleButton.addEventListener('click', () => {
    sidebar.classList.toggle('close');
});

let currentUser = ''
export async function fetchPosts() {
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

document.addEventListener('DOMContentLoaded', () => {
    fetchPosts();
    fetchCategories();
    addButton.addEventListener('click', NewPost);
    fetchAllUsers();
    fetchPostsByCategory();
});
