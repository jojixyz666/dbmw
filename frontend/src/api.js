const BASE_URL = '/api';

export async function request(endpoint, options = {}) {
  const url = `${BASE_URL}${endpoint}`;
  const config = {
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    ...options,
  };

  const response = await fetch(url, config);
  
  // Handle file downloads
  const contentType = response.headers.get('content-type');
  if (contentType && (contentType.includes('text/csv') || contentType.includes('application/octet-stream'))) {
    return response.blob();
  }

  const data = await response.json().catch(() => ({}));
  if (!response.ok) {
    throw new Error(data.error || `HTTP error ${response.status}`);
  }
  return data;
}

export const api = {
  // Config
  getConfig: () => request('/config'),
  saveConfig: (cfg) => request('/config', { method: 'POST', body: JSON.stringify(cfg) }),

  // Connections
  listConnections: () => request('/connections'),
  getConnection: (id) => request(`/connections/${id}`),
  saveConnection: (cfg) => request('/connections', { method: 'POST', body: JSON.stringify(cfg) }),
  createConnection: (cfg) => request('/connections', { method: 'POST', body: JSON.stringify(cfg) }),
  updateConnection: (id, cfg) => request('/connections', { method: 'POST', body: JSON.stringify(cfg) }),
  deleteConnection: (id) => request(`/connections/${id}`, { method: 'DELETE' }),
  testConnection: (cfg) => request('/connections/test', { method: 'POST', body: JSON.stringify(cfg) }),
  setActiveConnection: (id) => request('/connections/active', { method: 'POST', body: JSON.stringify({ id }) }),

  // Explorer
  listDatabases: (connId) => request(`/explorer/databases?connId=${connId || ''}`),
  listSchemas: (connId, database) => request(`/explorer/schemas?connId=${connId || ''}&database=${encodeURIComponent(database || '')}`),
  listTables: (connId, schema) => request(`/explorer/tables?connId=${connId || ''}&schema=${encodeURIComponent(schema || '')}`),
  getTableDetails: (connId, schema, table) => {
    const actualTable = table !== undefined ? table : schema;
    const actualSchema = table !== undefined ? schema : '';
    return request(`/explorer/tables/${encodeURIComponent(actualTable)}/details?connId=${connId || ''}&schema=${encodeURIComponent(actualSchema || '')}`);
  },
  getTableSchema: (connId, schema, table) => {
    const actualTable = table !== undefined ? table : schema;
    const actualSchema = table !== undefined ? schema : '';
    return request(`/explorer/tables/${encodeURIComponent(actualTable)}/details?connId=${connId || ''}&schema=${encodeURIComponent(actualSchema || '')}`);
  },
  listColumns: (connId, schema, table) => {
    const actualTable = table !== undefined ? table : schema;
    const actualSchema = table !== undefined ? schema : '';
    return request(`/explorer/columns/${encodeURIComponent(actualTable)}?connId=${connId || ''}&schema=${encodeURIComponent(actualSchema || '')}`);
  },
  listIndexes: (connId, schema, table) => {
    const actualTable = table !== undefined ? table : schema;
    const actualSchema = table !== undefined ? schema : '';
    return request(`/explorer/indexes/${encodeURIComponent(actualTable)}?connId=${connId || ''}&schema=${encodeURIComponent(actualSchema || '')}`);
  },
  listForeignKeys: (connId, schema, table) => {
    const actualTable = table !== undefined ? table : schema;
    const actualSchema = table !== undefined ? schema : '';
    return request(`/explorer/foreign-keys/${encodeURIComponent(actualTable)}?connId=${connId || ''}&schema=${encodeURIComponent(actualSchema || '')}`);
  },
  listViews: (connId, schema) => request(`/explorer/views?connId=${connId || ''}&schema=${encodeURIComponent(schema || '')}`),

  // Query
  executeQuery: (connId, query) => request('/query/execute', { method: 'POST', body: JSON.stringify({ connectionId: connId, query }) }),
  getHistory: (connId, limit = 50) => request(`/query/history?connId=${connId || ''}&limit=${limit}`),
  clearHistory: (connId) => request(`/query/history?connId=${connId || ''}`, { method: 'DELETE' }),
  exportCSV: (result) => request('/query/export/csv', { method: 'POST', body: JSON.stringify(result) }),
  exportJSON: (result) => request('/query/export/json', { method: 'POST', body: JSON.stringify(result) }),

  // Data
  browseData: (connId, schema, table, opts = {}) => {
    const actualTable = table !== undefined ? table : schema;
    const actualSchema = table !== undefined ? schema : '';
    return request(`/data/browse/${actualTable}?connId=${connId || ''}&schema=${actualSchema || ''}`, { method: 'POST', body: JSON.stringify(opts) });
  },
  getTableRows: (connId, schema, table, opts = {}) => {
    if (typeof table === 'object' && table !== null) {
      opts = table;
      table = schema;
      schema = '';
    }
    return request(`/data/browse/${table}?connId=${connId || ''}&schema=${schema || ''}`, { method: 'POST', body: JSON.stringify(opts) });
  },
  insertRow: (connId, schema, table, values = {}) => {
    if (typeof table === 'object' && table !== null) {
      values = table;
      table = schema;
      schema = '';
    }
    return request(`/data/insert/${table}?connId=${connId || ''}&schema=${schema || ''}`, { method: 'POST', body: JSON.stringify(values) });
  },
  updateRow: (connId, schema, table, pk = {}, values = {}) => {
    if (typeof table === 'object' && table !== null) {
      const payload = table;
      table = schema;
      schema = '';
      return request(`/data/update/${table}?connId=${connId || ''}&schema=${schema || ''}`, { method: 'POST', body: JSON.stringify(payload) });
    }
    return request(`/data/update/${table}?connId=${connId || ''}&schema=${schema || ''}`, { method: 'POST', body: JSON.stringify({ primaryKey: pk, values }) });
  },
  deleteRow: (connId, schema, table, pk = {}) => {
    if (typeof table === 'object' && table !== null) {
      pk = table;
      table = schema;
      schema = '';
    }
    return request(`/data/delete/${table}?connId=${connId || ''}&schema=${schema || ''}`, { method: 'POST', body: JSON.stringify(pk) });
  },

  // ERD
  generateERD: (connId, schema = '') => request(`/erd/generate?connId=${connId || ''}&schema=${encodeURIComponent(schema || '')}`),
  getErdGraph: (connId, schema = '') => request(`/erd/generate?connId=${connId || ''}&schema=${encodeURIComponent(schema || '')}`),

  // Project
  detectProject: (path) => request(`/project/detect?path=${encodeURIComponent(path || '.')}`),
  generateProjectConfig: (path, config) => request('/project/generate', { method: 'POST', body: JSON.stringify({ path, config }) }),
};
