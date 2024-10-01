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