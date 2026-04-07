<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';

// --- Estado de Usuario y Sesión ---
const nombreUsuario = ref('Cargando...');
const rolUsuario = ref('...');
const iniciales = ref('??');
const vistaActual = ref('dashboard');

// --- Estado de la Interfaz (Menús Desplegables) ---
const menus = ref({
  ficheros: false,
  transacciones: false,
  inventario_sub: false,
  caja_sub: false,
  clientes_sub: false,
  bancos_sub: false,
  Ventas: false,
  Configuracion: false
});

const mobileMenuVisible = ref(false);

// --- Lógica de CRUD de Usuarios ---
const usuarios = ref([]);
const roles = ref([]);
const cargando = ref(false);
const modalAbierto = ref(false);
const rolModalAbierto = ref(false);
const rolEditando = ref(false);
const rolForm = ref({ id: null, nombre: '', descripcion: '' });
const API_URL = 'http://localhost:8080/api';

// Formulario para Nuevo/Editar Usuario
const form = ref({
  id: null,
  nombre_completo: '',
  username: '',
  password: '',
  rol_id: 1
});

// 1. Obtener Usuarios (Read)
const fetchUsuarios = async () => {
  cargando.value = true;
  try {
    const res = await axios.get(`${API_URL}/usuarios`);
    usuarios.value = res.data || [];
  } catch (err) {
    console.error("Error al obtener usuarios:", err);
  } finally {
    cargando.value = false;
  }
};

// 2. Guardar Usuario (Create)
const guardarUsuario = async () => {
  try {
    await axios.post(`${API_URL}/usuarios`, form.value);
    modalAbierto.value = false;
    resetForm();
    fetchUsuarios();
  } catch (err) {
    alert("Error al guardar el usuario. Revisa la consola.");
    console.error(err);
  }
};

// 3. Eliminar Usuario (Delete)
const eliminarUsuario = async (id) => {
  if (confirm("¿Estás seguro de eliminar este usuario?")) {
    try {
      await axios.delete(`${API_URL}/usuarios/${id}`);
      fetchUsuarios();
    } catch (err) {
      alert("No se pudo eliminar el usuario.");
    }
  }
};


const resetForm = () => {
  form.value = { id: null, nombre_completo: '', username: '', password: '', rol_id: 1 };
};

const resetRolForm = () => {
  rolForm.value = { id: null, nombre: '', descripcion: '' };
  rolEditando.value = false;
};

const abrirNuevoRol = () => {
  resetRolForm();
  rolModalAbierto.value = true;
};

const editarRol = (rol) => {
  rolForm.value = { id: rol.id, nombre: rol.nombre, descripcion: rol.descripcion || '' };
  rolEditando.value = true;
  rolModalAbierto.value = true;
};

const guardarRol = async () => {
  try {
    if (!rolForm.value.nombre) {
      alert('El nombre del rol es obligatorio.');
      return;
    }

    if (rolForm.value.id) {
      await axios.put(`${API_URL}/roles/${rolForm.value.id}`, rolForm.value);
    } else {
      await axios.post(`${API_URL}/roles`, rolForm.value);
    }

    rolModalAbierto.value = false;
    resetRolForm();
    fetchRoles();
  } catch (err) {
    alert('Error al guardar el rol. Revisa la consola.');
    console.error(err);
  }
};

const eliminarRol = async (id) => {
  if (confirm('¿Estás seguro de eliminar este rol?')) {
    try {
      await axios.delete(`${API_URL}/roles/${id}`);
      fetchRoles();
    } catch (err) {
      alert('No se pudo eliminar el rol.');
      console.error(err);
    }
  }
};

const cerrarRolModal = () => {
  rolModalAbierto.value = false;
  resetRolForm();
};

