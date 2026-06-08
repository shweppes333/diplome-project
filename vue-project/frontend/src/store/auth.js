// frontend/src/store/auth.js
import { reactive, computed } from 'vue'

// Состояние (данные)
const state = reactive({
  isAuthenticated: false,    // Авторизован ли пользователь
  currentUser: null,         // Данные текущего пользователя
  token: null                // JWT токен
})

// Actions (функции для изменения состояния)
export const useAuthStore = () => {
  
  // Проверка авторизации при загрузке
  const checkAuth = () => {
    const token = localStorage.getItem('token')
    const user = localStorage.getItem('user')
    
    if (token && user) {
      state.isAuthenticated = true
      state.token = token
      state.currentUser = JSON.parse(user)
    }
  }
  
  // Вход в систему
  const login = async (login, password) => {
    try {
      const response = await fetch('http://localhost:8080/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ login, password })
      })
      
      const data = await response.json()
      
      if (response.ok) {
        // Обновляем состояние
        state.isAuthenticated = true
        state.token = data.token
        state.currentUser = data.user
        
        // Сохраняем в localStorage
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
  
  // Регистрация
  const register = async (userData) => {
    try {
      const response = await fetch('http://localhost:8080/api/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(userData)
      })
      
      if (response.ok) {
        return { success: true }
      } else {
        const data = await response.json()
        return { success: false, error: data.error }
      }
    } catch (error) {
      return { success: false, error: 'Ошибка соединения' }
    }
  }
  
  // Выход из системы
  const logout = () => {
    state.isAuthenticated = false
    state.token = null
    state.currentUser = null
    
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    
    window.location.href = '/'
  }
  
  // Получение заголовков для API запросов
  const getAuthHeader = () => ({
    'Authorization': `Bearer ${state.token}`,
    'Content-Type': 'application/json'
  })
  
  // Инициализация
  checkAuth()
  
  // Возвращаем реактивные данные и методы
  return {
    // Данные (computed для реактивности)
    isAuthenticated: computed(() => state.isAuthenticated),
    currentUser: computed(() => state.currentUser),
    token: computed(() => state.token),
    
    // Методы
    login,
    register,
    logout,
    getAuthHeader
  }
}