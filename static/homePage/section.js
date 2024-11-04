import { fetchPersonnalComment, fetchPersonnalPosts, fetchPersonnalResponse } from "./dashboard.js";
import { initReportEventListeners } from "./requestHandlers.js";
import { getPPFromID } from "./utils.js";

// Sélection des liens et des sections
const homeSection = document.getElementById('home-section');
const dashboardSection = document.getElementById('dashboard-section');
const searchSection = document.getElementById('search-section');
const notificationsSection = document.getElementById('notifications-section');
const trendingSection = document.getElementById('trending-section');
const requestSection = document.getElementById('request-section');
const moderationSection = document.getElementById('moderation-section');

// Masquer toutes les sections
export function hideAllSections() {
    homeSection.style.display = 'none';
    dashboardSection.style.display = 'none';
    searchSection.style.display = 'none';
    notificationsSection.style.display = 'none';
    trendingSection.style.display = 'none';
    requestSection.style.display = 'none';
    moderationSection.style.display = 'none';
}

// Afficher une section spécifique
export function showSection(section) {
    hideAllSections();
    section.style.display = 'block';
}

// Gestion des sections standards
export function initSectionEvents(UserInfo) {
    document.getElementById('home-link').addEventListener('click', () => showSection(homeSection));
    document.getElementById('search-link').addEventListener('click', () => showSection(searchSection));


    if (UserInfo) {
        document.getElementById('request-link').addEventListener('click', () => {
            showSection(requestSection);
            initReportEventListeners();
        });

        document.getElementById('notifications-link').addEventListener('click', () => {
            if (!UserInfo) {
                alert("You must be logged in to see notifications.");
                return;
            }
            showSection(notificationsSection);
        });

        document.getElementById('moderation-link').addEventListener('click', () => showSection(moderationSection));
    }
    document.getElementById('trend-link').addEventListener('click', () => showSection(trendingSection));

}

// Afficher et configurer le dashboard
export function showDashboard(UserInfo) {
    if (!UserInfo) {
        alert("You must be logged in to access the dashboard.");
        return;
    }

    hideAllSections();
    dashboardSection.style.display = 'flex';

    const currentActiveItem = document.querySelector('#nav-bar li.active');
    if (currentActiveItem) {
        const activeId = currentActiveItem.id;

        // Fetch en fonction de l'élément actif
        switch (activeId) {
            case 'personnal-post':
                fetchPersonnalPosts();
                break;
            case 'personnal-comment':
                fetchPersonnalComment();
                break;
            case 'personnal-reaction':
                fetchPersonnalResponse();
                break;
            default:
                break;
        }

        const profilPicture = document.getElementById('profile-picture');
        getPPFromID(UserInfo.user_uuid).then(img => { profilPicture.src = img });
    }

    const navItems = document.querySelectorAll('#nav-bar li');
    navItems.forEach(item => {
        item.addEventListener('click', () => {
            navItems.forEach(nav => nav.classList.remove('active'));
            item.classList.add('active');
        });
    });
}