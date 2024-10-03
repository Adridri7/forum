export async function getPPFromID(id) {
    var pp = "";

    try {
        const response = await fetch("/api/get-pp", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({user_uuid: id})  // Send the data as JSON
        });

        const data = await response.json();

        if (response.ok) {
            console.log(data);
            pp = data.profile_picture;
        } else {
            alert("Erreur lors de l'inscription : " + data.message);
        }
    } catch (error) {
        console.error("Erreur lors de l'inscription", error.message);
    }

    return pp;
}

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

export function isUserInfoValid() {
    var userInfo = getUserInfoFromCookie();
    var uuidRegex = /([0-9a-f]{8})-([0-9a-f]{4})-([0-9a-f]{4})-([0-9a-f]{4})-([0-9a-f]{12})/;
    var emailRegex = /([0-9A-Za-z]+[\.-]*)+@([0-9A-Za-z]+-*)+.(com|org|fr)/;

    // Is UUID valid?
    if (!uuidRegex.test(userInfo.uuid)) {
        return false;
    }

    // Is email valid?
    if (!emailRegex.test(userInfo.email)) {
        return false;
    }

    // Is role valid?
    switch (userInfo.role) {
        case "", "user", "mod", "admin":
            break;
        default:
            return false;
    }

    return true;
}

function removeQuotes(uuid) {
    return uuid.replace(/"/g, '');
}