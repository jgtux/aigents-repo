// Dados dos agentes
const agents = [
    {
        name: "Claude AI",
        description: "Good for texts and programming",
        image: "../img/Claude_logo.webp",
        url: "https://claude.ai/new",
        categories: ["programming", "texts"]
    },
    {
        name: "Chat GPT",
        description: "Good for any task image, video, texts, programming...",
        image: "../img/GPT_logo.png",
        url: "https://chatgpt.com/?openaicom_referred=true",
        categories: ["programming", "images", "videos", "texts"]
    },
    {
        name: "Gemini AI",
        description: "Good for any task image, video, texts, programming...",
        image: "../img/Gemini_logo.webp",
        url: "https://gemini.google.com/app?hl=pt-BR",
        categories: ["programming", "images", "videos", "texts"]
    },
    {
        name: "Leonardo AI",
        description: "Good for image Creating at all",
        image: "/ui/img/Leonardo_logo.png",
        url: "https://app.leonardo.ai/",
        categories: ["images"]
    }
];

// Quantidade de cards vazios para completar o grid
const TOTAL_CARDS = 12;

// Função para criar um card de agente
function createAgentCard(agent) {
    return `
        <article class="card__agent" data-name="${agent.name.toLowerCase()}" data-categories="${agent.categories.join(',')}">
            <a href="${agent.url}" target="_blank">
                <div class="agent__placeholder">
                    <img class="img__agents" src="${agent.image}" alt="${agent.name}">
                </div>
                <h3 class="agent__nome">${agent.name}</h3>
                <p class="agent__descricao">${agent.description}</p>
            </a>
        </article>
    `;
}

// Função para criar um card vazio
function createEmptyCard() {
    return `
        <article class="card__agent card__agent--empty">
            <div class="agent__placeholder">IA</div>
            <h3 class="agent__nome">Agent Name</h3>
            <p class="agent__descricao">Brief description of the AI agent</p>
        </article>
    `;
}

// Função para renderizar todos os cards
function renderCards() {
    const grid = document.querySelector('.grid__agents');
    if (!grid) return;

    let html = '';
    
    // Adiciona os agentes reais
    agents.forEach(agent => {
        html += createAgentCard(agent);
    });
    
    // Adiciona cards vazios para completar o grid
    const emptyCardsCount = TOTAL_CARDS - agents.length;
    for (let i = 0; i < emptyCardsCount; i++) {
        html += createEmptyCard();
    }
    
    grid.innerHTML = html;
}

// Função de busca
function setupSearch() {
    const searchInput = document.querySelector('.search__input');
    if (!searchInput) return;

    searchInput.addEventListener('input', (e) => {
        const searchTerm = e.target.value.toLowerCase().trim();
        const cards = document.querySelectorAll('.card__agent');

        cards.forEach(card => {
            // Se houver busca ativa, esconde os cards vazios
            if (card.classList.contains('card__agent--empty')) {
                card.style.display = searchTerm ? 'none' : 'block';
                return;
            }

            const agentName = card.dataset.name || '';
            const agentDescription = card.querySelector('.agent__descricao')?.textContent.toLowerCase() || '';

            // Verifica se o termo de busca está no nome ou descrição
            if (agentName.includes(searchTerm) || agentDescription.includes(searchTerm)) {
                card.style.display = 'block';
            } else {
                card.style.display = 'none';
            }
        });
    });
}

// Inicializa quando o DOM estiver carregado
document.addEventListener('DOMContentLoaded', () => {
    renderCards();
    setupSearch();
});

// Adiciona funcionalidade de ícone de busca
document.addEventListener('DOMContentLoaded', () => {
    const searchIcon = document.querySelector('.search__icon');
    const searchInput = document.querySelector('.search__input');
    
    if (searchIcon && searchInput) {
        searchIcon.addEventListener('click', () => {
            searchInput.focus();
        });
    }
});