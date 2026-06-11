<template>
  <div class="auth-card">
    <div class="card p-4 bg-cream">
      <h3 class="text-center mb-4">Регистрация</h3>
      <form @submit.prevent="handleRegister">
        <div class="mb-3">
          <label class="form-label">Имя</label>
          <input type="text" class="form-control" v-model="name" required>
        </div>
        <div class="mb-3">
          <label class="form-label">Email</label>
          <input type="email" class="form-control" v-model="email" required>
        </div>
        <div class="mb-3">
          <label class="form-label">Логин (мин. 6 символов)</label>
          <input type="text" class="form-control" v-model="login" required>
        </div>
        <div class="mb-3">
          <label class="form-label">Пароль (мин. 6 символов)</label>
          <input type="password" class="form-control" v-model="password" required>
        </div>
        <button type="submit" class="btn btn-primary-custom w-100 rounded-pill" :disabled="loading">
          {{ loading ? 'Регистрация...' : 'Зарегистрироваться' }}
        </button>
        <p v-if="error" class="text-danger mt-2">{{ error }}</p>
      </form>
      <p class="mt-3 text-center">Уже есть аккаунт? <router-link to="/auth">Войти</router-link></p>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Register',
  data() {
    return {
      name: '',
      email: '',
      login: '',
      password: '',
      loading: false,
      error: ''
    }
  },
  methods: {
    async handleRegister() {
      this.loading = true
      this.error = ''
      
      try {
        const response = await fetch('http://localhost:8080/api/register', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            name: this.name,
            email: this.email,
            login: this.login,
            password: this.password
          })
        })
        
        const data = await response.json()
        
        if (response.ok) {
          alert('Регистрация успешна! Теперь войдите в аккаунт.')
          this.$router.push('/auth')
        } else {
          this.error = data.error || 'Ошибка регистрации'
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