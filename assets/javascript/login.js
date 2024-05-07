function login() {
    var modal = document.getElementById("loginModal");
    modal.style.display = "block";
}

function closeModal() {
    var modal = document.getElementById("loginModal");
    modal.style.display = "none";
}

document.getElementById('loginForm').addEventListener('submit', function (event) {
    event.preventDefault();

    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    fetch('/api/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            username: username,
            password: password,
        }),
    })
        .then(response => response.json())
        .then(data => {
            // Close the modal and update the UI
            closeModal();
            document.getElementById('user-not-logged').style.display = 'none';
            document.getElementById('logged-username').textContent = data.username;
            document.getElementById('user-logged').style.display = 'block';
        })
        .catch((error) => {
            console.error('Error:', error);
        });
});