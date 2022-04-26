<script>
    import { onMount } from 'svelte';
    import { Link } from "svelte-routing";
    import { api, parseHtml } from '../Utils';
    
    export let site;
    export let post;
    
    let data = [];
    let title = "";
    // https://meta.stackexchange.com/questions/2677/database-schema-documentation-for-the-public-data-dump-and-sede
    const actions = {
        2: "creation",
        4: "edit title",
        5: "edit body",
        6: "edit tags",
        7: "rollback title",
        8: "rollback body",
        9: "rollback tags",
        10: "post closed",
        11: "post reopened",
        12: "post deleted",
        13: "post undeleted",
        14: "post locked",
        15: "post unlocked",
        16: "community owned",
        18: "question merged",
        19: "question protected",
        20: "question unprotected",
        22: "question unmerged",
        23: "suggested edit applied",
        25: "tweeted",
        31: "comment discussion moved to chat",
        33: "post notice added",
        34: "post notice removed",
        50: "bumped by community user",
        52: "question became hot network question",
        53: "question removed from hot network question"
    };
    const actKeys = Object.keys(actions);

    onMount(() => {
      api(site+"/post/"+post).then((ret) => {
        let parent = post;
        if(ret !== null && ret.parentID > 0) {
            parent = ret.parentID;
        }
        api(site+"/post/"+parent+"/history").then((r) => {
          r.forEach(e => {
            if(e.postHistoryTypeID === 1) {
                title = e.text;
                return;
            }
          });
        });
      });

        
      api(site+"/post/"+post+"/history").then((ret) => {
        if(ret !== null) {
          ret.forEach(e => {
            if(actKeys.includes(e.postHistoryTypeID.toString())) {
                api(site+"/user/"+e.userID).then((r) => {
                    if(ret !== null) {
                        e.user = r.displayName;
                        data = [...data, e];
                    }
                })
            }
          });
        }
      })
    });
</script>

<main>
    <div class="mt-4 mx-4">
      <div class="card bg-dark text-white mb-2">
          <div class="card-header">
            History for {@html parseHtml(title)}
          </div>
          <div class="card-body">
            <table class="table table-dark table-striped table-hover">
                <thead>
                    <tr>
                        <th scope="col">Date</th>
                        <th scope="col">Action</th>
                        <th scope="col">User</th>
                        <th scope="col">License</th>
                        <th scope="col">Comment</th>
                        <th scope="col">Text</th>
                    </tr>
                </thead>
                <tbody>
                    {#each data as h}
                    <tr>
                        <td>{h.creationDate.substring(0,16).replace('T', ' ')}</td>
                        <td>{actions[h.postHistoryTypeID]}</td>
                        <td><Link to={"/site/"+site+"/user/"+h.userID}>{h.user}</Link></td>
                        <td>{h.contentLicense}</td>
                        <td>{h.comment}</td>
                        <td>{@html parseHtml(h.text)}</td>
                    </tr>
                    {/each}
                </tbody>
            </table>
          </div>
      </div>
</main>
