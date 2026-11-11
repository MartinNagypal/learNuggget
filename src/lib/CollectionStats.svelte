<script>
    let {id, version} = $props();

    let len = $state(0);
    let marked = $state(0);
    let avgScore = $state(0);
    let lastScore = $state(0);
    let totalPlayed = $state(0);

    async function getCollectionStats(){
        const res = await fetch(`/api/collection/stats?id=${encodeURIComponent(id)}`);
        if(!res.ok){
            return;
        }
        const data = await res.json();
        avgScore = data.avgScore;
        lastScore = data.lastScore;
        totalPlayed = data.totalPlayed;
        len = data.len;
        marked = data.marked;
    }
    getCollectionStats();

    $effect(() => {
        version;
        getCollectionStats();
    });

</script>
<div class="collectionStats">
    <div class="collectionStatElement">
        <div class="statHeader">{len}</div>
        <div class="statFooter">Cards</div>
    </div>
    <div class="collectionStatElement">
        <div class="statHeader">{marked}</div>
        <div class="statFooter">Marked</div>
    </div>
    <div class="collectionStatElement">
        <div class="statHeader">{avgScore}%</div>
        <div class="statFooter">Average </div>
    </div>
    <div class="collectionStatElement">
        <div class="statHeader">{lastScore}%</div>
        <div class="statFooter">Last</div>
    </div>
</div>

<style>
    .collectionStats{
        display: grid;
        grid-template-columns: repeat(4, 1fr);
        gap: 5pt;
    }

    .collectionStatElement{
        padding: 3pt 0;
        display: flex;
        flex-direction: column;
        align-items: center;
        align-content: center;
    }

    .statHeader{
        font-family: var(--font-header);
        color: var(--accent-primary);
    }
    .statFooter{
        font-size: var(--font-size-small);
        color: var(--text-color-secondary);
    }
</style>