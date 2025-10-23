-- ============================================================
-- ENUM para papéis de autenticação
-- ============================================================
CREATE TYPE role_enum AS ENUM ('USER', 'ADMIN', 'AGENT');

-- ============================================================
-- Função e trigger para atualizar o campo updated_at
-- ============================================================
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ============================================================
-- Tabela de autenticação
-- ============================================================
CREATE TABLE auths (
  auth_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email VARCHAR(255) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  role role_enum DEFAULT 'USER',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP DEFAULT NULL
);

-- ============================================================
-- Tabela de usuários
-- ============================================================
CREATE TABLE users (
  user_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  auth_uuid UUID NOT NULL UNIQUE,
  first_name VARCHAR(128) NOT NULL,
  last_name VARCHAR(128) NOT NULL,
  id_document CHAR(14) NOT NULL UNIQUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP DEFAULT NULL,
  FOREIGN KEY (auth_uuid) REFERENCES auths(auth_uuid) ON DELETE CASCADE
);

-- ============================================================
-- Tabela de sistemas dos agentes
-- ============================================================
CREATE TABLE agent_systems (
  agent_system_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  category_system_preset JSONB NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================
-- Tabela de categorias de agentes
-- ============================================================
CREATE TABLE agent_categories (
  category_id SERIAL PRIMARY KEY,
  category_name VARCHAR(32) NOT NULL,
  agent_system_uuid_preset UUID NOT NULL UNIQUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP DEFAULT NULL,
  FOREIGN KEY (agent_system_uuid_preset) REFERENCES agent_systems(agent_system_uuid)
);

-- ============================================================
-- Tabela de configuração dos agentes
-- ============================================================
CREATE TABLE agents_config (
  agent_config_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  category_id INT NOT NULL,
  category_preset_enabled BOOL DEFAULT TRUE,
  agent_system_uuid UUID NOT NULL UNIQUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (category_id) REFERENCES agent_categories(category_id),
  FOREIGN KEY (agent_system_uuid) REFERENCES agent_systems(agent_system_uuid)
);

-- ============================================================
-- Tabela de agentes
-- ============================================================
CREATE TABLE agents (
  agent_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(100) NOT NULL,
  description VARCHAR(512),
  agent_config_uuid UUID NOT NULL UNIQUE,
  creator_uuid UUID NOT NULL,
  creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (agent_config_uuid) REFERENCES agents_config(agent_config_uuid) ON DELETE CASCADE,
  FOREIGN KEY (creator_uuid) REFERENCES users(user_uuid) ON DELETE CASCADE
);

-- ============================================================
-- Tabela de avaliações
-- ============================================================
CREATE TABLE reviews (
  review_id SERIAL PRIMARY KEY,
  rating DECIMAL(2,1) NOT NULL CHECK (rating >= 0 AND rating <= 5),
  comment VARCHAR(512),
  user_uuid UUID NOT NULL,
  agent_uuid UUID NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_uuid) REFERENCES users(user_uuid) ON DELETE CASCADE,
  FOREIGN KEY (agent_uuid) REFERENCES agents(agent_uuid) ON DELETE CASCADE
);

-- ============================================================
-- Tabela de chats
-- ============================================================
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

-- ============================================================
-- Tabela de conteúdos das mensagens
-- ============================================================
CREATE TABLE message_contents (
  message_content_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  message_content TEXT NOT NULL
);

-- ============================================================
-- Tabela de mensagens
-- ============================================================
CREATE TABLE messages (
  message_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  sender_uuid UUID NOT NULL,
  receiver_uuid UUID NOT NULL,
  chat_uuid UUID NOT NULL,
  message_content_uuid UUID NOT NULL UNIQUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (chat_uuid) REFERENCES chats(chat_uuid) ON DELETE CASCADE,
  FOREIGN KEY (message_content_uuid) REFERENCES message_contents(message_content_uuid) ON DELETE CASCADE
);

-- ============================================================
-- ÍNDICES PARA OTIMIZAÇÃO
-- ============================================================

-- Usuários
CREATE INDEX idx_users_auth_uuid ON users(auth_uuid);

-- Agentes
CREATE INDEX idx_agents_creator_uuid ON agents(creator_uuid);
CREATE INDEX idx_agents_config_uuid ON agents(agent_config_uuid);

-- Avaliações
CREATE INDEX idx_reviews_agent_uuid ON reviews(agent_uuid);
CREATE INDEX idx_reviews_user_uuid ON reviews(user_uuid);

-- Chats
CREATE INDEX idx_chats_user_uuid ON chats(user_uuid);
CREATE INDEX idx_chats_agent_uuid ON chats(agent_uuid);

-- Mensagens
CREATE INDEX idx_messages_chat_uuid ON messages(chat_uuid);
CREATE INDEX idx_messages_sender_uuid ON messages(sender_uuid);

-- ============================================================
-- TRIGGERS PARA updated_at
-- ============================================================

-- auths
CREATE TRIGGER trg_auths_updated
BEFORE UPDATE ON auths
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- users
CREATE TRIGGER trg_users_updated
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- agent_systems
CREATE TRIGGER trg_agent_systems_updated
BEFORE UPDATE ON agent_systems
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- agent_categories
CREATE TRIGGER trg_agent_categories_updated
BEFORE UPDATE ON agent_categories
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- agents_config
CREATE TRIGGER trg_agents_config_updated
BEFORE UPDATE ON agents_config
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- agents
CREATE TRIGGER trg_agents_updated
BEFORE UPDATE ON agents
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- reviews
CREATE TRIGGER trg_reviews_updated
BEFORE UPDATE ON reviews
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- chats
CREATE TRIGGER trg_chats_updated
BEFORE UPDATE ON chats
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();
