<script>

    let {CollectionID, toggleAddSingleCardActive, reloadCollectionLen} = $props();
    let Question = $state("");
    let Answer = $state("");
    let hasError = $state(0);

    async function postNewSingleCard(){
        const data = {
            question: Question,
            answer: Answer,
            collectionID: CollectionID
        };
        const res = await fetch("/api/addSingleCard", {
            method: "POST",
            credentials: "include",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(data)
        });

        if(!res.ok){
            hasError = 1;
            return;
        }
        await reloadCollectionLen();
        hasError = 0;
        toggleAddSingleCardActive();

    }

</script>

<div class="addSingleCardContent" class:addSingleCardError={hasError}>
    <span class="AddSingleCardCategoryName">Frage</span>
    <input class="mainInputStyle" placeholder="Frage" type="text" bind:value={Question}>
    <span class="AddSingleCardCategoryName">Antwort</span>
    <input class="mainInputStyle" placeholder="Antwort" type="text" bind:value={Answer}>
    <button class="successButton" onclick={postNewSingleCard}>Hinzufügen</button>
</div>

<style>
    .addSingleCardContent{
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 5pt;
        border: 1.5pt solid var(--border-color);
        border-radius: var(--border-radius-primary);

        padding: 5pt;
        margin-top: 5pt;
        transition: all ease-out 0.25s;
    }
    .AddSingleCardCategoryName{
        font-family: var(--font-header);
        color: var(--accent-secondary);
    }
    .addSingleCardError {
        border-color: var(--failure-color) !important;
    }
</style>