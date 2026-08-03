<template>
  <div class="fixed bottom-5 right-5 z-50 flex flex-col gap-2.5 pointer-events-none font-mono">
    <div
      v-for="toast in store.toasts"
      :key="toast.id"
      class="pointer-events-auto px-4 py-3 rounded-xl shadow-2xl flex items-center gap-3 border text-xs transition-all transform animate-in slide-in-from-bottom duration-200 backdrop-blur-md"
      :class="toast.type === 'error' ? 'bg-red-950/90 border-red-800 text-red-100' : 'bg-surface border-primary/40 text-on-surface shadow-lg'"
    >
      <span v-if="toast.type === 'error'" class="material-symbols-outlined text-red-400 text-base">error</span>
      <span v-else class="material-symbols-outlined text-primary text-base">check_circle</span>
      <span class="font-bold max-w-sm break-words leading-relaxed">{{ toast.message }}</span>
      <button @click="store.removeToast(toast.id)" class="ml-2 opacity-60 hover:opacity-100 text-on-surface-variant hover:text-on-surface">✕</button>
    </div>
  </div>
</template>

<script setup>
import { useAppStore } from '../stores/app';
const store = useAppStore();
</script>