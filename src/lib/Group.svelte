<script>
    import Collection from "./Collection.svelte";
    import {getContext, onDestroy} from "svelte";
    import { fly, fade } from 'svelte/transition';
    import NewCollection from "./NewCollection.svelte";

    const showInfo = getContext('showInfo');
    let {name, id, create, learningPage, reloadGroups, collectionsVersion = 0} = $props();
    let open = $state(0);
    let newGroupName = $state("");

    let collections = $state([]);
    let newCollection = $state(0);

    let delGroup = $state(0);
    let delGroupTimer;

    async function deleteGroup(){
        if(delGroup === 0){
            delGroup = 1;
            showInfo("Click again to approve", 1);
            clearTimeout(delGroupTimer);
            delGroupTimer = setTimeout(() => {
                delGroup = 0;
            }, 3000);
            return;
        }
        // reset, so a failed delete needs a new confirmation
        clearTimeout(delGroupTimer);
        delGroup = 0;
        const res = await fetch (`/api/group/delete?gid=${id}`, {
            method: 'DELETE'
        })
        const data = await res.json();
        if (!res.ok){
            showInfo(data.message, 1);
            return;
        }
        showInfo(data.message, 0);
        await reloadGroups();
    }

    class collectionElement{
        constructor(name, id, gID) {
            this.name = name;
            this.id = id;
            this.gid = gID;

        }
    }

    async function getCollections(){
        const response = await fetch(`/api/collections?gid=${id}`);
        const data = await response.json();
        let newCollections = data.map(collection => ({
            name: collection.Name,
            id: collection.Id,
            gid: collection.Gid
        }));
        if (JSON.stringify(collections) !== JSON.stringify(newCollections)) {
            collections = 0;
            collections = newCollections;
        }
    }
    // Runs once at start and again whenever the homepage counts
    // collectionsVersion up (after a group was created or deleted)
    $effect(() => {
        collectionsVersion;
        getCollections();
    });
    const interval = setInterval(getCollections, 25000);
    onDestroy(() => clearInterval(interval));

    function toggleOpen(){
        if(open === 0){
            open = 1;
        }
        else {
            open = 0;
        }
    }

    function toggleNewCollection(){
        if(newCollection === 0){
            newCollection = 1;
        }
    }

    function closeNewCollection(){
        if(newCollection === 1){
            newCollection = 0;
        }
    }

    async function createGroup(){
        if(newGroupName === "" || newGroupName === " "){
            closeNewCollection();
            showInfo("No name provided", 1);
            return;
        }

        const params = new URLSearchParams({name: newGroupName.trim()});
        const res = await fetch(`/api/group/create?${params}`, {
            method: "POST",
        });
        const data = await res.json();
        if(!res.ok){
            showInfo(data.message, 1);
            return;
        }
        showInfo(data.message, 0);
        newGroupName = "";
        open = 0;
        await reloadGroups();
    }

</script>
<div class="group">
    <div class="groupHeader">
        <div onclick={toggleOpen}>
            {#if open === 0}
                <i class="fa-solid fa-angle-right"></i>
            {:else}
                <i class="fa-solid fa-angle-down"></i>
            {/if}
            <span>{name}</span>
        </div>
        {#if id != 0}
            <div class="delGroup">
                <i class="fa-solid fa-trash" onclick={deleteGroup}></i>
            </div>
        {/if}
    </div>
    {#if open === 1}
        <div out:fly={{ y: -40, duration: 300, opacity: 0 }} in:fly={{ y: -40, duration: 200, opacity: 0 }}>
            {#if create === 1}
                <div class="createGroup">
                    <div class="createGroupHeader">Name</div>
                    <input type="text" class="mainInputStyle" placeholder="Name" bind:value={newGroupName}>
                    <button class="standardButton" onclick={createGroup}>Hinzufügen</button>
                </div>
            {:else}
                <div class="groupCollections">
                    {#each collections as collection (collection.id)}
                        {#if collection.gid === id}
                            <Collection name={collection.name} id={collection.id} learningpage={learningPage} reloadCollections={getCollections}/>
                        {/if}
                    {/each}
                    <div class="newCollection" onclick={toggleNewCollection}>
                        {#if newCollection === 0}
                            Neue Kollektion
                        {:else }
                            <div  in:fade={{duration: 200}}><NewCollection closeNewCollection={closeNewCollection} reloadCollections={getCollections} gid={id}/></div>
                        {/if}
                    </div>
                </div>
            {/if}
        </div>
    {/if}
</div>

<style>
    .group {
        display: flex;
        flex-direction: column;
        gap: 10pt;
        min-width: 250px;
        max-width: 375px;
        width: 90%;
    }

    .groupHeader{
        font-family: var(--font-header);
        color: var(--text-color-secondary);
        user-select: none;

        border-bottom: 1.5pt solid var(--border-color);
        border-radius: var(--border-radius-primary);
        padding: 2.5pt 7.5pt;

        display: flex;
        flex-direction: row;
        justify-content: space-between;
    }

    .groupCollections{
        display: flex;
        flex-direction: column;
        gap: 10pt;
        width: 100%;
    }

    .createGroup{
        display: flex;
        flex-direction: column;
        gap: 5pt;
        align-items: center;
        justify-content: center;
    }

    .createGroupHeader {
        font-family: var(--font-header);
        color: var(--accent-secondary);
    }

    .groupCollections {
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .delGroup {
        transition: ease 0.2s;
    }

    .delGroup:active {
        transform: scale(0.8);
    }

    .newCollection {
        background-color: var(--accent-surface);
        width: 90%;
        min-width: 250px;
        max-width: 375px;
        border-style: solid;
        border-width: 1.25pt;
        border-radius: var(--border-radius-primary);
        border-color: var(--border-color);
        padding: 10pt;

        box-shadow: 0 0 2.5pt 1.5pt var(--accent-surface);
        color: var(--success-color);
        font-family: var(--font-header);
        user-select: none;

        display:flex;
        align-items: center;
        justify-content: center;
        transition: all 0.2s;
    }

    .newCollection:active {
        transform: scale(0.95);
    }
</style>