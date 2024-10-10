export async function fetchAllUsers() {
    try {
        const response = await fetch('/api/users/fetchAllUsers');
        if (!response.ok) {
            throw new Error('Erreur lors de la récupération des utilisateurs');
        }
        const users = await response.json();

        displayUsers(users);
        console.log(users);
    } catch (error) {
        console.error(error);
    }
}

// Fonction pour afficher les utilisateurs dans la liste
function displayUsers(users) {
    const usersList = document.getElementById('user-list');
    usersList.innerHTML = '';

    users.forEach(user => {
        const li = document.createElement('li');
        li.style.display = 'flex';
        li.style.alignItems = 'center';
        li.style.marginBottom = '10px';

        // Ajout de l'image de profil
        const img = document.createElement('img');
        img.src = user.profile_picture || "https://koreus.cdn.li/media/201404/90-insolite-34.jpg";
        img.alt = `${user.username}'s profile picture`;
        img.style.width = '40px';
        img.style.height = '40px';
        img.style.borderRadius = '50%';
        img.style.marginRight = '10px';

        const username = document.createElement('span');
        username.textContent = user.username || 'Nom d\'utilisateur inconnu';

        li.appendChild(img);
        li.appendChild(username);

        usersList.appendChild(li);
    });
}
