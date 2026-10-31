<script>
    import katex from 'katex';
    import 'katex/dist/katex.min.css';
    let { text } = $props();

    const parts = $derived(
        (text ?? '')
            .split(/\$\$([\s\S]+?)\$\$/g)
            .map((value, i) => ({ isLatex: i % 2 === 1, value }))
            .filter(p => p.value.trim() !== '')
    );
</script>

<div class="rich">
    {#each parts as part}
        {#if part.isLatex}
            <div class="formula">
                {@html katex.renderToString(part.value, { displayMode: false, throwOnError: false })}
            </div>
        {:else}
            <div>{part.value.trim()}</div>
        {/if}
    {/each}
</div>

<style>
    .rich {
        display: flex;
        flex-direction: row;
        flex-wrap: wrap;
        align-items: center;
        gap: 5pt;
    }

    .formula {
        color: var(--formula-color);
        font-weight: bold;
        margin-left: 5pt;
        margin-right: 5pt;
    }
</style>