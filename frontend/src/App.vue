<template>
  <div class="app-main">
    <LoginComponent v-if="!isAuthenticated" @login-success="entrarAlSistema" />
    <DashboardComponent v-else @logout="salirDelSistema" />
  </div>
</template>

<script setup>
import { ref } from 'vue';
import LoginComponent from './components/Login.vue'; // Verifica que la ruta sea correcta
import DashboardComponent from './components/Dashboard.vue';

const isAuthenticated = ref(localStorage.getItem('isLogged') === 'true');

const entrarAlSistema = () => {
  // Cuando el Login emite 'login-success', esta función se activa
  isAuthenticated.value = true;
};

const salirDelSistema = () => {
  isAuthenticated.value = false;
  localStorage.removeItem('isLogged');
};
</script>