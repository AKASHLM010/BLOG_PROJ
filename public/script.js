// Register Page
const registerForm = document.getElementById('registerForm');
const registerMessage = document.getElementById('message');

registerForm.addEventListener('submit', (e) => {
    e.preventDefault();

    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    const data = {
        username,
        password
    };

    fetch('/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
        .then(response => response.json())
        .then(result => {
            if (result.error) {
                registerMessage.textContent = result.error;
            } else {
                registerMessage.textContent = 'Registered successfully!';
            }
        })
        .catch(error => {
            console.error('Error:', error);
        });
});

// Login Page
const loginForm = document.getElementById('loginForm');
const loginMessage = document.getElementById('message');

loginForm.addEventListener('submit', (e) => {
    e.preventDefault();

    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    const data = {
        username,
        password
    };

    fetch('/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
        .then(response => response.json())
        .then(result => {
            if (result.error) {
                loginMessage.textContent = result.error;
            } else {
                loginMessage.textContent = 'Logged in successfully!';
            }
        })
        .catch(error => {
            console.error('Error:', error);
        });
});
