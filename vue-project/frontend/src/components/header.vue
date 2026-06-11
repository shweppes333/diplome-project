<template>
  <nav class="navbar navbar-expand-lg navbar-light py-3 px-4 shadow-sm" style="background-color: #f1e2c9;">
    <div class="container">
      <router-link to="/" class="navbar-brand fs-3 fw-semibold" style="color:#3a2c1f;">
        <i class="bi bi-scissors"></i> HandMadeStudio
      </router-link>

      <div class="ms-auto d-flex gap-3 align-items-center">
    
        <router-link v-if="isAuthenticated" to="/cart" class="btn btn-outline-custom rounded-pill px-3">
          <i class="bi bi-bag fs-5"></i> Корзина
        </router-link>

      
        <template v-if="!isAuthenticated">
          <router-link to="/auth" class="btn btn-light rounded-pill px-4 border">
            <i class="bi bi-box-arrow-in-right"></i> Войти
          </router-link>
          <router-link to="/register" class="btn btn-primary-custom rounded-pill px-4">
            <i class="bi bi-person-plus"></i> Регистрация
          </router-link>
        </template>

     
        <template v-if="isAuthenticated">
          <span class="text-dark">
            <i class="bi bi-person-circle"></i> {{ userName }}
          </span>
          <button @click="logout" class="btn btn-danger rounded-pill px-4">
            <i class="bi bi-box-arrow-right"></i> Выйти
          </button>
        </template>
      </div>
    </div>
  </nav>
</template>

<script>
import { useAuthStore } from '../store/auth.js'

export default {
  name: 'Header',
  data() {
    return {
      isAuthenticated: false,
      userName: ''
    }
  },
  mounted() {
    this.checkAuth()
   
    window.addEventListener('storage', this.checkAuth)
  },
  beforeUnmount() {
    window.removeEventListener('storage', this.checkAuth)
  },
  methods: {
    checkAuth() {
      const token = localStorage.getItem('token')
      const user = localStorage.getItem('user')
      
      this.isAuthenticated = !!(token && user)
      
      if (user) {
        try {
          const userData = JSON.parse(user)
          this.userName = userData.name || userData.login || ''
        } catch(e) {
          this.userName = ''
        }
      }
    },
    logout() {
      if (confirm('Выйти из аккаунта?')) {
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        localStorage.removeItem('handmade_cart')
        this.isAuthenticated = false
        this.userName = ''
        this.$router.push('/')
        window.location.reload() 
      }
    }
  }
}
</script>

<style scoped>
.btn-outline-custom {
  border: 1px solid #b47c4a;
  color: #b47c4a;
  background: transparent;
}

.btn-outline-custom:hover {
  background-color: #b47c4a;
  color: white;
}

.btn-primary-custom {
  background-color: #b47c4a;
  border-color: #a06c3c;
  color: white;
}

.btn-primary-custom:hover {
  background-color: #9b623a;
  color: white;
}

.btn-light {
  background-color: #fffcf7;
  border: 1px solid #e0cfb4;
}

.btn-danger {
  background-color: #dc3545;
  border: none;
}

.btn-danger:hover {
  background-color: #c82333;
}
</style>