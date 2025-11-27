import { useState, useEffect } from 'react'

// --- TEMA ---
const theme = {
  bg: '#121212',
  cardBg: '#1e1e1e',
  textMain: '#e0e0e0',
  textMuted: '#777777ff',
  accent: '#bb86fc',
  success: '#03dac6', 
  danger: '#cf6679',
  border: '#4f4f4fff'
}

// --- PANTALLA DE LOGIN ---
const LoginScreen = ({ onLogin, loading, error }) => {
  const [uid, setUid] = useState('')
  const [pass, setPass] = useState('')

  const handleSubmit = (e) => {
    e.preventDefault();
    if(!uid || !pass) return;
    onLogin(uid, pass);
  }

  // Estilos específicos para centrar
  const loginStyles = {
    overlay: {
      position: 'fixed',     
      top: 0,
      left: 0,
      width: '100vw',        
      height: '100vh',       
      backgroundColor: theme.bg,
      display: 'flex',       
      alignItems: 'center',  
      justifyContent: 'center', 
      zIndex: 9999            
    },
    card: {
      maxWidth: '450px',
      width: '90%',          
      backgroundColor: theme.cardBg,
      border: `1px solid ${theme.border}`,
      boxShadow: '0 0 50px rgba(0,0,0,0.8)'
    }
  }

  return (
    <div style={loginStyles.overlay}>
      <div className="card p-5" style={loginStyles.card}>
        
        <div className="text-center mb-5">
          <div className="display-1 mb-3">👤</div>
          <h2 className="fw-bold tracking-wider" style={{ color: theme.accent, letterSpacing: '2px' }}>LOGIN</h2>
          <small style={{ color: theme.textMuted }}>ACCESO AL SISTEMA DE VENTAS</small>
        </div>

        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label className="form-label small fw-bold" style={{ color: theme.textMuted }}>IDENTIFICADOR (ID)</label>
            <input 
              type="text" 
              className="form-control form-control-lg bg-dark text-white border-secondary" 
              placeholder="usuario"
              value={uid}
              onChange={(e) => setUid(e.target.value)}
              autoFocus
            />
          </div>
          <div className="mb-4">
            <label className="form-label small fw-bold" style={{ color: theme.textMuted }}>CLAVE DE ACCESO</label>
            <input 
              type="password" 
              className="form-control form-control-lg bg-dark text-white border-secondary" 
              placeholder="••••••••"
              value={pass}
              onChange={(e) => setPass(e.target.value)}
            />
          </div>

          {error && (
            <div className="alert alert-dismissible fade show d-flex align-items-center" role="alert" 
                 style={{ backgroundColor: 'rgba(207, 102, 121, 0.2)', color: theme.danger, border: `1px solid ${theme.danger}` }}>
              ⚠️ {error}
            </div>
          )}

          <button type="submit" className="btn w-100 fw-bold py-3 mt-2" 
                  style={{ backgroundColor: theme.accent, color: '#000', fontSize: '1.1rem', boxShadow: `0 0 15px ${theme.accent}40` }} 
                  disabled={loading}>
            {loading ? 'AUTENTICANDO...' : 'INICIAR SESIÓN ➜'}
          </button>
        </form>
      
      </div>
    </div>
  )
}

