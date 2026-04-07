<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';

// --- Estado de Usuario y Sesión ---
const nombreUsuario = ref('Cargando...');
const rolUsuario = ref('...');
const iniciales = ref('??');

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
const vistaActual = ref('inicio'); 
const listaRoles = ref([]);        
const cargando = ref(false);       

// --- Lógica de CRUD de Usuarios ---
const usuarios = ref([]);
const roles = ref([]);
const listaInvoices = ref([]);     
const listaClientesValue = ref([]); 
const listaProductosValue = ref([]); 
const listaTasas = ref({}); // Nueva: Tasas de cambio (USD/VES, USD/COP)
const modalAbierto = ref(false);
const rolModalAbierto = ref(false);
const rolEditando = ref(false);
const rolForm = ref({ id: null, nombre: '', descripcion: '' });
const modalVentaAbierto = ref(false); // Nueva: Modal de venta
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

// --- Gestión de Vistas ---
const cambiarVista = (nuevaVista) => {
  vistaActual.value = nuevaVista;
  if (nuevaVista === 'usuarios') {
    fetchUsuarios();
  } else if (nuevaVista === 'roles') {
    obtenerRoles();
  } else if (nuevaVista.startsWith('ventas_')) {
    const statusMap = {
      'ventas_facturas': 'FINAL',
      'ventas_devoluciones': 'RETURNED',
      'ventas_anulaciones': 'VOIDED'
    };
    fetchInvoices(statusMap[nuevaVista] || '');
  }
  if (mobileMenuVisible.value) mobileMenuVisible.value = false;
};

const logout = () => {
  localStorage.clear();
  window.location.href = '/';
};





const crearUsuario = async () => {
  const nombre_completo = prompt("Nombre completo:");
  if (!nombre_completo) return;
  const username = prompt("Nombre de usuario:");
  if (!username) return;
  const password = prompt("Contraseña:");
  if (!password) return;
  const rol_id_str = prompt("ID de Rol (1: Sistemas, etc.):");
  const rol_id = parseInt(rol_id_str);
  if (isNaN(rol_id)) return alert("ID de rol inválido");

  try {
    const response = await fetch('http://localhost:8080/api/usuarios', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ nombre_completo, username, password, rol_id })
    });
    if (response.ok) {
      alert("Usuario creado con éxito");
      obtenerUsuarios();
    } else {
      const data = await response.json();
      alert("Error: " + (data.error || "No se pudo crear"));
    }
  } catch (error) {
    console.error("Error creando usuario:", error);
  }
};


const fetchInvoices = async (status = '') => {
  cargando.value = true;
  try {
    const res = await axios.get(`${API_URL}/invoices?status=${status}`);
    listaInvoices.value = res.data || [];
  } catch (err) {
    console.error("Error al obtener facturas:", err);
  } finally {
    cargando.value = false;
  }
};

const abrirModalVenta = async () => {
  ventaForm.value = {
    client_id: null, currency: 'USD', apply_iva: true, apply_igtf: true, notes: '', items: []
  };
  totalesVenta.value = { subtotal: 0, iva: 0, igtf: 0, total: 0, equivalents: { USD: 0, VES: 0, COP: 0 } };
  modalVentaAbierto.value = true;
  
  try {
    const [c, p, r] = await Promise.all([
      axios.get(`${API_URL}/clients`),
      axios.get(`${API_URL}/products`),
      axios.get(`${API_URL}/rates`)
    ]);
    listaClientesValue.value = c.data || [];
    listaProductosValue.value = p.data || [];
    
    // Organizar tasas en objeto { 'VES': 36.5, 'COP': 3950 }
    const ratesObj = {};
    (r.data || []).forEach(rate => {
      ratesObj[rate.currency_code] = rate.rate;
    });
    listaTasas.value = ratesObj;
  } catch (err) { console.error(err); }
};

const actualizarPrecioItem = (idx) => {
  const item = ventaForm.value.items[idx];
  const prod = listaProductosValue.value.find(p => p.id === item.product_id);
  if (prod) {
    item.price = prod.price;
    calcularTotales();
  }
};

// Formulario de Venta
const ventaForm = ref({
  client_id: null,
  currency: 'USD',
  apply_iva: true,
  apply_igtf: true,
  notes: '',
  items: []
});

const agregarItemVenta = () => {
  ventaForm.value.items.push({ product_id: null, quantity: 1, price: 0 });
};

const eliminarItemVenta = (index) => {
  ventaForm.value.items.splice(index, 1);
};

const totalesVenta = ref({ subtotal: 0, iva: 0, igtf: 0, total: 0, equivalents: { USD: 0, VES: 0, COP: 0 } });

