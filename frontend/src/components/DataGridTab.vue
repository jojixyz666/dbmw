<template>
  <div class="h-full flex flex-col bg-dark-900 overflow-hidden">
    <!-- Top Action Bar -->
    <div class="px-6 py-3 border-b border-slate-800 flex items-center justify-between bg-dark-900/80 flex-wrap gap-3">
      <div class="flex items-center gap-3">
        <div class="p-2 bg-emerald-500/10 border border-emerald-500/20 rounded-lg text-emerald-400">
          📋
        </div>
        <div>
          <div class="flex items-center gap-2">
            <h2 class="text-sm font-bold text-slate-100 font-mono">
              {{ store.selectedTable || 'Select a table' }}
            </h2>
            <span v-if="totalRows >= 0" class="text-xs text-slate-500 font-mono">
              ({{ totalRows }} rows)
            </span>
          </div>
          <p class="text-[11px] text-slate-400">Inline editable spreadsheet data browser</p>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <input
          v-model="searchTerm"
          @keyup.enter="fetchRows(1)"
          placeholder="Search columns..."
          class="bg-dark-800 border border-slate-700/80 rounded-md px-3 py-1 text-xs text-slate-100 placeholder-slate-500 focus:outline-none focus:border-emerald-500 w-44 font-mono"
        />

        <button
          @click="showFilterModal = true"
          class="px-3 py-1 bg-dark-800 hover:bg-dark-700 text-slate-200 border border-slate-700 rounded-md text-xs font-medium flex items-center gap-1 transition"
          :class="filters.length > 0 ? 'border-emerald-500 text-emerald-300' : ''"
        >
          <span>🎯</span> Filter {{ filters.length > 0 ? `(${filters.length})` : '' }}
        </button>

        <button
          @click="openInsertModal"
          :disabled="!store.selectedTable"
          class="px-3 py-1 bg-emerald-600 hover:bg-emerald-500 disabled:opacity-50 text-white rounded-md text-xs font-semibold flex items-center gap-1 transition shadow"
        >
          <span>+</span> Add Row
        </button>

        <button
          @click="fetchRows(page)"
          class="p-1 px-2 text-slate-400 hover:text-white bg-dark-800 border border-slate-700 rounded-md transition text-xs"
          title="Refresh"
        >
          🔄
        </button>
      </div>
    </div>

    <!-- Data Table Container -->
    <div class="flex-1 overflow-auto bg-dark-900 relative">
      <div v-if="!store.selectedTable" class="flex-1 h-full flex flex-col items-center justify-center text-slate-500 p-8 text-center">
        <div class="text-4xl mb-3">📋</div>
        <h3 class="text-base font-semibold text-slate-300 mb-1">No table selected</h3>
        <p class="text-xs max-w-sm">Select a table from the sidebar to inspect records, sort columns, filter rows, and edit values inline.</p>
      </div>

      <div v-else-if="loading" class="flex items-center justify-center h-full text-slate-500 text-sm">
        Loading table records...
      </div>

      <div v-else-if="rows.length === 0" class="flex flex-col items-center justify-center h-full text-slate-500 text-sm">
        <p class="mb-2">No records found matching criteria.</p>
        <button @click="openInsertModal" class="text-xs text-emerald-400 hover:underline">+ Insert a new record</button>
      </div>

      <table v-else class="w-full text-left text-xs border-collapse">
        <thead class="sticky top-0 bg-dark-800 text-slate-300 uppercase text-[10px] tracking-wider z-10 shadow border-b border-slate-700 font-mono select-none">
          <tr>
            <th class="py-2.5 px-3 bg-dark-900 border-r border-slate-800 text-center w-12 text-slate-500">Actions</th>
            <th
              v-for="col in columns"
              :key="col"
              @click="toggleSort(col)"
              class="py-2.5 px-4 cursor-pointer hover:bg-dark-700/80 transition border-r border-slate-800 whitespace-nowrap"
            >
              <div class="flex items-center justify-between gap-2">
                <span class="flex items-center gap-1.5 align-middle">
                  <span v-if="pkColumns.includes(col)" class="text-amber-400" title="Primary Key">🔑</span>
                  {{ col }}
                </span>
                <span v-if="sortBy?.column === col" class="text-emerald-400 font-bold">
                  {{ sortBy.direction === 'asc' ? '▲' : '▼' }}
                </span>
                <span v-else class="text-slate-600 opacity-40 hover:opacity-100">⇅</span>
              </div>
            </th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800/80 text-slate-200 font-mono">
          <tr
            v-for="(row, rIdx) in rows"
            :key="rIdx"
            class="hover:bg-dark-800/60 transition group"
          >
            <!-- Delete action button -->
            <td class="py-1.5 px-3 text-center border-r border-slate-800/60 bg-dark-900/40">
              <button
                @click="confirmDeleteRow(row)"
                class="opacity-40 group-hover:opacity-100 text-red-400 hover:text-red-300 transition text-xs p-0.5 rounded"
                title="Delete row"
              >
                🗑️
              </button>
            </td>

            <!-- Editable cells -->
            <td
              v-for="col in columns"
              :key="col"
              class="py-2 px-4 border-r border-slate-800/60 whitespace-nowrap max-w-xs truncate relative cursor-pointer group/cell"
              @click="startEdit(rIdx, col, row[col])"
            >
              <div v-if="editingCell?.row === rIdx && editingCell?.col === col" class="flex items-center gap-1">
                <input
                  v-model="editingCell.value"
                  @blur="saveEdit"
                  @keyup.enter="saveEdit"
                  @keyup.esc="cancelEdit"
                  ref="editInput"
                  class="bg-dark-950 border border-emerald-500 rounded px-2 py-0.5 text-xs text-slate-100 w-full focus:outline-none"
                />
              </div>
              <div v-else class="flex items-center justify-between">
                <span v-if="row[col] === null" class="text-slate-600 italic">NULL</span>
                <span v-else-if="typeof row[col] === 'boolean'" class="px-1.5 py-0.5 rounded text-[10px] font-bold" :class="row[col] ? 'bg-emerald-500/20 text-emerald-300' : 'bg-red-500/20 text-red-300'">
                  {{ row[col] ? 'TRUE' : 'FALSE' }}
                </span>
                <span v-else>{{ row[col] }}</span>
                <span class="opacity-0 group-hover/cell:opacity-60 text-[10px] text-slate-500 ml-1">✏️</span>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Pagination Footer -->
    <div class="px-6 py-2.5 border-t border-slate-800 flex items-center justify-between bg-dark-800/60 text-xs font-mono">
      <div class="text-slate-400 flex items-center gap-3">
        <span>Page {{ page }} of {{ totalPages || 1 }}</span>
        <span>•</span>
        <div class="flex items-center gap-1">
          <span>Rows per page:</span>
          <select
            v-model.number="pageSize"
            @change="fetchRows(1)"
            class="bg-dark-900 border border-slate-700 rounded px-2 py-0.5 text-slate-200 focus:outline-none focus:border-emerald-500"
          >
            <option :value="10">10</option>
            <option :value="25">25</option>
            <option :value="50">50</option>
            <option :value="100">100</option>
          </select>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <button
          @click="fetchRows(page - 1)"
          :disabled="page <= 1"
          class="px-3 py-1 bg-dark-800 hover:bg-dark-700 disabled:opacity-40 text-slate-200 border border-slate-700 rounded transition"
        >
          Previous
        </button>
        <button
          @click="fetchRows(page + 1)"
          :disabled="page >= totalPages"
          class="px-3 py-1 bg-dark-800 hover:bg-dark-700 disabled:opacity-40 text-slate-200 border border-slate-700 rounded transition"
        >
          Next
        </button>
      </div>
    </div>

    <!-- Insert Row Modal -->
    <div v-if="showInsertModal" class="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4">
      <div class="bg-dark-800 border border-slate-700 rounded-xl w-full max-w-xl shadow-2xl p-6 space-y-4 max-h-[85vh] overflow-y-auto">
        <div class="flex items-center justify-between pb-3 border-b border-slate-700">
          <h3 class="text-sm font-bold text-slate-100 font-mono">Insert Record into {{ store.selectedTable }}</h3>
          <button @click="showInsertModal = false" class="text-slate-400 hover:text-white">✕</button>
        </div>

        <div class="space-y-3 font-mono">
          <div v-for="col in metaColumns" :key="col.name">
            <label class="block text-xs font-semibold text-slate-300 mb-1 flex items-center justify-between">
              <span>{{ col.name }} <span class="text-emerald-400 font-normal text-[11px]">({{ col.dataType }})</span></span>
              <span v-if="col.isPrimaryKey" class="text-amber-400 text-[10px]">PK</span>
            </label>
            <input
              v-model="insertForm[col.name]"
              :placeholder="`Enter ${col.name}...`"
              class="w-full bg-dark-900 border border-slate-700 rounded-md px-3 py-1.5 text-xs text-slate-100 placeholder-slate-600 focus:outline-none focus:border-emerald-500"
            />
          </div>
        </div>

        <div class="flex items-center justify-end gap-3 pt-3 border-t border-slate-700">
          <button
            @click="showInsertModal = false"
            class="px-4 py-1.5 bg-dark-700 text-slate-300 rounded text-xs font-semibold"
          >
            Cancel
          </button>
          <button
            @click="submitInsert"
            class="px-4 py-1.5 bg-emerald-600 hover:bg-emerald-500 text-white rounded text-xs font-semibold"
          >
            Insert Row
          </button>
        </div>
      </div>
    </div>

    <!-- Filter Modal -->
    <div v-if="showFilterModal" class="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4">
      <div class="bg-dark-800 border border-slate-700 rounded-xl w-full max-w-lg shadow-2xl p-6 space-y-4">
        <div class="flex items-center justify-between pb-3 border-b border-slate-700">
          <h3 class="text-sm font-bold text-slate-100">Filter Records</h3>
          <button @click="showFilterModal = false" class="text-slate-400 hover:text-white">✕</button>
        </div>

        <div class="space-y-3 font-mono">
          <div v-for="(f, idx) in filters" :key="idx" class="flex items-center gap-2">
            <select
              v-model="f.column"
              class="bg-dark-900 border border-slate-700 rounded px-2 py-1 text-xs text-slate-200"
            >
              <option v-for="c in columns" :key="c" :value="c">{{ c }}</option>
            </select>

            <select
              v-model="f.operator"
              class="bg-dark-900 border border-slate-700 rounded px-2 py-1 text-xs text-slate-200"
            >
              <option value="eq">=</option>
              <option value="neq">!=</option>
              <option value="contains">contains</option>
              <option value="starts_with">starts with</option>
              <option value="gt">&gt;</option>
              <option value="lt">&lt;</option>
              <option value="is_null">is null</option>
              <option value="is_not_null">is not null</option>
            </select>

            <input
              v-if="f.operator !== 'is_null' && f.operator !== 'is_not_null'"
              v-model="f.value"
              placeholder="Value..."
              class="flex-1 bg-dark-900 border border-slate-700 rounded px-2 py-1 text-xs text-slate-100"
            />

            <button @click="filters.splice(idx, 1)" class="text-red-400 hover:text-red-300 px-1">✕</button>
          </div>

          <button
            @click="addFilterClause"
            class="text-xs text-emerald-400 hover:underline font-semibold"
          >
            + Add Filter Rule
          </button>
        </div>

        <div class="flex items-center justify-between pt-3 border-t border-slate-700">
          <button
            @click="filters = []; fetchRows(1); showFilterModal = false;"
            class="text-xs text-slate-400 hover:text-white"
          >
            Reset Filters
          </button>
          <div class="flex items-center gap-2">
            <button
              @click="showFilterModal = false"
              class="px-3 py-1.5 bg-dark-700 text-slate-300 rounded text-xs"
            >
              Cancel
            </button>
            <button
              @click="fetchRows(1); showFilterModal = false;"
              class="px-4 py-1.5 bg-emerald-600 hover:bg-emerald-500 text-white rounded text-xs font-semibold"
            >
              Apply Filter
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, watch, nextTick, onMounted } from 'vue';
import { useAppStore } from '../stores/app';
import { api } from '../api';

