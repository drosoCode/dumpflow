<script>
    import { createEventDispatcher } from 'svelte';
    export let page; // current page
    export let max; // total number of results
    export let limit; // results per pages
    const dispatch = createEventDispatcher();

    let maxPage = Math.round(max/limit)

    const click = (nb) => {
        page = nb;
        dispatch('click');
    }
</script>

<main class="d-flex justify-content-around">
    <nav>
        <ul class="pagination">
          {#if page-1 >= 0}
          <li class="page-item">
            <button class="page-link" on:click={() => {click(page-1)}}>Previous</button>
          </li>
          {/if}

          
          {#if page-3 >= 0}
          <li class="page-item"><button class="page-link" on:click={() => {click(0)}}>{0}</button></li>
          <li class="page-item disabled"><button class="page-link">...</button></li>
          {/if}


          {#if page-2 >= 0}
          <li class="page-item"><button class="page-link" on:click={() => {click(page-2)}}>{page-2}</button></li>
          {/if}
          {#if page-1 >= 0}
          <li class="page-item"><button class="page-link" on:click={() => {click(page-1)}}>{page-1}</button></li>
          {/if}

          <li class="page-item active"><button class="page-link" on:click={() => {click(page)}}>{page}</button></li>

          {#if page+1 <= maxPage}
          <li class="page-item"><button class="page-link" on:click={() => {click(page+1)}}>{page+1}</button></li>
          {/if}
          {#if page+2 <= maxPage}
          <li class="page-item"><button class="page-link" on:click={() => {click(page+2)}}>{page+2}</button></li>
          {/if}

          {#if page+3 < maxPage}
          <li class="page-item disabled"><button class="page-link">...</button></li>
          <li class="page-item"><button class="page-link" on:click={() => {click(maxPage)}}>{maxPage}</button></li>
          {/if}
          
          {#if page+1 < maxPage}        
          <li class="page-item">
            <button class="page-link" on:click={() => {click(page+1)}}>Next</button>
          </li>
          {/if}
        </ul>
      </nav>
</main>