<template>
  <div class="fixed inset-0 z-40 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4">
    <div class="bg-dark-800 border border-slate-700/80 rounded-xl w-full max-w-2xl shadow-2xl flex flex-col max-h-[90vh] overflow-hidden">
      <!-- Header -->
      <div class="px-6 py-4 border-b border-slate-700/80 flex items-center justify-between bg-dark-900/50">
        <div class="flex items-center gap-2">
          <div class="p-2 bg-emerald-500/10 rounded-lg text-emerald-400">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/></svg>
          </div>
          <div>
            <h3 class="font-semibold text-slate-100">Database Connection Profiles</h3>
            <p class="text-xs text-slate-400">Manage local or remote database instances</p>
          </div>
        </div>
        <button @click="store.isConnectionsModalOpen = false" class="text-slate-400 hover:text-white p-1 rounded">✕</button>
      </div>

      <!-- Body -->
      <div class="p-6 flex-1 overflow-y-auto space-y-6">
        <!-- New / Edit Form -->
        <div class="bg-dark-900/80 p-5 rounded-lg border border-slate-700/60 space-y-4">
          <div class="flex items-center justify-between">
            <h4 class="text-sm font-semibold text-emerald-400 uppercase tracking-wider">
              {{ form.id ? 'Edit Connection' : 'New Connection Profile' }}
            </h4>
            <button v-if="form.id" @click="resetForm" class="text-xs text-slate-400 hover:text-white underline">
              Cancel Edit
            </button>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div class="col-span-2 sm:col-span-1">
              <label class="block text-xs font-medium text-slate-300 mb-1">Name</label>
              <input
                v-model="form.name"
                placeholder="e.g. Local Dev Database"
                class="w-full bg-dark-800 border border-slate-700 rounded-md px-3 py-2 text-sm text-slate-100 placeholder-slate-500 focus:outline-none focus:border-emerald-500"
              />
            </div>
            <div class="col-span-2 sm:col-span-1">
              <label class="block text-xs font-medium text-slate-300 mb-1">Driver</label>
              <select
                v-model="form.driver"
                @change="onDriverChange"
                class="w-full bg-dark-800 border border-slate-700 rounded-md px-3 py-2 text-sm text-slate-100 focus:outline-none focus:border-emerald-500"
              >
                <option value="postgres">PostgreSQL</option>
                <option value="mysql">MySQL</option>
                <option value="mariadb">MariaDB</option>
                <option value="sqlite">SQLite (Pure Go)</option>
              </select>
            </div>
          </div>

          <!-- SQLite Specific -->
          <div v-if="form.driver === 'sqlite'" class="space-y-2">
            <label class="block text-xs font-medium text-slate-300">SQLite File Path</label>
            <input
              v-model="form.filePath"
              placeholder="./database.sqlite or :memory:"
              class="w-full bg-dark-800 border border-slate-700 rounded-md px-3 py-2 text-sm text-slate-100 placeholder-slate-500 focus:outline-none focus:border-emerald-500 font-mono"
            />
          </div>

          <!-- Network DB Specific -->
          <div v-else class="space-y-4">
            <div class="grid grid-cols-3 gap-4">
              <div class="col-span-2">
                <label class="block text-xs font-medium text-slate-300 mb-1">Host</label>
                <input
                  v-model="form.host"
                  placeholder="127.0.0.1 or localhost"
                  class="w-full bg-dark-800 border border-slate-700 rounded-md px-3 py-2 text-sm text-slate-100 placeholder-slate-500 focus:outline-none focus:border-emerald-500 font-mono"
                />
              </div>
              <div>
                <label class="block text-xs font-medium text-slate-300 mb-1">Port</label>
                <input
                  v-model.number="form.port"
                  type="number"
                  placeholder="5432"
                  class="w-full bg-dark-800 border border-slate-700 rounded-md px-3 py-2 text-sm text-slate-100 placeholder-slate-500 focus:outline-none focus:border-emerald-500 font-mono"
                />
              </div>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-xs font-medium text-slate-300 mb-1">Username</label>
                <input
                  v-model="form.user"
                  placeholder="postgres / root"
                  class="w-full bg-dark-800 border border-slate-700 rounded-md px-3 py-2 text-sm text-slate-100 placeholder-slate-500 focus:outline-none focus:border-emerald-500 font-mono"
                />
              </div>
              <div>
                <label class="block text-xs font-medium text-slate-300 mb-1">Password</label>
                <input
                  v-model="form.password"
                  type="password"
                  placeholder="••••••••"
                  class="w-full bg-dark-800 border border-slate-700 rounded-md px-3 py-2 text-sm text-slate-100 placeholder-slate-500 focus:outline-none focus:border-emerald-500 font-mono"
                />
              </div>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-xs font-medium text-slate-300 mb-1">Database Name</label>
                <input
                  v-model="form.database"
                  placeholder="postgres / my_app_dev"
                  class="w-full bg-dark-800 border border-slate-700 rounded-md px-3 py-2 text-sm text-slate-100 placeholder-slate-500 focus:outline-none focus:border-emerald-500 font-mono"
                />
              </div>
              <div v-if="form.driver === 'postgres'">
                <label class="block text-xs font-medium text-slate-300 mb-1">SSL Mode</label>
                <select
                  v-model="form.sslMode"
                  class="w-full bg-dark-800 border border-slate-700 rounded-md px-3 py-2 text-sm text-slate-100 focus:outline-none focus:border-emerald-500"
                >
                  <option value="disable">disable</option>
                  <option value="require">require</option>
                  <option value="verify-full">verify-full</option>
                </select>
              </div>
            </div>
          </div>

          <!-- Form Actions -->
          <div class="flex items-center justify-end gap-3 pt-2">
            <button
              type="button"
              @click="testConn"
              :disabled="isTesting"
              class="px-4 py-2 bg-dark-700 hover:bg-dark-600 text-slate-200 rounded-md text-xs font-semibold border border-slate-600 transition"
            >
              {{ isTesting ? 'Testing...' : 'Test Connection' }}
            </button>
            <button
              type="button"
              @click="saveConn"
              class="px-4 py-2 bg-emerald-600 hover:bg-emerald-500 text-white rounded-md text-xs font-semibold transition"
            >
              {{ form.id ? 'Update Profile' : 'Save Connection' }}
            </button>
          </div>
        </div>

        <!-- Saved Connections List -->
        <div>
          <h4 class="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-3">Saved Connections</h4>
          <div v-if="store.connections.length === 0" class="text-center py-6 text-slate-500 text-sm">
            No connection profiles saved yet. Fill out the form above to add one.
          </div>
          <div v-else class="space-y-2">
            <div
              v-for="c in store.connections"
              :key="c.id"
              class="p-3 bg-dark-900 border border-slate-700/60 rounded-lg flex items-center justify-between hover:border-slate-600 transition"
            >
              <div class="flex items-center gap-3">
                <span
                  class="px-2 py-0.5 text-[10px] font-bold rounded uppercase tracking-wider"
                  :class="c.driver === 'postgres' ? 'bg-blue-500/20 text-blue-400 border border-blue-500/30' : c.driver === 'mysql' || c.driver === 'mariadb' ? 'bg-amber-500/20 text-amber-400 border border-amber-500/30' : 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30'"
                >
                  {{ c.driver }}
                </span>
                <div>
                  <h5 class="text-sm font-semibold text-slate-200">{{ c.name }}</h5>
                  <p class="text-xs text-slate-500 font-mono">
                    {{ c.driver === 'sqlite' ? c.filePath : `${c.host}:${c.port} / ${c.database || 'default'}` }}
                  </p>
                </div>
              </div>
              <div class="flex items-center gap-2">
                <button
                  v-if="store.activeConnectionId !== c.id"
                  @click="activateConn(c.id)"
                  class="px-3 py-1 bg-emerald-600/80 hover:bg-emerald-500 text-white rounded text-xs transition"
                >
                  Activate
                </button>
                <span v-else class="px-2 py-1 bg-emerald-500/20 text-emerald-400 border border-emerald-500/40 rounded text-[11px] font-medium">
                  Active
                </span>
                <button @click="editConn(c)" class="p-1 text-slate-400 hover:text-white rounded">
                  ✏️
                </button>
                <button @click="deleteConn(c.id)" class="p-1 text-red-400 hover:text-red-300 rounded">
                  🗑️
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue';
import { useAppStore } from '../stores/app';
import { api } from '../api';

