// export function getUserInfoFromCookie() {
//     let userInfo = null;

//     if (document.cookie.substring(0, 10) === 'UserLogged') {
//         const parts = document.cookie.substring(11).split('|');

//         if (parts.length < 4) {
//             return null;
//         }

//         userInfo = {
//             uuid: removeQuotes(parts[0]),   // UUID
//             username: parts[1],             // Nom d'utilisateur
//             email: parts[2],                // Email
//             role: removeQuotes(parts[3])    // Rôle
//         };
//     }

//     return userInfo;
// }
// function removeQuotes(uuid) {
//     return uuid.replace(/"/g, '');
// }

export async function getPPFromID(id) {
    var pp = "";

    try {
        const response = await fetch("/api/get-pp", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ user_uuid: id })  // Send the data as JSON
        });

        const data = await response.json();

        if (response.ok) {
            pp = data;
        } else {
            alert("Erreur lors de l'inscription : " + data.message);
        }
    } catch (error) {
        console.error("Erreur lors de l'inscription", error.message);
    }

    return pp;
}

// export function isUserInfoValid() {
//     var userInfo = getUserInfoFromCookie();
//     var uuidRegex = /([0-9a-f]{8})-([0-9a-f]{4})-([0-9a-f]{4})-([0-9a-f]{4})-([0-9a-f]{12})/;
//     var emailRegex = /([0-9A-Za-z]+[\.-]*)+@([0-9A-Za-z]+-*)+.(com|org|fr)/;

//     // Is UUID valid?
//     if (!uuidRegex.test(userInfo.uuid)) {
//         return false;
//     }

//     // Is email valid?
//     if (!emailRegex.test(userInfo.email)) {
//         return false;
//     }

//     // Is role valid?
//     switch (userInfo.role) {
//         case "user":
//         case "mod":
//         case "admin":
//             break;
//         default:
//             return false;
//     }

//     return true;
// }

export function resetUsersPost(section) {
    const usersPost = document.querySelector(`.users-post[data-section="${section}"]`);
    console.log(usersPost)

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