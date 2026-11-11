<script>
  import {getContext, setContext} from "svelte";
  import {fade} from "svelte/transition";
  import Login from "./pages/Login.svelte";
  import Homepage from './pages/Homepage.svelte';
  import LearningPage from "./pages/LearningPage.svelte";
  import Topbar from './lib/Topbar.svelte'
  import InfoMessage from "./lib/InfoMessage.svelte";
  import NavBar from "./lib/NavBar.svelte";
  import SocialPage from "./pages/SocialPage.svelte";
  import LearnPage from "./pages/LearnPage.svelte";
  import ConfirmDeleteAccount from "./lib/ConfirmDeleteAccount.svelte";

  let isLoggedIn = $state(0);
  let page = $state(0);
  let selectedCollection = $state("")
  let selectedId = $state(0);
  let learningSession = $state(0);
  let gameMode = $state(0);
  let markedOnly = $state(false);
  let returnPage = $state(0);
  let infoVisible = $state(0);
  let infoMessage = $state("");
  let infoError = $state(0);
  let name = $state("");
  let showConfirmDeleteAccount = $state(0);

  async function verifySession(){
      const res = await fetch("/api/auth/verifySession");
      if(res.ok){
          isLoggedIn = 1;
          getUsername();
          return
      }
      isLoggedIn = 0;
      page = 0;
  }
  verifySession();
  setInterval(verifySession, 120000);
  setContext('verifySession', verifySession);
  function homePage(){
   page = 0;
  }
  setContext('homepage', homePage);

  function toggleIsLoggedIn(){
      isLoggedIn = (isLoggedIn+1) % 2;
  }

  setContext('toggleIsLoggedIn', toggleIsLoggedIn)

  function learningPage(name, id, gamemode, onlyMarked = false){
      if (page !== 1) {
          returnPage = page;
      }
      selectedCollection = name;
      selectedId = id;
      markedOnly = onlyMarked;
      homePage();
      learningSession++;
      page = 1;
      gameMode = gamemode;
  }
  setContext('learningPage', learningPage)

  function leaveLearningPage(){
      page = returnPage;
  }
  setContext('leaveLearningPage', leaveLearningPage);

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

  function setPage(pageNr){
      page = pageNr;
  }
  setContext('setPage', setPage);

  async function getUsername(){
      const res = await fetch("/api/auth/username");
      const data = await res.json();
      name = data.username;
  }
  getUsername();

  function toggleConfirmAccountDeletion(){
      if (showConfirmDeleteAccount === 0){
          showConfirmDeleteAccount = 1;
      }
      else {
          showConfirmDeleteAccount = 0;
      }
  }

  setContext('toggleConfirmAccountDeletion', toggleConfirmAccountDeletion);

  async function logout(){
      const res = await fetch("/api/auth/logout")
      const data = await res.json()
      if(!res.ok){
          showInfo(data.message, 1);
          return
      }
      await verifySession();
      showInfo(data.message, 0);
  }

  setContext('logout', logout);

</script>

<main class="app">
    {#if showConfirmDeleteAccount === 1}
        <div class="confirmDeleteContainer" transition:fade={{duration: 200}}><ConfirmDeleteAccount/></div>
    {/if}
    <div class="topbar"><Topbar homepage={homePage} isLoggedIn={isLoggedIn}/></div>
    {#if infoVisible === 1}
        <div class="infoScreen"><InfoMessage infoMesage={infoMessage} error={infoError}/></div>
    {/if}

    <div class="pageContent">
        {#if isLoggedIn === 1}
            {#if page === 0}
                <Homepage learningpage={learningPage} />
            {:else if page === 1}
                {#key learningSession}
                    <LearningPage name={selectedCollection} id={selectedId} gameMode={gameMode} markedOnly={markedOnly}/>
                {/key}
            {:else if page === 2}
                <SocialPage name={name}/>
            {:else if page === 3}
                <LearnPage/>
            {/if}
        {:else}
            <Login/>
        {/if}
    </div>
    {#if isLoggedIn === 1}
        <div class="mainNavBar"><NavBar/></div>
    {/if}
</main>