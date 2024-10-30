import { toggleRequestReaction } from "./API_request.js";
import { UserInfo } from "./app.js";
import { promoteUser } from "./role.js";

export function DisplayRequest(request, historyRequestContainer = false) {
    const displayTimeStamp = request.created_at ? new Date(request.created_at).toLocaleString() : new Date().toLocaleString();
    const svgApprove = `<svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960"
        width="24px" fill="#3c85de">
        <path
            d="m424-296 282-282-56-56-226 226-114-114-56 56 170 170Zm56 216q-83 0-156-31.5T197-197q-54-54-85.5-127T80-480q0-83 31.5-156T197-763q54-54 127-85.5T480-880q83 0 156 31.5T763-763q54 54 85.5 127T880-480q0 83-31.5 156T763-197q-54 54-127 85.5T480-80Zm0-80q134 0 227-93t93-227q0-134-93-227t-227-93q-134 0-227 93t-93 227q0 134 93 227t227 93Zm0-320Z" />
    </svg>`

    const svgReject = `<svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960"
        width="24px" fill="#B3001B">
        <path
            d="M480-80q-83 0-156-31.5T197-197q-54-54-85.5-127T80-480q0-83 31.5-156T197-763q54-54 127-85.5T480-880q83 0 156 31.5T763-763q54 54 85.5 127T880-480q0 83-31.5 156T763-197q-54 54-127 85.5T480-80Zm0-80q54 0 104-17.5t92-50.5L228-676q-33 42-50.5 92T160-480q0 134 93 227t227 93Zm252-124q33-42 50.5-92T800-480q0-134-93-227t-227-93q-54 0-104 17.5T284-732l448 448Z" />
    </svg>`

    const usersRequest = document.querySelector('.users-request');
    const historyRequest = document.getElementById('history-users-request');


    const container = document.createElement('div');
    container.className = 'container';

    const profilePictureRequest = document.createElement('div');
    profilePictureRequest.className = 'profile-picture-request';

    const profileImg = document.createElement('img');
    profileImg.src = request.profile_picture;

    profilePictureRequest.appendChild(profileImg);

    const messageContainer = document.createElement('div');
    messageContainer.className = 'message-container';
    messageContainer.setAttribute('request-uuid', request.request_uuid)

    const headerMessage = document.createElement('div');
    headerMessage.className = 'header-message';

    const usernameSpan = document.createElement('span');
    usernameSpan.className = 'username';
    usernameSpan.textContent = request.username;

    const timestampSpan = document.createElement('span');
    timestampSpan.className = 'timestamp';
    timestampSpan.textContent = displayTimeStamp;

    headerMessage.appendChild(usernameSpan);
    headerMessage.appendChild(timestampSpan);

    const messageContent = document.createElement('div');
    messageContent.className = 'message-content';

    const messageSpan = document.createElement('span');
    messageSpan.textContent = request.content;

    messageContent.appendChild(messageSpan);


    messageContainer.appendChild(headerMessage);
    messageContainer.appendChild(messageContent);

    if (!historyRequestContainer) {
        const reactionBtnRequest = document.createElement('div');
        reactionBtnRequest.className = 'reaction-btn-request';

        const approveButton = document.createElement('button');
        approveButton.className = 'approuve-btn'
        approveButton.innerHTML = svgApprove

        const rejectButton = document.createElement('button');
        rejectButton.className = 'reject-btn'
        rejectButton.innerHTML = svgReject;

        reactionBtnRequest.appendChild(approveButton);
        reactionBtnRequest.appendChild(rejectButton);
        messageContainer.appendChild(reactionBtnRequest);
    }


    container.appendChild(profilePictureRequest);
    container.appendChild(messageContainer);

    if (historyRequestContainer && request.action !== "") {

        const action = document.createElement('span');
        action.className = 'action'
        action.innerHTML = `Your request : <strong style="color: #3c85de;">${request.content}</strong> has been "${request.action}"`;
        messageContainer.appendChild(action)
    }

    if (historyRequestContainer) {
        historyRequest.appendChild(container)
    } else {
        usersRequest.appendChild(container)
    }
}

export function DisplayUsersAdmin(users) {
    const admin = document.getElementById('admin-list');
    const modo = document.getElementById('modo-list');
    const user = document.getElementById('members-list');
    const goat = document.getElementById('GOAT-list');

    const memberList = document.createElement('li');

    const members = document.createElement('div');
    members.classList.add('members');

    const image = document.createElement('img')
    image.src = users.profile_picture || "https://koreus.cdn.li/media/201404/90-insolite-34.jpg";
    image.alt = `${users.username}'s profile picture`

    const username = document.createElement('span');
    username.classList.add('username')
    username.textContent = users.username;

    const menu = document.createElement('ul');
    menu.classList.add('menu');
    menu.style.display = 'none'; // caché par défaut
    menu.setAttribute('data-user-uuid', users.user_uuid);

    if (users.role !== "GOAT" && users.role !== "admin" && users.role !== UserInfo.role) {
        const promoteButton = document.createElement('li');
        promoteButton.classList.add('menu-item');
        promoteButton.textContent = 'Promote';
        promoteButton.addEventListener('click', () => promoteUser(users.user_uuid, "promote"));
        menu.appendChild(promoteButton);
    }

    if ((users.role !== "user" && users.role !== "GOAT" && users.role !== UserInfo.role)) {
        const demoteButton = document.createElement('li');
        demoteButton.classList.add('menu-item');
        demoteButton.textContent = 'Demote';
        demoteButton.addEventListener('click', () => promoteUser(users.user_uuid, "demote"));
        menu.appendChild(demoteButton);
    }

    members.appendChild(image)
    members.appendChild(username)
    members.appendChild(menu)
    memberList.appendChild(members)

    if (users.role === "GOAT") {
        goat.appendChild(memberList)
    }
    else if (users.role === "admin") {
        admin.appendChild(memberList)
    } else if (users.role === "modo") {
        modo.appendChild(memberList)
    } else {
        user.appendChild(memberList)
    }
    members.addEventListener('click', function (event) {
        toggleMenus(event, users.user_uuid);
    });
}

