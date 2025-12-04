------------------------------------------------------------
-- POPULAÇÃO INICIAL DO BANCO
-- Somente usuários USER
------------------------------------------------------------

-------------------------------
-- 1. auths
-------------------------------
INSERT INTO auths (auth_uuid, email, password, role) VALUES
('11111111-1111-1111-1111-111111111111', 'ana@example.com', 'senha123', 'USER'),
('22222222-2222-2222-2222-222222222222', 'bruno@example.com', 'senha123', 'USER'),
('33333333-3333-3333-3333-333333333333', 'carla@example.com', 'senha123', 'USER'),
('44444444-4444-4444-4444-444444444444', 'diego@example.com', 'senha123', 'USER'),
('55555555-5555-5555-5555-555555555555', 'erika@example.com', 'senha123', 'USER');


-------------------------------
-- 2. agent_systems
-------------------------------
INSERT INTO agent_systems (agent_system_uuid, system_preset) VALUES
('aaaaaaa1-aaaa-aaaa-aaaa-aaaaaaaaaaa1', '{"persona": "Especialista em marketing digital", "temperature": 0.6}'),
('aaaaaaa2-aaaa-aaaa-aaaa-aaaaaaaaaaa2', '{"persona": "Programador full-stack", "temperature": 0.2}'),
('aaaaaaa3-aaaa-aaaa-aaaa-aaaaaaaaaaa3', '{"persona": "Nutricionista virtual", "temperature": 0.8}'),
('aaaaaaa4-aaaa-aaaa-aaaa-aaaaaaaaaaa4', '{"persona": "Professor de matemática", "temperature": 0.4}'),
('aaaaaaa5-aaaa-aaaa-aaaa-aaaaaaaaaaa5', '{"persona": "Consultor jurídico", "temperature": 0.3}');


-------------------------------
-- 3. agent_categories
-------------------------------
INSERT INTO agent_categories (category_id, category_name, agent_system_uuid_preset) VALUES
(1, 'Marketing',    'aaaaaaa1-aaaa-aaaa-aaaa-aaaaaaaaaaa1'),
(2, 'Programação',  'aaaaaaa2-aaaa-aaaa-aaaa-aaaaaaaaaaa2'),
(3, 'Saúde',        'aaaaaaa3-aaaa-aaaa-aaaa-aaaaaaaaaaa3'),
(4, 'Educação',     'aaaaaaa4-aaaa-aaaa-aaaa-aaaaaaaaaaa4'),
(5, 'Jurídico',     'aaaaaaa5-aaaa-aaaa-aaaa-aaaaaaaaaaa5');


-------------------------------
-- 4. agents_config
-------------------------------
INSERT INTO agents_config (agent_config_uuid, category_id, agent_system_uuid) VALUES
('bbbbbbb1-bbbb-bbbb-bbbb-bbbbbbbbbbb1', 1, 'aaaaaaa1-aaaa-aaaa-aaaa-aaaaaaaaaaa1'),
('bbbbbbb2-bbbb-bbbb-bbbb-bbbbbbbbbbb2', 2, 'aaaaaaa2-aaaa-aaaa-aaaa-aaaaaaaaaaa2'),
('bbbbbbb3-bbbb-bbbb-bbbb-bbbbbbbbbbb3', 3, 'aaaaaaa3-aaaa-aaaa-aaaa-aaaaaaaaaaa3'),
('bbbbbbb4-bbbb-bbbb-bbbb-bbbbbbbbbbb4', 4, 'aaaaaaa4-aaaa-aaaa-aaaa-aaaaaaaaaaa4'),
('bbbbbbb5-bbbb-bbbb-bbbb-bbbbbbbbbbb5', 5, 'aaaaaaa5-aaaa-aaaa-aaaa-aaaaaaaaaaa5');


-------------------------------
-- 5. agents
-------------------------------
INSERT INTO agents (
  agent_uuid, name, description, image_url,
  agent_config_uuid, auth_uuid
)
VALUES
('ccccccc1-cccc-cccc-cccc-ccccccccccc1',
 'SocialBoost AI',
 'Especialista em marketing digital, criação de posts e estratégias de engajamento.',
 'https://example.com/agents/marketing.png',
 'bbbbbbb1-bbbb-bbbb-bbbb-bbbbbbbbbbb1',
 '22222222-2222-2222-2222-222222222222'),

