<template>
  <div class="fixed inset-0 z-40 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4">
    <div class="bg-dark-800 border border-slate-700/80 rounded-xl w-full max-w-md shadow-2xl p-6 space-y-4">
      <div class="flex items-center justify-between pb-3 border-b border-slate-700">
        <h3 class="text-sm font-bold text-slate-100 flex items-center gap-2">
          <span>⚙️</span> Application Settings
        </h3>
        <button @click="store.isSettingsModalOpen = false" class="text-slate-400 hover:text-white">✕</button>
      </div>

      <div class="space-y-4 text-xs">
        <div>
          <label class="block font-medium text-slate-300 mb-1">Theme</label>
          <select
            v-model="store.config.theme"
            class="w-full bg-dark-900 border border-slate-700 rounded-md px-3 py-2 text-slate-200 focus:outline-none focus:border-emerald-500"
          >
            <option value="dark">Dark (Emerald / Midnight)</option>
            <option value="light">Light</option>
          </select>
        </div>

        <div>
          <label class="block font-medium text-slate-300 mb-1">Default Grid Page Size</label>
          <select
            v-model.number="store.config.defaultPageSize"
            class="w-full bg-dark-900 border border-slate-700 rounded-md px-3 py-2 text-slate-200 focus:outline-none focus:border-emerald-500 font-mono"
          >
            <option :value="10">10 rows</option>
            <option :value="25">25 rows</option>
            <option :value="50">50 rows</option>
            <option :value="100">100 rows</option>
          </select>
        </div>

        <div>
          <label class="block font-medium text-slate-300 mb-1">Server Port</label>
          <input
            v-model.number="store.config.serverPort"
            type="number"
            class="w-full bg-dark-900 border border-slate-700 rounded-md px-3 py-2 text-slate-200 focus:outline-none focus:border-emerald-500 font-mono"
          />
        </div>
      </div>

      <div class="flex items-center justify-end gap-3 pt-3 border-t border-slate-700">
        <button
          @click="store.isSettingsModalOpen = false"
          class="px-4 py-1.5 bg-dark-700 text-slate-300 rounded text-xs"
        >
          Cancel
        </button>
        <button
          @click="saveSettings"
          class="px-4 py-1.5 bg-emerald-600 hover:bg-emerald-500 text-white rounded text-xs font-semibold"
        >
          Save Preferences
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useAppStore } from '../stores/app';
import { api } from '../api';

const store = useAppStore();

async function saveSettings() {
  try {
    await api.saveConfig(store.config);
    store.addToast('Settings saved successfully', 'success');
    store.isSettingsModalOpen = false;
  } catch (err) {
    store.addToast(`Failed to save settings: ${err.message}`, 'error');
  }
}
</script>
