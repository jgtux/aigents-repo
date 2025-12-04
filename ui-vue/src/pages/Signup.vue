<template>
  <div>
    <section class="corpo">    
      <nav class="div__logo">
        <img src="@/assets/images/Logo_grande.svg" alt="Logo da AIgents"><br>
        <h1 class="titulo__logo"><strong>AIgents</strong></h1>
      </nav>
      <nav class="div__formulario">
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
          <input type="submit" value="Sign-up" class="botao__login">
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
export default {
  name: 'SignupPage',
  data() {
    return {
      email: '',
      password: '',
      confirmPassword: ''
    }
  },
  metaInfo: {
    title: 'AIgents Sign-up',
    meta: [
      { name: 'description', content: 'Create your AIgents account and start discovering and building AI agents. Join our artificial intelligence platform today.' },
      { name: 'keywords', content: 'AIgents signup, create account, register AI platform, join AIgents' },
      { name: 'author', content: 'AIgents' },
      { name: 'robots', content: 'noindex, nofollow' },
      { property: 'og:title', content: 'Sign Up - AIgents' },
      { property: 'og:description', content: 'Create your account and start building AI agents' },
      { property: 'og:type', content: 'website' }
    ]
  },
  methods: {
    handleSignup() {
      // Validar se as senhas coincidem
      if (this.password !== this.confirmPassword) {
        alert('Passwords do not match!')
        return
      }
      
      // Implementar lógica de signup aqui
      // Substitua a lógica do scriptSign.js e auth.js aqui
      console.log('Signup attempt:', {
        email: this.email,
        password: this.password
      })
      
      // Exemplo de redirecionamento após signup bem-sucedido:
      // this.$router.push('/login')
    }
  }
}
</script>

<style scoped>
/* Global styles are imported in main.js */
</style>
