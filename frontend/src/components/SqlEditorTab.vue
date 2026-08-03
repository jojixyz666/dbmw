<template>
  <div class="h-full flex flex-col bg-dark-900 overflow-hidden">
    <!-- Action Bar -->
    <div class="px-6 py-3 border-b border-slate-800 flex items-center justify-between bg-dark-900/90 flex-wrap gap-3">
      <div class="flex items-center gap-3">
        <div class="p-2 bg-emerald-500/10 border border-emerald-500/20 rounded-lg text-emerald-400">
          ⚡
        </div>
        <div>
          <h2 class="text-sm font-bold text-slate-100">SQL Query Console</h2>
          <p class="text-[11px] text-slate-400">Execute queries, inspect results, and export data</p>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <button
          @click="showHistory = !showHistory"
          class="px-3 py-1.5 bg-dark-800 hover:bg-dark-700 text-slate-200 border border-slate-700 rounded-md text-xs font-semibold flex items-center gap-1.5 transition"
          :class="showHistory ? 'border-emerald-500 text-emerald-400' : ''"
        >
          <span>📜</span> History
        </button>

        <button
          @click="exportCSV"
          :disabled="!queryResult?.rows || queryResult.rows.length === 0"
          class="px-3 py-1.5 bg-dark-800 hover:bg-dark-700 disabled:opacity-40 text-slate-200 border border-slate-700 rounded-md text-xs font-semibold flex items-center gap-1 transition"
        >
          <span>📥</span> CSV
        </button>

        <button
          @click="exportJSON"
          :disabled="!queryResult?.rows || queryResult.rows.length === 0"
          class="px-3 py-1.5 bg-dark-800 hover:bg-dark-700 disabled:opacity-40 text-slate-200 border border-slate-700 rounded-md text-xs font-semibold flex items-center gap-1 transition"
        >
          <span>📥</span> JSON
        </button>

        <button
          @click="runQuery"
          :disabled="isExecuting || !sqlText.trim()"
          class="px-4 py-1.5 bg-emerald-600 hover:bg-emerald-500 disabled:opacity-50 text-white rounded-md text-xs font-bold flex items-center gap-1.5 transition shadow"
        >
          <span v-if="isExecuting">⏳</span>
          <span v-else>▶</span>
          {{ isExecuting ? 'Running...' : 'Run Query (Ctrl+Enter)' }}
        </button>
      </div>
    </div>

    <!-- Main Workspace Split: Editor on Top, Results on Bottom -->
    <div class="flex-1 flex overflow-hidden">
      <div class="flex-1 flex flex-col overflow-hidden">
        <!-- Editor Input Panel -->
        <div class="h-44 border-b border-slate-800 flex flex-col bg-dark-950">
          <div class="px-4 py-1.5 bg-dark-900/60 border-b border-slate-800/80 flex items-center justify-between text-[11px] font-mono text-slate-400">
            <span>SQL Editor [{{ store.activeConnection?.driver || 'SQL' }}]</span>
            <span>Press <kbd class="px-1.5 py-0.5 bg-dark-800 border border-slate-700 rounded text-[10px]">Ctrl+Enter</kbd> to run</span>
          </div>

          <textarea
            v-model="sqlText"
            @keydown.ctrl.enter="runQuery"
            @keydown.meta.enter="runQuery"
            placeholder="SELECT * FROM users LIMIT 10;"
            class="flex-1 w-full bg-dark-950 p-4 font-mono text-xs text-emerald-300 placeholder-slate-600 focus:outline-none resize-none leading-relaxed"
            spellcheck="false"
          ></textarea>
        </div>

        <!-- Query Execution Banner / Stats -->
        <div class="px-6 py-2 bg-dark-900 border-b border-slate-800 flex items-center justify-between text-xs font-mono">
          <div class="flex items-center gap-3">
            <span class="text-slate-400 font-bold uppercase text-[10px]">Result</span>
            <span v-if="queryResult?.executionTimeMs !== undefined" class="text-emerald-400">
              ⚡ {{ queryResult.executionTimeMs.toFixed(2) }} ms
            </span>
            <span v-if="queryResult?.rowsAffected !== undefined" class="text-slate-400">
              • {{ queryResult.rowsAffected }} row(s)
            </span>
          </div>

          <div v-if="queryResult?.error" class="text-red-400 font-semibold truncate max-w-lg">
            ❌ {{ queryResult.error }}
          </div>
        </div>

        <!-- Result View Table -->
        <div class="flex-1 overflow-auto bg-dark-900">
          <div v-if="isExecuting" class="flex items-center justify-center h-full text-slate-500 text-sm">
            Executing query against {{ store.activeConnection?.name }}...
          </div>

          <div v-else-if="!queryResult" class="flex flex-col items-center justify-center h-full text-slate-600 p-8 text-center">
            <div class="text-3xl mb-2">⚡</div>
            <p class="text-xs">Type a query in the editor above and hit <strong>Run Query</strong></p>
          </div>

          <div v-else-if="queryResult.error" class="p-6 text-red-300 text-xs font-mono bg-red-950/20 m-4 border border-red-800/40 rounded-lg">
            <div class="font-bold mb-1">SQL Execution Error:</div>
            <div>{{ queryResult.error }}</div>
          </div>

          <table v-else-if="queryResult.columns && queryResult.columns.length > 0" class="w-full text-left text-xs font-mono border-collapse">
            <thead class="sticky top-0 bg-dark-800 text-slate-300 uppercase text-[10px] tracking-wider z-10 border-b border-slate-700 shadow">
              <tr>
                <th
                  v-for="(col, idx) in queryResult.columns"
                  :key="col"
                  class="py-2.5 px-4 border-r border-slate-800 whitespace-nowrap"
                >
                  <div>{{ col }}</div>
                  <div class="text-[9px] text-slate-500 font-normal">
                    {{ queryResult.columnTypes ? queryResult.columnTypes[idx] : '' }}
                  </div>
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-800 text-slate-200">
              <tr v-for="(row, rIdx) in queryResult.rows" :key="rIdx" class="hover:bg-dark-800/60 transition">
                <td v-for="col in queryResult.columns" :key="col" class="py-2 px-4 border-r border-slate-800/60 whitespace-nowrap max-w-sm truncate">
                  <span v-if="row[col] === null" class="text-slate-600 italic">NULL</span>
                  <span v-else>{{ row[col] }}</span>
                </td>
              </tr>
            </tbody>
          </table>

          <div v-else class="text-center py-12 text-slate-500 text-xs font-mono">
            Query executed successfully. (0 result rows)
          </div>
        </div>
      </div>

      <!-- History Drawer -->
      <div v-if="showHistory" class="w-80 border-l border-slate-800 bg-dark-950 flex flex-col">
        <div class="px-4 py-3 border-b border-slate-800 flex items-center justify-between">
          <span class="text-xs font-bold text-slate-200 uppercase tracking-wider">Query History</span>
          <button @click="clearHistory" class="text-[11px] text-red-400 hover:underline">
            Clear
          </button>
        </div>

        <div class="flex-1 overflow-y-auto p-3 space-y-2 font-mono">
          <div v-if="historyList.length === 0" class="text-center py-8 text-slate-600 text-xs">
            No history yet
          </div>
          <div
            v-for="item in historyList"
            :key="item.id"
            @click="sqlText = item.query"
            class="p-2.5 bg-dark-900 hover:bg-dark-800 border border-slate-800 rounded-lg cursor-pointer transition text-xs group"
          >
            <div class="flex items-center justify-between text-[10px] text-slate-500 mb-1">
              <span :class="item.status === 'success' ? 'text-emerald-400' : 'text-red-400'">
                {{ item.status === 'success' ? '✔' : '✖' }} {{ item.durationMs.toFixed(1) }}ms
              </span>
              <span>{{ formatTime(item.executedAt) }}</span>
            </div>
            <div class="text-slate-300 font-mono text-[11px] line-clamp-2 break-all group-hover:text-emerald-300">
              {{ item.query }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue';
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
    if (res.error) {
      store.addToast(`Query error: ${res.error}`, 'error');
    } else {
      store.addToast(`Executed in ${res.executionTimeMs.toFixed(1)} ms`, 'success');
    }
    loadHistory();
  } catch (err) {
    store.addToast(err.message, 'error');
    queryResult.value = {
      error: err.message,
      executionTimeMs: 0,
    };
  } finally {
    isExecuting.value = false;
  }
}

