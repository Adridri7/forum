import { isUserInfoValid } from "../utils.js";

if (isUserInfoValid()) {
    window.location.href = "/";
}

document.addEventListener('DOMContentLoaded', () => {
    const loginBtn = document.getElementById('login-btn');
    const signupBtn = document.getElementById('signup-btn');
    const loginForm = document.getElementById('login-form');
    const signupForm = document.getElementById('signup-form');
    const homeBtn = document.getElementById('back-home');


    signupBtn.addEventListener('click', () => {
        loginForm.classList.add('inactive');
        signupForm.classList.remove('inactive');
        signupForm.classList.add('active');
        signupBtn.classList.add('active');
        loginBtn.classList.remove('active');
    });

    loginBtn.addEventListener('click', () => {
        loginForm.classList.remove('inactive');
        signupForm.classList.remove('active');
        signupForm.classList.add('inactive');
        loginBtn.classList.add('active');
        signupBtn.classList.remove('active');
    });

    homeBtn.addEventListener('click', () => {
        window.location.href = '/';
    });


    document.getElementById('submit-login').addEventListener('click', async (event) => {
        event.preventDefault();
        console.log("Bouton cliqué !");

        const email = document.getElementById("email").value;
        const password = document.getElementById("password").value;

        const data = {
            email: email,
            password: password
        }

        try {
            const response = await fetch("/api/login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(data)
            });

            if (response.ok) {
                window.location.href = "/";
            } else {
                const error = await response.json();
                alert("Erreur lors du login: " + error.message);
            }
        } catch (error) {
            console.error("Erreur lors du login", error);
        }
    });

    document.getElementById('submit-register').addEventListener('click', async (event) => {
        event.preventDefault();
        const username = document.getElementById("new-username").value;
        const password = document.getElementById("new-password").value;
        const email = document.getElementById('new-email').value;

        const data = {
            username: username,
            password: password,
            email: email
        }

        try {
            const response = await fetch("/api/registration", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(data)
            });
            //console.log(data)
            if (response.ok) {
                window.location.href = "/"
            } else {
                const error = await response.json();
                alert("Erreur lors du login", + error.message);
            }
        } catch (error) {
            console.error("Erreur lors du login", error.message)
        }
    });

    document.getElementById('profile-image-input').addEventListener('change', function (event) {
        const fileNameElement = document.getElementById('file-name');
        const files = event.target.files;

        if (files.length > 0) {
            fileNameElement.textContent = files[0].name;
        } else {
            fileNameElement.textContent = "Aucun fichier choisi";
        }
    });

    document.getElementById('google-btn').addEventListener('click', () => {
        window.location.href = "/api/google_login"
    })
});