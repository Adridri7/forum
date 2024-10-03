export async function handleLogout(event) {
    event.preventDefault();

    // Demande de confirmation avant de procéder à la déconnexion
    const confirmLogout = window.confirm("Êtes-vous sûr de vouloir vous déconnecter ?");

    if (confirmLogout) {
        try {
            const response = await fetch('/logout', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'same-origin'
            });

            if (response.ok) {
                window.location.href = '/'; // Redirige vers la page d'accueil après la déconnexion
            } else {
                // Gérer les erreurs de réponse ici
                console.error('Erreur lors de la déconnexion');
            }
        } catch (error) {
            // Gérer les erreurs de réseau ici
            console.error('Erreur réseau :', error);
        }
    }
}