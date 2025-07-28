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
      path: '/recipes',
      name: 'recipes',
      component: () => import('../views/meals/recipes/RecipesView.vue'),
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/recipes/create',
      name: 'recipeCreate',
      component: () => import('../views/meals/recipes/RecipeCreateView.vue'),
      meta: {
        requiresAuth: true,
        breadcrumbs: true,
      },
    },
    {
      path: '/circles/:circleId/recipes/create',
      name: 'circleRecipeCreate',
      component: () => import('../views/meals/recipes/RecipeCreateView.vue'),
      meta: {
        requiresAuth: true,
        breadcrumbs: true,
      },
    },
    {
      path: '/recipes/:recipeId',
      name: 'recipe',
      component: () => import('../views/meals/recipes/RecipeView.vue'),
      meta: {
        requiresAuth: true,
        breadcrumbs: true,
      },
    },
    {
      path: '/users/:userId/recipes/:recipeId',
      name: 'userRecipe',
      component: () => import('../views/meals/recipes/RecipeView.vue'),
      meta: {
        requiresAuth: true,
        breadcrumbs: true,
      },
    },
    {
      path: '/circles/:circleId/recipes/:recipeId',
      name: 'circleRecipe',
      component: () => import('../views/meals/recipes/RecipeView.vue'),
      meta: {
        requiresAuth: true,
        breadcrumbs: true,
      },
    },
    {
      path: '/recipes/:recipeId/edit',
      name: 'recipeEdit',
      component: () => import('../views/meals/recipes/RecipeEditView.vue'),
      meta: {
        requiresAuth: true,
        breadcrumbs: true,
      },
    },
    {
      path: '/circles/:circleId/recipes/:recipeId/edit',
      name: 'circleRecipeEdit',
      component: () => import('../views/meals/recipes/RecipeEditView.vue'),
      meta: {
        requiresAuth: true,
        breadcrumbs: true,
      },
    },
    {
      path: '/ingredients',
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
      path: '/users/:userId/edit',
      name: 'user-edit',
      component: () => import('../views/users/UserEditView.vue'),
      meta: {
        requiresAuth: true,
        breadcrumbs: true,
      },
    },
    {
      path: '/users/:userId',
      name: 'user',
      component: () => import('../views/users/UserView.vue'),
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/users/:userId/users/:friendUserId',
      name: 'userFriend',
      component: () => import('../views/users/UserView.vue'),
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/users',
      name: 'users',
      component: () => import('../views/users/UsersView.vue'),
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/circles/:circleId/edit',
      name: 'circle-edit',
      component: () => import('../views/circles/CircleEditView.vue'),
      meta: {
        requiresAuth: true,
        breadcrumbs: true,
      },
    },
    {
      path: '/circles/create',
      name: 'circleCreate',
      component: () => import('../views/circles/CircleCreateView.vue'),
      meta: {
        requiresAuth: true,
        breadcrumbs: true,
      },
    },
    {
      path: '/circles',
      name: 'circles',
      component: () => import('../views/circles/CirclesView.vue'),
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/circles/:circleId',
      name: 'circle',
      component: () => import('../views/circles/CircleView.vue'),
      meta: {
        requiresAuth: true,
        breadcrumbs: true,
      },
    },
    {
      path: '/api-docs',
      name: 'api-docs',
      component: () => import('../views/ApiDocsView.vue'),
      meta: {
        requiresNoAuth: true,
      },
    },
  ],
})

router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  
  // Wait for auth initialization to complete
  await authStore.waitForAuthInit()
  
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
      next('/recipes')
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
