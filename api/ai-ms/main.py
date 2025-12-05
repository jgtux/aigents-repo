import uuid
import time
import os

from collections import OrderedDict
from typing import Dict, List, Optional
import asyncio
import json
import websockets
from dotenv import load_dotenv

from langchain_groq import ChatGroq
from langchain_core.messages import SystemMessage, HumanMessage, AIMessage
from langchain_core.callbacks import AsyncCallbackHandler

from websockets.exceptions import ConnectionClosed
from typing import Set
import contextlib

# ==========================
# Load Environment Variables
# ==========================
load_dotenv()

# ==========================
# Config
# ==========================


HEARTBEAT_INTERVAL = int(os.getenv("HEARTBEAT_INTERVAL", "30"))
CONNECTION_TIMEOUT = int(os.getenv("CONNECTION_TIMEOUT", "300"))

STREAM_MIN_CHUNK_SIZE = int(os.getenv("STREAM_MIN_CHUNK_SIZE", "50"))
STREAM_MAX_DELAY = float(os.getenv("STREAM_MAX_DELAY", "0.3"))

GROQ_API_KEY = os.getenv("GROQ_API_KEY")
if not GROQ_API_KEY:
    raise ValueError("GROQ_API_KEY not found in environment variables")

# WebSocket configuration
WS_HOST = os.getenv("WS_HOST", "localhost")
WS_PORT = int(os.getenv("WS_PORT", "8765"))

# Agent configuration
MAX_AGENT_CACHE_SIZE = int(os.getenv("MAX_AGENT_CACHE_SIZE", "50"))

# Chat cache configuration
MAX_CHAT_CACHE_SIZE = int(os.getenv("MAX_CHAT_CACHE_SIZE", "100"))
MAX_CHAT_MESSAGES = int(os.getenv("MAX_CHAT_MESSAGES", "200"))
MAX_CHAT_TOKENS = int(os.getenv("MAX_CHAT_TOKENS", "50000"))

# LLM configuration
LLM_MODEL = os.getenv("LLM_MODEL", "llama-3.3-70b-versatile")
LLM_TEMPERATURE = float(os.getenv("LLM_TEMPERATURE", "0.7"))
LLM_MAX_TOKENS = int(os.getenv("LLM_MAX_TOKENS", "2000"))

# Context window management
MAX_CONTEXT_MESSAGES = int(os.getenv("MAX_CONTEXT_MESSAGES", "20"))
CONTEXT_STRATEGY = os.getenv("CONTEXT_STRATEGY", "sliding_window")

# ==========================
# Data Models
# ==========================
class AgentSystem:
    """Represents agent_systems table"""
    def __init__(self, agent_system_uuid: str, category_system_preset: dict):
        self.agent_system_uuid = agent_system_uuid
        self.category_system_preset = category_system_preset
        self.updated_at = time.time()


class AgentConfig:
    """Represents agents_config table"""
    def __init__(self, agent_config_uuid: str, category_id: int, 
                 category_preset_enabled: bool, agent_system_uuid: str):
        self.agent_config_uuid = agent_config_uuid
        self.category_id = category_id
        self.category_preset_enabled = category_preset_enabled
        self.agent_system_uuid = agent_system_uuid
        self.created_at = time.time()
        self.updated_at = time.time()


class Agent:
    """Represents agents table"""
    def __init__(self, agent_uuid: str, name: str, description: str, 
                 agent_config_uuid: str, auth_uuid: str,
                 agent_config: Optional[AgentConfig] = None,
                 agent_system: Optional[AgentSystem] = None):
        self.agent_uuid = agent_uuid
        self.name = name
        self.description = description
        self.agent_config_uuid = agent_config_uuid
        self.auth_uuid = auth_uuid
        self.agent_config = agent_config
        self.agent_system = agent_system
        self.creation_date = time.time()
        self.last_used = time.time()

    def touch(self):
        """Update last access time for LRU"""
        self.last_used = time.time()
    
    def get_system_prompt(self) -> str:
        """Extract system prompt from agent system configuration"""
        if self.agent_system and self.agent_system.category_system_preset:
            return self.agent_system.category_system_preset.get(
                "system_prompt", 
                "You are a helpful assistant."
            )
        return "You are a helpful assistant."


class Chat:
    """Represents chats table"""
    def __init__(self, chat_uuid: str, agent_uuid: str, auth_uuid: str):
        self.chat_uuid = chat_uuid
        self.agent_uuid = agent_uuid
        self.auth_uuid = auth_uuid
        self.created_at = time.time()
        self.updated_at = time.time()
        self.deleted_at = None
        self.last_accessed = time.time()
    
    def touch(self):
        """Update last access time for LRU"""
        self.last_accessed = time.time()
        self.updated_at = time.time()


