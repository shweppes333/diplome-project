<template>
  <div>
    <h2 class="mb-4">🛒 Моя корзина</h2>
    <div class="row g-4">
      <!-- Левая часть: товары -->
      <div class="col-md-7">
        <div class="bg-cream p-3 rounded-4" style="min-height: 300px;">
          <div v-if="loading" class="text-center py-5">
            <div class="spinner-border text-primary" role="status">
              <span class="visually-hidden">Загрузка...</span>
            </div>
          </div>
          <div v-else-if="cartItems.length === 0" class="text-center py-5">
            <i class="bi bi-bag-x fs-1 text-muted"></i>
            <p class="mt-2">Ваша корзина пуста. Добавьте товары на главной странице.</p>
            <router-link to="/" class="btn btn-outline-custom rounded-pill">На главную</router-link>
          </div>
          
          <div v-else>
            <div v-for="item in cartItems" :key="item.id" class="d-flex justify-content-between align-items-center border-bottom pb-3 mb-3">
              <div class="d-flex gap-3 align-items-center">
                <div class="product-img-placeholder" style="width: 70px; height: 70px; font-size: 1.5rem;">
                  <i :class="item.product_icon || 'bi bi-box'"></i>
                </div>
                <div>
                  <h5 class="mb-0">{{ item.product_name }}</h5>
                  <small>{{ item.price }} ₽ x {{ item.quantity }}</small>
                </div>
              </div>
              <div>
                <strong>{{ item.price * item.quantity }} ₽</strong>
                <button class="btn btn-sm btn-outline-danger ms-2 rounded-circle" @click="removeItem(item.product_id)">
                  <i class="bi bi-trash"></i>
                </button>
              </div>
            </div>
            <div class="text-end mt-3 fw-bold fs-4">Итого: {{ totalAmount }} ₽</div>
          </div>
        </div>
      </div>
      
      <!-- Правая часть: форма оплаты -->
      <div class="col-md-5">
        <div class="cart-summary">
          <h4><i class="bi bi-credit-card"></i> Данные оплаты</h4>
          <form @submit.prevent="payOrder">
            <div class="mb-3">
              <label class="form-label">Email для чека</label>
              <input type="email" class="form-control" v-model="email" required :placeholder="userEmail || 'customer@example.com'">
              <small class="text-muted">На этот email придет чек об оплате</small>
            </div>
            <div class="mb-3">
              <label class="form-label">Способ оплаты</label>
              <select class="form-select" v-model="paymentMethod" required>
                <option value="bank_card">Банковская карта</option>
                <option value="yoomoney">ЮMoney</option>
                <option value="sberpay">СберПэй</option>
              </select>
            </div>
            <div class="mb-3 p-3 bg-light rounded">
              <h5>Сумма к оплате: <span class="text-primary">{{ totalAmount }} ₽</span></h5>
            </div>
            <button type="submit" class="btn btn-primary-custom w-100 rounded-pill py-2" :disabled="cartItems.length === 0 || loading">
              <span v-if="!loading">Оплатить <i class="bi bi-arrow-right-short"></i></span>
              <span v-else>Обработка...</span>
            </button>
            <p class="small text-muted mt-2">После нажатия произойдет редирект на платежный шлюз ЮKassa</p>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '../store/auth.js'

export default {
  name: 'CartPage',
  setup() {
    const authStore = useAuthStore()
    const cartItems = ref([])
    const loading = ref(false)
    const email = ref('')
    const paymentMethod = ref('bank_card')
    
    // Получаем email пользователя из стора
    const userEmail = computed(() => authStore.currentUser?.email || '')
    
    // Вычисляем общую сумму
    const totalAmount = computed(() => {
      return cartItems.value.reduce((sum, item) => sum + (item.price * item.quantity), 0)
    })
    
    // Загрузка корзины с сервера
    const loadCart = async () => {
      loading.value = true
      try {
        const response = await fetch('http://localhost:8080/api/cart', {
          headers: authStore.getAuthHeader()
        })
        
        if (response.ok) {
          const data = await response.json()
          cartItems.value = data
        } else if (response.status === 401) {
          alert('Сессия истекла, войдите снова')
          authStore.logout()
        }
      } catch (error) {
        console.error('Ошибка загрузки корзины:', error)
      } finally {
        loading.value = false
      }
    }
    
    // Удаление товара из корзины
    const removeItem = async (productId) => {
      try {
        const response = await fetch(`http://localhost:8080/api/cart/${productId}`, {
          method: 'DELETE',
          headers: authStore.getAuthHeader()
        })
        
        if (response.ok) {
          await loadCart() // Перезагружаем корзину
        } else {
          alert('Ошибка при удалении товара')
        }
      } catch (error) {
        console.error('Ошибка удаления:', error)
        alert('Ошибка при удалении товара')
      }
    }
    
    // Оформление заказа и оплата
    const payOrder = async () => {
      if (!email.value && !userEmail.value) {
        alert('Укажите email для отправки чека')
        return
      }
      
      if (cartItems.value.length === 0) {
        alert('Корзина пуста')
        return
      }
      
      loading.value = true
      
      try {
        // Здесь должен быть запрос к вашему бэкенду для создания платежа
        const paymentEmail = email.value || userEmail.value
        const total = totalAmount.value
        
        // Симуляция редиректа на ЮKassa
        const paymentUrl = `https://demo.yookassa.ru/payment?amount=${total}&description=Оплата+заказа+HandMadeStudio&email=${encodeURIComponent(paymentEmail)}`
        
        // Очищаем корзину после успешного заказа
        const clearResponse = await fetch('http://localhost:8080/api/cart/clear', {
          method: 'POST',
          headers: authStore.getAuthHeader()
        })
        
        if (clearResponse.ok) {
          // Редирект на страницу оплаты
          window.location.href = paymentUrl
        } else {
          alert('Ошибка при оформлении заказа')
        }
      } catch (error) {
        console.error('Ошибка оплаты:', error)
        alert('Ошибка при оформлении заказа')
      } finally {
        loading.value = false
      }
    }
    
    // Загружаем корзину при монтировании компонента
    onMounted(() => {
      if (authStore.isAuthenticated) {
        loadCart()
      }
    })
    
    return {
      cartItems,
      loading,
      email,
      paymentMethod,
      userEmail,
      totalAmount,
      removeItem,
      payOrder
    }
  }
}
</script>

<style scoped>
.bg-cream {
  background-color: #f9f2e3;
}

.cart-summary {
  background-color: #f9f2e3;
  border-radius: 24px;
  padding: 1.5rem;
}

.product-img-placeholder {
  background-color: #e9dccd;
  width: 70px;
  height: 70px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  font-size: 1.5rem;
  color: #a9845c;
}

.btn-outline-custom {
  border-color: #b47c4a;
  color: #b47c4a;
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
  border-color: #8b552f;
  color: white;
}

.btn-primary-custom:disabled {
  background-color: #d4a373;
  cursor: not-allowed;
}
</style>