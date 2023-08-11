document.getElementById("registerForm").addEventListener("submit", function(event) {
  event.preventDefault(); // Prevent form submission

  var formData = {
    first_name: document.getElementById("firstName").value,
    last_name: document.getElementById("lastName").value,
    email: document.getElementById("email").value,
    password: document.getElementById("password").value,
    phone: document.getElementById("phone").value
  };

  if (!validatePassword(formData.password)) {
    alert('Password must be at least 6 characters long.');
    return;
  }

  if (!validateEmail(formData.email)) {
    alert('Please enter a valid email address.');
    return;
  }

  if (!validatePhoneNumber(formData.phone)) {
    alert('Please enter a valid phone number.');
    return;
  }

  fetch("/register", {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(formData)
  })
  .then(response => {
    if (response.status === 409) {
      throw new Error("Email or Phone number already registered");
    } else if (!response.ok) {
      throw new Error("Registration failed");
    }
    return response.json();
  })
  .then(data => {
    window.alert("Registered successfully");
    window.location.href = "/login";
  })
  .catch(error => {
    console.error("Error:", error);
    window.alert(error.message);
  });
});

function validatePassword(password) {
  return password.length >= 6;
}

function validateEmail(email) {
  // Basic email validation regex (you can use a more complex regex for stricter validation)
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
}

function validatePhoneNumber(phone) {
  // Basic phone number validation regex (you can use a more complex regex for stricter validation)
  const phoneRegex = /^\d{10}$/;
  return phoneRegex.test(phone);
}