// projects.js - Script atualizado para a página My Projects

document.addEventListener('DOMContentLoaded', function() {
    const btnNewProject = document.getElementById('btnNewProject');
    const newProjectCard = document.getElementById('newProjectCard');
    const modal = document.getElementById('modalNewProject');
    const modalClose = document.querySelector('.modal__close');
    const formNewProject = document.getElementById('formNewProject');
    const searchInput = document.getElementById('searchProjects');
    const projectsGrid = document.getElementById('projectsGrid');

    // Carregar projetos salvos ao carregar a página
    carregarProjetos();

    // Redirecionar para a página criar.html quando clicar no botão "+ novo"
    if (btnNewProject) {
        btnNewProject.addEventListener('click', function() {
            window.location.href = 'criar.html';
        });
    }

    // Redirecionar para a página criar.html quando clicar no card "+ novo"
    if (newProjectCard) {
        newProjectCard.addEventListener('click', function() {
            window.location.href = 'criar.html';
        });
    }

    // Funcionalidade de Editar Projeto
    document.addEventListener('click', function(e) {
        // Botão Editar
        if (e.target.classList.contains('btn__edit') || e.target.closest('.btn__edit')) {
            const card = e.target.closest('.card__project');
            const nome = card.querySelector('.project__nome');
            const descricao = card.querySelector('.project__descricao');
            const btnEdit = card.querySelector('.btn__edit');
            const btnSave = card.querySelector('.btn__save');

            // Tornar editável
            nome.contentEditable = true;
            descricao.contentEditable = true;
            nome.focus();

            // Trocar botões
            btnEdit.style.display = 'none';
            btnSave.style.display = 'inline-block';
        }

        // Botão Salvar
        if (e.target.classList.contains('btn__save') || e.target.closest('.btn__save')) {
            const card = e.target.closest('.card__project');
            const nome = card.querySelector('.project__nome');
            const descricao = card.querySelector('.project__descricao');
            const btnEdit = card.querySelector('.btn__edit');
            const btnSave = card.querySelector('.btn__save');
            const dataElement = card.querySelector('.project__data');

            // Desabilitar edição
            nome.contentEditable = false;
            descricao.contentEditable = false;

            // Atualizar data
            const hoje = new Date().toLocaleDateString('pt-BR');
            dataElement.textContent = `Última edição: ${hoje}`;

            // Trocar botões
            btnEdit.style.display = 'inline-block';
            btnSave.style.display = 'none';

            // Salvar no localStorage de projetos
            const projectId = card.dataset.projectId;
            let projetos = JSON.parse(localStorage.getItem('projetos')) || [];
            
            // Encontrar e atualizar o projeto
            const projetoIndex = projetos.findIndex(p => p.id === projectId);
            if (projetoIndex !== -1) {
                projetos[projetoIndex].nome = nome.textContent;
                projetos[projetoIndex].descricao = descricao.textContent;
                projetos[projetoIndex].dataEdicao = hoje;
                localStorage.setItem('projetos', JSON.stringify(projetos));
            }
            
            // Salvar no localStorage antigo também (para compatibilidade)
            const projects = JSON.parse(localStorage.getItem('projects')) || {};
            projects[projectId] = {
                nome: nome.textContent,
                descricao: descricao.textContent,
                dataEdicao: hoje
            };
            localStorage.setItem('projects', JSON.stringify(projects));

            alert('Projeto salvo com sucesso!');
        }

        // Botão Deletar
        if (e.target.classList.contains('btn__delete') || e.target.closest('.btn__delete')) {
            if (confirm('Tem certeza que deseja deletar este projeto?')) {
                const card = e.target.closest('.card__project');
                const projectId = card.dataset.projectId;

                // Remover do localStorage de projetos
                let projetos = JSON.parse(localStorage.getItem('projetos')) || [];
                projetos = projetos.filter(p => p.id !== projectId);
                localStorage.setItem('projetos', JSON.stringify(projetos));

                // Remover do localStorage antigo (para compatibilidade)
                const projects = JSON.parse(localStorage.getItem('projects')) || {};
                delete projects[projectId];
                localStorage.setItem('projects', JSON.stringify(projects));

                // Animação de remoção
                card.style.transform = 'scale(0)';
                card.style.opacity = '0';
                setTimeout(() => {
                    card.remove();
                }, 300);
            }
        }
    });

    // Funcionalidade de Busca
    if (searchInput) {
        searchInput.addEventListener('input', function(e) {
            const searchTerm = e.target.value.toLowerCase();
            const projectCards = document.querySelectorAll('.card__project');

            projectCards.forEach(card => {
                const nome = card.querySelector('.project__nome').textContent.toLowerCase();
                const descricao = card.querySelector('.project__descricao').textContent.toLowerCase();
                const id = card.querySelector('.project__id').textContent.toLowerCase();

                if (nome.includes(searchTerm) || descricao.includes(searchTerm) || id.includes(searchTerm)) {
                    card.style.display = 'flex';
                } else {
                    card.style.display = 'none';
                }
            });
        });
    }

    // Fechar modal ao clicar fora
    window.addEventListener('click', function(e) {
        if (e.target === modal) {
            modal.style.display = 'none';
        }
    });

    // Carregar projetos salvos do localStorage
    loadSavedProjects();
});

