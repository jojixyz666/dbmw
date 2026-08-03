<template>
  <div class="fixed inset-0 z-40 bg-black/75 backdrop-blur-md flex items-center justify-center p-4">
    <div class="glass-panel w-full max-w-2xl rounded-xl z-10 flex flex-col overflow-hidden border border-outline-variant/50 shadow-2xl max-h-[90vh]">
      <!-- Modal Header -->
      <div class="px-6 py-4 border-b border-outline-variant flex justify-between items-center bg-surface">
        <div class="flex items-center gap-3">
          <div class="w-9 h-9 rounded-lg bg-primary/10 border border-primary/30 flex items-center justify-center text-primary">
            <span class="material-symbols-outlined text-xl">database</span>
          </div>
          <div>
            <h2 class="font-bold text-on-surface text-sm m-0">Database Connection Profiles</h2>
            <p class="text-xs text-on-surface-variant m-0">Manage local or remote database instances</p>
          </div>
        </div>
        <button
          @click="store.isConnectionsModalOpen = false"
          class="text-on-surface-variant hover:text-on-surface transition-colors rounded-full p-1 hover:bg-surface-container-low"
        >
          <span class="material-symbols-outlined text-xl">close</span>
        </button>
      </div>

      <!-- Modal Body -->
      <div class="p-6 overflow-y-auto bg-surface-container-lowest flex-1 space-y-6">
        <!-- New / Edit Form Panel -->
        <div class="border border-outline-variant rounded-lg p-5 bg-surface space-y-4 shadow-sm">
          <div class="flex items-center justify-between">
            <h3 class="text-xs font-bold text-primary uppercase tracking-wider font-mono m-0">
              {{ form.id ? 'Edit Connection Profile' : 'New Connection Profile' }}
            </h3>
            <button v-if="form.id" @click="resetForm" class="text-xs text-on-surface-variant hover:text-on-surface underline font-mono">
              Cancel Edit
            </button>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4">
            <div class="space-y-1.5">
              <label class="font-mono text-xs text-on-surface-variant block font-semibold">Name</label>
              <input
                v-model="form.name"
                placeholder="e.g. Local Dev Database"
                class="w-full h-9 px-3 rounded border border-outline-variant bg-surface-container-lowest text-on-surface font-mono text-xs focus:border-primary focus:ring-1 focus:ring-primary outline-none transition-colors placeholder:text-on-surface-variant/50"
              />
            </div>
            <div class="space-y-1.5">
              <label class="font-mono text-xs text-on-surface-variant block font-semibold">Driver</label>
              <select
                v-model="form.driver"
                @change="onDriverChange"
                class="w-full h-9 px-3 rounded border border-outline-variant bg-surface-container-lowest text-on-surface font-mono text-xs focus:border-primary focus:ring-1 focus:ring-primary outline-none transition-colors"
              >
                <option value="postgres">PostgreSQL</option>
                <option value="mysql">MySQL</option>
                <option value="mariadb">MariaDB</option>
                <option value="sqlite">SQLite (Pure Go)</option>
              </select>
            </div>
          </div>

          <!-- SQLite Specific -->
          <div v-if="form.driver === 'sqlite'" class="space-y-1.5">
            <label class="font-mono text-xs text-on-surface-variant block font-semibold">SQLite File Path</label>
            <input
              v-model="form.filePath"
              placeholder="./database.sqlite or :memory:"
              class="w-full h-9 px-3 rounded border border-outline-variant bg-surface-container-lowest text-on-surface font-mono text-xs focus:border-primary focus:ring-1 focus:ring-primary outline-none transition-colors placeholder:text-on-surface-variant/50"
            />
          </div>

          <!-- Network DB Specific -->
          <div v-else class="space-y-4">
            <div class="grid grid-cols-3 gap-4">
              <div class="col-span-2 space-y-1.5">
                <label class="font-mono text-xs text-on-surface-variant block font-semibold">Host</label>
                <input
                  v-model="form.host"
                  placeholder="127.0.0.1 or localhost"
                  class="w-full h-9 px-3 rounded border border-outline-variant bg-surface-container-lowest text-on-surface font-mono text-xs focus:border-primary focus:ring-1 focus:ring-primary outline-none transition-colors placeholder:text-on-surface-variant/50"
                />
              </div>
              <div class="space-y-1.5">
                <label class="font-mono text-xs text-on-surface-variant block font-semibold">Port</label>
                <input
                  v-model.number="form.port"
                  type="number"
                  placeholder="5432"
                  class="w-full h-9 px-3 rounded border border-outline-variant bg-surface-container-lowest text-on-surface font-mono text-xs focus:border-primary focus:ring-1 focus:ring-primary outline-none transition-colors"
                />
              </div>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div class="space-y-1.5">
                <label class="font-mono text-xs text-on-surface-variant block font-semibold">Username</label>
                <input
                  v-model="form.user"
                  placeholder="postgres / root"
                  class="w-full h-9 px-3 rounded border border-outline-variant bg-surface-container-lowest text-on-surface font-mono text-xs focus:border-primary focus:ring-1 focus:ring-primary outline-none transition-colors placeholder:text-on-surface-variant/50"
                />
              </div>
              <div class="space-y-1.5">
                <label class="font-mono text-xs text-on-surface-variant block font-semibold">Password</label>
                <input
                  v-model="form.password"
                  type="password"
                  placeholder="••••••••"
                  class="w-full h-9 px-3 rounded border border-outline-variant bg-surface-container-lowest text-on-surface font-mono text-xs focus:border-primary focus:ring-1 focus:ring-primary outline-none transition-colors"
                />
              </div>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div class="space-y-1.5">
                <label class="font-mono text-xs text-on-surface-variant block font-semibold">Database Name</label>
                <input
                  v-model="form.database"
                  placeholder="postgres / my_app_dev"
                  class="w-full h-9 px-3 rounded border border-outline-variant bg-surface-container-lowest text-on-surface font-mono text-xs focus:border-primary focus:ring-1 focus:ring-primary outline-none transition-colors placeholder:text-on-surface-variant/50"
                />
              </div>
              <div v-if="form.driver === 'postgres'" class="space-y-1.5">
                <label class="font-mono text-xs text-on-surface-variant block font-semibold">SSL Mode</label>
                <select
                  v-model="form.sslMode"
                  class="w-full h-9 px-3 rounded border border-outline-variant bg-surface-container-lowest text-on-surface font-mono text-xs focus:border-primary focus:ring-1 focus:ring-primary outline-none transition-colors"
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
              class="h-9 px-4 rounded bg-surface-container hover:bg-surface-container-high text-on-surface border border-outline-variant transition-colors text-xs font-semibold font-mono"
            >
              {{ isTesting ? 'Testing...' : 'Test Connection' }}
            </button>
            <button
              type="button"
              @click="saveConn"
              class="h-9 px-4 rounded bg-primary hover:bg-primary-container text-on-primary font-bold transition-colors text-xs shadow-sm font-mono active:scale-95"
            >
              {{ form.id ? 'Update Profile' : 'Save Connection' }}
            </button>
          </div>
        </div>

        <!-- Saved Connections List -->
        <div>
          <h4 class="text-xs font-bold text-on-surface-variant uppercase tracking-wider mb-3 font-mono">Saved Connections</h4>
          <div v-if="store.connections.length === 0" class="text-center py-6 text-on-surface-variant text-xs font-mono">
            No connection profiles saved yet. Fill out the form above to add one.
          </div>
          <div v-else class="space-y-2">
            <div
              v-for="c in store.connections"
              :key="c.id"
              class="p-3.5 bg-surface border border-outline-variant rounded-lg flex items-center justify-between hover:border-primary transition-colors shadow-sm group"
            >
              <div class="flex items-center gap-3">
                <span
                  class="px-2.5 py-0.5 text-[10px] font-bold rounded uppercase tracking-wider font-mono border"
                  :class="c.driver === 'postgres' ? 'bg-blue-500/20 text-blue-400 border-blue-500/30' : c.driver === 'mysql' || c.driver === 'mariadb' ? 'bg-amber-500/20 text-amber-400 border-amber-500/30' : 'bg-primary/20 text-primary border-primary/30'"
                >
                  {{ c.driver }}
                </span>
                <div>
                  <h5 class="text-xs font-bold text-on-surface m-0">{{ c.name }}</h5>
                  <p class="text-[11px] text-on-surface-variant font-mono mt-0.5 m-0">
                    {{ c.driver === 'sqlite' ? c.filePath : `${c.user || 'root'}@${c.host}:${c.port}/${c.database}` }}
                  </p>
                </div>
              </div>

              <div class="flex items-center gap-2">
                <button
                  @click="store.switchConnection(c.id)"
                  class="px-3 py-1 bg-primary/20 hover:bg-primary text-primary hover:text-on-primary rounded text-xs font-bold border border-primary/40 transition-colors font-mono"
                  :class="store.activeConnectionId === c.id ? 'bg-primary text-on-primary font-extrabold' : ''"
                >
                  {{ store.activeConnectionId === c.id ? 'Active' : 'Connect' }}
                </button>
                <button
                  @click="editConn(c)"
                  class="p-1.5 text-on-surface-variant hover:text-on-surface rounded hover:bg-surface-container transition-colors"
                  title="Edit"
                >
                  <span class="material-symbols-outlined text-base">edit</span>
                </button>
                <button
                  @click="deleteConn(c.id)"
                  class="p-1.5 text-on-surface-variant hover:text-red-400 rounded hover:bg-red-950/30 transition-colors"
                  title="Delete"
                >
                  <span class="material-symbols-outlined text-base">delete</span>
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
import { ref } from 'vue';
import { useAppStore } from '../stores/app';
import { api } from '../api';

