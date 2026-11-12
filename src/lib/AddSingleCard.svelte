<script>
    import {getContext} from "svelte";
    const showInfo = getContext('showInfo');

    let {CollectionID, toggleAddSingleCardActive, reloadCollectionStats} = $props();
    let Question = $state("");
    let Answer = $state("");
    let hasError = $state(0);

    async function postNewSingleCard(){
        // same check as the backend (ValidateCard)
        if(Question.trim() === "" || Answer.trim() === ""){
            hasError = 1;
            showInfo("Question and Answer should not be empty", 1);
            return;
        }
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
            const errData = await res.json().catch(() => ({message: "Card could not be saved"}));
            showInfo(errData.message, 1);
            return;
        }
        reloadCollectionStats();
        hasError = 0;
        toggleAddSingleCardActive();

    }

</script>

<div class="addSingleCardContent" class:addSingleCardError={hasError}>
    <span class="AddSingleCardCategoryName">Question</span>
    <input class="mainInputStyle" placeholder="Question" type="text" bind:value={Question}>
    <span class="AddSingleCardCategoryName">Answer</span>
    <input class="mainInputStyle" placeholder="Answer" type="text" bind:value={Answer}>
    <button class="successButton" onclick={postNewSingleCard}>Add</button>
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