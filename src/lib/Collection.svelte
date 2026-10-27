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

    function toggleAddSingleCardActive() {
        if(addSingleCardActive === 0){
            addSingleCardActive = 1;
        }
        else {
            addSingleCardActive = 0;
        }
    }

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
                <button class="successButton" onclick={toggleAddSingleCardActive} in:scale={{ duration: 150, start: 0.8 }}>+</button>
            {:else}
                <button class="dangerButton" onclick={toggleAddSingleCardActive} in:scale={{ duration: 150, start: 0.8 }}>x</button>
            {/if}
            <button class="standardButton" onclick={toggleShowCollectionSettings}><i class="fa-solid fa-pen"></i></button>
            <button class="dangerButton" onclick={removeCollection}><i class="fa-solid fa-trash"></i></button>
        </div>
        <button class="standardButton" onclick={() => learningpage(name, id)}>Lernen</button>
    </div>
    {#if showCollectionSettings === 1}
        <CollectionSettings id={id} toggleShowCollectionSettings={toggleShowCollectionSettings}/>
    {/if}
    {#if addSingleCardActive === 1}
        <AddSingleCard CollectionID={id} toggleAddSingleCardActive={toggleAddSingleCardActive} reloadCollectionLen={getCollectionLen}/>
    {/if}
</div>