<script>
  import Card from './lib/Card.svelte';
  import Homepage from './pages/Homepage.svelte';
  import LearningPage from "./pages/LearningPage.svelte";
  import Topbar from './lib/Topbar.svelte'
  import Footer from './lib/Footer.svelte'

  let page = $state(0);
  let selectedCollection = $state("")
  let selectedId = $state(0)

  function homePage(){
   page = 0;
  }

  function learningPage(name, id){
      selectedCollection = name;
      selectedId = id;
      page = 1;
  }

</script>

<main class="app">
    <div class="topbar"><Topbar homepage={homePage}/></div>
 {#if page === 0}
     <Homepage learningpage={learningPage} />
  {:else}
      <LearningPage name={selectedCollection} id={selectedId}/>
  {/if}
    <div class="footer"><Footer /></div>
</main>