<script setup>
import { ref, onMounted } from 'vue';

const searchQuery = ref('');
const filterStatus = ref('all');
const inventory = ref([]);

const fetchInventory = async () => {
  try {
    const response = await fetch('http://localhost:8080/api/productos');
    if (response.ok) {
      inventory.value = await response.json();
    } else {
      console.error('Error al cargar inventario');
    }
  } catch (error) {
    console.error('Error de conexión con el backend:', error);
  }
};

onMounted(() => {
  fetchInventory();
});

const getStatusClass = (status) => {
  switch (status) {
    case 'Disponible': return 'bg-emerald-500/10 text-emerald-500 border-emerald-500/20';
    case 'Stock Bajo': return 'bg-amber-500/10 text-amber-500 border-amber-500/20';
    case 'Agotado': return 'bg-red-500/10 text-red-500 border-red-500/20';
    default: return 'bg-neutral-500/10 text-neutral-500 border-neutral-500/20';
  }
};
</script>

<template>
  <div class="space-y-8 animate-in fade-in duration-500">
    <!-- Header de la Vista -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div class="space-y-1">
        <h1 class="text-3xl font-black text-white tracking-tighter uppercase italic">Control de Inventario</h1>
        <p class="text-sm font-medium text-neutral-500 uppercase tracking-widest">Gestión integral de productos y existencias</p>
      </div>
      
      <button class="flex items-center justify-center space-x-2 bg-brand-royal hover:bg-blue-500 text-white font-bold py-3 px-6 rounded-2xl shadow-lg shadow-brand-royal/20 transition-all active:scale-95 text-sm uppercase">
        <span>+ Nuevo Producto</span>
      </button>
    </div>

    <!-- Barra de Herramientas (Filtros y Búsqueda) -->
    <div class="p-6 rounded-3xl bg-white/5 border border-white/10 backdrop-blur-md flex flex-col lg:flex-row gap-4">
      <div class="flex-1 relative group">
        <svg class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-neutral-500 group-focus-within:text-brand-royal transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
        <input 
          v-model="searchQuery"
          type="text" 
          placeholder="Buscar por descripción o código..." 
          class="w-full bg-white/5 border border-white/10 rounded-2xl pl-12 pr-6 py-3 text-sm text-white placeholder:text-neutral-600 focus:outline-none focus:ring-2 focus:ring-brand-royal/30 transition-all"
        />
      </div>
      
      <div class="flex gap-4">
        <select v-model="filterStatus" class="bg-white/5 border border-white/10 rounded-2xl px-6 py-3 text-sm text-neutral-300 focus:outline-none focus:ring-2 focus:ring-brand-royal/30 transition-all cursor-pointer">
          <option value="all" class="bg-brand-navy-dark">Todos los Estados</option>
          <option value="disponible" class="bg-brand-navy-dark">Disponible</option>
          <option value="bajo" class="bg-brand-navy-dark">Bajo Stock</option>
          <option value="agotado" class="bg-brand-navy-dark">Agotado</option>
        </select>
        
        <button class="p-3 bg-white/5 hover:bg-white/10 border border-white/10 rounded-2xl text-neutral-400 hover:text-white transition-all">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z"/></svg>
        </button>
      </div>
    </div>

    <!-- Tabla de Inventario -->
    <div class="rounded-3xl bg-white/5 border border-white/10 overflow-hidden backdrop-blur-md shadow-2xl">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-white/5 border-b border-white/10 text-[10px] font-black text-neutral-500 uppercase tracking-[0.2em]">
              <th class="px-8 py-5">Código</th>
              <th class="px-8 py-5">Descripción</th>
              <th class="px-8 py-5 text-center">Existencia</th>
              <th class="px-8 py-5">Precio</th>
              <th class="px-8 py-5">Estado</th>
              <th class="px-8 py-5 text-right">Acciones</th>
            </tr>
          </thead>
          <tbody class="text-sm font-medium">
            <tr v-for="item in inventory" :key="item.id" class="group border-b border-white/5 hover:bg-white/[0.02] transition-colors">
              <td class="px-8 py-4 font-mono text-xs text-brand-royal font-bold">{{ String(item.id).padStart(5, '0') }}</td>
              <td class="px-8 py-4">
                <div class="flex flex-col">
                  <span class="text-white font-bold">{{ item.nombre }}</span>
                  <span class="text-[10px] text-neutral-600 uppercase font-black">{{ item.unidad || 'Unidad' }}</span>
                </div>
              </td>
              <td class="px-8 py-4">
                <div class="flex flex-col items-center">
                  <span :class="item.stock <= item.stockMin ? 'text-amber-500' : 'text-neutral-300'" class="font-bold font-mono">
                    {{ item.stock.toFixed(2) }}
                  </span>
                  <div class="w-16 h-1 bg-white/5 rounded-full mt-2 overflow-hidden">
                    <div 
                      class="h-full rounded-full transition-all duration-500"
                      :class="item.stock <= 0 ? 'bg-red-500' : (item.stock <= item.stockMin ? 'bg-amber-500' : 'bg-emerald-500')"
                      :style="{ width: Math.min((item.stock / (item.stockMin * 2)) * 100, 100) + '%' }"
                    ></div>
                  </div>
                </div>
              </td>
              <td class="px-8 py-4 text-white font-bold font-mono">${{ item.precio.toFixed(2) }}</td>
              <td class="px-8 py-4">
                <span :class="getStatusClass(item.status)" class="px-3 py-1 rounded-full text-[10px] font-black uppercase tracking-wider border">
                  {{ item.status }}
                </span>
              </td>
              <td class="px-8 py-4 text-right">
                <div class="flex items-center justify-end space-x-2 opacity-0 group-hover:opacity-100 transition-opacity">
                  <button class="p-2 hover:bg-brand-royal/20 text-neutral-400 hover:text-brand-royal rounded-xl transition-all" title="Editar">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/></svg>
                  </button>
                  <button class="p-2 hover:bg-red-500/20 text-neutral-400 hover:text-red-500 rounded-xl transition-all" title="Eliminar">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
