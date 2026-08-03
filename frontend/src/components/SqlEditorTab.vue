<template>
  <div class="h-full flex flex-col bg-surface overflow-hidden">
    <!-- Header Actions Area -->
    <header class="h-[72px] border-b border-outline-variant flex items-center justify-between px-6 bg-surface flex-shrink-0">
      <div class="flex items-center gap-4">
        <div class="w-10 h-10 rounded-lg bg-surface-container-low border border-outline-variant flex items-center justify-center text-primary shadow-sm">
          <span class="material-symbols-outlined text-2xl">bolt</span>
        </div>
        <div>
          <h1 class="font-bold text-on-surface text-base m-0 font-sans">SQL Query Console</h1>
          <p class="text-xs text-on-surface-variant m-0">Execute queries, inspect tabular results, and export dataset</p>
        </div>
      </div>

      <div class="flex items-center gap-3">
        <button
          @click="showHistory = !showHistory"
          class="h-9 px-4 rounded border border-outline-variant text-on-surface-variant flex items-center gap-2 hover:bg-surface-container-low transition-colors text-xs font-medium bg-transparent"
          :class="showHistory ? 'border-primary text-primary bg-primary/10' : ''"
        >
          <span class="material-symbols-outlined text-[18px]">history</span>
          History
        </button>

        <div class="flex rounded border border-outline-variant overflow-hidden h-9">
          <button
            @click="exportCSV"
            :disabled="!queryResult?.rows || queryResult.rows.length === 0"
            class="px-3 border-r border-outline-variant text-on-surface-variant disabled:opacity-40 flex items-center gap-1.5 hover:bg-surface-container-low transition-colors text-xs font-medium bg-transparent"
          >
            <span class="material-symbols-outlined text-[16px]">download</span>
            CSV
          </button>
          <button
            @click="exportJSON"
            :disabled="!queryResult?.rows || queryResult.rows.length === 0"
            class="px-3 text-on-surface-variant disabled:opacity-40 flex items-center gap-1.5 hover:bg-surface-container-low transition-colors text-xs font-medium bg-transparent"
          >
            <span class="material-symbols-outlined text-[16px]">data_object</span>
            JSON
          </button>
        </div>

        <button
          @click="runQuery"
          :disabled="isExecuting || !sqlText.trim()"
          class="h-9 px-4 rounded bg-primary text-on-primary disabled:opacity-50 flex items-center gap-2 hover:bg-primary-container transition-colors text-xs font-bold shadow-sm ml-2 active:scale-95"
        >
          <span v-if="isExecuting" class="material-symbols-outlined text-[18px] animate-spin">sync</span>
          <span v-else class="material-symbols-outlined text-[18px]">play_arrow</span>
          {{ isExecuting ? 'Running...' : 'Run Query (Ctrl+Enter)' }}
        </button>
      </div>
    </header>

    <!-- Main Workspace Split: Editor on Top, Results on Bottom -->
    <div class="flex-1 flex overflow-hidden">
      <div class="flex-1 flex flex-col overflow-hidden">
        <!-- Editor Input Panel -->
        <div class="h-48 border-b border-outline-variant flex flex-col bg-surface-container-lowest">
          <div class="px-4 py-2 bg-surface-container-low border-b border-outline-variant flex items-center justify-between text-xs font-mono text-on-surface-variant">
            <span class="flex items-center gap-2">
              <span class="material-symbols-outlined text-sm text-primary">terminal</span>
              SQL Editor [<span class="text-primary font-bold">{{ store.activeConnection?.driver || 'SQL' }}</span>]
            </span>
            <span>Press <kbd class="px-1.5 py-0.5 rounded border border-outline-variant bg-surface text-[10px] text-on-surface font-mono">Ctrl + Enter</kbd> to execute</span>
          </div>

          <textarea
            v-model="sqlText"
            @keydown.ctrl.enter="runQuery"
            @keydown.meta.enter="runQuery"
            placeholder="SELECT * FROM users LIMIT 10;"
            class="flex-1 w-full bg-surface-container-lowest p-4 font-mono text-xs text-primary placeholder:text-on-surface-variant/50 focus:outline-none resize-none leading-relaxed tracking-wide"
            spellcheck="false"
          ></textarea>
        </div>

        <!-- Query Execution Banner / Stats -->
        <div class="px-6 py-2 bg-surface-container-low border-b border-outline-variant flex items-center justify-between text-xs font-mono">
          <div class="flex items-center gap-3">
            <span class="text-on-surface-variant font-bold uppercase text-[10px]">Result</span>
            <span v-if="queryResult?.executionTimeMs !== undefined" class="text-primary font-semibold bg-primary/10 px-2 py-0.5 rounded border border-primary/20 flex items-center gap-1">
              <span class="material-symbols-outlined text-sm">schedule</span> {{ queryResult.executionTimeMs.toFixed(2) }} ms
            </span>
            <span v-if="queryResult?.rowsAffected !== undefined" class="text-on-surface-variant bg-surface-container px-2 py-0.5 rounded border border-outline-variant">
              • {{ queryResult.rowsAffected }} row(s)
            </span>
          </div>

          <div v-if="queryResult?.error" class="text-red-400 font-semibold truncate max-w-lg flex items-center gap-1">
            <span class="material-symbols-outlined text-sm">error</span> {{ queryResult.error }}
          </div>
        </div>

        <!-- Result View Table -->
        <div class="flex-1 overflow-auto bg-surface-container-lowest">
          <div v-if="isExecuting" class="flex items-center justify-center h-full text-on-surface-variant text-xs font-mono gap-2">
            <span class="material-symbols-outlined text-primary text-xl animate-spin">sync</span> Executing query against {{ store.activeConnection?.name }}...
          </div>

          <div v-else-if="!queryResult" class="flex flex-col items-center justify-center h-full text-on-surface-variant p-8 text-center">
            <span class="material-symbols-outlined text-5xl mb-3 text-primary/40">terminal</span>
            <p class="text-xs font-mono">Type your query above and click <strong class="text-primary">Run Query</strong></p>
          </div>

          <div v-else-if="queryResult.error" class="p-6 text-red-300 text-xs font-mono bg-red-950/20 m-4 border border-red-800/40 rounded-xl shadow-lg">
            <div class="font-bold mb-1.5 flex items-center gap-2">
              <span class="material-symbols-outlined text-red-400 text-sm">error</span> SQL Execution Error:
            </div>
            <div class="leading-relaxed">{{ queryResult.error }}</div>
          </div>

          <table v-else-if="queryResult.columns && queryResult.columns.length > 0" class="w-full text-left text-xs font-mono border-collapse">
            <thead class="sticky top-0 bg-surface-container-low text-on-surface-variant uppercase text-[10px] tracking-wider z-10 border-b border-outline-variant shadow-sm">
              <tr>
                <th class="py-2.5 px-3 border-r border-outline-variant text-on-surface-variant text-center w-12 bg-surface-container-low">#</th>
                <th
                  v-for="(col, idx) in queryResult.columns"
                  :key="col"
                  class="py-2.5 px-4 border-r border-outline-variant whitespace-nowrap"
                >
                  <div class="font-bold text-on-surface">{{ col }}</div>
                  <div class="text-[9px] text-primary font-normal">
                    {{ queryResult.columnTypes ? queryResult.columnTypes[idx] : '' }}
                  </div>
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-outline-variant/60 text-on-surface">
              <tr v-for="(row, rIdx) in queryResult.rows" :key="rIdx" class="hover:bg-surface-container/60 transition-colors">
                <td class="py-2 px-3 border-r border-outline-variant text-on-surface-variant text-center font-mono text-[11px] bg-surface-container-lowest select-none">{{ rIdx + 1 }}</td>
                <td v-for="col in queryResult.columns" :key="col" class="py-2 px-4 border-r border-outline-variant whitespace-nowrap max-w-sm truncate">
                  <span v-if="row[col] === null" class="text-on-surface-variant italic">NULL</span>
                  <span v-else>{{ row[col] }}</span>
                </td>
              </tr>
            </tbody>
          </table>

          <div v-else class="text-center py-12 text-on-surface-variant text-xs font-mono">
            Query executed successfully with 0 result rows.
          </div>
        </div>
      </div>

      <!-- History Drawer -->
      <div v-if="showHistory" class="w-80 border-l border-outline-variant bg-surface flex flex-col shadow-xl">
        <div class="px-4 py-3 border-b border-outline-variant flex items-center justify-between bg-surface-container-low">
          <span class="text-xs font-bold text-on-surface uppercase tracking-wider font-mono flex items-center gap-1.5">
            <span class="material-symbols-outlined text-sm">history</span> Query History
          </span>
          <button @click="clearHistory" class="text-[11px] text-red-400 hover:underline font-mono">
            Clear
          </button>
        </div>

        <div class="flex-1 overflow-y-auto p-3 space-y-2 font-mono">
          <div v-if="historyList.length === 0" class="text-center py-8 text-on-surface-variant text-xs">
            No history yet
          </div>
          <div
            v-for="item in historyList"
            :key="item.id"
            @click="sqlText = item.query"
            class="p-3 bg-surface-container-lowest hover:bg-surface-container-low border border-outline-variant rounded-xl cursor-pointer transition-colors text-xs group shadow-sm"
          >
            <div class="flex items-center justify-between text-[10px] text-on-surface-variant mb-1.5">
              <span :class="item.status === 'success' ? 'text-primary font-bold' : 'text-red-400 font-bold'">
                {{ item.status === 'success' ? '✔' : '✖' }} {{ item.durationMs.toFixed(1) }}ms
              </span>
              <span>{{ formatTime(item.executedAt) }}</span>
            </div>
            <div class="text-on-surface font-mono text-[11px] line-clamp-2 break-all group-hover:text-primary transition-colors">
              {{ item.query }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useAppStore } from '../stores/app';
