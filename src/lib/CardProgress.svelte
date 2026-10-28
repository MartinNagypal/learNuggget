<script>
    let {mode, current, len, progress} = $props();
    import {fade} from "svelte/transition";
</script>

<div class="cardProgress">
    <span>{current}</span>
    {#if mode != 2}
        <span> <span style="color:var(--success-color);">{progress}</span> / {len}</span>
        <progress class="cardProgressBar" value={progress} max={len} transition:fade={{duration: 200}}></progress>
    {/if}
</div>

<style>

    .cardProgress{
        width: 90%;
        display: grid;
        grid-template-columns: 1fr 3fr 4fr;
        align-items: center;
        gap: 10pt;
    }

    .cardProgressBar {
        -webkit-appearance: none;
        appearance: none;

        border: none;
        height: 10px;
        width: 80%;
        accent-color: var(--success-color);
        border-radius: var(--border-radius-primary);
        background-color: var(--text-color-white);
        transition: all 0.3s;
    }

    .cardProgressBar::-webkit-progress-bar {
        background-color: var(--text-color-white);
        border-radius: var(--border-radius-primary);
    }

    .cardProgressBar::-webkit-progress-value {
        background-color: var(--success-color);
        border-radius: var(--border-radius-primary);
    }

    .cardProgressBar::-moz-progress-bar {
        background-color: var(--success-color);
        border-radius: var(--border-radius-primary);
    }
</style>