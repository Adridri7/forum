import { DisplayMessages } from "./displayMessage.js";
import { initEventListeners } from "./comment.js";
import { resetUsersPost, startGradientAnimation } from "./utils.js";

// Fonction pour récupérer toutes les catégories depuis l'API
export async function fetchCategories() {
    try {
        const response = await fetch('/api/post/fetchAllCategories');
        if (!response.ok) {
            throw new Error('Erreur lors de la récupération des catégories');
        }
        const categories = await response.json();

        displayCategories(categories);
    } catch (error) {
        console.error(error);
    }
}

// Fonction pour afficher les catégories dans la liste
function displayCategories(categories) {
    const usersPost = document.getElementById('users-post');

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
        li.setAttribute('data-category', category);
        li.classList.add('categorie');
        li.style.cursor = 'pointer'

        startGradientAnimation(li);

        li.addEventListener('click', () => {
            fetchPostsByCategory(category);
        });

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
