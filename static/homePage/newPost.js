import { createPost } from "./createdPost.js";
import { getUserInfoFromCookie } from "./utils.js";


const addButton = document.getElementById('add-button');
const modalPost = document.getElementById('modal-post');
const userPost = document.getElementById('users-post');
const loginBtn = document.getElementById('login-btn');

let isModal = false;

export function NewPost() {
    CreatedModal();
    const newpost = document.getElementById('created-post');
    newpost.style.display = 'flex';
    modalPost.style.display = 'flex';
    isModal = true;
    const inputField = document.getElementById('message');
    inputField.focus();

    // Ajouter un écouteur d'événement pour fermer le modal lorsqu'un clic se produit
    document.addEventListener('click', closeModal);
}

function CreatedModal() {
    const userInfo = getUserInfoFromCookie();
    const createdPost = document.createElement('div');
    createdPost.classList.add('created-post');
    createdPost.id = 'created-post';

    const postHeader = document.createElement('div');
    postHeader.classList.add('post');
    postHeader.textContent = 'New Post';

    const profilInfo = document.createElement('div')
    profilInfo.classList.add('profile-info');

    const userNameSpan = document.createElement('div');
    userNameSpan.classList.add('user-name');
    userNameSpan.textContent = userInfo.username;

    const profile_picture = document.createElement('div');
    profile_picture.classList.add('profile-picture');

    const image = document.createElement('img');
    image.src = userInfo.profileImageURL;

    profile_picture.appendChild(image)

    profilInfo.appendChild(profile_picture)
    profilInfo.appendChild(userNameSpan)

    const form = document.createElement('form');
    form.classList.add('message-input');
    form.id = 'message-form';

    const inputField = document.createElement('input');
    inputField.type = 'text';
    inputField.id = 'message';
    inputField.placeholder = "What's new ?";

    const sendBox = document.createElement('div');
    sendBox.classList.add('sendBox');

    const postIcon = document.createElement('div');
    postIcon.classList.add('post-icon');

    const imageUpload = document.createElement('button');
    imageUpload.classList.add('image-upload');
    imageUpload.innerHTML = `
    <svg xmlns="http://www.w3.org/2000/svg" height="28px" viewBox="0 -960 960 960" width="28px" fill="currentcolor">
        <path d="M200-120q-33 0-56.5-23.5T120-200v-560q0-33 23.5-56.5T200-840h560q33 0 56.5 23.5T840-760v560q0 33-23.5 56.5T760-120H200Zm0-80h560v-560H200v560Zm40-80h480L570-480 450-320l-90-120-120 160Zm-40 80v-560 560Z" />
    </svg>`;

    const postButton = document.createElement('button');
    postButton.id = 'post-btn';
    postButton.textContent = 'Post';
    postButton.type = 'submit';

    // Ajout de l'élément au formulaire
    form.appendChild(inputField);

    postIcon.appendChild(imageUpload)

    sendBox.appendChild(postIcon)
    sendBox.appendChild(postButton)

    // Ajout des sections à la modal
    createdPost.appendChild(postHeader);
    createdPost.appendChild(profilInfo);
    createdPost.appendChild(form);
    createdPost.appendChild(sendBox)

    // Ajout du formulaire dans le modal
    userPost.appendChild(createdPost);

    postButton.addEventListener('click', createPost);

    inputField.addEventListener('keydown', (event) => {
        if (event.key === 'Enter') {
            event.preventDefault();
            createPost(); // Appelle la fonction pour créer le post
        }
    });
}


// Fonction pour fermer le modal si on clique à l'extérieur de 'created-post'
function closeModal(event) {
    const newpost = document.getElementById('created-post');
    const modalpost = document.getElementById('modal-post');
    // Vérifie si le clic a eu lieu en dehors de 'created-post'
    if (!newpost.contains(event.target) && event.target !== addButton && isModal) {
        newpost.style.display = 'none'; // Ferme le post
        modalpost.style.display = 'none'; // Ferme aussi le fond modal
        isModal = false;

        document.removeEventListener('click', closeModal);
    }
}

loginBtn.addEventListener('click', () => {
    window.location.href = "/authenticate"
})