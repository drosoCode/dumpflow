<script>
    import { api } from '../Utils';
    import { onMount } from 'svelte';
    import PostPreview from '../components/PostPreview.svelte';
    import Pagination from '../components/Pagination.svelte';
    import { navigate } from 'svelte-routing';

    export let site;
    let data = [];
    let tags = [];
    let postID = "";

    let searchStr = "";
    let username = "";
    let includeComments = false;
    let includeHistory = false;
    let tag = "";

    let searchResults = {};
    let page = -1;

    const limit = 5;

	onMount(async () => {
		const res = await api("site/"+site);
        if(res != null) {
            data = res;
        }
        
		const ret = await api(site+"/tag");
        if(ret != null) {
            tags = ret;
        }
	});

    const search = () => {
        let lst = [];
        if(tag !== "") {
            lst.push(tag)
        }
        api(site+"/post/search", "POST", {
            search: searchStr,
            username: username,
            comments: includeComments,
            history: includeHistory,
            tags: lst,
            start: page == -1 ? 0 : page,
            limit: limit
        }).then((data) => {
            searchResults = data;
            if(page == -1)
                page = 0;
        })
    }

    const gotoPost = () => {
        if(postID !== "")
            navigate("/site/"+site+"/post/"+postID)
    }

</script>

<main>
    {#if page == -1}
    <div class="mt-4 mx-4">
        <div class="card bg-dark text-white mb-1" on:keydown={(e) => {if (e.key == "Enter") search(); }}>
            <h3 class="card-header">{data.link}</h3>
            <div class="card-body">
                <div class="d-flex flex-column flex-md-row flex-lg-row flex-xl-row flex-xxl-row justify-content-between">
                    <div class="mb-3" style="width: 100%">
                        <input type="text" class="form-control" id="search" placeholder="Search" bind:value={searchStr}>
                        <button class="btn btn-primary mt-5" type="button" style="width: 95%" on:click={search}>
                            <i class="fa-solid fa-magnifying-glass"></i>&nbsp;&nbsp;Search
                        </button>
                    </div>
                    <div>
                        <input type="text" class="form-control" id="user" placeholder="Username" bind:value={username}>
                        <select class="form-select mt-3" id="tags" bind:value={tag}>
                            <option value="">Select Tag</option>
                            {#each tags as tag}
                                <option value={tag.tagName}>{tag.tagName} [{tag.count}]</option>
                            {/each}
                        </select>
                        <div class="form-check mt-3">
                            <input class="form-check-input" type="checkbox" value="" id="search_comments" bind:checked={includeComments}>
                            <label class="form-check-label" for="search_comments">
                                <i class="fa-solid fa-comment"></i>&nbsp;&nbsp;Search in comments
                            </label>
                        </div>
                        <div class="form-check mt-3">
                            <input class="form-check-input" type="checkbox" value="" id="search_histroy" bind:checked={includeHistory}>
                            <label class="form-check-label" for="search_histroy">
                                <i class="fa-solid fa-clock-rotate-left"></i>&nbsp;&nbsp;Search in history
                            </label>
                        </div>
                    </div>
                </div>
                
                <hr style="width: 95%; margin: 5px, auto;">
                
                <form class="row g-3">
                    <div class="col-auto">
                    <input type="text" class="form-control" placeholder="Post ID" bind:value={postID}>
                    </div>
                    <div class="col-auto">
                    <button type="submit" class="btn btn-primary mb-3" on:click={gotoPost}><i class="fa-solid fa-arrow-right"></i>&nbsp;&nbsp;Go</button>
                    </div>
                </form>
            </div>            
        </div>
    </div>
    {:else}
        {#if searchResults.results > 0}
        {#each searchResults.posts as post}
            <PostPreview site={site} post={post}/>
        {/each}
        <div style="position:fixed;bottom:0;width: 100%;">
            <Pagination bind:page={page}  max={searchResults.results} limit={limit} on:click={search} />
        </div>
        {:else}
        
        <div class="mt-4 mx-4">
            <div class="card bg-dark text-white mb-2">
                <div class="card-body">
                    <h3>No Results</h3>
                </div>
            </div>
        </div>
        {/if}
    {/if}
</main>