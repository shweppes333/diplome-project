<template>
  <div>
    <Header />
    <main class="container py-4">
      <router-view />
    </main>
    <Footer />
  </div>
</template>

<script>
import Header from './components/header.vue'
import Footer from './components/footer.vue'

export default {
  name: 'App',
  components: {
    Header,
    Footer
  }
}
</script>