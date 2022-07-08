import { createRouter, createWebHistory } from 'vue-router'
import routes from './routes'

const history = createWebHistory()

const router = createRouter({ history, routes })

// Authorize (Make sure that is the first hook.)
router.beforeEach(to => {

})

router.afterEach(to => {
  const items = [import.meta.env.VITE_TITLE]
  to.meta.title != null && items.unshift(to.meta.title)
  document.title = items.join(' · ')
})

export default router
