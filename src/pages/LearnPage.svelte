<script>
    import { fade } from "svelte/transition";
</script>

<div class="learnPage">
    <span in:fade={{duration: 200}}>
        This page is not implemented yet.
    </span>
</div>

<style>
    .learnPage{
        display: flex;
        flex: 1;
        align-items: center;
        justify-content: center;

        color: var(--text-color-secondary);
        font-family: var(--font-header);
    }
</style>