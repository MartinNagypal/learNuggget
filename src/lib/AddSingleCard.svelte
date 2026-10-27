<script>

    let {CollectionID, toggleAddSingleCardActive} = $props();
    let Question = $state("");
    let Answer = $state("");

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

        const resData = await res.json();
        console.log(res)
        toggleAddSingleCardActive()
    }

</script>

<div class="addSingleCardContent">
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
    }
    .AddSingleCardCategoryName{
        font-family: var(--font-header);
        color: var(--accent-secondary);
    }
</style>