<template>
  <div class="h-full flex flex-col bg-surface-container-lowest overflow-y-auto p-8 max-w-4xl mx-auto space-y-6">
    <!-- Header -->
    <div class="flex items-center gap-4 border-b border-outline-variant pb-4">
      <div class="w-10 h-10 rounded-lg bg-surface-container-low border border-outline-variant flex items-center justify-center text-primary shadow-sm">
        <span class="material-symbols-outlined text-2xl">folder_open</span>
      </div>
      <div>
        <h2 class="text-base font-bold text-on-surface font-sans tracking-tight m-0">Project Framework Detector & Config Generator</h2>
        <p class="text-xs text-on-surface-variant m-0">Inspect codebase for framework configurations and generate <code class="text-primary font-mono font-bold">dbmw.yml</code></p>
      </div>
    </div>

    <!-- Inspector Card -->
    <div class="bg-surface border border-outline-variant rounded-xl p-6 space-y-5 shadow-sm">
      <div>
        <label class="block text-xs font-bold text-on-surface uppercase tracking-wider mb-2 font-mono">
          Project Folder Path
        </label>
        <div class="flex items-center gap-2">
          <input
            v-model="projectPath"
            placeholder="e.g. . (current folder) or D:/my-projects/laravel-app"
            class="flex-1 bg-surface-container-lowest border border-outline-variant rounded-lg px-4 py-2.5 text-xs text-on-surface placeholder:text-on-surface-variant/50 focus:outline-none focus:border-primary font-mono transition-colors"
          />
          <button
            @click="detectFramework"
            :disabled="isDetecting"
            class="px-5 py-2.5 bg-primary hover:bg-primary-container disabled:opacity-40 text-on-primary rounded-lg text-xs font-bold shrink-0 transition-colors shadow-sm active:scale-95 flex items-center gap-1.5"
          >
            <span v-if="isDetecting" class="material-symbols-outlined text-sm animate-spin">sync</span>
            <span v-else class="material-symbols-outlined text-sm">search</span>
            {{ isDetecting ? 'Inspecting...' : 'Detect Framework' }}
          </button>
        </div>
      </div>

      <!-- Detection Results -->
      <div v-if="detected" class="bg-surface-container-low p-5 rounded-lg border border-outline-variant space-y-4 font-mono">
        <div class="flex items-center justify-between border-b border-outline-variant pb-3">
          <div class="flex items-center gap-2">
            <span class="text-xs font-bold text-on-surface">Framework:</span>
            <span class="px-2.5 py-0.5 bg-primary/20 text-primary font-bold rounded text-xs border border-primary/30">
              {{ detected.name }} ({{ detected.language }})
            </span>
          </div>
          <span class="text-[10px] text-on-surface-variant">
            Detected: {{ detected.detectedFiles?.join(', ') || 'none' }}
          </span>
        </div>

        <div class="grid grid-cols-2 gap-4 text-xs">
          <div>
            <span class="text-on-surface-variant">Suggested Driver:</span>
            <span class="text-primary font-bold ml-2">{{ detected.suggestedDriver || 'postgres' }}</span>
          </div>
          <div>
            <span class="text-on-surface-variant">Suggested Host:</span>
            <span class="text-on-surface ml-2 font-mono">{{ detected.suggestedHost || 'localhost' }}:{{ detected.suggestedPort || 5432 }}</span>
          </div>
          <div>
            <span class="text-on-surface-variant">Suggested Database:</span>
            <span class="text-on-surface ml-2 font-mono">{{ detected.suggestedDb || 'app_dev' }}</span>
          </div>
          <div>
            <span class="text-on-surface-variant">Suggested User:</span>
            <span class="text-on-surface ml-2 font-mono">{{ detected.suggestedUser || 'root' }}</span>
          </div>
        </div>

        <div class="pt-3 border-t border-outline-variant flex items-center justify-between">
          <span class="text-[11px] text-on-surface-variant">Ready to create project descriptor file</span>
          <button
            @click="generateYaml"
            class="px-4 py-2 bg-primary text-on-primary rounded text-xs font-bold transition-colors shadow-sm active:scale-95"
          >
            Generate dbmw.yml
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useAppStore } from '../stores/app';
import { api } from '../api';

const store = useAppStore();
const projectPath = ref('.');
const isDetecting = ref(false);
const detected = ref(null);

async function detectFramework() {
  try {
    isDetecting.value = true;
    const res = await api.detectProject(projectPath.value);
    detected.value = res;
    store.addToast(`Detected framework: ${res.name}`, 'success');
  } catch (err) {
    store.addToast(`Detection failed: ${err.message}`, 'error');
  } finally {
    isDetecting.value = false;
  }
}

async function generateYaml() {
  if (!detected.value) return;

  const cfg = {
    version: '1',
    project_name: detected.value.name,
    default_connection: 'primary',
    connections: [
      {
        name: 'primary',
        driver: detected.value.suggestedDriver || 'postgres',
        host: detected.value.suggestedHost || 'localhost',
        port: detected.value.suggestedPort || 5432,
        user: detected.value.suggestedUser || '',
        database: detected.value.suggestedDb || 'dev_db',
      },
    ],
  };

  try {
    const res = await api.generateProjectConfig(projectPath.value, cfg);
    store.addToast(`Created ${res.path}`, 'success');
  } catch (err) {
    store.addToast(err.message, 'error');
  }
}
</script>