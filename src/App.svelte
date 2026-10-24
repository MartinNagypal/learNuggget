<script>
  import Card from './lib/Card.svelte';
  import Homepage from './pages/Homepage.svelte';
  import LearningPage from "./pages/LearningPage.svelte";

  let page = $state(0);

  function homePage(){
   page = 0;
  }

  function learningPage(){
   page = 1;
  }

</script>

<main>
 {#if page === 0}
     <Homepage homepage={homePage} learningpage={learningPage} />
  {:else}
      <LearningPage homepage={homePage}/>
  {/if}
</main>