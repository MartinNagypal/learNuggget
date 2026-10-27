<script>
    let {current, len, known} = $props();
    import {fade} from "svelte/transition";
</script>

<div class="cardProgress">
    <span>{current}</span>
    <span> <span style="color:var(--success-color);">{known}</span> / {len}</span>
    <progress class="cardProgressBar" value={known} max={len} transition:fade={{duration: 200}}></progress>
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