document.getElementById("registerForm").addEventListener("submit", function(event) {
  event.preventDefault(); // Prevent form submission

  var formData = {
    first_name: document.getElementById("firstName").value,
    last_name: document.getElementById("lastName").value,
    email: document.getElementById("email").value,
    password: document.getElementById("password").value,
    phone: document.getElementById("phone").value
  };

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
