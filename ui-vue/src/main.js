import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

// Import global styles
import '@/assets/styles/styles.css'

// Optional: Import vue-meta for meta tags (install with: npm install vue-meta@next)
// import { createMetaManager } from 'vue-meta'

const app = createApp(App)

app.use(router)

// Optional: If using vue-meta for meta tags
// app.use(createMetaManager())

app.mount('#app')
