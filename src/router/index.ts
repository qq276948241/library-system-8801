import { createRouter, createWebHistory, type RouteRecordRaw, type RouteLocationNormalized, type NavigationGuardNext } from 'vue-router'
import FrontLayout from '@/layouts/FrontLayout.vue'
import AdminLayout from '@/layouts/AdminLayout.vue'
import LoginPage from '@/pages/LoginPage.vue'
import BooksPage from '@/pages/BooksPage.vue'
import BorrowsPage from '@/pages/BorrowsPage.vue'
import ProfilePage from '@/pages/ProfilePage.vue'
import AdminBooksPage from '@/pages/AdminBooksPage.vue'
import AdminUsersPage from '@/pages/AdminUsersPage.vue'
import AdminBorrowsPage from '@/pages/AdminBorrowsPage.vue'
import { useAuth } from '@/composables/useAuth'

let restorePromise: Promise<void> | null = null

async function ensureReady() {
  const { restore } = useAuth()
  if (!restorePromise) {
    restorePromise = restore()
  }
  await restorePromise
}

async function requireAuth(to: RouteLocationNormalized, _from: RouteLocationNormalized, next: NavigationGuardNext) {
  await ensureReady()
  const { isLogged } = useAuth()
  if (!isLogged.value) {
    next({ name: 'login', query: { redirect: to.fullPath } })
  } else {
    next()
  }
}

async function requireAdmin(to: RouteLocationNormalized, _from: RouteLocationNormalized, next: NavigationGuardNext) {
  await ensureReady()
  const { isLogged, isAdmin } = useAuth()
  if (!isLogged.value) {
    next({ name: 'login', query: { redirect: to.fullPath } })
  } else if (!isAdmin.value) {
    next({ name: 'books' })
  } else {
    next()
  }
}

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'login',
    component: LoginPage,
  },
  {
    path: '/',
    component: FrontLayout,
    beforeEnter: requireAuth,
    children: [
      { path: '', redirect: { name: 'books' } },
      { path: 'books', name: 'books', component: BooksPage },
      { path: 'borrows', name: 'borrows', component: BorrowsPage },
      { path: 'profile', name: 'profile', component: ProfilePage },
    ],
  },
  {
    path: '/admin',
    component: AdminLayout,
    beforeEnter: requireAdmin,
    children: [
      { path: '', redirect: { name: 'admin-borrows' } },
      { path: 'borrows', name: 'admin-borrows', component: AdminBorrowsPage },
      { path: 'books', name: 'admin-books', component: AdminBooksPage },
      { path: 'users', name: 'admin-users', component: AdminUsersPage },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