// Obtener Roles de la DB
const fetchRoles = async () => {
  try {
    const res = await axios.get(`${API_URL}/roles`);
    roles.value = res.data || [];
    
    // Opcional: Si hay roles, poner el primero por defecto en el form
    if (roles.value.length > 0 && !form.value.rol_id) {
      form.value.rol_id = roles.value[0].id;
    }
  } catch (err) {
    console.error("Error al obtener roles:", err);
  }
};
// --- Inicialización ---
onMounted(() => {
  const nombre = localStorage.getItem('usuario') || 'Jesus Sarmiento';
  const rol = localStorage.getItem('rol') || 'Administrador';
  
  nombreUsuario.value = nombre;
  rolUsuario.value = rol;
  iniciales.value = nombre.split(' ').filter(n => n).map(n => n[0]).join('').toUpperCase().substring(0, 2);

  fetchUsuarios();
  fetchRoles();
});

// --- Funciones de Navegación ---
const toggleSubmenu = (menuKey) => {
  menus.value[menuKey] = !menus.value[menuKey];
};

const navegar = (vista) => {
  vistaActual.value = vista;
  if (window.innerWidth < 1024) mobileMenuVisible.value = false;
};

const logout = () => {
  localStorage.clear();
  window.location.href = '/';
};



</script>