('ccccccc2-cccc-cccc-cccc-ccccccccccc2',
 'DevHelper Pro',
 'Assistente de programação full-stack, ajuda com código, melhorias e bugs.',
 'https://example.com/agents/dev.png',
 'bbbbbbb2-bbbb-bbbb-bbbb-bbbbbbbbbbb2',
 '33333333-3333-3333-3333-333333333333'),

('ccccccc3-cccc-cccc-cccc-ccccccccccc3',
 'NutriVida AI',
 'Nutricionista virtual que recomenda dietas e orientações de saúde.',
 'https://example.com/agents/nutri.png',
 'bbbbbbb3-bbbb-bbbb-bbbb-bbbbbbbbbbb3',
 '33333333-3333-3333-3333-333333333333'),

('ccccccc4-cccc-cccc-cccc-ccccccccccc4',
 'MathTutor AI',
 'Professor de matemática especializado em exercícios passo a passo.',
 'https://example.com/agents/math.png',
 'bbbbbbb4-bbbb-bbbb-bbbb-bbbbbbbbbbb4',
 '22222222-2222-2222-2222-222222222222'),

('ccccccc5-cccc-cccc-cccc-ccccccccccc5',
 'LegalMind',
 'Consultor jurídico com foco em direitos do consumidor.',
 'https://example.com/agents/legal.png',
 'bbbbbbb5-bbbb-bbbb-bbbb-bbbbbbbbbbb5',
 '22222222-2222-2222-2222-222222222222');


-------------------------------
-- 6. chats
-------------------------------
INSERT INTO chats (chat_uuid, agent_uuid, auth_uuid) VALUES
('ddddddd1-dddd-dddd-dddd-ddddddddddd1', 'ccccccc1-cccc-cccc-cccc-ccccccccccc1', '11111111-1111-1111-1111-111111111111'),
('ddddddd2-dddd-dddd-dddd-ddddddddddd2', 'ccccccc2-cccc-cccc-cccc-ccccccccccc2', '44444444-4444-4444-4444-444444444444'),
('ddddddd3-dddd-dddd-dddd-ddddddddddd3', 'ccccccc3-cccc-cccc-cccc-ccccccccccc3', '55555555-5555-5555-5555-555555555555'),
('ddddddd4-dddd-dddd-dddd-ddddddddddd4', 'ccccccc4-cccc-cccc-cccc-ccccccccccc4', '11111111-1111-1111-1111-111111111111'),
('ddddddd5-dddd-dddd-dddd-ddddddddddd5', 'ccccccc5-cccc-cccc-cccc-ccccccccccc5', '44444444-4444-4444-4444-444444444444');


-------------------------------
-- 7. message_contents
-------------------------------
INSERT INTO message_contents (message_content_uuid, message_content) VALUES
('eeeeeee1-eeee-eeee-eeee-eeeeeeeeeee1', 'Oi! Pode me ajudar com ideias de posts?'),
('eeeeeee2-eeee-eeee-eeee-eeeeeeeeeee2', 'Claro! Que tipo de negócio você tem?'),
('eeeeeee3-eeee-eeee-eeee-eeeeeeeeeee3', 'Estou com erro no meu código JavaScript.'),
('eeeeeee4-eeee-eeee-eeee-eeeeeeeeeee4', 'Envie o código que analiso para você.'),
('eeeeeee5-eeee-eeee-eeee-eeeeeeeeeee5', 'Quero melhorar minha alimentação.'),
('eeeeeee6-eeee-eeee-eeee-eeeeeeeeeee6', 'Perfeito! Quais são seus objetivos?'),
('eeeeeee7-eeee-eeee-eeee-eeeeeeeeeee7', 'Pode me explicar derivadas?'),
('eeeeeee8-eeee-eeee-eeee-eeeeeeeeeee8', 'Claro, vamos começar pelo conceito básico.'),
('eeeeeee9-eeee-eeee-eeee-eeeeeeeeeee9', 'Tenho um problema com uma loja online.'),
('eeeeee10-eeee-eeee-eeee-eeeeeeeeeee0', 'Me conte o que aconteceu e eu oriento você.');


