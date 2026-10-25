<script>
    import '../style/Collection.css';
    import { fade, fly } from 'svelte/transition';
    let {name, learningpage, reloadCollections} = $props();

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
            <button class="standardButton" onclick={printName}><i class="fa-solid fa-pen"></i></button>
            <button class="dangerButton" onclick={removeCollection}><i class="fa-solid fa-trash"></i></button>
        </div>
        <button class="standardButton" onclick={() => learningpage(name)}>Lernen</button>
    </div>
</div>