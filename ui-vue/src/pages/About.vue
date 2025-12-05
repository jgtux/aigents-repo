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
          <li>
            <router-link class="itens__lista" to="/agents">Agents</router-link>
          </li>
          <li>
            <router-link class="itens__lista" to="/myprojects">My Projects</router-link>
          </li>
          <li class="itens__lista">
            <router-link class="About" to="/about">About</router-link>
          </li>
        </ul>
      </nav>

      <ul class="header__links__lista">
        <router-link to="/login">
          <img src="@/assets/images/Perfil.svg" alt="Profile Icon">
        </router-link>
      </ul>
    </header>

    <main class="corpo__about">
      <section class="textos_about">
        <h2>AIgents — Uniting the Minds of Artificial Intelligence</h2>
        <p>AIgents was created with the purpose of bringing together the world's most advanced artificial intelligences in one single place. Our goal is to offer an integrated experience where users can access different assistants and compare their responses, making the most of what each one has to offer.</p>
        <p>More than just a platform, AIgents is an ecosystem of collaboration between artificial intelligences. We aim to simplify access to technology and expand the possibilities of interaction, making the use of AI more efficient, accessible, and intelligent. Here, all AI voices come together — working side by side to empower what you can create.</p>
      </section>
      <figure class="about__img">
        <img class="imagem__about" src="@/assets/images/Enzo.jpg" alt="CEO da empresa pensado">
        <figcaption>Our pretty CEO</figcaption>
      </figure>
    </main>
  </div>
</template>

<script>
export default {
  name: 'AboutPage',
  metaInfo: {
    title: 'AIgents About Us',
    meta: [
      { name: 'description', content: 'Learn about AIgents - A platform uniting the world\'s most advanced artificial intelligences. Compare responses and access different AI assistants in one place.' },
      { name: 'keywords', content: 'about AIgents, AI platform, artificial intelligence ecosystem, multiple AI assistants, AI comparison' },
      { name: 'author', content: 'AIgents' }
    ]
  }
}
</script>

<style scoped>
/* Global styles are imported in main.js */
</style>
