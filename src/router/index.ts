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
    },
    {
      path: '/account/:accountId/settings',
      name: 'account-settings',
      component: () => import('../views/accounts/AccountSettingsView.vue'),
    },
    {
      path: '/account/:accountId',
      name: 'account',
      component: () => import('../views/accounts/AccountView.vue'),
    },
    {
      path: '/circle/:circleId/settings',
      name: 'circle-settings',
      component: () => import('../views/circles/CircleSettingsView.vue'),
    },
    {
      path: '/circle/:circleId',
      name: 'circle',
      component: () => import('../views/circles/CircleView.vue'),
    },
  ],
})

export default router
