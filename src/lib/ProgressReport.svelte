<script>
    import {getContext} from "svelte";
    import {fade} from "svelte/transition";

    const goToHomepage = getContext('homepage');
    const goToLearningPage = getContext('learningPage');
    const showInfo = getContext('showInfo')

    let {mode, len, wrongCnt, rightCnt, name, id} = $props();
    let score = $state(0);

    let avgScore = 0;
    let lastScore = 0;

    if(mode === 0){
        score = Number((100 / (1 + (wrongCnt*2) / len)).toFixed(1));
    }
    else {
        score = Number(((rightCnt/len) * 100).toFixed(1));
    }
    function learnAgain() {
        goToLearningPage(name, id, mode);
    }

    async function getStatsUpdateStats(){
        const res = await fetch(`/api/collection/stats?id=${encodeURIComponent(id)}`);
        if(!res.ok){
            showInfo("Failed to update stats.", 1);
            return;
        }
        const data = await res.json()
        const newAvgScore = Number(((data.avgScore * data.totalPlayed + score) / (data.totalPlayed + 1)).toFixed(1));
        const update = await fetch("/api/collection/stats/update", {
            method: "PATCH",
            credentials: "include",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ collectionID: id, avgScore: newAvgScore, lastScore:score})
        });
        if(!update.ok){
            showInfo("Failed to update stats.", 1);
            return;
        }
    }
    getStatsUpdateStats();

</script>

<div class="progressReport" transition:fade={{duration: 250}}>
    <span class="progressReportHeader">Geschafft!</span>
    <span>Score</span>
    <span class="progressReportScore">{score}%</span>
    {#if mode === 0}
        <span>Anzahl Falsch</span>
        <span class="progressReportFalse">{wrongCnt}</span>
    {/if}
    <span>Richtig</span>
    <span class="progressReportTrue">{rightCnt}/{len}</span>

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
        grid-template-rows: repeat(3, 1fr);
        gap: 5pt;
        width: 70%;
    }
</style>