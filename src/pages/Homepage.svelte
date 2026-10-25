<script>
    import '../style/Homepage.css';
    import Collection from '../lib/Collection.svelte'

    let name = "Martin"
    let {learningpage} = $props();
    let collections = $state([]);

    async function getCollections(){
        const response = await fetch("/api/collections");
        const data = await response.json();
        for(const collection of data){
            if(!collections.includes(collection)){
                collections[collections.length] = collection;
            }
        }
    }
    getCollections();
    setInterval(getCollections, 5000);

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
            <Collection name="{collection}" learningpage={learningpage}/>
        {/each}
        <div class="newCollection">Neue Kollektion</div>
    </div>


</main>