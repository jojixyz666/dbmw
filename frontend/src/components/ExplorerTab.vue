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
            <div class="text-[11px] font-mono text-on-surface-variant flex items-center gap-3">
              <span v-if="pkCount > 0" class="flex items-center gap-1 text-amber-500 font-semibold">
                <span class="material-symbols-outlined text-xs">key</span> {{ pkCount }} PK
              </span>
              <span v-if="foreignKeys.length > 0" class="flex items-center gap-1 text-primary font-semibold">
                <span class="material-symbols-outlined text-xs">link</span> {{ foreignKeys.length }} FK
              </span>
            </div>
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
                    <span>{{ col.name }}</span>
                    <!-- FK indicator on column -->
                    <span
                      v-if="isFkColumn(col.name)"
                      class="px-1.5 py-0.2 rounded text-[9px] bg-primary/10 text-primary border border-primary/20 font-bold"
                      :title="getFkTooltip(col.name)"
                    >
                      FK
                    </span>
                  </td>
                  <td class="py-2.5 px-4 text-primary font-bold">{{ col.dataType }}</td>
                  <td class="py-2.5 px-4">
                    <span
                      class="px-2 py-0.5 rounded text-[10px] font-bold"
                      :class="col.isPrimaryKey ? 'bg-amber-500/10 text-amber-500 border border-amber-500/20' : 'text-on-surface-variant'"
                    >
                      {{ col.isPrimaryKey ? 'YES' : 'NO' }}
                    </span>
                  </td>
                  <td class="py-2.5 px-4 text-on-surface-variant">{{ col.isNullable ? 'YES' : 'NO' }}</td>
                  <td class="py-2.5 px-4 text-on-surface-variant">{{ col.defaultValue ?? 'NULL' }}</td>
                  <td class="py-2.5 px-4 text-on-surface-variant">{{ col.extra || col.comment || '-' }}</td>
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
              <div v-if="indexes.length === 0" class="text-xs text-on-surface-variant italic py-2 text-center">
                No indexes defined
              </div>
              <div v-else class="space-y-2.5">
                <div
                  v-for="idx in indexes"
                  :key="idx.name"
                  class="p-3 bg-surface-container-lowest border border-outline-variant/60 rounded-lg flex items-center justify-between"
                >
                  <div class="space-y-1">
                    <div class="text-xs font-bold text-on-surface flex items-center gap-1.5">
                      <span v-if="idx.isPrimary" class="material-symbols-outlined text-amber-400 text-xs">key</span>
                      <span>{{ idx.name }}</span>
                      <span v-if="idx.type" class="text-[9px] px-1.5 py-0.2 rounded bg-surface-container text-on-surface-variant border border-outline-variant/50">
                        {{ idx.type }}
                      </span>
                    </div>
                    <div class="text-[10px] text-on-surface-variant flex items-center gap-1 flex-wrap">
                      <span>Columns:</span>
                      <span
                        v-for="col in (idx.columns || [])"
                        :key="col"
                        class="px-1.5 py-0.5 rounded bg-primary/10 text-primary border border-primary/20 font-semibold"
                      >
                        {{ col }}
                      </span>
                      <span v-if="!idx.columns || idx.columns.length === 0" class="text-on-surface-variant italic">
                        -
                      </span>
                    </div>
                  </div>
                  <div class="flex items-center gap-1.5">
                    <span
                      v-if="idx.isPrimary"
                      class="px-1.5 py-0.5 text-[9px] bg-amber-500/10 text-amber-500 rounded border border-amber-500/20 font-bold"
                    >
                      PRIMARY
                    </span>
                    <span
                      v-if="idx.isUnique"
                      class="px-1.5 py-0.5 text-[9px] bg-emerald-500/10 text-emerald-500 rounded border border-emerald-500/20 font-bold"
                    >
                      UNIQUE
                    </span>
                  </div>
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
              <div v-if="foreignKeys.length === 0" class="text-xs text-on-surface-variant italic py-2 text-center">
                No foreign keys defined for this table
              </div>
              <div v-else class="space-y-2.5">
                <div
                  v-for="fk in foreignKeys"
                  :key="fk.name || `${fk.column}_${fk.refTable}`"
                  class="p-3 bg-surface-container-lowest border border-outline-variant/60 rounded-lg space-y-2"
                >
                  <div class="flex items-center justify-between flex-wrap gap-1">
                    <div class="text-xs font-bold text-on-surface flex items-center gap-1.5">
                      <span class="material-symbols-outlined text-primary text-sm">link</span>
                      <span class="px-1.5 py-0.5 rounded bg-surface-container text-on-surface font-mono font-bold">{{ fk.column || fk.columnName }}</span>
                      <span class="text-primary font-bold">➔</span>
                      <button
                        @click="jumpToTable(fk.refTable || fk.referencedTable)"
                        class="px-1.5 py-0.5 rounded bg-primary/10 hover:bg-primary/20 text-primary border border-primary/20 font-mono font-bold transition-colors cursor-pointer"
                        :title="`View schema for ${fk.refTable || fk.referencedTable}`"
                      >
                        {{ fk.refTable || fk.referencedTable }}.{{ fk.refColumn || fk.referencedColumn }}
                      </button>
                    </div>

                    <!-- Actions / Rules badges -->
                    <div class="flex items-center gap-1 text-[9px] font-mono">
                      <span v-if="fk.onDelete && fk.onDelete !== 'NO ACTION'" class="px-1.5 py-0.5 rounded bg-rose-500/10 text-rose-400 border border-rose-500/20" :title="'ON DELETE ' + fk.onDelete">
                        DEL: {{ fk.onDelete }}
                      </span>
                      <span v-if="fk.onUpdate && fk.onUpdate !== 'NO ACTION'" class="px-1.5 py-0.5 rounded bg-amber-500/10 text-amber-400 border border-amber-500/20" :title="'ON UPDATE ' + fk.onUpdate">
                        UPD: {{ fk.onUpdate }}
                      </span>
                    </div>
                  </div>

                  <div class="text-[10px] text-on-surface-variant flex items-center justify-between">
                    <div>Constraint: <span class="font-mono text-on-surface">{{ fk.name }}</span></div>
                    <span v-if="fk.refSchema" class="text-[9px] text-on-surface-variant">Schema: {{ fk.refSchema }}</span>
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
import { ref, computed, watch } from 'vue';
import { useAppStore } from '../stores/app';
import { api } from '../api';
import TableSelector from './TableSelector.vue';
import TableSelectorEmptyState from './TableSelectorEmptyState.vue';

