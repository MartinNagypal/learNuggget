<script>
    import {getContext, setContext} from "svelte";
  import Homepage from './pages/Homepage.svelte';
  import LearningPage from "./pages/LearningPage.svelte";
  import Topbar from './lib/Topbar.svelte'
  import Footer from './lib/Footer.svelte'

  let page = $state(0);
  let selectedCollection = $state("")
  let selectedId = $state(0);
  let learningSession = $state(0);
  let gameMode = $state(0);

  function homePage(){
   page = 0;
  }

  setContext('homepage', homePage);

  function learningPage(name, id, gamemode){
      selectedCollection = name;
      selectedId = id;
      homePage();
      learningSession++;
      page = 1;
      gameMode = gamemode;
  }

  setContext('learningPage', learningPage)

</script>

<main class="app">
    <div class="topbar"><Topbar homepage={homePage}/></div>
 {#if page === 0}
     <Homepage learningpage={learningPage} />
  {:else}
     {#key learningSession}
         <LearningPage name={selectedCollection} id={selectedId} gameMode={gameMode}/>
     {/key}
  {/if}
    <div class="footer"><Footer /></div>
</main>