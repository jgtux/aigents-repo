// chat.js - Chat functionality for AIgents

document.addEventListener('DOMContentLoaded', function() {
    const chatMessages = document.getElementById('chatMessages');
    const chatInput = document.getElementById('chatInput');
    const btnSend = document.getElementById('btnSend');
    const btnBack = document.getElementById('btnBack');
    const chatForm = document.getElementById('chatForm');
    
    // Elementos de informação do agente
    const agentName = document.getElementById('agentName');
    const agentCategory = document.getElementById('agentCategory');
    const agentDescription = document.getElementById('agentDescription');
    const agentId = document.getElementById('agentId');
    const agentAvatar = document.getElementById('agentAvatar');

    // Carregar informações do agente
    loadAgentInfo();

    // Event Listeners
    chatForm.addEventListener('submit', handleSendMessage);
    btnBack.addEventListener('click', () => window.location.href = 'Myprojects.html');
    
    // Auto-resize textarea
    chatInput.addEventListener('input', autoResizeTextarea);
    
    // Ctrl/Cmd + Enter para enviar
    chatInput.addEventListener('keydown', handleKeyDown);

    // Funções
    function loadAgentInfo() {
        // Pegar ID do projeto da URL
        const urlParams = new URLSearchParams(window.location.search);
        const projectId = urlParams.get('id');
        
        if (!projectId) {
            console.error('No project ID provided');
            return;
        }

        // Buscar projeto no localStorage
        const projetos = JSON.parse(localStorage.getItem('projetos')) || [];
        const projeto = projetos.find(p => p.id === projectId);
        
        if (projeto) {
            agentName.textContent = projeto.nome;
            agentCategory.textContent = projeto.categoria || 'General';
            agentDescription.textContent = projeto.descricao;
            agentId.textContent = `ID: ${projeto.id}`;
            
            // Se tiver imagem, mostrar no avatar
            if (projeto.imagem) {
                agentAvatar.innerHTML = `<img src="${projeto.imagem}" alt="${projeto.nome}">`;
            } else {
                agentAvatar.textContent = projeto.nome.substring(0, 2).toUpperCase();
            }
            
            // Armazenar o agente atual
            sessionStorage.setItem('currentAgent', JSON.stringify(projeto));
        } else {
            console.error('Project not found');
            agentName.textContent = 'Agent Not Found';
            agentDescription.textContent = 'This agent could not be loaded.';
        }
    }

    function handleSendMessage(e) {
        e.preventDefault();
        
        const message = chatInput.value.trim();
        
        if (!message) return;
        
        // Adicionar mensagem do usuário
        addMessage(message, 'user');
        
        // Limpar input
        chatInput.value = '';
        autoResizeTextarea();
        
        // Desabilitar envio enquanto processa
        btnSend.disabled = true;
        
        // Simular resposta do agente (typing indicator)
        showTypingIndicator();
        
        // Simular delay de resposta
        setTimeout(() => {
            hideTypingIndicator();
            const agentResponse = generateAgentResponse(message);
            addMessage(agentResponse, 'agent');
            btnSend.disabled = false;
            chatInput.focus();
        }, 1500 + Math.random() * 1000);
    }

    function addMessage(text, sender) {
        // Remover mensagem de boas-vindas se existir
        const welcomeMessage = chatMessages.querySelector('.welcome__message');
        if (welcomeMessage) {
            welcomeMessage.remove();
        }

        const messageElement = document.createElement('article');
        messageElement.className = `message message__${sender}`;
        
        const avatar = document.createElement('figure');
        avatar.className = 'message__avatar';
        
        if (sender === 'user') {
            avatar.textContent = 'You';
        } else {
            const currentAgent = JSON.parse(sessionStorage.getItem('currentAgent'));
            if (currentAgent && currentAgent.imagem) {
                avatar.innerHTML = `<img src="${currentAgent.imagem}" alt="${currentAgent.nome}" style="width: 100%; height: 100%; object-fit: cover; border-radius: 50%;">`;
            } else {
                avatar.textContent = currentAgent ? currentAgent.nome.substring(0, 2).toUpperCase() : 'AI';
            }
        }
        
        const contentWrapper = document.createElement('section');
        contentWrapper.className = 'message__content';
        
        const bubble = document.createElement('p');
        bubble.className = 'message__bubble';
        bubble.textContent = text;
        
        const time = document.createElement('time');
        time.className = 'message__time';
        time.textContent = getCurrentTime();
        
        contentWrapper.appendChild(bubble);
        contentWrapper.appendChild(time);
        
        messageElement.appendChild(avatar);
        messageElement.appendChild(contentWrapper);
        
        chatMessages.appendChild(messageElement);
        
        // Scroll para o final
        scrollToBottom();
        
        // Salvar no histórico
        saveMessageToHistory(text, sender);
    }

    function showTypingIndicator() {
        const indicator = document.createElement('article');
        indicator.className = 'message message__agent';
        indicator.id = 'typingIndicator';
        
        const avatar = document.createElement('figure');
        avatar.className = 'message__avatar';
        const currentAgent = JSON.parse(sessionStorage.getItem('currentAgent'));
        if (currentAgent && currentAgent.imagem) {
            avatar.innerHTML = `<img src="${currentAgent.imagem}" alt="${currentAgent.nome}" style="width: 100%; height: 100%; object-fit: cover; border-radius: 50%;">`;
        } else {
            avatar.textContent = currentAgent ? currentAgent.nome.substring(0, 2).toUpperCase() : 'AI';
        }
        
        const typingDiv = document.createElement('section');
        typingDiv.className = 'typing__indicator';
        typingDiv.innerHTML = '<span class="typing__dot"></span><span class="typing__dot"></span><span class="typing__dot"></span>';
        
        indicator.appendChild(avatar);
        indicator.appendChild(typingDiv);
        
        chatMessages.appendChild(indicator);
        scrollToBottom();
    }

    function hideTypingIndicator() {
        const indicator = document.getElementById('typingIndicator');
        if (indicator) {
            indicator.remove();
        }
    }

    function generateAgentResponse(userMessage) {
        // Respostas simuladas - aqui você pode integrar com APIs de IA reais
        const responses = [
            "That's an interesting question! Let me help you with that.",
            "I understand what you're asking. Here's what I think...",
            "Great question! Based on what you've told me, I'd suggest...",
            "I'm here to help! Let me break this down for you.",
            "Thanks for asking! Here's my take on that...",
            "I can definitely assist with that. Let me explain...",
            "That's a good point. From my perspective...",
            "I appreciate you sharing that. Here's what I recommend..."
        ];
        
        const randomResponse = responses[Math.floor(Math.random() * responses.length)];
        
        // Adicionar contexto baseado na mensagem do usuário
        if (userMessage.toLowerCase().includes('help')) {
            return "I'm here to help! What specific task would you like assistance with?";
        } else if (userMessage.toLowerCase().includes('code') || userMessage.toLowerCase().includes('program')) {
            return "I'd be happy to help with coding! What programming language or concept are you working with?";
        } else if (userMessage.toLowerCase().includes('thank')) {
            return "You're welcome! Is there anything else I can help you with?";
        }
        
        return randomResponse;
    }

    function getCurrentTime() {
        const now = new Date();
        return now.toLocaleTimeString('en-US', { 
            hour: '2-digit', 
            minute: '2-digit' 
        });
    }

    function scrollToBottom() {
        chatMessages.scrollTop = chatMessages.scrollHeight;
    }

    function autoResizeTextarea() {
        chatInput.style.height = 'auto';
        chatInput.style.height = Math.min(chatInput.scrollHeight, 150) + 'px';
    }

    function handleKeyDown(e) {
        // Enter sem Shift envia a mensagem
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            chatForm.dispatchEvent(new Event('submit'));
        }
    }

    function saveMessageToHistory(message, sender) {
        const urlParams = new URLSearchParams(window.location.search);
        const projectId = urlParams.get('id');
        
        if (!projectId) return;
        
        // Buscar histórico existente
        const historyKey = `chat_history_${projectId}`;
        let history = JSON.parse(localStorage.getItem(historyKey)) || [];
        
        // Adicionar nova mensagem
        history.push({
            message: message,
            sender: sender,
            timestamp: new Date().toISOString()
        });
        
        // Limitar histórico a 100 mensagens
        if (history.length > 100) {
            history = history.slice(-100);
        }
        
        // Salvar no localStorage
        localStorage.setItem(historyKey, JSON.stringify(history));
    }

    function loadChatHistory() {
        const urlParams = new URLSearchParams(window.location.search);
        const projectId = urlParams.get('id');
        
        if (!projectId) return;
        
        const historyKey = `chat_history_${projectId}`;
        const history = JSON.parse(localStorage.getItem(historyKey)) || [];
        
        // Remover mensagem de boas-vindas se houver histórico
        if (history.length > 0) {
            const welcomeMessage = chatMessages.querySelector('.welcome__message');
            if (welcomeMessage) {
                welcomeMessage.remove();
            }
            
            // Carregar mensagens anteriores
            history.forEach(item => {
                addMessage(item.message, item.sender);
            });
        }
    }

    // Carregar histórico ao iniciar
    loadChatHistory();
    
    // Focar no input ao carregar
    chatInput.focus();
});