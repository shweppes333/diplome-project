<template>
  <div>
    <h2 class="mb-4">🛒 Моя корзина</h2>
    <div class="row g-4">
      <!-- Левая часть: товары -->
      <div class="col-md-7">
        <div class="bg-cream p-3 rounded-4" style="min-height: 300px;">
          <!-- Состояние загрузки -->
          <div v-if="loading" class="text-center py-5">
            <div class="spinner-border text-primary" role="status">
              <span class="visually-hidden">Загрузка...</span>
            </div>
            <p class="mt-2">Загрузка корзины...</p>
          </div>
          
          <!-- Пустая корзина -->
          <div v-else-if="!cartItems || cartItems.length === 0" class="text-center py-5">
            <i class="bi bi-bag-x fs-1 text-muted"></i>
            <p class="mt-2">Ваша корзина пуста. Добавьте товары на главной странице.</p>
            <router-link to="/" class="btn btn-outline-custom rounded-pill">На главную</router-link>
          </div>
          
          <!-- Список товаров -->
          <div v-else>
            <div v-for="item in cartItems" :key="item.id" class="d-flex justify-content-between align-items-center border-bottom pb-3 mb-3">
              <div class="d-flex gap-3 align-items-center">
                <div class="product-img-placeholder" style="width: 70px; height: 70px; font-size: 1.5rem;">
                  <i :class="item.product_icon || 'bi bi-box'"></i>
                </div>
                <div>
                  <h5 class="mb-0">{{ item.product_name }}</h5>
                  <small class="text-muted">{{ item.price }} ₽ x {{ item.quantity }}</small>
                </div>
              </div>
              <div>
                <strong class="text-primary">{{ item.price * item.quantity }} ₽</strong>
                <button 
                  class="btn btn-sm btn-outline-danger ms-2 rounded-circle" 
                  @click="removeItem(item.product_id)"
                  :disabled="removingItemId === item.product_id"
                >
                  <span v-if="removingItemId === item.product_id">
                    <span class="spinner-border spinner-border-sm" role="status"></span>
                  </span>
                  <span v-else>
                    <i class="bi bi-trash"></i>
                  </span>
                </button>
              </div>
            </div>
            
            <!-- Итоговая сумма -->
            <div class="text-end mt-4 pt-3 border-top">
              <h4>Итого: <span class="text-primary">{{ totalAmount }} ₽</span></h4>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Правая часть: форма оплаты -->
      <div class="col-md-5">
        <div class="cart-summary">
          <h4><i class="bi bi-credit-card"></i> Данные для заказа</h4>
          
          <div class="mb-3">
            <label class="form-label">Email для чека</label>
            <input 
              type="email" 
              class="form-control" 
              v-model="email" 
              required 
              :placeholder="userEmail || 'customer@example.com'"
            >
            <small class="text-muted">На этот email придет подтверждение заказа</small>
          </div>
          
          <div class="mb-3">
            <label class="form-label">Способ оплаты</label>
            <select class="form-select" v-model="paymentMethod" required>
              <option value="bank_card">Банковская карта</option>
              <option value="yoomoney">ЮMoney</option>
              <option value="sberpay">СберПэй</option>
              <option value="cash">Наличными при получении</option>
            </select>
          </div>
          
          <div class="mb-3">
            <label class="form-label">Адрес доставки</label>
            <input 
              type="text" 
              class="form-control" 
              v-model="address" 
              placeholder="г. Москва, ул. Примерная, д. 1, кв. 1"
            >
          </div>
          
          <div class="mb-3">
            <label class="form-label">Телефон</label>
            <input 
              type="tel" 
              class="form-control" 
              v-model="phone" 
              placeholder="+7 (999) 123-45-67"
            >
          </div>
          
          <!-- Блок с суммой заказа -->
          <div class="mb-3 p-3 bg-light rounded">
            <div class="d-flex justify-content-between mb-2">
              <span>Товаров в корзине:</span>
              <span class="fw-bold">{{ cartItems && cartItems.length ? cartItems.length : 0 }} шт.</span>
            </div>
            <div class="d-flex justify-content-between">
              <span>Сумма заказа:</span>
              <span class="fw-bold fs-5 text-primary">{{ totalAmount }} ₽</span>
            </div>
          </div>
          
          <button 
            type="button" 
            class="btn btn-primary-custom w-100 rounded-pill py-2" 
            :disabled="!cartItems || cartItems.length === 0 || processing"
            @click="placeOrder"
          >
            <span v-if="!processing">
              <i class="bi bi-check-circle"></i> Оформить заказ
            </span>
            <span v-else>
              <span class="spinner-border spinner-border-sm me-2" role="status"></span>
              Оформление...
            </span>
          </button>
          
          <p class="small text-muted mt-2 text-center">
            Нажимая кнопку, вы соглашаетесь с условиями обработки заказа
          </p>
        </div>
      </div>
    </div>
    
    <!-- Модальное окно для отображения информации о заказе -->
    <div class="modal fade" id="orderModal" tabindex="-1" aria-hidden="true" ref="orderModal">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
          <div class="modal-header bg-success text-white">
            <h5 class="modal-title">
              <i class="bi bi-check-circle-fill"></i> Заказ оформлен!
            </h5>
            <button type="button" class="btn-close btn-close-white" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <div class="text-center mb-3">
              <i class="bi bi-emoji-smile" style="font-size: 4rem; color: #b47c4a;"></i>
            </div>
            <h5 class="text-center mb-3">Спасибо за ваш заказ!</h5>
            <p><strong>Номер заказа:</strong> {{ orderNumber }}</p>
            <p><strong>Сумма заказа:</strong> {{ totalAmount }} ₽</p>
            <p><strong>Способ оплаты:</strong> {{ getPaymentMethodText() }}</p>
            <p><strong>Email для связи:</strong> {{ email || userEmail }}</p>
            <p v-if="address"><strong>Адрес доставки:</strong> {{ address }}</p>
            <p v-if="phone"><strong>Телефон:</strong> {{ phone }}</p>
            <hr>
            <h6>Состав заказа:</h6>
            <ul class="list-unstyled">
              <li v-for="item in lastOrderItems" :key="item.id" class="mb-1">
                {{ item.product_name }} - {{ item.quantity }} шт. x {{ item.price }} ₽ = 
                <strong>{{ item.price * item.quantity }} ₽</strong>
              </li>
            </ul>
            <div class="alert alert-info mt-3 mb-0">
              <i class="bi bi-info-circle"></i> 
              Подтверждение заказа отправлено на вашу почту. Наш менеджер свяжется с вами в ближайшее время.
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-primary-custom" data-bs-dismiss="modal" @click="goToHome">
              На главную
            </button>
            <button type="button" class="btn btn-outline-secondary" data-bs-dismiss="modal">
              Закрыть
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../store/auth.js'
import { Modal } from 'bootstrap'

