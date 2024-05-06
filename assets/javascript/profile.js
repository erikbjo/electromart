document.addEventListener('DOMContentLoaded', function () {
    // Fetch user profile when the page is loaded
    fetchUserProfile();

    // Add event listener to the form for updating profile
    document.getElementById('profileForm').addEventListener('submit', function (event) {
        event.preventDefault(); // Prevent default form submission behavior
        updateUserProfile();
    });

    // Add event listeners for modal related functionality
    setupModalEventListeners();
});

function setupModalEventListeners() {
    // Add event listener to the delete user button
    document.getElementById('deleteUserBtn').addEventListener('click', function () {
        deleteUser();
    });

    // Add event listener to the change password button
    document.getElementById('changePasswordBtn').addEventListener('click', function () {
        // Display the change password modal
        document.getElementById('changePasswordModal').style.display = "block";
    });

    // Add event listener to the close button in the change password modal
    document.getElementsByClassName('close')[0].addEventListener('click', function () {
        // Hide the change password modal
        document.getElementById('changePasswordModal').style.display = "none";
    });

    // Add event listener to close the change password modal when pressing ESC key
    document.addEventListener('keydown', function (event) {
        if (event.key === 'Escape') {
            // Hide the change password modal
            document.getElementById('changePasswordModal').style.display = "none";
        }
    });

    // Add event listener to the form for changing password
    document.getElementById('changePasswordForm').addEventListener('submit', function (event) {
        event.preventDefault(); // Prevent default form submission behavior
        changePassword();
    });
}


function fetchUserProfile() {
    console.log("Fetching user profile...");
    var xhr = new XMLHttpRequest();
    xhr.open('GET', '/api/profile', true);
    xhr.setRequestHeader('Content-type', 'application/json');
    xhr.onload = function () {
        if (xhr.status === 200) {
            var respJson = JSON.parse(xhr.responseText);
            fillUserProfileForm(respJson);
        } else if (xhr.status === 401) {
            window.location.href = '/login';
        } else {
            showError('Failed to get profile information. Please try again later.');
        }
    };
    xhr.onerror = function () {
        showError('Failed to fetch profile. Please check your internet connection and try again.');
    };
    xhr.send();
}

function fillUserProfileForm(user) {
    console.log("Filling user profile form...");
    document.getElementById('username').value = user.Username;
    document.getElementById('email').value = user.Email;
    document.getElementById('firstName').value = user.FirstName;
    document.getElementById('lastName').value = user.LastName;
    document.getElementById('phone').value = user.Phone;
    if (user.Address && user.Address.Valid) {
        document.getElementById('address').value = user.Address.String;
    }
    if (user.PostCode && user.PostCode.Valid) {
        document.getElementById('postCode').value = user.PostCode.String;
    }
}

function updateUserProfile() {
    var username = document.getElementById('username').value.trim();
    var email = document.getElementById('email').value.trim();
    var firstName = document.getElementById('firstName').value.trim();
    var lastName = document.getElementById('lastName').value.trim();
    var phone = document.getElementById('phone').value.trim();
    var address = document.getElementById('address').value.trim();
    var postCode = document.getElementById('postCode').value.trim();

    if (username === '' || email === '' || firstName === '' || lastName === '' || address === '' || postCode === '') {
        showErrorMessage('Please fill in all required fields.');
        return;
    }

    var data = {
        "Username": username,
        "Email": email,
        "FirstName": firstName,
        "LastName": lastName,
        "Phone": phone,
        "Address": {"String": address, "Valid": true},
        "PostCode": {"String": postCode, "Valid": true}
    };

    var xhr = new XMLHttpRequest();
    xhr.open('PATCH', '/api/profile', true);
    xhr.setRequestHeader('Content-type', 'application/json');
    xhr.onload = function () {
        if (this.status == 200) {
            showSuccessMessage('Profile updated successfully');
            fetchUserProfile();
        } else if (this.status == 400) {
            showErrorMessage('Bad request: The server could not understand the request due to invalid syntax.');
        } else if (this.status == 401) {
            window.location.href = '/login';
        } else if (this.status == 403) {
            showErrorMessage('Forbidden: You are not allowed to update this profile.');
        } else if (this.status == 404) {
            showErrorMessage('Not found: The requested resource could not be found.');
        } else if (this.status == 500) {
            showErrorMessage('Internal Server Error: Failed to update profile. Please try again later.');
        } else {
            showErrorMessage('Unknown error occurred. Status code: ' + this.status);
        }
    };
    xhr.onerror = function () {
        showErrorMessage('Failed to update profile. Please check your internet connection and try again.');
    }
    xhr.send(JSON.stringify(data));
}

function deleteUser() {
    var xhr = new XMLHttpRequest();
    xhr.open('DELETE', '/api/profile', true);
    xhr.setRequestHeader('Content-type', 'application/json');
    xhr.onload = function () {
        if (xhr.status === 204) {
            showSuccessMessage('User deleted successfully');
            window.location.href = '/';
        } else {
            showErrorMessage('Failed to delete user. Please try again later.');
        }
    };
    xhr.onerror = function () {
        showErrorMessage('Failed to delete user. Please check your internet connection and try again.');
    };
    xhr.send();
}

function changePassword() {
    var oldPassword = document.getElementById('oldPassword').value.trim();
    var newPassword = document.getElementById('newPassword').value.trim();

    if (oldPassword === '' || newPassword === '') {
        showErrorMessage('Please fill in all required fields.');
        return;
    }

    var data = {
        "oldPassword": oldPassword,
        "newPassword": newPassword
    };

    var xhr = new XMLHttpRequest();
    xhr.open('PATCH', '/api/change_password', true);
    xhr.setRequestHeader('Content-type', 'application/json');
    xhr.onload = function () {
        if (xhr.status === 200) {
            showSuccessMessage('Password changed successfully');
            document.getElementById('changePasswordModal').style.display = "none";
            location.reload();
        } else if (xhr.status === 400) {
            showErrorMessage('Bad request: The server could not understand the request due to invalid syntax.');
        } else if (xhr.status === 401) {
            showErrorMessage('Unauthorized: Please login to change your password.');
            window.location.href = '/login';
        } else if (xhr.status === 500) {
            showErrorMessage('Internal Server Error: Failed to change password. Please try again later.');
        } else {
            showErrorMessage('Unknown error occurred. Status code: ' + xhr.status);
        }
    };
    xhr.onerror = function () {
        showErrorMessage('Failed to change password. Please check your internet connection and try again.');
    }
    xhr.send(JSON.stringify(data));
}

function showSuccessMessage(message) {
    var successMessageElement = document.createElement('div');
    successMessageElement.classList.add('success-message');
    successMessageElement.textContent = message;

    var messagesContainer = document.getElementById('messages');
    messagesContainer.appendChild(successMessageElement);

    // Show the messages container
    messagesContainer.style.display = 'block';

    setTimeout(function () {
        successMessageElement.remove();
        // Hide the messages container if it's empty
        if (messagesContainer.childElementCount === 0) {
            messagesContainer.style.display = 'none';
        }
    }, 5000);
}

function showErrorMessage(message) {
    var errorMessageElement = document.createElement('div');
    errorMessageElement.classList.add('error-message');
    errorMessageElement.textContent = message;

    var messagesContainer = document.getElementById('messages');
    messagesContainer.appendChild(errorMessageElement);

    // Show the messages container
    messagesContainer.style.display = 'block';

    setTimeout(function () {
        errorMessageElement.remove();
        // Hide the messages container if it's empty
        if (messagesContainer.childElementCount === 0) {
            messagesContainer.style.display = 'none';
        }
    }, 5000);
}