const calcularTotales = () => {
  let sub = 0;
  ventaForm.value.items.forEach(item => {
    sub += (item.quantity || 0) * (item.price || 0);
  });
  
  let iva = ventaForm.value.apply_iva ? sub * 0.16 : 0;
  let igtf = (ventaForm.value.apply_igtf && ventaForm.value.currency !== 'VES') ? (sub + iva) * 0.03 : 0;
  
  const finalTotal = sub + iva + igtf;
  
  // Calcular equivalencias
  const rates = listaTasas.value;
  let totalUSD = 0;
  
  if (ventaForm.value.currency === 'USD') totalUSD = finalTotal;
  else if (ventaForm.value.currency === 'VES') totalUSD = finalTotal / (rates.VES || 1);
  else if (ventaForm.value.currency === 'COP') totalUSD = finalTotal / (rates.COP || 1);

  totalesVenta.value = {
    subtotal: sub,
    iva: iva,
    igtf: igtf,
    total: finalTotal,
    equivalents: {
      USD: totalUSD,
      VES: totalUSD * (rates.VES || 0),
      COP: totalUSD * (rates.COP || 0)
    }
  };
};

const guardarVenta = async () => {
  if (!ventaForm.value.client_id || ventaForm.value.items.length === 0) {
    return alert("Completa los datos del cliente y al menos un producto");
  }
  try {
    const res = await axios.post(`${API_URL}/invoices`, {
      ...ventaForm.value,
      user_id: 1 // TODO: get from logged user
    });
    alert(res.data.message);
    modalVentaAbierto.value = false;
    cambiarVista('ventas_facturas');
  } catch (err) {
    alert("Error al guardar venta");
  }
};

const obtenerRoles = async () => {
  cargando.value = true;
  try {
    const response = await fetch('http://localhost:8080/api/roles');
    if (response.ok) {
      listaRoles.value = await response.json();
    }
  } catch (error) {
    console.error("Error cargando roles:", error);
  } finally {
    cargando.value = false;
  }
};

const crearRol = async () => {
  const nombre = prompt("Nombre del nuevo rol:");
  if (!nombre) return;
  const descripcion = prompt("Descripción (opcional):");

  try {
    const response = await fetch('http://localhost:8080/api/roles', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ nombre, descripcion })
    });
    if (response.ok) {
      alert("Rol creado con éxito");
      obtenerRoles();
    }
  } catch (error) {
    console.error("Error creando rol:", error);
  }
};

