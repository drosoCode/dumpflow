<script>
    import { Link } from "svelte-routing";
    import { api, parseHtml } from '../Utils';
    import PostCard from '../components/PostCard.svelte';
    
    export let site;
    export let post;
    
    let data = {};
    let answers = [];
    let related = [];
    let title = "";
    let creationDate = 10000000000000;
    let creationDateTxt = "";
    let siteUrl = "";

    const update = (postID) => {
		  api("site/"+site).then((ret) => {
        siteUrl = ret.link;
      })

      api(site+"/post/"+postID).then((ret) => {
        if(ret !== null) {
          data = ret;
        }
      })

      api(site+"/post/"+postID+"/answers").then((ret) => {
        if(ret !== null) {
            answers = ret;
        }
      })
      
      api(site+"/post/"+postID+"/related").then((ret) => {
        if(ret !== null) {
          related = [];
          ret.forEach(e => {
            api(site+"/post/"+e.relatedPostID+"/history").then((r) => {
              if (r!== null) {
                  let t = "";
                  let titleDate = 0;
                  r.forEach(e => {
                    let d = Date.parse(e.creationDate)
                    if(d > titleDate && [1,4,7].includes(e.postHistoryTypeID)) {
                      titleDate = d
                      t = e.text
                    }
                  })
                  related = [...related, [e.relatedPostID, t]]
              }
            })
          });
        }
      })

      api(site+"/post/"+postID+"/history").then((ret) => {
        if(ret !== null) {
          let titleDate = 0
          ret.forEach(e => {
            let d = Date.parse(e.creationDate)
            if(d > titleDate && [1,4,7].includes(e.postHistoryTypeID)) {
              titleDate = d
              title = e.text
            }
            if(d < creationDate) {
              creationDate = d;
              creationDateTxt = e.creationDate;
            }
          });
        }
      })
    };

    $: update(post);
</script>

<main>
    <div class="mt-4 mx-4">
      <div class="card bg-dark text-white mb-2">
          <div class="card-header">
            {@html parseHtml(title)}
          </div>
          <div class="card-body d-flex justify-content-around">
            <span>Asked on {creationDateTxt.substring(0,10)}</span>
            <span>Viewed: {data.viewCount} times</span>
            <span>Answers: {answers.length}</span>
            <a href={"https://"+siteUrl+"/questions/"+post} target="_blank">View online</a>
          </div>
      </div>

      <div class="d-flex flex-column flex-md-row flex-lg-row flex-xl-row flex-xxl-row">
        <div class={related.length > 0 ? "posts-list" : ""}>
          {#if data.id !== undefined}
            <PostCard site={site} postData={data} answer={false}/>
            {#each answers as answer}
              <PostCard site={site} postData={answer} answer={true}/>
            {/each}
          {/if}
        </div>
        {#if related.length > 0}
          <div class="mt-4 related-links" style="">
            <div class="card bg-dark text-white mb-2">
              <div class="card-header">
                <i class="fa-solid fa-link"></i>&nbsp;&nbsp;Related
              </div>
              <div class="card-body">
                {#each related as r}
                  <Link to={"/site/"+site+"/post/"+r[0]}>{r[1]}</Link><br><br>
                {/each}
              </div>
            </div>
          </div>
        {/if}
      </div>
    </div>
</main>

<style>

@media (min-width: 768px) { 
  .related-links {
    width: 20%;
  }
  .posts-list {
    width: 80%;
  }
}

</style>
