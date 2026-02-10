import { useEffect, useMemo, useState } from 'react'
import { api, bookContentUrl, clearAuth, getEmail, getRole, setAuth } from './api/client'

function App() {
  const [auth, setAuthState] = useState({ token: localStorage.getItem('token') || '' })
  const [books, setBooks] = useState([])
  const [users, setUsers] = useState([])
  const [selectedBook, setSelectedBook] = useState(null)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [activeModule, setActiveModule] = useState('books')
  const [booksView, setBooksView] = useState('list')
  const [usersView, setUsersView] = useState('list')

  const role = useMemo(() => getRole(), [auth])
  const email = useMemo(() => getEmail(), [auth])

  useEffect(() => {
    if (!auth.token) return
    refreshBooks()
    if (role === 'admin') refreshUsers()
  }, [auth.token, role])

  async function refreshBooks() {
    setLoading(true)
    try {
      const data = await api.listBooks()
      setBooks(data)
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  async function refreshUsers() {
    try {
      const data = await api.listUsers()
      setUsers(data)
    } catch (err) {
      setError(err.message)
    }
  }

  function handleLogout() {
    clearAuth()
    setAuthState({ token: '' })
    setBooks([])
    setUsers([])
    setSelectedBook(null)
  }

  async function handleLogin(e) {
    e.preventDefault()
    setError('')
    const form = new FormData(e.currentTarget)
    const payload = { email: form.get('email'), password: form.get('password') }
    try {
      const data = await api.login(payload)
      setAuth(data.token, data.role, data.email)
      setAuthState({ token: data.token })
    } catch (err) {
      setError(err.message)
    }
  }

  async function handleCreateUser(e) {
    e.preventDefault()
    setError('')
    const form = new FormData(e.currentTarget)
    const payload = {
      email: form.get('email'),
      password: form.get('password'),
      role: form.get('role')
    }
    try {
      await api.createUser(payload)
      e.currentTarget.reset()
      refreshUsers()
    } catch (err) {
      setError(err.message)
    }
  }

  async function handleCreateBook(e) {
    e.preventDefault()
    setError('')
    const form = new FormData(e.currentTarget)
    const payload = {
      title: form.get('title'),
      slug: form.get('slug')
    }
    try {
      await api.createBook(payload)
      e.currentTarget.reset()
      refreshBooks()
    } catch (err) {
      setError(err.message)
    }
  }

  async function handleBuildBook(id) {
    setError('')
    try {
      await api.buildBook(id)
      refreshBooks()
    } catch (err) {
      setError(err.message)
    }
  }

  async function handleUploadBook(id, file) {
    if (!file) return
    setError('')
    try {
      await api.uploadBook(id, file)
      refreshBooks()
    } catch (err) {
      setError(err.message)
    }
  }

  async function handleDeleteBook(id) {
    setError('')
    try {
      await api.deleteBook(id)
      refreshBooks()
    } catch (err) {
      setError(err.message)
    }
  }

  async function handleUserRoleChange(id, role) {
    setError('')
    try {
      await api.updateUser(id, { role })
      refreshUsers()
    } catch (err) {
      setError(err.message)
    }
  }

  async function handleUserActiveToggle(id, active) {
    setError('')
    try {
      await api.updateUser(id, { active })
      refreshUsers()
    } catch (err) {
      setError(err.message)
    }
  }

  async function handleUserDelete(id) {
    setError('')
    try {
      await api.deleteUser(id)
      refreshUsers()
    } catch (err) {
      setError(err.message)
    }
  }

  if (!auth.token) {
    return (
      <div className="page">
        <header className="hero">
          <div>
            <p className="eyebrow">mdBook Portal</p>
            <h1>Secure reading with role-based access.</h1>
            <p className="lead">Admin manages books and users. Readers access authorized content only.</p>
          </div>
        </header>
        <section className="card-grid">
          <div className="card">
            <h2>Login</h2>
            <form onSubmit={handleLogin} className="form">
              <label>
                Email
                <input name="email" type="email" required />
              </label>
              <label>
                Password
                <input name="password" type="password" required />
              </label>
              <button type="submit">Sign in</button>
            </form>
          </div>
        </section>
        {error && <p className="error">{error}</p>}
      </div>
    )
  }

  return (
    <div className="page">
      <header className="topbar">
        <div>
          <p className="eyebrow">Signed in</p>
          <h2>{email}</h2>
          <p className="muted">Role: {role}</p>
        </div>
        <div className="topbar-actions">
          <nav className="nav">
            <button
              className={activeModule === 'books' ? 'nav-active' : 'ghost'}
              onClick={() => setActiveModule('books')}
            >
              Books
            </button>
            {role === 'admin' && (
              <button
                className={activeModule === 'users' ? 'nav-active' : 'ghost'}
                onClick={() => setActiveModule('users')}
              >
                Users
              </button>
            )}
          </nav>
          <button onClick={handleLogout} className="ghost">Logout</button>
        </div>
      </header>

      {error && <p className="error">{error}</p>}

      {activeModule === 'books' && (
        <section className="section">
          <div className="section-header">
            <div>
              <h3>Books</h3>
              {loading && <span className="muted">Loading...</span>}
            </div>
            {role === 'admin' && (
              <nav className="subnav">
                <button
                  className={booksView === 'list' ? 'nav-active' : 'ghost'}
                  onClick={() => setBooksView('list')}
                >
                  List Books
                </button>
                <button
                  className={booksView === 'create' ? 'nav-active' : 'ghost'}
                  onClick={() => setBooksView('create')}
                >
                  Create Book
                </button>
              </nav>
            )}
          </div>
          {booksView === 'list' && (
            <div className="books">
              {books.map((book) => (
                <div key={book.id} className="book-card">
                  <div>
                    <h4>{book.title}</h4>
                    <p className="muted">{book.slug}</p>
                  </div>
                  <div className="book-actions">
                    <button onClick={() => setSelectedBook(book)}>Open</button>
                    {role === 'admin' && (
                      <>
                        <label className="upload">
                          Upload zip
                          <input
                            type="file"
                            accept=".zip"
                            onChange={(e) => handleUploadBook(book.id, e.target.files[0])}
                          />
                        </label>
                        <button className="ghost" onClick={() => handleBuildBook(book.id)}>Build</button>
                        <button className="danger" onClick={() => handleDeleteBook(book.id)}>Delete</button>
                      </>
                    )}
                  </div>
                </div>
              ))}
              {books.length === 0 && <p className="muted">No books yet.</p>}
            </div>
          )}
          {role === 'admin' && booksView === 'create' && (
            <div className="card-grid">
              <div className="card">
                <h3>Create Book</h3>
                <form onSubmit={handleCreateBook} className="form">
                  <label>
                    Title
                    <input name="title" type="text" required />
                  </label>
                  <label>
                    Slug (optional)
                    <input name="slug" type="text" />
                  </label>
                  <button type="submit">Create</button>
                </form>
              </div>
            </div>
          )}
        </section>
      )}

      {selectedBook && (
        <section className="section viewer">
          <div className="section-header">
            <h3>Viewer</h3>
            <button className="ghost" onClick={() => setSelectedBook(null)}>Close</button>
          </div>
          <iframe title="mdbook" src={bookContentUrl(selectedBook.id)} />
        </section>
      )}

      {role === 'admin' && activeModule === 'users' && (
        <section className="section admin">
          <div className="section-header">
            <div>
              <h3>Users</h3>
            </div>
            <nav className="subnav">
              <button
                className={usersView === 'list' ? 'nav-active' : 'ghost'}
                onClick={() => setUsersView('list')}
              >
                List Users
              </button>
              <button
                className={usersView === 'create' ? 'nav-active' : 'ghost'}
                onClick={() => setUsersView('create')}
              >
                Create User
              </button>
            </nav>
          </div>

          {usersView === 'create' && (
            <div className="card-grid">
              <div className="card">
                <h3>Create User</h3>
                <form onSubmit={handleCreateUser} className="form">
                  <label>
                    Email
                    <input name="email" type="email" required />
                  </label>
                  <label>
                    Password
                    <input name="password" type="password" required />
                  </label>
                  <label>
                    Role
                    <select name="role" defaultValue="reader">
                      <option value="reader">reader</option>
                      <option value="admin">admin</option>
                    </select>
                  </label>
                  <button type="submit">Create</button>
                </form>
              </div>
            </div>
          )}

          {usersView === 'list' && (
            <div className="card">
              <h3>Users</h3>
              <div className="users">
                {users.map((user) => (
                  <div key={user.id} className="user-row">
                    <div>
                      <strong>{user.email}</strong>
                      <p className="muted">{user.role} · {user.active ? 'active' : 'inactive'}</p>
                    </div>
                    <div className="user-actions">
                      <select
                        value={user.role}
                        onChange={(e) => handleUserRoleChange(user.id, e.target.value)}
                      >
                        <option value="reader">reader</option>
                        <option value="admin">admin</option>
                      </select>
                      <button
                        className="ghost"
                        onClick={() => handleUserActiveToggle(user.id, !user.active)}
                      >
                        {user.active ? 'Deactivate' : 'Activate'}
                      </button>
                      <button className="danger" onClick={() => handleUserDelete(user.id)}>
                        Delete
                      </button>
                    </div>
                  </div>
                ))}
                {users.length === 0 && <p className="muted">No users found.</p>}
              </div>
            </div>
          )}
        </section>
      )}

    </div>
  )
}

export default App
