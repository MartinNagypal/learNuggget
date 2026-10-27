<script>

    let {id, toggleShowCollectionSettings} = $props();
    let text = $state("")

    async function addBulkCards(){
        const res = await fetch(`/api/addBulk?id=${id}`,{
            method: 'POST',
            credentials: 'include',
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                bulk: text
            })
        });
        if(res.ok){
            toggleShowCollectionSettings();
        }
    }

</script>

<div class="collectionSettings">
    <div class="collectionSettingsHeader">Bulk Hinzufügen</div>
    <textarea bind:value={text}></textarea>
    <button class="successButton" onclick={addBulkCards}>Hinzufügen</button>
</div>

<style>
    .collectionSettings{
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

    .collectionSettingsHeader{
        color: var(--accent-secondary);
        font-family: var(--font-header);
    }

    .collectionSettings textarea {
        border: 1.5pt solid var(--border-color);
        border-radius: var(--border-radius-primary);
        padding: 5pt;
        box-sizing: border-box;
        transition: border-color ease 0.3s;
        width: 90%;
        max-width: 90%;
        max-height: 100%;
    }

    .collectionSettings textarea:focus {
        outline: none;
        border-color: var(--accent-secondary);
    }
</style>