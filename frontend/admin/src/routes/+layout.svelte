<script lang="ts">
  import '../app.css';
  import { Link, PageContainer } from '@dinnerdonebetter/ui';

  let { children } = $props();
</script>

<div class="layout-root">
  <header class="layout-header">
    <div class="container">
      <div><Link href="/">Admin</Link></div>
      <nav class="layout-links">
        <Link href="/">Dashboard</Link>
      </nav>
    </div>
  </header>

  <main class="layout-main">
    <PageContainer>
      {@render children()}
    </PageContainer>
  </main>
</div>

<style>
  .layout-root {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
  }
  .layout-header {
    background: var(--color-surface);
    border-bottom: 1px solid var(--color-border);
    padding: var(--space-md);
  }
  .container {
    max-width: 64rem;
    margin: 0 auto;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .layout-links {
    display: flex;
    gap: var(--space-md);
  }
  .layout-main {
    flex: 1;
    padding: var(--space-lg);
  }
</style>
