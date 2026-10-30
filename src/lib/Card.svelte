<script>
    import '../style/Card.css';
    import CardProgress from "./CardProgress.svelte";
    import { scale } from 'svelte/transition';
    import { fly } from 'svelte/transition';
    let {mode, text, answer, site, changeSite, direction, current, len, progress, toggleKnownCard, toggleWrong} = $props();

</script>

<div class="card" in:fly={{ x: direction * 80, duration: 300, opacity: 0 }} out:fly={{ x: direction * -80, duration: 200, opacity: 0 }}>
    {#if site === 0}
        <div class="content front">
            <div class="cardTop"><CardProgress mode={mode} current={current} len={len} progress={progress}/></div>
            <div onclick={changeSite} class="cardContent" in:scale={{ duration: 320, start: 0.75, opacity: 0 }} out:scale={{ duration: 180, start: 1.15, opacity: 0 }}>{text}</div>
            <div onclick={changeSite} class="answer hidden"> {answer}</div>
            <div onclick={changeSite} class="cardBottom" in:scale={{ duration: 320, start: 0.75, opacity: 0 }} out:scale={{ duration: 180, start: 1.15, opacity: 0 }}>
                Tippe zum Umdrehen
            </div>
        </div>
    {:else}
        <div class="content">
            <div class="cardTop"><CardProgress mode={mode} current={current} len={len} progress={progress}/></div>
            <div onclick={changeSite} class="question" in:scale={{ duration: 320, start: 0.75, opacity: 0 }} out:scale={{ duration: 180, start: 1.15, opacity: 0 }}>{text}</div>
            <div onclick={changeSite} class="answer" in:scale={{ duration: 320, start: 0.75, opacity: 0 }} out:scale={{ duration: 180, start: 1.15, opacity: 0 }} > {answer}</div>
            <div style="width: 100%;" in:scale={{ duration: 320, start: 0.75, opacity: 0 }} out:scale={{ duration: 180, start: 1.15, opacity: 0 }}>
                {#if mode !== 2}
                    <div class="cardActions">
                        <button class="failedButton button" onclick={toggleWrong}>Falsch</button>
                        <button class="successButton button" onclick={toggleKnownCard}>Richtig</button>
                    </div>
                {/if}
            </div>

        </div>
    {/if}

</div>