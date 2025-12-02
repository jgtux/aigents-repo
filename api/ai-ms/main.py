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

# ==========================
# Load Environment Variables
# ==========================
load_dotenv()

# ==========================
# Config
# ==========================
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
        # If incoming count is less than cache, chat was modified elsewhere
        # If incoming count is much greater, we're missing messages
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
                m["content"],
                m.get("created_at", time.time())
            )
            for m in messages
        ]
        
        if mode == "auto":
            # Decide based on cache state
            if len(session.messages) == 0:
                # Empty cache - do full load
                session.load_history(msg_objects)
                self.stats['full_reloads'] += 1
            elif session.needs_full_reload(len(msg_objects)):
                # Cache desync detected - full reload
                print(f"[Chat Cache] Desync detected for {chat_uuid[:8]}... "
                      f"(cache: {len(session.messages)}, incoming: {len(msg_objects)})")
                session.load_history(msg_objects)
                self.stats['full_reloads'] += 1
            else:
                # Incremental update - only add new messages
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
            # Recreate with only recent messages
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
# Streaming Callback Handler
# ==========================
class WebSocketStreamingCallback(AsyncCallbackHandler):
    """Handles streaming LLM tokens to WebSocket client"""
    def __init__(self, websocket, chat_uuid: str, agent_uuid: str):
        self.websocket = websocket
        self.chat_uuid = chat_uuid
        self.agent_uuid = agent_uuid
        self.full_response = ""

    async def on_llm_new_token(self, token: str, **kwargs):
        """Called when a new token is generated"""
        self.full_response += token
        await self.websocket.send(json.dumps({
            "chat_uuid": self.chat_uuid,
            "agent_uuid": self.agent_uuid,
            "content": token,
            "partial": True
        }))

    async def on_llm_end(self, response, **kwargs):
        """Called when LLM finishes generating"""
        await self.websocket.send(json.dumps({
            "chat_uuid": self.chat_uuid,
            "agent_uuid": self.agent_uuid,
            "content": self.full_response,
            "partial": False,
            "message_uuid": str(uuid.uuid4()),
            "message_content_uuid": str(uuid.uuid4())
        }))


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

async def handle_connection(websocket):
    """Handle WebSocket connections and message routing"""
    async for message in websocket:
        try:
            data = json.loads(message)
            
            # Handle stats request
            if data.get("command") == "stats":
                agent_stats = agent_manager.cache.get_stats()
                chat_stats = chat_cache.get_cache_stats()
                await websocket.send(json.dumps({
                    "type": "stats",
                    "agent_cache": agent_stats,
                    "chat_cache": chat_stats
                }))
                continue
            
            # Extract message data
            chat_uuid = data.get("chat_uuid")
            content = data.get("content")
            sender_uuid = data.get("sender_uuid")
            sender_type = data.get("sender_type", "AUTH")
            receiver_uuid = data.get("receiver_uuid")
            receiver_type = data.get("receiver_type", "AGENT")
            
            # Chat history sync (from API)
            chat_history = data.get("chat_history", [])  # List of message dicts
            sync_mode = data.get("sync_mode", "auto")  # auto, incremental, or full
            
            # Agent configuration
            agent_uuid = data.get("agent_uuid")
            agent_name = data.get("agent_name")
            agent_description = data.get("agent_description")
            category_id = data.get("category_id", 1)
            system_prompt = data.get("system_prompt")
            
            # 1️⃣ Get or create agent
            agent = agent_manager.get_or_create(
                agent_uuid=agent_uuid or receiver_uuid,
                auth_uuid=sender_uuid,
                name=agent_name,
                description=agent_description,
                category_id=category_id,
                system_prompt=system_prompt
            )

            # 2️⃣ Sync chat history if provided (smart caching)
            if chat_history:
                chat_cache.sync_messages(
                    chat_uuid,
                    agent.agent_uuid,
                    sender_uuid,
                    chat_history,
                    mode=sync_mode
                )

            # 3️⃣ Create and add user message
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

            # 4️⃣ Get system prompt
            agent_system_prompt = agent.get_system_prompt()

            # 5️⃣ Build LangChain messages
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
                print(f"[Chat {chat_uuid[:8]}...] Context: {len(messages)-1} messages, "
                      f"~{stats['estimated_tokens']} tokens")

            # 6️⃣ Stream LLM response
            callback = WebSocketStreamingCallback(websocket, chat_uuid, agent.agent_uuid)

            try:
                response = await llm.ainvoke(
                    messages,
                    config={"callbacks": [callback]}
                )
                
                # 7️⃣ Save agent response
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
                
                print(f"[Chat {chat_uuid[:8]}...] Agent '{agent.name}' responded "
                      f"({len(callback.full_response)} chars)")

            except Exception as e:
                print(f"[Error] LLM error: {str(e)}")
                await websocket.send(json.dumps({
                    "error": str(e),
                    "chat_uuid": chat_uuid
                }))
                
        except json.JSONDecodeError as e:
            print(f"[Error] Invalid JSON: {str(e)}")
            await websocket.send(json.dumps({
                "error": f"Invalid JSON: {str(e)}"
            }))
        except Exception as e:
            print(f"[Error] Server error: {str(e)}")
            await websocket.send(json.dumps({
                "error": f"Server error: {str(e)}"
            }))


async def main():
    """Start WebSocket server"""
    async with websockets.serve(handle_connection, WS_HOST, WS_PORT):
        print(f"WebSocket LLM server running on ws://{WS_HOST}:{WS_PORT}")
        print(f"Using Groq model: {LLM_MODEL}")
        print(f"️Temperature: {LLM_TEMPERATURE}")
        print(f"Max agent cache: {MAX_AGENT_CACHE_SIZE}")
        print(f"Max chat cache: {MAX_CHAT_CACHE_SIZE}")
        print(f"Max messages per chat: {MAX_CHAT_MESSAGES}")
        print(f"Max tokens per chat: {MAX_CHAT_TOKENS}")
        print(f"Context window: {MAX_CONTEXT_MESSAGES} messages")
        print(f"Ready to handle agent requests with smart caching...")
        await asyncio.Future()

if __name__ == "__main__":
    asyncio.run(main())
