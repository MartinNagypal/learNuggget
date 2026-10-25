<script>
    import '../style/Topbar.css'

    let {homepage} = $props();

</script>

<div class="topbar">
    <img src="../../public/logo.png" height="40px" alt="logo">
    <span class="brandName homebutton" onclick={homepage}><i class="fa-solid fa-house"></i></span>
</div>