<template>
  <div>
    <header>
      <nav class="div__logo__pequena" aria-label="Logo e navegação principal">
        <img src="@/assets/images/Logo_pequena.svg" alt="Logo AIgents">
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
          <li class="itens__lista">
            <router-link class="projects__active" to="/myprojects" aria-current="page">My Projects</router-link>
          </li>
          <li>
            <router-link class="itens__lista" to="/about">About</router-link>
          </li>
        </ul>
      </nav>

      <aside class="header__links__lista" aria-label="Busca e perfil">
        <search class="search__container">
          <img src="@/assets/images/Lupinha.svg" alt="" class="search__icon" aria-hidden="true">
          <input 
            type="search" 
            placeholder="Search" 
            class="search__input" 
            aria-label="Buscar"
            v-model="searchQuery"
          >
        </search>
        <router-link to="/login" aria-label="Perfil do usuário">
          <img src="@/assets/images/Perfil.svg" alt="Ícone de perfil">
        </router-link>
      </aside>
    </header>

    <main class="main__criar">
      <section class="criar__container">
        <h2 class="visually-hidden">Criar Novo Agente de IA</h2>
        
        <form class="form__criar" @submit.prevent="handleCreateAgent" aria-labelledby="form-title">
          <fieldset class="form__section">
            <legend class="visually-hidden">Informações básicas do agente</legend>
            
            <article class="form__group__criar">
              <label for="nomeAgente">Nome</label>
              <input 
                type="text" 
                id="nomeAgente" 
                class="input__criar" 
                required 
                placeholder="Nome do seu agente"
                aria-required="true"
                v-model="agent.name"
              >
            </article>

            <article class="form__group__criar">
              <label for="categoriaAgente">Categoria</label>
              <select 
                id="categoriaAgente" 
                class="input__criar" 
                required 
                aria-required="true"
                v-model="agent.category"
              >
                <option value="">Selecione uma categoria</option>
                <option value="programing">Programing</option>
                <option value="images">Images</option>
                <option value="videos">Videos</option>
                <option value="texts">Texts</option>
                <option value="music">Music</option>
                <option value="other">Other</option>
              </select>
            </article>

            <article class="form__group__criar form__group__image">
              <label for="urlImagem">URL da Imagem</label>
              <input 
                type="url" 
                id="urlImagem" 
                class="input__criar" 
                placeholder="https://exemplo.com/imagem.jpg"
                aria-describedby="url-helper"
                v-model="agent.imageUrl"
                @input="handleImagePreview"
              >
              <small id="url-helper" class="visually-hidden">Cole a URL completa da imagem do seu agente</small>
              
              <figure 
                class="preview__container" 
                v-show="showPreview"
              >
                <img 
                  :src="agent.imageUrl" 
                  class="image__preview" 
                  alt="Preview da imagem do agente"
                  @error="handleImageError"
                >
              </figure>
            </article>
          </fieldset>

          <fieldset class="form__section form__section__description">
            <legend class="visually-hidden">Descrição detalhada do agente</legend>
            
            <article class="form__group__criar form__group__full">
              <label for="descricaoAgente">Descrição</label>
              <textarea 
                id="descricaoAgente" 
                class="textarea__criar" 
                required 
                placeholder="Escreva um pouco sobre seu agente"
                rows="8"
                aria-required="true"
                aria-describedby="desc-helper"
                v-model="agent.description"
              ></textarea>
              <small id="desc-helper" class="visually-hidden">Descreva as funcionalidades e características do seu agente de IA</small>
            </article>

            <button type="submit" class="btn__criar" aria-label="Criar agente de IA">Criar</button>
          </fieldset>
        </form>
      </section>
    </main>
  </div>
</template>

<script>
export default {
  name: 'CreatePage',
  data() {
    return {
      searchQuery: '',
      showPreview: false,
      agent: {
        name: '',
        category: '',
        imageUrl: '',
        description: ''
      }
    }
  },
  metaInfo: {
    title: 'Criar Agente de IA - AIgents',
    meta: [
      { name: 'description', content: 'Create your personalized AI agent on AIgents platform. Configure name, category, image and description of your intelligent assistant.' },
      { name: 'keywords', content: 'create AI agent, virtual assistant, artificial intelligence, custom AI, AIgents creator' },
      { name: 'author', content: 'AIgents' }
    ],
    htmlAttrs: {
      lang: 'pt-BR'
    }
  },
  methods: {
    handleImagePreview() {
      // Mostrar preview se a URL não estiver vazia e parecer válida
      if (this.agent.imageUrl && this.agent.imageUrl.trim() !== '') {
        this.showPreview = true
      } else {
        this.showPreview = false
      }
    },
    handleImageError() {
      // Esconder preview se a imagem falhar ao carregar
      this.showPreview = false
      console.warn('Failed to load image preview')
    },
    handleCreateAgent() {
      // Validar dados antes de enviar
      if (!this.agent.name || !this.agent.category || !this.agent.description) {
        alert('Por favor, preencha todos os campos obrigatórios')
        return
      }

      // Implementar lógica de criação do agente aqui
      // Substitua a lógica do criar.js aqui
      console.log('Creating agent:', this.agent)

      // Exemplo: Salvar no backend e redirecionar
      // await this.saveAgent(this.agent)
      // this.$router.push('/myprojects')

      // Por enquanto, apenas mostrar confirmação
      alert(`Agente "${this.agent.name}" criado com sucesso!`)
      
      // Resetar formulário
      this.resetForm()
    },
    resetForm() {
      this.agent = {
        name: '',
        category: '',
        imageUrl: '',
        description: ''
      }
      this.showPreview = false
    }
  }
}
</script>

<style scoped>
/* Global styles are imported in main.js */
</style>
