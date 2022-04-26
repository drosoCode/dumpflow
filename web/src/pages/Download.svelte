<script>
    import { api } from '../Utils';
    import { onMount } from 'svelte';
    import { navigate } from "svelte-routing";
    import Spinner from '../components/Spinner.svelte';

    let sites = [];
    let date = "";

	onMount(async () => {
		const res = await api("site/download");
        if(res != null) {
            sites = Object.keys(res.sites);
            date = res.date;
        }
	});

    const download = (site) => {
        api("site/download?site="+site, "POST").then((ret) => {
            if(ret != null) {
                navigate("/status", { replace: true });
            }
        })
    }
</script>

<main>
    <div class="mt-4 mx-4">
        <div class="alert alert-primary" role="alert">
            <i class="fa-solid fa-arrows-rotate"></i>&nbsp;&nbsp;Last Update: {date.substring(0,10)}
        </div>

        {#each sites as site}
            <div class="card bg-dark text-white mb-1">
                <div class="card-body d-flex justify-content-between">
                    {site} <button type="button" class="btn btn-primary" on:click={() => {download(site)}}><i class="fa-solid fa-download"></i>&nbsp;&nbsp;Download</button>
                </div>
            </div>
        {:else}
            <Spinner/>
        {/each}
    </div>
</main>