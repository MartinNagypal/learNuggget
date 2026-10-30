<script>
  import {getContext, setContext} from "svelte";
  import Login from "./pages/Login.svelte";
  import Homepage from './pages/Homepage.svelte';
  import LearningPage from "./pages/LearningPage.svelte";
  import Topbar from './lib/Topbar.svelte'
  import Footer from './lib/Footer.svelte'
  import InfoMessage from "./lib/InfoMessage.svelte";

  let isLoggedIn = $state(0);
  let page = $state(0);
  let selectedCollection = $state("")
  let selectedId = $state(0);
  let learningSession = $state(0);
  let gameMode = $state(0);
  let infoVisible = $state(0);
  let infoMessage = $state("");
  let infoError = $state(0);

  async function verifySession(){
      const res = await fetch("/api/auth/verifySession");
      if(res.ok){
          isLoggedIn = 1;
          return
      }
      isLoggedIn = 0;
  }
  verifySession();
  setInterval(verifySession, 60000);
  setContext('verifySession', verifySession);
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

  function showInfo(info, error){
      infoMessage = info;
      infoError = error;
      infoVisible = 1;
      setTimeout(() => {
          infoVisible = 0;
          infoMessage = "";
          infoError = 0;
      }, 2500);
  }

  setContext('showInfo', showInfo);

</script>

<main class="app">
    <div class="topbar"><Topbar homepage={homePage} isLoggedIn={isLoggedIn}/></div>
    {#if infoVisible === 1}
        <div class="infoScreen"><InfoMessage infoMesage={infoMessage} error={infoError}/></div>
    {/if}
    {#if isLoggedIn === 1}
        {#if page === 0}
            <Homepage learningpage={learningPage} />
        {:else}
             {#key learningSession}
                 <LearningPage name={selectedCollection} id={selectedId} gameMode={gameMode}/>
             {/key}
        {/if}
    {:else}
        <Login/>
    {/if}
    <div class="footer"><Footer /></div>
</main>