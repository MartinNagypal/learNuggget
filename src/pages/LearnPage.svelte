<script>
    import { fade } from "svelte/transition";
    import { getContext } from "svelte";
    import Collection from "../lib/Collection.svelte";
    import { ALL_MARKED } from "../lib/markedCards.js";

    const learningPage = getContext('learningPage');
    const showInfo = getContext('showInfo');

    let collections = $state([]);
    let loaded = $state(false);
    let total = $derived(collections.reduce((sum, collection) => sum + collection.marked, 0));

    async function getMarkedCollections(){
        const res = await fetch("/api/collections/marked");
        const data = await res.json();
        if(!res.ok){
            showInfo(data.message, 1);
            return;
        }
        collections = data;
        loaded = true;
    }
    getMarkedCollections();

    function learnAll(mode){
        learningPage("Markierte Karten", ALL_MARKED, mode, true);
    }
</script>

<main class="learnPage" in:fade={{duration: 200}}>
    {#if loaded}
        {#if total === 0}
            <div class="emptyMarked">
                <span>Noch keine Karten markiert.</span>
                <span class="emptyMarkedHint">Beim Lernen mit <i class="fa-regular fa-star"></i> markieren, dann tauchen sie hier auf.</span>
            </div>
        {:else}
            <div class="collection">
                <div class="name">Alle markierten Karten</div>
                <div class="allMarkedCount">
                    {total} {total === 1 ? "Karte" : "Karten"} aus {collections.length} {collections.length === 1 ? "Kollektion" : "Kollektionen"}
                </div>
                <div class="gamemodeButtons">
                    <button class="standardButton" onclick={() => learnAll(2)}><i class="fa-solid fa-infinity"></i></button>
                    <button class="standardButton" onclick={() => learnAll(1)}><i class="fa-solid fa-bullseye"></i></button>
                    <button class="standardButton" onclick={() => learnAll(0)}><i class="fa-solid fa-dumbbell"></i></button>
                </div>
            </div>

            <div class="learnPageHeader">Nach Kollektion</div>
            <div class="markedCollections">
                {#each collections as collection (collection.id)}
                    <Collection name={collection.name} id={collection.id} learningpage={learningPage} reloadCollections={getMarkedCollections} markedMode={true}/>
                {/each}
            </div>
        {/if}
    {/if}
</main>

<style>
    .learnPage{
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 15pt;
        padding: 25pt 0;
    }

    .allMarkedCount{
        text-align: center;
        color: var(--text-color-secondary);
        font-size: var(--font-size-small);
        padding-bottom: 4pt;
    }

    .learnPageHeader{
        font-family: var(--font-header);
        color: var(--text-color-secondary);
        margin-top: 10pt;
    }

    .markedCollections{
        width: 95%;
        display: flex;
        flex-direction: row;
        flex-wrap: wrap;
        justify-content: center;
        gap: 20pt;
    }

    .emptyMarked{
        flex: 1;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 5pt;
        padding-top: 30%;
        text-align: center;

        font-family: var(--font-header);
        color: var(--text-color-secondary);
    }

    .emptyMarkedHint{
        font-size: var(--font-size-small);
    }
</style>
