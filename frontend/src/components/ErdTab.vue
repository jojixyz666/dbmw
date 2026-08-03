<template>
  <div class="h-full flex flex-col bg-dark-900 overflow-hidden relative select-none">
    <!-- Action Bar -->
    <div class="px-6 py-3 border-b border-slate-800 flex items-center justify-between bg-dark-900/90 z-20">
      <div class="flex items-center gap-3">
        <div class="p-2 bg-emerald-500/10 border border-emerald-500/20 rounded-lg text-emerald-400">
          🕸️
        </div>
        <div>
          <h2 class="text-sm font-bold text-slate-100">Visual ERD (Entity-Relationship Diagram)</h2>
          <p class="text-[11px] text-slate-400">Foreign key relationships auto-mapped across schema tables</p>
        </div>
      </div>

      <!-- Zoom & Pan Controls -->
      <div class="flex items-center gap-2">
        <button
          @click="zoomIn"
          class="px-2.5 py-1 bg-dark-800 hover:bg-dark-700 text-slate-200 border border-slate-700 rounded-md text-xs font-bold transition"
          title="Zoom In"
        >
          ＋
        </button>
        <span class="text-xs text-slate-400 font-mono w-12 text-center">{{ Math.round(zoom * 100) }}%</span>
        <button
          @click="zoomOut"
          class="px-2.5 py-1 bg-dark-800 hover:bg-dark-700 text-slate-200 border border-slate-700 rounded-md text-xs font-bold transition"
          title="Zoom Out"
        >
          －
        </button>
        <button
          @click="resetView"
          class="px-3 py-1 bg-dark-800 hover:bg-dark-700 text-slate-200 border border-slate-700 rounded-md text-xs font-medium transition"
        >
          Reset View
        </button>
        <button
          @click="loadGraph"
          class="px-3 py-1 bg-emerald-600 hover:bg-emerald-500 text-white rounded-md text-xs font-semibold flex items-center gap-1 transition shadow"
        >
          <span>🔄</span> Re-layout
        </button>
      </div>
    </div>

    <!-- Canvas Container -->
    <div
      class="flex-1 overflow-hidden relative bg-[#090d16] cursor-grab active:cursor-grabbing"
      @mousedown="startPan"
      @mousemove="doPanOrDrag"
      @mouseup="endPanOrDrag"
      @wheel="onWheel"
    >
      <!-- Grid background -->
      <svg class="absolute inset-0 w-full h-full pointer-events-none opacity-20" xmlns="http://www.w3.org/2000/svg">
        <defs>
          <pattern id="erd-grid" width="30" height="30" patternUnits="userSpaceOnUse">
            <path d="M 30 0 L 0 0 0 30" fill="none" stroke="#475569" stroke-width="0.5"/>
          </pattern>
        </defs>
        <rect width="100%" height="100%" fill="url(#erd-grid)" />
      </svg>

      <div
        class="absolute origin-top-left transition-transform duration-75 ease-out"
        :style="{
          transform: `translate(${panX}px, ${panY}px) scale(${zoom})`,
        }"
      >
        <!-- SVG Edges connecting FKs -->
        <svg class="absolute inset-0 overflow-visible pointer-events-none" style="width: 5000px; height: 5000px;">
          <defs>
            <marker id="arrow" viewBox="0 0 10 10" refX="5" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
              <path d="M 0 0 L 10 5 L 0 10 z" fill="#10b981" />
            </marker>
          </defs>

          <g v-for="edge in edgesWithCoordinates" :key="edge.id">
            <path
              :d="edge.pathSvg"
              fill="none"
              stroke="#10b981"
              stroke-width="2"
              stroke-dasharray="4 2"
              marker-end="url(#arrow)"
            />
            <text
              :x="edge.midX"
              :y="edge.midY - 8"
              fill="#94a3b8"
              font-size="10"
              font-family="monospace"
              text-anchor="middle"
              class="bg-dark-900"
            >
              {{ edge.label || '1:N' }}
            </text>
          </g>
        </svg>

        <!-- Table Nodes -->
        <div
          v-for="node in nodes"
          :key="node.id"
          :style="{
            left: `${node.position.x}px`,
            top: `${node.position.y}px`,
          }"
          @mousedown.stop="startDragNode($event, node)"
          class="absolute w-64 bg-dark-800/95 border border-slate-700/90 rounded-xl shadow-xl overflow-hidden cursor-move transition-shadow hover:shadow-2xl hover:border-emerald-500/80"
        >
          <!-- Table Node Header -->
          <div class="px-3 py-2 bg-dark-900/90 border-b border-slate-700/80 flex items-center justify-between">
            <div class="flex items-center gap-1.5 font-mono font-bold text-xs text-slate-100 truncate">
              <span class="text-emerald-400">📊</span>
              <span class="truncate">{{ node.label }}</span>
            </div>
            <span class="text-[9px] font-mono px-1.5 py-0.5 bg-slate-800 text-slate-400 rounded">
              {{ node.columns?.length || 0 }} cols
            </span>
          </div>

          <!-- Column Attributes -->
          <div class="p-2 space-y-1 max-h-56 overflow-y-auto font-mono text-[11px]">
            <div
              v-for="col in node.columns"
              :key="col.name"
              class="flex items-center justify-between px-1.5 py-0.5 rounded hover:bg-dark-700/60"
            >
              <div class="flex items-center gap-1.5 truncate text-slate-200">
                <span v-if="col.isPrimaryKey" class="text-amber-400 text-[10px]">🔑</span>
                <span v-else class="text-slate-600 text-[10px]">•</span>
                <span class="truncate">{{ col.name }}</span>
              </div>
              <span class="text-emerald-400/80 text-[10px] shrink-0">{{ col.dataType }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useAppStore } from '../stores/app';
