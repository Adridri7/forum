export function getUserInfoFromCookie() {
    const cookies = document.cookie.split(';');
    let userInfo = {};

    cookies.forEach(cookie => {
        const [name, value] = cookie.trim().split('=');
        if (name === 'UserLogged') {
            const decodedValue = decodeURIComponent(value);
            const parts = decodedValue.split('|');

            if (parts.length >= 5) {
                userInfo = {
                    uuid: removeQuotes(parts[0]),          // UUID
                    username: parts[1],      // Nom d'utilisateur
                    email: parts[2],         // Email
                    role: parts[3],          // Rôle
                    profileImageURL: removeQuotes(parts[4]) // URL de l'image de profil
                };
            }
        }
    });

    return userInfo;
}

function removeQuotes(uuid) {
    return uuid.replace(/"/g, '');
}