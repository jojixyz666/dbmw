<template>
  <div class="flex-1 flex flex-col h-full overflow-y-auto p-6 max-w-5xl mx-auto w-full space-y-6">
    <!-- Header Hero Banner -->
    <div class="bg-surface border border-outline-variant rounded-2xl p-6 shadow-sm flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
      <div class="flex items-center gap-4">
        <div class="w-12 h-12 rounded-xl bg-primary/10 border border-primary/20 flex items-center justify-center text-primary shadow-sm flex-shrink-0">
          <span class="material-symbols-outlined text-3xl">{{ icon }}</span>
        </div>
        <div>
          <h2 class="text-base font-bold text-on-surface tracking-tight m-0 font-sans">
            {{ title }}
          </h2>
          <p class="text-xs text-on-surface-variant m-0 mt-0.5 font-sans">
            {{ subtitle }}
          </p>
        </div>
      </div>

      <div class="flex items-center gap-2 font-mono text-xs text-on-surface-variant bg-surface-container px-3 py-1.5 rounded-lg border border-outline-variant/60">
        <span class="w-2 h-2 rounded-full bg-primary animate-pulse"></span>
        <span class="font-semibold">{{ store.activeConnection?.name || 'Active Database' }}</span>
        <span>•</span>
        <span>{{ store.tables.length }} Tables, {{ store.views.length }} Views</span>
      </div>
    </div>

    <!-- Filter & Search Bar -->
    <div class="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3 bg-surface border border-outline-variant rounded-xl p-3 shadow-sm">
      <div class="relative flex-1">
        <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant text-[18px]">search</span>
        <input
          v-model="search"
          type="text"
          placeholder="Filter tables & views by name (e.g. orders, users, logs)..."
          class="w-full bg-surface-container-lowest border border-outline-variant rounded-lg pl-9 pr-8 py-2 text-xs text-on-surface placeholder:text-on-surface-variant/60 focus:outline-none focus:border-primary font-mono transition-colors"
        />
        <button
          v-if="search"
          @click="search = ''"
          class="absolute right-3 top-1/2 -translate-y-1/2 text-on-surface-variant hover:text-on-surface text-xs"
        >
          ✕
        </button>
      </div>

      <!-- Filter Types -->
      <div class="flex items-center gap-1.5 font-mono text-xs shrink-0">
        <button
          v-for="cat in categories"
          :key="cat.id"
          @click="filterType = cat.id"
          class="px-3 py-1.5 rounded-lg font-medium transition-all"
          :class="filterType === cat.id ? 'bg-primary text-on-primary font-bold shadow-sm' : 'bg-surface-container hover:bg-surface-container-high text-on-surface-variant'"
        >
          {{ cat.label }} ({{ cat.count }})
        </button>
      </div>
    </div>

    <!-- Tables Grid -->
    <div v-if="filteredItems.length > 0" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3.5 pb-8">
      <div
        v-for="item in filteredItems"
        :key="item.name"
        @click="handleCardClick(item.name)"
        class="bg-surface border border-outline-variant hover:border-primary rounded-xl p-4 cursor-pointer transition-all hover:shadow-md group flex flex-col justify-between space-y-3 relative overflow-hidden"
      >
        <!-- Top row: Icon, Name, Row count badge -->
        <div class="flex items-start justify-between gap-2">
          <div class="flex items-center gap-2.5 min-w-0">
            <div class="w-8 h-8 rounded-lg bg-surface-container-low border border-outline-variant flex items-center justify-center group-hover:bg-primary/10 group-hover:border-primary/30 transition-colors flex-shrink-0">
              <span
                class="material-symbols-outlined text-[18px]"
                :class="item.type === 'VIEW' ? 'text-blue-400' : 'text-primary'"
              >
                {{ item.type === 'VIEW' ? 'visibility' : 'table_rows' }}
              </span>
            </div>
            <div class="min-w-0">
              <div class="font-mono font-bold text-xs text-on-surface truncate group-hover:text-primary transition-colors">
                {{ item.name }}
              </div>
              <div class="text-[10px] text-on-surface-variant font-mono uppercase mt-0.5">
                {{ item.type === 'VIEW' ? 'Database View' : (item.schema || 'Table') }}
              </div>
            </div>
          </div>

          <span
            v-if="item.rowCount !== undefined && item.rowCount !== null"
            class="text-[10px] font-mono bg-surface-container-low px-2 py-0.5 rounded border border-outline-variant text-on-surface-variant font-semibold flex-shrink-0"
          >
            {{ item.rowCount }} rows
          </span>
        </div>

        <!-- Comment / Description if available -->
        <p v-if="item.comment" class="text-xs text-on-surface-variant line-clamp-2 m-0 font-sans italic">
          {{ item.comment }}
        </p>

        <!-- Actions Footer -->
        <div class="pt-2 border-t border-outline-variant/60 flex items-center justify-between font-mono text-xs">
          <span class="text-[11px] text-primary font-bold group-hover:underline flex items-center gap-1">
            <span class="material-symbols-outlined text-sm">arrow_forward</span> Select Table
          </span>

          <div class="flex items-center gap-1">
            <button
              @click.stop="selectTab(item.name, 'data')"
              class="p-1.5 text-on-surface-variant hover:text-primary hover:bg-surface-container rounded-lg transition-colors"
              title="Browse Data Table"
            >
              <span class="material-symbols-outlined text-sm">table_chart</span>
            </button>
            <button
              @click.stop="selectTab(item.name, 'explorer')"
              class="p-1.5 text-on-surface-variant hover:text-primary hover:bg-surface-container rounded-lg transition-colors"
              title="Inspect Schema"
            >
              <span class="material-symbols-outlined text-sm">schema</span>
            </button>
            <button
              @click.stop="querySql(item.name)"
              class="p-1.5 text-on-surface-variant hover:text-primary hover:bg-surface-container rounded-lg transition-colors"
              title="Query SELECT * FROM"
            >
              <span class="material-symbols-outlined text-sm">terminal</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty Search Results State -->
    <div v-else-if="search" class="bg-surface border border-outline-variant rounded-2xl p-12 text-center text-on-surface-variant font-mono space-y-3 shadow-sm">
      <span class="material-symbols-outlined text-4xl text-on-surface-variant/40">search_off</span>
      <h3 class="text-sm font-bold text-on-surface m-0">No tables or views found matching "{{ search }}"</h3>
      <p class="text-xs max-w-sm mx-auto text-on-surface-variant/80">Try tweaking your search keywords or clear the filter to see all database tables.</p>
      <button
        @click="search = ''"
        class="px-4 py-2 bg-primary text-on-primary rounded-lg text-xs font-bold font-sans shadow-sm hover:bg-primary-container transition-colors inline-flex items-center gap-1.5"
      >
        <span class="material-symbols-outlined text-sm">clear_all</span> Show All Tables
      </button>
    </div>

    <!-- No Tables in Database State -->
    <div v-else class="bg-surface border border-outline-variant rounded-2xl p-12 text-center text-on-surface-variant font-mono space-y-3 shadow-sm">
      <span class="material-symbols-outlined text-4xl text-primary/50">layers_clear</span>
      <h3 class="text-sm font-bold text-on-surface m-0">No tables discovered in schema {{ store.currentSchema || 'default' }}</h3>
      <p class="text-xs max-w-sm mx-auto text-on-surface-variant/80">Check connected database or switch schemas, or execute create table queries.</p>
      <button
        @click="store.refreshExplorer()"
        class="px-4 py-2 bg-primary text-on-primary rounded-lg text-xs font-bold font-sans shadow-sm hover:bg-primary-container transition-colors inline-flex items-center gap-1.5"
      >
        <span class="material-symbols-outlined text-sm">refresh</span> Reload Tables
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useAppStore } from '../stores/app';