class Message:
    """Represents messages table"""
    def __init__(self, message_uuid: str, sender_uuid: str, sender_type: str,
                 receiver_uuid: str, receiver_type: str,
                 chat_uuid: str, message_content_uuid: str, content: str,
                 created_at: float = None):
        self.message_uuid = message_uuid
        self.sender_uuid = sender_uuid
        self.sender_type = sender_type  # 'AUTH' or 'AGENT'
        self.receiver_uuid = receiver_uuid
        self.receiver_type = receiver_type
        self.chat_uuid = chat_uuid
        self.message_content_uuid = message_content_uuid
        self.content = content
        self.created_at = created_at or time.time()


class ChatSession:
    """In-memory representation of a chat with its messages"""
    def __init__(self, chat: Chat):
        self.chat = chat
        self.messages: List[Message] = []
        self.last_message_count = 0  # Track for incremental updates
    
    def add_message(self, msg: Message):
        """Add a message to this chat session"""
        self.messages.append(msg)
        self.chat.touch()
    
    def load_history(self, messages: List[Message]):
        """Load full message history (for cache miss or re-init)"""
        self.messages = sorted(messages, key=lambda m: m.created_at)
        self.last_message_count = len(self.messages)
        self.chat.touch()
        print(f"[Chat Cache] Loaded {len(messages)} messages for chat {self.chat.chat_uuid[:8]}...")
    
    def append_incremental(self, new_messages: List[Message]):
        """Append only new messages (for incremental updates)"""
        if new_messages:
            self.messages.extend(new_messages)
            self.last_message_count = len(self.messages)
            self.chat.touch()
            print(f"[Chat Cache] Added {len(new_messages)} incremental messages to {self.chat.chat_uuid[:8]}...")
    
    def needs_full_reload(self, incoming_message_count: int) -> bool:
        """Determine if we need full history from API"""
        return (incoming_message_count < self.last_message_count or 
                incoming_message_count > self.last_message_count + 10)
    
    def get_recent_messages(self, limit: int) -> List[Message]:
        """Get the most recent N messages"""
        return self.messages[-limit:] if len(self.messages) > limit else self.messages
    
    def get_stats(self) -> dict:
        """Get statistics about this chat session"""
        total_chars = sum(len(m.content) for m in self.messages)
        estimated_tokens = total_chars // 4
        return {
            "chat_uuid": self.chat.chat_uuid,
            "agent_uuid": self.chat.agent_uuid,
            "auth_uuid": self.chat.auth_uuid,
            "total_messages": len(self.messages),
            "total_characters": total_chars,
            "estimated_tokens": estimated_tokens,
            "created_at": self.chat.created_at,
            "last_accessed": self.chat.last_accessed,
            "age_seconds": time.time() - self.chat.created_at
        }
    
    def is_oversized(self, max_messages: int, max_tokens: int) -> bool:
        """Check if chat session exceeds size limits"""
        if len(self.messages) > max_messages:
            return True
        stats = self.get_stats()
        if stats['estimated_tokens'] > max_tokens:
            return True
        return False


# ==========================
# Agent Cache (LRU)
# ==========================
class AgentCache:
    """LRU cache for agents"""
    def __init__(self, max_size: int = 50):
        self.max_size = max_size
        self.cache: Dict[str, Agent] = OrderedDict()
        self.evictions = 0

    def get(self, agent_uuid: str) -> Optional[Agent]:
        """Get agent from cache, update LRU order"""
        agent = self.cache.get(agent_uuid)
        if agent:
            agent.touch()
            self.cache.move_to_end(agent_uuid)
        return agent

    def put(self, agent: Agent):
        """Add or update agent in cache"""
        if agent.agent_uuid in self.cache:
            self.cache.move_to_end(agent.agent_uuid)
        else:
            if len(self.cache) >= self.max_size:
                evicted_uuid, evicted_agent = self.cache.popitem(last=False)
                self.evictions += 1
                print(f"[Agent Cache] Evicted agent '{evicted_agent.name}' ({evicted_uuid[:8]}...)")
            self.cache[agent.agent_uuid] = agent
    
    def size(self) -> int:
        return len(self.cache)
    
    def get_stats(self) -> dict:
        """Get cache statistics"""
        return {
            "agents_in_cache": self.size(),
            "max_cache_size": self.max_size,
            "cache_utilization": f"{(self.size() / self.max_size) * 100:.1f}%",
            "total_evictions": self.evictions
        }


