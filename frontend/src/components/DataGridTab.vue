<template>
  <div class="h-full flex flex-col bg-surface overflow-hidden">
    <!-- Top Action Bar -->
    <header class="h-[72px] px-6 border-b border-outline-variant flex items-center justify-between bg-surface flex-wrap gap-3 flex-shrink-0">
      <div class="flex items-center gap-3">
        <div class="w-10 h-10 rounded-lg bg-surface-container-low border border-outline-variant flex items-center justify-center text-primary shadow-sm flex-shrink-0">
          <span class="material-symbols-outlined text-2xl">table_chart</span>
        </div>
        <div class="flex flex-col">
          <div class="flex items-center gap-2">
            <!-- Accessible Table Selector Dropdown -->
            <TableSelector placeholder="Select a Table to Browse..." variant="header" />
            <span v-if="store.selectedTable && totalRows >= 0" class="text-xs text-primary font-mono bg-primary/10 px-2 py-0.5 rounded border border-primary/20 font-bold hidden sm:inline">
              {{ totalRows }} rows
            </span>
          </div>
          <p class="text-xs text-on-surface-variant m-0 mt-0.5">Inline editable spreadsheet data browser</p>
        </div>
      </div>

      <div class="flex items-center gap-2" v-if="store.selectedTable">
        <div class="relative">
          <span class="material-symbols-outlined absolute left-2.5 top-1/2 -translate-y-1/2 text-on-surface-variant text-[16px]">search</span>
          <input
            v-model="searchTerm"
            @keyup.enter="fetchRows(1)"
            placeholder="Search columns..."
            class="bg-surface-container-lowest border border-outline-variant rounded-md pl-8 pr-3 py-1.5 text-xs text-on-surface placeholder:text-on-surface-variant/70 focus:outline-none focus:border-primary w-48 font-mono transition-colors"
          />
        </div>

        <button
          @click="showFilterModal = true"
          class="h-9 px-3 bg-surface-container-low hover:bg-surface-container text-on-surface border border-outline-variant rounded text-xs font-semibold flex items-center gap-1.5 transition-colors shadow-sm"
          :class="filters.length > 0 ? 'border-primary text-primary bg-primary/10' : ''"
        >
          <span class="material-symbols-outlined text-[16px]">filter_alt</span> Filter {{ filters.length > 0 ? `(${filters.length})` : '' }}
        </button>

        <button
          @click="openInsertModal"
          :disabled="!store.selectedTable"
          class="h-9 px-3.5 bg-primary hover:bg-primary-container disabled:opacity-40 text-on-primary rounded text-xs font-bold flex items-center gap-1 transition-colors shadow-sm active:scale-95"
        >
          <span class="material-symbols-outlined text-[16px]">add</span> Add Row
        </button>

        <button
          @click="fetchRows(page)"
          class="h-9 px-2.5 text-on-surface-variant hover:text-on-surface bg-surface-container-low border border-outline-variant rounded transition-colors text-xs shadow-sm flex items-center justify-center"
          title="Refresh"
        >
          <span class="material-symbols-outlined text-[16px]">refresh</span>
        </button>
      </div>
    </header>

    <!-- Data Table Container -->
    <div class="flex-1 overflow-auto bg-surface-container-lowest relative flex flex-col">
      <TableSelectorEmptyState
        v-if="!store.selectedTable"
        title="Select a Table to Browse Data"
        subtitle="Choose any database table or view to inspect records, sort columns, and edit values"
        icon="table_chart"
        defaultTab="data"
        @select="fetchRows(1)"
      />

      <div v-else-if="loading" class="flex items-center justify-center h-full text-on-surface-variant text-xs font-mono gap-2">
        <span class="material-symbols-outlined text-primary text-xl animate-spin">sync</span> Loading table records...
      </div>

      <div v-else-if="rows.length === 0" class="flex flex-col items-center justify-center h-full text-on-surface-variant text-xs font-mono space-y-3">
        <p class="text-xs text-on-surface-variant">No records found matching criteria.</p>
        <button @click="openInsertModal" class="px-3 py-1.5 bg-primary text-on-primary rounded text-xs font-bold transition shadow">
          + Insert a new record
        </button>
      </div>

      <table v-else class="w-full text-left text-xs border-collapse">
        <thead class="sticky top-0 bg-surface-container-low text-on-surface-variant uppercase text-[10px] tracking-wider z-10 shadow-sm border-b border-outline-variant font-mono select-none">
          <tr>
            <th class="py-2.5 px-3 border-r border-outline-variant w-10 text-center bg-surface-container-low">#</th>
            <th
              v-for="col in columns"
              :key="col"
              @click="toggleSort(col)"
              class="py-2.5 px-4 cursor-pointer hover:bg-surface-container transition-colors border-r border-outline-variant whitespace-nowrap"
            >
              <div class="flex items-center justify-between gap-2">
                <span class="flex items-center gap-1.5 align-middle">
                  <span v-if="pkColumns.includes(col)" class="material-symbols-outlined text-amber-400 text-sm" title="Primary Key">key</span>
                  <span class="font-bold text-on-surface">{{ col }}</span>
                </span>
                <span v-if="sortBy?.column === col" class="text-primary font-bold">
                  {{ sortBy.direction === 'asc' ? '▲' : '▼' }}
                </span>
                <span v-else class="text-on-surface-variant opacity-40 hover:opacity-100">⇅</span>
              </div>
            </th>
          </tr>
        </thead>
        <tbody class="divide-y divide-outline-variant/60 text-on-surface font-mono">
          <tr
            v-for="(row, rIdx) in rows"
            :key="rIdx"
            class="hover:bg-surface-container/60 transition-colors group"
          >
            <!-- Delete action button -->
            <td class="py-1.5 px-3 text-center border-r border-outline-variant bg-surface-container-lowest">
              <button
                @click="confirmDeleteRow(row)"
                class="opacity-40 group-hover:opacity-100 text-red-400 hover:text-red-300 transition-opacity text-xs p-0.5 rounded flex items-center justify-center mx-auto"
                title="Delete row"
              >
                <span class="material-symbols-outlined text-sm">delete</span>
              </button>
            </td>

            <!-- Editable cells -->
            <td
              v-for="col in columns"
              :key="col"
              class="py-2 px-4 border-r border-outline-variant whitespace-nowrap max-w-xs truncate cursor-pointer hover:bg-primary/10 transition-colors"
              @dblclick="startEdit(row, col)"
            >
              <template v-if="editingCell?.row === row && editingCell?.col === col">
                <input
                  v-model="editValue"
                  @keyup.enter="saveEdit"
                  @keyup.esc="cancelEdit"
                  @blur="saveEdit"
                  ref="editInput"
                  class="w-full bg-surface-container border border-primary text-primary font-mono px-2 py-1 text-xs rounded focus:outline-none shadow"
                />
              </template>
              <template v-else>
                <span v-if="row[col] === null" class="text-on-surface-variant italic">NULL</span>
                <span v-else-if="typeof row[col] === 'boolean'" class="px-1.5 py-0.5 rounded text-[10px] font-bold" :class="row[col] ? 'bg-primary/20 text-primary' : 'bg-red-950/30 text-red-400'">
                  {{ row[col] }}
                </span>
                <span v-else>{{ row[col] }}</span>
              </template>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Pagination Footer -->
    <div v-if="store.selectedTable" class="px-6 py-2.5 bg-surface-container-low border-t border-outline-variant flex items-center justify-between text-xs font-mono text-on-surface-variant flex-shrink-0">
      <div class="flex items-center gap-2">
        <span>Page {{ page }} of {{ totalPages }}</span>
        <span class="text-outline-variant">•</span>
        <select
          v-model.number="pageSize"
          @change="fetchRows(1)"
          class="bg-surface-container-lowest border border-outline-variant text-on-surface rounded px-2 py-0.5 text-xs focus:outline-none focus:border-primary"
        >
          <option :value="10">10 / page</option>
          <option :value="25">25 / page</option>
          <option :value="50">50 / page</option>
          <option :value="100">100 / page</option>
        </select>
      </div>

      <div class="flex items-center gap-1.5">
        <button
          @click="fetchRows(1)"
          :disabled="page === 1"
          class="px-2.5 py-1 bg-surface-container hover:bg-surface-container-high disabled:opacity-30 rounded border border-outline-variant transition-colors"
        >
          « First
        </button>
        <button
          @click="fetchRows(page - 1)"
          :disabled="page === 1"
          class="px-2.5 py-1 bg-surface-container hover:bg-surface-container-high disabled:opacity-30 rounded border border-outline-variant transition-colors"
        >
          ‹ Prev
        </button>
        <button
          @click="fetchRows(page + 1)"
          :disabled="page >= totalPages"
          class="px-2.5 py-1 bg-surface-container hover:bg-surface-container-high disabled:opacity-30 rounded border border-outline-variant transition-colors"
        >
          Next ›
        </button>
        <button
          @click="fetchRows(totalPages)"
          :disabled="page >= totalPages"
          class="px-2.5 py-1 bg-surface-container hover:bg-surface-container-high disabled:opacity-30 rounded border border-outline-variant transition-colors"
        >
          Last »
        </button>
      </div>
    </div>

    <!-- Insert Row Modal -->
    <div v-if="showInsertModal" class="fixed inset-0 z-50 bg-black/75 backdrop-blur-sm flex items-center justify-center p-4">
      <div class="bg-surface border border-outline-variant rounded-xl w-full max-w-lg shadow-2xl overflow-hidden flex flex-col max-h-[85vh]">
        <div class="px-6 py-4 border-b border-outline-variant flex items-center justify-between bg-surface-container-low">
          <h3 class="font-bold text-on-surface flex items-center gap-2 font-mono text-sm">
            <span class="material-symbols-outlined text-primary">add</span> Insert New Record into <span class="text-primary">{{ store.selectedTable }}</span>
          </h3>
          <button @click="showInsertModal = false" class="text-on-surface-variant hover:text-on-surface font-mono">✕</button>
        </div>

        <div class="p-6 overflow-y-auto space-y-3 font-mono">
          <div v-for="col in columns" :key="col" class="space-y-1">
            <label class="block text-xs font-bold text-on-surface flex items-center gap-1">
              {{ col }}
              <span v-if="pkColumns.includes(col)" class="text-amber-400 text-[10px]">(Primary Key)</span>
            </label>
            <input
              v-model="newRowData[col]"
              :placeholder="pkColumns.includes(col) ? 'Auto-increment or required ID' : 'Value...'"
              class="w-full bg-surface-container-lowest border border-outline-variant rounded-md px-3 py-1.5 text-xs text-on-surface placeholder:text-on-surface-variant/50 focus:outline-none focus:border-primary font-mono"
            />
          </div>
        </div>

        <div class="px-6 py-4 bg-surface-container-low border-t border-outline-variant flex justify-end gap-2">
          <button @click="showInsertModal = false" class="px-4 py-2 bg-surface-container hover:bg-surface-container-high text-on-surface rounded text-xs font-semibold">Cancel</button>
          <button @click="saveNewRecord" class="px-4 py-2 bg-primary text-on-primary rounded text-xs font-bold shadow">Save Record</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue';
