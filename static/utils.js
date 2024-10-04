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
            pp = data;
        } else {
            alert("Erreur lors de l'inscription : " + data.message);
        }
    } catch (error) {
        console.error("Erreur lors de l'inscription", error.message);
    }

    return pp;
}

export function getUserInfoFromCookie() {
    let userInfo = {};

    if (document.cookie.substring(0, 10) === 'UserLogged') {
        const parts = document.cookie.substring(11).split('|');

        if (parts.length < 4) {
            return {};
        }

        userInfo = {
            uuid: removeQuotes(parts[0]),          // UUID
            username: parts[1],      // Nom d'utilisateur
            email: parts[2],         // Email
            role: parts[3]           // Rôle
        };
    }

    return userInfo;
}

export function isUserInfoValid() {
    var userInfo = getUserInfoFromCookie();
    var uuidRegex = /([0-9a-f]{8})-([0-9a-f]{4})-([0-9a-f]{4})-([0-9a-f]{4})-([0-9a-f]{12})/;
    var emailRegex = /([0-9A-Za-z]+[\.-]*)+@([0-9A-Za-z]+-*)+.(com|org|fr)/;

    console.log(userInfo);

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
        case "user":
        case "mod":
        case "admin":
            break;
        default:
            return false;
    }

    return true;
}

function removeQuotes(uuid) {
    return uuid.replace(/"/g, '');
}