// --- PANEL DE ADMINISTRADOR ---
const AdminPanel = ({ user, editingProduct, onCancelEdit, onProductSaved }) => {
  const [name, setName] = useState('')
  const [price, setPrice] = useState('')
  const [stock, setStock] = useState('')
  
  useEffect(() => {
    if (editingProduct) {
      setName(editingProduct.name)
      setPrice(editingProduct.price)
      setStock(editingProduct.stock)
    } else {
      setName(''); setPrice(''); setStock('')
    }
  }, [editingProduct])

  const handleSubmit = async (e) => {
    e.preventDefault()
    const endpoint = editingProduct ? '/api/admin/update-product' : '/api/admin/create-product'
    
    const body = {
      user_id: user.id,
      name, 
      price: parseFloat(price), 
      stock: parseInt(stock)
    }
    if (editingProduct) body.id = editingProduct.id

    try {
      await fetch(endpoint, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify(body)
      })
      onProductSaved()
    } catch(err) { alert('Error guardando') }
  }

  return (
    <div className="card border-0 shadow-lg sticky-top" style={{ top: '100px', backgroundColor: '#1a1a1a', borderTop: `4px solid ${theme.accent}` }}>
      <div className="card-header bg-transparent border-bottom border-secondary py-3">
        <div className="d-flex justify-content-between align-items-center">
            <h6 className="m-0 fw-bold" style={{color: theme.accent}}>
            {editingProduct ? `EDITANDO ID: ${editingProduct.id}` : 'NUEVO PRODUCTO'}
            </h6>
            {editingProduct && <button onClick={onCancelEdit} className="btn btn-sm btn-outline-secondary" style={{fontSize: '0.7rem'}}>CANCELAR</button>}
        </div>
      </div>
      
      <div className="card-body">
        <form onSubmit={handleSubmit} className="row g-3">
          <div className="col-12">
            <label className="form-label small" style={{ color: theme.textMuted }}>NOMBRE</label>
            <input type="text" className="form-control bg-dark text-white border-secondary" 
                   placeholder="Ej. RTX 4090" value={name} onChange={e=>setName(e.target.value)} required />
          </div>
          <div className="col-12">
            <label className="form-label small" style={{ color: theme.textMuted }}>PRECIO</label>
            <input type="number" className="form-control bg-dark text-white border-secondary" 
                   placeholder="0.00" value={price} onChange={e=>setPrice(e.target.value)} required />
          </div>
          <div className="col-12">
            <label className="form-label small" style={{ color: theme.textMuted }}>STOCK</label>
            <input type="number" className="form-control bg-dark text-white border-secondary" 
                   placeholder="0" value={stock} onChange={e=>setStock(e.target.value)} required />
          </div>
          <div className="col-12 mt-4">
            <button className={`btn w-100 fw-bold py-2 ${editingProduct ? 'btn-info' : 'btn-primary'}`} 
               style={{backgroundColor: theme.accent, border:'none', color:'#000'}}>
               {editingProduct ? 'GUARDAR CAMBIOS' : 'CREAR PRODUCTO'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

// --- APLICACIÓN PRINCIPAL ---
function App() {
  const [user, setUser] = useState(null)
  const [loginError, setLoginError] = useState(null)

  const [editingProduct, setEditingProduct] = useState(null)
  
  const [products, setProducts] = useState([])
  const [cart, setCart] = useState([])
  const [orders, setOrders] = useState([])
  const [alert, setAlert] = useState(null)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (alert) {
      const timer = setTimeout(() => setAlert(null), 4000);
      return () => clearTimeout(timer);
    }
  }, [alert]);

  useEffect(() => {
    if (user) {
      fetchProducts().then((loadedProducts) => {
         fetchOrders(loadedProducts);
      });
    }
  }, [user])

  const handleLogin = async (uid, pass) => {
    setLoading(true)
    setLoginError(null)
    try {
      const res = await fetch('/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ user_id: uid, password: pass })
      })
      const data = await res.json()
      
      if (data.success) {
        setUser(data.user)
      } else {
        setLoginError(data.error || 'Credenciales inválidas')
      }
    } catch (err) {
      setLoginError('No se pudo conectar con el Gateway')
    } finally {
      setLoading(false)
    }
  }

  const handleLogout = () => {
    setUser(null); setCart([]); setOrders([]);
  }

  const deleteProduct = async (id) => {
    if(!confirm('¿Estás seguro de eliminar este producto?')) return;
    await fetch('/api/admin/delete-product', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({ user_id: user.id, id: id })
    })
    fetchProducts() // Recargar
}

  const fetchProducts = async () => {
    try {
      const res = await fetch('/api/products')
      if (!res.ok) throw new Error('Error')
      const data = await res.json()
      
      const sortedData = data.sort((a, b) => a.id - b.id);

      setProducts(sortedData.map(p => ({ ...p, currentStock: p.stock })))

      return sortedData;
    } catch (error) {
      setAlert({ type: 'danger', msg: 'Error de conexión con microservicios' })
    }
  }

  const fetchOrders = async (currentProducts) => {
    if (!user) return;
    try {
      const catalog = currentProducts || products;

      const res = await fetch(`/api/orders?user_id=${user.id}`)
      const data = await res.json()
      
      if (!Array.isArray(data)) return;

      const formattedOrders = data.map(o => {
        const dateStr = new Date(Number(o.date) * 1000).toLocaleString()
        
        const itemsWithNames = o.items.map(item => {
            const productInfo = catalog.find(p => p.id === item.id) 
            return {
                ...item,
                name: productInfo ? productInfo.name : `Producto #${item.id}`
            }
        })

        return {
            id: o.id,
            total: o.total,
            date: dateStr,
            items: itemsWithNames
        }
      })
      
      setOrders(formattedOrders.reverse())

    } catch (error) {
      console.error("Error cargando historial:", error)
    }
  }

  const addToCart = (product) => {
    if (product.currentStock <= 0) return;
    setCart(prev => {
      const exists = prev.find(item => item.id === product.id);
      if (exists) return prev.map(item => item.id === product.id ? { ...item, quantity: item.quantity + 1 } : item)
      return [...prev, { ...product, quantity: 1 }]
    })
    setProducts(prev => prev.map(p => p.id === product.id ? { ...p, currentStock: p.currentStock - 1 } : p))
  }

  const clearCart = () => {
    setCart([]); fetchProducts();
  }

  const handleCheckout = async () => {
    if (cart.length === 0) return;
    setLoading(true)
    try {
      const res = await fetch('/api/buy', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          user_id: user.id,
          items: cart.map(i => ({ product_id: i.id, quantity: i.quantity }))
        })
      })
      const data = await res.json()

      if (data.success) {
        setAlert({ type: 'success', msg: `¡Compra Exitosa! Orden #${data.order.id}` })
        setOrders(prev => [{ id: data.order.id, total: data.order.total, date: new Date().toLocaleTimeString(), items: [...cart] }, ...prev])
        setCart([]); fetchProducts();
      } else {
        setAlert({ type: 'danger', msg: `Error: ${data.error}` })
        fetchProducts()
      }
    } catch (error) {
      setAlert({ type: 'danger', msg: 'Fallo crítico en la red' })
    } finally {
      setLoading(false)
    }
  }

  const cartTotal = cart.reduce((acc, item) => acc + (item.price * item.quantity), 0);
  
  if (!user) return <LoginScreen onLogin={handleLogin} loading={loading} error={loginError} />

  return (
    <div style={{ backgroundColor: theme.bg, minHeight: '100vh', color: theme.textMain, fontFamily: "'Segoe UI', sans-serif" }}>
      
      {/* NAVBAR */}
      <nav className="navbar navbar-dark sticky-top shadow-sm px-4 py-3" style={{ backgroundColor: '#000', borderBottom: `1px solid ${theme.border}` }}>
        <div className="container-fluid">
          <div className="d-flex align-items-center">
            <span className="fs-2 me-2"></span>
            <div>
              <span className="navbar-brand fw-bold fs-4 m-0" style={{ color: theme.accent, letterSpacing: '1px' }}>MARKETPLACE</span>
              <div style={{ fontSize: '0.75rem', color: theme.textMuted }}>MICROSERVICES ARCHITECTURE</div>
            </div>
          </div>
          <div className="d-flex align-items-center bg-dark rounded-pill px-3 py-2 border border-secondary">
            <div className="d-flex flex-column text-end me-3" style={{ lineHeight: '1.2' }}>
              <span className="fw-bold text-white">{user.name.toUpperCase()}</span>
              <small style={{ color: theme.accent, fontSize: '0.7rem' }}>{user.id} • {user.email}</small>
            </div>
            <button onClick={handleLogout} className="btn btn-sm btn-outline-danger rounded-pill px-3">SALIR</button>
          </div>
        </div>
      </nav>

      {/* ALERTA */}
      {alert && (
        <div className="fixed-top d-flex justify-content-center" style={{ top: '90px', zIndex: 9999 }}>
          <div className={`shadow-lg px-4 py-3 rounded-3 d-flex align-items-center border border-${alert.type === 'success' ? 'success' : 'danger'}`} 
               style={{ backgroundColor: '#222', minWidth: '300px' }}>
            <span className="fs-4 me-3">{alert.type === 'success' ? '' : '❌'}</span>
            <div>
              <h6 className="m-0 fw-bold text-white">{alert.type === 'success' ? 'Operación Exitosa' : 'Error'}</h6>
              <small className="text-muted">{alert.msg}</small>
            </div>
          </div>
        </div>
      )}

      <div className="container-fluid px-4 py-4">
        <div className="row g-4">
          
          {/* CATÁLOGO */}
          <div className="col-lg-9">
            <div className="d-flex justify-content-between align-items-end mb-4 border-bottom border-secondary pb-2">
              <h4 className="m-0 fw-light text-white">PRODUCTOS <strong style={{ color: theme.accent }}>DISPONIBLES</strong></h4>
            </div>

            <div className="row row-cols-1 row-cols-md-2 row-cols-xl-3 row-cols-xxl-4 g-4 mb-5">
              {products.length === 0 ? (
                <div className="col-12 text-center py-5 text-muted">
                  <div className="spinner-border text-primary mb-3" role="status"></div>
                  <p>Sincronizando con Microservicio de Productos...</p>
                </div>
              ) : (
                products.map((p) => (
                  <div key={p.id} className="col">
                    <div className="card h-100 shadow-sm border-0" 
                         style={{ backgroundColor: theme.cardBg, transition: 'transform 0.2s' }}
                         onMouseEnter={e => e.currentTarget.style.transform = 'translateY(-5px)'}
                         onMouseLeave={e => e.currentTarget.style.transform = 'translateY(0)'}
                    >
                      <div className="card-body d-flex flex-column p-4">
                        <div className="d-flex justify-content-between mb-3">
                           <span className="badge bg-black border border-secondary text-secondary font-monospace">{p.id}</span>
                           <span className={`badge ${p.currentStock > 0 ? '' : 'bg-danger text-white'}`} 
                                 style={p.currentStock > 0 ? { color: theme.success, border: `1px solid ${theme.success}`, backgroundColor: 'transparent' } : {}}>
                             STOCK: {p.currentStock}
                           </span>
                        </div>
                        <h5 className="card-title fw-bold text-white mb-1">{p.name}</h5>
                        <small style={{ color: theme.textMuted }}>Hardware Component</small>
                        <div className="mt-auto pt-4">
                          <div className="d-flex justify-content-between align-items-center mb-3">
                              <span className="fs-3 fw-light text-white">${p.price}</span>
                          </div>

                          {/* LÓGICA DE BOTONES */}
                          {user.id === 'u1' ? (
                              <div className="d-flex gap-2">
                                  <button className="btn w-50 btn-outline-info fw-bold" onClick={() => setEditingProduct(p)}>
                                      EDITAR
                                  </button>
                                  <button className="btn w-50 btn-outline-danger fw-bold" onClick={() => deleteProduct(p.id)}>
                                      BORRAR
                                  </button>
                              </div>
                          ) : (
                              <button className="btn w-100 py-2 fw-bold" onClick={() => addToCart(p)} disabled={p.currentStock < 1}
                              style={{ 
                                  backgroundColor: p.currentStock > 0 ? 'transparent' : '#333', 
                                  color: p.currentStock > 0 ? theme.accent : '#666',
                                  border: `2px solid ${p.currentStock > 0 ? theme.accent : '#333'}` 
                              }}>
                              {p.currentStock > 0 ? 'AGREGAR AL CARRITO' : 'AGOTADO'}
                              </button>
                          )}
                      </div>
                      </div>
                    </div>
                  </div>
                ))
              )}
            </div>

            {orders.length > 0 && (
              <div className="mt-5 p-4 rounded-3" style={{ backgroundColor: '#181818', border: `1px solid ${theme.border}` }}>
                <h5 className="fw-bold mb-4" style={{ color: theme.textMain }}>HISTORIAL DE ÓRDENES</h5>
                <div className="table-responsive">
                  <table className="table table-dark table-hover align-middle mb-0" style={{ backgroundColor: 'transparent' }}>
                    <thead>
                      <tr style={{ color: theme.textMuted, fontSize: '0.85rem', textTransform: 'uppercase' }}>
                        <th>ID Orden</th>
                        <th>Hora</th>
                        <th>Detalle</th>
                        <th className="text-end">Total</th>
                        <th className="text-end">Estado</th>
                      </tr>
                    </thead>
                    <tbody>
                      {orders.map((o, i) => (
                        <tr key={i} style={{ borderBottom: '1px solid #333' }}>
                          <td className="font-monospace" style={{ color: theme.accent }}>{o.id}</td>
                          <td className="text-white">{o.date}</td>
                          <td>
                            {o.items.map(it => (
                              <span key={it.id} className="badge bg-transparent border border-info text-info me-1">
                                {it.quantity}x {it.name}
                              </span>
                            ))}
                          </td>
                          <td className="text-end fw-bold fs-5 text-white">${o.total}</td>
                          <td className="text-end"><span className="badge bg-success text-dark">COMPLETADO</span></td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            )}
          </div>

          {/* CARRITO */}
          <div className="col-lg-3">
            {user.is_admin ? (
                <AdminPanel 
                    user={user} 
                    editingProduct={editingProduct} 
                    onCancelEdit={() => setEditingProduct(null)}
                    onProductSaved={() => { setEditingProduct(null); fetchProducts(); }} 
                />
            ) : (
            <div className="sticky-top" style={{ top: '100px', maxHeight: '85vh', overflowY: 'auto' }}>
              <div className="card border-0 shadow-lg" style={{ backgroundColor: theme.cardBg, borderTop: `4px solid ${theme.accent}` }}>
                <div className="card-header bg-transparent border-bottom border-secondary py-3">
                  <h5 className="m-0 fw-bold text-white d-flex align-items-center">
                    <span className="me-2"></span> CARRITO
                  </h5>
                </div>
                <div className="card-body p-0">
                  {cart.length === 0 ? (
                    <div className="p-5 text-center">
                      <p className="m-0" style={{ color: theme.textMuted }}>No has seleccionado items.</p>
                    </div>
                  ) : (
                    <ul className="list-group list-group-flush">
                      {cart.map((item, idx) => (
                        <li key={idx} className="list-group-item bg-transparent border-secondary d-flex justify-content-between align-items-center py-3">
                          <div>
                            <div className="fw-bold text-white">{item.name}</div>
                            <small style={{ color: theme.accent }}>${item.price} x {item.quantity}</small>
                          </div>
                          <span className="fw-bold fs-5 text-white">
                            ${(item.price * item.quantity).toFixed(2)}
                          </span>
                        </li>
                      ))}
                    </ul>
                  )}
                </div>
                <div className="card-footer bg-dark border-top border-secondary p-4">
                  <div className="d-flex justify-content-between align-items-end mb-4">
                    <span style={{ color: theme.textMuted }}>TOTAL A PAGAR</span>
                    <span className="display-6 fw-bold text-white">${cartTotal.toFixed(2)}</span>
                  </div>
                  <button 
                    className="btn w-100 py-3 fw-bold fs-5 shadow" 
                    onClick={handleCheckout}
                    disabled={cart.length === 0 || loading}
                    style={{ 
                      backgroundColor: cart.length > 0 ? theme.success : '#333',
                      color: cart.length > 0 ? '#000' : '#555',
                      border: 'none',
                      transition: 'all 0.3s'
                    }}
                  >
                    {loading ? <span>PROCESANDO...</span> : 'CONFIRMAR PAGO'}
                  </button>
                  {cart.length > 0 && (
                    <button onClick={clearCart} className="btn w-100 mt-3 text-danger text-decoration-none btn-sm">
                      VACIAR CARRITO
                    </button>
                  )}
                </div>
              </div>
            </div>
            )}
          </div>

        </div>
      </div>
    </div>
  )
}

export default App