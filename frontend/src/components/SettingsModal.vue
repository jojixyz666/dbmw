<template>
  <div class="fixed inset-0 z-40 bg-black/75 backdrop-blur-md flex items-center justify-center p-4">
    <div class="glass-panel w-full max-w-md rounded-xl shadow-2xl p-6 space-y-5 border border-outline-variant/50">
      <div class="flex items-center justify-between pb-3 border-b border-outline-variant">
        <h3 class="text-sm font-bold text-on-surface flex items-center gap-2 font-mono m-0">
          <span class="material-symbols-outlined text-primary">settings</span> Application Settings
        </h3>
        <button @click="store.isSettingsModalOpen = false" class="text-on-surface-variant hover:text-on-surface font-mono">✕</button>
      </div>

      <div class="space-y-4 text-xs font-mono">
        <div>
          <label class="block font-bold text-on-surface mb-1.5">Appearance Mode</label>
          <select
            v-model="store.config.theme"
            @change="store.applyTheme()"
            class="w-full bg-surface-container-lowest border border-outline-variant rounded-lg px-3.5 py-2 text-on-surface focus:outline-none focus:border-primary transition-colors"
          >
            <option value="light">Precision Light (Default Clean White)</option>
            <option value="dark">Precision Dark (Midnight)</option>
          </select>
        </div>

        <div>
          <label class="block font-bold text-on-surface mb-1.5">Color Palette Accent</label>
          <select
            v-model="store.config.palette"
            @change="store.applyTheme()"
            class="w-full bg-surface-container-lowest border border-outline-variant rounded-lg px-3.5 py-2 text-on-surface focus:outline-none focus:border-primary transition-colors font-bold"
          >
            <option value="emerald">💚 Emerald Mint (Supabase Green)</option>
            <option value="indigo">💙 Indigo Blue (DataGrip Sapphire)</option>
            <option value="amber">💛 Amber Gold (Precision Amber)</option>
            <option value="violet">💜 Violet Deep (Modern Purple)</option>
          </select>
        </div>

        <div>
          <label class="block font-bold text-on-surface mb-1.5">Default Grid Page Size</label>
          <select
            v-model.number="store.config.defaultPageSize"
            class="w-full bg-surface-container-lowest border border-outline-variant rounded-lg px-3.5 py-2 text-on-surface focus:outline-none focus:border-primary font-mono transition-colors"
          >
            <option :value="10">10 rows</option>
            <option :value="25">25 rows</option>
            <option :value="50">50 rows</option>
            <option :value="100">100 rows</option>
          </select>
        </div>

        <div>
          <label class="block font-bold text-on-surface mb-1.5">Server Port</label>
          <input
            v-model.number="store.config.serverPort"
            type="number"
            class="w-full bg-surface-container-lowest border border-outline-variant rounded-lg px-3.5 py-2 text-on-surface focus:outline-none focus:border-primary font-mono transition-colors"
          />
        </div>
      </div>

      <div class="flex items-center justify-end gap-3 pt-3 border-t border-outline-variant">
        <button
          @click="store.isSettingsModalOpen = false"
          class="px-4 py-2 bg-surface-container hover:bg-surface-container-high text-on-surface rounded text-xs font-semibold font-mono"
        >
          Cancel
        </button>
        <button
          @click="saveSettings"
          class="px-4 py-2 bg-primary text-on-primary rounded text-xs font-bold shadow-sm font-mono active:scale-95 transition-colors"
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
    store.applyTheme();
    store.addToast('Settings saved successfully', 'success');
    store.isSettingsModalOpen = false;
  } catch (err) {
    store.addToast(`Failed to save settings: ${err.message}`, 'error');
  }
}
</script>