# ==========================
# Chat Cache (Smart LRU with incremental loading)
# ==========================
class ChatCache:
    """LRU cache for chat sessions with incremental update support"""
    def __init__(self, max_cache_size: int = 100, max_context_messages: int = 20,
                 max_chat_messages: int = 200, max_chat_tokens: int = 50000):
        self.cache: Dict[str, ChatSession] = OrderedDict()
        self.max_cache_size = max_cache_size
        self.max_context_messages = max_context_messages
        self.max_chat_messages = max_chat_messages
        self.max_chat_tokens = max_chat_tokens
        
        self.evictions = {
            "lru_evictions": 0,
            "size_evictions": 0,
            "total_evictions": 0
        }
        
        self.stats = {
            "cache_hits": 0,
            "cache_misses": 0,
            "full_reloads": 0,
            "incremental_updates": 0
        }

    def _evict_session(self, chat_uuid: str, reason: str):
        """Evict a chat session from cache"""
        if chat_uuid in self.cache:
            session = self.cache[chat_uuid]
            stats = session.get_stats()
            print(f"[Chat Cache] Evicted chat {chat_uuid[:8]}... ({reason}) - "
                  f"{stats['total_messages']} messages, ~{stats['estimated_tokens']} tokens")
            del self.cache[chat_uuid]
            self.evictions['total_evictions'] += 1

    def get_or_create_session(self, chat_uuid: str, agent_uuid: str, auth_uuid: str) -> ChatSession:
        """Get existing session or create new empty one"""
        if chat_uuid in self.cache:
            session = self.cache[chat_uuid]
            session.chat.touch()
            self.cache.move_to_end(chat_uuid)
            self.stats['cache_hits'] += 1
            return session
        
        # Cache miss - create new empty session
        self.stats['cache_misses'] += 1
        chat = Chat(chat_uuid, agent_uuid, auth_uuid)
        session = ChatSession(chat)
        
        # LRU eviction if needed
        if len(self.cache) >= self.max_cache_size:
            evicted_uuid = next(iter(self.cache))
            self._evict_session(evicted_uuid, "LRU cache full")
            self.evictions['lru_evictions'] += 1
        
        self.cache[chat_uuid] = session
        return session

    def sync_messages(self, chat_uuid: str, agent_uuid: str, auth_uuid: str,
                     messages: List[dict], mode: str = "auto") -> ChatSession:
        """
        Sync messages from API to cache
        
        Modes:
        - "auto": Automatically decide (incremental or full)
        - "incremental": Only add new messages
        - "full": Replace entire history
        """
        session = self.get_or_create_session(chat_uuid, agent_uuid, auth_uuid)
        
        # Convert dict messages to Message objects
        msg_objects = [
            Message(
                m.get("message_uuid", str(uuid.uuid4())),
                m["sender_uuid"],
                m["sender_type"],
                m["receiver_uuid"],
                m["receiver_type"],
                chat_uuid,
                m.get("message_content_uuid", str(uuid.uuid4())),
               m.get("content")
               or m.get("MessageContent", {}).get("Content")
               or "",
               m.get("created_at", time.time())
            )
            for m in messages
        ]

        if mode == "auto":
            # Decide based on cache state
            if len(session.messages) == 0:
                session.load_history(msg_objects)
                self.stats['full_reloads'] += 1
            elif session.needs_full_reload(len(msg_objects)):
                print(f"[Chat Cache] Desync detected for {chat_uuid[:8]}... "
                      f"(cache: {len(session.messages)}, incoming: {len(msg_objects)})")
                session.load_history(msg_objects)
                self.stats['full_reloads'] += 1
            else:
                new_messages = msg_objects[len(session.messages):]
                session.append_incremental(new_messages)
                self.stats['incremental_updates'] += 1
        elif mode == "full":
            session.load_history(msg_objects)
            self.stats['full_reloads'] += 1
        elif mode == "incremental":
            new_messages = msg_objects[len(session.messages):]
            session.append_incremental(new_messages)
            self.stats['incremental_updates'] += 1
        
        # Check for oversized session
        if session.is_oversized(self.max_chat_messages, self.max_chat_tokens):
            print(f"[Chat Cache] Warning: Oversized chat {chat_uuid[:8]}...")
            self._evict_session(chat_uuid, "exceeded size limits")
            self.evictions['size_evictions'] += 1
            session = self.get_or_create_session(chat_uuid, agent_uuid, auth_uuid)
            recent_messages = msg_objects[-self.max_context_messages:]
            session.load_history(recent_messages)
        
        return session

    def add_new_message(self, msg: Message, agent_uuid: str, auth_uuid: str):
        """Add a single new message (for agent responses)"""
        session = self.get_or_create_session(msg.chat_uuid, agent_uuid, auth_uuid)
        session.add_message(msg)

    def get_langchain_messages(self, chat_uuid: str, agent_uuid: str, 
                              auth_uuid: str, system_prompt: str, 
                              use_sliding_window: bool = True) -> List:
        """Convert chat history to LangChain message format"""
        messages = [SystemMessage(content=system_prompt)]
        
        session = self.get_or_create_session(chat_uuid, agent_uuid, auth_uuid)
        
        recent_msgs = (session.get_recent_messages(self.max_context_messages) 
                      if use_sliding_window else session.messages)
        
        for m in recent_msgs:
            if m.sender_type == "AGENT":
                messages.append(AIMessage(content=m.content))
            else:
                messages.append(HumanMessage(content=m.content))
        
        return messages
    
    def get_session_stats(self, chat_uuid: str) -> Optional[dict]:
        """Get statistics about a specific chat session"""
        if chat_uuid not in self.cache:
            return None
        return self.cache[chat_uuid].get_stats()
    
    def get_cache_stats(self) -> dict:
        """Get overall cache statistics"""
        total_sessions = len(self.cache)
        total_messages = sum(len(s.messages) for s in self.cache.values())
        total_chars = sum(sum(len(m.content) for m in s.messages) for s in self.cache.values())
        
        largest_session = None
        if self.cache:
            largest_session = max(self.cache.values(), 
                                key=lambda s: s.get_stats()['estimated_tokens'])
            largest_stats = largest_session.get_stats()
        
        hit_rate = (self.stats['cache_hits'] / 
                   (self.stats['cache_hits'] + self.stats['cache_misses']) * 100
                   if (self.stats['cache_hits'] + self.stats['cache_misses']) > 0 else 0)
        
        return {
            "total_chats_in_cache": total_sessions,
            "max_cache_size": self.max_cache_size,
            "cache_utilization": f"{(total_sessions / self.max_cache_size) * 100:.1f}%",
            "total_messages": total_messages,
            "total_characters": total_chars,
            "estimated_tokens": total_chars // 4,
            "avg_messages_per_chat": total_messages / total_sessions if total_sessions > 0 else 0,
            "cache_hit_rate": f"{hit_rate:.1f}%",
            "performance_stats": self.stats,
            "size_limits": {
                "max_messages_per_chat": self.max_chat_messages,
                "max_tokens_per_chat": self.max_chat_tokens,
                "context_window": self.max_context_messages
            },
            "largest_chat": {
                "chat_uuid": largest_stats['chat_uuid'][:8] + "..." if largest_session else None,
                "messages": largest_stats['total_messages'] if largest_session else 0,
                "tokens": largest_stats['estimated_tokens'] if largest_session else 0
            } if largest_session else None,
            "eviction_stats": self.evictions
        }


