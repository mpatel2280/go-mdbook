const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api'

export function getToken() {
  return localStorage.getItem('token')
}

export function setAuth(token, role, email) {
  localStorage.setItem('token', token)
  localStorage.setItem('role', role)
  localStorage.setItem('email', email)
}

export function clearAuth() {
  localStorage.removeItem('token')
  localStorage.removeItem('role')
  localStorage.removeItem('email')
}

export function getRole() {
  return localStorage.getItem('role') || ''
}

export function getEmail() {
  return localStorage.getItem('email') || ''
}

async function request(path, options = {}) {
  const headers = { 'Content-Type': 'application/json', ...(options.headers || {}) }
  const token = getToken()
  if (token) headers.Authorization = `Bearer ${token}`
  const res = await fetch(`${API_URL}${path}`, { ...options, headers })
  if (!res.ok) {
    let errMsg = 'Request failed'
    try {
      const data = await res.json()
      errMsg = data.error || errMsg
    } catch {
      // ignore
    }
    throw new Error(errMsg)
  }
  if (res.status === 204) return null
  return res.json()
}

export const api = {
  login: (payload) => request('/auth/login', { method: 'POST', body: JSON.stringify(payload) }),
  me: () => request('/me'),
  listBooks: () => request('/books'),
  getBook: (id) => request(`/books/${id}`),
  listUsers: () => request('/admin/users'),
  createUser: (payload) => request('/admin/users', { method: 'POST', body: JSON.stringify(payload) }),
  updateUser: (id, payload) => request(`/admin/users/${id}`, { method: 'PATCH', body: JSON.stringify(payload) }),
  deleteUser: (id) => request(`/admin/users/${id}`, { method: 'DELETE' }),
  createBook: (payload) => request('/admin/books', { method: 'POST', body: JSON.stringify(payload) }),
  updateBook: (id, payload) => request(`/admin/books/${id}`, { method: 'PATCH', body: JSON.stringify(payload) }),
  deleteBook: (id) => request(`/admin/books/${id}`, { method: 'DELETE' }),
  buildBook: (id) => request(`/admin/books/${id}/build`, { method: 'POST' }),
  uploadBook: async (id, file) => {
    const token = getToken()
    const form = new FormData()
    form.append('file', file)
    const res = await fetch(`${API_URL}/admin/books/${id}/upload`, {
      method: 'POST',
      headers: token ? { Authorization: `Bearer ${token}` } : {},
      body: form
    })
    if (!res.ok) {
      let errMsg = 'Upload failed'
      try {
        const data = await res.json()
        errMsg = data.error || errMsg
      } catch {
        // ignore
      }
      throw new Error(errMsg)
    }
    return res.json()
  }
}

export function bookContentUrl(id, path = '') {
  const suffix = path ? `/${path.replace(/^\/+/, '')}` : ''
  return `${API_URL}/books/${id}/content${suffix}`
}
