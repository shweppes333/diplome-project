<template>
  <div class="auth-card">
    <div class="card p-4 bg-cream">
      <h3 class="text-center mb-4">Авторизация</h3>
      <form @submit.prevent="handleLogin">
        <div class="mb-3">
          <label class="form-label">Логин</label>
          <input type="text" class="form-control" v-model="form.login" required>
        </div>
        <div class="mb-3">
          <label class="form-label">Пароль</label>
          <input type="password" class="form-control" v-model="form.password" required>
        </div>
        <button type="submit" class="btn btn-primary-custom w-100 rounded-pill" :disabled="loading">
          {{ loading ? 'Загрузка...' : 'Авторизоваться' }}
        </button>
        <p v-if="error" class="text-danger mt-2 small">{{ error }}</p>
      </form>
      <p class="mt-3 text-center">Нет аккаунта? <router-link to="/register">Регистрация</router-link></p>
    </div>
  </div>
</template>

<script>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../store/auth.js'

export default {
  name: 'LoginPage',
  setup() {
    const router = useRouter()
    const authStore = useAuthStore()
    const loading = ref(false)
    const error = ref('')
    
    const form = reactive({
      login: '',
      password: ''
    })
    
    const handleLogin = async () => {
      loading.value = true
      error.value = ''
      
      const result = await authStore.login(form.login, form.password)
      
      if (result.success) {
        router.push('/')
      } else {
        error.value = result.error || 'Ошибка авторизации'
      }
      
      loading.value = false
    }
    
    return {
      form,
      loading,
      error,
      handleLogin
    }
  }
}
</script>