// Função para carregar projetos do localStorage e exibir na tela
function carregarProjetos() {
    const projetos = JSON.parse(localStorage.getItem('projetos')) || [];
    const projectsGrid = document.getElementById('projectsGrid');
    
    // Limpar projetos existentes (manter apenas o card de novo projeto)
    const existingProjects = projectsGrid.querySelectorAll('.card__project');
    existingProjects.forEach(card => card.remove());
    
    // Adicionar cada projeto salvo
    projetos.forEach(projeto => {
        criarCardProjeto(projeto);
    });
}

// Função para criar um card de projeto
function criarCardProjeto(projeto) {
    const projectsGrid = document.getElementById('projectsGrid');
    const newProjectCard = document.getElementById('newProjectCard');
    
    const cardHTML = `
        <div class="card__project" data-project-id="${projeto.id}">
            <div class="project__placeholder">
                ${projeto.imagem ? 
                    `<img src="${projeto.imagem}" alt="${projeto.nome}" style="width: 100%; height: 100%; object-fit: cover; border-radius: 8px;">` : 
                    'IA'
                }
            </div>
            <div class="project__info">
                <h3 class="project__nome" contenteditable="false">${projeto.nome}</h3>
                <p class="project__descricao" contenteditable="false">${projeto.descricao}</p>
                <div class="project__meta">
                    <span class="project__id">ID: ${projeto.id}</span>
                    <span class="project__data">Última edição: ${projeto.dataEdicao}</span>
                </div>
            </div>
            <div class="project__actions">
                <button class="btn__edit" title="Editar">✏️</button>
                <button class="btn__save" title="Salvar" style="display: none;">💾</button>
                <button class="btn__delete" title="Deletar">🗑️</button>
            </div>
        </div>
    `;
    
    // Inserir antes do card de novo projeto
    newProjectCard.insertAdjacentHTML('beforebegin', cardHTML);
}

// Função para carregar projetos salvos (mantida para compatibilidade)
function loadSavedProjects() {
    const projects = JSON.parse(localStorage.getItem('projects')) || {};
    
    Object.keys(projects).forEach(projectId => {
        const project = projects[projectId];
        const card = document.querySelector(`[data-project-id="${projectId}"]`);
        
        if (card) {
            card.querySelector('.project__nome').textContent = project.nome;
            card.querySelector('.project__descricao').textContent = project.descricao;
            if (project.dataEdicao) {
                card.querySelector('.project__data').textContent = `Última edição: ${project.dataEdicao}`;
            }
        }
    });
}