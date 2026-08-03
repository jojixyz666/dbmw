<template>
  <div class="flex h-screen w-screen bg-dark-900 text-slate-100 overflow-hidden font-['Plus_Jakarta_Sans',sans-serif]">
    <!-- Left Navigation Sidebar -->
    <Sidebar />

    <!-- Center Workspace -->
    <main class="flex-1 flex flex-col h-full overflow-hidden bg-dark-900 border-l border-slate-800">
      <ExplorerTab v-if="store.activeTab === 'explorer'" />
      <SqlEditorTab v-else-if="store.activeTab === 'editor'" />
      <DataGridTab v-else-if="store.activeTab === 'data'" />
      <ErdTab v-else-if="store.activeTab === 'erd'" />
      <ProjectTab v-else-if="store.activeTab === 'project'" />
    </main>

    <!-- Modals -->
    <ConnectionsModal v-if="store.isConnectionsModalOpen" />
    <SettingsModal v-if="store.isSettingsModalOpen" />

    <!-- Interactive Notifications -->
    <Toast />
  </div>
</template>

<script setup>
import { onMounted } from 'vue';
import { useAppStore } from './stores/app';

import Sidebar from './components/Sidebar.vue';
import ExplorerTab from './components/ExplorerTab.vue';
import SqlEditorTab from './components/SqlEditorTab.vue';
import DataGridTab from './components/DataGridTab.vue';
import ErdTab from './components/ErdTab.vue';
import ProjectTab from './components/ProjectTab.vue';
import ConnectionsModal from './components/ConnectionsModal.vue';
import SettingsModal from './components/SettingsModal.vue';
import Toast from './components/Toast.vue';

const store = useAppStore();

onMounted(async () => {
  await store.loadConfig();
  await store.loadConnections();
});
</script>
