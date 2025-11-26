// scriptSign.js - Updated signup handler
// Replace the content of: ui/Scripts/scriptSign.js

document.addEventListener("DOMContentLoaded", () => {
  const form = document.querySelector("form");
  const emailInput = document.querySelector('input[type="email"]');
  const passwordInput = document.querySelectorAll('input[type="password"]')[0];
  const confirmPasswordInput = document.querySelectorAll('input[type="password"]')[1];
  const btnSignup = document.getElementById("btnSignup");

  // Prevent form default submission
  form.addEventListener("submit", async (event) => {
    event.preventDefault();

    const email = emailInput.value.trim();
    const password = passwordInput.value.trim();
    const confirmPassword = confirmPasswordInput.value.trim();

    // Validate inputs
    if (!email || !password || !confirmPassword) {
      alert("Por favor, preencha todos os campos.");
      return;
    }

    if (password !== confirmPassword) {
      alert("As senhas não coincidem.");
      return;
    }

    if (password.length < 8 || password.length > 25) {
      alert("A senha deve ter entre 8 e 25 caracteres.");
      return;
    }

    // Disable button and show loading state
    btnSignup.disabled = true;
    const originalValue = btnSignup.value;
    btnSignup.value = "Cadastrando...";

    try {
      // Call the signup function from auth.js
      const result = await signUp(email, password);

      if (result.success) {
        alert("Cadastro realizado com sucesso!");
        window.location.href = "/Index.html/Login.html";
      } else {
        alert(`Erro ao cadastrar: ${result.error}`);
        btnSignup.disabled = false;
        btnSignup.value = originalValue;
      }
    } catch (err) {
      console.error("Erro na requisição:", err);
      alert("Erro ao conectar ao servidor.");
      btnSignup.disabled = false;
      btnSignup.value = originalValue;
    }
  });
});
