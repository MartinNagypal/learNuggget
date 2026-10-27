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
        }
        else {
            closeNewCollection();
        }
    }

</script>


<main id="newCollectionContent" class="newCollectionContent">
    <span class="newCollectionHeader">Name</span>
    <div class="newCollectionForm">
        <input id="inputCollectionName" class="mainInputStyle" type="text" placeholder="Kollektion" bind:value={name}>
        <button class="successButton" onclick={postNewCollection}>+</button>
        <button class="dangerButton" onclick={(event) => { event.stopPropagation(); closeNewCollection();}}>x</button>
    </div>
</main>


<style>
    .newCollectionContent{
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        gap: 5pt;
        color: var(--text-color-secondary);
        transition: all ease 0.3s;

        input {
            height: 36px;
            box-sizing: border-box;
        }

        input:focus {
            border: 1.5pt solid var(--success-color) !important;
        }

        button {
            padding: 5pt !important;
            padding-left: 10pt !important;
            padding-right: 10pt !important;
            box-sizing: border-box !important;
            height: 35px !important;
            text-align: center;
        }
    }

    .newCollectionForm{
        display: flex;
        flex-direction: row;
        gap: 5pt;
        justify-content: center;
        align-items: center;
    }

</style>