const store = useAppStore();
const isTesting = ref(false);

const form = reactive({
  id: '',
  name: '',
  driver: 'postgres',
  host: '127.0.0.1',
  port: 5432,
  user: 'postgres',
  password: '',
  database: 'postgres',
  filePath: '',
  sslMode: 'disable',
});

function onDriverChange() {
  if (form.driver === 'postgres') {
    form.port = 5432;
    form.user = 'postgres';
    form.database = 'postgres';
  } else if (form.driver === 'mysql' || form.driver === 'mariadb') {
    form.port = 3306;
    form.user = 'root';
    form.database = 'test';
  } else if (form.driver === 'sqlite') {
    form.filePath = './database.sqlite';
  }
}

function resetForm() {
  form.id = '';
  form.name = '';
  form.driver = 'postgres';
  form.host = '127.0.0.1';
  form.port = 5432;
  form.user = 'postgres';
  form.password = '';
  form.database = 'postgres';
  form.filePath = '';
  form.sslMode = 'disable';
}

function editConn(c) {
  form.id = c.id;
  form.name = c.name;
  form.driver = c.driver;
  form.host = c.host || '127.0.0.1';
  form.port = c.port || (c.driver === 'postgres' ? 5432 : 3306);
  form.user = c.user || '';
  form.password = c.password || '';
  form.database = c.database || '';
  form.filePath = c.filePath || '';
  form.sslMode = c.sslMode || 'disable';
}

async function testConn() {
  try {
    isTesting.value = true;
    await api.testConnection(form);
    store.addToast('Connection successful!', 'success');
  } catch (err) {
    store.addToast(`Connection failed: ${err.message}`, 'error');
  } finally {
    isTesting.value = false;
  }
}

async function saveConn() {
  try {
    if (!form.name) {
      store.addToast('Connection name is required', 'error');
      return;
    }
    const saved = await api.saveConnection(form);
    await store.loadConnections();
    if (!store.activeConnectionId) {
      await store.setActiveConnection(saved.id);
    }
    resetForm();
    store.addToast('Saved connection profile', 'success');
  } catch (err) {
    store.addToast(err.message, 'error');
  }
}

async function deleteConn(id) {
  if (!confirm('Are you sure you want to delete this connection profile?')) return;
  try {
    await api.deleteConnection(id);
    await store.loadConnections();
    store.addToast('Connection profile deleted', 'success');
  } catch (err) {
    store.addToast(err.message, 'error');
  }
}

async function activateConn(id) {
  await store.setActiveConnection(id);
  store.isConnectionsModalOpen = false;
}
</script>
