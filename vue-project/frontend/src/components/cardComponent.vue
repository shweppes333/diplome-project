<template>
  <div class="card h-100 p-3 bg-beige-card">
    <div class="product-img-placeholder">
      <i :class="product.icon"></i>
    </div>
    <div class="card-body text-center">
      <h5 class="card-title">{{ product.name }}</h5>
      <p class="card-text text-secondary small">{{ product.description }}</p>
      <p class="fw-bold fs-5">{{ product.price }} ₽</p>
      <button class="btn btn-primary-custom rounded-pill w-100" @click="addToCart" :disabled="!isAuthenticated || adding">
        <span v-if="!adding"><i class="bi bi-cart-plus"></i> В корзину</span>
        <span v-else><i class="bi bi-hourglass-split"></i> Добавление...</span>
      </button>
      <small v-if="!isAuthenticated" class="text-danger d-block mt-1">Авторизуйтесь для покупок</small>
    </div>
  </div>
</template>

<script>
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../store/auth.js'

export default {
  name: 'ProductCard',
  props: {
    product: {
      type: Object,
      required: true
    }
  },
  setup(props) {
    const router = useRouter()
    const authStore = useAuthStore()
    const adding = ref(false)
    const isAuthenticated = computed(() => authStore.isAuthenticated)
    
    const addToCart = async () => {
      if (!isAuthenticated.value) {
        alert("Пожалуйста, авторизуйтесь, чтобы добавить товар в корзину.")
        router.push('/login')
        return
      }
      
      adding.value = true
      
      try {
        const response = await fetch('http://localhost:8080/api/cart', {
          method: 'POST',
          headers: authStore.getAuthHeader(),
          body: JSON.stringify({
            product_id: props.product.id,
            quantity: 1
          })
        })
        
        if (response.ok) {
          alert(`Товар "${props.product.name}" добавлен в корзину!`)
        } else if (response.status === 401) {
          alert('Сессия истекла, войдите снова')
          authStore.logout()
        } else {
          alert('Ошибка при добавлении товара')
        }
      } catch (error) {
        console.error('Ошибка добавления в корзину:', error)
        alert('Ошибка соединения с сервером')
      } finally {
        adding.value = false
      }
    }
    
    return {
      isAuthenticated,
      adding,
      addToCart
    }
  }
}
</script>