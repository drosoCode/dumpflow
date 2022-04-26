<script>
    import { api } from '../Utils';
    import { onMount, onDestroy } from 'svelte';

    let imports = {};
    let downloads = {};
    let status = [];
    let interval = null;

    const nbFinished = (data) => {
        return (data.badges.current >= data.badges.total ? 1 : 0) + (data.comments.current >= data.comments.total ? 1 : 0) + (data.postHistory.current >= data.postHistory.total ? 1 : 0) + (data.postLinks.current >= data.postLinks.total ? 1 : 0) + (data.posts.current >= data.posts.total ? 1 : 0) + (data.tags.current >= data.tags.total ? 1 : 0) + (data.users.current >= data.users.total ? 1 : 0) + (data.votes.current >= data.votes.total ? 1 : 0) + (data.unzipping.current >= data.unzipping.total ? 1 : 0)
    }

    const update = () => {
		api("site/status").then((res) => {
            if(res !== null) {
                status = Object.keys(res);
                for (const [key, value] of Object.entries(res)) {
                    if(value) {
                        api("site/import/status?path="+key).then((data) => {
                            if(data !== null) {
                                imports[key] = data
                            }
                        })
                    }
                    else {
                        api("site/download/status?site="+key).then((data) => {
                            if(data !== null) {
                                downloads[key] = data
                            }
                        })
                    }
                }
                for(const x of Object.keys(imports)) {
                    if(nbFinished(imports[x]) === 8 || !status.includes(x)) {
                        delete imports[x];
                        imports = imports;
                    }
                }
                for(const x of Object.keys(downloads)) {
                    if(downloads[x].finishedFiles === downloads[x].totalFiles || !status.includes(x)) {
                        delete downloads[x];
                        downloads = downloads;
                    }
                }
            }
        })
	}

	onMount(() => {
        update();
        interval = setInterval(update, 1000)
    });

    onDestroy(() => clearInterval(interval));
</script>

<main>
    <div class="mt-4 mx-4">
    
        <div class="alert alert-primary" role="alert">
          <i class="fa-solid fa-download"></i>&nbsp;&nbsp;Downloads
        </div>

        {#each Object.keys(downloads) as i}
            <div class="card bg-dark text-white mb-2">
                <div class="card-header">{i}</div>
                <div class="card-body">
                    {#if downloads[i].totalFiles > 1}
                    <div class="progress mb-1">
                      <div class="progress-bar bg-success progress-bar-striped progress-bar-animated" role="progressbar" style="width: {(downloads[i].finishedFiles+1)/(downloads[i].totalFiles+1)*100}%;">{downloads[i].finishedFiles+1}/{downloads[i].totalFiles}</div>
                    </div>
                    {/if}
                    <div class="progress">
                      <div class="progress-bar progress-bar-striped progress-bar-animated" role="progressbar" style="width: {downloads[i].currentFileDownloadedSize/downloads[i].currentFileTotalSize*100}%;">{downloads[i].currentFile}: {Math.round(downloads[i].currentFileDownloadedSize/downloads[i].currentFileTotalSize*100)}%</div>
                    </div>
                </div>
            </div>
        {/each}

        <div class="alert alert-primary mt-2" role="alert">
          <i class="fa-solid fa-file-import"></i>&nbsp;&nbsp;Imports
        </div>

        {#each Object.keys(imports) as i}
            <div class="card bg-dark text-white mb-2">
                <div class="card-header">{i}</div>
                <div class="card-body">
                    <div class="progress mb-1">
                      <div class="progress-bar bg-success progress-bar-striped progress-bar-animated" role="progressbar" style="width: {nbFinished(imports[i])/9*100}%;">{nbFinished(imports[i])}/9</div>
                    </div>
                    <div class="progress mb-1">
                      <div class="progress-bar progress-bar-striped progress-bar-animated" role="progressbar" style="width: {imports[i].unzipping.current/imports[i].unzipping.total*100}%;">Unzipping: {imports[i].unzipping.current}/{imports[i].unzipping.total}</div>
                    </div>
                    <div class="progress mb-1">
                      <div class="progress-bar progress-bar-striped progress-bar-animated" role="progressbar" style="width: {imports[i].badges.current/imports[i].badges.total*100}%;">Badges: {Math.round(imports[i].badges.current/imports[i].badges.total*100)}%</div>
                    </div>
                    <div class="progress mb-1">
                      <div class="progress-bar progress-bar-striped progress-bar-animated" role="progressbar" style="width: {imports[i].comments.current/imports[i].comments.total*100}%;">Comments: {Math.round(imports[i].comments.current/imports[i].comments.total*100)}%</div>
                    </div>
                    <div class="progress mb-1">
                      <div class="progress-bar progress-bar-striped progress-bar-animated" role="progressbar" style="width: {imports[i].postHistory.current/imports[i].postHistory.total*100}%;">Post History: {Math.round(imports[i].postHistory.current/imports[i].postHistory.total*100)}%</div>
                    </div>
                    <div class="progress mb-1">
                      <div class="progress-bar progress-bar-striped progress-bar-animated" role="progressbar" style="width: {imports[i].postLinks.current/imports[i].postLinks.total*100}%;">Post Links: {Math.round(imports[i].postLinks.current/imports[i].postLinks.total*100)}%</div>
                    </div>
                    <div class="progress mb-1">
                      <div class="progress-bar progress-bar-striped progress-bar-animated" role="progressbar" style="width: {imports[i].posts.current/imports[i].posts.total*100}%;">Posts: {Math.round(imports[i].posts.current/imports[i].posts.total*100)}%</div>
                    </div>
                    <div class="progress mb-1">
                      <div class="progress-bar progress-bar-striped progress-bar-animated" role="progressbar" style="width: {imports[i].tags.current/imports[i].tags.total*100}%;">Tags: {Math.round(imports[i].tags.current/imports[i].tags.total*100)}%</div>
                    </div>
                    <div class="progress mb-1">
                      <div class="progress-bar progress-bar-striped progress-bar-animated" role="progressbar" style="width: {imports[i].users.current/imports[i].users.total*100}%;">Users: {Math.round(imports[i].users.current/imports[i].users.total*100)}%</div>
                    </div>
                    <div class="progress mb-1">
                      <div class="progress-bar progress-bar-striped progress-bar-animated" role="progressbar" style="width: {imports[i].votes.current/imports[i].votes.total*100}%;">Votes: {Math.round(imports[i].votes.current/imports[i].votes.total*100)}%</div>
                    </div>
                </div>
            </div>
        {/each}
    </div>
</main>