import { useAppStore } from '../stores/app';
import { api } from '../api';
import TableSelector from './TableSelector.vue';
import TableSelectorEmptyState from './TableSelectorEmptyState.vue';

const store = useAppStore();
const rows = ref([]);
const columns = ref([]);
const pkColumns = ref([]);
const loading = ref(false);
const page = ref(1);
const pageSize = ref(25);
const totalRows = ref(0);
const searchTerm = ref('');
const sortBy = ref(null);
const filters = ref([]);

const editingCell = ref(null);
const editValue = ref('');
const editInput = ref(null);

const showInsertModal = ref(false);
const newRowData = ref({});
const showFilterModal = ref(false);

const totalPages = computed(() => Math.ceil(totalRows.value / pageSize.value) || 1);

async function fetchRows(p = 1) {
  if (!store.selectedTable) return;
  page.value = p;
  loading.value = true;
  try {
    const res = await api.getTableRows(store.activeConnectionId, store.currentSchema, store.selectedTable, {
      page: page.value,
      pageSize: pageSize.value,
      searchTerm: searchTerm.value || undefined,
      sortBy: sortBy.value ? { column: sortBy.value.column, direction: sortBy.value.direction } : undefined,
    });
    rows.value = res.rows || [];
    columns.value = res.columns || [];
    pkColumns.value = res.primaryKey || [];
    totalRows.value = res.totalRows || 0;
  } catch (err) {
    store.addToast(err.message || 'Failed to fetch rows', 'error');
  } finally {
    loading.value = false;
  }
}