-------------------------------
-- 8. messages
-------------------------------
INSERT INTO messages (
  message_uuid, sender_uuid, sender_type,
  receiver_uuid, receiver_type,
  chat_uuid, message_content_uuid
) VALUES
-- CHAT 1 — Marketing
('fffffff1-ffff-ffff-ffff-fffffffffff1',
 '11111111-1111-1111-1111-111111111111', 'AUTH',
 'ccccccc1-cccc-cccc-cccc-ccccccccccc1', 'AGENT',
 'ddddddd1-dddd-dddd-dddd-ddddddddddd1',
 'eeeeeee1-eeee-eeee-eeee-eeeeeeeeeee1'),

('fffffff2-ffff-ffff-ffff-fffffffffff2',
 'ccccccc1-cccc-cccc-cccc-ccccccccccc1', 'AGENT',
 '11111111-1111-1111-1111-111111111111', 'AUTH',
 'ddddddd1-dddd-dddd-dddd-ddddddddddd1',
 'eeeeeee2-eeee-eeee-eeee-eeeeeeeeeee2'),

-- CHAT 2 — Programação
('fffffff3-ffff-ffff-ffff-fffffffffff3',
 '44444444-4444-4444-4444-444444444444', 'AUTH',
 'ccccccc2-cccc-cccc-cccc-ccccccccccc2', 'AGENT',
 'ddddddd2-dddd-dddd-dddd-ddddddddddd2',
 'eeeeeee3-eeee-eeee-eeee-eeeeeeeeeee3'),

('fffffff4-ffff-ffff-ffff-fffffffffff4',
 'ccccccc2-cccc-cccc-cccc-ccccccccccc2', 'AGENT',
 '44444444-4444-4444-4444-444444444444', 'AUTH',
 'ddddddd2-dddd-dddd-dddd-ddddddddddd2',
 'eeeeeee4-eeee-eeee-eeee-eeeeeeeeeee4'),

-- CHAT 3 — Nutrição
('fffffff5-ffff-ffff-ffff-fffffffffff5',
 '55555555-5555-5555-5555-555555555555', 'AUTH',
 'ccccccc3-cccc-cccc-cccc-ccccccccccc3', 'AGENT',
 'ddddddd3-dddd-dddd-dddd-ddddddddddd3',
 'eeeeeee5-eeee-eeee-eeee-eeeeeeeeeee5'),

('fffffff6-ffff-ffff-ffff-fffffffffff6',
 'ccccccc3-cccc-cccc-cccc-ccccccccccc3', 'AGENT',
 '55555555-5555-5555-5555-555555555555', 'AUTH',
 'ddddddd3-dddd-dddd-dddd-ddddddddddd3',
 'eeeeeee6-eeee-eeee-eeee-eeeeeeeeeee6'),

-- CHAT 4 — Matemática
('fffffff7-ffff-ffff-ffff-fffffffffff7',
 '11111111-1111-1111-1111-111111111111', 'AUTH',
 'ccccccc4-cccc-cccc-cccc-ccccccccccc4', 'AGENT',
 'ddddddd4-dddd-dddd-dddd-ddddddddddd4',
 'eeeeeee7-eeee-eeee-eeee-eeeeeeeeeee7'),

('fffffff8-ffff-ffff-ffff-fffffffffff8',
 'ccccccc4-cccc-cccc-cccc-ccccccccccc4', 'AGENT',
 '11111111-1111-1111-1111-111111111111', 'AUTH',
 'ddddddd4-dddd-dddd-dddd-ddddddddddd4',
 'eeeeeee8-eeee-eeee-eeee-eeeeeeeeeee8'),

-- CHAT 5 — Jurídico
('fffffff9-ffff-ffff-ffff-fffffffffff9',
 '44444444-4444-4444-4444-444444444444', 'AUTH',
 'ccccccc5-cccc-cccc-cccc-ccccccccccc5', 'AGENT',
 'ddddddd5-dddd-dddd-dddd-ddddddddddd5',
 'eeeeeee9-eeee-eeee-eeee-eeeeeeeeeee9'),

('fffffff0-ffff-ffff-ffff-fffffffffff0',
 'ccccccc5-cccc-cccc-cccc-ccccccccccc5', 'AGENT',
 '44444444-4444-4444-4444-444444444444', 'AUTH',
 'ddddddd5-dddd-dddd-dddd-ddddddddddd5',
 'eeeeee10-eeee-eeee-eeee-eeeeeeeeeee0');
