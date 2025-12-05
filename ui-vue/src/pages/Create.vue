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
                :disabled="loading"
              >
            </article>

            <article class="form__group__criar">
              <label for="categoriaAgente">Categoria</label>
              <select 
                id="categoriaAgente" 
                class="input__criar" 
                required 
                aria-required="true"
                v-model="agent.categoryId"
                :disabled="loading || loadingCategories"
              >
                <option value="">{{ loadingCategories ? 'Carregando...' : 'Selecione uma categoria' }}</option>
                <option 
                  v-for="category in categories" 
                  :key="category.category_id" 
                  :value="category.category_id"
                >
                  {{ category.category_name }}
                </option>
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
                :disabled="loading"
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
                :disabled="loading"
              ></textarea>
              <small id="desc-helper" class="visually-hidden">Descreva as funcionalidades e características do seu agente de IA</small>
            </article>

            <button 
              type="submit" 
              class="btn__criar" 
              aria-label="Criar agente de IA"
              :disabled="loading"
            >
              {{ loading ? 'Criando...' : 'Criar' }}
            </button>
          </fieldset>
        </form>

        <!-- Error Display -->
        <div v-if="error" class="msg error">
          {{ error }}
        </div>
      </section>
    </main>
  </div>
</template>

<script>
import api from '@/api/api'

export default {
  name: 'CreatePage',
  data() {
    return {
      searchQuery: '',
      showPreview: false,
      loading: false,
      loadingCategories: false,
      error: null,
      categories: [],
      agent: {
        name: '',
        categoryId: '',
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
    async loadCategories() {
      this.loadingCategories = true
      this.error = null

      try {
        const response = await api.get('/agents/categories')
        
        if (response.data.status === 200) {
          this.categories = response.data.data || []
        } else {
          this.error = response.data.message || 'Failed to load categories'
        }
      } catch (err) {
        this.error = err.response?.data?.message || 'Failed to load categories. Please try again.'
        console.error('Error loading categories:', err)
      } finally {
        this.loadingCategories = false
      }
    },

    handleImagePreview() {
      if (this.agent.imageUrl && this.agent.imageUrl.trim() !== '') {
        this.showPreview = true
      } else {
        this.showPreview = false
      }
    },

    handleImageError() {
      this.showPreview = false
      console.warn('Failed to load image preview')
    },

    async handleCreateAgent() {
      // Validar dados antes de enviar
      if (!this.agent.name || !this.agent.categoryId || !this.agent.description) {
        this.error = 'Por favor, preencha todos os campos obrigatórios'
        return
      }

      this.loading = true
      this.error = null

      try {
        const payload = {
            name: this.agent.name,
            description: this.agent.description,
            image_url: this.agent.imageUrl || null,
            category_id: parseInt(this.agent.categoryId),
        }

        const response = await api.post('/agents/create', payload)

        if (response.data.status === 201) {
          alert(`Agente "${this.agent.name}" criado com sucesso!`)
          this.resetForm()
          // Redirecionar para my projects
          this.$router.push('/myprojects')
        } else {
          this.error = response.data.message || 'Failed to create agent'
        }
      } catch (err) {
        this.error = err.response?.data?.message || 'Failed to create agent. Please try again.'
        console.error('Error creating agent:', err)
      } finally {
        this.loading = false
      }
    },

    resetForm() {
      this.agent = {
        name: '',
        categoryId: '',
        imageUrl: '',
        description: ''
      }
      this.showPreview = false
      this.error = null
    }
  },

  mounted() {
    this.loadCategories()
  }
}
</script>

<style scoped>
/* Global styles are imported in main.js */
</style>
