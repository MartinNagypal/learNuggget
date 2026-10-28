<script>
    import '../style/Topbar.css'

    let {homepage, isLoggedIn} = $props();

</script>

<div class="topbar">
    <img src="../../public/logo.png" height="40px" alt="logo">

    {#if isLoggedIn === 1}
        <span class="brandName homebutton" onclick={homepage}><i class="fa-solid fa-house"></i></span>
    {/if}
</div>