<script>
    import '../style/Homepage.css';
    import Collection from '../lib/Collection.svelte'
    import Topbar from '../lib/Topbar.svelte'

    let name = "Martin"
    let collections = ["Algorithmen", "Mathe", "Binärsysteme"]
</script>

<main class="homepage">
    <div class="topbar"><Topbar/></div>
    <div class="welcome">
        <span>Hallo <span style="color:var(--success-color)">{name}</span></span>
        <span class="welcomeOrange">Übe weiter</span>
        <span> da wo du aufgehört hast!</span>
    </div>
    <h3>Kollektionen</h3>
    {#each collections as collection}
        <Collection name="{collection}"/>
    {/each}
</main>