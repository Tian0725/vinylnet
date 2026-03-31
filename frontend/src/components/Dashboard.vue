<script setup>
import { ref, onMounted } from 'vue';

// --- Estado de Usuario ---
const nombreUsuario = ref('Cargando...');
const rolUsuario = ref('...');
const iniciales = ref('??');

// --- Estado de la Interfaz (Menús) ---
// Agregamos todas las llaves para que la reactividad funcione en cada sección
const menus = ref({
  ficheros: false,
  transacciones: false,
  inventario_sub: false,
  caja_sub: false,
  clientes_sub: false,
  bancos_sub: false,
  Ventas: false
});

const mobileMenuVisible = ref(false);

onMounted(() => {
  const isLogged = localStorage.getItem('isLogged');
  if (!isLogged) {
    // Si no hay login, redirigir (comentado para que no falle en pruebas locales)
    // window.location.href = '/'; 
    nombreUsuario.value = "Usuario Demo";
    rolUsuario.value = "Administrador";
    iniciales.value = "UD";
    return;
  }

  const nombre = localStorage.getItem('usuario') || 'Usuario';
  const rol = localStorage.getItem('rol') || 'Operador';
  
  nombreUsuario.value = nombre;
  rolUsuario.value = rol;
  iniciales.value = nombre.split(' ').filter(n => n).map(n => n[0]).join('').toUpperCase().substring(0, 2);
});

// --- Funciones de Control ---
const toggleSubmenu = (menuKey) => {
  menus.value[menuKey] = !menus.value[menuKey];
};

const toggleMobileMenu = () => {
  mobileMenuVisible.value = !mobileMenuVisible.value;
};

const logout = () => {
  localStorage.clear();
  window.location.href = '/';
};
</script>