function toggleMenus(event, user_uuid) {
    console.log('ici')
    event.stopPropagation();
    const menu = document.querySelector(`.menu[data-user-uuid="${user_uuid}"]`);
    console.log(menu)

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

export function initRequestEventListeners() {
    const reactionBtnRequest = document.querySelectorAll('.reaction-btn-request');

    reactionBtnRequest.forEach(container => {
        // Ajout de l'écouteur d'événement pour chaque conteneur
        container.addEventListener('click', (event) => {
            // Récupérer le conteneur parent qui a l'attribut request-uuid
            const messageContainer = container.closest('.message-container');
            const request_uuid = messageContainer.getAttribute('request-uuid');
            console.log("ID :", request_uuid);
            toggleRequestReaction(event, request_uuid);
        });
    });
}

export function DisplayReport(report, request, historyRequestContainer = false) {
    const displayTimeStamp = report.created_at ? new Date(report.created_at).toLocaleString() : new Date().toLocaleString();
    const svgApprove = `<svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960"
        width="24px" fill="#3c85de">
        <path
            d="m424-296 282-282-56-56-226 226-114-114-56 56 170 170Zm56 216q-83 0-156-31.5T197-197q-54-54-85.5-127T80-480q0-83 31.5-156T197-763q54-54 127-85.5T480-880q83 0 156 31.5T763-763q54 54 85.5 127T880-480q0 83-31.5 156T763-197q-54 54-127 85.5T480-80Zm0-80q134 0 227-93t93-227q0-134-93-227t-227-93q-134 0-227 93t-93 227q0 134 93 227t227 93Zm0-320Z" />
    </svg>`

    const svgReject = `<svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960"
        width="24px" fill="#B3001B">
        <path
            d="M480-80q-83 0-156-31.5T197-197q-54-54-85.5-127T80-480q0-83 31.5-156T197-763q54-54 127-85.5T480-880q83 0 156 31.5T763-763q54 54 85.5 127T880-480q0 83-31.5 156T763-197q-54 54-127 85.5T480-80Zm0-80q54 0 104-17.5t92-50.5L228-676q-33 42-50.5 92T160-480q0 134 93 227t227 93Zm252-124q33-42 50.5-92T800-480q0-134-93-227t-227-93q-54 0-104 17.5T284-732l448 448Z" />
    </svg>`

    const usersRequest = document.querySelector('.users-request');
    const historyRequest = document.getElementById('history-users-request');

    const container = document.createElement('div');
    container.className = 'container';

    const profilePictureRequest = document.createElement('div');
    profilePictureRequest.className = 'profile-picture-request';

    const profileImg = document.createElement('img');
    profileImg.src = report.profile_picture;

    profilePictureRequest.appendChild(profileImg);

    const messageContainer = document.createElement('div');
    messageContainer.className = 'message-container';
    messageContainer.setAttribute('request-uuid', request.request_uuid)

    const headerMessage = document.createElement('div');
    headerMessage.className = 'header-message';

    const usernameSpan = document.createElement('span');
    usernameSpan.className = 'username';
    usernameSpan.textContent = report.username;

    const timestampSpan = document.createElement('span');
    timestampSpan.className = 'timestamp';
    timestampSpan.textContent = displayTimeStamp;

    headerMessage.appendChild(usernameSpan);
    headerMessage.appendChild(timestampSpan);

    const messageContent = document.createElement('div');
    messageContent.className = 'message-content';

    const messageSpan = document.createElement('span');
    messageSpan.textContent = report.content;

    messageContent.appendChild(messageSpan);


    messageContainer.appendChild(headerMessage);
    messageContainer.appendChild(messageContent);

    if (report.post_image && report.post_image !== "") {
        const imageMessage = document.createElement('img');
        imageMessage.classList.add('image-message');
        imageMessage.src = report.post_image;
        messageContainer.appendChild(imageMessage);
    }

    if (!historyRequestContainer) {
        const reactionBtnRequest = document.createElement('div');
        reactionBtnRequest.className = 'reaction-btn-request';

        const approveButton = document.createElement('button');
        approveButton.className = 'approuve-btn'
        approveButton.innerHTML = svgApprove

        const rejectButton = document.createElement('button');
        rejectButton.className = 'reject-btn'
        rejectButton.innerHTML = svgReject;

        reactionBtnRequest.appendChild(approveButton);
        reactionBtnRequest.appendChild(rejectButton);

        messageContainer.appendChild(reactionBtnRequest);

    }

    // Crée le span pour le message
    const messageReport = document.createElement('span');
    messageReport.innerHTML = `<strong style="color: #3c85de;">${request.username}</strong> report the following message for : "${request.content}"`;


    container.appendChild(profilePictureRequest);
    container.appendChild(messageContainer);

    if (historyRequestContainer && request.action !== "") {

        const action = document.createElement('span');
        action.className = 'action'
        action.innerHTML = `Your request : <strong style="color: #3c85de;">${request.content}</strong> has been "${request.action}"`;
        messageContainer.appendChild(action)
    }


    if (historyRequestContainer) {
        historyRequest.appendChild(messageReport)
        historyRequest.appendChild(container)
    } else {
        usersRequest.appendChild(messageReport)
        usersRequest.appendChild(container)
    }

}

// Appeler fetch post pour avoir les infos du post puis utiliser ces infos pour contruire le post et en même temps mettre le message de report