const store = useAppStore();

const rows = ref([]);
const columns = ref([]);
const pkColumns = ref([]);
const metaColumns = ref([]);
const page = ref(1);
const pageSize = ref(25);
const totalRows = ref(0);
const totalPages = ref(1);
const loading = ref(false);
const searchTerm = ref('');
const sortBy = ref(null);
const filters = ref([]);

const showInsertModal = ref(false);
const insertForm = reactive({});
const showFilterModal = ref(false);

const editingCell = ref(null);
const editInput = ref(null);

async function fetchRows(targetPage = 1) {
  if (!store.selectedTable || !store.activeConnectionId) return;

  try {
    loading.value = true;
    page.value = targetPage;

    const opts = {
      page: page.value,
      pageSize: pageSize.value,
      searchTerm: searchTerm.value,
      filters: filters.value,
      sortBy: sortBy.value,
    };

    const res = await api.browseData(store.activeConnectionId, store.currentSchema, store.selectedTable, opts);
    rows.value = res.rows || [];
    columns.value = res.columns || [];
    pkColumns.value = res.primaryKey || [];
    totalRows.value = res.totalRows || 0;
    totalPages.value = res.totalPages || 1;
  } catch (err) {
    store.addToast(`Failed to browse table: ${err.message}`, 'error');
  } finally {
    loading.value = false;
  }
}

