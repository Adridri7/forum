export async function promoteUser(user_uuid, action) {
    try {
        const response = await fetch(`/api/post/updateUserRole`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ user_uuid: user_uuid, action })
        });

        // Vérifie si la réponse est un succès (status 2xx)
        if (!response.ok) {
            const errorData = await response.json(); // Récupère les données d'erreur du serveur
            alert(`Error while updating user role: ${errorData.message || response.statusText}`);
        } else {
            alert("Success updating role");
        }
    } catch (error) {
        alert("Network error while updating user role: " + error.message);
    }
}