const store = useAppStore();

const form = ref({
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

const isTesting = ref(false);

function onDriverChange() {
  if (form.value.driver === 'postgres') form.value.port = 5432;
  else if (form.value.driver === 'mysql' || form.value.driver === 'mariadb') form.value.port = 3306;
}

function resetForm() {
  form.value = {
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
  };
}

function editConn(c) {
  form.value = { ...c };
}

async function testConn() {
  isTesting.value = true;
  try {
    const res = await api.testConnection(form.value);
    if (res.success) {
      store.addToast('Connection test succeeded!', 'success');
    } else {
      store.addToast(`Connection failed: ${res.message}`, 'error');
    }
  } catch (err) {
    store.addToast(`Test error: ${err.message}`, 'error');
  } finally {
    isTesting.value = false;
  }
}

async function saveConn() {
  if (!form.value.name) {
    store.addToast('Please enter a profile name', 'error');
    return;
  }
  try {
    if (form.value.id) {
      await api.updateConnection(form.value.id, form.value);
      store.addToast('Connection updated', 'success');
    } else {
      await api.createConnection(form.value);
      store.addToast('Connection saved', 'success');
    }
    await store.loadConnections();
    resetForm();
  } catch (err) {
    store.addToast(err.message, 'error');
  }
}

async function deleteConn(id) {
  if (!confirm('Are you sure you want to delete this profile?')) return;
  try {
    await api.deleteConnection(id);
    store.addToast('Connection profile deleted', 'success');
    await store.loadConnections();
  } catch (err) {
    store.addToast(err.message, 'error');
  }
}
</script>