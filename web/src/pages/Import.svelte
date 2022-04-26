<script>
    import { api } from '../Utils';
    import { onMount } from 'svelte';
    import { navigate } from "svelte-routing";

    let sites = [];

	onMount(async () => {
		const res = await api("site/import");
        if(res != null) {
            sites = res;
        }
	});

    const importSite = (path) => {
        api("site/import?path="+path, "POST").then((ret) => {
            if(ret != null) {
                navigate("/status", { replace: true });
            }
        })
    }
</script>

<main>
    <div class="mt-4 mx-4">
        {#each sites as site}
            <div class="card bg-dark text-white mb-1">
                <div class="card-body d-flex justify-content-between">
                    {site} <button type="button" class="btn btn-primary" on:click={() => {importSite(site)}}><i class="fa-solid fa-file-import"></i>&nbsp;&nbsp;Import</button>
                </div>
            </div>
        {/each}
    </div>
</main>