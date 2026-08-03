import { defineStore } from 'pinia';
import { api } from '../api';

export const useAppStore = defineStore('app', {
  state: () => ({
    connections: [],
    activeConnectionId: '',
    currentDatabase: '',
    currentSchema: '',
    databases: [],
    schemas: [],
    tables: [],
    views: [],
    selectedTable: null,
    activeTab: 'explorer', // 'explorer' | 'editor' | 'data' | 'erd' | 'project'
    toasts: [],
    isConnectionsModalOpen: false,
    isSettingsModalOpen: false,
    config: {
      theme: 'dark',
      defaultPageSize: 25,
      serverPort: 8085,
    },
    loading: false,
  }),

  getters: {
    activeConnection(state) {
      return state.connections.find((c) => c.id === state.activeConnectionId) || null;
    },
  },

  actions: {
    addToast(message, type = 'success') {
      const id = Date.now() + Math.random();
      this.toasts.push({ id, message, type });
      setTimeout(() => {
        this.toasts = this.toasts.filter((t) => t.id !== id);
      }, 4000);
    },

    removeToast(id) {
      this.toasts = this.toasts.filter((t) => t.id !== id);
    },

    async loadConfig() {
      try {
        const cfg = await api.getConfig();
        if (cfg) this.config = cfg;
      } catch (err) {
        console.error('Failed to load config:', err);
      }
    },

    async loadConnections() {
      try {
        this.loading = true;
        const res = await api.listConnections();
        this.connections = res.connections || [];
        if (res.activeId) {
          this.activeConnectionId = res.activeId;
        } else if (this.connections.length > 0 && !this.activeConnectionId) {
          await this.setActiveConnection(this.connections[0].id);
        }
        if (this.activeConnectionId) {
          await this.refreshExplorer();
        }
      } catch (err) {
        this.addToast(`Failed to load connections: ${err.message}`, 'error');
      } finally {
        this.loading = false;
      }
    },

    async setActiveConnection(id) {
      try {
        this.loading = true;
        await api.setActiveConnection(id);
        this.activeConnectionId = id;
        this.selectedTable = null;
        await this.refreshExplorer();
        this.addToast('Connection activated', 'success');
      } catch (err) {
        this.addToast(err.message, 'error');
      } finally {
        this.loading = false;
      }
    },

    async refreshExplorer() {
      if (!this.activeConnectionId) return;
      try {
        const dbs = await api.listDatabases(this.activeConnectionId).catch(() => []);
        this.databases = dbs || [];

        const schemas = await api.listSchemas(this.activeConnectionId, this.currentDatabase).catch(() => []);
        this.schemas = schemas || [];
        if (this.schemas.length > 0 && !this.currentSchema) {
          this.currentSchema = this.schemas[0].name;
        }

        const tbls = await api.listTables(this.activeConnectionId, this.currentSchema).catch(() => []);
        this.tables = (tbls || []).filter(t => t.type !== 'VIEW');
        this.views = (tbls || []).filter(t => t.type === 'VIEW');
      } catch (err) {
        console.error('Failed to refresh explorer:', err);
      }
    },

    selectTable(table) {
      this.selectedTable = table;
      if (this.activeTab !== 'data' && this.activeTab !== 'explorer') {
        this.activeTab = 'data';
      }
    },
  },
});
