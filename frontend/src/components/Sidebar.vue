<template>
  <aside class="w-72 bg-dark-900 border-r border-slate-800 flex flex-col h-full select-none">
    <!-- Brand Header -->
    <div class="px-4 py-3.5 border-b border-slate-800 flex items-center justify-between">
      <div class="flex items-center gap-2.5">
        <div class="p-1.5 bg-emerald-500/10 border border-emerald-500/30 rounded-lg text-emerald-400">
          <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <ellipse cx="12" cy="5" rx="9" ry="3"></ellipse>
            <path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"></path>
            <path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"></path>
          </svg>
        </div>
        <div>
          <span class="font-bold text-slate-100 text-sm tracking-tight flex items-center gap-1.5">
            DBMW
            <span class="text-[10px] uppercase font-mono px-1.5 py-0.2 bg-emerald-500/20 text-emerald-400 rounded">v0.0.1</span>
          </span>
          <p class="text-[11px] text-slate-400 leading-tight">Database Workspace</p>
        </div>
      </div>
      <button
        @click="store.isSettingsModalOpen = true"
        class="text-slate-400 hover:text-white p-1.5 rounded-md hover:bg-dark-800 transition"
        title="Settings"
      >
        ⚙️
      </button>
    </div>

    <!-- Active Connection Selector -->
    <div class="p-3 border-b border-slate-800 bg-dark-800/40">
      <div class="flex items-center justify-between mb-1.5">
        <span class="text-[10px] uppercase font-bold tracking-wider text-slate-400">Active Connection</span>
        <button
          @click="store.isConnectionsModalOpen = true"
          class="text-[11px] text-emerald-400 hover:text-emerald-300 font-medium flex items-center gap-1"
        >
          + Manage
        </button>
      </div>

      <div
        v-if="store.activeConnection"
        @click="store.isConnectionsModalOpen = true"
        class="bg-dark-800 hover:bg-dark-700/80 border border-slate-700/80 rounded-lg p-2.5 cursor-pointer transition flex items-center justify-between"
      >
        <div class="flex items-center gap-2 overflow-hidden">
          <span
            class="w-2 h-2 rounded-full ring-2 ring-emerald-500/20 shrink-0"
            :class="store.activeConnection.driver === 'postgres' ? 'bg-blue-400 ring-blue-500/20' : store.activeConnection.driver === 'mysql' ? 'bg-amber-400 ring-amber-500/20' : 'bg-emerald-400 ring-emerald-500/20'"
          ></span>
          <div class="truncate">
            <div class="text-xs font-semibold text-slate-100 truncate">{{ store.activeConnection.name }}</div>
            <div class="text-[10px] text-slate-400 font-mono truncate">
              {{ store.activeConnection.driver }} • {{ store.activeConnection.database || store.activeConnection.filePath || 'default' }}
            </div>
          </div>
        </div>
        <span class="text-slate-400 text-xs">▾</span>
      </div>

      <div
        v-else
        @click="store.isConnectionsModalOpen = true"
        class="bg-dark-800/70 border border-dashed border-slate-700 text-slate-400 rounded-lg p-3 text-center cursor-pointer hover:border-emerald-500/60 hover:text-emerald-400 transition text-xs"
      >
        + Connect to Database
      </div>
    </div>

    <!-- Navigation Tabs -->
    <div class="p-2 border-b border-slate-800 grid grid-cols-5 gap-1 bg-dark-900/60">
      <button
        v-for="t in tabs"
        :key="t.id"
        @click="store.activeTab = t.id"
        :title="t.label"
        class="flex flex-col items-center justify-center py-2 rounded-md transition text-xs relative"
        :class="store.activeTab === t.id ? 'bg-emerald-600/20 text-emerald-300 font-semibold border border-emerald-500/40' : 'text-slate-400 hover:bg-dark-800 hover:text-slate-200'"
      >
        <span class="text-sm mb-0.5">{{ t.icon }}</span>
        <span class="text-[10px] truncate max-w-[45px]">{{ t.short }}</span>
      </button>
    </div>

    <!-- Schema Tree / Explorer -->
    <div class="flex-1 overflow-y-auto p-3 space-y-4">
      <div v-if="store.activeConnection">
        <!-- Search filter -->
        <div class="relative mb-3">
          <input
            v-model="searchQuery"
            placeholder="Filter tables & views..."
            class="w-full bg-dark-800 border border-slate-700/80 rounded-md pl-7 pr-3 py-1.5 text-xs text-slate-100 placeholder-slate-500 focus:outline-none focus:border-emerald-500 font-mono"
          />
          <span class="absolute left-2.5 top-2 text-slate-500 text-xs">🔍</span>
          <button v-if="searchQuery" @click="searchQuery = ''" class="absolute right-2.5 top-1.5 text-slate-500 hover:text-white text-xs">✕</button>
        </div>

        <!-- Tables Section -->
        <div>
          <div class="flex items-center justify-between text-[11px] font-bold text-slate-400 uppercase tracking-wider mb-2 px-1">
            <span>Tables ({{ filteredTables.length }})</span>
            <button @click="store.refreshExplorer()" class="text-slate-500 hover:text-emerald-400 transition" title="Refresh">
              🔄
            </button>
          </div>

          <div v-if="filteredTables.length === 0" class="text-xs text-slate-500 px-1 italic">
            No tables found
          </div>
          <div v-else class="space-y-0.5">
            <button
              v-for="t in filteredTables"
              :key="t.name"
              @click="store.selectTable(t.name)"
              class="w-full flex items-center justify-between px-2.5 py-1.5 rounded-md text-xs transition text-left group"
              :class="store.selectedTable === t.name ? 'bg-emerald-600/20 text-emerald-300 font-semibold border border-emerald-500/30' : 'text-slate-300 hover:bg-dark-800 hover:text-slate-100'"
            >
              <div class="flex items-center gap-2 truncate">
                <span class="text-slate-500 group-hover:text-emerald-400">📊</span>
                <span class="truncate font-mono">{{ t.name }}</span>
              </div>
              <span v-if="t.rowCount" class="text-[10px] text-slate-500 font-mono">{{ t.rowCount }}</span>
            </button>
          </div>
        </div>

        <!-- Views Section -->
        <div v-if="filteredViews.length > 0" class="pt-2">
          <div class="text-[11px] font-bold text-slate-400 uppercase tracking-wider mb-2 px-1">
            Views ({{ filteredViews.length }})
          </div>
          <div class="space-y-0.5">
            <button
              v-for="v in filteredViews"
              :key="v.name"
              @click="store.selectTable(v.name)"
              class="w-full flex items-center gap-2 px-2.5 py-1.5 rounded-md text-xs transition text-left"
              :class="store.selectedTable === v.name ? 'bg-emerald-600/20 text-emerald-300 font-semibold border border-emerald-500/30' : 'text-slate-300 hover:bg-dark-800 hover:text-slate-100'"
            >
              <span class="text-slate-500">👁️</span>
              <span class="truncate font-mono">{{ v.name }}</span>
            </button>
          </div>
        </div>
      </div>

      <div v-else class="text-center py-12 text-slate-500 text-xs px-4">
        <p class="mb-3">No database connected.</p>
        <button
          @click="store.isConnectionsModalOpen = true"
          class="px-3 py-1.5 bg-emerald-600 hover:bg-emerald-500 text-white rounded-md text-xs transition font-semibold"
        >
          Add Connection
        </button>
      </div>
    </div>
  </aside>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useAppStore } from '../stores/app';

const store = useAppStore();
const searchQuery = ref('');

const tabs = [
  { id: 'explorer', label: 'Schema Explorer', short: 'Schema', icon: '🗂️' },
  { id: 'editor', label: 'SQL Query Editor', short: 'SQL', icon: '⚡' },
  { id: 'data', label: 'Data Browser Grid', short: 'Data', icon: '📋' },
  { id: 'erd', label: 'Visual ERD Diagram', short: 'ERD', icon: '🕸️' },
  { id: 'project', label: 'Project Config', short: 'Project', icon: '📦' },
];

const filteredTables = computed(() => {
  if (!searchQuery.value) return store.tables;
  const q = searchQuery.value.toLowerCase();
  return store.tables.filter(t => t.name.toLowerCase().includes(q));
});

const filteredViews = computed(() => {
  if (!searchQuery.value) return store.views;
  const q = searchQuery.value.toLowerCase();
  return store.views.filter(v => v.name.toLowerCase().includes(q));
});
</script>
