import { hideAllSections, UserInfo } from "./app.js";
import { clearList, fetchAllUsers } from "./displayUser.js";
import { DisplayReport, DisplayRequest, initRequestEventListeners } from "./request.js";


let post_uuid_report = null

export async function CreateRequest(user_uuid, content) {
    console.log("content de la request", content, "Et id du post :", post_uuid_report)
    const data = {
        user_uuid: user_uuid,
        content: content,
        post_uuid: post_uuid_report
    };

    try {
        const response = await fetch('/api/requests/createRequest', {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(data)
        });

        // Vérifie si la réponse est OK
        if (!response.ok) {
            const errorMessage = await response.text();
            throw new Error(`Error: ${response.status} - ${errorMessage}`);
        }

        alert("Request created successfully!")
        post_uuid_report = null
    } catch (error) {
        console.error('There was a problem with the fetch operation:', error);
    }
}

export async function FetchAdminRequest() {
    const request = document.querySelector('.users-request[data-section="moderation"]');
    clearList(request);
    try {
        const response = await fetch('/api/requests/fetchRequest', {
            headers: {
                "Content-Type": "application/json"
            },
        });

        if (!response.ok) {
            const errorMessage = await response.text();
            throw new Error(`Error: ${response.status} - ${errorMessage}`);
        }

        const datas = await response.json();
        console.log("yes :", datas)
        for (const data of datas) {
            if (data.post_content !== "") {
                const post = await FetchSinglePosts(data.post_uuid);
                DisplayReport(post[0], data);
            } else {
                DisplayRequest(data);
            }
        }

        fetchAllUsers();
        initRequestEventListeners();
    } catch (error) {
        console.error('There was a problem with the fetch operation:', error);
    }
}

export async function toggleRequestReaction(event, request_uuid) {
    // Vérifier s'il s'agit d'un clic sur un bouton d'approbation ou de rejet
    const button = event.target.closest('.approuve-btn, .reject-btn');
    if (!button) return; // Si le clic n'est pas sur l'un des boutons, on quitte la fonction
    const action = button.classList.contains('approuve-btn') ? 'approuve' : 'reject';

    try {
        const response = await fetch('/api/requests/action', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ request_uuid: request_uuid, action: action }),
            credentials: 'same-origin'
        });

        // Vérifie si la réponse est un succès
        if (!response.ok) {
            // Tente de récupérer l'erreur sous forme de texte ou JSON
            const errorMessage = await response.text();
            throw new Error(`Erreur lors de la mise à jour de la réaction: ${errorMessage}`);
        }

        const responseMessage = await response.text();

        alert(`Request processed :${response.status} - ${responseMessage} !`)
        MarkAsReadRequest(request_uuid)

    } catch (error) {
        // Affiche l'erreur complète ici
        console.error(error);
    }

}

async function MarkAsReadRequest(request_uuid) {

    try {
        const response = await fetch('/api/requests/isreadRequest', {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ request_uuid })
        });

        // Vérifie si la réponse est OK
        if (!response.ok) {
            const errorMessage = await response.text();
            throw new Error(`${response.status} - ${errorMessage}`);
        }

    } catch (error) {
        console.error('There was a problem with the fetch operation:', error);
    }
}

export function ReportMessage(post_uuid) {
    const requestSection = document.getElementById('request-section');
    hideAllSections();
    requestSection.style.display = 'block';
    FetchHistoryRequest();
    post_uuid_report = post_uuid
    // Initialisation des écouteurs d'événements
    initReportEventListeners();
}

export function handleKeyDown(event) {
    console.log("Suivi de Keydown l'id du post :", post_uuid_report)
    if (event.key === 'Enter') {
        event.stopPropagation();
        // Vérifie si la valeur n'est pas vide avant d'envoyer la requête
        if (event.target.value.trim() !== '') {
            CreateRequest(UserInfo.user_uuid, event.target.value);
            event.target.value = ''; // Réinitialiser le champ après l'envoi
        }
    }
}

export function handleSendClick(event) {
    console.log("Suivi de HandleClick l'id du post :", post_uuid_report)

    event.stopPropagation();
    const inputField = document.getElementById('request-input');
    // Vérifie si la valeur n'est pas vide avant d'envoyer la requête
    if (inputField.value.trim() !== '') {
        CreateRequest(UserInfo.user_uuid, inputField.value);
        inputField.value = ''; // Réinitialiser le champ après l'envoi
    }
}

export function initReportEventListeners() {
    const inputField = document.getElementById('request-input');
    const sendButton = document.getElementById('request-btn');
    inputField.focus();


    // Retirer d'abord les écouteurs précédents pour éviter les doublons
    inputField.removeEventListener('keydown', (event) => handleKeyDown(event));
    sendButton.removeEventListener('click', (event) => handleSendClick(event));

    // Ajouter les écouteurs d'événements
    inputField.addEventListener('keydown', (event) => handleKeyDown(event));
    sendButton.addEventListener('click', (event) => handleSendClick(event));
}

async function FetchSinglePosts(post_uuid) {
    console.log(post_uuid)
    try {
        const response = await fetch('/api/post/getPostDetails', {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ post_uuid })
        });

        // Vérifie si la réponse est OK
        if (!response.ok) {
            const errorMessage = await response.text();
            throw new Error(`${response.status} - ${errorMessage}`);
        }

        const data = response.json();

        return data;

    } catch (error) {
        console.error('There was a problem with the fetch operation:', error);
    }
}

export async function FetchHistoryRequest() {
    const historyRequest = document.getElementById('history-users-request');
    historyRequest.innerHTML = '';
    const user_uuid = UserInfo.user_uuid
    try {
        const response = await fetch('/api/requests/historyRequest', {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ user_uuid })
        });

        // Vérifie si la réponse est OK
        if (!response.ok) {
            const errorMessage = await response.text();
            throw new Error(`${response.status} - ${errorMessage}`);
        }

        const datas = await response.json();
        console.log("data recu : ", datas)

        for (const data of datas) {
            if (data.post_content !== "") {
                const post = await FetchSinglePosts(data.post_uuid);
                DisplayReport(post[0], data, true);
            } else {
                DisplayRequest(data, true);
            }
        }

    } catch (error) {
        console.error('There was a problem with the fetch operation:', error);
    }
}