<template>
  <div class="h-full flex flex-col bg-dark-900 overflow-hidden">
    <!-- Header -->
    <div class="px-6 py-4 border-b border-slate-800 flex items-center justify-between bg-dark-900/60">
      <div class="flex items-center gap-3">
        <div class="p-2 bg-emerald-500/10 border border-emerald-500/20 rounded-lg text-emerald-400">
          🗂️
        </div>
        <div>
          <div class="flex items-center gap-2">
            <h2 class="text-base font-bold text-slate-100 font-mono">
              {{ store.selectedTable || 'Select a Table' }}
            </h2>
            <span v-if="store.selectedTable" class="px-2 py-0.5 bg-slate-800 text-slate-400 text-xs rounded border border-slate-700">
              SCHEMA
            </span>
          </div>
          <p class="text-xs text-slate-400">Structure, Columns, Indexes, and Constraints</p>
        </div>
      </div>

      <div v-if="store.selectedTable" class="flex items-center gap-2">
        <button
          @click="store.activeTab = 'data'"
          class="px-3.5 py-1.5 bg-emerald-600 hover:bg-emerald-500 text-white rounded-md text-xs font-semibold flex items-center gap-1.5 transition shadow"
        >
          <span>📋</span> Browse Table Data
        </button>
        <button
          @click="loadDetails"
          class="p-1.5 text-slate-400 hover:text-white bg-dark-800 border border-slate-700 rounded-md transition"
          title="Reload structure"
        >
          🔄
        </button>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="!store.selectedTable" class="flex-1 flex flex-col items-center justify-center text-slate-500 p-8 text-center">
      <div class="text-4xl mb-3">🔍</div>
      <h3 class="text-base font-semibold text-slate-300 mb-1">No table selected</h3>
      <p class="text-xs max-w-sm">Pick a table from the sidebar navigation to inspect its columns, primary keys, data types, and index constraints.</p>
    </div>

    <!-- Main Details View -->
    <div v-else class="flex-1 overflow-y-auto p-6 space-y-6">
      <div v-if="loading" class="text-center py-12 text-slate-500 text-sm">
        Loading schema metadata...
      </div>

      <div v-else class="space-y-6">
        <!-- Columns Section -->
        <div class="bg-dark-800 border border-slate-700/80 rounded-xl overflow-hidden shadow">
          <div class="px-4 py-3 border-b border-slate-700/80 bg-dark-900/40 flex items-center justify-between">
            <h3 class="text-xs font-bold text-slate-200 uppercase tracking-wider flex items-center gap-2">
              <span>Columns ({{ columns.length }})</span>
            </h3>
          </div>

          <div class="overflow-x-auto">
            <table class="w-full text-left text-xs">
              <thead class="bg-dark-900/80 text-slate-400 uppercase text-[10px] tracking-wider font-mono">
                <tr>
                  <th class="py-2.5 px-4">#</th>
                  <th class="py-2.5 px-4">Name</th>
                  <th class="py-2.5 px-4">Type</th>
                  <th class="py-2.5 px-4">Primary Key</th>
                  <th class="py-2.5 px-4">Nullable</th>
                  <th class="py-2.5 px-4">Default</th>
                  <th class="py-2.5 px-4">Extra</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-800 text-slate-300 font-mono">
                <tr v-for="col in columns" :key="col.name" class="hover:bg-dark-700/50 transition">
                  <td class="py-2 px-4 text-slate-500">{{ col.position }}</td>
                  <td class="py-2 px-4 font-semibold text-slate-100 flex items-center gap-2">
                    <span v-if="col.isPrimaryKey" class="text-amber-400" title="Primary Key">🔑</span>
                    {{ col.name }}
                  </td>
                  <td class="py-2 px-4 text-emerald-400">{{ col.dataType }}</td>
                  <td class="py-2 px-4">
                    <span v-if="col.isPrimaryKey" class="px-2 py-0.5 bg-amber-500/20 text-amber-300 border border-amber-500/30 rounded text-[10px] uppercase font-bold">
                      PK
                    </span>
                    <span v-else class="text-slate-600">-</span>
                  </td>
                  <td class="py-2 px-4">
                    <span :class="col.isNullable ? 'text-amber-400/80' : 'text-slate-500'">
                      {{ col.isNullable ? 'YES' : 'NO' }}
                    </span>
                  </td>
                  <td class="py-2 px-4 text-slate-400">{{ col.defaultValue || '—' }}</td>
                  <td class="py-2 px-4 text-slate-500">{{ col.extra || '—' }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Indexes & Foreign Keys Grid -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <!-- Indexes -->
          <div class="bg-dark-800 border border-slate-700/80 rounded-xl overflow-hidden shadow">
            <div class="px-4 py-3 border-b border-slate-700/80 bg-dark-900/40">
              <h3 class="text-xs font-bold text-slate-200 uppercase tracking-wider">
                Indexes ({{ indexes.length }})
              </h3>
            </div>
            <div class="p-3">
              <div v-if="indexes.length === 0" class="text-center py-4 text-slate-500 text-xs italic">
                No indexes defined
              </div>
              <div v-else class="space-y-2">
                <div
                  v-for="idx in indexes"
                  :key="idx.name"
                  class="p-2.5 bg-dark-900/80 border border-slate-800 rounded-lg flex items-center justify-between text-xs"
                >
                  <div>
                    <div class="font-semibold text-slate-200 font-mono">{{ idx.name }}</div>
                    <div class="text-[11px] text-slate-400 font-mono mt-0.5">
                      Columns: ({{ (idx.columns || []).join(', ') }})
                    </div>
                  </div>
                  <div class="flex items-center gap-1.5 text-[10px] font-bold">
                    <span v-if="idx.isPrimary" class="px-1.5 py-0.5 bg-amber-500/20 text-amber-300 rounded border border-amber-500/30">
                      PRIMARY
                    </span>
                    <span v-if="idx.isUnique && !idx.isPrimary" class="px-1.5 py-0.5 bg-blue-500/20 text-blue-300 rounded border border-blue-500/30">
                      UNIQUE
                    </span>
                    <span class="px-1.5 py-0.5 bg-slate-800 text-slate-400 rounded">
                      {{ idx.type || 'BTREE' }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Foreign Keys -->
          <div class="bg-dark-800 border border-slate-700/80 rounded-xl overflow-hidden shadow">
            <div class="px-4 py-3 border-b border-slate-700/80 bg-dark-900/40">
              <h3 class="text-xs font-bold text-slate-200 uppercase tracking-wider">
                Foreign Keys ({{ foreignKeys.length }})
              </h3>
            </div>
            <div class="p-3">
              <div v-if="foreignKeys.length === 0" class="text-center py-4 text-slate-500 text-xs italic">
                No foreign keys defined
              </div>
              <div v-else class="space-y-2">
                <div
                  v-for="fk in foreignKeys"
                  :key="fk.name"
                  class="p-2.5 bg-dark-900/80 border border-slate-800 rounded-lg text-xs"
                >
                  <div class="font-semibold text-slate-200 font-mono flex items-center gap-1.5">
                    <span class="text-emerald-400">🔗</span>
                    {{ fk.name || 'Foreign Key' }}
                  </div>
                  <div class="text-[11px] text-slate-300 font-mono mt-1 flex items-center gap-1">
                    <span class="text-slate-400">{{ store.selectedTable }}.{{ fk.column }}</span>
                    <span class="text-emerald-400">➔</span>
                    <span class="text-emerald-300 font-semibold">{{ fk.refTable }}.{{ fk.refColumn }}</span>
                  </div>
                  <div v-if="fk.onDelete || fk.onUpdate" class="text-[10px] text-slate-500 font-mono mt-1">
                    ON DELETE: {{ fk.onDelete || 'NO ACTION' }} | ON UPDATE: {{ fk.onUpdate || 'NO ACTION' }}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue';
import { useAppStore } from '../stores/app';
import { api } from '../api';

const store = useAppStore();
const columns = ref([]);
const indexes = ref([]);
const foreignKeys = ref([]);
const loading = ref(false);

async function loadDetails() {
  if (!store.selectedTable || !store.activeConnectionId) {
    columns.value = [];
    indexes.value = [];
    foreignKeys.value = [];
    return;
  }

  try {
    loading.value = true;
    const details = await api.getTableDetails(store.activeConnectionId, store.currentSchema, store.selectedTable);
    columns.value = details.columns || [];
    indexes.value = details.indexes || [];
    foreignKeys.value = details.foreignKeys || [];
  } catch (err) {
    store.addToast(`Failed to load table details: ${err.message}`, 'error');
  } finally {
    loading.value = false;
  }
}

watch(() => [store.selectedTable, store.activeConnectionId, store.currentSchema], () => {
  loadDetails();
});

onMounted(() => {
  if (store.selectedTable) {
    loadDetails();
  }
});
</script>
