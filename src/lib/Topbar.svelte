<script>
    import '../style/Topbar.css'
    import {getContext} from "svelte";

    const verifySession = getContext("verifySession")
    const showInfo = getContext('showInfo')
    let {homepage, isLoggedIn} = $props();

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

</script>

<div class="topbar">
    <img src="../../public/logo.png" height="40px" alt="logo">

    {#if isLoggedIn === 1}
        <div class="topbarActions">
            <i class="fa-solid fa-right-from-bracket fa-xl homebutton" onclick={logout}></i>
            <span class="brandName homebutton" onclick={homepage}><i class="fa-solid fa-house"></i></span>
        </div>

    {/if}
</div>