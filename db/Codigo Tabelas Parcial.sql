-- Create ENUM types first
CREATE TYPE user_type_enum AS ENUM ('USER', 'ADMIN');

-- Tabela da autenticacao
CREATE TABLE auths (
    auth_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);

-- Tabela dados usuario
CREATE TABLE users (
    user_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    auth_uuid UUID NOT NULL UNIQUE, -- UNIQUE adicionado
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    id_document CHAR(14) NOT NULL UNIQUE,
    user_type user_type_enum DEFAULT 'USER',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (auth_uuid) REFERENCES auths(auth_uuid) ON DELETE CASCADE
);

-- Tabela sistemas para os agentes
CREATE TABLE agent_systems (
    agent_system_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_system_preset JSONB NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabela categoria dos agentes
CREATE TABLE agent_categories (
    category_id SERIAL PRIMARY KEY,
    category_name VARCHAR(32),
    agent_system_uuid_preset UUID NOT NULL UNIQUE, -- UNIQUE adicionado
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (agent_system_uuid_preset) REFERENCES agent_systems(agent_system_uuid)
);

-- Configuracao dos agentes
CREATE TABLE agents_config (
    agent_config_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id INT NOT NULL,
    category_preset_enabled BOOL DEFAULT TRUE,
    agent_system_uuid UUID NOT NULL UNIQUE, -- UNIQUE adicionado
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES agent_categories(category_id),
    FOREIGN KEY (agent_system_uuid) REFERENCES agent_systems(agent_system_uuid)
);
  
-- Tabela dos agentes
CREATE TABLE agents (
    agent_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description VARCHAR(300),
    agent_config_uuid UUID NOT NULL UNIQUE, -- UNIQUE adicionado
    creator_uuid UUID NOT NULL,
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (agent_config_uuid) REFERENCES agents_config(agent_config_uuid) ON DELETE CASCADE,
    FOREIGN KEY (creator_uuid) REFERENCES users(user_uuid) ON DELETE CASCADE
);

-- Tabela das avaliacoes
CREATE TABLE reviews (
    review_id SERIAL PRIMARY KEY,
    rating DECIMAL(2,1) NOT NULL,
    comment VARCHAR(300),
    user_uuid UUID NOT NULL,
    agent_uuid UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_uuid) REFERENCES users(user_uuid) ON DELETE CASCADE,
    FOREIGN KEY (agent_uuid) REFERENCES agents(agent_uuid) ON DELETE CASCADE
);

-- Tabela para salvar o chat
CREATE TABLE chats (
    chat_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent_uuid UUID NOT NULL,
    user_uuid UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (user_uuid) REFERENCES users(user_uuid) ON DELETE CASCADE,
    FOREIGN KEY (agent_uuid) REFERENCES agents(agent_uuid) ON DELETE CASCADE
);

-- Tabela para o conteúdo das mensagens
CREATE TABLE message_contents (
    message_content_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_content TEXT NOT NULL
);

-- Tabela mensagem
CREATE TABLE messages (
    message_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_uuid UUID NOT NULL,
    receiver_uuid UUID NOT NULL,
    chat_uuid UUID NOT NULL,
    message_content_uuid UUID NOT NULL UNIQUE,	    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_uuid) REFERENCES users(user_uuid) ON DELETE CASCADE,
    FOREIGN KEY (receiver_uuid) REFERENCES users(user_uuid) ON DELETE CASCADE,
    FOREIGN KEY (chat_uuid) REFERENCES chats(chat_uuid) ON DELETE CASCADE,
    FOREIGN KEY (message_content_uuid) REFERENCES message_contents(message_content_uuid) ON DELETE CASCADE
);