<template>
  <div class="app-container">
    <div v-if="mobileMenuVisible" class="sidebar-overlay" @click="toggleMobileMenu"></div>

    <aside class="sidebar" :class="{ 'sidebar-active': mobileMenuVisible }">
      <div class="brand">
        <div class="brand-flex">
          <span class="logo-dot"></span>
          <h2 class="brand-name">VINYL<span>NET</span></h2>
        </div>
        <button class="close-menu" @click="toggleMobileMenu">✕</button>
      </div>
      
      <nav class="menu">
        <div class="menu-section">Operaciones</div>
        <button class="nav-btn active">
          <span class="btn-content"><span>🏠</span> Inicio</span>
        </button>

        <div class="submenu-container">
          <button class="nav-btn" @click="toggleSubmenu('ficheros')">
            <span class="btn-content"><span>🗂️</span> Ficheros</span>
            <span class="arrow" :class="{ rotated: menus.ficheros }">▼</span>
          </button>
          <div v-show="menus.ficheros" class="submenu-list">
            <button class="nav-sub-btn">Departamentos</button>
            <button class="nav-sub-btn">Depósitos</button>
            <button class="nav-sub-btn">Inventarios</button>
            <button class="nav-sub-btn">Clientes</button>
            <button class="nav-sub-btn">Vendedores</button>
            <button class="nav-sub-btn">Bancos</button>
          </div>
        </div>

        <div class="submenu-container">
          <button class="nav-btn" @click="toggleSubmenu('transacciones')">
            <span class="btn-content"><span>📦</span> Transacciones</span>
            <span class="arrow" :class="{ rotated: menus.transacciones }">▼</span>
          </button>
          
          <div v-show="menus.transacciones" class="submenu-list">
            
            <div class="submenu-container">
              <button class="nav-sub-btn nested-btn" @click="toggleSubmenu('inventario_sub')">
                <span class="btn-content"><span>📋</span> Inventario</span>
                <span class="arrow small" :class="{ rotated: menus.inventario_sub }">▼</span>
              </button>
              <div v-show="menus.inventario_sub" class="sub-submenu-list">
                <button class="nav-sub-btn final-item">➕ Cargos</button>
                <button class="nav-sub-btn final-item">➖ Descargos</button>
                <button class="nav-sub-btn final-item">↔️ Transferencia</button>
              </div>
            </div>

            <div class="submenu-container">
              <button class="nav-sub-btn nested-btn" @click="toggleSubmenu('caja_sub')">
                <span class="btn-content"><span>📠</span> Mov. de Caja</span>
                <span class="arrow small" :class="{ rotated: menus.caja_sub }">▼</span>
              </button>
              <div v-show="menus.caja_sub" class="sub-submenu-list">
                <button class="nav-sub-btn final-item">➕ Préstamos</button>
                <button class="nav-sub-btn final-item">➖ Avances</button>
              </div>
            </div>

            <div class="submenu-container">
              <button class="nav-sub-btn nested-btn" @click="toggleSubmenu('clientes_sub')">
                <span class="btn-content"><span>👤</span> Clientes</span>
                <span class="arrow small" :class="{ rotated: menus.clientes_sub }">▼</span>
              </button>
              <div v-show="menus.clientes_sub" class="sub-submenu-list">
                <button class="nav-sub-btn final-item">💰 Cuentas x Cobrar</button>
              </div>
            </div>

            <div class="submenu-container">
              <button class="nav-sub-btn nested-btn" @click="toggleSubmenu('bancos_sub')">
                <span class="btn-content"><span>🏦</span> Bancos</span>
                <span class="arrow small" :class="{ rotated: menus.bancos_sub }">▼</span>
              </button>
              <div v-show="menus.bancos_sub" class="sub-submenu-list">
                <button class="nav-sub-btn final-item">💸 Movimientos</button>
              </div>
            </div>
          </div>
        </div>

        <div class="submenu-container">
          <button class="nav-btn" @click="toggleSubmenu('Ventas')">
            <span class="btn-content"><span>🛒</span> Ventas</span>
            <span class="arrow" :class="{ rotated: menus.Ventas }">▼</span>
          </button>
          <div v-show="menus.Ventas" class="submenu-list">
            <button class="nav-sub-btn">🧾 Facturas</button>
            <button class="nav-sub-btn">🧾 Facturas Nac.</button>
            <button class="nav-sub-btn">↩️ Devoluciones</button>
            <button class="nav-sub-btn">↩️ Devoluciones Nacionales</button>
            <button class="nav-sub-btn">🚫 Anulaciones</button>
          </div>
        </div>

        <div class="menu-section">Informes</div>
        <button class="nav-btn">
          <span class="btn-content"><span>📊</span> Informes</span>
        </button>

        <div class="menu-section">Sistema</div>
        <button @click="logout" class="nav-btn logout-item">
          <span class="btn-content"><span>🚪</span> Cerrar Sesión</span>
        </button>
      </nav>
    </aside>

    <main class="main-content">
      <header class="header">
        <div class="header-left">
          <button class="hamburger" @click="toggleMobileMenu">☰</button>
          <div class="search-bar">
            <input type="text" placeholder="Buscar..." />
          </div>
        </div>
        <div class="user-profile">
          <div class="avatar">{{ iniciales }}</div>
          <div class="user-meta">
            <span class="name">{{ nombreUsuario }}</span>
            <span class="role">{{ rolUsuario }}</span>
          </div>
        </div>
      </header>

      <section class="dashboard-view">
        <h1 class="view-title">Panel de Control</h1>
        <div class="stats-grid">
          <div class="stat-card">
            <span class="stat-label">Ventas del Día</span>
            <div class="stat-main"><span class="stat-value">$0.00</span></div>
          </div>
          <div class="stat-card">
            <span class="stat-label">Equipos Activos</span>
            <div class="stat-main"><span class="stat-value">0</span></div>
          </div>
        </div>
      </section>
    </main>
  </div>
</template>

<style scoped>
.app-container {
  --primary: #2563eb;
  --bg-main: #f8fafc;
  --bg-card: #ffffff;
  --text-dark: #1e293b;
  --text-light: #64748b;
  --border: #e2e8f0;
  display: flex; height: 100vh; background: var(--bg-main);
  color: var(--text-dark); font-family: 'Inter', sans-serif;
}

