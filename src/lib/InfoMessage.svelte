<script>
    import {fly} from "svelte/transition";
    let {infoMesage, error} = $props();
</script>
<div class="infoMessage" in:fly={{ y: -50, duration: 300}} out:fly={{ y: -50, duration: 300}}>
    {#if error === 0}
        <div class="infoMessageInfo">
            <i class="fa-solid fa-circle-check"></i>
            <span>{infoMesage}</span>
        </div>
    {:else}
        <div class="infoMessageError">
            <i class="fa-solid fa-circle-xmark"></i>
            <span>{infoMesage}</span>
        </div>
    {/if}
</div>

<style>
    .infoMessage {
        font-family: var(--font-header);
    }

    .infoMessageInfo, .infoMessageError {
        width: 100%;
        height: 100%;
        display: flex;
        flex-direction: row;
        flex-wrap: wrap;
        justify-content: center;
        align-items: center;
        gap: 5pt;

        padding: 10pt;
        padding-right: 20pt;
        padding-left: 20pt;
        box-sizing: border-box;
        border-radius: var(--border-radius-primary);
        background-color: var(--accent-surface);

        i {
            text-align: center;
        }
    }

    .infoMessageInfo {
        color: var(--success-color);
        border: 2pt solid var(--success-color);
        box-shadow: 0 0 150px 2px var(--success-color);
    }

    .infoMessageError {
        color: var(--failure-color);
        border: 2pt solid var(--failure-color);
        box-shadow: 0 0 100px 1px var(--failure-color);
    }
</style>