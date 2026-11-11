<script>
    import {getContext} from "svelte";
    import {fade} from "svelte/transition";

    const leaveLearningPage = getContext('leaveLearningPage');
    const goToLearningPage = getContext('learningPage');
    const showInfo = getContext('showInfo')

    let {mode, len, wrongCnt, rightCnt, name, id, markedOnly = false} = $props();
    let score = $state(0);

    let progess = $state(0);

    if(mode === 0){
        score = Number((100 / (1 + (wrongCnt*2) / len)).toFixed(1));
    }
    else {
        score = Number(((rightCnt/len) * 100).toFixed(1));
    }
    function learnAgain() {
        goToLearningPage(name, id, mode, markedOnly);
    }

    function learnMarked() {
        goToLearningPage(name, id, mode, true);
    }

    async function getStatsUpdateStats(){
        const update = await fetch("/api/collection/stats/update", {
            method: "PATCH",
            credentials: "include",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ collectionID: id, mode: mode, rightCnt: rightCnt, wrongCnt: wrongCnt })
        });
        if(!update.ok){
            showInfo("Failed to update stats.", 1);
            return;
        }
        const data = await update.json();

        score = data.score;
        if(data.previousPlayed > 0){
            progess = score - data.previousLastScore;
        }
    }
    if (!markedOnly) {
        getStatsUpdateStats();
    }

</script>

<div class="progressReport" transition:fade={{duration: 250}}>
    <span class="progressReportHeader">Geschafft!</span>
    <span>Score</span>
    <span class="progressReportScore">{score}%</span>
    {#if progess > 0}
        <span style="color: var(--success-color); font-size: var(--font-size-small)"><i class="fa-solid fa-angles-up"></i> +{progess.toFixed(1)}%</span>
    {:else if progess < 0}
        <span style="color: var(--failure-color); font-size: var(--font-size-small)"><i class="fa-solid fa-angles-down"></i> {progess.toFixed(1)}%</span>
    {/if}
    {#if mode === 0}
        <span>Anzahl Falsch</span>
        <span class="progressReportFalse">{wrongCnt}</span>
    {/if}
    <span>Richtig</span>
    <span class="progressReportTrue">{rightCnt}/{len}</span>

    <div class="progressReportActions">
        <button class="successButton" onclick={learnAgain}>Nochmal lernen</button>
        {#if !markedOnly}
            <button class="accentButton" onclick={learnMarked}>Markierte Karten lernen</button>
        {/if}
        <button class="standardButton" onclick={leaveLearningPage}> Zum Hauptmenü</button>
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

        button {
            padding: 10pt !important;
        }
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
        grid-auto-rows: 1fr;
        gap: 5pt;
        width: 70%;
    }
</style>