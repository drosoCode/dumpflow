<script>
    import { Link } from "svelte-routing";
    import { api } from '../Utils';
    
    export let site;
    export let post;
    export let author;
    export let answer;

    let user = {};
    let userID = 0;
    let date = "";
    let label = "";
    let badges = [0,0,0];

    $: {
      if(author) {
        label = answer ? "answered" : "asked";
      } else {
        label = "edited";
      }
    }
    
    const update = (postID) => {
      api(site+"/post/"+postID+"/history").then((ret) => {
        if(ret !== null) {
          let editDate = 0
          ret.forEach(e => {
            if(author) {
                // search for author
                if(e.postHistoryTypeID == 2) { // if initial body
                    userID = e.userID;
                    date = e.creationDate;
                }
            } else {
                // search for last edit
                let d = Date.parse(e.creationDate)
                if(d > editDate && e.postHistoryTypeID == 5) {
                    editDate = d
                    userID = e.userID;
                    date = e.creationDate;
                }
            }
          });

          if(userID !== 0) {
            api(site+"/user/"+userID).then((ret) => {
                if(ret !== null) {
                user = ret;
                }
            })
            api(site+"/user/"+userID+"/badges").then((ret) => {
                if(ret !== null) {
                  ret.forEach(e => {
                    if(e.class === 1)
                      badges[0]++;
                    else if(e.class === 2)
                      badges[1]++;
                    else if(e.class === 3)
                      badges[2]++;
                  });
                }
            })
          }
        }
      })
    };

    $: update(post);
</script>

<div class="mb-2 me-3 usercard">
  {#if userID !== 0}
    <span class="text-muted">{label} - {date.substring(0,19).replace("T", " ")}</span>
    <div class="d-flex mt-2">
      <img src={user.profileImageUrl} width="40px" height="40px" class="me-2"/>
      <div class="d-flex flex-column">
        <Link to={"/site/"+site+"/user/"+userID}>{user.displayName}</Link>
        <div>
          <span class="fw-bold">{user.reputation}</span>
          {#if badges[0] > 0}<span class="badge" style="background-color:#FFD700">{badges[0]}</span>{/if}
          {#if badges[1] > 0}<span class="badge" style="background-color:#C0C0C0">{badges[1]}</span>{/if}
          {#if badges[2] > 0}<span class="badge" style="background-color:#6A5211">{badges[2]}</span>{/if}
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  .usercard {
    background-color: #242627;
    font-size: 13px;
    max-width: 200px;
  }
</style>