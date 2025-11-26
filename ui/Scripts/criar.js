// criar.js - Script para a página de criação de agentes

document.addEventListener('DOMContentLoaded', function() {
    const formCriar = document.getElementById('formCriarAgente');
    const inputImagem = document.getElementById('imagemAgente');
    const uploadArea = document.getElementById('uploadArea');
    const previewContainer = document.getElementById('previewContainer');
    const imagePreview = document.getElementById('imagePreview');
    const btnRemoveImage = document.getElementById('btnRemoveImage');
    const uploadContent = uploadArea.querySelector('.upload__content');

    // Preview da imagem quando selecionada
    inputImagem.addEventListener('change', function(e) {
        const file = e.target.files[0];
        
        if (file && file.type.startsWith('image/')) {
            const reader = new FileReader();
            
            reader.onload = function(e) {
                imagePreview.src = e.target.result;
                uploadContent.style.display = 'none';
                previewContainer.style.display = 'flex';
            };
            
            reader.readAsDataURL(file);
        }
    });

    // Remover imagem
    btnRemoveImage.addEventListener('click', function(e) {
        e.stopPropagation();
        inputImagem.value = '';
        imagePreview.src = '';
        uploadContent.style.display = 'flex';
        previewContainer.style.display = 'none';
    });

    // Drag and drop para upload de imagem
    uploadArea.addEventListener('dragover', function(e) {
        e.preventDefault();
        uploadArea.style.borderColor = 'var(--cor-botao)';
        uploadArea.style.backgroundColor = 'rgba(127, 140, 170, 0.5)';
    });

    uploadArea.addEventListener('dragleave', function(e) {
        e.preventDefault();
        uploadArea.style.borderColor = '';
        uploadArea.style.backgroundColor = '';
    });

    uploadArea.addEventListener('drop', function(e) {
        e.preventDefault();
        uploadArea.style.borderColor = '';
        uploadArea.style.backgroundColor = '';
        
        const file = e.dataTransfer.files[0];
        
        if (file && file.type.startsWith('image/')) {
            // Atualizar o input file
            const dataTransfer = new DataTransfer();
            dataTransfer.items.add(file);
            inputImagem.files = dataTransfer.files;
            
            // Trigger change event
            const event = new Event('change', { bubbles: true });
            inputImagem.dispatchEvent(event);
        }
    });

    // Submit do formulário
    formCriar.addEventListener('submit', function(e) {
        e.preventDefault();
        
        const nome = document.getElementById('nomeAgente').value;
        const categoria = document.getElementById('categoriaAgente').value;
        const descricao = document.getElementById('descricaoAgente').value;
        const imagem = inputImagem.files[0];
        
        // Converter imagem para base64 para salvar no localStorage
        if (imagem) {
            const reader = new FileReader();
            reader.onload = function(e) {
                salvarProjeto(nome, categoria, descricao, e.target.result);
            };
            reader.readAsDataURL(imagem);
        } else {
            salvarProjeto(nome, categoria, descricao, null);
        }
    });

    // Função para salvar o projeto
    function salvarProjeto(nome, categoria, descricao, imagemBase64) {
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
            imagem: imagemBase64
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
});