<script>
    import { onMount } from 'svelte';
    import { api, parseHtml } from '../Utils';
    
    export let site;
    export let user;
    
    let data = {};
    let siteUrl = "";

    onMount(() => {
	  api("site/"+site).then((ret) => {
        siteUrl = ret.link;
      })

      api(site+"/user/"+user).then((ret) => {
        if(ret !== null) {
          data = ret;
        }
      })
    });
</script>

<main>
    {#if data.id !== undefined}
      <div class="mt-4 mx-4">
        <div class="card bg-dark text-white mb-2">
            <div class="card-header d-flex justify-content-around">
                <img src={data.profileImageUrl} width="50px" height="50px"/>
                <h3>{data.displayName}</h3>
                <a href={"https://"+siteUrl+"/users/"+user} target="_blank">View online</a>
            </div>
            <div class="card-body">
                {#if data.aboutMe !== ""}
                <div class="alert" style="background-color: #202020;">
                    {@html parseHtml(data.aboutMe)}
                </div>
                {/if}
                
                <table class="table table-dark table-striped table-hover">
                    <thead>
                        <tr>
                            <th scope="col">Name</th>
                            <th scope="col">Value</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <td>Creation Date</td>
                            <td>{data.creationDate.substring(0,10)}</td>
                        </tr>
                        <tr>
                            <td>Last Access Date</td>
                            <td>{data.lastAccessDate.substring(0,10)}</td>
                        </tr>
                        <tr>
                            <td>Location</td>
                            <td>{data.location}</td>
                        </tr>
                        <tr>
                            <td>Reputation</td>
                            <td>{data.reputation}</td>
                        </tr>
                        <tr>
                            <td>Views</td>
                            <td>{data.views}</td>
                        </tr>
                        <tr>
                            <td>Upvotes</td>
                            <td>{data.upvotes}</td>
                        </tr>
                        <tr>
                            <td>Downvotes</td>
                            <td>{data.downvotes}</td>
                        </tr>
                        {#if data.websiteUrl !== ""}
                        <tr>
                            <td>Website</td>
                            <td><a href={data.websiteUrl} target="_blank">website</a></td>
                        </tr>
                        {/if}
                    </tbody>
                </table>
            </div>
        </div>
      </div>
    {/if}
</main>