export default {
  name: 'Cart',
  setup() {
    const router = useRouter()
    const authStore = useAuthStore()
    
    // Состояния - всегда инициализируем пустым массивом
    const cartItems = ref([])
    const loading = ref(false)
    const processing = ref(false)
    const removingItemId = ref(null)
    const email = ref('')
    const paymentMethod = ref('bank_card')
    const address = ref('')
    const phone = ref('')
    const orderNumber = ref('')
    const orderModal = ref(null)
    const lastOrderItems = ref([])
    
    // Email пользователя из store
    const userEmail = computed(() => authStore.currentUser?.email || '')
    
    // Общая сумма - защита от null
    const totalAmount = computed(() => {
      if (!cartItems.value || !Array.isArray(cartItems.value)) return 0
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
          // Всегда делаем массив, даже если пришел null
          cartItems.value = Array.isArray(data) ? data : []
          console.log('Корзина загружена:', cartItems.value.length, 'товаров')
        } else if (response.status === 401) {
          alert('Сессия истекла, войдите снова')
          authStore.logout()
        } else {
          cartItems.value = []
        }
      } catch (error) {
        console.error('Ошибка загрузки корзины:', error)
        cartItems.value = []
      } finally {
        loading.value = false
      }
    }
    
    // Удаление товара из корзины
    const removeItem = async (productId) => {
      // Проверяем, что товар существует
      if (!productId) {
        console.error('Некорректный ID товара')
        return
      }
      
      removingItemId.value = productId
      
      try {
        const response = await fetch(`http://localhost:8080/api/cart/${productId}`, {
          method: 'DELETE',
          headers: authStore.getAuthHeader()
        })
        
        if (response.ok) {
          // Перезагружаем корзину после удаления
          await loadCart()
          console.log('Товар удален, корзина обновлена')
        } else {
          const error = await response.json()
          alert(error.error || 'Ошибка при удалении товара')
        }
      } catch (error) {
        console.error('Ошибка удаления:', error)
        alert('Ошибка при удалении товара')
      } finally {
        removingItemId.value = null
      }
    }
    
    // Получение текста способа оплаты
    const getPaymentMethodText = () => {
      const methods = {
        'bank_card': 'Банковская карта',
        'yoomoney': 'ЮMoney',
        'sberpay': 'СберПэй',
        'cash': 'Наличными при получении'
      }
      return methods[paymentMethod.value] || paymentMethod.value
    }
    
    // Генерация номера заказа
    const generateOrderNumber = () => {
      const date = new Date()
      const year = date.getFullYear()
      const month = String(date.getMonth() + 1).padStart(2, '0')
      const day = String(date.getDate()).padStart(2, '0')
      const random = Math.floor(Math.random() * 10000).toString().padStart(4, '0')
      return `${year}${month}${day}-${random}`
    }
    
    // Оформление заказа
    const placeOrder = async () => {
      const orderEmail = email.value || userEmail.value
      
      if (!orderEmail) {
        alert('Укажите email для отправки подтверждения заказа')
        return
      }
      
      if (!cartItems.value || cartItems.value.length === 0) {
        alert('Корзина пуста')
        return
      }
      
      processing.value = true
      
      try {
        // Сохраняем копию заказа для отображения
        lastOrderItems.value = [...cartItems.value]
        
        // Генерируем номер заказа
        orderNumber.value = generateOrderNumber()
        
        // Создаем объект заказа
        const orderData = {
          order_number: orderNumber.value,
          email: orderEmail,
          phone: phone.value,
          address: address.value,
          payment_method: paymentMethod.value,
          total_amount: totalAmount.value,
          items: cartItems.value.map(item => ({
            id: item.product_id,
            name: item.product_name,
            price: item.price,
            quantity: item.quantity,
            total: item.price * item.quantity
          }))
        }
        
        console.log('Отправка заказа:', orderData)
        
        // Очищаем корзину на сервере
        const clearResponse = await fetch('http://localhost:8080/api/cart/clear', {
          method: 'POST',
          headers: authStore.getAuthHeader()
        })
        
        if (clearResponse.ok) {
          // Очищаем локальное состояние
          cartItems.value = []
          
          // Показываем модальное окно
          const modalElement = orderModal.value
          if (modalElement) {
            const modal = new Modal(modalElement)
            modal.show()
          } else {
            // Если модальное окно не найдено, показываем alert
            alert(`Заказ №${orderNumber.value} успешно оформлен на сумму ${totalAmount.value} ₽`)
            router.push('/')
          }
        } else {
          alert('Ошибка при оформлении заказа. Попробуйте снова.')
        }
      } catch (error) {
        console.error('Ошибка оформления заказа:', error)
        alert('Ошибка при оформлении заказа')
      } finally {
        processing.value = false
      }
    }
    
    // Переход на главную страницу
    const goToHome = () => {
      router.push('/')
    }
    
    // Загружаем корзину при монтировании компонента
    onMounted(() => {
      if (authStore.isAuthenticated) {
        loadCart()
      } else {
        cartItems.value = []
      }
    })
    
    return {
      cartItems,
      loading,
      processing,
      removingItemId,
      email,
      paymentMethod,
      address,
      phone,
      userEmail,
      totalAmount,
      orderNumber,
      orderModal,
      lastOrderItems,
      removeItem,
      placeOrder,
      getPaymentMethodText,
      goToHome
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
  position: sticky;
  top: 20px;
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

.btn-primary-custom:hover:not(:disabled) {
  background-color: #9b623a;
  color: white;
}

.btn-primary-custom:disabled {
  background-color: #d4a373;
  cursor: not-allowed;
}

.btn-outline-danger {
  border: 1px solid #dc3545;
  color: #dc3545;
  background: transparent;
}

.btn-outline-danger:hover:not(:disabled) {
  background-color: #dc3545;
  color: white;
}

.btn-outline-danger:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.border-bottom {
  border-bottom: 1px solid #e0cfb4 !important;
}

.border-top {
  border-top: 1px solid #e0cfb4 !important;
}

.bg-light {
  background-color: #fffcf7 !important;
}

.modal-header.bg-success {
  background-color: #28a745 !important;
}

.modal-content {
  border-radius: 16px;
  border: none;
}

.btn-outline-secondary {
  border: 1px solid #b47c4a;
  color: #b47c4a;
}

.btn-outline-secondary:hover {
  background-color: #b47c4a;
  color: white;
}
</style>