import { api } from '../api';

const store = useAppStore();

const nodes = ref([]);
const rawEdges = ref([]);
const zoom = ref(1.0);
const panX = ref(60);
const panY = ref(60);

const isPanning = ref(false);
const panStart = ref({ x: 0, y: 0 });

const draggingNode = ref(null);
const dragOffset = ref({ x: 0, y: 0 });

async function loadGraph() {
  if (!store.activeConnectionId) return;

  try {
    const res = await api.generateERD(store.activeConnectionId, store.currentSchema);
    nodes.value = res.nodes || [];
    rawEdges.value = res.edges || [];
  } catch (err) {
    store.addToast(`Failed to load ERD: ${err.message}`, 'error');
  }
}

const edgesWithCoordinates = computed(() => {
  const nodeMap = new Map();
  nodes.value.forEach(n => nodeMap.set(n.id, n));

  return rawEdges.value.map(e => {
    const srcNode = nodeMap.get(e.source);
    const tgtNode = nodeMap.get(e.target);
    if (!srcNode || !tgtNode) return null;

    const x1 = srcNode.position.x + 256;
    const y1 = srcNode.position.y + 40;
    const x2 = tgtNode.position.x;
    const y2 = tgtNode.position.y + 40;

    const dx = Math.abs(x2 - x1) * 0.5;
    const pathSvg = `M ${x1} ${y1} C ${x1 + dx} ${y1}, ${x2 - dx} ${y2}, ${x2} ${y2}`;

    return {
      ...e,
      pathSvg,
      midX: (x1 + x2) / 2,
      midY: (y1 + y2) / 2,
    };
  }).filter(Boolean);
});

function zoomIn() {
  zoom.value = Math.min(zoom.value + 0.15, 2.0);
}

function zoomOut() {
  zoom.value = Math.max(zoom.value - 0.15, 0.4);
}

function resetView() {
  zoom.value = 1.0;
  panX.value = 60;
  panY.value = 60;
}

function onWheel(e) {
  e.preventDefault();
  if (e.deltaY < 0) {
    zoomIn();
  } else {
    zoomOut();
  }
}

function startPan(e) {
  isPanning.value = true;
  panStart.value = { x: e.clientX - panX.value, y: e.clientY - panY.value };
}

function startDragNode(e, node) {
  draggingNode.value = node;
  dragOffset.value = {
    x: (e.clientX / zoom.value) - node.position.x,
    y: (e.clientY / zoom.value) - node.position.y,
  };
}

function doPanOrDrag(e) {
  if (draggingNode.value) {
    draggingNode.value.position.x = (e.clientX / zoom.value) - dragOffset.value.x;
    draggingNode.value.position.y = (e.clientY / zoom.value) - dragOffset.value.y;
  } else if (isPanning.value) {
    panX.value = e.clientX - panStart.value.x;
    panY.value = e.clientY - panStart.value.y;
  }
}

function endPanOrDrag() {
  isPanning.value = false;
  draggingNode.value = null;
}

watch(() => [store.activeConnectionId, store.currentSchema], () => {
  loadGraph();
});

onMounted(() => {
  loadGraph();
});
</script>