# ==========================
# Agent Manager
# ==========================
class AgentManager:
    """Manages agent lifecycle and caching"""
    def __init__(self, max_cache_size: int = 50):
        self.cache = AgentCache(max_cache_size)

    def create_agent(self, name: str, description: str, auth_uuid: str,
                    category_id: int = 1, system_prompt: str = None) -> Agent:
        """Create a new agent with full configuration"""
        agent_system_uuid = str(uuid.uuid4())
        category_system_preset = {
            "system_prompt": system_prompt or "You are a helpful assistant.",
            "temperature": LLM_TEMPERATURE,
            "max_tokens": LLM_MAX_TOKENS
        }
        agent_system = AgentSystem(agent_system_uuid, category_system_preset)
        
        agent_config_uuid = str(uuid.uuid4())
        agent_config = AgentConfig(
            agent_config_uuid,
            category_id,
            True,
            agent_system_uuid
        )
        
        agent_uuid = str(uuid.uuid4())
        agent = Agent(
            agent_uuid,
            name,
            description,
            agent_config_uuid,
            auth_uuid,
            agent_config,
            agent_system
        )
        
        self.cache.put(agent)
        print(f"[Agent Manager] Created agent '{name}' ({agent_uuid[:8]}...) for user {auth_uuid[:8]}...")
        return agent

    def get_or_create(self, agent_uuid: str = None, auth_uuid: str = None,
                     name: str = None, description: str = None,
                     category_id: int = 1, system_prompt: str = None) -> Agent:
        """Get existing agent from cache or create new one"""
        if agent_uuid:
            agent = self.cache.get(agent_uuid)
            if agent:
                return agent
        
        return self.create_agent(
            name or "Default Agent",
            description or "Auto-created agent",
            auth_uuid,
            category_id,
            system_prompt
        )

