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
    const fileInput = document.getElementById('file-input');
    const removeImageButton = document.getElementById('remove-image');
    const imageContainer = document.querySelector('.image-container'); // Select the container
    const files = event.target.files;

    if (files.length > 0) {
        const file = files[0];
        const reader = new FileReader();

        reader.onload = function (e) {
            fileInput.src = e.target.result;
            fileInput.style.display = 'flex';  // Show the image
            removeImageButton.style.display = 'block';  // Show the remove button
            imageContainer.style.display = 'flex';  // Show the image container
        }

        reader.readAsDataURL(file);

    } else {
        fileInput.src = "";  // Clear the image if no file is chosen
        fileInput.style.display = 'none';  // Hide the image
        removeImageButton.style.display = 'none';  // Hide the remove button
        imageContainer.style.display = 'none';  // Hide the image container
    }
});

// Remove image button logic
document.getElementById('remove-image').addEventListener('click', function () {
    const fileInput = document.getElementById('file-input');
    const removeImageButton = document.getElementById('remove-image');
    const fileInputElement = document.getElementById('profile-image-input');
    const imageContainer = document.querySelector('.image-container');

    // Clear the file input value
    fileInputElement.value = '';
    fileInput.src = '';  // Remove the src attribute of the image
    fileInput.style.display = 'none';  // Hide the image
    removeImageButton.style.display = 'none';  // Hide the remove button
    imageContainer.style.display = 'none';  // Hide the image container
});
