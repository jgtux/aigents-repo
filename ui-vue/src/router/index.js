import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/pages/Home.vue'
import About from '@/pages/About.vue'
import Agents from '@/pages/Agents.vue'
import Myprojects from '@/pages/Myprojects.vue'
import Create from '@/pages/Create.vue'
import Login from '@/pages/Login.vue'
import Signup from '@/pages/Signup.vue'
import Chat from '@/pages/Chat.vue'
import { checkAuth } from '@/api/auth'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home,
    meta: { title: 'AIgents Home' }
  },
  {
    path: '/about',
    name: 'About',
    component: About,
    meta: { title: 'AIgents About Us' }
  },
  {
    path: '/agents',
    name: 'Agents',
    component: Agents,
    meta: { title: 'AIgents - Agents' }
  },
  {
    path: '/myprojects',
    name: 'Myprojects',
    component: Myprojects,
    meta: { 
      title: 'AIgents - My Projects',
      requiresAuth: true
    }
  },
  {
    path: '/create',
    name: 'Create',
    component: Create,
    meta: { 
      title: 'Create Agent - AIgents',
      requiresAuth: true
    }
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { title: 'AIgents Log-in' }
  },
  {
    path: '/signup',
    name: 'Signup',
    component: Signup,
    meta: { title: 'AIgents Sign-up' }
  },
  {
    path: '/chat/:agent_uuid',
    name: 'Chat',
    component: Chat,
    meta: { 
      title: 'Chat - AIgents',
      requiresAuth: true
    }
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/'
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
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
  console.log('=== NAVIGATION GUARD ===')
  console.log('From:', from.path)
  console.log('To:', to.path)
  
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
  console.log('Requires auth:', requiresAuth)
  
  if (requiresAuth) {
    console.log('Checking authentication...')
    try {
      const result = await checkAuth()
      console.log('Auth check result:', result)
      console.log('✅ User authenticated, allowing access')
      next()
    } catch (error) {
      console.log('❌ Auth check failed:', error.message)
      console.log('Redirecting to /login with redirect:', to.fullPath)
      next({
        path: '/login',
        query: { redirect: to.fullPath }
      })
    }
  } else {
    console.log('Public route, allowing access')
    next()
  }
})

// Atualiza o título da página
router.afterEach((to) => {
  document.title = to.meta.title || 'AIgents'
})

export default router