<template>
  <div class="flex h-screen bg-brand-navy-dark font-sans text-neutral-200 overflow-hidden">
    
    <div v-if="mobileMenuVisible" class="fixed inset-0 bg-black/60 backdrop-blur-sm z-[60] lg:hidden" @click="mobileMenuVisible = false"></div>

    <aside 
      class="fixed lg:static inset-y-0 left-0 w-72 bg-brand-navy-base/80 backdrop-blur-xl border-r border-white/5 z-[70] transform transition-transform duration-300 lg:translate-x-0"
      :class="mobileMenuVisible ? 'translate-x-0' : '-translate-x-full'"
    >
      <div class="h-20 flex items-center justify-between px-8 border-b border-white/5">
        <div class="flex items-center space-x-3">
          <div class="w-2.5 h-2.5 bg-brand-royal rounded-full shadow-[0_0_10px_rgba(37,99,235,0.5)]"></div>
          <h2 class="text-xl font-black text-white uppercase italic tracking-tighter">VINYL<span class="text-brand-royal">NET</span></h2>
        </div>
      </div>
      
      <nav class="p-6 h-[calc(100vh-5rem)] overflow-y-auto space-y-2 scrollbar-none">
        <div class="px-2 pb-2 text-[10px] font-bold uppercase tracking-[0.2em] text-neutral-600">Operaciones</div>
        
        <button @click="navegar('dashboard')" class="w-full flex items-center px-4 py-3 hover:bg-white/5 rounded-2xl transition-all group">
          <span class="mr-3 text-lg group-hover:scale-110 transition-transform">🏠</span>
          <span class="font-bold text-sm">Inicio</span>
        </button>

        <div class="space-y-1">
          <button @click="toggleSubmenu('ficheros')" class="w-full flex items-center justify-between px-4 py-3 hover:bg-white/5 rounded-2xl transition-all group">
            <div class="flex items-center">
              <span class="mr-3 text-lg text-neutral-400">🗂️</span>
              <span class="font-bold text-sm text-neutral-300">Ficheros</span>
            </div>
            <svg class="w-4 h-4 transition-transform" :class="{ 'rotate-180': menus.ficheros }" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path d="M19 9l-7 7-7-7"/></svg>
          </button>
          <div v-show="menus.ficheros" class="ml-9 border-l border-white/5 pl-4 space-y-1">
            <button v-for="item in ['Departamentos', 'Depósitos', 'Inventarios', 'Clientes', 'Vendedores', 'Bancos']" :key="item" class="w-full text-left px-3 py-2 text-xs font-semibold text-neutral-500 hover:text-white transition-colors">{{ item }}</button>
          </div>
        </div>

        <div class="space-y-1">
          <button @click="toggleSubmenu('transacciones')" class="w-full flex items-center justify-between px-4 py-3 hover:bg-white/5 rounded-2xl transition-all group">
            <div class="flex items-center">
              <span class="mr-3 text-lg text-neutral-400">📦</span>
              <span class="font-bold text-sm text-neutral-300">Transacciones</span>
            </div>
            <svg class="w-4 h-4 transition-transform" :class="{ 'rotate-180': menus.transacciones }" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path d="M19 9l-7 7-7-7"/></svg>
          </button>
          <div v-show="menus.transacciones" class="ml-9 border-l border-white/5 pl-4 space-y-3 pt-2">
            <div v-for="(sub, subKey) in {
              'Inventario': { key: 'inventario_sub', icon: '📋', items: ['➕ Cargos', '➖ Descargos', '↔️ Transferencia'] },
              'Mov. de Caja': { key: 'caja_sub', icon: '📠', items: ['➕ Préstamos', '➖ Avances'] },
              'Clientes': { key: 'clientes_sub', icon: '👤', items: ['💰 Cuentas x Cobrar'] },
              'Bancos': { key: 'bancos_sub', icon: '🏦', items: ['💸 Movimientos'] }
            }" :key="subKey">
              <button @click="toggleSubmenu(sub.key)" class="w-full flex items-center justify-between text-xs font-bold text-neutral-400 hover:text-white">
                <span class="flex items-center"><span class="mr-2">{{ sub.icon }}</span> {{ subKey }}</span>
                <svg class="w-3 h-3 transition-transform" :class="{ 'rotate-180': menus[sub.key] }" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path d="M19 9l-7 7-7-7"/></svg>
              </button>
              <div v-show="menus[sub.key]" class="ml-4 border-l border-white/5 pl-3 space-y-1 mt-1">
                <button v-for="item in sub.items" :key="item" class="w-full text-left py-1.5 text-[10px] text-neutral-600 hover:text-brand-royal transition-colors">{{ item }}</button>
              </div>
            </div>
          </div>
        </div>

        <div class="space-y-1">
          <button @click="toggleSubmenu('Ventas')" class="w-full flex items-center justify-between px-4 py-3 hover:bg-white/5 rounded-2xl transition-all group">
            <div class="flex items-center">
              <span class="mr-3 text-lg text-neutral-400">🛒</span>
              <span class="font-bold text-sm text-neutral-300">Ventas</span>
            </div>
            <svg class="w-4 h-4 transition-transform" :class="{ 'rotate-180': menus.Ventas }" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path d="M19 9l-7 7-7-7"/></svg>
          </button>
          <div v-show="menus.Ventas" class="ml-9 border-l border-white/5 pl-4 space-y-1">
            <button v-for="item in ['🧾 Facturas', '🧾 Facturas Nac.', '↩️ Devoluciones', '↩️ Devoluciones Nac.', '🚫 Anulaciones']" :key="item" class="w-full text-left px-3 py-2 text-xs font-semibold text-neutral-500 hover:text-white transition-colors">{{ item }}</button>
          </div>
        </div>

        <div class="pt-6 px-2 pb-2 text-[10px] font-bold uppercase tracking-[0.2em] text-neutral-600">Sistema</div>
        <div class="space-y-1">
          <button @click="toggleSubmenu('Configuracion')" class="w-full flex items-center justify-between px-4 py-3 hover:bg-white/5 rounded-2xl transition-all group">
            <div class="flex items-center">
              <span class="mr-3 text-lg text-neutral-400">⚙️</span>
              <span class="font-bold text-sm text-neutral-300">Configuración</span>
            </div>
            <svg class="w-4 h-4 transition-transform" :class="{ 'rotate-180': menus.Configuracion }" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path d="M19 9l-7 7-7-7"/></svg>
          </button>
          <div v-show="menus.Configuracion" class="ml-9 border-l border-white/5 pl-4 space-y-1">
            <button @click="navegar('usuarios')" class="w-full text-left px-3 py-2 text-xs font-semibold transition-colors" :class="vistaActual === 'usuarios' ? 'text-brand-royal' : 'text-neutral-500 hover:text-white'">Usuarios</button>
            <button
              v-for="item in ['Roles', 'Permisos', 'Empresas', 'Sucursales']"
              :key="item"
              @click="item === 'Roles' ? navegar('roles') : null"
              class="w-full text-left px-3 py-2 text-xs font-semibold transition-colors"
              :class="item === 'Roles' ? (vistaActual === 'roles' ? 'text-brand-royal' : 'text-neutral-500 hover:text-white') : 'text-neutral-500 hover:text-white'"
            >
              {{ item }}
            </button>
          </div>
        </div>

        <button @click="logout" class="w-full flex items-center px-4 py-3 text-red-500 hover:bg-red-500/10 rounded-2xl transition-all mt-10">
          <span class="mr-3 text-lg">🚪</span>
          <span class="font-bold text-sm">Cerrar Sesión</span>
        </button>
      </nav>
    </aside>

    <div class="flex-1 flex flex-col h-screen overflow-hidden">
      <header class="h-20 bg-brand-navy-base/40 backdrop-blur-md border-b border-white/5 flex items-center justify-between px-6 lg:px-10 flex-shrink-0">
        <button class="lg:hidden p-2 text-neutral-400" @click="mobileMenuVisible = true">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path d="M4 6h16M4 12h16m-7 6h7"/></svg>
        </button>

        <div class="flex items-center space-x-4 ml-auto">
          <div class="text-right hidden sm:block">
            <div class="text-xs font-black text-white uppercase tracking-tighter">{{ nombreUsuario }}</div>
            <div class="text-[10px] font-bold text-brand-royal uppercase tracking-widest">{{ rolUsuario }}</div>
          </div>
          <div class="w-10 h-10 rounded-2xl bg-gradient-to-tr from-brand-royal to-indigo-600 flex items-center justify-center text-xs font-black text-white shadow-lg">
            {{ iniciales }}
          </div>
        </div>
      </header>

      <main class="flex-1 overflow-y-auto p-6 lg:p-10 bg-[radial-gradient(ellipse_at_bottom_left,_var(--tw-gradient-stops))] from-brand-navy-base/20 via-brand-navy-dark to-brand-navy-dark">
        <div class="max-w-7xl mx-auto">
          
          <div v-if="vistaActual === 'dashboard'" class="space-y-10">
            <div class="flex flex-col space-y-1">
              <h1 class="text-3xl font-black text-white tracking-tighter uppercase italic">Panel de Control</h1>
              <p class="text-sm font-medium text-neutral-500 uppercase tracking-widest">Resumen ejecutivo del sistema</p>
            </div>
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
               <div v-for="stat in [{ label: 'Usuarios en DB', value: usuarios.length, icon: '👤', color: 'from-brand-royal/20 to-brand-royal/5', border: 'border-brand-royal/20' }]" :key="stat.label" class="p-6 rounded-3xl bg-neutral-900/50 border border-white/5">
                  <div class="text-2xl mb-2">{{ stat.icon }}</div>
                  <div class="text-xs font-bold text-neutral-500 uppercase">{{ stat.label }}</div>
                  <div class="text-2xl font-black text-white">{{ stat.value }}</div>
               </div>
            </div>
          </div>

          <div v-if="vistaActual === 'usuarios'" class="space-y-6">
            <div class="flex justify-between items-center">
              <div>
                <h2 class="text-2xl font-black text-white italic uppercase">Gestión de Usuarios</h2>
                <p class="text-xs text-neutral-500 font-bold uppercase tracking-widest">Administración de accesos a Vinylnet</p>
              </div>
              <button @click="modalAbierto = true" class="bg-brand-royal hover:bg-brand-royal/80 text-white px-6 py-2 rounded-xl font-bold text-xs transition-all shadow-lg shadow-brand-royal/20">
                + NUEVO USUARIO
              </button>
            </div>

            <div class="bg-white/5 border border-white/10 rounded-3xl overflow-hidden backdrop-blur-md">
              <table class="w-full text-left border-collapse">
                <thead class="bg-white/5 text-[10px] uppercase text-neutral-500 font-bold tracking-widest">
                  <tr>
                    <th class="p-5">Nombre</th>
                    <th class="p-5">Usuario</th>
                    <th class="p-5">Rol</th>
                    <th class="p-5 text-right">Acciones</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-white/5">
                  <tr v-for="user in usuarios" :key="user.id" class="hover:bg-white/5 transition-colors">
                    <td class="p-5 font-bold text-white text-sm">{{ user.nombre_completo || user.nombre_complet }}</td>
                    <td class="p-5 text-brand-royal font-mono text-xs">@{{ user.username }}</td>
                    <td class="p-5"><span class="px-3 py-1 bg-brand-royal/10 text-brand-royal text-[10px] font-bold rounded-full">{{ user.rol }}</span></td>
                    <td class="p-5 text-right">
                      <button @click="eliminarUsuario(user.id)" class="text-red-400 hover:text-red-300 text-xs font-bold uppercase transition-colors">Eliminar</button>
                    </td>
                  </tr>
                </tbody>
              </table>
              <div v-if="cargando" class="p-10 text-center text-brand-royal animate-pulse uppercase font-black">Sincronizando con PostgreSQL...</div>
            </div>
          </div>

          <div v-if="vistaActual === 'roles'" class="space-y-6">
            <div class="flex justify-between items-center">
              <div>
                <h2 class="text-2xl font-black text-white italic uppercase">Gestión de Roles</h2>
                <p class="text-xs text-neutral-500 font-bold uppercase tracking-widest">Crear, editar o eliminar roles del sistema</p>
              </div>
              <button @click="abrirNuevoRol" class="bg-brand-royal hover:bg-brand-royal/80 text-white px-6 py-2 rounded-xl font-bold text-xs transition-all shadow-lg shadow-brand-royal/20">
                + NUEVO ROL
              </button>
            </div>

            <div class="bg-white/5 border border-white/10 rounded-3xl overflow-hidden backdrop-blur-md">
              <table class="w-full text-left border-collapse">
                <thead class="bg-white/5 text-[10px] uppercase text-neutral-500 font-bold tracking-widest">
                  <tr>
                    <th class="p-5">Nombre</th>
                    <th class="p-5">Descripción</th>
                    <th class="p-5 text-right">Acciones</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-white/5">
                  <tr v-for="rol in roles" :key="rol.id" class="hover:bg-white/5 transition-colors">
                    <td class="p-5 font-bold text-white text-sm">{{ rol.nombre }}</td>
                    <td class="p-5 text-neutral-400 text-sm">{{ rol.descripcion || 'Sin descripción' }}</td>
                    <td class="p-5 text-right space-x-2">
                      <button @click="editarRol(rol)" class="text-brand-royal hover:text-white text-xs font-bold uppercase transition-colors">Editar</button>
                      <button @click="eliminarRol(rol.id)" class="text-red-400 hover:text-red-300 text-xs font-bold uppercase transition-colors">Eliminar</button>
                    </td>
                  </tr>
                </tbody>
              </table>
              <div v-if="roles.length === 0" class="p-10 text-center text-brand-royal animate-pulse uppercase font-black">No hay roles registrados.</div>
            </div>
          </div>

        </div>
      </main>
    </div>

    <div v-if="modalAbierto" class="fixed inset-0 bg-black/80 backdrop-blur-xl z-[100] flex items-center justify-center p-4">
      <div class="bg-brand-navy-base border border-white/10 w-full max-w-md rounded-3xl p-8 shadow-2xl">
        <h3 class="text-xl font-black text-white italic uppercase mb-6 tracking-tighter">Registrar Nuevo Usuario</h3>
        <div class="space-y-4">
          <div class="space-y-1">
            <label class="text-[10px] font-bold text-neutral-500 uppercase ml-2">Nombre Completo</label>
            <input v-model="form.nombre_completo" type="text" class="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-sm focus:outline-none focus:ring-2 ring-brand-royal/50 text-white">
          </div>
          <div class="grid grid-cols-2 gap-4">
          <div class="space-y-1">
    <label class="text-[10px] font-bold text-neutral-500 uppercase ml-2">Username</label>
    <input v-model="form.username" type="text" class="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-sm focus:outline-none text-white">
  </div>

  <div class="space-y-1">
    <label class="text-[10px] font-bold text-neutral-500 uppercase ml-2">Rol</label>
    <select 
      v-model="form.rol_id" 
      class="w-full bg-brand-navy-dark border border-white/10 rounded-xl px-4 py-3 text-sm text-white outline-none focus:ring-2 ring-brand-royal/50"
    >
      <option v-if="roles.length === 0" disabled value="">Cargando roles...</option>
      
      <option v-for="rol in roles" :key="rol.id" :value="rol.id">
        {{ rol.nombre }}
      </option>
    </select>
  </div>
          </div>
          <div class="space-y-1">
            <label class="text-[10px] font-bold text-neutral-500 uppercase ml-2">Contraseña Inicial</label>
            <input v-model="form.password" type="password" class="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-sm focus:outline-none text-white">
          </div>
          <div class="flex space-x-3 pt-4">
            <button @click="modalAbierto = false" class="flex-1 py-3 text-xs font-bold text-neutral-400 hover:text-white uppercase transition-colors">Cancelar</button>
            <button @click="guardarUsuario" class="flex-1 bg-brand-royal py-3 rounded-xl text-xs font-black text-white hover:bg-brand-royal/80 shadow-lg shadow-brand-royal/20 transition-all uppercase">Guardar</button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="rolModalAbierto" class="fixed inset-0 bg-black/80 backdrop-blur-xl z-[100] flex items-center justify-center p-4">
      <div class="bg-brand-navy-base border border-white/10 w-full max-w-md rounded-3xl p-8 shadow-2xl">
        <h3 class="text-xl font-black text-white italic uppercase mb-6 tracking-tighter">{{ rolEditando ? 'Editar Rol' : 'Crear Rol' }}</h3>
        <div class="space-y-4">
          <div class="space-y-1">
            <label class="text-[10px] font-bold text-neutral-500 uppercase ml-2">Nombre del Rol</label>
            <input v-model="rolForm.nombre" type="text" class="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-sm focus:outline-none focus:ring-2 ring-brand-royal/50 text-white">
          </div>
          <div class="space-y-1">
            <label class="text-[10px] font-bold text-neutral-500 uppercase ml-2">Descripción</label>
            <textarea v-model="rolForm.descripcion" rows="3" class="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-sm focus:outline-none focus:ring-2 ring-brand-royal/50 text-white"></textarea>
          </div>
          <div class="flex space-x-3 pt-4">
            <button @click="cerrarRolModal" class="flex-1 py-3 text-xs font-bold text-neutral-400 hover:text-white uppercase transition-colors">Cancelar</button>
            <button @click="guardarRol" class="flex-1 bg-brand-royal py-3 rounded-xl text-xs font-black text-white hover:bg-brand-royal/80 shadow-lg shadow-brand-royal/20 transition-all uppercase">Guardar</button>
          </div>
        </div>
      </div>
    </div>

  </div>
</template>

<style>
.scrollbar-none::-webkit-scrollbar { display: none; }
.scrollbar-none { -ms-overflow-style: none; scrollbar-width: none; }
</style>