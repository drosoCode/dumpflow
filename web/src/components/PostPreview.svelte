<script>
    import { api, parseHtml } from '../Utils';
    import { Link } from "svelte-routing";
    export let site;
    export let post;

    let data = {};
    let title = "";
    let userID = 0;
    let user = "";
    let votes = 0;
    let tags = [];


    const update = (postID) => {
      api(site+"/post/"+postID).then((ret) => {
        if(ret !== null) {
          data = ret;
          tags = data.tags.replaceAll("&lt;", "").split("&gt;");
        }
      })
      
      api(site+"/post/"+postID+"/history").then((ret) => {
        if(ret !== null) {
          let titleDate = 0
          ret.forEach(e => {
            if(e.postHistoryTypeID == 1) {
              userID = e.userID
              api(site+"/user/"+userID).then((d) => {
                if(d !== null)
                  user = d.displayName
              })
            }
            let d = Date.parse(e.creationDate)
            if(d > titleDate && [1,4,7].includes(e.postHistoryTypeID)) {
              titleDate = d
              title = e.text
            }
          });
        }
      })
      
      api(site+"/post/"+postID+"/votes").then((ret) => {
        if(ret !== null) {
          votes = (ret[2] || 0) - (ret[3] || 0);
        }
      })
    };

    $: update(post);
</script>

<main>
  <div class="mt-4 mx-4">
    {#if user !== ""}
    <div class="card bg-dark text-white mb-1 d-flex flex-row">
      <div class="flex-1 d-flex flex-column justify-content-evenly px-2">
        <span class="badge bg-secondary">{votes} Votes</span>
        {#if data.answerCount > 0}
          {#if data.acceptedAnswerID > 0}
            <span class="badge bg-success">{data.answerCount} Answers</span>
          {:else}
            <span class="badge bg-secondary">{data.answerCount} Answers</span>
          {/if}
          <span class="badge bg-secondary">{data.viewCount} Views</span>
        {/if}
      </div>
      <div class="flex-grow-1">
        <div class="card-header"><Link class="text-white text-decoration-none" to={"/site/"+site+"/post/"+post}>{@html parseHtml(title)}</Link></div>
        <div class="card-body preview">
            {@html parseHtml(data.body)}
        </div>
        <div class="card-footer d-flex justify-content-between">
            <div>
              {#each tags as tag}
                <span class="badge bg-primary">{tag}</span>&nbsp;
              {/each}
            </div>
            <div>Asked on {data.creationDate.substring(0,10)} by <Link class="text-decoration-none" to={"/site/"+site+"/user/"+userID}>{user}</Link></div>
        </div>
      </div>
    </div>
    {/if}
  </div>
</main>

<style>
.card-body {
  text-overflow: ellipsis;
  overflow:hidden;
  max-height: 200px;
}
</style>