# ==========================
# Connection Pool Manager
# ==========================
class ConnectionPool:
    """Manages WebSocket connections with pooling support"""
    def __init__(self):
        self.connections: Dict[str, websockets.WebSocketServerProtocol] = {}
        self.connection_metadata: Dict[str, dict] = {}
        self.lock = asyncio.Lock()
    
    async def register(self, connection_id: str, websocket: websockets.WebSocketServerProtocol, 
                      auth_uuid: str = None):
        """Register a new connection"""
        async with self.lock:
            self.connections[connection_id] = websocket
            self.connection_metadata[connection_id] = {
                "auth_uuid": auth_uuid,
                "connected_at": time.time(),
                "last_activity": time.time(),
                "messages_sent": 0,
                "messages_received": 0
            }
            print(f"[Connection Pool] Registered connection {connection_id[:8]}... "
                  f"(Total: {len(self.connections)})")
    
    async def unregister(self, connection_id: str):
        """Unregister a connection"""
        async with self.lock:
            if connection_id in self.connections:
                metadata = self.connection_metadata.get(connection_id, {})
                duration = time.time() - metadata.get("connected_at", time.time())
                print(f"[Connection Pool] Unregistered {connection_id[:8]}... "
                      f"(Duration: {duration:.1f}s, Sent: {metadata.get('messages_sent', 0)}, "
                      f"Received: {metadata.get('messages_received', 0)})")
                del self.connections[connection_id]
                del self.connection_metadata[connection_id]
    
    async def update_activity(self, connection_id: str, increment_sent: bool = False, 
                             increment_received: bool = False):
        """Update connection activity"""
        async with self.lock:
            if connection_id in self.connection_metadata:
                self.connection_metadata[connection_id]["last_activity"] = time.time()
                if increment_sent:
                    self.connection_metadata[connection_id]["messages_sent"] += 1
                if increment_received:
                    self.connection_metadata[connection_id]["messages_received"] += 1
    
    def get_connection(self, connection_id: str) -> Optional[websockets.WebSocketServerProtocol]:
        """Get a connection by ID"""
        return self.connections.get(connection_id)
    
    async def get_stats(self) -> dict:
        """Get connection pool statistics"""
        async with self.lock:
            total_sent = sum(m.get("messages_sent", 0) for m in self.connection_metadata.values())
            total_received = sum(m.get("messages_received", 0) for m in self.connection_metadata.values())
            avg_duration = 0
            if self.connection_metadata:
                now = time.time()
                avg_duration = sum(now - m.get("connected_at", now) 
                                 for m in self.connection_metadata.values()) / len(self.connection_metadata)
            
            return {
                "active_connections": len(self.connections),
                "total_messages_sent": total_sent,
                "total_messages_received": total_received,
                "average_connection_duration": f"{avg_duration:.1f}s"
            }
    
    async def cleanup_stale_connections(self, timeout: int = CONNECTION_TIMEOUT):
        """Remove stale connections that haven't been active"""
        async with self.lock:
            now = time.time()
            stale = []
            for conn_id, metadata in self.connection_metadata.items():
                if now - metadata["last_activity"] > timeout:
                    stale.append(conn_id)
            
            for conn_id in stale:
                print(f"[Connection Pool] Closing stale connection {conn_id[:8]}...")
                ws = self.connections.get(conn_id)
                if ws:
                    try:
                        await ws.close(code=1001, reason="Connection timeout")
                    except:
                        pass
                del self.connections[conn_id]
                del self.connection_metadata[conn_id]

# ==========================
# Streaming Callback Handler
# ==========================