const store = useAppStore();
const columns = ref([]);
const indexes = ref([]);
const foreignKeys = ref([]);
const loading = ref(false);

const pkCount = computed(() => columns.value.filter((c) => c.isPrimaryKey).length);

function isFkColumn(colName) {
  return foreignKeys.value.some((fk) => (fk.column || fk.columnName) === colName);
}

function getFkTooltip(colName) {
  const matching = foreignKeys.value.filter((fk) => (fk.column || fk.columnName) === colName);
  return matching.map((fk) => `➔ ${(fk.refTable || fk.referencedTable)}.${(fk.refColumn || fk.referencedColumn)} (${fk.name})`).join('\n');
}

function jumpToTable(tableName) {
  if (!tableName) return;
  store.selectTable(tableName, 'explorer');
}

async function loadDetails() {
  if (!store.selectedTable) {
    columns.value = [];
    indexes.value = [];
    foreignKeys.value = [];
    return;
  }
  loading.value = true;
  try {
    const data = await api.getTableSchema(store.activeConnectionId, store.currentSchema, store.selectedTable);
    columns.value = (data && data.columns) || [];
    indexes.value = (data && data.indexes) || [];
    foreignKeys.value = (data && data.foreignKeys) || [];

    // If foreignKeys or indexes came back empty, try targeted fallback endpoints
    if (foreignKeys.value.length === 0) {
      try {
        const directFks = await api.listForeignKeys(store.activeConnectionId, store.currentSchema, store.selectedTable);
        if (directFks && directFks.length > 0) {
          foreignKeys.value = directFks;
        }
      } catch (_) {}
    }

    if (indexes.value.length === 0) {
      try {
        const directIdxs = await api.listIndexes(store.activeConnectionId, store.currentSchema, store.selectedTable);
        if (directIdxs && directIdxs.length > 0) {
          indexes.value = directIdxs;
        }
      } catch (_) {}
    }
  } catch (err) {
    store.addToast(err.message || 'Failed to load table details', 'error');
  } finally {
    loading.value = false;
  }
}

watch(
  [() => store.selectedTable, () => store.currentSchema, () => store.activeConnectionId],
  () => {
    loadDetails();
  },
  { immediate: true }
);
</script>
