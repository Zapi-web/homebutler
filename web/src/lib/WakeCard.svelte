<script>
  import { onMount } from 'svelte';
  import { getWake, postWake } from './api.js';

  let targets = $state([]);
  let error = $state('');
  let waking = $state({});

  onMount(async () => {
    try {
      targets = await getWake();
    } catch (err) {
      error = err.message;
    }
  });

  async function sendWake(name) {
    waking[name] = true;
    waking = waking;
    try {
      await postWake(name);
      // Show success briefly
      setTimeout(() => {
        waking[name] = false;
        waking = waking;
      }, 2000);
    } catch (err) {
      waking[name] = false;
      waking = waking;
      error = err.message;
    }
  }
</script>

<div class="card">
  <div class="card-header">
    <h2>Wake-on-LAN</h2>
    {#if targets.length > 0}
      <span class="badge">{targets.length} devices</span>
    {/if}
  </div>

  {#if error}
    <p class="error">{error}</p>
  {:else if targets.length === 0}
    <p class="empty">No WoL targets configured</p>
  {:else}
    <div class="wake-list">
      {#each targets as t}
        <div class="wake-row">
          <div class="wake-info">
            <span class="wake-name">{t.name}</span>
            <span class="wake-mac">{t.mac}</span>
          </div>
          <button
            class="wake-btn"
            class:sending={waking[t.name]}
            onclick={() => sendWake(t.name)}
            disabled={waking[t.name]}
          >
            {#if waking[t.name]}
              <span class="spinner"></span> Sent
            {:else}
              Wake
            {/if}
          </button>
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

  .wake-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .wake-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0.5rem 0;
    border-bottom: 1px solid var(--border);
  }

  .wake-row:last-child {
    border-bottom: none;
  }

  .wake-info {
    display: flex;
    flex-direction: column;
    gap: 0.125rem;
  }

  .wake-name {
    font-size: 0.8rem;
    color: var(--text-heading);
    font-weight: 500;
  }

  .wake-mac {
    font-size: 0.7rem;
    color: var(--text-secondary);
    font-family: monospace;
  }

  .wake-btn {
    background: var(--accent);
    color: #fff;
    border: none;
    border-radius: 6px;
    padding: 0.35rem 0.85rem;
    font-size: 0.75rem;
    font-weight: 600;
    cursor: pointer;
    transition: background 0.15s ease, opacity 0.15s ease;
    display: flex;
    align-items: center;
    gap: 0.35rem;
  }

  .wake-btn:hover:not(:disabled) {
    background: color-mix(in srgb, var(--accent) 80%, white);
  }

  .wake-btn:disabled {
    cursor: default;
  }

  .wake-btn.sending {
    background: var(--green);
  }

  .spinner {
    width: 10px;
    height: 10px;
    border: 2px solid rgba(255,255,255,0.3);
    border-top-color: #fff;
    border-radius: 50%;
    animation: spin 0.6s linear infinite;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
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
