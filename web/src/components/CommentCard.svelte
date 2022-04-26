<script>
    import { Link } from "svelte-routing";
    import { api, parseHtml } from '../Utils';
    
    export let site;
    export let comment;

    let user = {};
    
    const update = (commentData) => {
      api(site+"/user/"+commentData.userID).then((ret) => {
        if(ret !== null) {
          user = ret;
        }
      })
    };

    $: update(comment);
</script>

<div class="mb-2" style="background-color: #242627; font-size: 13px;">
    {#if comment.score > 0}{comment.score} - {/if}
    {@html parseHtml(comment.text)} - <Link to={"/site/"+site+"/user/"+comment.userID}>{user.displayName}</Link> 
    <span class="text-muted">{comment.creationDate.substring(0,19).replace("T", " ")}</span>
</div>