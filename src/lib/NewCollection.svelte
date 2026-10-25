<script>
    let name = $state("");
    let {closeNewCollection, reloadCollections} = $props();

    async function postNewCollection() {
        if(name === "" || name === " "){
            closeNewCollection();
            return;
        }

        const params = new URLSearchParams({name: name.trim()});
        const res = await fetch(`/api/collections/new?${params}`, {
            method: 'POST'
        })
        if (res.status === 200) {
            closeNewCollection();
            reloadCollections();
            name = "";
            document.getElementById("newCollectionContent").classList.remove("errorBorder");
        }
        else {
            closeNewCollection();
        }
    }

</script>


<main id="newCollectionContent" class="newCollectionContent">
    <span class="newCollectionHeader">Name</span>
    <input id="inputCollectionName" class="mainInputStyle" type="text" placeholder="Kollektion" bind:value={name}>
    <button class="standardButton" onclick={postNewCollection}>Hinzufügen</button>
</main>


<style>
    .newCollectionContent{
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        gap: 5pt;
        color: var(--text-color-secondary);

        input:focus {
            border: 1.5pt solid var(--success-color) !important;
        }

        button {
            background-color: var(--success-color) !important;
            padding: 7.5pt !important;
        }
    }

</style>