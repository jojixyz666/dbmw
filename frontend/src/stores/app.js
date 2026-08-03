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
    pendingSqlQuery: null, // query to load into SQL console when navigating
    activeTab: 'explorer', // 'explorer' | 'editor' | 'data' | 'erd' | 'project'
    toasts: [],
    isConnectionsModalOpen: false,
    isSettingsModalOpen: false,
    config: {
      theme: 'light',
      palette: 'emerald',
      defaultPageSize: 25,
      serverPort: 8085,
    },
    loading: false,
    tablesLoading: false,
  }),

  getters: {
    activeConnection(state) {
      return state.connections.find((c) => c.id === state.activeConnectionId) || null;
    },
    allTables(state) {
      const tableList = (state.tables || []).map((t) => ({
        ...t,
        type: 'TABLE',
      }));
      const viewList = (state.views || []).map((v) => ({
        ...v,
        type: 'VIEW',
      }));
      return [...tableList, ...viewList];
    },
    selectedTableDetails(state) {
      if (!state.selectedTable) return null;
      const found = state.tables.find((t) => t.name === state.selectedTable);
      if (found) return { ...found, type: 'TABLE' };
      const foundView = state.views.find((v) => v.name === state.selectedTable);
      if (foundView) return { ...foundView, type: 'VIEW' };
      return { name: state.selectedTable, type: 'TABLE' };
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

    applyTheme() {
      const theme = this.config.theme || 'light';
      const palette = this.config.palette || 'emerald';
      document.documentElement.className = `${theme} palette-${palette}`;
    },

    async loadConfig() {
      try {
        const cfg = await api.getConfig();
        if (cfg) {
          this.config = { ...this.config, ...cfg };
          if (!cfg.theme) this.config.theme = 'light';
          if (!cfg.palette) this.config.palette = 'emerald';
        }
        this.applyTheme();
      } catch (err) {
        console.error('Failed to load config:', err);
        this.applyTheme();
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
        this.currentDatabase = '';
        this.currentSchema = '';
        await this.refreshExplorer();
        this.addToast('Connection activated', 'success');
      } catch (err) {
        this.addToast(err.message, 'error');
      } finally {
        this.loading = false;
      }
    },

    async setDatabase(databaseName) {
      this.currentDatabase = databaseName;
      this.currentSchema = '';
      this.selectedTable = null;
      await this.refreshExplorer();
    },

    async setSchema(schemaName) {
      this.currentSchema = schemaName;
      this.selectedTable = null;
      await this.refreshTables();
    },

    async refreshExplorer() {
      if (!this.activeConnectionId) return;
      this.tablesLoading = true;
      try {
        const dbs = await api.listDatabases(this.activeConnectionId).catch(() => []);
        this.databases = dbs || [];

        const schemas = await api.listSchemas(this.activeConnectionId, this.currentDatabase).catch(() => []);
        this.schemas = schemas || [];

        const schemaNames = this.schemas.map((s) => s.name);
        if (!this.currentSchema || !schemaNames.includes(this.currentSchema)) {
          if (schemaNames.includes('public')) {
            this.currentSchema = 'public';
          } else if (this.activeConnection?.database && schemaNames.includes(this.activeConnection.database)) {
            this.currentSchema = this.activeConnection.database;
          } else if (schemaNames.includes('main')) {
            this.currentSchema = 'main';
          } else if (schemaNames.length > 0) {
            const nonSystem = schemaNames.find(
              (s) => !['information_schema', 'pg_catalog', 'pg_toast', 'mysql', 'performance_schema', 'sys'].includes(s.toLowerCase())
            );
            this.currentSchema = nonSystem || schemaNames[0];
          } else {
            this.currentSchema = this.activeConnection?.database || '';
          }
        }

        await this.refreshTables();
      } catch (err) {
        console.error('Failed to refresh explorer:', err);
      } finally {
        this.tablesLoading = false;
      }
    },

    async refreshTables() {
      if (!this.activeConnectionId) return;
      this.tablesLoading = true;
      try {
        let tbls = await api.listTables(this.activeConnectionId, this.currentSchema).catch(() => []);
        
        // Fallback: If 0 tables found with current schema and schema isn't standard public/main, try default schema
        if ((!tbls || tbls.length === 0) && this.currentSchema && this.currentSchema !== 'public' && this.currentSchema !== 'main') {
          const fallbackTbls = await api.listTables(this.activeConnectionId, '').catch(() => []);
          if (fallbackTbls && fallbackTbls.length > 0) {
            tbls = fallbackTbls;
          }
        }

        this.tables = (tbls || []).filter((t) => {
          const type = (t.type || '').toUpperCase();
          return !type.includes('VIEW');
        });
        this.views = (tbls || []).filter((t) => {
          const type = (t.type || '').toUpperCase();
          return type.includes('VIEW');
        });
      } catch (err) {
        console.error('Failed to refresh tables:', err);
      } finally {
        this.tablesLoading = false;
      }
    },

    selectTable(table, targetTab = null) {
      if (!table) {
        this.selectedTable = null;
        return;
      }
      const tableName = typeof table === 'string' ? table : table.name;
      this.selectedTable = tableName;
      if (targetTab) {
        this.activeTab = targetTab;
      } else if (this.activeTab !== 'data' && this.activeTab !== 'explorer') {
        this.activeTab = 'data';
      }
    },

    queryTable(table, queryType = 'select', limit = 50) {
      const tableName = typeof table === 'string' ? table : table.name;
      this.selectedTable = tableName;
      let sql = '';
      if (queryType === 'count') {
        sql = `SELECT COUNT(*) AS total FROM ${tableName};`;
      } else if (queryType === 'desc') {
        sql = `SELECT * FROM ${tableName} ORDER BY 1 DESC LIMIT ${limit};`;
      } else {
        sql = `SELECT * FROM ${tableName} LIMIT ${limit};`;
      }
      this.pendingSqlQuery = sql;
      this.activeTab = 'editor';
    },

    clearSelectedTable() {
      this.selectedTable = null;
    },
  },
});