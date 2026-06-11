// frontend/src/store/auth.js
import { reactive, computed } from 'vue'

const API_URL = 'http://localhost:8080/api'

const state = reactive({
  isAuthenticated: false,
  currentUser: null,
  token: null
})

export const useAuthStore = () => {
  const checkAuth = () => {
    const token = localStorage.getItem('token')
    const user = localStorage.getItem('user')
    
    if (token && user) {
      state.isAuthenticated = true
      state.token = token
      state.currentUser = JSON.parse(user)
    } else {
      state.isAuthenticated = false
      state.token = null
      state.currentUser = null
    }
  }
  
  const login = async (login, password) => {
    try {
      const response = await fetch(`${API_URL}/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ login, password })
      })
      
      const data = await response.json()
      
      if (response.ok) {
        state.isAuthenticated = true
        state.token = data.token
        state.currentUser = data.user
        
        localStorage.setItem('token', data.token)
        localStorage.setItem('user', JSON.stringify(data.user))
        
        return { success: true }
      } else {
        return { success: false, error: data.error }
      }
    } catch (error) {
      return { success: false, error: 'Ошибка соединения' }
    }
  }
  
  const logout = () => {
    state.isAuthenticated = false
    state.token = null
    state.currentUser = null
    
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    localStorage.removeItem('handmade_cart')
  }
  
  const getAuthHeader = () => ({
    'Authorization': `Bearer ${state.token}`,
    'Content-Type': 'application/json'
  })
  
  checkAuth()
  
  return {
    isAuthenticated: computed(() => state.isAuthenticated),
    currentUser: computed(() => state.currentUser),
    login,
    logout,
    getAuthHeader,
    checkAuth
  }
}