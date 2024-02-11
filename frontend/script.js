function signup() {
    // Get user input
    var username = document.getElementById("username").value;
    var email = document.getElementById("email").value;
    var password = document.getElementById("password").value;

    // Create JSON object with user data
    var userData = {
        "username": username,
        "email": email,
        "password": password
    };

    // Make a POST request to the backend /signup endpoint
    fetch('/signup', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(userData),
    })
    .then(response => response.json())
    .then(data => {
        // Handle the response from the server
        console.log(data);
        alert('User registered successfully!');
        window.location.replace('login.html');
    })
    .catch((error) => {
        console.error('Error:', error);
        alert('An error occurred while registering the user.');
    });
}

function login() {
    // Get user input for login
    var loginEmail = document.getElementById("loginEmail").value;
    var loginPassword = document.getElementById("loginPassword").value;

    // Create JSON object with user credentials
    var loginData = {
        "email": loginEmail,
        "password": loginPassword
    };

    // Make a POST request to the backend /login endpoint
    fetch('/login', {
        method: 'POST'{
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(loginData),
    })
    .then(response => response.json())
    .then(data => {
        // Handle the response from the server
        console.log(data);

        // Save token in local storage
        localStorage.setItem('token', data.token);

        alert('Login successful!');
        // Redirect to another page (replace 'target_page.html' with your actual target page)
        // window.location.replace('target.html');
    })
    .catch((error) => {
        console.error('Error:', error);
        alert('Login failed. Please check your email and password.');
    });
}
