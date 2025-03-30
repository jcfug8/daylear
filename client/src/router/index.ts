import { createRouter, createWebHistory } from 'vue-router'
import LogInView from '../views/login/LogInView.vue'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: LogInView,
      meta: {
        requiresNoAuth: true,
      },
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/login/LogInView.vue'),
      meta: {
        requiresNoAuth: true,
      },
    },
    {
      path: '/meals/recipes',
      name: 'recipes',
      component: () => import('../views/meals/RecipesView.vue'),
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/meals/recipes/:recipeId',
      name: 'recipe',
      component: () => import('../views/meals/RecipeView.vue'),
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/meals/ingredients',
      name: 'ingredients',
      component: () => import('../views/meals/IngredientsView.vue'),
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/calendar',
      name: 'calendar',
      component: () => import('../views/calendar/CalendarView.vue'),
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/account/:accountId/settings',
      name: 'account-settings',
      component: () => import('../views/accounts/AccountSettingsView.vue'),
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/account/:accountId',
      name: 'account',
      component: () => import('../views/accounts/AccountView.vue'),
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/circle/:circleId/settings',
      name: 'circle-settings',
      component: () => import('../views/circles/CircleSettingsView.vue'),
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/circle/:circleId',
      name: 'circle',
      component: () => import('../views/circles/CircleView.vue'),
      meta: {
        requiresAuth: true,
      },
    },
  ],
})

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  if (to.meta.requiresAuth) {
    if (authStore.isLoggedIn) {
      // User is authenticated, proceed to the route
      next()
    } else {
      // User is not authenticated, redirect to login
      next('/login')
    }
  } else if (to.meta.requiresNoAuth) {
    if (authStore.isLoggedIn) {
      // User is authenticated, redirect to home
      next('/calendar')
    } else {
      // User is not authenticated, allow access
      next()
    }
  } else {
    // Non-protected route, allow access
    next()
  }
})

export default router
