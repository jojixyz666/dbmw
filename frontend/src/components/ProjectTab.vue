<template>
  <div class="h-full flex flex-col bg-dark-900 overflow-y-auto p-8 max-w-4xl mx-auto space-y-6">
    <!-- Header -->
    <div class="flex items-center gap-3 border-b border-slate-800 pb-4">
      <div class="p-2.5 bg-emerald-500/10 border border-emerald-500/20 rounded-xl text-emerald-400">
        📦
      </div>
      <div>
        <h2 class="text-base font-bold text-slate-100">Project Framework Detector & Config Generator</h2>
        <p class="text-xs text-slate-400">Inspect codebase for framework configurations and generate <code class="text-emerald-400">dbmw.yml</code></p>
      </div>
    </div>

    <!-- Inspector Card -->
    <div class="bg-dark-800 border border-slate-700/80 rounded-xl p-6 space-y-4 shadow">
      <div>
        <label class="block text-xs font-bold text-slate-300 uppercase tracking-wider mb-1.5">
          Project Folder Path
        </label>
        <div class="flex items-center gap-2">
          <input
            v-model="projectPath"
            placeholder="e.g. . (current folder) or D:/my-projects/laravel-app"
            class="flex-1 bg-dark-900 border border-slate-700 rounded-md px-3.5 py-2 text-xs text-slate-100 placeholder-slate-500 focus:outline-none focus:border-emerald-500 font-mono"
          />
          <button
            @click="detectFramework"
            :disabled="isDetecting"
            class="px-4 py-2 bg-emerald-600 hover:bg-emerald-500 text-white rounded-md text-xs font-semibold shrink-0 transition"
          >
            {{ isDetecting ? 'Inspecting...' : 'Detect Framework' }}
          </button>
        </div>
      </div>

      <!-- Detection Results -->
      <div v-if="detected" class="bg-dark-950 p-5 rounded-lg border border-slate-800 space-y-4 font-mono text-xs">
        <div class="flex items-center justify-between border-b border-slate-800 pb-3">
          <div class="flex items-center gap-2">
            <span class="text-xs font-bold text-slate-200">Framework:</span>
            <span class="px-2 py-0.5 bg-emerald-500/20 text-emerald-300 font-bold rounded border border-emerald-500/30">
              {{ detected.name }} ({{ detected.language }})
            </span>
          </div>
          <span class="text-[10px] text-slate-400">
            Detected: {{ detected.detectedFiles.join(', ') || 'none' }}
          </span>
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <span class="text-slate-400">Suggested Driver:</span>
            <span class="text-slate-100 font-bold ml-2">{{ detected.suggestedDriver || 'postgres' }}</span>
          </div>
          <div>
            <span class="text-slate-400">Suggested Host:</span>
            <span class="text-slate-100 ml-2 font-mono">{{ detected.suggestedHost || 'localhost' }}:{{ detected.suggestedPort || 5432 }}</span>
          </div>
          <div>
            <span class="text-slate-400">Suggested Database:</span>
            <span class="text-slate-100 ml-2 font-mono">{{ detected.suggestedDb || 'app_dev' }}</span>
          </div>
          <div>
            <span class="text-slate-400">Suggested User:</span>
            <span class="text-slate-100 ml-2 font-mono">{{ detected.suggestedUser || 'root' }}</span>
          </div>
        </div>

        <div class="pt-3 border-t border-slate-800 flex items-center justify-between">
          <span class="text-[11px] text-slate-400">Ready to create project descriptor file</span>
          <button
            @click="generateYaml"
            class="px-4 py-2 bg-emerald-600 hover:bg-emerald-500 text-white rounded text-xs font-bold transition"
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
