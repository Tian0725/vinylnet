<template>
  <div class="container">
    <div class="card">
      <h1 class="title">VINYL<span class="blue">NET</span></h1>
      <p class="subtitle">Sistema de Gestión</p>

      <form @submit.prevent="ejecutarLogin" class="form">
        <div class="group">
          <label>Usuario</label>
          <input v-model="form.username" type="text" placeholder="Tu usuario" required />
        </div>

        <div class="group">
          <label>Contraseña</label>
          <div class="password-box">
            <input v-model="form.password" :type="showPassword ? 'text' : 'password'" placeholder="••••••••" required />
            <button type="button" @click="showPassword = !showPassword" class="eye">👁️</button>
          </div>
        </div>

        <button type="submit" class="btn">{{ loading ? 'Ingresando...' : 'Iniciar sesión' }}</button>
      </form>

      <div class="footer">© 2026 VINYLNET</div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';

// 1. Definimos el objeto form (ESTO CORRIGE EL TYPEERROR)
const form = ref({
  username: '',
  password: ''
});

// Definimos las otras variables que usas en el template
const showPassword = ref(false);
const loading = ref(false);

const emit = defineEmits(['login-success']);

const ejecutarLogin = async () => {
  loading.value = true;
  
  try {
    // 2. DEBES DEFINIR 'response' (ESTO CORRIGE EL REFERENCEERROR)
    // Reemplaza la URL con la de tu API real
const response = await fetch('http://192.168.1.7:8080/api/login', { 
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
    // Aquí es donde caía el error antes
    console.error("Error en el login:", error);
    alert("No se pudo conectar con el servidor");
  } finally {
    loading.value = false;
  }
}
</script>

<style scoped>
.container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #0f172a, #1e293b, #020617);
  padding: 20px;
}

.card {
  width: 100%;
  max-width: 400px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 15px 40px rgba(0,0,0,0.4);
  overflow: hidden;
  text-align: center;
}

.title {
  font-size: 28px;
  font-weight: 800;
  color: #0f172a;
}

.blue { color: #2563eb; }

.subtitle {
  font-size: 12px;
  color: #64748b;
  margin-top: 5px;
  letter-spacing: 1px;
}

.form {
  padding: 20px 30px 30px;
}

.group {
  margin-bottom: 20px;
  text-align: left;
}

.group label {
  font-size: 11px;
  font-weight: bold;
  color: #475569;
  display: block;
  margin-bottom: 5px;
}

input {
  width: 100%;
  padding: 12px;
  border-radius: 10px;
  border: 1px solid #cbd5f5;
  background: #f1f5f9;
  transition: 0.3s;
}

input:focus {
  outline: none;
  border-color: #2563eb;
  background: white;
}

.password-box { position: relative; }

.eye {
  position: absolute;
  right: 10px;
  top: 8px;
  background: none;
  border: none;
  cursor: pointer;
}

.btn {
  width: 100%;
  padding: 12px;
  background: #2563eb;
  color: white;
  border: none;
  border-radius: 10px;
  font-weight: bold;
  cursor: pointer;
  transition: 0.2s;
}

.btn:hover { background: #1d4ed8; }

.btn:disabled { opacity: 0.6; cursor: not-allowed; }

.footer {
  text-align: center;
  font-size: 11px;
  color: #94a3b8;
  padding-bottom: 20px;
}
</style>