import { api } from '../api';

const store = useAppStore();
const sqlText = ref('SELECT 1 as id, "Hello DBMW" as greeting;');
const queryResult = ref(null);
const isExecuting = ref(false);
const showHistory = ref(false);
const historyList = ref([]);

async function runQuery() {
  if (!sqlText.value.trim() || isExecuting.value) return;

  try {
    isExecuting.value = true;
    const res = await api.executeQuery(store.activeConnectionId, sqlText.value);
    queryResult.value = res;

    saveToHistory({
      id: Date.now(),
      query: sqlText.value,
      executedAt: new Date().toISOString(),
      durationMs: res.executionTimeMs || 0,
      status: res.error ? 'error' : 'success'
    });
  } catch (err) {
    queryResult.value = { error: err.message || 'Execution failed' };
  } finally {
    isExecuting.value = false;
  }
}

function saveToHistory(item) {
  historyList.value.unshift(item);
  if (historyList.value.length > 50) historyList.value.pop();
  localStorage.setItem('dbmw_sql_history', JSON.stringify(historyList.value));
}

function clearHistory() {
  historyList.value = [];
  localStorage.removeItem('dbmw_sql_history');
}

function formatTime(isoStr) {
  try {
    return new Date(isoStr).toLocaleTimeString();
  } catch {
    return '';
  }
}

function exportCSV() {
  if (!queryResult.value?.rows || !queryResult.value?.columns) return;
  const cols = queryResult.value.columns;
  const rows = queryResult.value.rows;
  let csv = cols.join(',') + '\n';
  rows.forEach(r => {
    csv += cols.map(c => JSON.stringify(r[c] ?? '')).join(',') + '\n';
  });
  downloadFile(csv, 'query_export.csv', 'text/csv');
}

function exportJSON() {
  if (!queryResult.value?.rows) return;
  const json = JSON.stringify(queryResult.value.rows, null, 2);
  downloadFile(json, 'query_export.json', 'application/json');
}

function downloadFile(content, filename, type) {
  const blob = new Blob([content], { type });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = filename;
  a.click();
  URL.revokeObjectURL(url);
}

onMounted(() => {
  const saved = localStorage.getItem('dbmw_sql_history');
  if (saved) {
    try { historyList.value = JSON.parse(saved); } catch {}
  }
});
</script>