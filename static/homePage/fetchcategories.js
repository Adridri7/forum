import { DisplayMessages } from "./displayMessage.js";
import { initEventListeners } from "./comment.js";
import { resetUsersPost, startGradientAnimation } from "./utils.js";

// Fonction pour récupérer toutes les catégories depuis l'API
export async function fetchCategories() {
    try {
        const response = await fetch('/api/post/fetchAllCategories'); // Remplace par l'URL de ton API
        if (!response.ok) {
            throw new Error('Erreur lors de la récupération des catégories');
        }
        const categories = await response.json();

        // Affiche les catégories dans la liste et ajoute des écouteurs d'événements
        displayCategories(categories);
    } catch (error) {
        console.error(error);
    }
}

// Fonction pour afficher les catégories dans la liste
function displayCategories(categories) {
    const usersPost = document.getElementById('users-post');

    // Vide le contenu existant de users-post
    usersPost.innerHTML = '';

    // Modifier le style de #users-post
    usersPost.style.display = 'grid';
    usersPost.style.gridTemplateColumns = 'repeat(2, 1fr)';
    usersPost.style.gridAutoRows = 'min-content';
    usersPost.style.rowGap = '10px';
    usersPost.style.columnGap = '10px';
    usersPost.style.border = 'none';

    categories.forEach(category => {
        const li = document.createElement('li');
        li.textContent = category;
        li.setAttribute('data-category', category); // Définit un attribut personnalisé pour la catégorie
        li.classList.add('categorie'); // Ajoute la classe 'categorie'
        li.style.cursor = 'pointer'

        // Démarre l'animation de dégradé
        startGradientAnimation(li); // Appel à la fonction d'animation de dégradé

        li.addEventListener('click', () => {
            fetchPostsByCategory(category);
        });

        // Ajoute l'élément à users-post
        usersPost.appendChild(li);
    });
}

export async function fetchPostsByCategory(category) {
    resetUsersPost();

    try {
        const response = await fetch(`/api/post/fetchPostsByCategories`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ categories: category })
        });

        if (!response.ok) {
            throw new Error('Erreur lors de la récupération des articles pour cette catégorie');
        }

        const posts = await response.json();
        const postsContainer = document.getElementById('users-post');

        // Vérifiez que postsContainer n'est pas null avant de l'utiliser
        if (postsContainer) {
            postsContainer.innerHTML = '';

            posts.forEach(post => {
                DisplayMessages(post);
            });
            initEventListeners();
        } else {
            console.error("L'élément avec l'ID 'posts-container' n'a pas été trouvé dans le DOM.");
        }

    } catch (error) {
        console.error("Erreur:", error);
    }
}
