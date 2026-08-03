<template>
  <div class="h-full flex flex-col bg-surface overflow-hidden relative select-none">
    <!-- Action Bar -->
    <header class="h-[72px] px-6 border-b border-outline-variant flex items-center justify-between bg-surface z-20 flex-shrink-0">
      <div class="flex items-center gap-4">
        <div class="w-10 h-10 rounded-lg bg-surface-container-low border border-outline-variant flex items-center justify-center text-primary shadow-sm">
          <span class="material-symbols-outlined text-2xl">account_tree</span>
        </div>
        <div>
          <h2 class="text-sm font-bold text-on-surface tracking-tight m-0">Visual ERD (Entity-Relationship Diagram)</h2>
          <p class="text-xs text-on-surface-variant m-0">Foreign key relationships auto-mapped across schema tables</p>
        </div>
      </div>

      <!-- Zoom & Pan Controls -->
      <div class="flex items-center gap-2">
        <button
          @click="zoomIn"
          class="h-9 px-3 bg-surface-container-low hover:bg-surface-container text-on-surface border border-outline-variant rounded text-xs font-bold transition-colors shadow-sm"
          title="Zoom In"
        >
          ＋
        </button>
        <span class="text-xs text-on-surface-variant font-mono w-14 text-center font-bold bg-surface-container-lowest px-2 py-1.5 rounded border border-outline-variant">{{ Math.round(zoom * 100) }}%</span>
        <button
          @click="zoomOut"
          class="h-9 px-3 bg-surface-container-low hover:bg-surface-container text-on-surface border border-outline-variant rounded text-xs font-bold transition-colors shadow-sm"
          title="Zoom Out"
        >
          －
        </button>
        <button
          @click="resetView"
          class="h-9 px-3 bg-surface-container-low hover:bg-surface-container text-on-surface border border-outline-variant rounded text-xs font-semibold transition-colors shadow-sm"
        >
          Reset View
        </button>
        <button
          @click="loadGraph"
          class="h-9 px-3.5 bg-primary hover:bg-primary-container text-on-primary rounded text-xs font-bold flex items-center gap-1.5 transition-colors shadow-sm active:scale-95"
        >
          <span class="material-symbols-outlined text-[16px]">refresh</span> Re-layout
        </button>
      </div>
    </header>

    <!-- Canvas Container -->
    <div
      class="flex-1 overflow-hidden relative bg-surface-container-lowest cursor-grab active:cursor-grabbing"
      @mousedown="startPan"
      @mousemove="doPanOrDrag"
      @mouseup="endPanOrDrag"
      @wheel="onWheel"
    >
      <!-- Grid background -->
      <svg class="absolute inset-0 w-full h-full pointer-events-none opacity-25" xmlns="http://www.w3.org/2000/svg">
        <defs>
          <pattern id="erd-grid" width="30" height="30" patternUnits="userSpaceOnUse">
            <path d="M 30 0 L 0 0 0 30" fill="none" stroke="var(--border-outline-variant)" stroke-width="0.5"/>
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
              <path d="M 0 0 L 10 5 L 0 10 z" fill="var(--color-primary)" />
            </marker>
          </defs>

          <g v-for="edge in edgesWithCoordinates" :key="edge.id">
            <path
              :d="edge.pathSvg"
              fill="none"
              stroke="var(--color-primary)"
              stroke-width="2"
              stroke-dasharray="4 2"
              marker-end="url(#arrow)"
            />
            <text
              :x="edge.midX"
              :y="edge.midY - 8"
              fill="var(--color-primary)"
              font-size="10"
              font-family="monospace"
              text-anchor="middle"
              class="bg-surface"
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
          class="absolute w-64 bg-surface border border-outline-variant rounded-xl shadow-xl overflow-hidden cursor-move transition-all hover:border-primary hover:shadow-primary/10"
        >
          <!-- Table Node Header -->
          <div class="px-3.5 py-2.5 bg-surface-container-low border-b border-outline-variant flex items-center justify-between">
            <div class="flex items-center gap-2 font-mono font-bold text-xs text-on-surface truncate">
              <span class="material-symbols-outlined text-primary text-sm">table_rows</span>
              <span class="truncate">{{ node.label }}</span>
            </div>
            <span class="text-[9px] font-mono px-1.5 py-0.5 bg-primary/20 text-primary rounded border border-primary/30 font-bold">
              {{ node.columns?.length || 0 }} cols
            </span>
          </div>

          <!-- Column Attributes -->
          <div class="p-2 space-y-1 max-h-60 overflow-y-auto font-mono text-[11px]">
            <div
              v-for="col in node.columns"
              :key="col.name"
              class="flex items-center justify-between px-2 py-1 rounded hover:bg-surface-container-low transition-colors"
            >
              <div class="flex items-center gap-1.5 truncate">
                <span v-if="col.isPrimaryKey" class="material-symbols-outlined text-amber-400 text-[12px]" title="PK">key</span>
                <span :class="col.isPrimaryKey ? 'font-bold text-on-surface' : 'text-on-surface-variant'" class="truncate">
                  {{ col.name }}
                </span>
              </div>
              <span class="text-[10px] text-primary font-semibold shrink-0">{{ col.dataType }}</span>
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
const edges = ref([]);

