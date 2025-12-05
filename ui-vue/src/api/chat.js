// src/api/chat.js

// Get base URL from environment (same as api.js)
const API_BASE_URL = process.env.VUE_APP_API_URL || 'http://localhost:8080';

/**
 * Initialize a new chat with SSE streaming via POST
 * The backend expects POST with JSON body and returns SSE stream
 */
export const createChat = (agentUuid, messageContent, onChunk, onComplete, onError) => {
  console.log('[DEBUG FRONTEND] createChat called with:', { agentUuid, messageContent });
  
  const controller = new AbortController();
  const timeoutId = setTimeout(() => {
    console.log('[DEBUG FRONTEND] Request timeout');
    controller.abort();
    onError({ message: 'Request timeout - AI service may not be ready yet' });
  }, 30000); // 30 second timeout (increased for AI processing)

  fetch(`${API_BASE_URL}/chat/create`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify({
      agent_uuid: agentUuid,
      message_content: messageContent
    }),
    signal: controller.signal
  })
  .then(response => {
    clearTimeout(timeoutId);
    console.log('[DEBUG FRONTEND] Response received, status:', response.status);
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    
    const reader = response.body.getReader()
    const decoder = new TextDecoder()
    let buffer = ''
    let currentEvent = ''
    
    function readStream() {
      reader.read().then(({ done, value }) => {
        if (done) {
          console.log('[DEBUG FRONTEND] Stream ended');
          return
        }
        
        buffer += decoder.decode(value, { stream: true })
        const lines = buffer.split('\n')
        buffer = lines.pop() || ''
        
        for (const line of lines) {
          console.log('[DEBUG FRONTEND] Processing line:', line);
          
          if (line.startsWith('event:')) {
            currentEvent = line.substring(6).trim()
            console.log('[DEBUG FRONTEND] Event type:', currentEvent);
            continue
          }
          
          if (line.startsWith('data:')) {
            const data = line.substring(5).trim()
            console.log('[DEBUG FRONTEND] Data for event', currentEvent, ':', data.substring(0, 50));
            
            if (currentEvent === 'message') {
              onChunk(data)
            } else if (currentEvent === 'done') {
              console.log('[DEBUG FRONTEND] Done event received, data:', data);
              try {
                const finalData = JSON.parse(data)
                onComplete(finalData)
              } catch (err) {
                console.error('[DEBUG FRONTEND] Error parsing done event:', err);
                onComplete({ chat_uuid: data })
              }
              return
            } else if (currentEvent === 'error') {
              console.log('[DEBUG FRONTEND] Error event:', data);
              onError({ message: data })
              return
            }
          }
        }
        
        readStream()
      }).catch(err => {
        console.error('[DEBUG FRONTEND] Stream reading error:', err);
        onError(err)
      })
    }
    
    readStream()
  })
  .catch(err => {
    clearTimeout(timeoutId);
    console.error('[DEBUG FRONTEND] Fetch error:', err);
    if (err.name === 'AbortError') {
      onError({ message: 'Request timeout - AI service may not be ready yet' });
    } else {
      onError(err);
    }
  })
  
  return () => {
    clearTimeout(timeoutId);
    controller.abort();
    console.log('[DEBUG FRONTEND] Stream cleanup requested')
  }
}

/**
 * Send a message to an existing chat with SSE streaming via POST
 */
export const sendMessage = (chatUuid, messageContent, onChunk, onComplete, onError) => {
  console.log('[DEBUG FRONTEND] sendMessage called with:', { chatUuid, messageContent });
  
  const controller = new AbortController();
  const timeoutId = setTimeout(() => {
    console.log('[DEBUG FRONTEND] Request timeout');
    controller.abort();
    onError({ message: 'Request timeout - AI service may not be ready yet' });
  }, 30000); // 30 second timeout

  fetch(`${API_BASE_URL}/chat/send-new-message`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify({
      chat_uuid: chatUuid,
      message_content: messageContent
    }),
    signal: controller.signal
  })
  .then(response => {
    clearTimeout(timeoutId);
    console.log('[DEBUG FRONTEND] Response received, status:', response.status);
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    
    const reader = response.body.getReader()
    const decoder = new TextDecoder()
    let buffer = ''
    let currentEvent = ''
    
    function readStream() {
      reader.read().then(({ done, value }) => {
        if (done) {
          console.log('[DEBUG FRONTEND] Stream ended');
          return
        }
        
        buffer += decoder.decode(value, { stream: true })
        const lines = buffer.split('\n')
        buffer = lines.pop() || ''
        
        for (const line of lines) {
          console.log('[DEBUG FRONTEND] Processing line:', line);
          
          if (line.startsWith('event:')) {
            currentEvent = line.substring(6).trim()
            console.log('[DEBUG FRONTEND] Event type:', currentEvent);
            continue
          }
          
          if (line.startsWith('data:')) {
            const data = line.substring(5).trim()
            console.log('[DEBUG FRONTEND] Data for event', currentEvent, ':', data.substring(0, 50));
            
            if (currentEvent === 'message') {
              onChunk(data)
            } else if (currentEvent === 'done') {
              console.log('[DEBUG FRONTEND] Done event received');
              try {
                const finalData = JSON.parse(data)
                onComplete(finalData)
              } catch (err) {
                console.error('[DEBUG FRONTEND] Error parsing done event:', err);
                onComplete({})
              }
              return
            } else if (currentEvent === 'error') {
              console.log('[DEBUG FRONTEND] Error event:', data);
              onError({ message: data })
              return
            }
          }
        }
        
        readStream()
      }).catch(err => {
        console.error('[DEBUG FRONTEND] Stream reading error:', err);
        onError(err)
      })
    }
    
    readStream()
  })
  .catch(err => {
    clearTimeout(timeoutId);
    console.error('[DEBUG FRONTEND] Fetch error:', err);
    if (err.name === 'AbortError') {
      onError({ message: 'Request timeout - AI service may not be ready yet' });
    } else {
      onError(err);
    }
  })
  
  return () => {
    clearTimeout(timeoutId);
    controller.abort();
    console.log('[DEBUG FRONTEND] Stream cleanup requested')
  }
}

export default {
  createChat,
  sendMessage
}
