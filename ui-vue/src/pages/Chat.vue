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
            v-if="agent.image_url"
            :src="agent.image_url" 
            :alt="agent.name"
            @error="handleImageError"
            style="width: 100%; height: 100%; object-fit: cover; border-radius: 50%;"
          >
          <span v-else>{{ getInitials(agent.name) }}</span>
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
          <button class="btn__back" id="btnBack" @click="goBack" aria-label="Return to projects page">
            ← Back to Projects
          </button>
        </nav>
      </article>
    </aside>

    <!-- Chat Area -->
    <section class="chat__container" aria-label="Chat conversation">
      <article class="chat__messages" id="chatMessages" role="log" aria-live="polite" aria-atomic="false">
        <!-- Welcome Message -->
        <header class="welcome__message" v-if="messages.length === 0">
          <h3>Welcome! 👋</h3>
          <p>Start a conversation with {{ agent.name || 'your AI agent' }}</p>
        </header>
        
        <!-- Messages -->
        <article 
          v-for="(msg, index) in messages" 
          :key="index" 
          :class="['message', `message__${msg.sender}`]"
        >
          <!-- Avatar -->
          <figure class="message__avatar">
            <img 
              v-if="msg.sender === 'agent' && agent.image_url"
              :src="agent.image_url" 
              :alt="agent.name"
              style="width: 100%; height: 100%; object-fit: cover; border-radius: 50%;"
            >
            <span v-else-if="msg.sender === 'agent'">{{ getInitials(agent.name) }}</span>
            <span v-else>You</span>
          </figure>
          
          <!-- Message Content -->
          <section class="message__content">
            <p class="message__bubble">
              {{ msg.text }}<span v-if="msg.streaming" class="cursor__blink">|</span>
            </p>
            <time class="message__time">{{ msg.time }}</time>
          </section>
        </article>

        <!-- Typing Indicator -->
        <article v-if="isProcessing" class="message message__agent" id="typingIndicator">
          <figure class="message__avatar">
            <img 
              v-if="agent.image_url"
              :src="agent.image_url" 
              :alt="agent.name"
              style="width: 100%; height: 100%; object-fit: cover; border-radius: 50%;"
            >
            <span v-else>{{ getInitials(agent.name) }}</span>
          </figure>
          <section class="typing__indicator">
            <span class="typing__dot"></span>
            <span class="typing__dot"></span>
            <span class="typing__dot"></span>
          </section>
        </article>
      </article>

      <!-- Input Area -->
      <footer class="chat__input__container">
        <form class="chat__input__wrapper" id="chatForm" @submit.prevent="sendMessage" aria-label="Message input form">
          <label for="chatInput" class="visually-hidden">Type your message</label>
          <textarea
            class="chat__input"
            id="chatInput"
            placeholder="Type your message here..."
            rows="1"
            aria-label="Message input"
            v-model="input"
            @input="autoResizeTextarea"
            @keydown.enter.exact.prevent="sendMessage"
            @keydown.enter.shift="newLine"
            :disabled="isProcessing"
            required
          ></textarea>
          <button 
            class="btn__send" 
            id="btnSend"
            type="submit" 
            title="Send message" 
            aria-label="Send message"
            :disabled="isProcessing || !input.trim()"
          >
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
import { createChat, sendMessage as sendChatMessage } from '@/api/chat'

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
      chatUuid: null,
      isProcessing: false,
      currentStreamCleanup: null,
      currentBotMessageIndex: null
    };
  },
  mounted() {
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

    // Focus on input after load
    this.$nextTick(() => {
      const input = document.getElementById('chatInput')
      if (input) input.focus()
    })
  },
  beforeUnmount() {
    // Clean up any active SSE connection
    if (this.currentStreamCleanup) {
      this.currentStreamCleanup()
    }
  },
  methods: {
    async loadAgentDetails(agentUuid) {
      this.loading = true
      this.error = null
      
      try {
        const response = await api.get(`/agents/${agentUuid}`)
        
        if (response.data.status === 200) {
          this.agent = response.data.data
          // Store current agent in session storage
          sessionStorage.setItem('currentAgent', JSON.stringify(this.agent))
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
    
    getInitials(name) {
      if (!name) return 'AI'
      return name.substring(0, 2).toUpperCase()
    },
    
    handleImageError(event) {
      // Hide image and show initials instead
      event.target.style.display = 'none'
      const parent = event.target.parentElement
      if (parent) {
        const span = document.createElement('span')
        span.textContent = this.getInitials(this.agent.name)
        parent.appendChild(span)
      }
    },
    
    goBack() {
      this.$router.push("/agents");
    },
    
    newLine() {
      this.input += "\n"
    },
    
    getCurrentTime() {
      const now = new Date()
      return now.toLocaleTimeString('en-US', { 
        hour: '2-digit', 
        minute: '2-digit' 
      })
    },
    
    autoResizeTextarea() {
      this.$nextTick(() => {
        const textarea = document.getElementById('chatInput')
        if (textarea) {
          textarea.style.height = 'auto'
          textarea.style.height = Math.min(textarea.scrollHeight, 150) + 'px'
        }
      })
    },
    
    scrollToBottom() {
      this.$nextTick(() => {
        const chatMessages = document.getElementById('chatMessages')
        if (chatMessages) {
          chatMessages.scrollTop = chatMessages.scrollHeight
        }
      })
    },
    
    sendMessage() {
      const text = this.input.trim()
      if (!text || this.isProcessing) return
      
      // Add user message
      this.messages.push({ 
        sender: "user", 
        text,
        time: this.getCurrentTime(),
        streaming: false
      })
      
      this.input = ""
      this.autoResizeTextarea()
      this.scrollToBottom()
      this.isProcessing = true
      
      // Determine if this is the first message or a continuation
      if (!this.chatUuid) {
        this.initializeChat(text)
      } else {
        this.continueChat(text)
      }
    },
    
    initializeChat(messageContent) {
      this.scrollToBottom()
      
      // Start SSE stream
      this.currentStreamCleanup = createChat(
        this.agent.agent_uuid,
        messageContent,
        // onChunk
        (chunk) => {
          console.log('[DEBUG VUE] Received chunk:', JSON.stringify(chunk));
          
          // On first chunk, add the message
          if (this.currentBotMessageIndex === null) {
            this.currentBotMessageIndex = this.messages.length
            this.messages.push({
              sender: "agent",
              text: chunk,
              time: "",
              streaming: true
            })
          } else {
            // Just concatenate - LLM should send properly formatted text
            this.messages[this.currentBotMessageIndex].text += chunk
          }
          this.scrollToBottom()
        },
        // onComplete
        (responseData) => {
          console.log('[DEBUG VUE] Chat complete, final data:', responseData);
          
          if (this.currentBotMessageIndex !== null && this.messages[this.currentBotMessageIndex]) {
            this.messages[this.currentBotMessageIndex].streaming = false
            this.messages[this.currentBotMessageIndex].time = this.getCurrentTime()
          }
          
          if (responseData && responseData.chat_uuid) {
            this.chatUuid = responseData.chat_uuid
            console.log('[DEBUG VUE] Chat initialized with UUID:', this.chatUuid)
          }
          
          this.isProcessing = false
          this.currentStreamCleanup = null
          this.currentBotMessageIndex = null
          this.scrollToBottom()
          
          const input = document.getElementById('chatInput')
          if (input) input.focus()
        },
        // onError
        (error) => {
          console.error('[DEBUG VUE] Chat creation error:', error)
          
          if (this.currentBotMessageIndex === null) {
            this.messages.push({
              sender: "agent",
              text: error.message || 'Sorry, there was an error processing your message.',
              time: this.getCurrentTime(),
              streaming: false
            })
          } else if (this.messages[this.currentBotMessageIndex]) {
            this.messages[this.currentBotMessageIndex].text = error.message || 'Sorry, there was an error processing your message.'
            this.messages[this.currentBotMessageIndex].streaming = false
            this.messages[this.currentBotMessageIndex].time = this.getCurrentTime()
          }
          
          this.isProcessing = false
          this.currentStreamCleanup = null
          this.currentBotMessageIndex = null
          
          const input = document.getElementById('chatInput')
          if (input) input.focus()
        }
      )
    },
    
    continueChat(messageContent) {
      this.scrollToBottom()
      
      // Start SSE stream
      this.currentStreamCleanup = sendChatMessage(
        this.chatUuid,
        messageContent,
        // onChunk
        (chunk) => {
          console.log('[DEBUG VUE] Received chunk:', JSON.stringify(chunk));
          
          // On first chunk, add the message
          if (this.currentBotMessageIndex === null) {
            this.currentBotMessageIndex = this.messages.length
            this.messages.push({
              sender: "agent",
              text: chunk,
              time: "",
              streaming: true
            })
          } else {
            // Just concatenate - LLM should send properly formatted text
            this.messages[this.currentBotMessageIndex].text += chunk
          }
          this.scrollToBottom()
        },
        // onComplete
        () => {
          console.log('[DEBUG VUE] Message complete');
          
          if (this.currentBotMessageIndex !== null && this.messages[this.currentBotMessageIndex]) {
            this.messages[this.currentBotMessageIndex].streaming = false
            this.messages[this.currentBotMessageIndex].time = this.getCurrentTime()
          }
          
          this.isProcessing = false
          this.currentStreamCleanup = null
          this.currentBotMessageIndex = null
          this.scrollToBottom()
          
          const input = document.getElementById('chatInput')
          if (input) input.focus()
        },
        // onError
        (error) => {
          console.error('[DEBUG VUE] Message send error:', error)
          
          if (this.currentBotMessageIndex === null) {
            this.messages.push({
              sender: "agent",
              text: error.message || 'Sorry, there was an error processing your message.',
              time: this.getCurrentTime(),
              streaming: false
            })
          } else if (this.messages[this.currentBotMessageIndex]) {
            this.messages[this.currentBotMessageIndex].text = error.message || 'Sorry, there was an error processing your message.'
            this.messages[this.currentBotMessageIndex].streaming = false
            this.messages[this.currentBotMessageIndex].time = this.getCurrentTime()
          }
          
          this.isProcessing = false
          this.currentStreamCleanup = null
          this.currentBotMessageIndex = null
          
          const input = document.getElementById('chatInput')
          if (input) input.focus()
        }
      )
    }
  }
}
</script>

<style scoped>
/* Cursor blink animation for streaming */
.cursor__blink {
  animation: blink 1s infinite;
  font-weight: 100;
  margin-left: 2px;
}

@keyframes blink {
  0%, 49% { opacity: 1; }
  50%, 100% { opacity: 0; }
}

.visually-hidden {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border: 0;
}
</style>
