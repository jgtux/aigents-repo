<template>
  <div>
    <header>
      <nav class="div__logo__pequena">
        <img src="@/assets/images/Logo_pequena.svg" alt="AIgents Logo">
        <h1 class="titulo__logo__pequena"><strong>AIgents</strong></h1>
      </nav>
      <nav aria-label="Menu principal">
        <ul class="header__lista">
          <li class="itens__lista">
            <router-link class="home" to="/">Home</router-link>
          </li>
          <li>
            <router-link class="itens__lista" to="/agents">Agents</router-link>
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
        <router-link to="/login">
          <img src="@/assets/images/Perfil.svg" alt="Profile">
        </router-link>
      </ul>
    </header>
    
    <main class="corpo__landing">
      <div class="div__texto">
        <h1 class="titulo__home">ALL <strong> AI AGENTS</strong> IN ONE PLACE</h1>
        <p class="paragrafo__home">Discover and create your AI agents with our functional tool</p>
      </div>
      <ul class="lista__creating__discover">
        <li class="itens__dis__cre">
          <img src="@/assets/images/Robosinho.svg" alt="Icone de robosinho">
          <router-link to="/agents" class="discover">Discover AIs</router-link>
        </li>
        <li class="itens__dis__cre">
          <img src="@/assets/images/Engrenagens.svg" alt="icone de ingrenagens">
          <router-link to="/myprojects" class="discover">Start Creating</router-link>
        </li>
      </ul>
    </main>
  </div>
</template>

<script>
export default {
  name: 'HomePage',
  metaInfo: {
    title: 'AIgents Home',
    meta: [
      { name: 'description', content: 'AIgents - All AI agents in one place. Discover and create your AI agents with our functional tool. Access multiple artificial intelligences on a single platform.' },
      { name: 'keywords', content: 'AI agents, artificial intelligence, AI platform, create AI agents, discover AI, multiple AI, AIgents' },
      { name: 'author', content: 'AIgents' },
      { property: 'og:title', content: 'AIgents - All AI Agents in One Place' },
      { property: 'og:description', content: 'Discover and create your AI agents with our functional tool' },
      { property: 'og:type', content: 'website' }
    ]
  }
}
</script>

<style scoped>
/* Global styles are imported in main.js or index.html */
</style>
