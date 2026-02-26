<script>
  import { onMount, onDestroy } from 'svelte';
  import { getServers, getServerStatus } from './api.js';

  let servers = $state([]);
  let statuses = $state({});
  let error = $state('');
  let timer;

  async function refresh() {
    try {
      servers = await getServers();
      error = '';
      // Fetch status for each server
      for (const srv of servers) {
        try {
          const status = await getServerStatus(srv.name);
          statuses[srv.name] = { ok: true, data: status };
          statuses = statuses;
        } catch {
          statuses[srv.name] = { ok: false };
          statuses = statuses;
        }
      }
    } catch (err) {
      error = err.message;
    }
  }

  onMount(() => {
    refresh();
    timer = setInterval(refresh, 15000);
  });

  onDestroy(() => clearInterval(timer));

  function cpuOf(name) {
    return statuses[name]?.data?.cpu?.usage_percent ?? null;
  }

  function memOf(name) {
    return statuses[name]?.data?.memory?.usage_percent ?? null;
  }

  function uptimeOf(name) {
    return statuses[name]?.data?.uptime ?? null;
  }
</script>

<div class="card">
  <div class="card-header">
    <h2>Server Overview</h2>
    {#if servers.length > 0}
      <span class="badge">{servers.length} servers</span>
    {/if}
  </div>

  {#if error}
    <p class="error">{error}</p>
  {:else if servers.length === 0}
    <p class="empty">Loading...</p>
  {:else}
    <div class="server-grid">
      {#each servers as srv}
        {@const s = statuses[srv.name]}
        <div class="server-item">
          <div class="server-top">
            <span class="dot" style="background:{s ? (s.ok ? 'var(--green)' : 'var(--red)') : 'var(--text-secondary)'}"></span>
            <span class="server-name">{srv.name}</span>
            <span class="server-type">{srv.local ? 'local' : srv.host}</span>
          </div>
          {#if s?.ok}
            <div class="server-metrics">
              <span class="metric">CPU {cpuOf(srv.name)?.toFixed(0) ?? '—'}%</span>
              <span class="metric">MEM {memOf(srv.name)?.toFixed(0) ?? '—'}%</span>
              <span class="metric">{uptimeOf(srv.name) ?? '—'}</span>
            </div>
          {:else if s}
            <div class="server-metrics">
              <span class="metric offline">offline</span>
            </div>
          {:else}
            <div class="server-metrics">
              <span class="metric loading-text">connecting...</span>
            </div>
          {/if}
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  .card {
    background: var(--bg-card);
    border: 1px solid var(--border);
    border-radius: 8px;
    padding: 1rem 1.25rem;
    transition: border-color 0.2s ease, box-shadow 0.2s ease;
  }

  .card:hover {
    border-color: color-mix(in srgb, var(--accent) 40%, transparent);
    box-shadow: 0 0 12px color-mix(in srgb, var(--accent) 10%, transparent);
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.75rem;
  }

  h2 {
    font-size: 0.875rem;
    font-weight: 600;
    color: var(--text-heading);
  }

  .badge {
    font-size: 0.75rem;
    color: var(--text-secondary);
    background: var(--bg-primary);
    padding: 0.15rem 0.5rem;
    border-radius: 10px;
  }

  .server-grid {
    display: flex;
    gap: 0.75rem;
    flex-wrap: wrap;
  }

  .server-item {
    flex: 1;
    min-width: 180px;
    background: var(--bg-primary);
    border-radius: 6px;
    padding: 0.75rem;
    border: 1px solid var(--border);
  }

  .server-top {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-bottom: 0.5rem;
  }

  .dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    flex-shrink: 0;
  }

  .server-name {
    font-size: 0.8rem;
    font-weight: 600;
    color: var(--text-heading);
  }

  .server-type {
    font-size: 0.7rem;
    color: var(--text-secondary);
    margin-left: auto;
  }

  .server-metrics {
    display: flex;
    gap: 0.75rem;
    flex-wrap: wrap;
  }

  .metric {
    font-size: 0.7rem;
    color: var(--text-secondary);
    font-variant-numeric: tabular-nums;
  }

  .metric.offline {
    color: var(--red);
  }

  .metric.loading-text {
    color: var(--text-secondary);
    opacity: 0.6;
  }

  .error {
    color: var(--red);
    font-size: 0.875rem;
  }

  .empty {
    color: var(--text-secondary);
    font-size: 0.875rem;
  }
</style>
