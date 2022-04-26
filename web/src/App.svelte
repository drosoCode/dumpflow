<script>
  import { Router, Link, Route } from "svelte-routing";
  import Main from "./pages/Main.svelte";
  import Download from "./pages/Download.svelte";
  import Import from "./pages/Import.svelte";
  import Status from "./pages/Status.svelte";
  import User from "./pages/User.svelte";
  import Search from "./pages/Search.svelte";
  import Post from "./pages/Post.svelte";
  import PostHistory from "./pages/PostHistory.svelte";
  import "bootstrap";
  import "bootstrap/dist/css/bootstrap.min.css";
  import '@fortawesome/fontawesome-free/js/all.js';

  export let url = "";

  function getProps({ location, href, isPartiallyCurrent, isCurrent }) {
    const isActive = href === "/" ? isCurrent : isPartiallyCurrent || isCurrent;

    if (isActive) {
      return { class: "nav-link active" };
    }
    return { class: "nav-link"};
  }
</script>

<style>
  :global(body) {
    background-color: #303030;
    height: 100vh;
  }
</style>


<Router url="{url}">
  <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
    <div class="container-fluid">
      <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarTogglerDemo01" aria-controls="navbarTogglerDemo01" aria-expanded="false" aria-label="Toggle navigation">
        <span class="navbar-toggler-icon"></span>
      </button>
      <div class="collapse navbar-collapse" id="navbarTogglerDemo01">
        <Link class="navbar-brand" to="/"><i class="fa-brands fa-stack-overflow"></i>&nbsp;&nbsp;Dumpflow</Link>
        <ul class="navbar-nav me-auto mb-2 mb-lg-0">
          <li class="nav-item">
            <Link to="/" getProps="{getProps}"><i class="fa-solid fa-house"></i>&nbsp;&nbsp;Home</Link>
          </li>
          <li class="nav-item">
            <Link  to="/import" getProps="{getProps}"><i class="fa-solid fa-file-import"></i>&nbsp;&nbsp;Import</Link>
          </li>
          <li class="nav-item">
            <Link to="/download" getProps="{getProps}"><i class="fa-solid fa-download"></i>&nbsp;&nbsp;Download</Link>
          </li>
          <li class="nav-item">
            <Link to="/status" getProps="{getProps}"><i class="fa-solid fa-signal"></i>&nbsp;&nbsp;Status</Link>
          </li>
          <li class="nav-item">
            <a href="/swagger" class="nav-link" target="_blank"><i class="fa-solid fa-book"></i>&nbsp;&nbsp;API</a>
          </li>
        </ul>
      </div>
    </div>
  </nav>
  <div>
    <Route path="/status" component="{Status}" />
    <Route path="/import" component="{Import}" />
    <Route path="/download" component="{Download}" />
    <Route path="/site/:site/user/:user" let:params>
      <User site={params.site} user={params.user}/>
    </Route>
    <Route path="/site/:site" let:params>
      <Search site={params.site}/>
    </Route>
    <Route path="/site/:site/post/:post" let:params>
      <Post site={params.site} post={params.post}/>
    </Route>
    <Route path="/site/:site/post/:post/history" let:params>
      <PostHistory site={params.site} post={params.post}/>
    </Route>
    <Route path="/" component="{Main}"/>
  </div>
</Router>