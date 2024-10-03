import { DisplayMessages } from "./displayMessage.js";
import { initEventListeners } from "./comment.js";

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
    const categoriesList = document.getElementById('categories-list');
    categoriesList.innerHTML = ''; // Vide la liste existante

    categories.forEach(category => {
        const li = document.createElement('li');
        li.textContent = category;
        li.setAttribute('data-category', category); // Définit un attribut personnalisé pour la catégorie

        // Ajoute un écouteur d'événement pour gérer le clic sur chaque catégorie
        li.addEventListener('click', () => {
            fetchPostsByCategory(category); // Appelle la fonction pour récupérer les posts de cette catégorie
        });

        categoriesList.appendChild(li);
    });
}

export async function fetchPostsByCategory(category) {

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
