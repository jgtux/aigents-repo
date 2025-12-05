import axios from 'axios'

const api = axios.create({
  baseURL: process.env.VUE_APP_API_URL,
  withCredentials: true, // Importante: envia cookies HTTP-only
  headers: {
    'Content-Type': 'application/json'
  }
})

let isRefreshing = false
let failedQueue = []

const processQueue = (error, token = null) => {
  failedQueue.forEach(prom => {
    if (error) {
      prom.reject(error)
    } else {
      prom.resolve(token)
    }
  })
  
  failedQueue = []
}

// Interceptor de resposta para lidar com refresh token
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config
    
    // Skip interceptor for auth check and refresh endpoints
    if (originalRequest.url?.includes('/auth/check') || 
        originalRequest.url?.includes('/auth/refresh') ||
        originalRequest._skipAuthRefresh) {
      return Promise.reject(error)
    }

    // Se o erro for 401 e ainda não tentamos refresh
    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        // Se já está refreshing, coloca na fila
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject })
        }).then(() => {
          return api(originalRequest)
        }).catch(err => {
          return Promise.reject(err)
        })
      }

      originalRequest._retry = true
      isRefreshing = true

      try {
        // Tenta fazer refresh do token
        await api.post('/auth/refresh')
        
        processQueue(null)
        isRefreshing = false
        
        // Refaz a requisição original
        return api(originalRequest)
      } catch (refreshError) {
        processQueue(refreshError)
        isRefreshing = false
        
        // Don't redirect here - let the router guard handle it
        // Just reject the promise so the calling code can handle it
        return Promise.reject(refreshError)
      }
    }

    return Promise.reject(error)
  }
)

export default api
