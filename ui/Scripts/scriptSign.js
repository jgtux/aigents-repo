const apiUrl = "http://localhost:8080/api/v1"; // ✅ Always include protocol

document.addEventListener("DOMContentLoaded", () => {
  const form = document.querySelector("form");
  const emailInput = document.querySelector('input[type="email"]');
  const passwordInput = document.querySelectorAll('input[type="password"]')[0];
  const confirmPasswordInput = document.querySelectorAll('input[type="password"]')[1];

  form.addEventListener("submit", async (event) => {
    event.preventDefault();

    const email = emailInput.value.trim();
    const password = passwordInput.value.trim();
    const confirmPassword = confirmPasswordInput.value.trim();

    // Validations
    if (!email || !password || !confirmPassword) {
      alert("Por favor, preencha todos os campos.");
      return;
    }

    if (password !== confirmPassword) {
      alert("As senhas não coincidem.");
      return;
    }

    try {
      const response = await fetch(`${apiUrl}/auth/create`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, password }),
      });

      if (response.status === 201) {
        alert("Cadastro realizado com sucesso!");
        window.location.href = "/Index.html/Login.html";
      } else {
        const errorText = await response.text();
        alert(`Erro ao cadastrar: ${errorText || response.statusText}`);
      }
    } catch (err) {
      console.error("Erro na requisição:", err);
      alert("Erro ao conectar ao servidor.");
    }
  });
});
