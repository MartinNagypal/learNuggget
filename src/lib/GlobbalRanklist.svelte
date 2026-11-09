<script>
    let num = $state(0);
    let topUsers = $state([]);

    const rankClasses = ["gold", "silver", "bronze"];

    async function getTopUsers() {
        const res = await fetch("/api/user/top");
        if (!res.ok) return;
        const data = await res.json();
        topUsers = data.topUsers ?? [];
    }

    getTopUsers();

</script>

<div class="globalRanklist">
    <div class="globalRanklistHeader">
        Global Top 5
    </div>

    {#each topUsers as user, i}
        <div class="rankListElement {rankClasses[i] ?? ''}">
            <span>#{i + 1}</span>
            <span>{user.username}</span>
            <span>{user.avgScore.toFixed(1)}%</span>
        </div>
    {:else}
        <div class="rankListElement">Noch keine Einträge</div>
    {/each}
</div>

<style>
    .globalRanklist {
        display: flex;
        flex-direction: column;
        gap: 10pt;
        padding: 20pt;
    }

    .globalRanklistHeader{
        font-family: var(--font-header);
        font-size: var(--font-size-middle);
        color: var(--accent-primary);
        border-bottom: 1.5pt solid var(--border-color);
        padding-bottom: 2pt;
    }

    .rankListElement{
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: center;
        gap: 50pt;

        border: 1.5pt solid var(--border-color);
        border-radius: var(--border-radius-primary);
        padding: 5pt;

        font-family: var(--font-header);
        text-align: center;

    }

    .gold {
        background-color: var(--rank-gold-surface);
        color: var(--rank-gold-text);
    }
    .silver {
        background-color: var(--rank-silver-surface);
        color: var(--rank-silver-text);
    }
    .bronze{
        background-color: var(--rank-bronze-surface);
        color: var(--rank-bronze-text);
    }

</style>