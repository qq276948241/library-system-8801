import { reactive, computed, readonly } from 'vue'
import {
  authApi,
  setToken,
  clearToken,
  getToken,
  registerUnauthorizedHandler,
  type User,
  type Role,
} from '@/lib/api'

interface AuthState {
  user: User | null
  ready: boolean
}

const state = reactive<AuthState>({
  user: null,
  ready: false,
})

let restorePromise: Promise<void> | null = null

// 401 全局处理：清空登录态，跳转交给路由守卫 / App.vue 监听
registerUnauthorizedHandler(() => {
  state.user = null
  state.ready = true
  clearToken()
})

export const authState = readonly(state)

export function useAuth() {
  const user = computed(() => state.user)
  const isLogged = computed(() => !!state.user)
  const isAdmin = computed(() => state.user?.role === 'admin')
  const role = computed<Role | null>(() => state.user?.role ?? null)
  const ready = computed(() => state.ready)

  function restore(): Promise<void> {
    if (restorePromise) return restorePromise
    restorePromise = (async () => {
      const token = getToken()
      if (!token) {
        state.ready = true
        return
      }
      try {
        state.user = await authApi.me()
      } catch {
        clearToken()
        state.user = null
      } finally {
        state.ready = true
      }
    })()
    return restorePromise
  }

  async function login(username: string, password: string) {
    const res = await authApi.login(username, password)
    setToken(res.token)
    state.user = res.user
    return res.user
  }

  function logout() {
    state.user = null
    clearToken()
  }

  async function refreshUser() {
    if (!getToken()) return
    try {
      state.user = await authApi.me()
    } catch {
      /* ignore */
    }
  }

  return {
    user,
    role,
    isLogged,
    isAdmin,
    ready,
    login,
    logout,
    restore,
    refreshUser,
  }
}
