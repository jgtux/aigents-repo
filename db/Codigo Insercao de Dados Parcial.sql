-- 1. Inserir dados na tabela auths
INSERT INTO auths (email, password) VALUES
('maria.silva@email.com', 'hashed_password_123'),
('joao.costa@email.com', 'hashed_password_456'),
('ana.santos@email.com', 'hashed_password_789'),
('carlos.oliveira@email.com', 'hashed_password_101'),
('admin@empresa.com', 'hashed_password_admin');

-- 2. Inserir dados na tabela users (COM CAST para user_type_enum)
INSERT INTO users (auth_uuid, first_name, last_name, id_document, user_type) 
SELECT auth_uuid, 'Maria', 'Silva', '123.456.789-00', 'USER'::user_type_enum FROM auths WHERE email = 'maria.silva@email.com'
UNION ALL
SELECT auth_uuid, 'João', 'Costa', '234.567.890-11', 'USER'::user_type_enum FROM auths WHERE email = 'joao.costa@email.com'
UNION ALL
SELECT auth_uuid, 'Ana', 'Santos', '345.678.901-22', 'USER'::user_type_enum FROM auths WHERE email = 'ana.santos@email.com'
UNION ALL
SELECT auth_uuid, 'Carlos', 'Oliveira', '456.789.012-33', 'USER'::user_type_enum FROM auths WHERE email = 'carlos.oliveira@email.com'
UNION ALL
SELECT auth_uuid, 'Admin', 'Sistema', '567.890.123-44', 'ADMIN'::user_type_enum FROM auths WHERE email = 'admin@empresa.com';

-- 3. Inserir dados na tabela agent_systems
INSERT INTO agent_systems (category_system_preset) VALUES
('{"model": "gpt-4", "temperature": 0.7, "max_tokens": 2000}'),
('{"model": "claude-2", "temperature": 0.5, "max_tokens": 1500}'),
('{"model": "llama-2", "temperature": 0.8, "max_tokens": 1000}'),
('{"model": "gemini-pro", "temperature": 0.6, "max_tokens": 1800}'),
('{"model": "mistral", "temperature": 0.4, "max_tokens": 1200}');

-- 4. Inserir dados na tabela agent_categories
INSERT INTO agent_categories (category_name, agent_system_uuid_preset) 
SELECT 'ODS 8 - Trabalho Decente', agent_system_uuid FROM agent_systems WHERE category_system_preset::jsonb->>'model' = 'gpt-4'
UNION ALL
SELECT 'ODS 8 - Crescimento Econômico', agent_system_uuid FROM agent_systems WHERE category_system_preset::jsonb->>'model' = 'claude-2'
UNION ALL
SELECT 'ODS 9 - Indústria e Inovação', agent_system_uuid FROM agent_systems WHERE category_system_preset::jsonb->>'model' = 'llama-2'
UNION ALL
SELECT 'ODS 9 - Infraestrutura', agent_system_uuid FROM agent_systems WHERE category_system_preset::jsonb->>'model' = 'gemini-pro'
UNION ALL
SELECT 'Finanças Corporativas', agent_system_uuid FROM agent_systems WHERE category_system_preset::jsonb->>'model' = 'mistral';

-- 5. Inserir dados na tabela agents_config
INSERT INTO agents_config (category_id, agent_system_uuid, category_preset_enabled)
SELECT category_id, agent_system_uuid_preset, TRUE FROM agent_categories WHERE category_name = 'ODS 8 - Trabalho Decente'
UNION ALL
SELECT category_id, agent_system_uuid_preset, TRUE FROM agent_categories WHERE category_name = 'ODS 8 - Crescimento Econômico'
UNION ALL
SELECT category_id, agent_system_uuid_preset, FALSE FROM agent_categories WHERE category_name = 'ODS 9 - Indústria e Inovação'
UNION ALL
SELECT category_id, agent_system_uuid_preset, TRUE FROM agent_categories WHERE category_name = 'ODS 9 - Infraestrutura'
UNION ALL
SELECT category_id, agent_system_uuid_preset, TRUE FROM agent_categories WHERE category_name = 'Finanças Corporativas';

-- 6. Inserir dados na tabela agents
INSERT INTO agents (name, description, agent_config_uuid, creator_uuid)
SELECT 
    'TrabalhoDecente AI',
    'Agente para ODS 8 - Promoção do trabalho decente e direitos trabalhistas',
    ac.agent_config_uuid,
    u.user_uuid
FROM agents_config ac, users u 
WHERE u.first_name = 'Maria' 
AND ac.category_id IN (SELECT category_id FROM agent_categories WHERE category_name = 'ODS 8 - Trabalho Decente')
UNION ALL
SELECT 
    'CrescimentoEco Advisor',
    'Agente para ODS 8 - Crescimento econômico sustentável e produtividade',
    ac.agent_config_uuid,
    u.user_uuid
FROM agents_config ac, users u 
WHERE u.first_name = 'João'
AND ac.category_id IN (SELECT category_id FROM agent_categories WHERE category_name = 'ODS 8 - Crescimento Econômico')
UNION ALL
SELECT 
    'InovaIndustria AI',
    'Agente para ODS 9 - Inovação industrial e modernização tecnológica',
    ac.agent_config_uuid,
    u.user_uuid
FROM agents_config ac, users u 
WHERE u.first_name = 'Ana'
AND ac.category_id IN (SELECT category_id FROM agent_categories WHERE category_name = 'ODS 9 - Indústria e Inovação')
UNION ALL
SELECT 
    'InfraTech Pro',
    'Agente para ODS 9 - Infraestrutura resiliente e desenvolvimento sustentável',
    ac.agent_config_uuid,
    u.user_uuid
