<script>
    import { api } from '../Utils';
    import { onMount } from 'svelte';
    import { navigate, Link } from "svelte-routing";
    import DeleteModal from '../components/DeleteModal.svelte'

    let sites = [];

	onMount(async () => {
		const res = await api("site");
        if(res != null) {
            sites = res;
        }
	});

</script>

<main>
    <div class="mt-4 mx-4">
        {#each sites as site}
            {#if site.enabled}
            <div class="card bg-dark text-white mb-1">
                <div class="card-body d-flex justify-content-between linkTxt">
                    <Link to={"/site/"+site.dbName} style="text-decoration: none; color:white;">{site.link}</Link>
                    <div>
                        <span class="text-muted">{site.updateDate.substring(0,10)}</span>&nbsp;&nbsp;
                        <button type="button" class="btn btn-primary" on:click={() => {navigate("/site/"+site.dbName)}}>
                            <i class="fa-solid fa-magnifying-glass"></i>&nbsp;&nbsp;Browse
                        </button>
                        <a type="button" class="btn btn-info" href={"https://"+site.link} target="_blank">
                            <i class="fa-solid fa-up-right-from-square"></i>&nbsp;&nbsp;Open
                        </a>
                        <DeleteModal site={site.dbName} link={site.link} on:delete={() => {sites.splice(sites.indexOf(site), 1); sites = sites;}}/>
                    </div>
                </div>
            </div>
            {/if}
        {/each}
    </div>
</main>
