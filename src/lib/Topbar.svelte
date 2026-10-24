<script>
    import '../style/Topbar.css'

    let {homepage} = $props();

</script>

<div class="topbar">
    <span class="brandName">learNuggget</span>
    <span class="brandName homebutton" onclick={homepage}><i class="fa-solid fa-house"></i></span>
</div>