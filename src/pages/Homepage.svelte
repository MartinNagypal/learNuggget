<script>
    import '../style/Homepage.css';
    import Collection from '../lib/Collection.svelte'
    import NewCollection from "../lib/NewCollection.svelte";

    let name = "Martin"
    let {learningpage} = $props();
    let collections = $state([]);
    let newCollection = $state(0);

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

    async function getCollections(){
        const response = await fetch("/api/collections");
        const data = await response.json();
        collections.length = 0;
        for(const collection of data){
            collections[collections.length] = collection;
        }
    }
    getCollections();
    setInterval(getCollections, 4000);

</script>

<main class="homepage">
    <div class="welcome">
        <span>Hallo <span style="color:var(--success-color)">{name}</span></span>
        <span class="welcomeOrange">Übe weiter</span>
        <span> da wo du aufgehört hast!</span>
    </div>
    <h3 >Kollektionen</h3>
    <div class="collectionsContainer">
        {#each collections as collection}
            <Collection name="{collection}" learningpage={learningpage} reloadCollections={getCollections}/>
        {/each}
        <div class="newCollection" onclick={toggleNewCollection}>
            {#if newCollection === 0}
                Neue Kollektion
            {:else }
                <NewCollection closeNewCollection={closeNewCollection} reloadCollections={getCollections}/>
            {/if}
        </div>
    </div>


</main>