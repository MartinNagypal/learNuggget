<script>
    import '../style/Homepage.css';
    import Collection from '../lib/Collection.svelte'
    import NewCollection from "../lib/NewCollection.svelte";
    import { fade } from 'svelte/transition';
    import { scale } from 'svelte/transition';

    let name = "Martin"
    let {learningpage} = $props();
    let collections = $state([]);
    let newCollection = $state(0);
    let sortBy = $state(0);

    class collectionElement{
        constructor(name, id) {
            this.name = name;
            this.id = id;
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

    function toggleSort(){
        sortBy += 1;
        if(sortBy === 4) {
            sortBy = 0;
        }
        collections = [];
        getCollections();
    }

    async function getCollections(){
        const params = new URLSearchParams({orderState: sortBy.toString()});
        const response = await fetch(`/api/collections?${params}`);
        const data = await response.json();
        let newCollections = data.map(collection => ({
            name: collection.Name,
            id: collection.Id
        }));
        if (JSON.stringify(collections) !== JSON.stringify(newCollections)) {
            collections = 0;
            collections = newCollections;
        }
    }
    getCollections();
    setInterval(getCollections, 25000);

</script>

<main class="homepage">
    <div class="welcome">
        <span>Hallo <span style="color:var(--success-color)">{name}</span></span>
        <span class="welcomeOrange">Übe weiter</span>
        <span> da wo du aufgehört hast!</span>
    </div>
    <div class="collectionsHeader">
        <h3 >Kollektionen</h3>
        {#if sortBy === 0}
            <h5 onclick={toggleSort}> <i class="fa-solid fa-arrow-down-short-wide fa-lg"></i> Neuste</h5>
        {:else if sortBy === 1}
            <h5 onclick={toggleSort}> <i class="fa-solid fa-arrow-up-short-wide fa-lg"></i> Älteste</h5>
        {:else if sortBy === 2}
            <h5 onclick={toggleSort}><i class="fa-solid fa-arrow-down-a-z fa-lg"></i> A-Z</h5>
        {:else if sortBy === 3}
            <h5 onclick={toggleSort}><i class="fa-solid fa-arrow-up-a-z fa-lg"></i> Z-A</h5>
        {/if}
    </div>
    <div class="collectionsContainer">
        {#each collections as collection (collection.id)}
            <Collection name="{collection.name}" id={collection.id} learningpage={learningpage} reloadCollections={getCollections}/>
        {/each}
        <div class="newCollection" onclick={toggleNewCollection}>
            {#if newCollection === 0}
                Neue Kollektion
            {:else }
                <div  in:fade={{duration: 200}}><NewCollection closeNewCollection={closeNewCollection} reloadCollections={getCollections}/></div>
            {/if}
        </div>
    </div>
</main>