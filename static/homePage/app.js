// Sélectionnez les éléments
const toggleButton = document.getElementById('toggle-menu-btn');
const sidebar = document.getElementById('sidebar');
const darkModeToggle = document.getElementById('dark-mode-toggle');

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

// Événement pour le menu
toggleButton.addEventListener('click', () => {
    sidebar.classList.toggle('close');
});


async function fetchPosts() {
    const messagesList = document.getElementById('users-post');
    messagesList.innerHTML = '<p>Loading...</p>'; // Show loading state
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
            posts.forEach(post => {
                DisplayMessages(post.post_uuid, post.user_uuid, post.content, post.created_at);
            });
        }
    } catch (error) {
        messagesList.innerHTML = '<p>Error loading posts. Please try again.</p>';
        console.error(error);
    }
}

function DisplayMessages(id, username, content, timestamp) {

    const displayTimeStamp = timestamp ? new Date(timestamp).toLocaleString() : new Date().toLocaleString();


    const messagesList = document.getElementById('users-post');

    const messageItem = document.createElement('div');
    messageItem.classList.add('message-item');
    messageItem.setAttribute('data-id', id);

    const messageHeader = document.createElement('div');
    messageHeader.classList.add('message-header');

    const userNameSpan = document.createElement('span');
    userNameSpan.classList.add('username');
    userNameSpan.textContent = username;

    const timeStampSpan = document.createElement('span');
    timeStampSpan.classList.add('timestamp');
    timeStampSpan.textContent = displayTimeStamp;

    messageHeader.appendChild(userNameSpan);
    messageHeader.appendChild(timeStampSpan);

    const messageContent = document.createElement('div');
    messageContent.classList.add('message-content');
    messageContent.textContent = content;

    messageItem.appendChild(messageHeader);
    messageItem.appendChild(messageContent);

    messagesList.appendChild(messageItem);
    messagesList.scrollTop = messagesList.scrollHeight;
}

document.addEventListener("DOMContentLoaded", () => {
    fetchPosts();
});