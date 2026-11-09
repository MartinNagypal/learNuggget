<script>
    import ProfileInfoSocial from "../lib/ProfileInfoSocial.svelte";
    import GlobbalRanklist from "../lib/GlobbalRanklist.svelte";
    import {fade} from "svelte/transition";

    let {name} = $props();
</script>

<div class="socialPage" in:fade={{duration: 200}}>
    <ProfileInfoSocial name={name}/>
    <GlobbalRanklist/>
</div>

<style>
    .socialPage{
        display: flex;
        flex: 1;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        gap: 30pt;
    }
</style>