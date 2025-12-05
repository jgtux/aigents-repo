<template>
  <div>
    <section class="corpo">    
      <nav class="div__logo">
        <img src="@/assets/images/Logo_grande.svg" alt="Logo da AIgents"><br>
        <h1 class="titulo__logo"><strong>AIgents</strong></h1>
      </nav>

      <nav class="div__formulario">
        
        <!-- MENSAGEM DE ERRO (usa msg + error, apenas isso) -->
        <p v-if="errorMsg" class="msg error">{{ errorMsg }}</p>

        <form class="formulario" @submit.prevent="handleLogin">
          
          <input 
            type="email" 
            placeholder="E-mail" 
            class="inputs"
            v-model="email"
            required
          >

          <input 
            type="password" 
            placeholder="Password" 
            class="inputs"
            v-model="password"
            required
          >

          <button type="submit" class="botao__login" :disabled="loading">
            {{ loading ? "Logging in..." : "Log-in" }}
          </button>

        </form>

        <a href="#" class="senha__signup" @click.prevent="handleForgotPassword">
          Forgot your password?
        </a>

        <p class="signup__texto">
          Don't have an account? 
          <router-link to="/signup" class="senha__signup">Sign-up</router-link>
        </p>

      </nav>
    </section>
  </div>
</template>

<script>
import api from "@/api/api";

export default {
  name: 'LoginPage',

  data() {
    return {
      email: "",
      password: "",
      loading: false,
      errorMsg: ""
    };
  },

  methods: {
    async handleLogin() {
      this.errorMsg = "";
      this.loading = true;

      try {
        await api.post("/auth/login", {
          email: this.email,
          password: this.password
        });

        this.$router.push("/agents");

      } catch (err) {
        if (err.response) {
          switch (err.response.status) {
            case 400:
            case 401:
              this.errorMsg = "Invalid email or password.";
              break;

            case 429:
              this.errorMsg = "Too many attempts. Try again later.";
              break;

            default:
              this.errorMsg = "Unable to login. Please try again.";
          }
        } else {
          this.errorMsg = "Network error. Check your connection.";
        }

      } finally {
        this.loading = false;
      }
    },

    handleForgotPassword() {
      console.log("Forgot password clicked");
    }
  }
};
</script>


<style scoped>
.error-msg {
  color: #ff4d4d;
  margin-top: 10px;
  font-weight: bold;
  text-align: center;
}
</style>
