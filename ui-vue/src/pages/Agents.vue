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
    
    <main class="main__agents">
      <!-- Loading State -->
      <div v-if="loading" class="loading__container">
        <p>Loading agents...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="error__container">
        <p>{{ error }}</p>
        <button @click="loadAgents" class="retry__button">Try Again</button>
      </div>

      <!-- Agents Grid -->
      <section v-else class="grid__agents">
        <article 
          v-for="agent in agents" 
          :key="agent.agent_uuid"
          class="card__agent"
          :data-name="agent.name.toLowerCase()"
        >
          <router-link :to="`/chat/${agent.agent_uuid}`">
            <div class="agent__placeholder">
              <img 
                class="img__agents"
                :src="agent.image_url || require('@/assets/images/default-agent.png')" 
                :alt="agent.name"
                @error="handleImageError"
              >
            </div>
            <h3 class="agent__nome">{{ agent.name }}</h3>
            <p class="agent__descricao">{{ agent.description }}</p>
          </router-link>
        </article>

        <!-- Empty State -->
        <div v-if="agents.length === 0" class="empty__state">
          <p>No agents found.</p>
        </div>
      </section>

      <!-- Pagination -->
      <div v-if="!loading && !error && hasMorePages" class="pagination">
        <button 
          @click="previousPage" 
          :disabled="currentPage === 0"
          class="pagination__button"
        >
          Previous
        </button>
        
        <span class="pagination__info">
          Page {{ currentPage + 1 }}
        </span>
        
        <button 
          @click="nextPage" 
          :disabled="agents.length < pageSize"
          class="pagination__button"
        >
          Next
        </button>
      </div>
    </main>
  </div>
</template>

<script>
import api from '@/api/api'

export default {
  name: 'AgentsPage',
  data() {
    return {
      searchQuery: '',
      agents: [],
      loading: false,
      error: null,
      currentPage: 0,
      pageSize: 20
    }
  },
  computed: {
    hasMorePages() {
      return this.agents.length === this.pageSize || this.currentPage > 0
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
    async loadAgents() {
      this.loading = true
      this.error = null
      
      try {
        const response = await api.post('/agents/all', {
          page: this.currentPage,
          page_size: this.pageSize
        })
        
        // Handle Go Response[[]Agent] structure
        if (response.data.status === 200) {
          this.agents = response.data.data || []
        } else {
          this.error = response.data.message || 'Failed to load agents'
        }
        
      } catch (err) {
        this.error = err.response?.data?.message || 'Failed to load agents. Please try again.'
        console.error('Error loading agents:', err)
      } finally {
        this.loading = false
      }
    },
    
    handleSearch() {
      const searchTerm = this.searchQuery.toLowerCase().trim()
      const cards = document.querySelectorAll('.card__agent')

      cards.forEach(card => {
        const agentName = card.dataset.name || ''
        const agentDescription = card.querySelector('.agent__descricao')?.textContent.toLowerCase() || ''

        if (agentName.includes(searchTerm) || agentDescription.includes(searchTerm)) {
          card.style.display = 'block'
        } else {
          card.style.display = 'none'
        }
      })
    },
    
    handleImageError(event) {
      event.target.src = require('@/assets/images/default-agent.png')
    },
    
    nextPage() {
      if (this.agents.length === this.pageSize) {
        this.currentPage++
        window.scrollTo({ top: 0, behavior: 'smooth' })
        this.loadAgents()
      }
    },
    
    previousPage() {
      if (this.currentPage > 0) {
        this.currentPage--
        window.scrollTo({ top: 0, behavior: 'smooth' })
        this.loadAgents()
      }
    }
  },
  
  mounted() {
    this.loadAgents()
  },
  
  beforeUnmount() {
    if (this.searchTimeout) {
      clearTimeout(this.searchTimeout)
    }
  },
  
  updated() {
    // Setup search functionality after DOM update
    this.$nextTick(() => {
      if (this.searchQuery) {
        this.handleSearch()
      }
    })
  }
}
</script>
