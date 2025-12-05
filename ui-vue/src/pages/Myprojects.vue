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
          <li class="itens__lista">
            <router-link to="/myprojects" class="projects__active">My Projects</router-link>
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

    <main class="main__projects">
      <div class="projects__header">
        <h2 class="projects__titulo">Meus Projetos:</h2>
        <button class="btn__create__project" @click="goToCreate">
          <span class="plus__icon">+</span> novo
        </button>
      </div>

      <div class="grid__projects">
        <!-- Projeto Exemplo -->
        <div 
          class="card__project" 
          v-for="project in filteredProjects" 
          :key="project.id"
          :data-project-id="project.id"
        >
          <div class="project__placeholder">
            <img 
              class="img__agents"
              :src="project.image_url || require('@/assets/images/default-agent.png')" 
              :alt="project.name"
              @error="handleImageError"
            >
          </div>
          <div class="project__info">
            <h3 
              class="project__nome" 
              :contenteditable="project.isEditing"
              @blur="updateProjectName(project, $event)"
            >
              {{ project.name }}
            </h3>
            <p 
              class="project__descricao" 
              :contenteditable="project.isEditing"
              @blur="updateProjectDesc(project, $event)"
            >
              {{ project.description }}
            </p>
            <div class="project__meta">
              <span class="project__id">ID: {{ project.id }}</span>
              <span class="project__data">Última edição: {{ project.lastEdited }}</span>
            </div>
          </div>
          <div class="project__actions">
            <button 
              class="btn__edit" 
              title="Editar"
              @click="toggleEdit(project)"
              v-show="!project.isEditing"
            >
              ✏️
            </button>
            <button 
              class="btn__save" 
              title="Salvar"
              @click="saveProject(project)"
              v-show="project.isEditing"
            >
              💾
            </button>
            <button 
              class="btn__delete" 
              title="Deletar"
              @click="deleteProject(project.id)"
            >
              🗑️
            </button>
          </div>
        </div>

        <!-- Card para criar novo projeto -->
        <div class="card__new__project" @click="goToCreate">
          <div class="new__project__content">
            <span class="plus__icon__large">+</span>
            <p>novo</p>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script>
import api from '@/api/api'

export default {
  name: 'MyprojectsPage',
  data() {
    return {
      searchQuery: '',
      projects: [],
      loading: false,
      error: null,
      currentPage: 0,
      pageSize: 20
    }
  },
  computed: {
    filteredProjects() {
      if (!this.searchQuery) return this.projects
      
      const query = this.searchQuery.toLowerCase()
      return this.projects.filter(project => 
        project.name.toLowerCase().includes(query) ||
        project.description.toLowerCase().includes(query) ||
        project.id.toLowerCase().includes(query)
      )
    }
  },
  metaInfo: {
    title: 'AIgents - My Projects',
    meta: [
      { name: 'description', content: 'Manage your AI agent projects on AIgents. Create, edit, and organize your custom artificial intelligence assistants in one place.' },
      { name: 'keywords', content: 'my AI projects, manage AI agents, create AI assistant, custom AI, AI project management' },
      { name: 'author', content: 'AIgents' },
      { name: 'robots', content: 'noindex, nofollow' }
    ]
  },
  mounted() {
    this.loadProjects()
  },
  methods: {
    async loadProjects() {
      this.loading = true
      this.error = null
      
      try {
        const response = await api.post('/agents/my-projects', {
          page: this.currentPage,
          page_size: this.pageSize
        })
        
        // Handle Go Response structure
        if (response.data.status === 200) {
          this.projects = (response.data.data || []).map(project => ({
            id: project.id || project.project_id,
            name: project.name,
            description: project.description,
            image_url: project.image_url,
            lastEdited: project.last_edited || project.updated_at || new Date().toLocaleDateString('pt-BR'),
            isEditing: false
          }))
        } else {
          this.error = response.data.message || 'Failed to load projects'
        }
        
      } catch (err) {
        this.error = err.response?.data?.message || 'Failed to load projects. Please try again.'
        console.error('Error loading projects:', err)
      } finally {
        this.loading = false
      }
    },
    handleSearch() {
      // A busca é feita automaticamente através do computed property filteredProjects
    },
    handleImageError(event) {
      event.target.src = require('@/assets/images/default-agent.png')
    },
    goToCreate() {
      this.$router.push('/create')
    },
    toggleEdit(project) {
      project.isEditing = !project.isEditing
    },
    updateProjectName(project, event) {
      project.name = event.target.textContent
    },
    updateProjectDesc(project, event) {
      project.description = event.target.textContent
    },
    saveProject(project) {
      project.isEditing = false
      project.lastEdited = new Date().toLocaleDateString('pt-BR')
      // Aqui você pode adicionar lógica para salvar no backend
      console.log('Project saved:', project)
    },
    deleteProject(projectId) {
      if (confirm('Tem certeza que deseja deletar este projeto?')) {
        const index = this.projects.findIndex(p => p.id === projectId)
        if (index !== -1) {
          this.projects.splice(index, 1)
        }
      }
    }
  }
}
</script>

<style scoped>
/* Global styles are imported in main.js */
</style>
