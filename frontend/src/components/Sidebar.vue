<template>
  <aside class="w-[270px] h-screen flex-shrink-0 border-r border-outline-variant bg-surface flex flex-col z-30 select-none shadow-sm">
    <!-- Brand Header -->
    <div class="px-4 py-3.5 border-b border-outline-variant flex items-center justify-between bg-surface-container-low/80">
      <div class="flex items-center gap-3">
        <div class="w-9 h-9 rounded-xl bg-primary/10 border border-primary/30 flex items-center justify-center text-primary shadow-sm">
          <span class="material-symbols-outlined text-xl">database</span>
        </div>
        <div class="flex flex-col leading-tight">
          <div class="flex items-center gap-1.5 font-bold text-on-surface text-sm tracking-tight font-sans">
            DBMW
            <span class="text-[9px] font-mono font-bold px-1.5 py-0.5 bg-primary/15 text-primary rounded border border-primary/25">v0.0.2</span>
          </div>
          <span class="text-[11px] text-on-surface-variant font-medium">Precision Workspace</span>
        </div>
      </div>

      <div class="flex items-center gap-1">
        <!-- Quick Theme Toggle Button -->
        <button
          @click="toggleTheme"
          class="text-on-surface-variant hover:text-on-surface p-1.5 rounded-lg hover:bg-surface-container transition-colors"
          :title="store.config.theme === 'dark' ? 'Switch to Light Mode' : 'Switch to Dark Mode'"
        >
          <span class="material-symbols-outlined text-lg">{{ store.config.theme === 'dark' ? 'light_mode' : 'dark_mode' }}</span>
        </button>

        <!-- Settings Button -->
        <button
          @click="store.isSettingsModalOpen = true"
          class="text-on-surface-variant hover:text-on-surface p-1.5 rounded-lg hover:bg-surface-container transition-colors"
          title="Settings"
        >
          <span class="material-symbols-outlined text-lg">settings</span>
        </button>
      </div>
    </div>

    <!-- Active Connection Card -->
    <div class="p-3 border-b border-outline-variant bg-surface-container-low/40">
      <div class="flex justify-between items-center mb-2">
        <span class="text-[10px] font-bold uppercase tracking-wider text-on-surface-variant font-mono">Active Connection</span>
        <button
          @click="store.isConnectionsModalOpen = true"
          class="text-primary text-[11px] font-bold hover:underline flex items-center gap-0.5 transition-colors"
        >
          <span class="material-symbols-outlined text-[14px]">add</span> Manage
        </button>
      </div>

      <div
        v-if="store.activeConnection"
        @click="store.isConnectionsModalOpen = true"
        class="border border-outline-variant rounded-xl p-2.5 bg-surface-container-lowest flex items-center justify-between cursor-pointer hover:border-primary transition-all shadow-sm group"
      >
        <div class="flex items-center gap-2.5 overflow-hidden">
          <span
            class="w-2.5 h-2.5 rounded-full ring-2 shrink-0 animate-pulse-subtle"
            :class="store.activeConnection.driver === 'postgres' ? 'bg-blue-500 ring-blue-500/30' : store.activeConnection.driver === 'mysql' ? 'bg-amber-500 ring-amber-500/30' : 'bg-primary ring-primary/30'"
          ></span>
          <div class="truncate">
            <div class="text-xs font-bold text-on-surface truncate group-hover:text-primary transition-colors">
              {{ store.activeConnection.name }}
            </div>
            <div class="text-[10px] text-on-surface-variant font-mono truncate mt-0.5">
              {{ store.activeConnection.driver }} • {{ store.activeConnection.database || store.activeConnection.filePath || 'default' }}
            </div>
          </div>
        </div>
        <span class="material-symbols-outlined text-on-surface-variant text-base group-hover:text-on-surface transition-colors">unfold_more</span>
      </div>

      <div
        v-else
        @click="store.isConnectionsModalOpen = true"
        class="border border-dashed border-outline-variant text-on-surface-variant rounded-xl p-3 text-center cursor-pointer hover:border-primary hover:text-primary transition-colors text-xs font-bold flex items-center justify-center gap-1.5 bg-surface-container-lowest"
      >
        <span class="material-symbols-outlined text-base text-primary">add_circle</span> Connect to Database
      </div>
    </div>

    <!-- Database / Schema Selector if applicable -->
    <div v-if="store.activeConnection && (store.databases.length > 1 || store.schemas.length > 1)" class="px-3 py-2 border-b border-outline-variant bg-surface-container-low/20 space-y-1.5 font-mono text-xs">
      <div v-if="store.databases.length > 1" class="flex items-center justify-between gap-1">
        <span class="text-[10px] font-bold text-on-surface-variant uppercase">DB:</span>
        <select
          :value="store.currentDatabase"
          @change="store.setDatabase($event.target.value)"
          class="flex-1 bg-surface-container-lowest border border-outline-variant rounded px-2 py-1 text-[11px] text-on-surface focus:outline-none focus:border-primary"
        >
          <option v-for="d in store.databases" :key="d.name" :value="d.name">{{ d.name }}</option>
        </select>
      </div>

      <div v-if="store.schemas.length > 1" class="flex items-center justify-between gap-1">
        <span class="text-[10px] font-bold text-on-surface-variant uppercase">Schema:</span>
        <select
          :value="store.currentSchema"
          @change="store.setSchema($event.target.value)"
          class="flex-1 bg-surface-container-lowest border border-outline-variant rounded px-2 py-1 text-[11px] text-on-surface focus:outline-none focus:border-primary"
        >
          <option v-for="s in store.schemas" :key="s.name" :value="s.name">{{ s.name }}</option>
        </select>
      </div>
    </div>

    <!-- Workspace Navigation Links Bar -->
    <div class="px-2 py-2 border-b border-outline-variant bg-surface-container-low/30 space-y-0.5">
      <button
        v-for="t in tabs"
        :key="t.id"
        @click="store.activeTab = t.id"
        class="w-full flex items-center gap-3 px-3 py-2 rounded-lg text-xs transition-all text-left font-medium relative group"
        :class="store.activeTab === t.id ? 'bg-primary text-on-primary font-bold shadow-sm' : 'text-on-surface-variant hover:text-on-surface hover:bg-surface-container-low'"
      >
        <span class="material-symbols-outlined text-lg leading-none">{{ t.icon }}</span>
        <span class="font-sans text-xs font-semibold">{{ t.name }}</span>
        <span v-if="store.activeTab === t.id" class="ml-auto w-1.5 h-1.5 rounded-full bg-on-primary"></span>
      </button>
    </div>

    <!-- Schema Tree / Explorer -->
    <div class="flex-1 overflow-y-auto p-3 space-y-3">
      <div v-if="store.activeConnection">
        <!-- Search filter -->
        <div class="relative mb-2">
          <span class="material-symbols-outlined absolute left-2.5 top-1/2 -translate-y-1/2 text-[16px] text-on-surface-variant">search</span>
          <input
            v-model="searchQuery"
            placeholder="Filter tables & views..."
            class="w-full bg-surface-container-lowest border border-outline-variant rounded-lg pl-8 pr-7 py-1.5 text-xs text-on-surface placeholder:text-on-surface-variant/60 focus:outline-none focus:border-primary font-mono transition-colors"
          />
          <button v-if="searchQuery" @click="searchQuery = ''" class="absolute right-2.5 top-1/2 -translate-y-1/2 text-on-surface-variant hover:text-on-surface text-xs">✕</button>
        </div>

        <!-- Selected Table Notice if any -->
        <div v-if="store.selectedTable" class="mb-2 p-2 bg-primary/10 border border-primary/25 rounded-lg flex items-center justify-between text-xs font-mono">
          <div class="flex items-center gap-1.5 truncate text-primary font-bold">
            <span class="material-symbols-outlined text-[14px]">check_circle</span>
            <span class="truncate">{{ store.selectedTable }}</span>
          </div>
          <button @click="store.clearSelectedTable()" class="text-on-surface-variant hover:text-primary text-[10px] underline ml-1" title="Deselect table">
            Clear
          </button>
        </div>

        <!-- Tables Section -->
        <div>
          <div class="flex items-center justify-between text-[10px] font-bold text-on-surface-variant uppercase tracking-wider mb-1.5 px-1 font-mono">
            <span>Tables ({{ filteredTables.length }})</span>
            <button
              @click="store.refreshExplorer()"
              class="text-on-surface-variant hover:text-primary transition-colors flex items-center gap-1"
              :class="{ 'animate-spin': store.tablesLoading }"
              title="Refresh Tables"
            >
              <span class="material-symbols-outlined text-sm">refresh</span>
            </button>
          </div>

          <!-- Loading state -->
          <div v-if="store.tablesLoading" class="text-xs text-on-surface-variant px-2 py-3 font-mono flex items-center gap-2">
            <span class="material-symbols-outlined text-primary text-base animate-spin">sync</span>
            <span>Loading tables...</span>
          </div>

          <!-- Empty search result -->
          <div v-else-if="searchQuery && filteredTables.length === 0" class="text-xs text-on-surface-variant px-1 italic font-mono py-1.5 space-y-1">
            <div>No tables match "{{ searchQuery }}"</div>
            <button @click="searchQuery = ''" class="text-primary text-[11px] underline font-bold">Clear search</button>
          </div>

          <!-- No tables in schema -->
          <div v-else-if="filteredTables.length === 0" class="text-xs text-on-surface-variant px-2 py-2 italic font-mono bg-surface-container-low/40 rounded-lg border border-outline-variant/30 text-center space-y-1 my-1">
            <div class="text-[11px]">No tables in {{ store.currentSchema || 'schema' }}</div>
            <button @click="store.refreshTables()" class="text-primary text-[10px] font-bold hover:underline flex items-center justify-center gap-1 mx-auto">
              <span class="material-symbols-outlined text-xs">sync</span> Reload tables
            </button>
          </div>

          <!-- Tables List -->
          <div v-else class="space-y-0.5">
            <div
              v-for="t in filteredTables"
              :key="t.name"
              @click="store.selectTable(t.name)"
              class="w-full flex items-center justify-between px-2.5 py-1.5 rounded-lg text-xs transition-colors text-left group cursor-pointer"
              :class="store.selectedTable === t.name ? 'bg-surface-container-high text-primary font-bold border border-primary/30 shadow-sm' : 'text-on-surface-variant hover:bg-surface-container-low hover:text-on-surface border border-transparent'"
            >
              <div class="flex items-center gap-2 truncate min-w-0 pr-1">
                <span class="material-symbols-outlined text-[16px] flex-shrink-0" :class="store.selectedTable === t.name ? 'text-primary' : 'text-on-surface-variant group-hover:text-primary transition-colors'">
                  table_rows
                </span>
                <span class="truncate font-mono">{{ t.name }}</span>
              </div>

              <div class="flex items-center gap-1 flex-shrink-0">
                <!-- Action hover icons -->
                <div class="opacity-0 group-hover:opacity-100 flex items-center gap-0.5">
                  <button
                    @click.stop="store.selectTable(t.name, 'data')"
                    class="p-0.5 hover:text-primary rounded text-on-surface-variant"
                    title="Browse Data"
                  >
                    <span class="material-symbols-outlined text-[13px]">table_chart</span>
                  </button>
                  <button
                    @click.stop="store.selectTable(t.name, 'explorer')"
                    class="p-0.5 hover:text-primary rounded text-on-surface-variant"
                    title="View Schema"
                  >
                    <span class="material-symbols-outlined text-[13px]">schema</span>
                  </button>
                  <button
                    @click.stop="store.queryTable(t.name)"
                    class="p-0.5 hover:text-primary rounded text-on-surface-variant"
                    title="Query SQL"
                  >
                    <span class="material-symbols-outlined text-[13px]">terminal</span>
                  </button>
                </div>

                <span v-if="t.rowCount" class="text-[10px] text-on-surface-variant font-mono bg-surface-container-low px-1.5 py-0.2 rounded border border-outline-variant/40">
                  {{ t.rowCount }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- Views Section -->
        <div v-if="filteredViews.length > 0" class="pt-2">
          <div class="text-[10px] font-bold text-on-surface-variant uppercase tracking-wider mb-1.5 px-1 font-mono">
            Views ({{ filteredViews.length }})
          </div>
          <div class="space-y-0.5">
            <div
              v-for="v in filteredViews"
              :key="v.name"
              @click="store.selectTable(v.name)"
              class="w-full flex items-center justify-between px-2.5 py-1.5 rounded-lg text-xs transition-colors text-left group cursor-pointer"
              :class="store.selectedTable === v.name ? 'bg-surface-container-high text-primary font-bold border border-primary/30 shadow-sm' : 'text-on-surface-variant hover:bg-surface-container-low hover:text-on-surface border border-transparent'"
            >
              <div class="flex items-center gap-2 truncate min-w-0 pr-1">
                <span class="material-symbols-outlined text-[16px] text-blue-400 flex-shrink-0">visibility</span>
                <span class="truncate font-mono">{{ v.name }}</span>
              </div>

              <!-- Action hover icons -->
              <div class="opacity-0 group-hover:opacity-100 flex items-center gap-0.5">
                <button
                  @click.stop="store.selectTable(v.name, 'data')"
                  class="p-0.5 hover:text-primary rounded text-on-surface-variant"
                  title="Browse Data"
                >
                  <span class="material-symbols-outlined text-[13px]">table_chart</span>
                </button>
                <button
                  @click.stop="store.selectTable(v.name, 'explorer')"
                  class="p-0.5 hover:text-primary rounded text-on-surface-variant"
                  title="View Schema"
                >
                  <span class="material-symbols-outlined text-[13px]">schema</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- No Connection Placeholder -->
      <div v-else class="text-center py-10 px-4 text-on-surface-variant space-y-3">
        <span class="material-symbols-outlined text-4xl opacity-50 text-primary">power_off</span>
        <p class="text-xs">No active database connection.</p>
        <button
          @click="store.isConnectionsModalOpen = true"
          class="px-3 py-1.5 bg-primary text-on-primary font-bold text-xs rounded-lg shadow transition-colors"
        >
          Configure Connection
        </button>
      </div>
    </div>

    <!-- Footer Stats / Status -->
    <div class="p-3 border-t border-outline-variant bg-surface-container-low flex items-center justify-between text-[11px] text-on-surface-variant font-mono">
      <span class="flex items-center gap-1.5">
        <span class="w-2 h-2 rounded-full" :class="store.activeConnection ? 'bg-primary' : 'bg-outline'"></span>
        {{ store.activeConnection ? 'Connected' : 'Offline' }}
      </span>
      <span>{{ store.tables.length }} Tables</span>
    </div>
  </aside>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useAppStore } from '../stores/app';

const store = useAppStore();
const searchQuery = ref('');

const tabs = [
  { id: 'editor', name: 'SQL Console', icon: 'terminal' },
  { id: 'data', name: 'Data Browser', icon: 'table_chart' },
  { id: 'explorer', name: 'Schema Explorer', icon: 'schema' },
  { id: 'erd', name: 'Visual ERD', icon: 'account_tree' },
  { id: 'project', name: 'Project Config', icon: 'folder_open' },
];

function toggleTheme() {
  const nextTheme = store.config.theme === 'dark' ? 'light' : 'dark';
  store.config.theme = nextTheme;
  store.applyTheme();
}

const filteredTables = computed(() => {
  if (!searchQuery.value) return store.tables;
  const q = searchQuery.value.toLowerCase();
  return store.tables.filter((t) => t.name.toLowerCase().includes(q));
});

const filteredViews = computed(() => {
  if (!searchQuery.value) return store.views;
  const q = searchQuery.value.toLowerCase();
  return store.views.filter((v) => v.name.toLowerCase().includes(q));
});
</script>