/* --- SIDEBAR --- */
.sidebar {
  width: 260px; background: var(--bg-card); border-right: 1px solid var(--border);
  padding: 1.5rem; display: flex; flex-direction: column; transition: all 0.3s ease;
}

.brand { display: flex; align-items: center; justify-content: space-between; margin-bottom: 2rem; }
.brand-flex { display: flex; align-items: center; gap: 10px; }
.logo-dot { width: 12px; height: 12px; background: var(--primary); border-radius: 50%; }
.brand-name { font-size: 1.25rem; font-weight: 800; }
.brand-name span { color: var(--primary); }

.menu-section {
  font-size: 0.7rem; font-weight: 700; text-transform: uppercase;
  color: var(--text-light); margin: 1.5rem 0 0.5rem 0.5rem;
}

/* --- BOTONES Y ALINEACIÓN --- */
.nav-btn {
  width: 100%; text-align: left; padding: 0.75rem 1rem;
  border: none; background: transparent; border-radius: 8px;
  cursor: pointer; font-weight: 500; transition: 0.2s;
  display: flex; justify-content: space-between; align-items: center;
}

.btn-content {
  display: flex; align-items: center; gap: 12px; /* Espacio fijo entre icono y texto */
}

/* Forzamos que el contenedor del emoji tenga un ancho fijo para evitar descuadres */
.btn-content span:first-child {
  display: inline-block; width: 20px; text-align: center;
}

.nav-btn:hover { background: #f1f5f9; }
.nav-btn.active { background: #eff6ff; color: var(--primary); }

/* --- SUBMENÚS --- */
.submenu-list { 
  padding-left: 1rem; margin-bottom: 0.5rem; 
  border-left: 1px solid var(--border); margin-left: 1.4rem;
}

.nav-sub-btn {
  width: 100%; text-align: left; padding: 0.6rem;
  background: none; border: none; font-size: 0.85rem;
  color: var(--text-light); cursor: pointer; border-radius: 4px;
  display: flex; align-items: center; gap: 10px;
}

.nested-btn {
  justify-content: space-between; font-weight: 600; color: var(--text-dark) !important;
}

.sub-submenu-list { padding-left: 0.8rem; margin-top: 4px; display: flex; flex-direction: column; }

.final-item {
  font-size: 0.8rem !important; padding: 0.4rem 0.8rem !important;
}

.arrow { font-size: 0.7rem; transition: transform 0.3s; color: #94a3b8; }
.arrow.small { font-size: 0.5rem; }
.arrow.rotated { transform: rotate(180deg); }

/* --- RESPONSIVO --- */
.hamburger, .close-menu { display: none; background: none; border: none; font-size: 1.5rem; cursor: pointer; }

@media (max-width: 1024px) {
  .sidebar { position: fixed; left: -260px; top: 0; bottom: 0; z-index: 1000; }
  .sidebar-active { left: 0; }
  .sidebar-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.4); z-index: 999; }
  .hamburger, .close-menu { display: block; }
}

/* --- ESTILOS GENERALES --- */
.main-content { flex: 1; overflow-y: auto; }
.header {
  padding: 1rem 2rem; background: var(--bg-card); border-bottom: 1px solid var(--border);
  display: flex; justify-content: space-between; align-items: center;
}
.search-bar input { padding: 0.5rem 1rem; border: 1px solid var(--border); border-radius: 20px; background: #f1f5f9; }
.user-profile { display: flex; align-items: center; gap: 12px; }
.avatar { width: 35px; height: 35px; background: var(--primary); color: white; border-radius: 50%; display: grid; place-items: center; font-weight: bold; font-size: 0.8rem; }
.logout-item { color: #ef4444; margin-top: 1rem; }
.dashboard-view { padding: 2rem; }
.stats-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 1.5rem; }
.stat-card { background: white; padding: 1.5rem; border-radius: 12px; border: 1px solid var(--border); }
.stat-value { font-size: 1.5rem; font-weight: 700; }
.stat-label { color: var(--text-light); font-size: 0.85rem; }
</style>