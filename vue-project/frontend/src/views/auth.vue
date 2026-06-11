<template>
  <div class="auth-card">
    <div class="card p-4 bg-cream">
      <h3 class="text-center mb-4">Авторизация</h3>
      <form @submit.prevent="handleLogin">
        <div class="mb-3">
          <label class="form-label">Логин</label>
          <input type="text" class="form-control" v-model="login" required>
        </div>
        <div class="mb-3">
          <label class="form-label">Пароль</label>
          <input type="password" class="form-control" v-model="password" required>
        </div>
        <button type="submit" class="btn btn-primary-custom w-100 rounded-pill" :disabled="loading">
          {{ loading ? 'Загрузка...' : 'Войти' }}
        </button>
        <p v-if="error" class="text-danger mt-2">{{ error }}</p>
      </form>
      <p class="mt-3 text-center">Нет аккаунта? <router-link to="/register">Регистрация</router-link></p>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Auth',
  data() {
    return {
      login: '',
      password: '',
      loading: false,
      error: ''
    }
  },
  methods: {
    async handleLogin() {
      this.loading = true
      this.error = ''
      
      try {
        const response = await fetch('http://localhost:8080/api/login', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            login: this.login,
            password: this.password
          })
        })
        
        const data = await response.json()
        
        if (response.ok) {
      
          localStorage.setItem('token', data.token)
          localStorage.setItem('user', JSON.stringify(data.user))
      
          window.location.href = '/'
        } else {
          this.error = data.error || 'Ошибка входа'
        }
      } catch (err) {
        this.error = 'Ошибка соединения с сервером'
      } finally {
        this.loading = false
      }
    }
  }
}
</script>