// 青檀图书馆 API 客户端 —— 与 Go 后端契约一一对应
// 统一响应信封：{ code, msg, data }，code === 0 视为成功

export type Role = 'admin' | 'student'

export interface User {
  id: number
  username: string
  name: string
  role: Role
  email: string
  created_at: string
}

export interface Book {
  id: number
  title: string
  author: string
  isbn: string
  category: string
  description: string
  cover_color: string
  publisher: string
  published_year: number
  total_copies: number
  available_copies: number
  created_at: string
}

export type BorrowStatus = 'borrowed' | 'overdue' | 'returned'

export interface Borrow {
  id: number
  user_id: number
  book_id: number
  borrow_date: string
  due_date: string
  return_date: string | null
  status: BorrowStatus
  user?: User
  book?: Book
}

export interface CategoryCount {
  category: string
  count: number
}

export interface Stats {
  total_books: number
  total_copies: number
  available_copies: number
  total_users: number
  active_borrows: number
  overdue_borrows: number
  categories: CategoryCount[]
}

export interface PageData<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}

export interface ApiEnvelope<T> {
  code: number
  msg: string
  data: T
}

export interface LoginResult {
  token: string
  user: User
}

const TOKEN_KEY = 'qt_library_token'

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

export function setToken(token: string) {
  localStorage.setItem(TOKEN_KEY, token)
}

export function clearToken() {
  localStorage.removeItem(TOKEN_KEY)
}

export class ApiError extends Error {
  code: number
  constructor(code: number, msg: string) {
    super(msg)
    this.code = code
    this.name = 'ApiError'
  }
}

// 401 时触发的全局回调（由 useAuth 注册），避免循环依赖
let onUnauthorized: (() => void) | null = null
export function registerUnauthorizedHandler(fn: () => void) {
  onUnauthorized = fn
}

interface RequestOptions {
  method?: string
  body?: unknown
  query?: Record<string, string | number | undefined | null>
}

function buildUrl(path: string, query?: RequestOptions['query']): string {
  if (!query) return path
  const params = new URLSearchParams()
  Object.entries(query).forEach(([k, v]) => {
    if (v !== undefined && v !== null && v !== '') params.append(k, String(v))
  })
  const qs = params.toString()
  return qs ? `${path}?${qs}` : path
}

async function request<T>(path: string, opts: RequestOptions = {}): Promise<T> {
  const headers: Record<string, string> = { 'Content-Type': 'application/json' }
  const token = getToken()
  if (token) headers['Authorization'] = `Bearer ${token}`

  let res: Response
  try {
    res = await fetch(buildUrl(path, opts.query), {
      method: opts.method || 'GET',
      headers,
      body: opts.body !== undefined ? JSON.stringify(opts.body) : undefined,
    })
  } catch (e) {
    throw new ApiError(-1, '网络请求失败，请检查服务是否启动')
  }

  let json: ApiEnvelope<T>
  try {
    json = (await res.json()) as ApiEnvelope<T>
  } catch {
    throw new ApiError(-1, `服务响应解析失败 (${res.status})`)
  }

  if (json.code === 401 && onUnauthorized) onUnauthorized()
  if (json.code !== 0) throw new ApiError(json.code, json.msg)
  return json.data
}

/* ---------- 认证 ---------- */
export const authApi = {
  login: (username: string, password: string) =>
    request<LoginResult>('/api/auth/login', { method: 'POST', body: { username, password } }),
  me: () => request<User>('/api/auth/me'),
  changePassword: (old_password: string, new_password: string) =>
    request<{ ok: boolean }>('/api/auth/password', { method: 'PUT', body: { old_password, new_password } }),
}

/* ---------- 图书 ---------- */
export interface BookQuery {
  page?: number
  pageSize?: number
  keyword?: string
  category?: string
}

export interface BookInput {
  title: string
  author: string
  isbn?: string
  category?: string
  description?: string
  cover_color?: string
  publisher?: string
  published_year?: number
  total_copies?: number
}

export const booksApi = {
  list: (q: BookQuery = {}) =>
    request<PageData<Book>>('/api/books', { query: q as Record<string, string | number> }),
  categories: () => request<CategoryCount[]>('/api/books/categories'),
  get: (id: number) => request<Book>(`/api/books/${id}`),
  create: (data: BookInput) => request<Book>('/api/books', { method: 'POST', body: data }),
  update: (id: number, data: BookInput) => request<Book>(`/api/books/${id}`, { method: 'PUT', body: data }),
  remove: (id: number) => request<{ ok: boolean }>(`/api/books/${id}`, { method: 'DELETE' }),
}

/* ---------- 借阅 ---------- */
export interface BorrowQuery {
  page?: number
  pageSize?: number
  status?: BorrowStatus | ''
  userId?: number | ''
}

export const borrowsApi = {
  list: (q: BorrowQuery = {}) =>
    request<PageData<Borrow>>('/api/borrows', { query: q as Record<string, string | number> }),
  borrow: (book_id: number) =>
    request<{ ok: boolean }>('/api/borrows', { method: 'POST', body: { book_id } }),
  return: (id: number) =>
    request<{ ok: boolean }>(`/api/borrows/${id}/return`, { method: 'POST' }),
}

/* ---------- 用户 ---------- */
export interface UserQuery {
  page?: number
  pageSize?: number
  keyword?: string
}

export interface UserInput {
  username?: string
  password?: string
  name?: string
  role?: Role
  email?: string
}

export const usersApi = {
  list: (q: UserQuery = {}) =>
    request<PageData<User>>('/api/users', { query: q as Record<string, string | number> }),
  create: (data: UserInput) => request<User>('/api/users', { method: 'POST', body: data }),
  update: (id: number, data: UserInput) =>
    request<User>(`/api/users/${id}`, { method: 'PUT', body: data }),
  remove: (id: number) => request<{ ok: boolean }>(`/api/users/${id}`, { method: 'DELETE' }),
}

/* ---------- 统计 ---------- */
export const statsApi = {
  overview: () => request<Stats>('/api/stats'),
}
