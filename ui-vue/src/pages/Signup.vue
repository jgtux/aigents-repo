<template>
  <div>
    <section class="corpo">    
      <nav class="div__logo">
        <img src="@/assets/images/Logo_grande.svg" alt="Logo da AIgents"><br>
        <h1 class="titulo__logo"><strong>AIgents</strong></h1>
      </nav>

      <nav class="div__formulario">

        <!-- Mensagens -->
        <p v-if="errorMsg" class="msg error">{{ errorMsg }}</p>
        <p v-if="successMsg" class="msg success">{{ successMsg }}</p>

        <form class="formulario" @submit.prevent="handleSignup">
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
          <input 
            type="password" 
            placeholder="Confirm your Password" 
            class="inputs"
            v-model="confirmPassword"
            required
          >

          <input type="submit" 
                 :value="loading ? 'Creating...' : 'Sign-up'" 
                 class="botao__login"
                 :disabled="loading">
        </form>

        <p class="login__texto">
          Already have an account?
          <router-link to="/login" class="senha__signup">Log-in</router-link>
        </p>
      </nav>
    </section>
  </div>
</template>

<script>
import api from "@/api/api"; // axios configurado

export default {
  name: "SignupPage",

  data() {
    return {
      email: "",
      password: "",
      confirmPassword: "",
      loading: false,
      errorMsg: "",
      successMsg: "",
    };
  },

  methods: {
    async handleSignup() {
      this.errorMsg = "";
      this.successMsg = "";

      if (this.password !== this.confirmPassword) {
        this.errorMsg = "Passwords do not match.";
        return;
      }

      this.loading = true;

      try {
        await api.post("/auth/create", {
          email: this.email,
          password: this.password,
        });

        this.successMsg = "Account created successfully! Redirecting...";
        
        // Espera 1.5s e vai para login
        setTimeout(() => {
          this.$router.push("/login");
        }, 1500);

      } catch (err) {
        console.error("Signup failed:", err);

        if (err.response) {
          switch (err.response.status) {
            case 400:
              this.errorMsg = "Invalid information. Check the fields.";
              break;

            case 409:
              this.errorMsg = "This email is already registered.";
              break;

            case 500:
              this.errorMsg = "Internal error. Try again later.";
              break;

            default:
              this.errorMsg = "Unable to complete signup.";
          }
        } else {
          this.errorMsg = "Network error. Check your connection.";
        }

      } finally {
        this.loading = false;
      }
    }
  }
};
</script>

<style scoped>
/* Global styles are imported in main.js */
</style>
