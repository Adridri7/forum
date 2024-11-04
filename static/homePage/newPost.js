import { UserInfo } from "./app.js";
import { createPost } from "./apiCreatedPost.js";
import { getPPFromID } from "./utils.js";


const addButton = document.getElementById('add-button');
const modalPost = document.getElementById('modal-post');
const userPost = document.querySelector(`.users-post[data-section="home"]`);
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
    userNameSpan.textContent = UserInfo.username;

    const profile_picture = document.createElement('div');
    profile_picture.classList.add('profile-picture');

    const image = document.createElement('img');

    getPPFromID(UserInfo.user_uuid).then(img => { image.src = img });

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

    const embedPreview = document.createElement('img');
    embedPreview.id = 'embed-preview';
    embedPreview.style.maxHeight = '300px';
    embedPreview.style.maxWidth = '200px';

    const removeImg = document.createElement('button');
    removeImg.id = 'remove-image';
    removeImg.innerHTML = '&times;';

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
    const imageUploadInput = document.createElement('input');
    imageUploadInput.type = 'file';
    imageUploadInput.setAttribute('accept', 'image/*');
    imageUploadInput.style.display = 'none';

    // Associer le bouton d'upload au champ input pour l'upload d'image
    imageUpload.addEventListener('click', (event) => {
        event.preventDefault();
        imageUploadInput.click();
    });

    const postButton = document.createElement('button');
    postButton.id = 'post-btn';
    postButton.textContent = 'Post';
    postButton.type = 'submit';
    postButton.style.opacity = '0.3'

    postButton.disabled = true; // Désactiver le bouton au départ

    // Fonction pour vérifier si le bouton doit être activé ou désactivé
    function checkPostButtonState() {
        const messageContent = inputField.value.trim();
        const imageSrc = embedPreview.src;

        if (messageContent !== "" || imageSrc !== "") {
            postButton.disabled = false;
            postButton.style.opacity = '1'
        } else {
            postButton.disabled = true;
            postButton.style.opacity = '0.3'
        }
    }

    // Écouter les changements dans le champ de message
    inputField.addEventListener('input', () => {
        checkPostButtonState();
    });


    imageUploadInput.addEventListener('change', () => {
        const maxFileSize = 20 * 1024 * 1024; // 20 Mo en octets
        const selectedFile = imageUploadInput.files[0];

        if (!selectedFile) {
            return;
        }

        if (selectedFile.size > maxFileSize) {
            alert('Le fichier est trop volumineux. La taille maximale est de 20 Mo.');
            imageUploadInput.value = ''; // Réinitialise l'input pour permettre un nouveau choix
            return;
        }

        createdPost.style.height = '580px';

        const fr = new FileReader();
        fr.onload = () => {
            embedPreview.src = fr.result;
            embedPreview.style.display = 'block';  // Rendre l'aperçu visible à nouveau
            removeImg.style.display = 'block';  // Affiche le bouton de suppression
        };
        fr.readAsDataURL(selectedFile);

        embedPreview.alt = selectedFile.name;  // Mettre à jour l'attribut alt de l'aperçu
    });

    removeImg.addEventListener('click', function (event) {
        event.preventDefault();
        event.stopPropagation();

        // Réinitialiser l'image
        embedPreview.src = '';  // Vide la source de l'image
        embedPreview.alt = '';  // Réinitialise l'attribut alt
        embedPreview.style.display = 'none';  // Masquer l'image
        removeImg.style.display = 'none';  // Masquer le bouton de suppression

        // Réinitialiser l'input de fichier
        imageUploadInput.value = '';  // Vide la sélection d'image
        createdPost.style.height = '260px';  // Réduire la taille de la boîte

        checkPostButtonState();
    });

    // Ajout de l'élément au formulaire
    form.appendChild(inputField);
    form.appendChild(embedPreview);
    form.appendChild(removeImg);

    postIcon.appendChild(imageUpload)

    sendBox.appendChild(postIcon)
    sendBox.appendChild(postButton)

    // Ajout des sections à la modal
    createdPost.appendChild(postHeader);
    createdPost.appendChild(profilInfo);
    createdPost.appendChild(form);
    createdPost.appendChild(sendBox);

    // Ajout du formulaire dans le modal
    userPost.appendChild(createdPost);

    postButton.addEventListener('click', (event) => {

        createPost(event);
    });

    inputField.addEventListener('keydown', (event) => {
        if (event.key === 'Enter') {
            event.preventDefault();
            createPost(event);
        }
    });
}


// Fonction pour fermer le modal si on clique à l'extérieur de 'created-post'
function closeModal(event) {
    const newpost = document.getElementById('created-post');
    const modalpost = document.getElementById('modal-post');

    // Vérifie si l'élément cliqué est la croix ou un enfant du modal
    const isClickInsideModal = newpost.contains(event.target);
    const isClickOnRemoveImg = event.target.id === 'remove-image';  // Ajoute une vérification pour la croix

    if (!isClickInsideModal && event.target !== addButton && isModal && !isClickOnRemoveImg) {
        newpost.style.display = 'none';
        modalpost.style.display = 'none';
        isModal = false;

        document.removeEventListener('click', closeModal);
    }
}
loginBtn.addEventListener('click', () => {
    window.location.href = "/authenticate"
})