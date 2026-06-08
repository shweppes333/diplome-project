import HomePage from '../views/startpage.vue'
import RegisterPage from '../views/register.vue'
import LoginPage from '../views/auth.vue'
import CartPage from '../views/cart.vue'

const routes = [
  { path: '/', component: HomePage, name: 'home' },
  { path: '/register', component: RegisterPage, name: 'register' },
  { path: '/auth', component: LoginPage, name: 'login' },
  { path: '/cart', component: CartPage, name: 'cart' }
]

export default routes