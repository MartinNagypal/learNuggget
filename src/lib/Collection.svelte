<script>
    import '../style/Collection.css';
    import AddSingleCard from "./AddSingleCard.svelte";
    import CollectionSettings from "./CollectionSettings.svelte";
    import CollectionStats from "./CollectionStats.svelte";
    import {getContext} from 'svelte';
    import { fade, fly, scale } from 'svelte/transition';

    const showInfo = getContext('showInfo');
    let {name, id, learningpage, reloadCollections, markedMode = false} = $props();
    let addSingleCardActive = $state(0);
    let showCollectionSettings = $state(0);

    function toggleShowCollectionSettings(){
        if(showCollectionSettings === 0){
            showCollectionSettings = 1;
        }
        else {
            showCollectionSettings = 0;
        }
    }

    let delCount = $state(0)
    let delTimer;
    async function removeCollection() {
        if(delCount === 0){
            delCount ++;
            showInfo("Click again to approve", 1);
            clearTimeout(delTimer);
            delTimer = setTimeout(() => {
                delCount = 0;
            }, 3000);
            return;
        }
        // reset, so a failed delete needs a new confirmation
        clearTimeout(delTimer);
        delCount = 0;
        if(!name?.trim()){
            return;
        }
        const res = await fetch(`/api/collections/delete?id=${id}`, {
            method: 'DELETE'
        })
        if (res.status === 200) {
            reloadCollections();
        }
    }

    let statsVersion = $state(0);
    function reloadStats(){
        statsVersion++;
    }

</script>

<div class="collection" out:fly={{ x: 80, duration: 300, opacity: 0 }} in:fly={{ x: -80, duration: 200, opacity: 0 }}>
    <div class="name">{name}</div>
    <div class="collectionStatComp">
        <CollectionStats id={id} version={statsVersion}/>
    </div>
    <div class="collectionActions">
        {#if !markedMode}
            <div class="modifyButtons">
                {#if showCollectionSettings === 0}
                    <button class="successButton" onclick={toggleShowCollectionSettings} in:scale={{ duration: 150, start: 0.8 }}><i class="fa-solid fa-angle-down"></i></button>
                {:else}
                    <button class="dangerButton" onclick={toggleShowCollectionSettings} in:scale={{ duration: 150, start: 0.8 }}><i class="fa-solid fa-angle-up"></i></button>
                {/if}
                <button class="dangerButton" onclick={removeCollection}><i class="fa-solid fa-trash"></i></button>
            </div>
        {/if}
        <div class="gamemodeButtons" class:gamemodeButtonsRight={markedMode}>
            <button class="standardButton"onclick={() => learningpage(name, id, 2, markedMode)}><i class="fa-solid fa-infinity"></i></button>
            <button class="standardButton"onclick={() => learningpage(name, id, 1, markedMode)}><i class="fa-solid fa-bullseye"></i></button>
            <button class="standardButton" onclick={() => learningpage(name, id, 0, markedMode)}><i class="fa-solid fa-dumbbell"></i></button>
        </div>

    </div>
    {#if showCollectionSettings === 1 && !markedMode}
        <AddSingleCard CollectionID={id} toggleAddSingleCardActive={toggleShowCollectionSettings} reloadCollectionStats={reloadStats}/>
        <CollectionSettings id={id} toggleShowCollectionSettings={toggleShowCollectionSettings} reloadCollectionLen={reloadStats}/>
    {/if}
</div>