function toggleSort(col) {
  if (sortBy.value?.column === col) {
    if (sortBy.value.direction === 'asc') {
      sortBy.value.direction = 'desc';
    } else {
      sortBy.value = null;
    }
  } else {
    sortBy.value = { column: col, direction: 'asc' };
  }
  fetchRows(1);
}

function startEdit(row, col) {
  editingCell.value = { row, col };
  editValue.value = row[col] === null ? '' : String(row[col]);
  nextTick(() => {
    if (editInput.value && editInput.value[0]) {
      editInput.value[0].focus();
    }
  });
}

function cancelEdit() {
  editingCell.value = null;
}

async function saveEdit() {
  if (!editingCell.value) return;
  const { row, col } = editingCell.value;
  const oldValue = row[col];
  const newValue = editValue.value;
  editingCell.value = null;

  if (String(oldValue) === String(newValue)) return;

  try {
    await api.updateRow(store.activeConnectionId, store.currentSchema, store.selectedTable, {
      primaryKey: getPkValues(row),
      values: { [col]: newValue },
    });
    row[col] = newValue;
    store.addToast(`Updated ${col} successfully`);
  } catch (err) {
    store.addToast(err.message || 'Failed to update row', 'error');
    fetchRows(page.value);
  }
}

function getPkValues(row) {
  const pks = {};
  if (pkColumns.value.length > 0) {
    pkColumns.value.forEach(pk => { pks[pk] = row[pk]; });
  } else {
    columns.value.forEach(c => { pks[c] = row[c]; });
  }
  return pks;
}

async function confirmDeleteRow(row) {
  if (!confirm('Are you sure you want to delete this row?')) return;
  try {
    await api.deleteRow(store.activeConnectionId, store.currentSchema, store.selectedTable, getPkValues(row));
    store.addToast('Row deleted successfully');
    fetchRows(page.value);
  } catch (err) {
    store.addToast(err.message || 'Failed to delete row', 'error');
  }
}

function openInsertModal() {
  newRowData.value = {};
  columns.value.forEach(c => { newRowData.value[c] = ''; });
  showInsertModal.value = true;
}

async function saveNewRecord() {
  try {
    await api.insertRow(store.activeConnectionId, store.currentSchema, store.selectedTable, newRowData.value);
    showInsertModal.value = false;
    store.addToast('Record inserted successfully');
    fetchRows(1);
  } catch (err) {
    store.addToast(err.message || 'Failed to insert record', 'error');
  }
}

watch(() => store.selectedTable, () => {
  page.value = 1;
  sortBy.value = null;
  searchTerm.value = '';
  fetchRows(1);
});
</script>