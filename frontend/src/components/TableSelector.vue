<template>
  <div class="relative inline-block text-left" ref="selectorRef">
    <!-- Trigger Button -->
    <button
      type="button"
      @click="toggleDropdown"
      :disabled="disabled || !store.activeConnection"
      class="flex items-center gap-2 px-3 py-1.5 rounded-lg border transition-all text-xs font-mono select-none"
      :class="[
        isOpen
          ? 'border-primary ring-2 ring-primary/20 bg-surface-container-high'
          : 'border-outline-variant bg-surface-container-lowest hover:border-primary/60 hover:bg-surface-container-low text-on-surface',
        disabled || !store.activeConnection ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer',
        variant === 'header' ? 'h-9 text-xs font-semibold' : 'text-xs'
      ]"
    >
      <span class="material-symbols-outlined text-[18px]" :class="selectedItem?.type === 'VIEW' ? 'text-blue-400' : 'text-primary'">
        {{ selectedItem?.type === 'VIEW' ? 'visibility' : 'table_rows' }}
      </span>

      <span class="font-bold truncate max-w-[160px] sm:max-w-[220px]">
        {{ currentTable || placeholder }}
      </span>

      <span
        v-if="selectedItem?.rowCount !== undefined && selectedItem?.rowCount !== null"
        class="text-[10px] text-on-surface-variant font-mono bg-surface-container-low px-1.5 py-0.5 rounded border border-outline-variant/40 hidden sm:inline"
      >
        {{ selectedItem.rowCount }}
      </span>

      <span class="material-symbols-outlined text-base text-on-surface-variant transition-transform duration-200" :class="{ 'rotate-180': isOpen }">
        expand_more
      </span>
    </button>

    <!-- Dropdown Popover -->
    <div
      v-if="isOpen"
      class="absolute left-0 mt-1.5 w-80 sm:w-96 rounded-xl bg-surface border border-outline-variant shadow-2xl z-50 overflow-hidden flex flex-col font-sans animate-in fade-in zoom-in-95 duration-100"
    >
      <!-- Popover Header & Search -->
      <div class="p-3 bg-surface-container-low border-b border-outline-variant space-y-2">
        <div class="flex items-center justify-between">
          <div class="text-[11px] font-bold uppercase tracking-wider text-on-surface-variant font-mono flex items-center gap-1.5">
            <span class="material-symbols-outlined text-[15px] text-primary">table_chart</span>
            Select Table / View
          </div>
          <span class="text-[10px] font-mono text-on-surface-variant bg-surface-container px-2 py-0.5 rounded border border-outline-variant/50">
            Total: {{ store.allTables.length }}
          </span>
        </div>

        <div class="relative">
          <span class="material-symbols-outlined absolute left-2.5 top-1/2 -translate-y-1/2 text-[16px] text-on-surface-variant">search</span>
          <input
            ref="searchInputRef"
            v-model="search"
            @keydown.esc="isOpen = false"
            @keydown.enter="selectFirstFiltered"
            placeholder="Search tables & views (e.g. users)..."
            class="w-full bg-surface-container-lowest border border-outline-variant rounded-lg pl-8 pr-7 py-1.5 text-xs text-on-surface placeholder:text-on-surface-variant/60 focus:outline-none focus:border-primary font-mono transition-colors"
          />
          <button
            v-if="search"
            @click="search = ''"
            class="absolute right-2.5 top-1/2 -translate-y-1/2 text-on-surface-variant hover:text-on-surface text-xs"
          >
            ✕
          </button>
        </div>

        <!-- Filter tabs: All, Tables, Views -->
        <div class="flex gap-1 pt-1">
          <button
            v-for="cat in categories"
            :key="cat.id"
            @click="filterType = cat.id"
            class="px-2.5 py-1 rounded text-[11px] font-mono transition-colors font-medium"
            :class="filterType === cat.id ? 'bg-primary text-on-primary font-bold shadow-sm' : 'bg-surface-container hover:bg-surface-container-high text-on-surface-variant'"
          >
            {{ cat.label }} ({{ cat.count }})
          </button>
        </div>
      </div>

      <!-- Item List Container -->
      <div class="max-h-64 overflow-y-auto p-1.5 divide-y divide-outline-variant/30 space-y-0.5">
        <div v-if="filteredItems.length === 0" class="text-center py-8 text-on-surface-variant text-xs font-mono space-y-1">
          <span class="material-symbols-outlined text-2xl text-on-surface-variant/50">search_off</span>
          <div>No matching tables or views</div>
          <button v-if="search" @click="search = ''" class="text-primary text-[11px] hover:underline font-bold">
            Clear search
          </button>
        </div>

        <div
          v-for="item in filteredItems"
          :key="item.name"
          @click="handleSelect(item.name)"
          class="flex items-center justify-between px-2.5 py-2 rounded-lg text-xs cursor-pointer transition-colors group"
          :class="currentTable === item.name ? 'bg-primary/15 text-primary font-bold border border-primary/30' : 'text-on-surface hover:bg-surface-container-low'"
        >
          <div class="flex items-center gap-2 truncate min-w-0 pr-2">
            <span
              class="material-symbols-outlined text-[16px] flex-shrink-0"
              :class="item.type === 'VIEW' ? 'text-blue-400' : 'text-on-surface-variant group-hover:text-primary transition-colors'"
            >
              {{ item.type === 'VIEW' ? 'visibility' : 'table_rows' }}
            </span>
            <span class="font-mono truncate select-text">{{ item.name }}</span>
            <span v-if="item.type === 'VIEW'" class="text-[9px] uppercase px-1 rounded bg-blue-500/10 text-blue-400 border border-blue-400/20">
              VIEW
            </span>
          </div>

          <div class="flex items-center gap-1.5 flex-shrink-0">
            <span
              v-if="item.rowCount !== undefined && item.rowCount !== null"
              class="text-[10px] text-on-surface-variant font-mono bg-surface-container px-1.5 py-0.5 rounded border border-outline-variant/40"
            >
              {{ item.rowCount }} rows
            </span>

            <!-- Quick Action Buttons -->
            <div v-if="showQuickActions" class="opacity-0 group-hover:opacity-100 flex items-center gap-0.5 transition-opacity">
              <button
                @click.stop="handleSelectTab(item.name, 'data')"
                class="p-1 text-on-surface-variant hover:text-primary hover:bg-surface-container rounded"
                title="Browse Data Table"
              >
                <span class="material-symbols-outlined text-[14px]">table_chart</span>
              </button>
              <button
                @click.stop="handleSelectTab(item.name, 'explorer')"
                class="p-1 text-on-surface-variant hover:text-primary hover:bg-surface-container rounded"
                title="View Table Schema"
              >
                <span class="material-symbols-outlined text-[14px]">schema</span>
              </button>
              <button
                @click.stop="handleQuerySql(item.name)"
                class="p-1 text-on-surface-variant hover:text-primary hover:bg-surface-container rounded"
                title="Query SELECT * FROM"
              >
                <span class="material-symbols-outlined text-[14px]">terminal</span>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Popover Footer -->
      <div class="px-3 py-2 bg-surface-container-low border-t border-outline-variant flex items-center justify-between text-[10px] text-on-surface-variant font-mono">
        <span>Click to select table</span>
        <button @click="store.refreshTables()" class="hover:text-primary flex items-center gap-1">
          <span class="material-symbols-outlined text-[13px]">refresh</span> Refresh list
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, nextTick, onMounted, onUnmounted } from 'vue';
import { useAppStore } from '../stores/app';

