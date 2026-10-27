<script>
    import '../style/Collection.css';
    import AddSingleCard from "./AddSingleCard.svelte";
    import CollectionSettings from "./CollectionSettings.svelte";
    import { onMount } from 'svelte';
    import { fade, fly, scale } from 'svelte/transition';
    let {name, id, learningpage, reloadCollections} = $props();
    let addSingleCardActive = $state(0);
    let showCollectionSettings = $state(0);
    let len = $state(0);

    function toggleShowCollectionSettings(){
        if(showCollectionSettings === 0){
            showCollectionSettings = 1;
        }
        else {
            showCollectionSettings = 0;
        }
    }

    function printName(){
        console.log("this is " + name);
    }

    async function removeCollection() {
        if(!name?.trim()){
            return;
        }
        const params = new URLSearchParams({name: name.trim()});
        const res = await fetch(`/api/collections/delete?${params}`, {
            method: 'DELETE'
        })
        if (res.status === 200) {
            reloadCollections();
        }
    }

    async function getCollectionLen(){
        const res = await fetch(`/api/collection/len?id=${id}`)
        const data = await res.json()
        console.log(data.len);
        len = data.len;
    }
    onMount(() => {
        getCollectionLen();
    });

</script>

<div class="collection" out:fly={{ x: 80, duration: 300, opacity: 0 }} in:fly={{ x: -80, duration: 200, opacity: 0 }}>
    <div class="name">{name}</div>
    <div class="collectionStats">
        <span><b>{len}</b> Karten</span>
    </div>
    <div class="collectionActions">
        <div class="modifyButtons">
            {#if addSingleCardActive === 0}
                <button class="successButton" onclick={toggleShowCollectionSettings} in:scale={{ duration: 150, start: 0.8 }}>+</button>
            {:else}
                <button class="dangerButton" onclick={toggleShowCollectionSettings} in:scale={{ duration: 150, start: 0.8 }}>x</button>
            {/if}
            <button class="dangerButton" onclick={removeCollection}><i class="fa-solid fa-trash"></i></button>
        </div>
        <div class="gamemodeButtons">
            <button class="standardButton"><i class="fa-solid fa-infinity"></i></button>
            <button class="standardButton"><i class="fa-solid fa-bullseye"></i></button>
            <button class="standardButton" onclick={() => learningpage(name, id)}><i class="fa-solid fa-dumbbell"></i></button>
        </div>

    </div>
    {#if showCollectionSettings === 1}
        <AddSingleCard CollectionID={id} toggleAddSingleCardActive={showCollectionSettings} reloadCollectionLen={getCollectionLen}/>
        <CollectionSettings id={id} toggleShowCollectionSettings={toggleShowCollectionSettings} reloadCollectionLen={getCollectionLen}/>
    {/if}
</div>