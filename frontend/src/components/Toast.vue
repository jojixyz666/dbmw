<template>
  <div class="fixed bottom-4 right-4 z-50 flex flex-col gap-2 pointer-events-none">
    <div
      v-for="toast in store.toasts"
      :key="toast.id"
      class="pointer-events-auto px-4 py-3 rounded-lg shadow-lg flex items-center gap-3 border text-sm transition-all transform animate-in slide-in-from-bottom duration-200"
      :class="toast.type === 'error' ? 'bg-red-950/90 border-red-800 text-red-100' : 'bg-slate-900 border-emerald-500/50 text-emerald-100'"
    >
      <span v-if="toast.type === 'error'" class="text-red-400">⚠️</span>
      <span v-else class="text-emerald-400">✨</span>
      <span class="font-medium max-w-sm break-words">{{ toast.message }}</span>
      <button @click="store.removeToast(toast.id)" class="ml-2 opacity-70 hover:opacity-100">✕</button>
    </div>
  </div>
</template>

<script setup>
import { useAppStore } from '../stores/app';
const store = useAppStore();
</script>
