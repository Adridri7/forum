import { DisplayUsersAdmin } from "./request.js";

export async function fetchAllUsers() {
    try {
        const response = await fetch('/api/users/fetchAllUsers');
        if (!response.ok) {
            throw new Error('Erreur lors de la récupération des utilisateurs');
        }
        const users = await response.json();

        // Vider les listes uniquement si `list` est `true`
        clearList(admin);
        clearList(modo);
        clearList(user);
        clearList(goat)

        users.forEach((userData) => {
            try {
                DisplayUsersAdmin(userData);
            } catch (error) {
                console.error(`Erreur lors de l'affichage de l'utilisateur ${userData.id}:`, error);
            }
        });

    } catch (error) {
        console.error('Erreur générale:', error);
    }
}

const admin = document.getElementById('admin-list');
const modo = document.getElementById('modo-list');
const user = document.getElementById('members-list');
const goat = document.getElementById('GOAT-list');

// Fonction pour vider les enfants sauf le h3
export function clearList(list) {
    Array.from(list.children).forEach(child => {
        if (child.tagName !== 'H3' && child.tagName !== 'H1') {
            list.removeChild(child);
        }
    });
}