import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/meals/recipes',
      name: 'recipes',
      component: () => import('../views/meals/RecipesView.vue'),
    },
    {
      path: '/meals/ingredients',
      name: 'ingredients',
      component: () => import('../views/meals/IngredientsView.vue'),
    },
    {
      path: '/calendar',
      name: 'calendar',
      component: () => import('../views/calendar/CalendarView.vue'),
    }
  ],
})

export default router