FROM agents_config ac, users u 
WHERE u.first_name = 'Carlos'
AND ac.category_id IN (SELECT category_id FROM agent_categories WHERE category_name = 'ODS 9 - Infraestrutura')
UNION ALL
SELECT 
    'Finance Corp Analyst',
    'Agente para análise financeira corporativa e gestão empresarial',
    ac.agent_config_uuid,
    u.user_uuid
FROM agents_config ac, users u 
WHERE u.first_name = 'Admin'
AND ac.category_id IN (SELECT category_id FROM agent_categories WHERE category_name = 'Finanças Corporativas');

-- 7. Inserir dados na tabela reviews
INSERT INTO reviews (rating, comment, user_uuid, agent_uuid)
SELECT 4.5, 'Excelente para políticas de trabalho decente', u.user_uuid, a.agent_uuid
FROM users u, agents a WHERE u.first_name = 'João' AND a.name = 'TrabalhoDecente AI'
UNION ALL
SELECT 5.0, 'Fundamental para nosso planejamento de crescimento', u.user_uuid, a.agent_uuid
FROM users u, agents a WHERE u.first_name = 'Ana' AND a.name = 'CrescimentoEco Advisor'
UNION ALL
SELECT 4.0, 'Muito útil para inovação na nossa indústria', u.user_uuid, a.agent_uuid
FROM users u, agents a WHERE u.first_name = 'Carlos' AND a.name = 'InovaIndustria AI'
UNION ALL
SELECT 4.8, 'Essencial para projetos de infraestrutura sustentável', u.user_uuid, a.agent_uuid
FROM users u, agents a WHERE u.first_name = 'Maria' AND a.name = 'InfraTech Pro'
UNION ALL
SELECT 4.7, 'Preciso nas análises financeiras da nossa empresa', u.user_uuid, a.agent_uuid
FROM users u, agents a WHERE u.first_name = 'Admin' AND a.name = 'Finance Corp Analyst';

-- 8. Inserir dados na tabela chats
INSERT INTO chats (agent_uuid, user_uuid)
SELECT a.agent_uuid, u.user_uuid FROM agents a, users u WHERE a.name = 'TrabalhoDecente AI' AND u.first_name = 'João'
UNION ALL
SELECT a.agent_uuid, u.user_uuid FROM agents a, users u WHERE a.name = 'CrescimentoEco Advisor' AND u.first_name = 'Ana'
UNION ALL
SELECT a.agent_uuid, u.user_uuid FROM agents a, users u WHERE a.name = 'InovaIndustria AI' AND u.first_name = 'Carlos'
UNION ALL
SELECT a.agent_uuid, u.user_uuid FROM agents a, users u WHERE a.name = 'InfraTech Pro' AND u.first_name = 'Maria'
UNION ALL
SELECT a.agent_uuid, u.user_uuid FROM agents a, users u WHERE a.name = 'Finance Corp Analyst' AND u.first_name = 'Admin';

-- 9. Inserir dados na tabela message_contents
INSERT INTO message_contents (message_content) VALUES
('Como implementar políticas de trabalho decente na nossa empresa?'),
('Podemos começar com direitos trabalhistas, segurança no trabalho e igualdade de oportunidades. Qual o setor da sua empresa?'),
('Preciso de estratégias para crescimento econômico sustentável'),
('Vamos analisar indicadores econômicos e desenvolver um plano baseado em inovação e produtividade'),
('Como modernizar nossa linha de produção industrial?');

-- 10. Inserir dados na tabela messages
INSERT INTO messages (sender_uuid, receiver_uuid, chat_uuid, message_content_uuid)
SELECT 
    u.user_uuid, 
    a.creator_uuid, 
    c.chat_uuid, 
    mc.message_content_uuid
FROM users u, agents a, chats c, message_contents mc
WHERE u.first_name = 'João' AND a.name = 'TrabalhoDecente AI' AND c.user_uuid = u.user_uuid
AND mc.message_content = 'Como implementar políticas de trabalho decente na nossa empresa?'
UNION ALL
SELECT 
    a.creator_uuid,
    u.user_uuid,
    c.chat_uuid,
    mc.message_content_uuid
FROM users u, agents a, chats c, message_contents mc
WHERE u.first_name = 'João' AND a.name = 'TrabalhoDecente AI' AND c.user_uuid = u.user_uuid
AND mc.message_content = 'Podemos começar com direitos trabalhistas, segurança no trabalho e igualdade de oportunidades. Qual o setor da sua empresa?'
UNION ALL
SELECT 
    u.user_uuid,
    a.creator_uuid,
    c.chat_uuid,
    mc.message_content_uuid
FROM users u, agents a, chats c, message_contents mc
WHERE u.first_name = 'Ana' AND a.name = 'CrescimentoEco Advisor' AND c.user_uuid = u.user_uuid
AND mc.message_content = 'Preciso de estratégias para crescimento econômico sustentável'
UNION ALL
SELECT 
    a.creator_uuid,
    u.user_uuid,
    c.chat_uuid,
    mc.message_content_uuid
FROM users u, agents a, chats c, message_contents mc
WHERE u.first_name = 'Ana' AND a.name = 'CrescimentoEco Advisor' AND c.user_uuid = u.user_uuid
AND mc.message_content = 'Vamos analisar indicadores econômicos e desenvolver um plano baseado em inovação e produtividade'
UNION ALL
SELECT 
    u.user_uuid,
    a.creator_uuid,
    c.chat_uuid,
    mc.message_content_uuid
FROM users u, agents a, chats c, message_contents mc
WHERE u.first_name = 'Carlos' AND a.name = 'InovaIndustria AI' AND c.user_uuid = u.user_uuid
AND mc.message_content = 'Como modernizar nossa linha de produção industrial?';