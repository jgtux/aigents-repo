<template>
  <div>
    <section class="corpo">    
      <nav class="div__logo">
        <img src="@/assets/images/Logo_grande.svg" alt="Logo da AIgents"><br>
        <h1 class="titulo__logo"><strong>AIgents</strong></h1>
      </nav>
      <nav class="div__formulario">
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
          <button type="submit" class="botao__login">Log-in</button>
        </form>
        <a href="#" class="senha__signup" @click.prevent="handleForgotPassword">Forgot your password?</a>
        <p class="signup__texto">Don't have an account? 
          <router-link to="/signup" class="senha__signup">Sign-up</router-link>
        </p>
      </nav>
    </section>
  </div>
</template>

<script>
export default {
  name: 'LoginPage',
  data() {
    return {
      email: '',
      password: ''
    }
  },
  metaInfo: {
    title: 'AIgents Log-in',
    meta: [
      { name: 'description', content: 'Log in to AIgents and access your AI agents. Manage your artificial intelligence projects and discover new AI assistants.' },
      { name: 'keywords', content: 'AIgents login, sign in, AI platform access, user login' },
      { name: 'author', content: 'AIgents' },
      { name: 'robots', content: 'noindex, nofollow' },
      { property: 'og:title', content: 'Login - AIgents' },
      { property: 'og:description', content: 'Access your AI agents platform' },
      { property: 'og:type', content: 'website' }
    ]
  },
  methods: {
    handleLogin() {
      // Implementar lógica de login aqui
      // Substitua a lógica do scriptLogin.js e auth.js aqui
      console.log('Login attempt:', {
        email: this.email,
        password: this.password
      })
      
      // Exemplo de redirecionamento após login bem-sucedido:
      // this.$router.push('/agents')
    },
    handleForgotPassword() {
      // Implementar lógica de recuperação de senha
      console.log('Forgot password clicked')
    }
  }
}
</script>

<style scoped>
/* Global styles are imported in main.js */
</style>
