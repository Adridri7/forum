export async function fetchCategories() {
    try {
        const response = await fetch('/api/post/fetchAllCategories'); // Remplace par l'URL de ton API
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
    const categoriesList = document.getElementById('categories-list');
    categoriesList.innerHTML = ''; // Vide la liste existante

    categories.forEach(category => {
        const li = document.createElement('li');
        li.textContent = category;
        categoriesList.appendChild(li);
    });
}