import { deletePost } from "./app.js";

export function DisplayMessages(post) {
    const displayTimeStamp = post.created_at ? new Date(post.created_at).toLocaleString() : new Date().toLocaleString();
    const messagesList = document.getElementById('users-post');

    // Créer l'élément 'li' pour le message
    const messageItem = document.createElement('li');
    messageItem.classList.add('message-item');
    messageItem.setAttribute('post_uuid', post.post_uuid);

    // Créer le conteneur de l'image de profil
    const profileContainer = document.createElement('div');
    profileContainer.classList.add('profil-picture');
    const profilePicture = document.createElement('img');
    profilePicture.src = post.profile_picture || 'default-profile-picture.jpg';
    profilePicture.alt = 'user-profil-picture';
    profileContainer.appendChild(profilePicture);

    // Créer le conteneur de message
    const messageContainer = document.createElement('div');
    messageContainer.classList.add('message-container');

    // Créer l'en-tête du message
    const messageHeader = document.createElement('div');
    messageHeader.classList.add('message-header');
    const userNameSpan = document.createElement('span');
    userNameSpan.classList.add('username');
    userNameSpan.textContent = post.username;

    const timeStampSpan = document.createElement('span');
    timeStampSpan.classList.add('timestamp');
    timeStampSpan.textContent = displayTimeStamp;

    // Bouton SVG (trois points)
    const menuButton = document.createElement('button');
    menuButton.classList.add('menu-btn');
    const menuSvg = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
    menuSvg.setAttribute('xmlns', 'http://www.w3.org/2000/svg');
    menuSvg.setAttribute('height', '28px');
    menuSvg.setAttribute('width', '28px');
    menuSvg.setAttribute('viewBox', '0 -960 960 960');
    menuSvg.setAttribute('fill', 'currentcolor');
    const menuPath = document.createElementNS('http://www.w3.org/2000/svg', 'path');
    menuPath.setAttribute('d', 'M240-400q-33 0-56.5-23.5T160-480q0-33 23.5-56.5T240-560q33 0 56.5 23.5T320-480q0 33-23.5 56.5T240-400Zm240 0q-33 0-56.5-23.5T400-480q0-33 23.5-56.5T480-560q33 0 56.5 23.5T560-480q0 33-23.5 56.5T480-400Zm240 0q-33 0-56.5-23.5T640-480q0-33 23.5-56.5T720-560q33 0 56.5 23.5T800-480q0 33-23.5 56.5T720-400Z');
    menuSvg.appendChild(menuPath);
    menuButton.appendChild(menuSvg);

    menuButton.addEventListener('click', function (event) {
        toggleMenu(event, post.post_uuid);
    });

    // Créer le menu dans messageItem
    const menu = document.createElement('ul');
    menu.classList.add('menu');
    menu.style.display = 'none'; // caché par défaut
    menu.setAttribute('data-post-uuid', post.post_uuid); // Ajouter un attribut data pour identifier le menu

    const deleteMenuItem = document.createElement('li');
    deleteMenuItem.classList.add('menu-item');
    deleteMenuItem.textContent = 'Delete';
    deleteMenuItem.addEventListener('click', () => {
        deletePost(post.post_uuid); // Appelle la fonction de suppression
        menu.style.display = 'none'; // Fermer le menu après la suppression
    });

    menu.appendChild(deleteMenuItem);
    messageHeader.appendChild(userNameSpan);
    messageHeader.appendChild(timeStampSpan);
    messageHeader.appendChild(menuButton);
    messageHeader.appendChild(menu); // Ajout du menu à l'en-tête

    // Créer le contenu du message
    const messageContent = document.createElement('div');
    messageContent.classList.add('message-content');
    messageContent.textContent = post.content;

    // Créer les boutons de réaction
    const reactionBtnContainer = document.createElement('div');
    reactionBtnContainer.classList.add('reaction-btn');

    // Bouton Like avec SVG
    const likeButton = document.createElement('button');
    const likeSvg = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
    likeSvg.setAttribute('xmlns', 'http://www.w3.org/2000/svg');
    likeSvg.setAttribute('height', '24px');
    likeSvg.setAttribute('width', '24px');
    likeSvg.setAttribute('viewBox', '0 -960 960 960');
    likeSvg.setAttribute('fill', 'currentcolor');
    const likePath = document.createElementNS('http://www.w3.org/2000/svg', 'path');
    likePath.setAttribute('d', 'm480-120-58-52q-101-91-167-157T150-447.5Q111-500 95.5-544T80-634q0-94 63-157t157-63q52 0 99 22t81 62q34-40 81-62t99-22q94 0 157 63t63 157q0 46-15.5 90T810-447.5Q771-395 705-329T538-172l-58 52Zm0-108q96-86 158-147.5t98-107q36-45.5 50-81t14-70.5q0-60-40-100t-100-40q-47 0-87 26.5T518-680h-76q-15-41-55-67.5T300-774q-60 0-100 40t-40 100q0 35 14 70.5t50 81q36 45.5 98 107T480-228Zm0-273Z');
    likeSvg.appendChild(likePath);
    likeButton.appendChild(likeSvg);

    // Bouton Comment avec SVG
    const commentButton = document.createElement('button');
    commentButton.classList.add('comment-btn');
    const commentSvg = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
    commentSvg.setAttribute('xmlns', 'http://www.w3.org/2000/svg');
    commentSvg.setAttribute('height', '24px');
    commentSvg.setAttribute('width', '24px');
    commentSvg.setAttribute('viewBox', '0 -960 960 960');
    commentSvg.setAttribute('fill', 'currentcolor');
    const commentPath = document.createElementNS('http://www.w3.org/2000/svg', 'path');
    commentPath.setAttribute('d', 'M440-400h80v-120h120v-80H520v-120h-80v120H320v80h120v120ZM80-80v-720q0-33 23.5-56.5T160-880h640q33 0 56.5 23.5T880-800v480q0 33-23.5 56.5T800-240H240L80-80Zm126-240h594v-480H160v525l46-45Zm-46 0v-480 480Z');
    commentSvg.appendChild(commentPath);
    commentButton.appendChild(commentSvg);

    reactionBtnContainer.appendChild(likeButton);
    reactionBtnContainer.appendChild(commentButton);

    // Ajout des éléments au conteneur principal
    messageContainer.appendChild(messageHeader);
    messageContainer.appendChild(messageContent);
    messageContainer.appendChild(reactionBtnContainer);
    messageItem.appendChild(profileContainer);
    messageItem.appendChild(messageContainer);

    // Ajout du message dans la liste
    messagesList.appendChild(messageItem);
    messagesList.scrollTop = messagesList.scrollHeight;
}

function toggleMenu(event, post) {
    event.stopPropagation();
    const menu = event.currentTarget.nextElementSibling;

    // Fermer tous les autres menus
    const allMenus = document.querySelectorAll('.menu');
    allMenus.forEach(m => {
        if (m !== menu) {
            m.style.display = 'none';
        }
    });

    // Toggle le menu actuel
    menu.style.display = 'block'
}

document.addEventListener('click', (event) => {
    if (!event.target.closest('.message-item')) {
        const allMenus = document.querySelectorAll('.menu');
        allMenus.forEach(menu => {
            menu.style.display = 'none';
        });
    }
});