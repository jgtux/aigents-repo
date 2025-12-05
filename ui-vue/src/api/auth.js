import api from './api'

export const authService = {
  // Verifica se o usuário está autenticado
  async checkAuth() {
    try {
      const response = await api.post('/auth/check')
      return response.data
    } catch (error) {
      return null
    }
  },

  // Login
  async login(credentials) {
    try {
      const response = await api.post('/auth/login', credentials)
      return { success: true, data: response.data }
    } catch (error) {
      return { 
        success: false, 
        error: error.response?.data?.message || 'Erro ao fazer login' 
      }
    }
  },

  // Signup
  async signup(userData) {
    try {
      const response = await api.post('/auth/signup', userData)
      return { success: true, data: response.data }
    } catch (error) {
      return { 
        success: false, 
        error: error.response?.data?.message || 'Erro ao criar conta' 
      }
    }
  },

  // Logout
  async logout() {
    try {
      await api.post('/auth/logout')
      return { success: true }
    } catch (error) {
      return { 
        success: false, 
        error: error.response?.data?.message || 'Erro ao fazer logout' 
      }
    }
  },

  // Refresh token (chamado automaticamente pelo interceptor)
  async refresh() {
    try {
      const response = await api.post('/auth/refresh')
      return { success: true, data: response.data }
    } catch (error) {
      return { success: false }
    }
  }
}

// Exporta função auxiliar para o router
export const checkAuth = async () => {
  try {
    console.log('checkAuth: Making API call...')
    const response = await api.get('/auth/check', { timeout: 5000 })
    console.log('checkAuth: Response received:', response.status, response.data)
    
    if (response.status === 200) {
      console.log('checkAuth: Authentication successful')
      return true
    }
    
    console.log('checkAuth: Unexpected status code')
    throw new Error('Not authenticated')
  } catch (error) {
    console.log('checkAuth: Error caught:', error.message, error.response?.status)
    // Re-throw para o router guard capturar
    throw new Error('Not authenticated')
  }
}
