<script lang="ts">
  // Define types for your data
  interface TableItem {
    id: number;
    title: string;
    body: string;
  }

  // 1. State Runes for pagination and data
  let items = $state<TableItem[]>([]);
  let currentPage = $state(1);
  let limit = $state(5);
  let totalItems = $state(0);
  let loading = $state(false);

  // 2. Derived Rune for total pages
  const totalPages = $derived(Math.ceil(totalItems / limit));

  // 3. Effect Rune to fetch data whenever currentPage changes
  $effect(() => {
    fetchData(currentPage);
  });

  async function fetchData(page: number) {
    loading = true;
    try {
      // Example API using skip/limit logic
      const skip = (page - 1) * limit;
      const response = await fetch(
        `jsonplaceholder.typicode.com{skip}&_limit=${limit}`
      );
      
      if (!response.ok) throw new Error('Fetch failed');
      
      // Get total count (often found in headers or separate API call)
      const totalCountHeader = response.headers.get('x-total-count');
      totalItems = totalCountHeader ? parseInt(totalCountHeader) : 100;
      
      items = await response.json();
    } catch (err) {
      console.error(err);
    } finally {
      loading = false;
    }
  }

  function goToPage(page: number) {
    if (page >= 1 && page <= totalPages) {
      currentPage = page;
    }
  }
</script>

<!-- 4. Table UI -->
<div class="table-container">
  {#if loading}
    <p>Loading...</p>
  {:else}
    <table>
      <thead>
        <tr>
          <th>ID</th>
          <th>Title</th>
        </tr>
      </thead>
      <tbody>
        {#each items as item (item.id)}
          <tr>
            <td>{item.id}</td>
            <td>{item.title}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}

  <!-- 5. Pagination Controls -->
  <div class="pagination">
    <button 
      onclick={() => goToPage(currentPage - 1)} 
      disabled={currentPage === 1 || loading}
    >
      Previous
    </button>
    
    <span>Page {currentPage} of {totalPages}</span>
    
    <button 
      onclick={() => goToPage(currentPage + 1)} 
      disabled={currentPage === totalPages || loading}
    >
      Next
    </button>
  </div>
</div>