const props = defineProps({
  title: {
    type: String,
    default: 'Select a Table',
  },
  subtitle: {
    type: String,
    default: 'Pick a table or view to inspect records, columns, and relationships',
  },
  icon: {
    type: String,
    default: 'table_chart',
  },
  defaultTab: {
    type: String,
    default: 'data',
  },
});

const emit = defineEmits(['select']);
const store = useAppStore();

const search = ref('');
const filterType = ref('all'); // 'all' | 'table' | 'view'

const categories = computed(() => [
  { id: 'all', label: 'All', count: store.allTables.length },
  { id: 'table', label: 'Tables', count: store.tables.length },
  { id: 'view', label: 'Views', count: store.views.length },
]);

const filteredItems = computed(() => {
  let items = store.allTables;
  if (filterType.value === 'table') {
    items = items.filter((i) => i.type !== 'VIEW');
  } else if (filterType.value === 'view') {
    items = items.filter((i) => i.type === 'VIEW');
  }

  if (!search.value) return items;
  const q = search.value.toLowerCase().trim();
  return items.filter((i) => i.name.toLowerCase().includes(q));
});

function handleCardClick(tableName) {
  store.selectTable(tableName, props.defaultTab);
  emit('select', tableName);
}

function selectTab(tableName, tabName) {
  store.selectTable(tableName, tabName);
  emit('select', tableName);
}

function querySql(tableName) {
  store.queryTable(tableName);
}
</script>
