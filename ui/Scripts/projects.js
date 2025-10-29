// Projects JavaScript - Gerenciamento de Projetos com armazenamento local

// Carregar projetos do armazenamento ou usar dados padrão
let projects = loadProjects();

// Elementos do DOM
const projectsGrid = document.getElementById('projectsGrid');
const newProjectCard = document.getElementById('newProjectCard');
const btnNewProject = document.getElementById('btnNewProject');
const modal = document.getElementById('modalNewProject');
const modalClose = document.querySelector('.modal__close');
const formNewProject = document.getElementById('formNewProject');
const searchInput = document.getElementById('searchProjects');

// Função para carregar projetos do armazenamento
function loadProjects() {
    const stored = localStorage.getItem('aigents_projects');
    if (stored) {
        return JSON.parse(stored);
    }
    // Projeto padrão se não houver nada salvo
    return [
        {
            id: 'PRJ001',
            name: 'Agent Name',
            description: 'Brief description of the AI agent',
            lastModified: '24/10/2025'
        }
    ];
}

// Função para salvar projetos no armazenamento
function saveProjects() {
    localStorage.setItem('aigents_projects', JSON.stringify(projects));
}

// Função para gerar ID único
function generateProjectId() {
    const num = projects.length + 1;
    return `PRJ${String(num).padStart(3, '0')}`;
}

// Função para obter data atual formatada
function getCurrentDate() {
    const date = new Date();
    const day = String(date.getDate()).padStart(2, '0');
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const year = date.getFullYear();
    return `${day}/${month}/${year}`;
}

/**
 * Cria um card HTML para exibir um projeto
 * @param {Object} project - Objeto contendo id, name, description e lastModified
 * @returns {HTMLElement} Elemento div com o card completo
 */

// Função para criar card de projeto
function createProjectCard(project) {
    // Cria um novo elemento div para o card
    const card = document.createElement('div');
    card.className = 'card__project';
    card.setAttribute('data-project-id', project.id);
    
    // Define o HTML interno do card usando template string
    // Backticks (`) permitem strings multi-linha e interpolação de variáveis com ${}
    card.innerHTML = `
        <div class="project__placeholder">IA</div>
        <div class="project__info">
            <h3 class="project__nome" contenteditable="false">${project.name}</h3>
            <p class="project__descricao" contenteditable="false">${project.description}</p>
            <div class="project__meta">
                <span class="project__id">ID: ${project.id}</span>
                <span class="project__data">Última edição: ${project.lastModified}</span>
            </div>
        </div>
        <div class="project__actions">
            <button class="btn__edit" title="Editar">✏️</button>
            <button class="btn__save" title="Salvar" style="display: none;">💾</button>
            <button class="btn__delete" title="Deletar">🗑️</button>
        </div>
    `;

    // CAPTURA ELEMENTOS INTERNOS DO CARD
    // querySelector busca elementos dentro do card específico
    const btnEdit = card.querySelector('.btn__edit');
    const btnSave = card.querySelector('.btn__save');
    const btnDelete = card.querySelector('.btn__delete');
    const nome = card.querySelector('.project__nome');
    const descricao = card.querySelector('.project__descricao');
    
    // addEventListener registra uma função que será executada quando o evento ocorrer
    btnEdit.addEventListener('click', () => {
        nome.contentEditable = 'true';
        descricao.contentEditable = 'true';
        nome.focus();
        btnEdit.style.display = 'none';
        btnSave.style.display = 'block';
    });
    
    // Salvar
    btnSave.addEventListener('click', () => {
        const newName = nome.textContent.trim();
        const newDesc = descricao.textContent.trim();
        
        if (newName && newDesc) {
            // Atualizar dados
            const projectIndex = projects.findIndex(p => p.id === project.id);
            projects[projectIndex].name = newName;
            projects[projectIndex].description = newDesc;
            projects[projectIndex].lastModified = getCurrentDate();
            
            // Salvar no armazenamento
            saveProjects();
            
            // Atualizar display da data
            card.querySelector('.project__data').textContent = `Última edição: ${getCurrentDate()}`;
            
            nome.contentEditable = 'false';
            descricao.contentEditable = 'false';
            btnSave.style.display = 'none';
            btnEdit.style.display = 'block';
            
            alert('Projeto atualizado com sucesso!');
        } else {
            alert('Nome e descrição não podem estar vazios!');
        }
    });
    
    // Deletar
    btnDelete.addEventListener('click', () => {
        if (confirm('Tem certeza que deseja deletar este projeto?')) {
            projects = projects.filter(p => p.id !== project.id);
            saveProjects(); // Salvar após deletar
            card.remove();
            alert('Projeto deletado com sucesso!');
        }
    });
    
    return card;
}

// Função para renderizar todos os projetos
function renderProjects(projectsToRender = projects) {
    // Limpar grid (manter apenas o card de criar novo)
    projectsGrid.innerHTML = '';
    
    // Adicionar cards dos projetos
    projectsToRender.forEach(project => {
        const card = createProjectCard(project);
        projectsGrid.appendChild(card);
    });
    
    // Adicionar card de criar novo no final
    projectsGrid.appendChild(newProjectCard);
}

// Abrir modal
function openModal() {
    modal.style.display = 'block';
    document.getElementById('projectName').value = '';
    document.getElementById('projectDesc').value = '';
}

// Fechar modal
function closeModal() {
    modal.style.display = 'none';
}

// Event listeners para abrir modal
btnNewProject.addEventListener('click', openModal);
newProjectCard.addEventListener('click', openModal);

// Event listener para fechar modal
modalClose.addEventListener('click', closeModal);

// Fechar modal ao clicar fora
window.addEventListener('click', (e) => {
    if (e.target === modal) {
        closeModal();
    }
});

// Criar novo projeto
formNewProject.addEventListener('submit', (e) => {
    e.preventDefault();
    
    // Captura os valores dos campos
    const name = document.getElementById('projectName').value.trim();
    const description = document.getElementById('projectDesc').value.trim();
    
    // Validação: verifica se os campos foram preenchidos
    if (name && description) {
        // Cria um novo objeto de projeto
        const newProject = {
            id: generateProjectId(),
            name: name,
            description: description,
            lastModified: getCurrentDate()
        };
        
        projects.push(newProject);
        saveProjects(); // Salvar após criar
        renderProjects();
        closeModal();
        alert('Projeto criado com sucesso!');
    }
});

// FUNCIONALIDADE DE BUSCA
searchInput.addEventListener('input', (e) => {
    const searchTerm = e.target.value.toLowerCase();
    
    if (searchTerm === '') {
        renderProjects();
    } else {
        const filteredProjects = projects.filter(project => 
            project.name.toLowerCase().includes(searchTerm) ||
            project.description.toLowerCase().includes(searchTerm) ||
            project.id.toLowerCase().includes(searchTerm)
        );
        renderProjects(filteredProjects);
    }
});

// Função para limpar todos os projetos (útil para testes)
function clearAllProjects() {
    if (confirm('Tem certeza que deseja limpar TODOS os projetos? Esta ação não pode ser desfeita!')) {
        localStorage.removeItem('aigents_projects');
        projects = loadProjects();
        renderProjects();
        alert('Todos os projetos foram limpos!');
    }
}

// Adicionar função global para limpar (opcional - pode usar no console)
window.clearAllProjects = clearAllProjects;

// Inicializar
renderProjects();