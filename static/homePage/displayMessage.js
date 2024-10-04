import { deletePost } from "./app.js";

export function DisplayMessages(post, isComment = false) {

    const svgLike = `<svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="currentcolor"><path d="M720-120H280v-520l280-280 50 50q7 7 11.5 19t4.5 23v14l-44 174h258q32 0 56 24t24 56v80q0 7-2 15t-4 15L794-168q-9 20-30 34t-44 14Zm-360-80h360l120-280v-80H480l54-220-174 174v406Zm0-406v406-406Zm-80-34v80H160v360h120v80H80v-520h200Z"/></svg>`
    const svgDislike = `<svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="currentcolor"><path d="M240-840h440v520L400-40l-50-50q-7-7-11.5-19t-4.5-23v-14l44-174H120q-32 0-56-24t-24-56v-80q0-7 2-15t4-15l120-282q9-20 30-34t44-14Zm360 80H240L120-480v80h360l-54 220 174-174v-406Zm0 406v-406 406Zm80 34v-80h120v-360H680v-80h200v520H680Z"/></svg>`;
    const displayTimeStamp = post.created_at ? new Date(post.created_at).toLocaleString() : new Date().toLocaleString();
    const messagesList = document.getElementById('users-post');

    // Créer l'élément 'li' pour le message
    const messageItem = document.createElement('li');
    messageItem.classList.add('message-item');

    if (isComment) {
        messageItem.setAttribute('post_uuid', post.comment_id);
    } else {
        messageItem.setAttribute('post_uuid', post.post_uuid);
    }

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
    likeButton.classList.add('like-btn');
    likeButton.id = 'like-btn';
    likeButton.innerHTML = svgLike;

    const countLike = document.createElement('div');
    countLike.classList.add('like-count');
    countLike.textContent = post.likes;

    likeButton.appendChild(countLike);


    // Bouton Dislike + compteur
    const dislikeButton = document.createElement('button');
    dislikeButton.classList.add('dislike-btn');
    dislikeButton.id = 'dislike-btn';
    dislikeButton.innerHTML = svgDislike;

    const countDislike = document.createElement('div');
    countDislike.classList.add('dislike-count');
    countDislike.textContent = post.dislikes;

    dislikeButton.appendChild(countDislike);


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

    const commentCount = document.createElement('div');
    commentCount.classList.add('comment-count');
    commentCount.textContent = post.comment_count;

    commentButton.appendChild(commentCount);


    // Ajout des éléments à reaction
    reactionBtnContainer.appendChild(likeButton);
    reactionBtnContainer.appendChild(dislikeButton);
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

export function toggleMenu(event, post) {
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
    if (!event.target.closest('.menu-btn')) {
        const allMenus = document.querySelectorAll('.menu');
        allMenus.forEach(menu => {
            menu.style.display = 'none';
        });
    }
});