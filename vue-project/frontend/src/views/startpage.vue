<template>
  <div>
    <!-- Баннер -->
    <div class="banner p-4 p-md-5 mb-5 text-center">
      <h1 class="display-5 fw-semibold">Уникальные изделия ручной работы</h1>
      <p class="lead">Керамика, текстиль, дерево — тепло ваших рук в каждой детали</p>
      <router-link to="/register" class="btn btn-primary-custom btn-lg rounded-pill px-5 mt-2">
        Стать покупателем
      </router-link>
    </div>
    
  
    <h2 class="mb-4"> Популярные товары</h2>
    <div class="row g-4 mb-5">
      <div class="col-md-6 col-lg-3" v-for="product in products" :key="product.id">
        <ProductCard :product="product" />
      </div>
    </div>
    
    
  </div>
</template>

<script>
import ProductCard from '../components/cardComponent.vue'
import MapComponent from '../components/mapComponent.vue'

export default {
  name: 'HomePage',
  components: {
    ProductCard,
    MapComponent
  },
  data() {
    return {
      products: [
        { id: 1, name: "Керамическая кружка", description: "Ручная лепка, глазурь", price: 1250, icon: "bi bi-cup-hot-fill" },
        { id: 2, name: "Вязаный плед", description: "Хлопок, мягкий", price: 3450, icon: "bi bi-border-width" },
        { id: 3, name: "Деревянная ложка", description: "Ручная резьба", price: 780, icon: "bi bi-egg-fried" },
        { id: 4, name: "Свеча ручной работы", description: "Аромат ванили", price: 540, icon: "bi bi-candle" }
      ]
    }
  }
}
</script>