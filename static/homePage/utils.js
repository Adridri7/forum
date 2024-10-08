export function getUserInfoFromCookie() {
    const cookies = document.cookie.split(';');
    let userInfo = null; // Commencer par null

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

export function resetUsersPost() {
    const usersPost = document.getElementById('users-post');

    // Vide le contenu de users-post
    usersPost.innerHTML = '';

    // Réinitialise les styles d'origine
    usersPost.style.display = 'block';
    usersPost.style.gridTemplateColumns = '';
    usersPost.style.gridAutoRows = '';
    usersPost.style.rowGap = '';
    usersPost.style.columnGap = '';
    usersPost.style.border = '1px solid var(--border-color)';

}

export function getRandomColor() {
    // Génère une couleur aléatoire au format hexadécimal
    const letters = '0123456789ABCDEF';
    let color = '#';
    for (let i = 0; i < 6; i++) {
        color += letters[Math.floor(Math.random() * 16)];
    }
    return color;
}

export function startGradientAnimation(element) {
    // Définit les couleurs de dégradé
    const color1 = getRandomColor();
    const color2 = getRandomColor();
    const color3 = getRandomColor();

    // Applique le dégradé
    element.style.background = `linear-gradient(90deg, ${color1}, ${color2}, ${color3})`;
    element.style.backgroundSize = 'cover';
}

// Applique la couleur à chaque élément ayant la classe 'categorie'
document.querySelectorAll('.categorie').forEach(categorie => {
    const randomColor = startGradientAnimation();
    categorie.style.backgroundColor = randomColor;
});