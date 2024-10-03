export async function fetchAllUsers() {
    try {
        const response = await fetch('/api/users/fetchAllUsers'); // Remplace par l'URL de ton API
        if (!response.ok) {
            throw new Error('Erreur lors de la récupération des utilisateurs');
        }
        const users = await response.json();

        displayUsers(users);
    } catch (error) {
        console.error(error);
    }
}

// Fonction pour afficher les utilisateurs dans la liste
function displayUsers(users) {
    const usersList = document.getElementById('user-list');
    usersList.innerHTML = ''; // Vide la liste existante

    users.forEach(user => {
        // Création de l'élément <li> pour chaque utilisateur
        const li = document.createElement('li');
        li.style.display = 'flex'; // Pour aligner l'image et le texte
        li.style.alignItems = 'center';
        li.style.marginBottom = '10px'; // Espacement entre les éléments

        // Ajout de l'image de profil
        const img = document.createElement('img');
        img.src = user.profile_picture || 'default-profile.png'; // Utilise une image par défaut si aucune n'est fournie
        img.alt = `${user.username}'s profile picture`;
        img.style.width = '40px'; // Définir la taille de l'image
        img.style.height = '40px';
        img.style.borderRadius = '50%'; // Forme circulaire
        img.style.marginRight = '10px'; // Espacement entre l'image et le texte

        // Ajout du nom d'utilisateur
        const username = document.createElement('span');
        username.textContent = user.username || 'Nom d\'utilisateur inconnu';

        // Ajout de l'image et du nom d'utilisateur dans l'élément <li>
        li.appendChild(img);
        li.appendChild(username);

        // Ajout de l'élément <li> dans la liste <ul>
        usersList.appendChild(li);
    });
}
