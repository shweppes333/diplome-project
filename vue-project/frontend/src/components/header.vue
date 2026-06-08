<template>
  <nav class="navbar navbar-expand-lg navbar-light py-3 px-4 shadow-sm" style="background-color: #f1e2c9;">
    <div class="container">
      <router-link to="/" class="navbar-brand fs-3 fw-semibold" style="color:#3a2c1f;">
        <i class="bi bi-scissors"></i> Тепло рук
      </router-link>
      <div class="ms-auto d-flex gap-3 align-items-center">
        <router-link v-if="isAuthenticated" to="/cart" class="btn btn-outline-custom rounded-pill px-3 position-relative">
          <i class="bi bi-bag fs-5"></i> Корзина
          <span v-if="cartItemsCount > 0" class="position-absolute top-0 start-100 translate-middle badge rounded-pill bg-danger" style="font-size: 0.7rem;">
            {{ cartItemsCount }}
          </span>
        </router-link>
        
        <div v-if="!isAuthenticated">
          <router-link to="/auth" class="btn btn-light me-2 rounded-pill px-3 border">Вход</router-link>
          <router-link to="/register" class="btn btn-primary-custom rounded-pill px-3">Регистрация</router-link>
        </div>
        
        <div v-else class="dropdown">
          <button class="btn btn-light dropdown-toggle rounded-pill px-3" type="button" data-bs-toggle="dropdown">
            <i class="bi bi-person-circle"></i> {{ currentUserName }}
          </button>
          <ul class="dropdown-menu dropdown-menu-end">
            <li><a class="dropdown-item" href="#" @click.prevent="logout">Выйти</a></li>
          </ul>
        </div>
      </div>
    </div>
  </nav>
</template>

<script>
import { computed } from 'vue'
import { useAuthStore } from '../store/auth.js'

export default {
  name: 'Header',
  setup() {
    const authStore = useAuthStore()
    
    const isAuthenticated = computed(() => authStore.isAuthenticated)
    const currentUserName = computed(() => authStore.currentUser?.name || authStore.currentUser?.login || '')
    
    const cartItemsCount = computed(() => {
      if (!authStore.isAuthenticated) return 0
      const cart = JSON.parse(localStorage.getItem('handmade_cart') || '[]')
      return cart.reduce((cnt, it) => cnt + (it.quantity || 1), 0)
    })
    
    const logout = () => {
      authStore.logout()
    }
    
    return {
      isAuthenticated,
      currentUserName,
      cartItemsCount,
      logout
    }
  }
}
</script>