<template>
  <div class="h-full flex flex-col bg-surface overflow-hidden">
    <!-- Header -->
    <header class="h-[72px] px-6 border-b border-outline-variant flex items-center justify-between bg-surface flex-shrink-0 flex-wrap gap-3">
      <div class="flex items-center gap-3">
        <div class="w-10 h-10 rounded-lg bg-surface-container-low border border-outline-variant flex items-center justify-center text-primary shadow-sm flex-shrink-0">
          <span class="material-symbols-outlined text-2xl">schema</span>
        </div>
        <div class="flex flex-col">
          <div class="flex items-center gap-2">
            <TableSelector placeholder="Select a Table to Inspect..." variant="header" />
            <span v-if="store.selectedTable" class="px-2 py-0.5 bg-primary/20 text-primary text-[10px] font-bold uppercase rounded border border-primary/30 hidden sm:inline">
              SCHEMA
            </span>
          </div>
          <p class="text-xs text-on-surface-variant m-0 mt-0.5">Structure, Columns, Indexes, and Constraints</p>
        </div>
      </div>

      <div v-if="store.selectedTable" class="flex items-center gap-2">
        <button
          @click="store.activeTab = 'data'"
          class="h-9 px-3.5 bg-primary hover:bg-primary-container text-on-primary rounded text-xs font-bold flex items-center gap-1.5 transition-colors shadow-sm active:scale-95"
        >
          <span class="material-symbols-outlined text-[16px]">table_chart</span> Browse Data
        </button>
        <button
          @click="store.queryTable(store.selectedTable)"
          class="h-9 px-3 bg-surface-container-low hover:bg-surface-container text-on-surface border border-outline-variant rounded text-xs font-semibold flex items-center gap-1.5 transition-colors shadow-sm"
          title="Query Table in SQL Console"
        >
          <span class="material-symbols-outlined text-[16px]">terminal</span> Query SQL
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
    <div v-if="!store.selectedTable" class="flex-1 overflow-auto flex flex-col">
      <TableSelectorEmptyState
        title="Select a Table to Inspect Schema"
        subtitle="Pick any database table or view to inspect columns, data types, primary keys, indexes, and foreign keys"
        icon="schema"
        defaultTab="explorer"
        @select="loadDetails"
      />
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
                  <td class="py-2.5 px-4 text-primary font-bold">{{ col.dataType }}</td>
                  <td class="py-2.5 px-4">
                    <span
                      class="px-2 py-0.5 rounded text-[10px] font-bold"
                      :class="col.isPrimaryKey ? 'bg-amber-500/10 text-amber-400 border border-amber-500/20' : 'text-on-surface-variant'"
                    >
                      {{ col.isPrimaryKey ? 'YES' : 'NO' }}
                    </span>
                  </td>
                  <td class="py-2.5 px-4 text-on-surface-variant">{{ col.isNullable ? 'YES' : 'NO' }}</td>
                  <td class="py-2.5 px-4 text-on-surface-variant">{{ col.defaultValue || 'NULL' }}</td>
                  <td class="py-2.5 px-4 text-on-surface-variant">{{ col.extra || '-' }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Indexes and Keys Grid -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <!-- Indexes Section -->
          <div class="bg-surface border border-outline-variant rounded-xl overflow-hidden shadow-sm">
            <div class="px-5 py-3.5 border-b border-outline-variant bg-surface-container-low flex items-center justify-between">
              <h3 class="text-xs font-bold text-on-surface uppercase tracking-wider font-mono flex items-center gap-2 m-0">
                <span class="material-symbols-outlined text-sm text-primary">grain</span> Indexes ({{ indexes.length }})
              </h3>
            </div>
            <div class="p-4 font-mono">
              <div v-if="indexes.length === 0" class="text-xs text-on-surface-variant italic">
                No custom indexes defined
              </div>
              <div v-else class="space-y-2">
                <div
                  v-for="idx in indexes"
                  :key="idx.name"
                  class="p-3 bg-surface-container-lowest border border-outline-variant/60 rounded-lg flex items-center justify-between"
                >
                  <div>
                    <div class="text-xs font-bold text-on-surface flex items-center gap-1.5">
                      <span v-if="idx.isPrimary" class="material-symbols-outlined text-amber-400 text-xs">key</span>
                      {{ idx.name }}
                    </div>
                    <div class="text-[10px] text-on-surface-variant mt-0.5">
                      Columns: <span class="text-primary font-semibold">{{ (idx.columns || []).join(', ') }}</span>
                    </div>
                  </div>
                  <span
                    v-if="idx.isUnique"
                    class="px-1.5 py-0.5 text-[9px] bg-emerald-500/10 text-emerald-400 rounded border border-emerald-500/20 font-bold"
                  >
                    UNIQUE
                  </span>
                </div>
              </div>
            </div>
          </div>

          <!-- Foreign Keys Section -->
          <div class="bg-surface border border-outline-variant rounded-xl overflow-hidden shadow-sm">
            <div class="px-5 py-3.5 border-b border-outline-variant bg-surface-container-low flex items-center justify-between">
              <h3 class="text-xs font-bold text-on-surface uppercase tracking-wider font-mono flex items-center gap-2 m-0">
                <span class="material-symbols-outlined text-sm text-primary">account_tree</span> Foreign Keys ({{ foreignKeys.length }})
              </h3>
            </div>
            <div class="p-4 font-mono">
              <div v-if="foreignKeys.length === 0" class="text-xs text-on-surface-variant italic">
                No foreign keys mapped
              </div>
              <div v-else class="space-y-2">
                <div
                  v-for="fk in foreignKeys"
                  :key="fk.name"
                  class="p-3 bg-surface-container-lowest border border-outline-variant/60 rounded-lg space-y-1"
                >
                  <div class="text-xs font-bold text-on-surface flex items-center gap-2">
                    <span class="material-symbols-outlined text-primary text-sm">link</span>
                    {{ fk.column || fk.columnName }} ➔
                    <span class="text-primary">{{ fk.refTable || fk.referencedTable }}.{{ fk.refColumn || fk.referencedColumn }}</span>
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
import TableSelector from './TableSelector.vue';
import TableSelectorEmptyState from './TableSelectorEmptyState.vue';

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