async function loadHistory() {
  try {
    const list = await api.getHistory(store.activeConnectionId);
    historyList.value = list || [];
  } catch (err) {
    console.error(err);
  }
}

async function clearHistory() {
  try {
    await api.clearHistory(store.activeConnectionId);
    historyList.value = [];
    store.addToast('History cleared', 'success');
  } catch (err) {
    store.addToast(err.message, 'error');
  }
}

async function exportCSV() {
  if (!queryResult.value) return;
  try {
    const blob = await api.exportCSV(queryResult.value);
    downloadBlob(blob, 'query_export.csv');
  } catch (err) {
    store.addToast(err.message, 'error');
  }
}

async function exportJSON() {
  if (!queryResult.value) return;
  try {
    const blob = await api.exportJSON(queryResult.value);
    downloadBlob(blob, 'query_export.json');
  } catch (err) {
    store.addToast(err.message, 'error');
  }
}

function downloadBlob(blob, filename) {
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  window.URL.revokeObjectURL(url);
  document.body.removeChild(a);
}

function formatTime(isoStr) {
  if (!isoStr) return '';
  const d = new Date(isoStr);
  return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' });
}

watch(() => store.activeConnectionId, () => {
  loadHistory();
});

onMounted(() => {
  loadHistory();
});
</script>