function toggleSort(col) {
  if (!sortBy.value || sortBy.value.column !== col) {
    sortBy.value = { column: col, direction: 'asc' };
  } else if (sortBy.value.direction === 'asc') {
    sortBy.value = { column: col, direction: 'desc' };
  } else {
    sortBy.value = null;
  }
  fetchRows(1);
}

function addFilterClause() {
  filters.value.push({
    column: columns.value[0] || '',
    operator: 'eq',
    value: '',
  });
}

function startEdit(rowIdx, col, val) {
  editingCell.value = {
    row: rowIdx,
    col,
    value: val ?? '',
    original: val,
  };
  nextTick(() => {
    if (editInput.value?.[0]) editInput.value[0].focus();
  });
}

function cancelEdit() {
  editingCell.value = null;
}

async function saveEdit() {
  if (!editingCell.value) return;
  const { row: rIdx, col, value, original } = editingCell.value;
  if (value === original) {
    cancelEdit();
    return;
  }

  const rowObj = rows.value[rIdx];
  const pkMap = {};
  if (pkColumns.value.length > 0) {
    pkColumns.value.forEach(k => { pkMap[k] = rowObj[k]; });
  } else {
    // If no PK, fallback to full row match
    Object.assign(pkMap, rowObj);
  }

  try {
    await api.updateRow(store.activeConnectionId, store.currentSchema, store.selectedTable, pkMap, { [col]: value });
    rows.value[rIdx][col] = value;
    store.addToast('Cell updated', 'success');
  } catch (err) {
    store.addToast(`Update failed: ${err.message}`, 'error');
  } finally {
    cancelEdit();
  }
}