class WebSocketStreamingCallback(AsyncCallbackHandler):
    """Handles streaming LLM tokens to WebSocket with smart buffering"""
    def __init__(self, websocket, chat_uuid: str, agent_uuid: str):
        self.websocket = websocket
        self.chat_uuid = chat_uuid
        self.agent_uuid = agent_uuid
        self.full_response = ""
        self.buffer = ""
        self.last_send = time.time()
        self.min_chunk_size = STREAM_MIN_CHUNK_SIZE
        self.max_delay = STREAM_MAX_DELAY
        self.word_boundary_chars = {' ', '\n', '\t', '.', ',', '!', '?', ';', ':', '-'}

    def _should_send_buffer(self) -> bool:
        """Determina se deve enviar o buffer agora"""
        now = time.time()
        
        # Sempre enviar se o delay máximo foi atingido
        if (now - self.last_send) >= self.max_delay:
            return True
        
        # Enviar se temos caracteres suficientes E estamos em um limite de palavra
        if len(self.buffer) >= self.min_chunk_size:
            # Verificar se o último caractere é um delimitador de palavra
            if self.buffer and self.buffer[-1] in self.word_boundary_chars:
                return True
            
            # Se temos muito conteúdo (2x o mínimo), enviar mesmo sem delimitador
            if len(self.buffer) >= self.min_chunk_size * 2:
                return True
        
        return False

    async def on_llm_new_token(self, token: str, **kwargs):
        """Called when a new token is generated"""
        self.full_response += token
        self.buffer += token
        
        if self._should_send_buffer():
            await self._send_buffer()

    async def _send_buffer(self):
        """Envia o buffer acumulado"""
        if self.buffer:
            try:
                await self.websocket.send(json.dumps({
                    "chat_uuid": self.chat_uuid,
                    "agent_uuid": self.agent_uuid,
                    "content": self.buffer,
                    "partial": True
                }))
                self.buffer = ""
                self.last_send = time.time()
            except ConnectionClosed:
                print(f"[Streaming] Connection closed while sending chunk")
            except Exception as e:
                print(f"[Streaming] Error sending buffer: {e}")

    async def on_llm_end(self, response, **kwargs):
        """Called when LLM finishes generating"""
        # Enviar qualquer conteúdo restante no buffer
        await self._send_buffer()
        
        # Enviar mensagem final
        try:
            await self.websocket.send(json.dumps({
                "chat_uuid": self.chat_uuid,
                "agent_uuid": self.agent_uuid,
                "content": self.full_response,
                "partial": False,
                "message_uuid": str(uuid.uuid4()),
                "message_content_uuid": str(uuid.uuid4())
            }))
        except ConnectionClosed:
            print(f"[Streaming] Connection closed while sending final message")
        except Exception as e:
            print(f"[Streaming] Error sending final message: {e}")

# ==========================
# WebSocket Server
# ==========================
agent_manager = AgentManager(max_cache_size=MAX_AGENT_CACHE_SIZE)
chat_cache = ChatCache(
    max_cache_size=MAX_CHAT_CACHE_SIZE,
    max_context_messages=MAX_CONTEXT_MESSAGES,
    max_chat_messages=MAX_CHAT_MESSAGES,
    max_chat_tokens=MAX_CHAT_TOKENS
)

llm = ChatGroq(
    model=LLM_MODEL,
    temperature=LLM_TEMPERATURE,
    max_tokens=LLM_MAX_TOKENS,
    streaming=True,
    api_key=GROQ_API_KEY
)

connection_pool = ConnectionPool()

# ==========================
# Heartbeat Task
# ==========================
async def heartbeat_task(websocket, connection_id: str):
    """Send periodic pings to keep connection alive"""
    try:
        while True:
            await asyncio.sleep(HEARTBEAT_INTERVAL)
            if connection_pool.get_connection(connection_id):
                try:
                    pong = await websocket.ping()
                    await asyncio.wait_for(pong, timeout=10)
                    await connection_pool.update_activity(connection_id)
                except asyncio.TimeoutError:
                    print(f"[Heartbeat] Ping timeout for {connection_id[:8]}...")
                    break
                except ConnectionClosed:
                    break
            else:
                break
    except asyncio.CancelledError:
        pass
    except Exception as e:
        print(f"[Heartbeat] Error for {connection_id[:8]}...: {e}")

# ==========================
# Connection Cleanup Task
# ==========================
async def cleanup_task():
    """Periodically cleanup stale connections"""
    while True:
        await asyncio.sleep(60)  # Check every minute
        try:
            await connection_pool.cleanup_stale_connections()
        except Exception as e:
            print(f"[Cleanup] Error: {e}")

