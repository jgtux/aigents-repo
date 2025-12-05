<template>
  <div class="main__chat">
    <!-- Sidebar -->
    <aside class="chat__sidebar" aria-label="Agent Information">
      <!-- Loading State -->
      <div v-if="loading" class="loading__container">
        <p>Loading agent...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="error__container">
        <p>{{ error }}</p>
        <button @click="loadAgentDetails" class="retry__button">Try Again</button>
      </div>

      <!-- Agent Info Card -->
      <article v-else class="agent__info__card">
        <figure class="agent__avatar" id="agentAvatar">
          <img 
            class="img__agents"
            :src="agent.image_url || require('@/assets/images/default-agent.png')" 
            :alt="agent.name"
            @error="handleImageError"
          >
        </figure>
        <header class="agent__header">
          <h2 class="agent__title" id="agentName">{{ agent.name }}</h2>
          <p class="agent__category" id="agentCategory">{{ agent.category }}</p>
        </header>
        <p class="agent__description" id="agentDescription">{{ agent.description }}</p>
        <footer class="agent__meta">
          <span class="agent__id__display" id="agentId">ID: {{ agent.agent_uuid }}</span>
        </footer>
        <nav>
          <button class="btn__back" @click="goBack" aria-label="Return to projects page">
            ← Back to Projects
          </button>
        </nav>
      </article>
    </aside>
    <!-- Chat Area -->
    <section class="chat__container" aria-label="Chat conversation">
      <article class="chat__messages" id="chatMessages" role="log" aria-live="polite" aria-atomic="false">
        <header class="welcome__message" v-if="messages.length === 0">
          <h3>Welcome! 👋</h3>
          <p>Start a conversation with your AI agent</p>
        </header>
        <div v-for="(msg, index) in messages" :key="index" class="chat__message">
          <div :class="['bubble', msg.sender]">
            {{ msg.text }}
          </div>
        </div>
      </article>
      <!-- Input Area -->
      <footer class="chat__input__container">
        <form class="chat__input__wrapper" @submit.prevent="sendMessage" aria-label="Message input form">
          <label for="chatInput" class="visually-hidden">Type your message</label>
          <textarea
            class="chat__input"
            id="chatInput"
            placeholder="Type your message here..."
            rows="1"
            aria-label="Message input"
            v-model="input"
            @keydown.enter.exact.prevent="sendMessage"
            @keydown.enter.shift="newLine"
            required
          ></textarea>
          <button class="btn__send" type="submit" title="Send message" aria-label="Send message">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
              <path d="M22 2L11 13" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M22 2L15 22L11 13L2 9L22 2Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </button>
        </form>
        <small class="chat__hints">Press Enter to send, Shift+Enter for new line</small>
      </footer>
    </section>
  </div>
</template>

<script>
import api from '@/api/api'

export default {
  name: "ChatPage",
  data() {
    return {
      agent: {
        name: "",
        category: "",
        description: "",
        agent_uuid: "",
        image_url: ""
      },
      loading: false,
      error: null,
      messages: [],
      input: "",
    };
  },
  mounted() {
    // Try different possible parameter names
    const agentUuid = this.$route.params.agent_uuid || 
                      this.$route.params.agentUuid || 
                      this.$route.params.id ||
                      this.$route.params.uuid
    
    console.log('Full route params:', this.$route.params)
    console.log('Agent UUID:', agentUuid)
    
    if (agentUuid && typeof agentUuid === 'string') {
      this.loadAgentDetails(agentUuid)
    } else {
      this.error = `No agent ID found. Available params: ${JSON.stringify(this.$route.params)}`
    }
  },
  methods: {
    async loadAgentDetails(agentUuid) {
      this.loading = true
      this.error = null
      
      try {
        const response = await api.get(`/agents/${agentUuid}`)
        
        // Handle Go Response structure
        if (response.data.status === 200) {
          this.agent = response.data.data
        } else {
          this.error = response.data.message || 'Failed to load agent'
        }
        
      } catch (err) {
        this.error = err.response?.data?.message || 'Failed to load agent. Please try again.'
        console.error('Error loading agent:', err)
      } finally {
        this.loading = false
      }
    },
    
    handleImageError(event) {
      event.target.src = require('@/assets/images/default-agent.png')
    },
    
    goBack() {
      this.$router.push("/agents");
    },
    
    newLine() {
      this.input += "\n"
    },
    
    sendMessage() {
      const text = this.input.trim();
      if (!text) return;
      this.messages.push({ sender: "user", text });
      this.input = "";
      // Placeholder for SSE/bot logic later
      setTimeout(() => {
        this.messages.push({ sender: "bot", text: "(Bot response placeholder)" });
      }, 600);
    },
  },
};
</script>

<style scoped>
/* Global styles are imported in main.js */
</style>
