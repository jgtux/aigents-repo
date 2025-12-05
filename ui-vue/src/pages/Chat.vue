<template>
  <div class="main__chat">
    <!-- Sidebar lateral com informações do agente -->
    <aside class="chat__sidebar" aria-label="Agent Information">
      <!-- Estado de carregamento -->
      <div v-if="loading" class="loading__container">
        <p>Loading agent...</p>
      </div>

      <!-- Estado de erro -->
      <div v-else-if="error" class="error__container">
        <p>{{ error }}</p>
        <button @click="loadAgentDetails" class="retry__button">Try Again</button>
      </div>

      <!-- Card com informações do agente -->
      <article v-else class="agent__info__card">
        <figure class="agent__avatar" id="agentAvatar">
          <!-- Foto do agente, se existir -->
          <img 
            v-if="agent.image_url"
            :src="agent.image_url" 
            :alt="agent.name"
            @error="handleImageError"
            style="width: 100%; height: 100%; object-fit: cover; border-radius: 50%;"
          >
          <!-- Caso não tenha imagem, mostra iniciais -->
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
          <!-- Botão para voltar à listagem de agentes/projetos -->
          <button class="btn__back" id="btnBack" @click="goBack" aria-label="Return to projects page">
            ← Back to Projects
          </button>
        </nav>
      </article>
    </aside>

    <!-- Área principal do chat -->
    <section class="chat__container" aria-label="Chat conversation">
      <!-- Lista de mensagens -->
      <article class="chat__messages" id="chatMessages" role="log" aria-live="polite" aria-atomic="false">
        <!-- Mensagem de boas-vindas quando não há mensagens -->
        <header class="welcome__message" v-if="messages.length === 0">
          <h3>Welcome! 👋</h3>
          <p>Start a conversation with {{ agent.name || 'your AI agent' }}</p>
        </header>
        
        <!-- Loop de mensagens (usuário e agente) -->
        <article 
          v-for="(msg, index) in messages" 
          :key="index" 
          :class="['message', `message__${msg.sender}`]"
        >
          <!-- Avatar de cada mensagem -->
          <figure class="message__avatar">
            <!-- Se for mensagem do agente e ele tiver imagem, mostra a imagem -->
            <img 
              v-if="msg.sender === 'agent' && agent.image_url"
              :src="agent.image_url" 
              :alt="agent.name"
              style="width: 100%; height: 100%; object-fit: cover; border-radius: 50%;"
            >
            <!-- Se for do agente sem imagem, mostra iniciais -->
            <span v-else-if="msg.sender === 'agent'">{{ getInitials(agent.name) }}</span>
            <!-- Caso contrário, é mensagem do usuário -->
            <span v-else>You</span>
          </figure>
          
          <!-- Conteúdo da mensagem -->
          <section class="message__content">
            <p class="message__bubble">
              {{ msg.text }}<span v-if="msg.streaming" class="cursor__blink">|</span>
            </p>
            <time class="message__time">{{ msg.time }}</time>
          </section>
        </article>

        <!-- Indicador de digitação do agente (enquanto processa) -->
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

      <!-- Área de input do chat -->
      <footer class="chat__input__container">
        <!-- Formulário de envio de mensagem -->
        <form class="chat__input__wrapper" id="chatForm" @submit.prevent="sendMessage" aria-label="Message input form">
          <label for="chatInput" class="visually-hidden">Type your message</label>
          <!-- Textarea controlando o input do usuário -->
          <textarea
            class="chat__input"
            id="chatInput"
            placeholder="Type your message here..."
            rows="1"
            aria-label="Message input"
            v-model="input"
            @input="autoResizeTextarea"
            @keydown.enter.exact.prevent="sendMessage"   <!-- Enter envia -->
            @keydown.enter.shift="newLine"              <!-- Shift+Enter quebra linha -->
            :disabled="isProcessing"
            required
          ></textarea>
          <!-- Botão de enviar mensagem -->
          <button 
            class="btn__send" 
            id="btnSend"
            type="submit" 
            title="Send message" 
            aria-label="Send message"
            :disabled="isProcessing || !input.trim()"
          >
            <!-- Ícone de avião de papel (enviar) -->
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
      // Dados do agente carregados pela API
      agent: {
        name: "",
        category: "",
        description: "",
        agent_uuid: "",
        image_url: ""
      },
      // Flags de estado da tela
      loading: false,          // indica se está carregando dados do agente
      error: null,             // mensagem de erro, se existir
      messages: [],            // histórico de mensagens do chat
      input: "",               // texto atual digitado pelo usuário
      chatUuid: null,          // identificador do chat retornado pela API
      isProcessing: false,     // indica se o backend está processando a mensagem
      currentStreamCleanup: null,     // função para encerrar a conexão de streaming (SSE)
      currentBotMessageIndex: null    // índice da mensagem atual do bot em streaming
    };
  },
  mounted() {
    // Recupera o UUID do agente a partir dos parâmetros da rota,
    // aceitando múltiplas possíveis chaves (flexibilidade)
    const agentUuid = this.$route.params.agent_uuid || 
                      this.$route.params.agentUuid || 
                      this.$route.params.id ||
                      this.$route.params.uuid
    
    console.log('Full route params:', this.$route.params)
    console.log('Agent UUID:', agentUuid)
    
    // Se encontrou um UUID válido, carrega detalhes do agente
    if (agentUuid && typeof agentUuid === 'string') {
      this.loadAgentDetails(agentUuid)
    } else {
      // Caso não tenha, mostra erro com os params disponíveis
      this.error = `No agent ID found. Available params: ${JSON.stringify(this.$route.params)}`
    }

    // Depois de montar o componente, foca o input do chat
    this.$nextTick(() => {
      const input = document.getElementById('chatInput')
      if (input) input.focus()
    })
  },
  beforeUnmount() {
    // Antes de destruir o componente, encerra qualquer conexão SSE ativa
    if (this.currentStreamCleanup) {
      this.currentStreamCleanup()
    }
  },
  methods: {
    // Carrega detalhes do agente a partir da API
    async loadAgentDetails(agentUuid) {
      this.loading = true
      this.error = null
      
      try {
        const response = await api.get(`/agents/${agentUuid}`)
        
        // Verifica se a API retornou sucesso
        if (response.data.status === 200) {
          this.agent = response.data.data
          // Armazena o agente atual no sessionStorage (para uso posterior)
          sessionStorage.setItem('currentAgent', JSON.stringify(this.agent))
        } else {
          this.error = response.data.message || 'Failed to load agent'
        }
        
      } catch (err) {
        // Tratamento de erros de rede/servidor
        this.error = err.response?.data?.message || 'Failed to load agent. Please try again.'
        console.error('Error loading agent:', err)
      } finally {
        this.loading = false
      }
    },
    
    // Gera iniciais do nome do agente (para avatar sem imagem)
    getInitials(name) {
      if (!name) return 'AI'
      return name.substring(0, 2).toUpperCase()
    },
    
    // Handler de erro de carregamento da imagem do agente
    handleImageError(event) {
      // Esconde a imagem quebrada e insere um span com as iniciais
      event.target.style.display = 'none'
      const parent = event.target.parentElement
      if (parent) {
        const span = document.createElement('span')
        span.textContent = this.getInitials(this.agent.name)
        parent.appendChild(span)
      }
    },
    
    // Navega de volta para a página de agentes
    goBack() {
      this.$router.push("/agents");
    },
    
    // Insere uma nova linha no textarea (usado com Shift+Enter)
    newLine() {
      this.input += "\n"
    },
    
    // Retorna horário atual formatado (hh:mm)
    getCurrentTime() {
      const now = new Date()
      return now.toLocaleTimeString('en-US', { 
        hour: '2-digit', 
        minute: '2-digit' 
      })
    },
    
    // Ajusta automaticamente o tamanho do textarea de acordo com o conteúdo
    autoResizeTextarea() {
      this.$nextTick(() => {
        const textarea = document.getElementById('chatInput')
        if (textarea) {
          textarea.style.height = 'auto'
          textarea.style.height = Math.min(textarea.scrollHeight, 150) + 'px'
        }
      })
    },
    
    // Faz o scroll da área de mensagens para o final
    scrollToBottom() {
      this.$nextTick(() => {
        const chatMessages = document.getElementById('chatMessages')
        if (chatMessages) {
          chatMessages.scrollTop = chatMessages.scrollHeight
        }
      })
    },
    
    // Envia a mensagem do usuário
    sendMessage() {
      const text = this.input.trim()
      // Não envia se estiver vazio ou já estiver processando
      if (!text || this.isProcessing) return
      
      // Adiciona mensagem do usuário na lista
      this.messages.push({ 
        sender: "user", 
        text,
        time: this.getCurrentTime(),
        streaming: false
      })
      
      // Limpa input e ajusta textarea/scroll
      this.input = ""
      this.autoResizeTextarea()
      this.scrollToBottom()
      this.isProcessing = true
      
      // Decide se cria um novo chat ou continua um existente
      if (!this.chatUuid) {
        this.initializeChat(text)
      } else {
        this.continueChat(text)
      }
    },
    
    // Inicializa um novo chat com o backend (primeira mensagem)
    initializeChat(messageContent) {
      this.scrollToBottom()
      
      // Inicia o streaming via SSE
      this.currentStreamCleanup = createChat(
        this.agent.agent_uuid,
        messageContent,
        // onChunk: callback chamado a cada pedaço de resposta recebido
        (chunk) => {
          console.log('[DEBUG VUE] Received chunk:', JSON.stringify(chunk));
          
          // No primeiro chunk, cria a mensagem do agente
          if (this.currentBotMessageIndex === null) {
            this.currentBotMessageIndex = this.messages.length
            this.messages.push({
              sender: "agent",
              text: chunk,
              time: "",
              streaming: true
            })
          } else {
            // Nos seguintes, só concatena texto ao que já está lá
            this.messages[this.currentBotMessageIndex].text += chunk
          }
          this.scrollToBottom()
        },
        // onComplete: chamado quando o streaming termina
        (responseData) => {
          console.log('[DEBUG VUE] Chat complete, final data:', responseData);
          
          // Finaliza o estado de streaming da mensagem do bot
          if (this.currentBotMessageIndex !== null && this.messages[this.currentBotMessageIndex]) {
            this.messages[this.currentBotMessageIndex].streaming = false
            this.messages[this.currentBotMessageIndex].time = this.getCurrentTime()
          }
          
          // Guarda o UUID do chat retornado pelo backend
          if (responseData && responseData.chat_uuid) {
            this.chatUuid = responseData.chat_uuid
            console.log('[DEBUG VUE] Chat initialized with UUID:', this.chatUuid)
          }
          
          // Limpa estados de processamento/stream
          this.isProcessing = false
          this.currentStreamCleanup = null
          this.currentBotMessageIndex = null
          this.scrollToBottom()
          
          // Volta o foco para o input
          const input = document.getElementById('chatInput')
          if (input) input.focus()
        },
        // onError: tratado quando há erro no streaming/criação do chat
        (error) => {
          console.error('[DEBUG VUE] Chat creation error:', error)
          
          if (this.currentBotMessageIndex === null) {
            // Se não havia resposta parcial, cria uma mensagem de erro do agente
            this.messages.push({
              sender: "agent",
              text: error.message || 'Sorry, there was an error processing your message.',
              time: this.getCurrentTime(),
              streaming: false
            })
          } else if (this.messages[this.currentBotMessageIndex]) {
            // Se já tinha resposta parcial, substitui pelo erro
            this.messages[this.currentBotMessageIndex].text = error.message || 'Sorry, there was an error processing your message.'
            this.messages[this.currentBotMessageIndex].streaming = false
            this.messages[this.currentBotMessageIndex].time = this.getCurrentTime()
          }
          
          // Reseta estados
          this.isProcessing = false
          this.currentStreamCleanup = null
          this.currentBotMessageIndex = null
          
          const input = document.getElementById('chatInput')
          if (input) input.focus()
        }
      )
    },
    
    // Continua um chat já existente (após a primeira mensagem)
    continueChat(messageContent) {
      this.scrollToBottom()
      
      // Inicia o streaming de resposta para uma nova mensagem
      this.currentStreamCleanup = sendChatMessage(
        this.chatUuid,
        messageContent,
        // onChunk: recebe partes da resposta
        (chunk) => {
          console.log('[DEBUG VUE] Received chunk:', JSON.stringify(chunk));
          
          // No primeiro chunk cria a mensagem do agente
          if (this.currentBotMessageIndex === null) {
            this.currentBotMessageIndex = this.messages.length
            this.messages.push({
              sender: "agent",
              text: chunk,
              time: "",
              streaming: true
            })
          } else {
            // Nos demais, concatena o texto
            this.messages[this.currentBotMessageIndex].text += chunk
          }
          this.scrollToBottom()
        },
        // onComplete: chamado quando o backend termina a resposta
        () => {
          console.log('[DEBUG VUE] Message complete');
          
          if (this.currentBotMessageIndex !== null && this.messages[this.currentBotMessageIndex]) {