async def handle_connection(websocket):
    """Handle WebSocket connections with pooling support"""
    connection_id = str(uuid.uuid4())
    heartbeat = None
    auth_uuid = None
    
    try:
        # Register connection
        await connection_pool.register(connection_id, websocket)
        
        # Start heartbeat
        heartbeat = asyncio.create_task(heartbeat_task(websocket, connection_id))
        
        print(f"[Connection] New connection established: {connection_id[:8]}...")
        
        async for message in websocket:
            try:
                await connection_pool.update_activity(connection_id, increment_received=True)
                data = json.loads(message)
                
                # Handle connection identification
                if data.get("command") == "identify":
                    auth_uuid = data.get("auth_uuid")
                    if auth_uuid:
                        async with connection_pool.lock:
                            if connection_id in connection_pool.connection_metadata:
                                connection_pool.connection_metadata[connection_id]["auth_uuid"] = auth_uuid
                        print(f"[Connection] {connection_id[:8]}... identified as user {auth_uuid[:8]}...")
                    
                    await websocket.send(json.dumps({
                        "type": "identified",
                        "connection_id": connection_id
                    }))
                    await connection_pool.update_activity(connection_id, increment_sent=True)
                    continue
                
                # Handle stats request
                if data.get("command") == "stats":
                    agent_stats = agent_manager.cache.get_stats()
                    chat_stats = chat_cache.get_cache_stats()
                    pool_stats = await connection_pool.get_stats()
                    await websocket.send(json.dumps({
                        "type": "stats",
                        "agent_cache": agent_stats,
                        "chat_cache": chat_stats,
                        "connection_pool": pool_stats
                    }))
                    await connection_pool.update_activity(connection_id, increment_sent=True)
                    continue
                
                # Require identification before processing messages
                if not auth_uuid:
                    await websocket.send(json.dumps({
                        "error": "Connection not identified. Send 'identify' command first.",
                        "connection_id": connection_id
                    }))
                    await connection_pool.update_activity(connection_id, increment_sent=True)
                    continue
                
                # Extract message data
                chat_uuid = data.get("chat_uuid")
                content = data.get("content")
                sender_uuid = data.get("sender_uuid")
                sender_type = data.get("sender_type", "AUTH")
                receiver_uuid = data.get("receiver_uuid")
                receiver_type = data.get("receiver_type", "AGENT")
                
                # Validation
                if not all([chat_uuid, content, sender_uuid]):
                    await websocket.send(json.dumps({
                        "error": "Missing required fields: chat_uuid, content, sender_uuid",
                        "connection_id": connection_id
                    }))
                    await connection_pool.update_activity(connection_id, increment_sent=True)
                    continue
                
                # Verify auth_uuid matches
                if sender_uuid != auth_uuid:
                    print(f"[Security] Auth mismatch: connection {auth_uuid[:8]} tried to send as {sender_uuid[:8]}")
                    await websocket.send(json.dumps({
                        "error": "Sender UUID does not match authenticated user",
                        "connection_id": connection_id
                    }))
                    await connection_pool.update_activity(connection_id, increment_sent=True)
                    continue
                
                print(f"[Connection {connection_id[:8]}] Processing message from {sender_uuid[:8]} "
                      f"for chat {chat_uuid[:8]}")
                
                # Chat history sync
                chat_history = data.get("chat_history", [])
                sync_mode = data.get("sync_mode", "auto")
                
                # Agent configuration
                agent_uuid = data.get("agent_uuid")
                agent_name = data.get("agent_name")
                agent_description = data.get("agent_description")
                category_id = data.get("category_id", 1)
                system_prompt = data.get("system_prompt")
                
                # Get or create agent
                agent = agent_manager.get_or_create(
                    agent_uuid=agent_uuid or receiver_uuid,
                    auth_uuid=sender_uuid,
                    name=agent_name,
                    description=agent_description,
                    category_id=category_id,
                    system_prompt=system_prompt
                )

                # Sync chat history if provided
                if chat_history:
                    print(f"[Chat {chat_uuid[:8]}] Syncing {len(chat_history)} messages (mode: {sync_mode})")
                    chat_cache.sync_messages(
                        chat_uuid,
                        agent.agent_uuid,
                        sender_uuid,
                        chat_history,
                        mode=sync_mode
                    )

                # Create and add user message
                message_uuid = str(uuid.uuid4())
                message_content_uuid = str(uuid.uuid4())
                user_msg = Message(
                    message_uuid,
                    sender_uuid,
                    sender_type,
                    agent.agent_uuid,
                    "AGENT",
                    chat_uuid,
                    message_content_uuid,
                    content
                )
                chat_cache.add_new_message(user_msg, agent.agent_uuid, sender_uuid)

                # Get system prompt
                agent_system_prompt = agent.get_system_prompt()

                # Build LangChain messages
                messages = chat_cache.get_langchain_messages(
                    chat_uuid,
                    agent.agent_uuid,
                    sender_uuid,
                    agent_system_prompt,
                    use_sliding_window=(CONTEXT_STRATEGY == "sliding_window")
                )
                
                # Log context
                stats = chat_cache.get_session_stats(chat_uuid)
                if stats:
                    print(f"[Chat {chat_uuid[:8]}] Context: {len(messages)-1} messages, "
                          f"~{stats['estimated_tokens']} tokens, agent: {agent.name}")

                # Stream LLM response
                callback = WebSocketStreamingCallback(websocket, chat_uuid, agent.agent_uuid)

                try:
                    print(f"[Chat {chat_uuid[:8]}] Invoking LLM...")
                    start_time = time.time()
                    
                    response = await llm.ainvoke(
                        messages,
                        config={"callbacks": [callback]}
                    )
                    
                    elapsed = time.time() - start_time
                    
                    # Save agent response
                    llm_message_uuid = str(uuid.uuid4())
                    llm_message_content_uuid = str(uuid.uuid4())
                    agent_msg = Message(
                        llm_message_uuid,
                        agent.agent_uuid,
                        "AGENT",
                        sender_uuid,
                        "AUTH",
                        chat_uuid,
                        llm_message_content_uuid,
                        callback.full_response
                    )
                    chat_cache.add_new_message(agent_msg, agent.agent_uuid, sender_uuid)
                    
                    print(f"[Chat {chat_uuid[:8]}] Agent '{agent.name}' responded "
                          f"({len(callback.full_response)} chars, {elapsed:.2f}s)")
                    
                    await connection_pool.update_activity(connection_id, increment_sent=True)

                except Exception as e:
                    print(f"[Error] LLM error in chat {chat_uuid[:8]}: {str(e)}")
                    import traceback
                    traceback.print_exc()
                    await websocket.send(json.dumps({
                        "error": str(e),
                        "chat_uuid": chat_uuid,
                        "connection_id": connection_id
                    }))
                    await connection_pool.update_activity(connection_id, increment_sent=True)
                    
            except json.JSONDecodeError as e:
                print(f"[Error] Invalid JSON: {str(e)}")
                await websocket.send(json.dumps({
                    "error": f"Invalid JSON: {str(e)}"
                }))
                await connection_pool.update_activity(connection_id, increment_sent=True)
            except Exception as e:
                print(f"[Error] Message handler error: {str(e)}")
                import traceback
                traceback.print_exc()
                try:
                    await websocket.send(json.dumps({
                        "error": f"Server error: {str(e)}"
                    }))
                    await connection_pool.update_activity(connection_id, increment_sent=True)
                except:
                    pass
                    
    except ConnectionClosed as e:
        print(f"[Info] Connection {connection_id[:8]}... closed: {e}")
    except Exception as e:
        print(f"[Error] Connection handler error: {str(e)}")
        import traceback
        traceback.print_exc()
    finally:
        # Cleanup
        if heartbeat:
            heartbeat.cancel()
            try:
                await heartbeat
            except asyncio.CancelledError:
                pass
        await connection_pool.unregister(connection_id)


