import { ReportMessage } from "./API_request.js";
import { UserInfo } from "./app.js";
import { deleteComment, deletePost, updateComment, updatePost } from "./messageAction.js";
import { promoteUser } from "./role.js";

export function DisplayMessages(post, section, isComment = false, isNotif = false) {
    const svgLike = `<svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="currentcolor"><path d="M720-120H280v-520l280-280 50 50q7 7 11.5 19t4.5 23v14l-44 174h258q32 0 56 24t24 56v80q0 7-2 15t-4 15L794-168q-9 20-30 34t-44 14Zm-360-80h360l120-280v-80H480l54-220-174 174v406Zm0-406v406-406Zm-80-34v80H160v360h120v80H80v-520h200Z"/></svg>`
    const svgDislike = `<svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="currentcolor"><path d="M240-840h440v520L400-40l-50-50q-7-7-11.5-19t-4.5-23v-14l44-174H120q-32 0-56-24t-24-56v-80q0-7 2-15t4-15l120-282q9-20 30-34t44-14Zm360 80H240L120-480v80h360l-54 220 174-174v-406Zm0 406v-406 406Zm80 34v-80h120v-360H680v-80h200v520H680Z"/></svg>`;
    const displayTimeStamp = post.created_at ? new Date(post.created_at).toLocaleString() : new Date().toLocaleString();
    const svgCertificate = `<svg xmlns="http://www.w3.org/2000/svg" class="certificate" fill="currentcolor" viewBox="0 0 512 512"><!--!Font Awesome Free 6.6.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license/free Copyright 2024 Fonticons, Inc.--><path d="M211 7.3C205 1 196-1.4 187.6 .8s-14.9 8.9-17.1 17.3L154.7 80.6l-62-17.5c-8.4-2.4-17.4 0-23.5 6.1s-8.5 15.1-6.1 23.5l17.5 62L18.1 170.6c-8.4 2.1-15 8.7-17.3 17.1S1 205 7.3 211l46.2 45L7.3 301C1 307-1.4 316 .8 324.4s8.9 14.9 17.3 17.1l62.5 15.8-17.5 62c-2.4 8.4 0 17.4 6.1 23.5s15.1 8.5 23.5 6.1l62-17.5 15.8 62.5c2.1 8.4 8.7 15 17.1 17.3s17.3-.2 23.4-6.4l45-46.2 45 46.2c6.1 6.2 15 8.7 23.4 6.4s14.9-8.9 17.1-17.3l15.8-62.5 62 17.5c8.4 2.4 17.4 0 23.5-6.1s8.5-15.1 6.1-23.5l-17.5-62 62.5-15.8c8.4-2.1 15-8.7 17.3-17.1s-.2-17.4-6.4-23.4l-46.2-45 46.2-45c6.2-6.1 8.7-15 6.4-23.4s-8.9-14.9-17.3-17.1l-62.5-15.8 17.5-62c2.4-8.4 0-17.4-6.1-23.5s-15.1-8.5-23.5-6.1l-62 17.5L341.4 18.1c-2.1-8.4-8.7-15-17.1-17.3S307 1 301 7.3L256 53.5 211 7.3z"/></svg>`
    const messagesList = document.querySelector(`.users-post[data-section="${section}"]`);

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


    const roleDiv = document.createElement('div');
    roleDiv.classList.add('user-role');
    roleDiv.innerHTML = svgCertificate;
    const svgElement = roleDiv.querySelector('.certificate path'); // Sélectionner le chemin du SVG
    if (post.role === "modo") {
        svgElement.setAttribute('fill', '#3c85de');
    } else if (post.role === "admin") {
        svgElement.setAttribute('fill', '#B3001B');
    } else if (post.role === "GOAT") {
        svgElement.setAttribute('fill', '#EEC643')
    }

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

    // Créer le menu dans messageItem
    const menu = document.createElement('ul');
    menu.classList.add('menu');
    menu.style.display = 'none'; // caché par défaut
    const postUuid = isComment ? post.comment_id : post.post_uuid;
    menu.setAttribute('data-post-uuid', postUuid);

    if (UserInfo && (UserInfo.user_uuid || UserInfo.role === "admin" || UserInfo.role === "GOAT" || UserInfo.role === "modo")) {
        const deleteMenuItem = document.createElement('li');
        deleteMenuItem.classList.add('menu-item');
        deleteMenuItem.textContent = 'Delete';
        deleteMenuItem.addEventListener('click', () => {
            if (isComment) {
                deleteComment(post.post_uuid, post.comment_id);
            } else {
                deletePost(post.post_uuid);
            }
            menu.style.display = 'none'; // Fermer le menu après la suppression
        });

        // Si l'utilisateur est le créateur du post, ajouter un bouton d'édition
        if (UserInfo.user_uuid === post.user_uuid) {
            const editButton = document.createElement('li');
            editButton.classList.add('menu-item');
            editButton.textContent = 'Edit';
            editButton.addEventListener('click', () => {
                // Remplacer le contenu par la barre d'input
                messageContent.style.display = 'none';
                editMessageInput.style.display = 'block';
                editMessageInput.focus();
                editMessageInput.value = messageContent.textContent;

                // Gérer la sauvegarde des modifications
                const originalMessage = messageContent.textContent;
                editMessageInput.addEventListener('blur', async () => {
                    const updatedContent = editMessageInput.value;
                    originalMessage = updatedContent;
                    if (isComment) {
                        await updateComment(post.comment_id, updatedContent);
                    } else {
                        await updatePost(post.post_uuid, updatedContent);
                    }
                    messageContent.textContent = updatedContent;
                    messageContent.style.display = 'block';
                    editMessageInput.style.display = 'none';
                });

                // Gérer les touches Entrée et Échappement
                editMessageInput.addEventListener('keydown', (event) => {
                    if (event.key === 'Enter') {
                        editMessageInput.blur(); // Déclenche l'événement de flou pour sauvegarder
                    } else if (event.key === 'Escape') {
                        messageContent.style.display = 'block';
                        editMessageInput.style.display = 'none';
                        editMessageInput.value = originalMessage;
                    }
                });
            });
            menu.appendChild(editButton);
        }

        // Bouton de promotion
        if (post.role != "GOAT" && post.role !== "admin" && post.role !== UserInfo.role) {
            const promoteButton = document.createElement('li');
            promoteButton.classList.add('menu-item');
            promoteButton.textContent = 'Promote';
            promoteButton.addEventListener('click', () => promoteUser(post.user_uuid, "promote"));
            menu.appendChild(promoteButton);
        }



        if (post.role !== "user" && post.role !== "GOAT" && post.role !== UserInfo.role) {
            const demoteButton = document.createElement('li');
            demoteButton.classList.add('menu-item');
            demoteButton.textContent = 'Demote';
            demoteButton.addEventListener('click', () => promoteUser(post.user_uuid, "demote"));
            menu.appendChild(demoteButton)
        }

        const report = document.createElement('li');
        report.className = 'menu-item';
        report.textContent = 'Report';
        report.addEventListener('click', (event) => {
            event.stopPropagation();
            ReportMessage(post.post_uuid);
        });
        menu.appendChild(report);


        menu.appendChild(deleteMenuItem);
    }

    messageHeader.appendChild(userNameSpan);
    messageHeader.appendChild(roleDiv);
    messageHeader.appendChild(timeStampSpan);
    messageHeader.appendChild(menuButton);
    if (menu.children.length > 0) {
        menuButton.appendChild(menuSvg);
        messageHeader.appendChild(menu);
    }

    // Créer le contenu du message
    const messageContent = document.createElement('div');
    messageContent.classList.add('message-content');
    messageContent.textContent = post.content;

    const edit = document.createElement('span');
    edit.classList.add('edit-span');
    edit.textContent = '(edited)';

    // Créer la barre d'input pour update un message
    const editMessageInput = document.createElement('input');
    editMessageInput.classList.add('input-edit-message');
    editMessageInput.type = 'text';
    editMessageInput.value = post.content;
    editMessageInput.style.display = 'none';

    // Créer les boutons de réaction
    const reactionBtnContainer = document.createElement('div');
    reactionBtnContainer.classList.add('reaction-btn');
    if (isComment) {
        reactionBtnContainer.id = "reaction-btn";
    } else {
        reactionBtnContainer.id = "reaction-comment-btn";
    }

    // Bouton Like avec SVG
    const likeButton = document.createElement('button');
    if (isComment) {
        likeButton.id = 'like-comment-btn';
        likeButton.classList.add('like-comment-btn');
    } else {
        likeButton.id = 'like-btn';
        likeButton.classList.add('like-btn');
    }

    likeButton.innerHTML = svgLike;

    const countLike = document.createElement('div');
    countLike.classList.add('like-count');
    countLike.textContent = post.likes;

    likeButton.appendChild(countLike);


    // Bouton Dislike + compteur
    const dislikeButton = document.createElement('button');

    if (isComment) {
        dislikeButton.id = 'dislike-comment-btn';
        dislikeButton.classList.add('dislike-comment-btn');
    } else {
        dislikeButton.id = 'dislike-btn';
        dislikeButton.classList.add('dislike-btn');
    }
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
    if (!isNotif) {
        reactionBtnContainer.appendChild(likeButton);
        reactionBtnContainer.appendChild(dislikeButton);
        reactionBtnContainer.appendChild(commentButton);
    }

    // Ajout des éléments au conteneur principal
    messageContainer.appendChild(messageHeader);
    messageContainer.appendChild(messageContent);
    messageContainer.appendChild(editMessageInput);

    if (post.isUpdated || (post.update_at && post.update_at !== "0001-01-01T00:00:00Z")) {
        messageContainer.appendChild(edit);
    }

    if (post.post_image && post.post_image !== "") {
        const imageMessage = document.createElement('img');
        imageMessage.classList.add('image-message');
        imageMessage.src = post.post_image;
        messageContainer.appendChild(imageMessage);
    }

    messageContainer.appendChild(reactionBtnContainer);
    messageItem.appendChild(profileContainer);
    messageItem.appendChild(messageContainer);

    // Ajout du message dans la liste
    messagesList.appendChild(messageItem);
    messagesList.scrollTop = messagesList.scrollHeight;
}

export function toggleMenu(event, post_uuid) {
    event.stopPropagation();
    const menu = document.querySelector(`.menu[data-post-uuid="${post_uuid}"]`);

    // Fermer tous les autres menus
    const allMenus = document.querySelectorAll('.menu');
    allMenus.forEach(m => {
        if (m !== menu) {
            m.style.display = 'none';
        }
    });

    // Positionner le menu à l'endroit du curseur
    menu.style.top = `${event.clientY}px`;
    menu.style.left = `${event.clientX}px`;

    // Toggle le menu actuel
    menu.style.display = 'block';
}

document.addEventListener('click', (event) => {
    if (!event.target.closest('.menu-btn')) {
        const allMenus = document.querySelectorAll('.menu');
        allMenus.forEach(menu => {
            menu.style.display = 'none';
        });
    }
});