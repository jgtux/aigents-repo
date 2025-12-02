// scriptLogin.js - Updated login handler
// Replace the content of: ui/Scripts/scriptLogin.js

document.addEventListener("DOMContentLoaded", () => {
  const form = document.querySelector("form");
  const emailInput = document.querySelector('input[type="email"]');
  const passwordInput = document.querySelector('input[type="password"]');
  const btnLogin = document.getElementById("btnLogin");

  // Prevent form default submission
  form.addEventListener("submit", (e) => {
    e.preventDefault();
  });

  btnLogin.addEventListener("click", async () => {
    const email = emailInput.value.trim();
    const password = passwordInput.value.trim();

    // Validate inputs
    if (!email || !password) {
      alert("Por favor, preencha todos os campos.");
      return;
    }

    // Disable button and show loading state
    btnLogin.disabled = true;
    const originalText = btnLogin.textContent;
    btnLogin.textContent = "Entrando...";

    try {
      // Call the login function from auth.js
      const result = await login(email, password);

      if (result.success) {
        // Success! Redirect to home
        alert("Login realizado com sucesso!");
        window.location.href = "../Index.html/Home.html";
      } else {
        // Show error
        alert(`Erro: ${result.error || "Falha no login"}`);
        btnLogin.disabled = false;
        btnLogin.textContent = originalText;
      }
    } catch (error) {
      console.error("Login error:", error);
      alert("Erro ao conectar ao servidor.");
      btnLogin.disabled = false;
      btnLogin.textContent = originalText;
    }
  });

  // Allow Enter key to submit
  passwordInput.addEventListener("keypress", (e) => {
    if (e.key === "Enter") {
      btnLogin.click();
    }
  });
});
