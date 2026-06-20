import { createRouter, createWebHistory } from 'vue-router'
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

const { isLogged, isAdmin, ready, restore } = useAuth()

let restorePromise = null

async function ensureReady() {
  if (!restorePromise) {
    restorePromise = restore()
  }
  await restorePromise
  if (!ready.value) {
    await new Promise((resolve) => {
      const check = () => {
        if (ready.value) resolve()
        else setTimeout(check, 50)
      }
      check()
    })
  }
}

async function requireAuth(to, _from, next) {
  await ensureReady()
  if (!isLogged.value) {
    next({ name: 'login', query: { redirect: to.fullPath } })
  } else {
    next()
  }
}

async function requireAdmin(to, _from, next) {
  await ensureReady()
  if (!isLogged.value) {
    next({ name: 'login', query: { redirect: to.fullPath } })
  } else if (!isAdmin.value) {
    next({ name: 'books' })
  } else {
    next()
  }
}

const routes = [
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
