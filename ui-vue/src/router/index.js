import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/pages/Home.vue'
import About from '@/pages/About.vue'
import Agents from '@/pages/Agents.vue'
import Myprojects from '@/pages/Myprojects.vue'
import Create from '@/pages/Create.vue'
import Login from '@/pages/Login.vue'
import Signup from '@/pages/Signup.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home,
    meta: {
      title: 'AIgents Home'
    }
  },
  {
    path: '/about',
    name: 'About',
    component: About,
    meta: {
      title: 'AIgents About Us'
    }
  },
  {
    path: '/agents',
    name: 'Agents',
    component: Agents,
    meta: {
      title: 'AIgents - Agents'
    }
  },
  {
    path: '/myprojects',
    name: 'Myprojects',
    component: Myprojects,
    meta: {
      title: 'AIgents - My Projects',
      requiresAuth: true // Página privada
    }
  },
  {
    path: '/create',
    name: 'Create',
    component: Create,
    meta: {
      title: 'Create Agent - AIgents',
      requiresAuth: true // Página privada
    }
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: {
      title: 'AIgents Log-in'
    }
  },
  {
    path: '/signup',
    name: 'Signup',
    component: Signup,
    meta: {
      title: 'AIgents Sign-up'
    }
  },
  {
    // Rota 404 - Redireciona para Home se a página não existir
    path: '/:pathMatch(.*)*',
    redirect: '/'
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
  // Scroll para o topo ao mudar de página
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  }
})

// Navigation Guard - Protege rotas que requerem autenticação
router.beforeEach(async (to, from, next) => {
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
  
  if (requiresAuth) {
    try {
      // Verifica autenticação fazendo uma requisição ao backend
      const response = await fetch('/api/auth/check', {
        method: 'GET',
        credentials: 'include', // Importante: envia cookies HTTP-only
        headers: {
          'Content-Type': 'application/json'
        }
      })
      
      if (response.ok) {
        // Usuário autenticado
        next()
      } else {
        // Não autenticado, redireciona para login
        next('/login')
      }
    } catch (error) {
      console.error('Auth check failed:', error)
      next('/login')
    }
  } else if (to.path === '/login' || to.path === '/signup') {
    // Opcional: verifica se já está autenticado antes de mostrar login/signup
    try {
      const response = await fetch('/api/auth/check', {
        method: 'GET',
        credentials: 'include'
      })
      
      if (response.ok) {
        // Já está autenticado, redireciona para home
        next('/')
      } else {
        next()
      }
    } catch (error) {
      // Erro na verificação, permite acesso ao login/signup
      next()
    }
  } else {
    next()
  }
})

// Atualiza o título da página
router.afterEach((to) => {
  document.title = to.meta.title || 'AIgents'
})

export default router
