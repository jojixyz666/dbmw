<template>
  <div class="h-full flex flex-col bg-surface overflow-hidden">
    <!-- Header -->
    <header class="h-[72px] px-6 border-b border-outline-variant flex items-center justify-between bg-surface flex-shrink-0">
      <div class="flex items-center gap-4">
        <div class="w-10 h-10 rounded-lg bg-surface-container-low border border-outline-variant flex items-center justify-center text-primary shadow-sm">
          <span class="material-symbols-outlined text-2xl">schema</span>
        </div>
        <div>
          <div class="flex items-center gap-2">
            <h2 class="text-sm font-bold text-on-surface font-mono tracking-tight m-0">
              {{ store.selectedTable || 'Select a Table' }}
            </h2>
            <span v-if="store.selectedTable" class="px-2 py-0.5 bg-primary/20 text-primary text-[10px] font-bold uppercase rounded border border-primary/30">
              SCHEMA
            </span>
          </div>
          <p class="text-xs text-on-surface-variant m-0">Structure, Columns, Indexes, and Constraints</p>
        </div>
      </div>

      <div v-if="store.selectedTable" class="flex items-center gap-2">
        <button
          @click="store.activeTab = 'data'"
          class="h-9 px-3.5 bg-primary hover:bg-primary-container text-on-primary rounded text-xs font-bold flex items-center gap-1.5 transition-colors shadow-sm active:scale-95"
        >
          <span class="material-symbols-outlined text-[16px]">table_chart</span> Browse Table Data
        </button>
        <button
          @click="loadDetails"
          class="h-9 px-2.5 text-on-surface-variant hover:text-on-surface bg-surface-container-low border border-outline-variant rounded transition-colors text-xs shadow-sm flex items-center justify-center"
          title="Reload structure"
        >
          <span class="material-symbols-outlined text-[16px]">refresh</span>
        </button>
      </div>
    </header>

    <!-- Empty State -->
    <div v-if="!store.selectedTable" class="flex-1 flex flex-col items-center justify-center text-on-surface-variant p-8 text-center">
      <span class="material-symbols-outlined text-5xl mb-3 text-primary/40">search</span>
      <h3 class="text-sm font-bold text-on-surface mb-1">No Table Selected</h3>
      <p class="text-xs max-w-sm text-on-surface-variant font-mono leading-relaxed">Pick a table from the sidebar navigation to inspect its columns, primary keys, data types, and index constraints.</p>
    </div>

    <!-- Main Details View -->
    <div v-else class="flex-1 overflow-y-auto p-6 space-y-6 bg-surface-container-lowest">
      <div v-if="loading" class="text-center py-12 text-on-surface-variant text-xs font-mono gap-2 flex items-center justify-center">
        <span class="material-symbols-outlined text-primary text-xl animate-spin">sync</span> Loading schema metadata...
      </div>

      <div v-else class="space-y-6">
        <!-- Columns Section -->
        <div class="bg-surface border border-outline-variant rounded-xl overflow-hidden shadow-sm">
          <div class="px-5 py-3.5 border-b border-outline-variant bg-surface-container-low flex items-center justify-between">
            <h3 class="text-xs font-bold text-on-surface uppercase tracking-wider font-mono flex items-center gap-2 m-0">
              <span class="material-symbols-outlined text-sm text-primary">view_column</span> Columns ({{ columns.length }})
            </h3>
          </div>

          <div class="overflow-x-auto">
            <table class="w-full text-left text-xs">
              <thead class="bg-surface-container-low text-on-surface-variant uppercase text-[10px] tracking-wider font-mono border-b border-outline-variant">
                <tr>
                  <th class="py-2.5 px-4 w-12 text-center">#</th>
                  <th class="py-2.5 px-4">Name</th>
                  <th class="py-2.5 px-4">Type</th>
                  <th class="py-2.5 px-4">Primary Key</th>
                  <th class="py-2.5 px-4">Nullable</th>
                  <th class="py-2.5 px-4">Default</th>
                  <th class="py-2.5 px-4">Extra</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-outline-variant/60 text-on-surface font-mono">
                <tr v-for="col in columns" :key="col.name" class="hover:bg-surface-container/60 transition-colors">
                  <td class="py-2.5 px-4 text-on-surface-variant text-center">{{ col.position }}</td>
                  <td class="py-2.5 px-4 font-bold text-on-surface flex items-center gap-2">
                    <span v-if="col.isPrimaryKey" class="material-symbols-outlined text-amber-400 text-sm" title="Primary Key">key</span>
                    {{ col.name }}
                  </td>
                  <td class="py-2.5 px-4 text-primary font-semibold">{{ col.dataType }}</td>
                  <td class="py-2.5 px-4">
                    <span v-if="col.isPrimaryKey" class="px-2 py-0.5 bg-amber-500/20 text-amber-600 dark:text-amber-300 border border-amber-500/30 rounded text-[10px] uppercase font-bold">
                      PK
                    </span>
                    <span v-else class="text-on-surface-variant">-</span>
                  </td>
                  <td class="py-2.5 px-4">
                    <span :class="col.isNullable ? 'text-amber-500 font-bold' : 'text-on-surface-variant'">
                      {{ col.isNullable ? 'YES' : 'NO' }}
                    </span>
                  </td>
                  <td class="py-2.5 px-4 text-on-surface-variant">{{ col.defaultValue || '—' }}</td>
                  <td class="py-2.5 px-4 text-on-surface-variant">{{ col.extra || '—' }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Indexes & Foreign Keys Grid -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <!-- Indexes -->
          <div class="bg-surface border border-outline-variant rounded-xl overflow-hidden shadow-sm">
            <div class="px-5 py-3.5 border-b border-outline-variant bg-surface-container-low">
              <h3 class="text-xs font-bold text-on-surface uppercase tracking-wider font-mono m-0 flex items-center gap-2">
                <span class="material-symbols-outlined text-sm text-primary">bolt</span> Indexes ({{ indexes.length }})
              </h3>
            </div>
            <div class="p-4">
              <div v-if="indexes.length === 0" class="text-center py-6 text-on-surface-variant text-xs italic font-mono">
                No indexes defined on table
              </div>
              <div v-else class="space-y-2 font-mono">
                <div
                  v-for="idx in indexes"
                  :key="idx.name"
                  class="p-3 bg-surface-container-lowest rounded-lg border border-outline-variant flex items-center justify-between"
                >
                  <div>
                    <div class="text-xs font-bold text-on-surface flex items-center gap-2">
                      <span class="material-symbols-outlined text-primary text-sm">bolt</span> {{ idx.name }}
                      <span v-if="idx.isUnique" class="px-1.5 py-0.2 bg-blue-500/20 text-blue-600 dark:text-blue-300 rounded text-[9px] uppercase font-bold">UNIQUE</span>
                    </div>
                    <div class="text-[11px] text-on-surface-variant mt-1">
                      Columns: <span class="text-primary font-bold">{{ idx.columns?.join(', ') }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Foreign Keys -->
          <div class="bg-surface border border-outline-variant rounded-xl overflow-hidden shadow-sm">
            <div class="px-5 py-3.5 border-b border-outline-variant bg-surface-container-low">
              <h3 class="text-xs font-bold text-on-surface uppercase tracking-wider font-mono m-0 flex items-center gap-2">
                <span class="material-symbols-outlined text-sm text-primary">link</span> Foreign Keys ({{ foreignKeys.length }})
              </h3>
            </div>
            <div class="p-4">
              <div v-if="foreignKeys.length === 0" class="text-center py-6 text-on-surface-variant text-xs italic font-mono">
                No foreign keys linked
              </div>
              <div v-else class="space-y-2 font-mono">
                <div
                  v-for="fk in foreignKeys"
                  :key="fk.name"
                  class="p-3 bg-surface-container-lowest rounded-lg border border-outline-variant space-y-1"
                >
                  <div class="text-xs font-bold text-on-surface flex items-center gap-2">
                    <span class="material-symbols-outlined text-primary text-sm">link</span> {{ fk.columnName }} ➔ <span class="text-primary">{{ fk.referencedTable }}.{{ fk.referencedColumn }}</span>
                  </div>
                  <div class="text-[10px] text-on-surface-variant">Constraint: {{ fk.name }}</div>
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
import { ref, watch } from 'vue';
import { useAppStore } from '../stores/app';
import { api } from '../api';

const store = useAppStore();
const columns = ref([]);
const indexes = ref([]);
const foreignKeys = ref([]);
const loading = ref(false);

async function loadDetails() {
  if (!store.selectedTable) return;
  loading.value = true;
  try {
    const data = await api.getTableSchema(store.activeConnectionId, store.selectedTable);
    columns.value = data.columns || [];
    indexes.value = data.indexes || [];
    foreignKeys.value = data.foreignKeys || [];
  } catch (err) {
    store.addToast(err.message || 'Failed to load table details', 'error');
  } finally {
    loading.value = false;
  }
}

watch(() => store.selectedTable, loadDetails, { immediate: true });
</script>