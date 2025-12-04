<template>
  <div class="main__chat">
    <!-- Sidebar -->
    <aside class="chat__sidebar" aria-label="Agent Information">
      <article class="agent__info__card">
        <figure class="agent__avatar" id="agentAvatar">IA</figure>

        <header class="agent__header">
          <h2 class="agent__title" id="agentName">{{ agent.name }}</h2>
          <p class="agent__category" id="agentCategory">{{ agent.category }}</p>
        </header>

        <p class="agent__description" id="agentDescription">{{ agent.description }}</p>

        <footer class="agent__meta">
          <span class="agent__id__display" id="agentId">ID: {{ agent.id }}</span>
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
export default {
  name: "ChatPage",
  data() {
    return {
      agent: {
        name: "Agent Name",
        category: "Category",
        description: "Agent description will appear here",
        id: "PRJ001",
      },
      messages: [],
      input: "",
    };
  },
  methods: {
    goBack() {
      this.$router.push("/projects");
    },

    newLine() {
      const textarea = this.$refs.chatInput;
      if (textarea) textarea.value += "\n";
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
