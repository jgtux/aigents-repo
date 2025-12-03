// criar.js - Script para a página de criação de agentes com URL de imagem

document.addEventListener('DOMContentLoaded', function() {
    const formCriar = document.getElementById('formCriarAgente');
    const urlImagem = document.getElementById('urlImagem');
    const previewContainer = document.getElementById('previewContainer');
    const imagePreview = document.getElementById('imagePreview');

    // Preview da imagem quando a URL é inserida
    urlImagem.addEventListener('input', function(e) {
        const url = e.target.value.trim();
        
        if (url) {
            // Validar se é uma URL válida
            try {
                new URL(url);
                
                // Mostrar preview
                imagePreview.src = url;
                previewContainer.style.display = 'flex';
                
                // Se a imagem falhar ao carregar, esconder o preview
                imagePreview.onerror = function() {
                    previewContainer.style.display = 'none';
                    alert('Não foi possível carregar a imagem. Verifique a URL.');
                };
                
            } catch (error) {
                previewContainer.style.display = 'none';
            }
        } else {
            previewContainer.style.display = 'none';
        }
    });

    // Submit do formulário
    formCriar.addEventListener('submit', function(e) {
        e.preventDefault();
        
        const nome = document.getElementById('nomeAgente').value;
        const categoria = document.getElementById('categoriaAgente').value;
        const descricao = document.getElementById('descricaoAgente').value;
        const imagemUrl = urlImagem.value.trim();
        
        // Validar URL da imagem se foi fornecida
        if (imagemUrl) {
            try {
                new URL(imagemUrl);
            } catch (error) {
                alert('Por favor, insira uma URL válida para a imagem.');
                return;
            }
        }
        
        salvarProjeto(nome, categoria, descricao, imagemUrl);
    });

    // Função para salvar o projeto
    function salvarProjeto(nome, categoria, descricao, imagemUrl) {
        // Gerar ID único
        const projectId = 'PRJ' + Date.now();
        const hoje = new Date().toLocaleDateString('pt-BR');
        
        // Criar objeto do projeto
        const novoProjeto = {
            id: projectId,
            nome: nome,
            categoria: categoria,
            descricao: descricao,
            dataCriacao: hoje,
            dataEdicao: hoje,
            imagem: imagemUrl || null // Salvar a URL ou null se não fornecida
        };
        
        // Buscar projetos existentes
        let projetos = JSON.parse(localStorage.getItem('projetos')) || [];
        
        // Adicionar novo projeto
        projetos.push(novoProjeto);
        
        // Salvar no localStorage
        localStorage.setItem('projetos', JSON.stringify(projetos));
        
        // Mostrar mensagem de sucesso
        alert('Agente criado com sucesso!');
        
        // Redirecionar para My Projects
        window.location.href = 'Myprojects.html';
    }

    // Validação em tempo real
    const inputs = document.querySelectorAll('.input__criar, .textarea__criar');
    
    inputs.forEach(input => {
        input.addEventListener('blur', function() {
            if (this.value.trim() === '' && this.hasAttribute('required')) {
                this.style.borderColor = '#dc3545';
            } else {
                this.style.borderColor = '';
            }
        });
        
        input.addEventListener('input', function() {
            if (this.style.borderColor === 'rgb(220, 53, 69)') {
                this.style.borderColor = '';
            }
        });
    });
});s