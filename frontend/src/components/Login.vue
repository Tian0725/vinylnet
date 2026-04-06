<template>
  <div class="min-h-screen flex items-center justify-center bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-brand-navy-base/60 via-brand-navy-dark to-brand-navy-dark p-6">
    <div class="w-full max-w-md">
      <!-- Logo/Brand Section -->
      <div class="text-center mb-10 space-y-2">
        <h1 class="text-4xl font-black tracking-tighter text-white">
          VINYL<span class="text-brand-royal">NET</span>
        </h1>
        <p class="text-sm font-medium uppercase tracking-[0.2em] text-neutral-500">
          Sistema de Gestión v2.0
        </p>
      </div>

      <!-- Login Card -->
      <div class="bg-white/5 backdrop-blur-2xl border border-white/10 rounded-3xl shadow-2xl p-8 md:p-10">
        <form @submit.prevent="ejecutarLogin" class="space-y-6">
          <div class="space-y-2">
            <label class="text-[10px] font-bold uppercase tracking-widest text-neutral-400 block px-1">
              Usuario
            </label>
            <input 
              v-model="form.username" 
              type="text" 
              placeholder="admin" 
              required 
              class="w-full bg-white/5 border border-white/10 rounded-2xl px-5 py-4 text-white placeholder:text-neutral-600 focus:outline-none focus:ring-2 focus:ring-brand-royal/40 focus:border-brand-royal/50 transition-all duration-300"
            />
          </div>

          <div class="space-y-2">
            <label class="text-[10px] font-bold uppercase tracking-widest text-neutral-400 block px-1">
              Contraseña
            </label>
            <div class="relative group">
              <input 
                v-model="form.password" 
                :type="showPassword ? 'text' : 'password'" 
                placeholder="••••••••" 
                required 
                class="w-full bg-white/5 border border-white/10 rounded-2xl px-5 py-4 text-white placeholder:text-neutral-600 focus:outline-none focus:ring-2 focus:ring-brand-royal/40 focus:border-brand-royal/50 transition-all duration-300"
              />
              <button 
                type="button" 
                @click="showPassword = !showPassword" 
                class="absolute right-4 top-1/2 -translate-y-1/2 p-2 text-neutral-500 hover:text-white transition-colors"
              >
                <span v-if="showPassword">🙈</span>
                <span v-else>👁️</span>
              </button>
            </div>
          </div>

          <button 
            type="submit" 
            :disabled="loading"
            class="w-full group relative flex items-center justify-center bg-brand-royal hover:bg-blue-500 disabled:bg-brand-navy-base disabled:opacity-50 text-white font-bold py-4 rounded-2xl shadow-lg shadow-brand-royal/20 transition-all duration-300 active:scale-[0.98]"
          >
            <span v-if="loading" class="flex items-center space-x-2">
              <svg class="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <span>Ingresando...</span>
            </span>
            <span v-else>Iniciar sesión</span>
          </button>
        </form>

        <div class="mt-8 text-center border-top border-white/5 pt-8">
           <p class="text-[10px] text-neutral-600 font-medium uppercase tracking-[0.2em]">
             © 2026 VINYLNET PLATFORM
           </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';

const form = ref({
  username: '',
  password: ''
});

const showPassword = ref(false);
const loading = ref(false);

const emit = defineEmits(['login-success']);

const ejecutarLogin = async () => {
  loading.value = true;
  
  try {
    const response = await fetch('http://localhost:8080/api/login', { 
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(form.value)
    });

    if (response.ok) {
      const data = await response.json();
      
      localStorage.setItem('isLogged', 'true'); 
      localStorage.setItem('usuario', data.usuario);
      localStorage.setItem('rol', data.rol);
      
      emit('login-success'); 
    } else {
      alert("Credenciales incorrectas");
    }
  } catch (error) {
    console.error("Error en el login:", error);
    alert("No se pudo conectar con el servidor");
  } finally {
    loading.value = false;
  }
}
</script>