const formatearFecha = (fechaRaw) => {
  if (!fechaRaw) return '-';
  const fecha = new Date(fechaRaw);
  return fecha.toLocaleDateString('es-ES', { 
    year: 'numeric', month: 'short', day: 'numeric',
    hour: '2-digit', minute: '2-digit'
  });
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
            <button @click="cambiarVista('ventas_facturas')" class="w-full text-left px-3 py-2 text-xs font-semibold text-neutral-500 hover:text-white transition-colors">🧾 Facturas</button>
            <button @click="cambiarVista('ventas_facturas_nac')" class="w-full text-left px-3 py-2 text-xs font-semibold text-neutral-500 hover:text-white transition-colors">🧾 Facturas Nac.</button>
            <button @click="cambiarVista('ventas_devoluciones')" class="w-full text-left px-3 py-2 text-xs font-semibold text-neutral-500 hover:text-white transition-colors">↩️ Devoluciones</button>
            <button @click="cambiarVista('ventas_devoluciones_nac')" class="w-full text-left px-3 py-2 text-xs font-semibold text-neutral-500 hover:text-white transition-colors">↩️ Devoluciones Nac.</button>
            <button @click="cambiarVista('ventas_anulaciones')" class="w-full text-left px-3 py-2 text-xs font-semibold text-neutral-500 hover:text-white transition-colors">🚫 Anulaciones</button>
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
            <button @click="cambiarVista('usuarios')" class="w-full text-left px-3 py-2 text-xs font-semibold text-neutral-500 hover:text-white transition-colors">Usuarios</button>
            <button @click="cambiarVista('roles')" class="w-full text-left px-3 py-2 text-xs font-semibold text-neutral-500 hover:text-white transition-colors">Roles</button>
            <button
              v-for="item in ['Permisos', 'Empresas', 'Sucursales']"
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
          
          <!-- VISTA: INICIO -->
          <div v-if="vistaActual === 'inicio'" class="space-y-10">
            <!-- View Title -->
            <div class="flex flex-col space-y-1">
              <h1 class="text-3xl font-black text-white tracking-tighter uppercase italic">Panel de Control</h1>
              <p class="text-sm font-medium text-neutral-500 uppercase tracking-widest">Resumen ejecutivo del sistema</p>
            </div>

            <!-- Stats Grid -->
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
              <div v-for="stat in [
                { label: 'Ventas del Día', value: '$0.00', icon: '💰', color: 'from-emerald-500/20 to-emerald-500/5', border: 'border-emerald-500/20' },
                { label: 'Equipos Activos', value: '0', icon: '💻', color: 'from-brand-royal/20 to-brand-royal/5', border: 'border-brand-royal/20' },
                { label: 'Alertas Sistema', value: 'Ninguna', icon: '⚠️', color: 'from-amber-500/20 to-amber-500/5', border: 'border-amber-500/20' },
                { label: 'Transacciones', value: '1,240', icon: '📊', color: 'from-purple-500/20 to-purple-500/5', border: 'border-purple-500/20' }
              ]" :key="stat.label" 
                class="relative group"
              >
                <div :class="`h-full p-6 rounded-3xl bg-neutral-900/50 border ${stat.border} overflow-hidden transition-all duration-300 hover:-translate-y-1 hover:shadow-2xl hover:shadow-black/50`">
                  <div :class="`absolute inset-0 bg-gradient-to-br ${stat.color} opacity-0 group-hover:opacity-100 transition-opacity`"></div>
                  <div class="relative flex flex-col justify-between h-full space-y-6">
                    <div class="flex items-center justify-between">
                      <span class="text-2xl">{{ stat.icon }}</span>
                      <span class="text-[10px] font-black uppercase tracking-widest text-neutral-500 group-hover:text-white transition-colors">{{ stat.label }}</span>
                    </div>
                    <div class="text-2xl font-black text-white tracking-tight">{{ stat.value }}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
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

    <!-- MODAL: NUEVA VENTA -->
    <div v-if="modalVentaAbierto" class="fixed inset-0 bg-black/90 backdrop-blur-3xl z-[150] flex items-center justify-center p-4">
      <div class="bg-brand-navy-base border border-white/10 w-full max-w-4xl rounded-[40px] p-10 shadow-2xl relative overflow-hidden">
        <div class="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-brand-royal via-indigo-500 to-purple-600"></div>
        
        <div class="flex justify-between items-start mb-10">
          <div>
            <h3 class="text-3xl font-black text-white italic uppercase tracking-tighter">Generar Nueva Transacción</h3>
            <p class="text-[10px] font-bold text-neutral-500 uppercase tracking-[0.3em]">Módulo de Facturación Automática</p>
          </div>
          <button @click="modalVentaAbierto = false" class="bg-white/5 hover:bg-red-500/20 p-3 rounded-2xl transition-colors text-neutral-400 hover:text-red-500">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path d="M6 18L18 6M6 6l12 12"/></svg>
          </button>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-3 gap-10">
          <!-- Left: Client & Config -->
          <div class="space-y-8">
            <div class="space-y-4">
              <label class="text-[10px] font-black text-brand-royal uppercase tracking-widest">1. Datos del Cliente</label>
              <select v-model="ventaForm.client_id" @change="fetchAuxiliaresVenta" class="w-full bg-white/5 border border-white/10 rounded-2xl px-5 py-4 text-sm text-white outline-none focus:ring-2 ring-brand-royal transition-all">
                <option :value="null">Seleccionar Cliente...</option>
                <option v-for="c in listaClientesValue" :key="c.id" :value="c.id">{{ c.name }} ({{ c.rif }})</option>
              </select>
            </div>

            <div class="space-y-4">
              <label class="text-[10px] font-black text-brand-royal uppercase tracking-widest">2. Configuración Fiscal</label>
              <div class="bg-white/5 p-6 rounded-3xl space-y-4 border border-white/5">
                <div class="flex items-center justify-between">
                  <span class="text-xs font-bold text-neutral-400">Moneda Cobro</span>
                  <select v-model="ventaForm.currency" @change="calcularTotales" class="bg-brand-navy-dark text-[10px] font-black text-white border-0 rounded-lg px-2 py-1 outline-none">
                    <option value="USD">USD ($)</option>
                    <option value="VES">VES (Bs.)</option>
                    <option value="COP">COP (Pesos)</option>
                  </select>
                </div>
                <div class="flex items-center justify-between">
                  <span class="text-xs font-bold text-neutral-400 italic">Aplicar IVA (16%)</span>
                  <input type="checkbox" v-model="ventaForm.apply_iva" @change="calcularTotales" class="accent-brand-royal">
                </div>
                <div v-if="ventaForm.currency !== 'VES'" class="flex items-center justify-between">
                  <span class="text-xs font-bold text-indigo-400 italic">Aplicar IGTF (3%)</span>
                  <input type="checkbox" v-model="ventaForm.apply_igtf" @change="calcularTotales" class="accent-indigo-500">
                </div>
              </div>
            </div>
          </div>

          <!-- Middle: Items -->
          <div class="lg:col-span-2 space-y-8">
            <div class="flex justify-between items-center">
              <label class="text-[10px] font-black text-brand-royal uppercase tracking-widest">3. Detalle de Productos</label>
              <button @click="agregarItemVenta" class="text-[10px] font-black text-emerald-400 hover:text-emerald-300 uppercase tracking-widest">+ Añadir Ítem</button>
            </div>

            <div class="max-h-[300px] overflow-y-auto space-y-3 pr-2 scrollbar-none">
              <div v-for="(item, idx) in ventaForm.items" :key="idx" class="flex items-center space-x-3 group">
                <select v-model="item.product_id" @change="actualizarPrecioItem(idx)" class="flex-1 bg-white/5 border border-white/10 rounded-2xl px-4 py-3 text-xs text-white outline-none">
                  <option :value="null">Producto...</option>
                  <option v-for="p in listaProductosValue" :key="p.id" :value="p.id">{{ p.name }} (${{ p.price }})</option>
                </select>
                <input type="number" v-model="item.quantity" @input="calcularTotales" placeholder="Cant" class="w-20 bg-white/5 border border-white/10 rounded-2xl px-4 py-3 text-xs text-white text-center">
                <input type="number" v-model="item.price" @input="calcularTotales" placeholder="Precio" class="w-24 bg-white/5 border border-white/10 rounded-2xl px-4 py-3 text-xs text-white text-right">
                <button @click="eliminarItemVenta(idx)" class="p-3 text-neutral-600 hover:text-red-500 transition-colors opacity-0 group-hover:opacity-100 italic font-black">X</button>
              </div>
              <div v-if="ventaForm.items.length === 0" class="text-center py-10 border-2 border-dashed border-white/5 rounded-3xl text-neutral-600 text-xs font-bold uppercase italic">
                La lista de productos está vacía
              </div>
            </div>

            <!-- Totals Footer -->
            <div class="pt-6 border-t border-white/10 grid grid-cols-1 lg:grid-cols-2 gap-8">
              <div class="space-y-4">
                <div class="p-4 bg-brand-navy-dark/50 rounded-2xl border border-white/5 space-y-2">
                  <p class="text-[8px] font-black text-neutral-500 uppercase tracking-widest mb-2">Previsualización Multimoneda</p>
                  <div class="flex justify-between items-center text-[10px] font-bold">
                    <span class="text-emerald-400">TOTAL USD:</span>
                    <span class="text-white">$ {{ totalesVenta.equivalents.USD.toFixed(2) }}</span>
                  </div>
                  <div class="flex justify-between items-center text-[10px] font-bold">
                    <span class="text-brand-royal">TOTAL VES:</span>
                    <span class="text-white">Bs. {{ totalesVenta.equivalents.VES.toLocaleString() }}</span>
                  </div>
                  <div class="flex justify-between items-center text-[10px] font-bold">
                    <span class="text-amber-400">TOTAL COP:</span>
                    <span class="text-white">$ {{ (totalesVenta.equivalents.COP / 1000).toFixed(0) }}k ({{ totalesVenta.equivalents.COP.toLocaleString() }})</span>
                  </div>
                </div>
                <div class="flex justify-between text-[10px] font-bold text-neutral-500 uppercase px-2"><span>Subtotal (Base)</span> <span class="text-white">{{ ventaForm.currency }} {{ totalesVenta.subtotal.toFixed(2) }}</span></div>
                <div class="flex justify-between text-[10px] font-bold text-neutral-500 uppercase px-2"><span>IVA (16%)</span> <span class="text-white">{{ ventaForm.currency }} {{ totalesVenta.iva.toFixed(2) }}</span></div>
                <div v-if="totalesVenta.igtf > 0" class="flex justify-between text-[10px] font-bold text-indigo-400 uppercase px-2"><span>IGTF (3%)</span> <span>{{ ventaForm.currency }} {{ totalesVenta.igtf.toFixed(2) }}</span></div>
              </div>
              
              <div class="flex flex-col items-end justify-center">
                <div class="text-[10px] font-black text-brand-royal uppercase tracking-[0.2em] mb-1">Total a Cobrar</div>
                <div class="text-5xl font-black text-white italic tracking-tighter">{{ ventaForm.currency }} {{ totalesVenta.total.toFixed(2) }}</div>
              </div>
            </div>

            <div class="flex space-x-4 pt-6">
              <button @click="modalVentaAbierto = false" class="flex-1 py-4 text-xs font-bold text-neutral-500 hover:text-white uppercase transition-colors">Abortar Operación</button>
              <button @click="guardarVenta" class="flex-[2] bg-gradient-to-r from-brand-royal to-indigo-600 py-4 rounded-2xl text-xs font-black text-white hover:scale-[1.02] active:scale-[0.98] transition-all shadow-xl shadow-brand-royal/20 uppercase tracking-widest">
                Confirmar y Registrar Venta
              </button>
            </div>
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