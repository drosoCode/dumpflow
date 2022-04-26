<script>
    import { onMount } from 'svelte';
    import { Link } from "svelte-routing";
    import { api, parseHtml } from '../Utils';
    import CommentCard from './CommentCard.svelte';
    import UserCard from './UserCard.svelte';
    
    export let site;
    export let postData;
    export let answer;

    let votes = 0;
    let acceptedAnswer = false;
    let favs = 0;
    let tags = [];
    let comments1 = [];
    let comments2 = [];

    onMount(() => {
      tags = postData.tags.replaceAll("&lt;", "").split("&gt;");
      
      api(site+"/post/"+postData.id+"/votes").then((ret) => {
        if(ret !== null) {
          votes = (ret[2] || 0) - (ret[3] || 0);
          favs = ret[5] || 0;
          acceptedAnswer = ret[1] !== undefined;
        }
      })
      
      api(site+"/post/"+postData.id+"/comments").then((ret) => {
        if(ret !== null) {
          if(ret.length > 2) {
            comments1 = ret.slice(0, 2);
            comments2 = ret.slice(3);
          } else {
            comments1 = ret;
          }
        }
      })
    });
</script>

<main>
    <div class="mt-4 mx-4">
      <div class="card bg-dark text-white mb-1 d-flex flex-row">
        <div class="flex-1 d-flex flex-column justify-content-evenly px-2">
          {#if acceptedAnswer}
            <span class="badge bg-success">{votes} Votes</span>
          {:else}
            <span class="badge bg-secondary">{votes} Votes</span>
          {/if}
          {#if favs > 0}
            <span class="badge bg-secondary">{favs}&nbsp;&nbsp;<i class="fa-solid fa-bookmark"></i></span>
          {/if}
          <Link class="btn btn-secondary" to={"/site/"+site+"/post/"+postData.id+"/history"}><i class="fa-solid fa-clock-rotate-left"></i></Link>
        </div>
        <div class="flex-grow-1">
          <div class="card-body">
            <div class="mb-2">
              {@html parseHtml(postData.body)}
            </div>
            <div class="d-flex">
              <div class="flex-grow-1">
                <!-- tags, if question -->
                {#if !answer}
                <div class="mt-4">
                  {#each tags as tag}
                  <span class="badge bg-primary">{tag}</span>&nbsp;
                  {/each}
                </div>
                {/if}
              </div>
              <!--author cards-->
              <UserCard post={postData.id} site={site} author={false} answer={answer}/>
              <UserCard post={postData.id} site={site} author={true} answer={answer}/>
            </div>
          </div>
          <div class="card-footer d-flex justify-content-between">
            <!-- comments -->
              <div class="ms-4 me-1">
                {#each comments1 as comment}
                  <CommentCard comment={comment} site={site}/>
                {/each}
                {#if comments2.length > 0}
                  <button class="btn btn-link" on:click={() => {comments1 = [...comments1, ...comments2]; comments2 = [];}}>Show More</button>
                {/if}
              </div>
          </div>
        </div>
      </div>
    </div>
</main>
