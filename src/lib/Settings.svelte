<script>
    import {getContext} from "svelte";
    import { fly } from "svelte/transition";
    const logout = getContext('logout');
    const toggleShowAccountDelete = getContext('toggleConfirmAccountDeletion');
    const showInfo = getContext('showInfo');
    const toggleIsLoggedIn = getContext('toggleIsLoggedIn');
    const isRoot = getContext('isRoot');
    const setPage = getContext('setPage');
    let showSettings = $state(0);

    function toggleShowsettings(){
        showSettings = (showSettings +1) % 2;
    }

    async function logoutFromAllSession(){
        const res = await fetch ("/api/auth/logout/all");
        const data = await res.json();
        if(!res.ok){
            showInfo(data.message, 1);
            return;
        }
        showInfo(data.message, 0);
        toggleIsLoggedIn();
    }

</script>
<div class="settings">
    <div class="settingsHeader" onclick={toggleShowsettings}>Settings</div>
    {#if showSettings === 1}
        <div class="settingsElements" out:fly={{ y: 40, duration: 300, opacity: 0 }} in:fly={{ y: -40, duration: 200, opacity: 0 }}>
            <div onclick={toggleShowAccountDelete} class="settingSelectionElement">Delete Account</div>
            <div onclick={logout} class="settingSelectionElement">Logout</div>
            <div class="settingSelectionElement" onclick={logoutFromAllSession}>Logout from all Sessions</div>
            {#if isRoot()}
                <div class="settingSelectionElement adminElement" onclick={() => setPage(4)}>Admin Panel</div>
            {/if}
        </div>
    {/if}
</div>

<style>
    .settings, .settingsElements{
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        color: var(--text-color-secondary);
        gap: 7.5pt;
        user-select: none;
    }

    .settingsElements{
        margin-bottom: 20pt;
    }

    .settingsHeader{
        font-family: var(--font-header);
        border-bottom: 1.5pt solid var(--border-color);
        padding: 0 50px;
        margin-bottom: 2.5pt;
    }

    .settingSelectionElement{
        cursor: pointer;
    }

    .adminElement{
        color: var(--accent-primary);
    }
</style>