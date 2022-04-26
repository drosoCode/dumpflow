<script>
  import { api } from '../Utils';
  import { createEventDispatcher } from 'svelte';
  export let site;
  export let link;
  const dispatch = createEventDispatcher();

  const deleteSite = () => {
    api("site/"+site, "DELETE").then((ret) => {
      if(ret !== null) {
        dispatch('delete');
      }
    })
  }
</script>

<span class="text-dark">
  <button type="button" class="btn btn-danger" data-bs-toggle="modal" data-bs-target="#modal_{site}">
    <i class="fa-solid fa-trash-can"></i>&nbsp;&nbsp;Delete
  </button>

  <div class="modal fade" id="modal_{site}" tabindex="-1" aria-labelledby="modal_{site}_label" aria-hidden="true">
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title" id="modal_{site}_label">Modal title</h5>
          <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
        </div>
        <div class="modal-body">
          Do you really want to delete {link} ?
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button>
          <button type="button" class="btn btn-danger" on:click={deleteSite} data-bs-dismiss="modal">Delete</button>
        </div>
      </div>
    </div>
  </div>
</span>