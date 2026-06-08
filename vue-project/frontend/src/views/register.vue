<template>
  <div class="auth-card">
    <div class="card p-4 bg-cream">
      <h3 class="text-center mb-4">Регистрация</h3>
      <form @submit.prevent="handleRegister">
        <div class="mb-3">
          <label class="form-label">Имя</label>
          <input type="text" class="form-control" v-model="form.name" required>
        </div>
        <div class="mb-3">
          <label class="form-label">Email</label>
          <input type="email" class="form-control" v-model="form.email" required>
          <small class="text-muted">Должен содержать @</small>
        </div>
        <div class="mb-3">
          <label class="form-label">Логин</label>
          <input type="text" class="form-control" v-model="form.login" required>
          <small class="text-muted">Не менее 6 символов</small>
        </div>
        <div class="mb-3">
          <label class="form-label">Пароль</label>
          <input type="password" class="form-control" v-model="form.password" required>
          <small class="text-muted">Не менее 6 символов</small>
        </div>
        <button type="submit" class="btn btn-primary-custom w-100 rounded-pill" :disabled="loading">
          {{ loading ? 'Регистрация...' : 'Зарегистрироваться' }}
        </button>
        <div v-if="error" class="alert alert-danger mt-3 small">{{ error }}</div>
      </form>
      <p class="mt-3 text-center">Уже есть аккаунт? <router-link to="/login">Войти</router-link></p>
    </div>
  </div>
</template>

<script>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../store/auth.js'

export default {
  name: 'RegisterPage',
  setup() {
    const router = useRouter()
    const authStore = useAuthStore()
    const loading = ref(false)
    const error = ref('')
    
    const form = reactive({
      name: '',
      email: '',
      login: '',
      password: ''
    })
    
    const handleRegister = async () => {
      loading.value = true
      error.value = ''
      
      const result = await authStore.register(form)
      
      if (result.success) {
        alert('Регистрация успешна! Теперь войдите.')
        router.push('/login')
      } else {
        error.value = result.error || 'Ошибка регистрации'
      }
      
      loading.value = false
    }
    
    return {
      form,
      loading,
      error,
      handleRegister
    }
  }
}
</script>