<script>
    import '../style/Collection.css';
    import AddSingleCard from "./AddSingleCard.svelte";
    import { fade, fly } from 'svelte/transition';
    let {name, id, learningpage, reloadCollections} = $props();
    let addSingleCardActive = $state(0);

    function toggleAddSingleCardActive() {
        if(addSingleCardActive == 0){
            addSingleCardActive = 1;
        }
        else {
            addSingleCardActive = 0;
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

</script>

<div class="collection">
    <div class="name">{name}</div>
    <div class="collectionActions">
        <div class="modifyButtons">
            <button class="successButton" onclick={toggleAddSingleCardActive}>+</button>
            <button class="standardButton" onclick={printName}><i class="fa-solid fa-pen"></i></button>
            <button class="dangerButton" onclick={removeCollection}><i class="fa-solid fa-trash"></i></button>
        </div>
        <button class="standardButton" onclick={() => learningpage(name)}>Lernen</button>
    </div>
    {#if addSingleCardActive == 1}
        <AddSingleCard CollectionID={id} toggleAddSingleCardActive={addSingleCardActive}/>
    {/if}
</div>