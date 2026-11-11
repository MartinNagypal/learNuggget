<script>
    import {getContext} from "svelte";
    import {fade} from "svelte/transition";
    import {get} from "svelte/store";

    const toggleShowDeleteAccount = getContext('toggleConfirmAccountDeletion');
    const showInfo = getContext('showInfo');
    const toggleIsLoggedIn = getContext('toggleIsLoggedIn');
    let input = $state("");

    async function deleteAccount(){
        if (input === "CONFIRM"){
            const res = await fetch("/api/auth/deleteUser");
            const data = await res.json();
            if(!res.ok){
                showInfo(data.message, 1);
                toggleShowDeleteAccount();
                return;
            }
            showInfo(data.message, 0);
            toggleShowDeleteAccount();
            toggleIsLoggedIn();
        }
        else{
            showInfo("Wrong input", 1);
            input = "";
        }
    }

</script>

<div class="confirmDeleteAccount">
    <div class="confirmDeleteAccountElements">
        <div class="confirmDeleteAccountHeader">Are you sure?</div>
        <div style="color: var(--failure-color); font-size: var(--font-size-small);">WARNING: All your data will be deleted</div>
        <div style="color: var(--failure-color); font-size: var(--font-size-small);">including your collections, cards and progress.</div>
        <div class="confirmDeleteAccountSubtext"> type "CONFIRM" to delete your account.</div>
        <input  type="text" class="mainInputStyle" placeholder="CONFIRM" bind:value={input}>
        <div class="confirmDeleteAccountButtons">
            <button class="dangerButton" onclick={deleteAccount}>Confirm</button>
            <button class="standardButton" onclick={toggleShowDeleteAccount}>Cancel</button>
        </div>

    </div>

</div>

<style>
    .confirmDeleteAccount{
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        flex: 1;
    }

    .confirmDeleteAccountElements{
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 7.5pt;

        padding: 25pt;
        background-color: var(--accent-surface);
        border: 1.5pt solid var(--border-radius-primary);
        border-radius: var(--border-radius-primary);
    }

    .confirmDeleteAccountHeader{
        font-family: var(--font-header);
        color: var(--failure-color);
        font-size: var(--font-size-middle);
    }

    .confirmDeleteAccountSubtext{
        font-family: var(--font-header);
        color: var(--accent-primary);
        font-size: var(--font-size-small);
    }

    .confirmDeleteAccountButtons{
        display: grid;
        grid-template-rows: repeat(2,1fr);
        gap: 5pt;
    }
</style>