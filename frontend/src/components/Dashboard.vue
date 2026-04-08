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

// --- Estado de Devoluciones ---
const searchResult = ref(null);
const searchNum = ref('');
const historialDevoluciones = ref([]);
const modalDevolucionConfirmar = ref(false);
const devolucionForm = ref({
  original_id: null,
  motivo: 'Garantía / Defecto',
  detalles: []
});

// --- Lógica de CRUD de Usuarios ---
const usuarios = ref([]);
const roles = ref([]);
const listaInvoices = ref([]);     
const listaClientesValue = ref([]); 
const listaProductosValue = ref([]); 
const listaTasas = ref({}); 
const modalAbierto = ref(false);
const rolModalAbierto = ref(false);
const rolEditando = ref(false);
const rolForm = ref({ id: null, nombre: '', descripcion: '' });
const modalVentaAbierto = ref(false); 
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

// 2. Guardar Usuario (Create/Update)
const guardarUsuario = async () => {
  try {
    if (form.value.id) {
        await axios.put(`${API_URL}/usuarios/${form.value.id}`, form.value);
    } else {
        await axios.post(`${API_URL}/usuarios`, form.value);
    }
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

const fetchRoles = async () => {
  try {
    const res = await axios.get(`${API_URL}/roles`);
    roles.value = res.data || [];
    if (roles.value.length > 0 && !form.value.rol_id) {
      form.value.rol_id = roles.value[0].id;
    }
  } catch (err) {
    console.error("Error al obtener roles:", err);
  }
};

// --- Inicialización ---
onMounted(() => {
  const isLogged = localStorage.getItem('isLogged');
  const nombre = localStorage.getItem('usuario') || (isLogged ? '' : 'Jesus Sarmiento');
  const rol = localStorage.getItem('rol') || (isLogged ? '' : 'Administrador');
  
  nombreUsuario.value = nombre ? nombre : "Usuario Demo";
  rolUsuario.value = rol ? rol : "Administrador";
  iniciales.value = nombreUsuario.value.split(' ').filter(n => n).map(n => n[0]).join('').toUpperCase().substring(0, 2);

  if (isLogged) {
    fetchUsuarios();
    fetchRoles();
  }
});

// --- Funciones de Navegación ---
const toggleSubmenu = (menuKey) => {
  menus.value[menuKey] = !menus.value[menuKey];
};

const navegar = (vista) => {
  vistaActual.value = vista;
  if (window.innerWidth < 1024) mobileMenuVisible.value = false;
};

const cambiarVista = (nuevaVista) => {
  vistaActual.value = nuevaVista;
  if (nuevaVista === 'usuarios') {
    fetchUsuarios();
  } else if (nuevaVista === 'roles') {
    fetchRoles();
  } else if (nuevaVista === 'ventas_devoluciones') {
    fetchHistorialDevoluciones();
  } else if (nuevaVista.startsWith('ventas_')) {
    const tipoMap = {
      'ventas_facturas': '1',
      'ventas_facturas_nac': '5'
    };
    fetchInvoices(tipoMap[nuevaVista] || '');
  }
  if (mobileMenuVisible.value) mobileMenuVisible.value = false;
};

const logout = () => {
  localStorage.clear();
  window.location.href = '/';
};

const fetchInvoices = async (tipo = '') => {
  cargando.value = true;
  try {
    const res = await axios.get(`${API_URL}/documentos?tipo=${tipo}`);
    listaInvoices.value = res.data || [];
  } catch (err) {
    console.error("Error al obtener documentos:", err);
  } finally {
    cargando.value = false;
  }
};

const buscarFacturaDevolucion = async () => {
  if (!searchNum.value) return;
  cargando.value = true;
  try {
    const res = await axios.get(`${API_URL}/documentos/buscar?numero=${searchNum.value}`);
    searchResult.value = res.data;
  } catch (err) {
    alert(err.response?.data?.error || "Error en la búsqueda");
    searchResult.value = null;
  } finally {
    cargando.value = false;
  }
};

const prepararDevolucion = () => {
  if (!searchResult.value) return;
  devolucionForm.value = {
    original_id: searchResult.value.documento.id,
    motivo: 'Garantía / Defecto',
    detalles: searchResult.value.detalles.map(d => ({
      producto_id: d.producto_id,
      nombre: d.producto_nombre,
      cantidad: d.cantidad,
      seleccionado: true
    }))
  };
  modalDevolucionConfirmar.value = true;
};

const ejecutarDevolucion = async () => {
  const itemsParaDevolver = devolucionForm.value.detalles.filter(i => i.seleccionado);
  if (itemsParaDevolver.length === 0) return alert("Selecciona al menos un ítem");

  try {
    const res = await axios.post(`${API_URL}/documentos/devolucion`, {
      original_id: devolucionForm.value.original_id,
      motivo: devolucionForm.value.motivo,
      detalles: itemsParaDevolver.map(i => ({
        producto_id: i.producto_id,
        cantidad: i.cantidad
      }))
    });
    alert(res.data.message);
    modalDevolucionConfirmar.value = false;
    searchResult.value = null;
    cambiarVista('ventas_devoluciones');
  } catch (err) {
    alert("Error al procesar devolución");
  }
};

const fetchHistorialDevoluciones = async () => {
  cargando.value = true;
  try {
    const res = await axios.get(`${API_URL}/devoluciones/historial`);
    historialDevoluciones.value = res.data || [];
  } catch (err) {
    console.error(err);
  } finally {
    cargando.value = false;
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

const totalesVenta = ref({ subtotal: 0, iva: 0, igtf: 0, total: 0, equivalents: { USD: 0, VES: 0, COP: 0 } });

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

const agregarItemVenta = () => {
  ventaForm.value.items.push({ product_id: null, quantity: 1, price: 0 });
};

const eliminarItemVenta = (index) => {
  ventaForm.value.items.splice(index, 1);
};

const calcularTotales = () => {
  let sub = 0;
  ventaForm.value.items.forEach(item => {
    sub += (item.quantity || 0) * (item.price || 0);
  });
  
  let iva = ventaForm.value.apply_iva ? sub * 0.16 : 0;
  let igtf = (ventaForm.value.apply_igtf && ventaForm.value.currency !== 'VES') ? (sub + iva) * 0.03 : 0;
  
  const finalTotal = sub + iva + igtf;
  
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
    await axios.post(`${API_URL}/invoices`, {
      ...ventaForm.value,
      user_id: 1 
    });
    alert("Venta registrada con éxito");
    modalVentaAbierto.value = false;
    cambiarVista('ventas_facturas');
  } catch (err) {
    alert("Error al guardar venta");
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
        
        <button @click="navegar('inicio')" class="w-full flex items-center px-4 py-3 hover:bg-white/5 rounded-2xl transition-all group">
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
          
          <div v-if="vistaActual === 'inicio'" class="space-y-10">
            <div class="flex flex-col space-y-1">
              <h1 class="text-3xl font-black text-white tracking-tighter uppercase italic">Panel de Control</h1>
              <p class="text-sm font-medium text-neutral-500 uppercase tracking-widest">Resumen ejecutivo del sistema</p>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
              <div v-for="stat in [
                { label: 'Ventas del Día', value: '$0.00', icon: '💰', color: 'from-emerald-500/20 to-emerald-500/5', border: 'border-emerald-500/20' },
                { label: 'Equipos Activos', value: '0', icon: '💻', color: 'from-brand-royal/20 to-brand-royal/5', border: 'border-brand-royal/20' },
                { label: 'Alertas Sistema', value: 'Ninguna', icon: '⚠️', color: 'from-amber-500/20 to-amber-500/5', border: 'border-amber-500/20' },
                { label: 'Usuarios en Sistema', value: usuarios.length, icon: '👤', color: 'from-purple-500/20 to-purple-500/5', border: 'border-purple-500/20' }
              ]" :key="stat.label" class="relative group">
                <div :class="`h-full p-6 rounded-3xl bg-neutral-900/50 border ${stat.border} overflow-hidden transition-all duration-300` text-white">
                  <div class="text-2xl mb-4">{{ stat.icon }}</div>
                  <div class="text-xs font-bold text-neutral-500 uppercase">{{ stat.label }}</div>
                  <div class="text-2xl font-black">{{ stat.value }}</div>
                </div>
              </div>
            </div>
          </div>

          <!-- VISTA: TRANSACCIONES (FACTURAS / NACIONALES) -->
          <div v-if="vistaActual.startsWith('ventas_') && vistaActual !== 'ventas_devoluciones'" class="space-y-6">
            <div class="flex justify-between items-center">
              <div>
                <h2 class="text-2xl font-black text-white italic uppercase">{{ vistaActual === 'ventas_facturas' ? 'Facturación Estándar' : 'Facturación Nacional' }}</h2>
                <p class="text-[10px] text-neutral-500 font-bold uppercase tracking-[0.3em]">Lista Maestra de Documentos</p>
              </div>
              <button @click="abrirModalVenta" class="bg-brand-royal px-6 py-3 rounded-2xl text-[10px] font-black text-white uppercase tracking-widest hover:scale-105 transition-all">+ Nueva Factura</button>
            </div>

            <div class="bg-white/5 border border-white/10 rounded-3xl overflow-hidden">
               <table class="w-full text-left">
                <thead class="bg-white/5 text-[10px] uppercase text-neutral-500 font-bold tracking-widest">
                  <tr>
                    <th class="p-5">Número</th>
                    <th class="p-5">Cliente</th>
                    <th class="p-5">Fecha</th>
                    <th class="p-5 text-right">Monto Total</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-white/5">
                  <tr v-for="doc in listaInvoices" :key="doc.id" class="hover:bg-white/5 transition-colors">
                    <td class="p-5 font-bold text-white text-sm">#{{ doc.numero_documento }}</td>
                    <td class="p-5 text-neutral-400 font-semibold text-xs">{{ doc.cliente_nombre?.String || 'Cliente Final' }}</td>
                    <td class="p-5 text-neutral-400 text-xs">{{ doc.fecha_emision }}</td>
                    <td class="p-5 text-right font-black text-white text-sm">$ {{ doc.total_neto.toLocaleString() }}</td>
                  </tr>
                   <tr v-if="listaInvoices.length === 0">
                    <td colspan="4" class="p-10 text-center text-neutral-600 font-bold uppercase text-xs italic">No hay documentos registrados para esta sección</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <!-- VISTA: DEVOLUCIONES -->
          <div v-if="vistaActual === 'ventas_devoluciones'" class="space-y-6">
            <div class="flex justify-between items-center">
              <div>
                <h2 class="text-2xl font-black text-white italic uppercase">Historial de Devoluciones</h2>
                <p class="text-[10px] text-neutral-500 font-bold uppercase tracking-[0.3em]">Gestión de Notas de Crédito</p>
              </div>
            </div>

            <div class="bg-white/5 p-8 rounded-[32px] border border-white/10 space-y-4 shadow-xl">
              <h3 class="text-xs font-black text-brand-royal uppercase italic tracking-widest">Ejecutar Nueva Devolución</h3>
              <div class="flex space-x-4">
                <input v-model="searchNum" type="text" placeholder="Ej: FAC-0001 o FAN-0001" class="flex-1 bg-brand-navy-dark border border-white/10 rounded-2xl px-6 py-4 text-sm text-white outline-none focus:ring-2 ring-brand-royal transition-all">
                <button @click="buscarFacturaDevolucion" class="bg-brand-royal px-10 rounded-2xl text-[10px] font-black text-white uppercase tracking-widest hover:bg-brand-royal/80 transition-all">Consultar</button>
              </div>
              
              <div v-if="searchResult" class="mt-6 p-6 bg-brand-navy-dark rounded-3xl border border-brand-royal/30 flex justify-between items-center animate-in fade-in zoom-in duration-300 shadow-2xl">
                <div>
                  <p class="text-[8px] font-black text-brand-royal uppercase tracking-[0.3em] mb-1">Documento Encontrado</p>
                  <p class="text-xl font-bold text-white italic uppercase tracking-tighter">{{ searchResult.documento.numero_documento }} — {{ searchResult.documento.cliente_nombre?.String || 'Cliente Especial' }}</p>
                  <p class="text-[10px] text-neutral-500 font-bold">MONTO BASE: $ {{ searchResult.documento.total_neto.toLocaleString() }} | FECHA: {{ searchResult.documento.fecha_emision }}</p>
                </div>
                <button @click="prepararDevolucion" class="bg-brand-royal/10 hover:bg-brand-royal text-brand-royal hover:text-white px-8 py-3 rounded-2xl text-[10px] font-black transition-all border border-brand-royal/20 uppercase tracking-widest">Proceder a Devolución</button>
              </div>
            </div>

            <div class="bg-white/5 border border-white/10 rounded-3xl overflow-hidden shadow-2xl">
              <table class="w-full text-left">
                <thead class="bg-white/5 text-[10px] uppercase text-neutral-500 font-bold tracking-widest">
                  <tr>
                    <th class="p-5">Número DEV</th>
                    <th class="p-5">Ref. Original</th>
                    <th class="p-5">Motivo</th>
                    <th class="p-5">Fecha</th>
                    <th class="p-5 text-right">Monto</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-white/5">
                   <tr v-for="doc in historialDevoluciones" :key="doc.id" class="hover:bg-white/5 transition-colors">
                    <td class="p-5 font-bold text-white text-sm">#{{ doc.numero_documento }}</td>
                    <td class="p-5 text-amber-500 font-black text-xs italic">{{ doc.numero_referencia?.String || '-' }}</td>
                    <td class="p-5 text-neutral-400 font-semibold text-xs">{{ doc.motivo?.String || 'Garantía' }}</td>
                    <td class="p-5 text-neutral-400 text-xs">{{ doc.fecha_emision }}</td>
                    <td class="p-5 text-right font-black text-white text-sm">$ {{ doc.total_neto.toLocaleString() }}</td>
                  </tr>
                  <tr v-if="historialDevoluciones.length === 0">
                    <td colspan="5" class="p-10 text-center text-neutral-600 font-bold uppercase text-xs italic">No hay historial de devoluciones</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <!-- VISTA: USUARIOS -->
          <div v-if="vistaActual === 'usuarios'" class="space-y-6">
            <div class="flex justify-between items-center">
              <div>
                <h2 class="text-2xl font-black text-white italic uppercase">Usuarios Sistema</h2>
                <p class="text-[10px] text-neutral-500 font-bold uppercase tracking-[0.3em]">Administración de Accesos</p>
              </div>
              <button @click="modalAbierto = true" class="bg-brand-royal px-6 py-3 rounded-2xl text-[10px] font-black text-white uppercase tracking-widest">+ Nuevo Usuario</button>
            </div>

            <div class="bg-white/5 border border-white/10 rounded-3xl overflow-hidden">
               <table class="w-full text-left">
                <thead class="bg-white/5 text-[10px] uppercase text-neutral-500 font-bold tracking-widest">
                  <tr>
                    <th class="p-5">Nombre</th>
                    <th class="p-5">Username</th>
                    <th class="p-5">Rol</th>
                    <th class="p-5 text-right">Acciones</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-white/5">
                  <tr v-for="u in usuarios" :key="u.id" class="hover:bg-white/5 transition-colors">
                    <td class="p-5 font-bold text-white text-sm">{{ u.nombre_completo }}</td>
                    <td class="p-5 text-neutral-400 font-semibold text-xs">{{ u.username }}</td>
                    <td class="p-5 text-neutral-400 text-xs">{{ u.rol_nombre }}</td>
                    <td class="p-5 text-right space-x-3">
                        <button @click="form = {...u}; modalAbierto = true" class="text-brand-royal text-[10px] font-black uppercase hover:text-white">Editar</button>
                        <button @click="eliminarUsuario(u.id)" class="text-red-500 text-[10px] font-black uppercase hover:text-red-400">Borrar</button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <!-- VISTA: ROLES -->
          <div v-if="vistaActual === 'roles'" class="space-y-6">
            <div class="flex justify-between items-center">
              <div>
                <h2 class="text-2xl font-black text-white italic uppercase">Roles y Permisos</h2>
                <p class="text-[10px] text-neutral-500 font-bold uppercase tracking-[0.3em]">Perfiles de Seguridad</p>
              </div>
              <button @click="abrirNuevoRol" class="bg-brand-royal px-6 py-3 rounded-2xl text-[10px] font-black text-white uppercase tracking-widest">+ Crear Rol</button>
            </div>

            <div class="bg-white/5 border border-white/10 rounded-3xl overflow-hidden">
               <table class="w-full text-left">
                <thead class="bg-white/5 text-[10px] uppercase text-neutral-500 font-bold tracking-widest">
                  <tr>
                    <th class="p-5">Rol</th>
                    <th class="p-5">Descripción</th>
                    <th class="p-5 text-right">Acciones</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-white/5">
                  <tr v-for="r in roles" :key="r.id" class="hover:bg-white/5 transition-colors">
                    <td class="p-5 font-bold text-white text-sm">{{ r.nombre }}</td>
                    <td class="p-5 text-neutral-400 font-semibold text-xs">{{ r.descripcion || '-' }}</td>
                    <td class="p-5 text-right space-x-3">
                        <button @click="editarRol(r)" class="text-brand-royal text-[10px] font-black uppercase hover:text-white">Editar</button>
                        <button @click="eliminarRol(r.id)" class="text-red-500 text-[10px] font-black uppercase hover:text-red-400">Borrar</button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

        </div>
      </main>
    </div>

    <!-- MODAL: USUARIO -->
    <div v-if="modalAbierto" class="fixed inset-0 bg-black/80 backdrop-blur-xl z-[100] flex items-center justify-center p-4">
      <div class="bg-brand-navy-base border border-white/10 w-full max-w-md rounded-3xl p-8 shadow-2xl">
        <h3 class="text-xl font-black text-white italic uppercase mb-6 tracking-tighter">Registrar Usuario</h3>
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
              <select v-model="form.rol_id" class="w-full bg-brand-navy-dark border border-white/10 rounded-xl px-4 py-3 text-sm text-white outline-none">
                <option v-for="r in roles" :key="r.id" :value="r.id">{{ r.nombre }}</option>
              </select>
            </div>
          </div>
          <div class="space-y-1">
            <label class="text-[10px] font-bold text-neutral-500 uppercase ml-2">Contraseña</label>
            <input v-model="form.password" type="password" class="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-sm focus:outline-none text-white">
          </div>
          <div class="flex space-x-3 pt-4">
            <button @click="modalAbierto = false; resetForm()" class="flex-1 py-3 text-xs font-bold text-neutral-400 hover:text-white uppercase transition-colors">Cancelar</button>
            <button @click="guardarUsuario" class="flex-1 bg-brand-royal py-3 rounded-xl text-xs font-black text-white hover:bg-brand-royal/80 shadow-lg shadow-brand-royal/20 transition-all uppercase">Guardar</button>
          </div>
        </div>
      </div>
    </div>

    <!-- MODAL: ROL -->
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

    <!-- MODAL: CONFIRMAR DEVOLUCION -->
    <div v-if="modalDevolucionConfirmar" class="fixed inset-0 bg-black/90 backdrop-blur-3xl z-[200] flex items-center justify-center p-4">
      <div class="bg-brand-navy-base border border-white/10 w-full max-w-2xl rounded-[40px] p-10 shadow-2xl relative overflow-hidden">
        <div class="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-amber-500 via-orange-500 to-red-600"></div>
        
        <div class="mb-8">
          <h3 class="text-3xl font-black text-white italic uppercase tracking-tighter">Procesar Devolución</h3>
          <p class="text-[10px] font-bold text-amber-500 uppercase tracking-[0.3em]">Nota de Crédito Vinculada</p>
        </div>

        <div class="space-y-6">
          <div class="space-y-2">
            <label class="text-[10px] font-black text-neutral-500 uppercase tracking-widest">Seleccionar Ítems a Devolver</label>
            <div class="bg-white/5 rounded-3xl border border-white/5 overflow-hidden max-h-[250px] overflow-y-auto scrollbar-none">
              <div v-for="(item, idx) in devolucionForm.detalles" :key="idx" class="flex items-center p-4 border-b border-white/5 last:border-0">
                <input type="checkbox" v-model="item.seleccionado" class="w-5 h-5 accent-brand-royal rounded-lg mr-4">
                <div class="flex-1">
                  <p class="text-sm font-bold text-white">{{ item.nombre }}</p>
                  <p class="text-[10px] text-neutral-500">Cantidad original: {{ item.cantidad }}</p>
                </div>
                <div class="flex items-center space-x-3">
                  <span class="text-[10px] font-black text-neutral-400 uppercase">DEVOLVER</span>
                  <input type="number" v-model="item.cantidad" class="w-16 bg-brand-navy-dark border border-white/10 rounded-xl px-3 py-1.5 text-xs text-white text-center">
                </div>
              </div>
            </div>
          </div>

          <div class="space-y-2">
            <label class="text-[10px] font-black text-neutral-500 uppercase tracking-widest">Motivo de la Devolución</label>
            <textarea v-model="devolucionForm.motivo" class="w-full bg-white/5 border border-white/10 rounded-2xl px-5 py-4 text-sm text-white outline-none h-24 focus:ring-2 ring-amber-500/50"></textarea>
          </div>

          <div class="flex space-x-4 pt-6">
            <button @click="modalDevolucionConfirmar = false" class="flex-1 py-4 text-xs font-bold text-neutral-500 hover:text-white uppercase transition-colors">Cancelar</button>
            <button @click="ejecutarDevolucion" class="flex-[2] bg-gradient-to-r from-amber-500 to-orange-600 py-4 rounded-2xl text-xs font-black text-white hover:scale-[1.02] active:scale-[0.98] transition-all shadow-xl shadow-amber-500/20 uppercase tracking-widest">
              Confirmar Devolución
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- MODAL: NUEVA VENTA (RESTORED) -->
    <div v-if="modalVentaAbierto" class="fixed inset-0 bg-black/90 backdrop-blur-3xl z-[150] flex items-center justify-center p-4">
       <!-- Content from stashed version previously seen -->
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
          <div class="space-y-8">
            <div class="space-y-4">
              <label class="text-[10px] font-black text-brand-royal uppercase tracking-widest">1. Datos del Cliente</label>
              <select v-model="ventaForm.client_id" class="w-full bg-white/5 border border-white/10 rounded-2xl px-5 py-4 text-sm text-white outline-none">
                <option :value="null">Seleccionar Cliente...</option>
                <option v-for="c in listaClientesValue" :key="c.id" :value="c.id">{{ c.name }}</option>
              </select>
            </div>
            <div class="space-y-4">
              <label class="text-[10px] font-black text-brand-royal uppercase tracking-widest">2. Configuración Fiscal</label>
              <div class="bg-white/5 p-6 rounded-3xl space-y-4">
                <div class="flex justify-between items-center"><span class="text-xs font-bold text-neutral-400">IVA (16%)</span><input type="checkbox" v-model="ventaForm.apply_iva" @change="calcularTotales"></div>
                <div class="flex justify-between items-center"><span class="text-xs font-bold text-neutral-400">IGTF (3%)</span><input type="checkbox" v-model="ventaForm.apply_igtf" @change="calcularTotales"></div>
              </div>
            </div>
          </div>
          <div class="lg:col-span-2 space-y-8">
            <div class="flex justify-between items-center"><label class="text-[10px] font-black text-brand-royal uppercase">3. Detalle</label><button @click="agregarItemVenta" class="text-[10px] font-black text-emerald-400 uppercase">+ Ítem</button></div>
            <div class="max-h-[250px] overflow-y-auto space-y-3 pr-2 scrollbar-none">
              <div v-for="(item, idx) in ventaForm.items" :key="idx" class="flex items-center space-x-3">
                 <select v-model="item.product_id" @change="actualizarPrecioItem(idx)" class="flex-1 bg-white/5 border-0 rounded-2xl px-4 py-3 text-xs text-white">
                   <option v-for="p in listaProductosValue" :key="p.id" :value="p.id">{{ p.name }}</option>
                 </select>
                 <input type="number" v-model="item.quantity" @input="calcularTotales" class="w-20 bg-white/5 border-0 rounded-2xl px-4 py-3 text-xs text-center text-white">
                 <button @click="eliminarItemVenta(idx)" class="text-red-500 font-black">X</button>
              </div>
            </div>
            <div class="pt-6 border-t border-white/10 flex justify-between items-center">
               <div class="text-3xl font-black text-white italic">{{ ventaForm.currency }} {{ totalesVenta.total.toFixed(2) }}</div>
               <button @click="guardarVenta" class="bg-brand-royal px-10 py-4 rounded-2xl text-[10px] font-black text-white uppercase tracking-widest">Registrar Venta</button>
            </div>
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