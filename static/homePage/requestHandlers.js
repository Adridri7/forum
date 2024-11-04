import { CreateRequest, FetchHistoryRequest, MarkAsReadRequest } from "./apiRequests.js";
import { hideAllSections } from "./section.js";

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
