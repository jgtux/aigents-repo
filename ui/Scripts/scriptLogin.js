document.addEventListener("DOMContentLoaded", () => {
    const btnLogin = document.getElementById("btnLogin");
  
    btnLogin.addEventListener("click", () => {
      const emailInput = document.querySelector('input[type="email"]');
      const passwordInput = document.querySelector('input[type="password"]');
  
      const email = emailInput.value.trim();
      const password = passwordInput.value.trim();
  
      if (email === "" || password === "") {
        alert("Por favor, preencha todos os campos.");
        return;
      }
  
      // Recupera usuário do sessionStorage ao invés de localStorage
      const usuarioCadastrado = JSON.parse(sessionStorage.getItem("usuarioCadastrado"));
  
      if (usuarioCadastrado && email === usuarioCadastrado.email && password === usuarioCadastrado.senha) {
        sessionStorage.setItem("usuarioLogado", "true");
        // Caminho corrigido para a nova estrutura
        window.location.href = "/Index.html/Home.html";
      } else {
        alert("E-mail ou senha inválidos.");
      }
      
    });
  });