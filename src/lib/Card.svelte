<script>
    import '../style/Card.css';
    import { scale } from 'svelte/transition';
    import { fly } from 'svelte/transition';
    let {collection, text, answer, site, changeSite, direction} = $props();


</script>

<div class="card" in:fly={{ x: direction * 80, duration: 300, opacity: 0 }} out:fly={{ x: direction * -80, duration: 200, opacity: 0 }}>
    {#if site === 0}
        <div class="content front">
            <div class="cardTop">{collection}</div>
            <div onclick={changeSite} class="cardContent" in:scale={{ duration: 320, start: 0.75, opacity: 0 }} out:scale={{ duration: 180, start: 1.15, opacity: 0 }}>{text}</div>
            <div onclick={changeSite} class="answer hidden"> {answer}</div>
            <div onclick={changeSite} class="cardBottom" in:scale={{ duration: 320, start: 0.75, opacity: 0 }} out:scale={{ duration: 180, start: 1.15, opacity: 0 }}>
                Tippe zum Umdrehen
            </div>
        </div>
    {:else}
        <div class="content">
            <div class="cardTop">{collection}</div>
            <div onclick={changeSite} class="question" in:scale={{ duration: 320, start: 0.75, opacity: 0 }} out:scale={{ duration: 180, start: 1.15, opacity: 0 }}>{text}</div>
            <div onclick={changeSite} class="answer" in:scale={{ duration: 320, start: 0.75, opacity: 0 }} out:scale={{ duration: 180, start: 1.15, opacity: 0 }} > {answer}</div>
            <div class="cardActions" in:scale={{ duration: 320, start: 0.75, opacity: 0 }} out:scale={{ duration: 180, start: 1.15, opacity: 0 }}>
                <button class="failedButton button">Falsch</button>
                <button class="successButton button">Richtig</button>
            </div>
        </div>
    {/if}

</div>