async def main():
    """Start WebSocket server with all tasks"""
    # Start cleanup task
    cleanup = asyncio.create_task(cleanup_task())
    
    try:
        async with websockets.serve(
            handle_connection, 
            WS_HOST, 
            WS_PORT,
            ping_interval=None,  # We handle pings manually
            ping_timeout=None
        ):
            print(f"╔═══════════════════════════════════════════════════════════════╗")
            print(f"║  WebSocket LLM Server with Connection Pooling               ║")
            print(f"╠═══════════════════════════════════════════════════════════════╣")
            print(f"║  Server:     ws://{WS_HOST}:{WS_PORT}                           ")
            print(f"║  Model:      {LLM_MODEL}                          ")
            print(f"║  Temperature: {LLM_TEMPERATURE}                                   ")
            print(f"╠═══════════════════════════════════════════════════════════════╣")
            print(f"║  Caching Configuration:                                       ║")
            print(f"║  - Max agent cache: {MAX_AGENT_CACHE_SIZE}                                 ")
            print(f"║  - Max chat cache: {MAX_CHAT_CACHE_SIZE}                                ")
            print(f"║  - Max messages per chat: {MAX_CHAT_MESSAGES}                       ")
            print(f"║  - Max tokens per chat: {MAX_CHAT_TOKENS}                     ")
            print(f"║  - Context window: {MAX_CONTEXT_MESSAGES} messages                     ")
            print(f"╠═══════════════════════════════════════════════════════════════╣")
            print(f"║  Streaming Configuration:                                     ║")
            print(f"║  - Min chunk size: {STREAM_MIN_CHUNK_SIZE} chars                         ")
            print(f"║  - Max delay: {STREAM_MAX_DELAY}s                                  ")
            print(f"╠═══════════════════════════════════════════════════════════════╣")
            print(f"║  Connection Pool:                                             ║")
            print(f"║  - Heartbeat interval: {HEARTBEAT_INTERVAL}s                           ")
            print(f"║  - Connection timeout: {CONNECTION_TIMEOUT}s                          ")
            print(f"╚═══════════════════════════════════════════════════════════════╝")
            print(f"\n✓ Server ready to accept connections...")
            
            await asyncio.Future()
    finally:
        cleanup.cancel()
        try:
            await cleanup
        except asyncio.CancelledError:
            pass

if __name__ == "__main__":
    asyncio.run(main())
