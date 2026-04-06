<script setup>
import { ref, onMounted } from 'vue';

// --- Estado de Usuario ---
const nombreUsuario = ref('Cargando...');
const rolUsuario = ref('...');
const iniciales = ref('??');

// --- Estado de la Interfaz (Menús) ---
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
  <div class="flex h-screen bg-brand-navy-dark font-sans text-neutral-200 overflow-hidden">
    <!-- Mobile Overlay -->
    <div 
      v-if="mobileMenuVisible" 
      class="fixed inset-0 bg-black/60 backdrop-blur-sm z-[60] lg:hidden" 
      @click="toggleMobileMenu"
    ></div>

    <!-- Sidebar -->
    <aside 
      class="fixed lg:static inset-y-0 left-0 w-72 bg-brand-navy-base/80 backdrop-blur-xl border-r border-white/5 z-[70] transform transition-transform duration-300 lg:translate-x-0"
      :class="mobileMenuVisible ? 'translate-x-0' : '-translate-x-full'"
    >
      <!-- Sidebar Header -->
      <div class="h-20 flex items-center justify-between px-8 border-b border-white/5">
        <div class="flex items-center space-x-3">
          <div class="w-2.5 h-2.5 bg-brand-royal rounded-full shadow-[0_0_10px_rgba(37,99,235,0.5)]"></div>
          <h2 class="text-xl font-black tracking-tighter text-white uppercase italic">
            VINYL<span class="text-brand-royal">NET</span>
          </h2>
        </div>
        <button class="lg:hidden text-neutral-400 hover:text-white" @click="toggleMobileMenu">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
        </button>
      </div>
      
      <!-- Sidebar Navigation -->
      <nav class="p-6 h-[calc(100vh-5rem)] overflow-y-auto space-y-2 scrollbar-none">
        <div class="px-2 pb-2 text-[10px] font-bold uppercase tracking-[0.2em] text-neutral-600">
          Operaciones
        </div>
        
        <button class="w-full flex items-center px-4 py-3 bg-brand-royal/10 border border-brand-royal/20 text-brand-royal rounded-2xl group transition-all">
          <span class="mr-3 text-lg opacity-80 group-hover:scale-110 transition-transform">🏠</span>
          <span class="font-bold text-sm">Inicio</span>
        </button>

        <!-- Ficheros Submenu -->
        <div class="space-y-1">
          <button 
            @click="toggleSubmenu('ficheros')"
            class="w-full flex items-center justify-between px-4 py-3 hover:bg-white/5 rounded-2xl transition-all group"
          >
            <div class="flex items-center">
              <span class="mr-3 text-lg opacity-80 group-hover:scale-110 transition-transform text-neutral-400">🗂️</span>
              <span class="font-bold text-sm text-neutral-300">Ficheros</span>
            </div>
            <svg 
              class="w-4 h-4 text-neutral-600 transition-transform" 
              :class="{ 'rotate-180': menus.ficheros }"
              fill="none" stroke="currentColor" viewBox="0 0 24 24"
            ><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/></svg>
          </button>
          
          <div v-show="menus.ficheros" class="ml-9 border-l border-white/5 pl-4 space-y-1">
            <button v-for="item in ['Departamentos', 'Depósitos', 'Inventarios', 'Clientes', 'Vendedores', 'Bancos']" :key="item"
              class="w-full text-left px-3 py-2 text-xs font-semibold text-neutral-500 hover:text-white transition-colors"
            >
              {{ item }}
            </button>
          </div>
        </div>

        <!-- Transacciones Submenu -->
        <div class="space-y-1">
          <button 
            @click="toggleSubmenu('transacciones')"
            class="w-full flex items-center justify-between px-4 py-3 hover:bg-white/5 rounded-2xl transition-all group"
          >
            <div class="flex items-center">
              <span class="mr-3 text-lg opacity-80 group-hover:scale-110 transition-transform text-neutral-400">📦</span>
              <span class="font-bold text-sm text-neutral-300">Transacciones</span>
            </div>
            <svg 
              class="w-4 h-4 text-neutral-600 transition-transform" 
              :class="{ 'rotate-180': menus.transacciones }"
              fill="none" stroke="currentColor" viewBox="0 0 24 24"
            ><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/></svg>
          </button>
          
          <div v-show="menus.transacciones" class="ml-9 border-l border-white/5 pl-4 space-y-3 pt-2">
            <!-- Nested Submenus for Transacciones -->
            <div v-for="(sub, subKey) in {
              'Inventario': { key: 'inventario_sub', icon: '📋', items: ['➕ Cargos', '➖ Descargos', '↔️ Transferencia'] },
              'Mov. de Caja': { key: 'caja_sub', icon: '📠', items: ['➕ Préstamos', '➖ Avances'] },
              'Clientes': { key: 'clientes_sub', icon: '👤', items: ['💰 Cuentas x Cobrar'] },
              'Bancos': { key: 'bancos_sub', icon: '🏦', items: ['💸 Movimientos'] }
            }" :key="subKey" class="space-y-1">
              <button 
                @click="toggleSubmenu(sub.key)"
                class="w-full flex items-center justify-between text-xs font-bold text-neutral-400 hover:text-white"
              >
                <span class="flex items-center">
                  <span class="mr-2">{{ sub.icon }}</span> {{ subKey }}
                </span>
                <svg class="w-3 h-3 transition-transform" :class="{ 'rotate-180': menus[sub.key] }" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path d="M19 9l-7 7-7-7"/></svg>
              </button>
              <div v-show="menus[sub.key]" class="ml-4 border-l border-white/5 pl-3 space-y-1 mt-1">
                <button v-for="item in sub.items" :key="item" class="w-full text-left py-1.5 text-[10px] font-medium text-neutral-600 hover:text-brand-royal transition-colors">
                  {{ item }}
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Ventas Submenu -->
        <div class="space-y-1">
          <button 
            @click="toggleSubmenu('Ventas')"
            class="w-full flex items-center justify-between px-4 py-3 hover:bg-white/5 rounded-2xl transition-all group"
          >
            <div class="flex items-center">
              <span class="mr-3 text-lg opacity-80 group-hover:scale-110 transition-transform text-neutral-400">🛒</span>
              <span class="font-bold text-sm text-neutral-300">Ventas</span>
            </div>
            <svg 
              class="w-4 h-4 text-neutral-600 transition-transform" 
              :class="{ 'rotate-180': menus.Ventas }"
              fill="none" stroke="currentColor" viewBox="0 0 24 24"
            ><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/></svg>
          </button>
          
          <div v-show="menus.Ventas" class="ml-9 border-l border-white/5 pl-4 space-y-1">
            <button v-for="item in ['🧾 Facturas', '🧾 Facturas Nac.', '↩️ Devoluciones', '↩️ Devoluciones Nacionales', '🚫 Anulaciones']" :key="item"
              class="w-full text-left px-3 py-2 text-xs font-semibold text-neutral-500 hover:text-white transition-colors"
            >
              {{ item }}
            </button>
          </div>
        </div>

        <div class="pt-6 px-2 pb-2 text-[10px] font-bold uppercase tracking-[0.2em] text-neutral-600">
          Informes
        </div>
        <button class="w-full flex items-center px-4 py-3 hover:bg-white/5 rounded-2xl transition-all group">
          <span class="mr-3 text-lg opacity-80 group-hover:scale-110 transition-transform text-neutral-400">📊</span>
          <span class="font-bold text-sm text-neutral-300">Dashboard General</span>
        </button>

        <div class="pt-6 px-2 pb-2 text-[10px] font-bold uppercase tracking-[0.2em] text-neutral-600">
          Sistema
        </div>
         <div class="space-y-1">
          <button 
            @click="toggleSubmenu('Configuracion')"
            class="w-full flex items-center justify-between px-4 py-3 hover:bg-white/5 rounded-2xl transition-all group"
          >
            <div class="flex items-center">
              <span class="mr-3 text-lg opacity-80 group-hover:scale-110 transition-transform text-neutral-400">⚙️</span>
              <span class="font-bold text-sm text-neutral-300">Configuracion</span>
            </div>
            <svg 
              class="w-4 h-4 text-neutral-600 transition-transform" 
              :class="{ 'rotate-180': menus.Configuracion }"
              fill="none" stroke="currentColor" viewBox="0 0 24 24"
            ><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/></svg>
          </button>
          
          <div v-show="menus.Configuracion" class="ml-9 border-l border-white/5 pl-4 space-y-1">
            <button v-for="item in ['Usuarios', 'Roles', 'Permisos', 'Empresas', 'Sucursales']" :key="item"
              class="w-full text-left px-3 py-2 text-xs font-semibold text-neutral-500 hover:text-white transition-colors"
            >
              {{ item }}
            </button>
          </div>
        </div>
        <button 
          @click="logout" 
          class="w-full flex items-center px-4 py-3 text-red-500 hover:bg-red-500/10 rounded-2xl transition-all group"
        >
          <span class="mr-3 text-lg opacity-80 group-hover:scale-110 transition-transform">🚪</span>
          <span class="font-bold text-sm">Cerrar Sesión</span>
        </button>
      </nav>
    </aside>

    <!-- Main Content Area -->
    <div class="flex-1 flex flex-col h-screen overflow-hidden">
      <!-- Header -->
      <header class="h-20 bg-brand-navy-base/40 backdrop-blur-md border-b border-white/5 flex items-center justify-between px-6 lg:px-10 flex-shrink-0">
        <div class="flex items-center space-x-6">
          <button class="lg:hidden p-2 text-neutral-400 hover:bg-white/5 rounded-xl transition-colors" @click="toggleMobileMenu">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16m-7 6h7"/></svg>
          </button>
          
          <div class="hidden md:flex items-center bg-white/5 border border-white/10 rounded-2xl px-4 py-2 w-96 group focus-within:ring-2 focus-within:ring-brand-royal/30 transition-all duration-300">
            <svg class="w-4 h-4 text-neutral-500 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
            <input type="text" placeholder="Buscar reporte o transacción..." class="bg-transparent border-none focus:ring-0 text-sm text-white placeholder:text-neutral-600 w-full" />
          </div>
        </div>

        <!-- User Profile -->
        <div class="flex items-center space-x-4">
          <div class="text-right hidden sm:block">
            <div class="text-xs font-black text-white uppercase tracking-tighter">{{ nombreUsuario }}</div>
            <div class="text-[10px] font-bold text-brand-royal uppercase tracking-widest">{{ rolUsuario }}</div>
          </div>
          <div class="w-10 h-10 rounded-2xl bg-gradient-to-tr from-brand-royal to-indigo-600 flex items-center justify-center text-xs font-black text-white shadow-lg shadow-brand-royal/20 ring-2 ring-white/10">
            {{ iniciales }}
          </div>
        </div>
      </header>

      <!-- Main Scrollable Content -->
      <main class="flex-1 overflow-y-auto p-6 lg:p-10 bg-[radial-gradient(ellipse_at_bottom_left,_var(--tw-gradient-stops))] from-brand-navy-base/20 via-brand-navy-dark to-brand-navy-dark">
        <div class="max-w-7xl mx-auto space-y-10">
          
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

          <!-- Feature Placeholder -->
          <div class="grid grid-cols-1 lg:grid-cols-3 gap-8 pb-10 text-center sm:text-left">
            <div class="lg:col-span-2 p-8 rounded-3xl bg-white/5 border border-white/10 flex flex-col justify-center items-center h-80 space-y-6">
              <div class="p-4 rounded-full bg-brand-royal/10 text-brand-royal">
                <svg class="w-12 h-12" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/></svg>
              </div>
              <div class="space-y-2">
                <h3 class="text-xl font-bold text-white tracking-tight">Listo para Datos en Tiempo Real</h3>
                <p class="text-sm text-neutral-500 max-w-sm">Conecta tus endpoints de backend para visualizar flujos de caja y analíticas de inventario aquí.</p>
              </div>
            </div>
            
            <div class="p-8 rounded-3xl bg-white/5 border border-white/10 flex flex-col justify-center items-center h-80 space-y-6">
              <div class="p-4 rounded-full bg-amber-500/10 text-amber-400 italic font-black text-2xl">?</div>
              <div class="space-y-2">
                <h3 class="text-xl font-bold text-white tracking-tight">Ayuda y Soporte</h3>
                <p class="text-sm text-neutral-500">¿Necesitas ayuda con el nuevo diseño? Consulta la guía rápida.</p>
                <button class="px-6 py-3 bg-white/10 hover:bg-white/20 rounded-2xl text-xs font-bold uppercase transition-all">Ver Guía</button>
              </div>
            </div>
          </div>

        </div>
      </main>
    </div>
  </div>
</template>

<style>
/* Utilities for transitions and layout */
.scrollbar-none::-webkit-scrollbar {
  display: none;
}
.scrollbar-none {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>  