const props = defineProps({
  modelValue: {
    type: String,
    default: undefined,
  },
  placeholder: {
    type: String,
    default: 'Select a Table...',
  },
  variant: {
    type: String,
    default: 'header', // 'header' | 'compact' | 'inline'
  },
  showQuickActions: {
    type: Boolean,
    default: true,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(['update:modelValue', 'select']);

const store = useAppStore();
const isOpen = ref(false);
const search = ref('');
const filterType = ref('all'); // 'all' | 'table' | 'view'
const selectorRef = ref(null);
const searchInputRef = ref(null);

const currentTable = computed(() => {
  if (props.modelValue !== undefined) return props.modelValue;
  return store.selectedTable;
});

const selectedItem = computed(() => {
  if (!currentTable.value) return null;
  return store.allTables.find((t) => t.name === currentTable.value) || null;
});

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

function toggleDropdown() {
  if (props.disabled || !store.activeConnection) return;
  isOpen.value = !isOpen.value;
  if (isOpen.value) {
    nextTick(() => {
      searchInputRef.value?.focus();
    });
  }
}

function handleSelect(tableName) {
  emit('update:modelValue', tableName);
  emit('select', { table: tableName });
  store.selectTable(tableName);
  isOpen.value = false;
}

function handleSelectTab(tableName, tabName) {
  store.selectTable(tableName, tabName);
  isOpen.value = false;
}

function handleQuerySql(tableName) {
  store.queryTable(tableName);
  isOpen.value = false;
}

function selectFirstFiltered() {
  if (filteredItems.value.length > 0) {
    handleSelect(filteredItems.value[0].name);
  }
}

function handleClickOutside(e) {
  if (selectorRef.value && !selectorRef.value.contains(e.target)) {
    isOpen.value = false;
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside);
});

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside);
});
</script>