async function confirmDeleteRow(row) {
  if (!confirm('Are you sure you want to delete this row?')) return;

  const pkMap = {};
  if (pkColumns.value.length > 0) {
    pkColumns.value.forEach(k => { pkMap[k] = row[k]; });
  } else {
    Object.assign(pkMap, row);
  }

  try {
    await api.deleteRow(store.activeConnectionId, store.currentSchema, store.selectedTable, pkMap);
    store.addToast('Row deleted', 'success');
    await fetchRows(page.value);
  } catch (err) {
    store.addToast(`Delete failed: ${err.message}`, 'error');
  }
}

async function openInsertModal() {
  try {
    const cols = await api.listColumns(store.activeConnectionId, store.currentSchema, store.selectedTable);
    metaColumns.value = cols || [];
    Object.keys(insertForm).forEach(k => delete insertForm[k]);
    metaColumns.value.forEach(c => {
      insertForm[c.name] = '';
    });
    showInsertModal.value = true;
  } catch (err) {
    store.addToast(err.message, 'error');
  }
}

async function submitInsert() {
  try {
    const nonNullValues = {};
    for (const [k, v] of Object.entries(insertForm)) {
      if (v !== '') nonNullValues[k] = v;
    }
    await api.insertRow(store.activeConnectionId, store.currentSchema, store.selectedTable, nonNullValues);
    store.addToast('Row inserted', 'success');
    showInsertModal.value = false;
    await fetchRows(page.value);
  } catch (err) {
    store.addToast(`Insert failed: ${err.message}`, 'error');
  }
}

watch(() => [store.selectedTable, store.activeConnectionId, store.currentSchema], () => {
  if (store.selectedTable) {
    fetchRows(1);
  }
});

onMounted(() => {
  if (store.selectedTable) {
    fetchRows(1);
  }
});
</script>
