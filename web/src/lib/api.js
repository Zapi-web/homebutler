const BASE = '';

async function fetchJSON(path, opts) {
  const res = await fetch(`${BASE}${path}`, opts);
  if (!res.ok) {
    const body = await res.text();
    throw new Error(body || res.statusText);
  }
  return res.json();
}

export function getStatus() {
  return fetchJSON('/api/status');
}

export function getDocker() {
  return fetchJSON('/api/docker');
}

export function getProcesses() {
  return fetchJSON('/api/processes');
}

export function getAlerts() {
  return fetchJSON('/api/alerts');
}

export function getPorts() {
  return fetchJSON('/api/ports');
}

export function getWake() {
  return fetchJSON('/api/wake');
}

export function postWake(name) {
  return fetchJSON(`/api/wake/${encodeURIComponent(name)}`, { method: 'POST' });
}

export function getServers() {
  return fetchJSON('/api/servers');
}

export function getServerStatus(name) {
  return fetchJSON(`/api/servers/${encodeURIComponent(name)}/status`);
}
