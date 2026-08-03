<template>
  <div class="h-full flex flex-col bg-surface overflow-hidden relative select-none">
    <!-- Action Bar -->
    <header class="h-[72px] px-6 border-b border-outline-variant flex items-center justify-between bg-surface z-20 flex-shrink-0 flex-wrap gap-3">
      <div class="flex items-center gap-3">
        <div class="w-10 h-10 rounded-lg bg-surface-container-low border border-outline-variant flex items-center justify-center text-primary shadow-sm flex-shrink-0">
          <span class="material-symbols-outlined text-2xl">account_tree</span>
        </div>
        <div>
          <h2 class="text-sm font-bold text-on-surface tracking-tight m-0">Visual ERD (Entity-Relationship Diagram)</h2>
          <p class="text-xs text-on-surface-variant m-0">Foreign key relationships auto-mapped across schema tables</p>
        </div>
      </div>

      <!-- Table Focus & Zoom & Pan Controls -->
      <div class="flex items-center gap-2">
        <TableSelector
          placeholder="Focus Table..."
          variant="header"
          @select="focusTable"
        />

        <div class="flex items-center gap-1 bg-surface-container-lowest border border-outline-variant rounded-lg p-1">
          <button
            @click="zoomIn"
            class="h-7 px-2.5 bg-surface-container-low hover:bg-surface-container text-on-surface border border-outline-variant rounded text-xs font-bold transition-colors shadow-sm"
            title="Zoom In"
          >
            ＋
          </button>
          <span class="text-xs text-on-surface-variant font-mono w-12 text-center font-bold px-1">{{ Math.round(zoom * 100) }}%</span>
          <button
            @click="zoomOut"
            class="h-7 px-2.5 bg-surface-container-low hover:bg-surface-container text-on-surface border border-outline-variant rounded text-xs font-bold transition-colors shadow-sm"
            title="Zoom Out"
          >
            －
          </button>
          <button
            @click="resetView"
            class="h-7 px-2 bg-surface-container-low hover:bg-surface-container text-on-surface border border-outline-variant rounded text-xs font-medium transition-colors shadow-sm"
            title="Reset Canvas View"
          >
            Reset
          </button>
        </div>

        <button
          @click="loadGraph"
          class="h-9 px-2.5 text-on-surface-variant hover:text-on-surface bg-surface-container-low border border-outline-variant rounded-lg transition-colors text-xs shadow-sm flex items-center justify-center"
          title="Reload ERD Graph"
        >
          <span class="material-symbols-outlined text-[18px]">refresh</span>
        </button>
      </div>
    </header>

    <!-- Canvas Workspace -->
    <div
      ref="canvasRef"
      class="flex-1 relative overflow-hidden bg-surface-container-lowest cursor-grab active:cursor-grabbing"
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
              font-weight="bold"
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
          class="absolute w-64 bg-surface border rounded-xl shadow-xl overflow-hidden cursor-move transition-all"
          :class="store.selectedTable === node.label ? 'border-primary ring-2 ring-primary/40 shadow-primary/20' : 'border-outline-variant hover:border-primary/60 hover:shadow-lg'"
        >
          <!-- Table Node Header -->
          <div class="px-3.5 py-2 bg-surface-container-low border-b border-outline-variant flex items-center justify-between">
            <div class="flex items-center gap-2 font-mono font-bold text-xs text-on-surface truncate min-w-0">
              <span class="material-symbols-outlined text-primary text-sm flex-shrink-0">table_rows</span>
              <span class="truncate">{{ node.label }}</span>
            </div>
            <div class="flex items-center gap-1 flex-shrink-0">
              <span class="text-[9px] font-mono px-1.5 py-0.5 bg-primary/20 text-primary rounded border border-primary/30 font-bold">
                {{ node.columns?.length || 0 }} cols
              </span>
            </div>
          </div>

          <!-- Column Attributes -->
          <div class="p-2 space-y-1 max-h-56 overflow-y-auto font-mono text-[11px]">
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

          <!-- Node Footer Actions -->
          <div class="px-3 py-1.5 bg-surface-container-low/80 border-t border-outline-variant flex items-center justify-between text-[11px] font-mono">
            <button
              @click.stop="store.selectTable(node.label, 'data')"
              class="hover:text-primary transition-colors flex items-center gap-0.5 text-on-surface-variant font-bold text-[10px]"
            >
              <span class="material-symbols-outlined text-xs">table_chart</span> Data
            </button>
            <button
              @click.stop="store.selectTable(node.label, 'explorer')"
              class="hover:text-primary transition-colors flex items-center gap-0.5 text-on-surface-variant font-bold text-[10px]"
            >
              <span class="material-symbols-outlined text-xs">schema</span> Schema
            </button>
            <button
              @click.stop="store.queryTable(node.label)"
              class="hover:text-primary transition-colors flex items-center gap-0.5 text-on-surface-variant font-bold text-[10px]"
            >
              <span class="material-symbols-outlined text-xs">terminal</span> SQL
            </button>
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
import TableSelector from './TableSelector.vue';

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
        },
      };
    });
  } catch (err) {
    store.addToast(err.message || 'Failed to load ERD graph', 'error');
  }
}

function focusTable({ table }) {
  if (!table) return;
  store.selectedTable = table;
  const targetNode = nodes.value.find((n) => n.label === table || n.id === table);
  if (targetNode) {
    panX.value = -targetNode.position.x * zoom.value + 300;
    panY.value = -targetNode.position.y * zoom.value + 200;
  }
}

const edgesWithCoordinates = computed(() => {
  const nodeMap = new Map();
  nodes.value.forEach((n) => nodeMap.set(n.id, n));

  return edges.value.map((edge) => {
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

function zoomIn() {
  zoom.value = Math.min(zoom.value + 0.15, 2.5);
}
function zoomOut() {
  zoom.value = Math.max(zoom.value - 0.15, 0.4);
}
function resetView() {
  zoom.value = 1;
  panX.value = 60;
  panY.value = 60;
}

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
  nodeDragOffsetX.value = (e.clientX - panX.value) / zoom.value - node.position.x;
  nodeDragOffsetY.value = (e.clientY - panY.value) / zoom.value - node.position.y;
}

function doPanOrDrag(e) {
  if (isPanning.value) {
    panX.value = e.clientX - startPanX.value;
    panY.value = e.clientY - startPanY.value;
  } else if (draggingNode.value) {
    const x = (e.clientX - panX.value) / zoom.value - nodeDragOffsetX.value;
    const y = (e.clientY - panY.value) / zoom.value - nodeDragOffsetY.value;
    draggingNode.value.position = { x, y };
  }
}

function endPanOrDrag() {
  isPanning.value = false;
  draggingNode.value = null;
}

watch(() => store.activeConnectionId, loadGraph, { immediate: true });
</script>
