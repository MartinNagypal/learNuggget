<script>
    import {getContext} from "svelte";
    import {fade} from "svelte/transition";

    const goToHomepage = getContext('homepage');
    const goToLearningPage = getContext('learningPage');
    console.log(len);
    console.log(rightCnt);
    console.log(wrongCnt)

    let {mode, len, wrongCnt, rightCnt, name, id} = $props();
    let score = $state(0);

    if(mode === 0){
        score = (100 / (1 + (wrongCnt*2) / len)).toFixed(2);
    }
    else {
        score = ((rightCnt/len) * 100).toFixed(2);
    }
    function learnAgain() {
        goToLearningPage(name, id, mode);
    }

</script>

<div class="progressReport" transition:fade={{duration: 250}}>
    <span class="progressReportHeader">Geschafft!</span>
    <span>Score</span>
    <span class="progressReportScore">{score}%</span>
    {#if mode === 0}
        <span>Anzahl Falsch</span>
        <span class="progressReportFalse">{wrongCnt}</span>
    {:else if mode === 1}
        <span>Richtig</span>
        <span class="progressReportTrue">{rightCnt}/{len}</span>
    {/if}

    <div class="progressReportActions">
        <button class="successButton" onclick={learnAgain}>Nochmal lernen</button>
        <button class="accentButton">Markierte Karten lernen</button>
        <button class="standardButton" onclick={goToHomepage}> Zum Hauptmenü</button>
    </div>
</div>

<style>
    .progressReport{
        font-family: var(--font-header);
        font-size: var(--font-size-middle);
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        gap: 5pt;

        width: 80%;
        max-width: min(75%, 400px);

        border: 1.5pt solid var(--border-color);
        padding: 15pt;
        border-radius: var(--border-radius-primary);
        background-color: var(--accent-surface);
    }

    .progressReportHeader{
        color: var(--accent-primary);
        font-size: var(--font-size-big);
    }

    .progressReportScore {
        color: var(--success-color);
    }

    .progressReportFalse {
        color: var(--failure-color);
    }
    .progressReportTrue {
        color: var(--success-color);
    }

    .progressReportActions {
        display: grid;
        grid-template-rows: repeat(3, 1fr);
        gap: 5pt;
        width: 70%;
    }
</style>