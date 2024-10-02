function IsValidCookie(cookie) {
    if (cookie.length === 0) {
        return false;
    }

    var values = cookie.split("|");

    if (values.length !== 5) {
        return false;
    }

    // Is first value a valid UUID?
    var uuidRegex = /([0-9a-f]{8})-([0-9a-f]{4})-([0-9a-f]{4})-([0-9a-f]{4})-([0-9a-f]{12})/;

    if (!uuidRegex.test(values[0])) {
        return false;
    }

    // Do I really check the username?

    // Is third value a valid email?
    var emailRegex = /([0-9A-Za-z]+[\.-]*)+@([0-9A-Za-z]+-*)+.(com|org|fr)/;

    if (!emailRegex.test(values[2])) {
        return false;
    }

    // Is fourth value a valid role?
    switch (values[3]) {
        case "", "user", "mod", "admin":
            break;
        default:
            return false;
    }

    // Can I even check the profile picture??

    return true;
}