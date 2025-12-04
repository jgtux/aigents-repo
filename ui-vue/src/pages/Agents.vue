<template>
  <div>
    <header>
      <nav class="div__logo__pequena">
        <img src="@/assets/images/Logo_pequena.svg" alt="AIgents Logo">
        <h1 class="titulo__logo__pequena"><strong>AIgents</strong></h1>
      </nav>
      <nav aria-label="Menu principal">
        <ul class="header__lista">
          <li>
            <router-link class="itens__lista" to="/">Home</router-link>
          </li>
          <li class="itens__lista">
            <router-link to="/agents" class="agents__active">Agents</router-link>
          </li>
          <li>
            <router-link class="itens__lista" to="/myprojects">My Projects</router-link>
          </li>
          <li>
            <router-link class="itens__lista" to="/about">About</router-link>
          </li>
        </ul>
      </nav>
      <ul class="header__links__lista">
        <div class="search__container">
          <img src="@/assets/images/Lupinha.svg" alt="Search icon" class="search__icon">
          <input 
            type="text" 
            placeholder="Search" 
            class="search__input"
            v-model="searchQuery"
            @input="handleSearch"
          >
        </div>
        <router-link to="/login">
          <img src="@/assets/images/Perfil.svg" alt="Profile Icon">
        </router-link>
      </ul>
    </header>
    
    <!-- Busca por categorias implementação futura
    <nav class="nav__categorias">
      <ul class="lista__categorias">
        <li><button class="categoria__item categoria__active">Programming</button></li>
      </ul>
    </nav>-->
    
    <main class="main__agents">
      <section class="grid__agents">
        <!-- Os cards serão gerados dinamicamente pelo Vue -->
      </section>
    </main>
  </div>
</template>

<script>
export default {
  name: 'AgentsPage',
  data() {
    return {
      searchQuery: '',
      agents: []
    }
  },
  metaInfo: {
    title: 'AIgents - Agents',
    meta: [
      { name: 'description', content: 'Browse and discover AI agents for programming and other categories. Find the perfect artificial intelligence assistant for your needs on AIgents.' },
      { name: 'keywords', content: 'AI agents, programming AI, browse AI assistants, artificial intelligence tools, AI categories' },
      { name: 'author', content: 'AIgents' }
    ],
    htmlAttrs: {
      lang: 'pt-BR'
    }
  },
  methods: {
    handleSearch() {
      // Implementar lógica de busca aqui
      console.log('Searching for:', this.searchQuery)
    },
    loadAgents() {
      // Aqui você pode carregar os agents de uma API ou dados locais
      // Substitua a lógica do agents.js aqui
    }
  },
  mounted() {
    this.loadAgents()
  }
}
</script>

<style scoped>
/* Global styles are imported in main.js */
</style>
