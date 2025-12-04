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
        <button class="btn__create__project" @click="openModal">
          <span class="plus__icon">+</span> novo
        </button>
      </div>

      <div class="grid__projects">
        <!-- Projeto Exemplo 1 -->
        <div 
          class="card__project" 
          v-for="project in filteredProjects" 
          :key="project.id"
          :data-project-id="project.id"
        >
          <div class="project__placeholder">IA</div>
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
        <div class="card__new__project" @click="openModal">
          <div class="new__project__content">
            <span class="plus__icon__large">+</span>
            <p>novo</p>
          </div>
        </div>
      </div>
    </main>

    <!-- Modal para criar novo projeto -->
    <div class="modal" v-show="showModal" @click.self="closeModal">
      <div class="modal__content">
        <span class="modal__close" @click="closeModal">&times;</span>
        <h2>Criar Novo Projeto</h2>
        <form @submit.prevent="createProject">
          <div class="form__group">
            <label for="projectName">Nome do Projeto:</label>
            <input type="text" id="projectName" v-model="newProject.name" required>
          </div>
          <div class="form__group">
            <label for="projectDesc">Descrição:</label>
            <textarea id="projectDesc" rows="4" v-model="newProject.description" required></textarea>
          </div>
          <button type="submit" class="btn__submit">Criar Projeto</button>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'MyprojectsPage',
  data() {
    return {
      searchQuery: '',
      showModal: false,
      newProject: {
        name: '',
        description: ''
      },
      projects: [
        {
          id: 'PRJ001',
          name: 'Agent Name',
          description: 'Brief description of the AI agent',
          lastEdited: '24/10/2025',
          isEditing: false
        }
      ]
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
  methods: {
    handleSearch() {
      // A busca é feita automaticamente através do computed property filteredProjects
    },
    openModal() {
      this.showModal = true
    },
    closeModal() {
      this.showModal = false
      this.newProject = { name: '', description: '' }
    },
    createProject() {
      const newId = `PRJ${String(this.projects.length + 1).padStart(3, '0')}`
      const today = new Date().toLocaleDateString('pt-BR')
      
      this.projects.push({
        id: newId,
        name: this.newProject.name,
        description: this.newProject.description,
        lastEdited: today,
        isEditing: false
      })
      
      this.closeModal()
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