const zoom = ref(1);
const panX = ref(60);
const panY = ref(60);

const isPanning = ref(false);
const startPanX = ref(0);
const startPanY = ref(0);

const draggingNode = ref(null);
const nodeDragOffsetX = ref(0);
const nodeDragOffsetY = ref(0);

async function loadGraph() {
  if (!store.activeConnectionId) return;
  try {
    const data = await api.getErdGraph(store.activeConnectionId);
    edges.value = data.edges || [];

    const rawNodes = data.nodes || [];
    const colsCount = Math.ceil(Math.sqrt(rawNodes.length)) || 1;
    nodes.value = rawNodes.map((n, idx) => {
      const col = idx % colsCount;
      const row = Math.floor(idx / colsCount);
      return {
        ...n,
        position: {
          x: col * 300 + 40,
          y: row * 300 + 40,
        }
      };
    });
  } catch (err) {
    store.addToast(err.message || 'Failed to load ERD graph', 'error');
  }
}

const edgesWithCoordinates = computed(() => {
  const nodeMap = new Map();
  nodes.value.forEach(n => nodeMap.set(n.id, n));

  return edges.value.map(edge => {
    const sourceNode = nodeMap.get(edge.source);
    const targetNode = nodeMap.get(edge.target);

    if (!sourceNode || !targetNode) {
      return { ...edge, pathSvg: '', midX: 0, midY: 0 };
    }

    const sX = sourceNode.position.x + 128;
    const sY = sourceNode.position.y + 60;
    const tX = targetNode.position.x + 128;
    const tY = targetNode.position.y + 60;

    const dx = tX - sX;
    const dy = tY - sY;
    const cx1 = sX + dx / 2;
    const cy1 = sY;
    const cx2 = sX + dx / 2;
    const cy2 = tY;

    const pathSvg = `M ${sX} ${sY} C ${cx1} ${cy1}, ${cx2} ${cy2}, ${tX} ${tY}`;

    return {
      ...edge,
      pathSvg,
      midX: (sX + tX) / 2,
      midY: (sY + tY) / 2,
    };
  });
});

function zoomIn() { zoom.value = Math.min(zoom.value + 0.15, 2.5); }
function zoomOut() { zoom.value = Math.max(zoom.value - 0.15, 0.4); }
function resetView() { zoom.value = 1; panX.value = 60; panY.value = 60; }

function onWheel(e) {
  e.preventDefault();
  if (e.deltaY < 0) zoomIn();
  else zoomOut();
}

function startPan(e) {
  if (e.target.tagName !== 'DIV' && e.target.tagName !== 'svg' && e.target.tagName !== 'rect') return;
  isPanning.value = true;
  startPanX.value = e.clientX - panX.value;
  startPanY.value = e.clientY - panY.value;
}

function startDragNode(e, node) {
  draggingNode.value = node;
  nodeDragOffsetX.value = (e.clientX / zoom.value) - node.position.x;
  nodeDragOffsetY.value = (e.clientY / zoom.value) - node.position.y;
}

function doPanOrDrag(e) {
  if (draggingNode.value) {
    draggingNode.value.position.x = (e.clientX / zoom.value) - nodeDragOffsetX.value;
    draggingNode.value.position.y = (e.clientY / zoom.value) - nodeDragOffsetY.value;
  } else if (isPanning.value) {
    panX.value = e.clientX - startPanX.value;
    panY.value = e.clientY - startPanY.value;
  }
}

function endPanOrDrag() {
  isPanning.value = false;
  draggingNode.value = null;
}

onMounted(loadGraph);
watch(() => store.activeConnectionId, loadGraph);
</script>