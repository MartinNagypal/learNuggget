<script>
    import { fade } from "svelte/transition";
    import { getContext } from "svelte";

    const showInfo = getContext('showInfo');

    let stats = $state({users: 0, collections: 0, groups: 0, cards: 0});
    let users = $state([]);
    let delUserId = $state(0);
    let delUserTimer;

    async function getStats(){
        const res = await fetch("/api/admin/stats");
        const data = await res.json();
        if(!res.ok){
            showInfo(data.message, 1);
            return;
        }
        stats = data;
    }

    async function getUsers(){
        const res = await fetch("/api/admin/users");
        const data = await res.json();
        if(!res.ok){
            showInfo(data.message, 1);
            return;
        }
        users = data;
    }

    getStats();
    getUsers();

    function formatLastActive(timestamp){
        if(timestamp === 0){
            return "never";
        }
        const seconds = Math.floor(Date.now() / 1000) - timestamp;
        if(seconds < 60){
            return "just now";
        }
        if(seconds < 3600){
            return `${Math.floor(seconds / 60)} min ago`;
        }
        if(seconds < 86400){
            return `${Math.floor(seconds / 3600)} h ago`;
        }
        const days = Math.floor(seconds / 86400);
        return days === 1 ? "1 day ago" : `${days} days ago`;
    }

    async function logoutUser(user){
        const res = await fetch(`/api/admin/user/logout?id=${user.userID}`, {
            method: 'POST'
        });
        const data = await res.json();
        if(!res.ok){
            showInfo(data.message, 1);
            return;
        }
        showInfo(`${user.username} logged out`, 0);
        await getUsers();
    }

    async function deleteUser(user){
        if(delUserId !== user.userID){
            delUserId = user.userID;
            showInfo(`Click again to delete ${user.username}`, 1);
            // cancel the timer of a previous click (e.g. on another user)
            clearTimeout(delUserTimer);
            delUserTimer = setTimeout(() => {
                delUserId = 0;
            }, 3000);
            return;
        }
        clearTimeout(delUserTimer);
        delUserId = 0;
        const res = await fetch(`/api/admin/user/delete?id=${user.userID}`, {
            method: 'DELETE'
        });
        const data = await res.json();
        if(!res.ok){
            showInfo(data.message, 1);
            return;
        }
        showInfo(`${user.username} deleted`, 0);
        await getUsers();
        await getStats();
    }
</script>

<main class="adminPage" in:fade={{duration: 200}}>
    <section class="adminSection">
        <div class="sectionHeader"><i class="fa-solid fa-chart-simple"></i> Statistics</div>
        <div class="adminBox adminStats">
            <div class="adminStatElement">
                <div class="statHeader">{stats.users}</div>
                <div class="statFooter">Users</div>
            </div>
            <div class="adminStatElement">
                <div class="statHeader">{stats.collections}</div>
                <div class="statFooter">Collections</div>
            </div>
            <div class="adminStatElement">
                <div class="statHeader">{stats.groups}</div>
                <div class="statFooter">Groups</div>
            </div>
            <div class="adminStatElement">
                <div class="statHeader">{stats.cards}</div>
                <div class="statFooter">Cards</div>
            </div>
        </div>
    </section>

    <section class="adminSection">
        <div class="sectionHeader usersHeader">
            <span><i class="fa-solid fa-users"></i> Users</span>
            <span class="usersCount">{users.length}</span>
        </div>
        <div class="adminBox userList">
            {#each users as user (user.userID)}
                <div class="userRow">
                    <div>
                        <div>
                            <span class="userId">#{user.userID}</span>
                            <span class="userName">{user.username}</span>
                            {#if user.root}
                                <span class="rootBadge">Root</span>
                            {/if}
                        </div>
                        <div class="userInfo">
                            {formatLastActive(user.lastActive)} · {user.collections} {user.collections === 1 ? "collection" : "collections"} · {user.cards} {user.cards === 1 ? "card" : "cards"}
                        </div>
                    </div>
                    {#if !user.root}
                        <div class="userActions">
                            <button class="standardButton" onclick={() => logoutUser(user)}><i class="fa-solid fa-right-from-bracket"></i></button>
                            <button class="dangerButton" onclick={() => deleteUser(user)}><i class="fa-solid fa-trash"></i></button>
                        </div>
                    {/if}
                </div>
            {/each}
        </div>
    </section>
</main>

<style>
    .adminPage{
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 25pt;
        padding: 25pt 0;
    }

    .adminSection{
        width: 90%;
        max-width: 375px;
        display: flex;
        flex-direction: column;
        gap: 10pt;
    }

    .sectionHeader{
        font-family: var(--font-header);
        font-size: var(--font-size-middle);
        color: var(--accent-secondary);
        border-bottom: 1.5pt solid var(--border-color);
        padding: 2.5pt 7.5pt;

        i {
            color: var(--accent-primary);
        }
    }

    .usersHeader{
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .usersCount{
        font-family: var(--font-text);
        font-size: var(--font-size-small);
        color: var(--text-color-secondary);
    }

    .adminBox{
        background-color: var(--accent-surface);
        border: 1.25pt solid var(--border-color);
        border-radius: var(--border-radius-primary);
        padding: 10pt;
    }

    .adminStats{
        display: grid;
        grid-template-columns: repeat(4, 1fr);
        gap: 5pt;
    }

    .adminStatElement{
        display: flex;
        flex-direction: column;
        align-items: center;
        padding: 3pt 0;
    }

    .statHeader{
        font-family: var(--font-header);
        color: var(--accent-primary);
    }

    .statFooter{
        font-size: var(--font-size-small);
        color: var(--text-color-secondary);
    }

    .userList{
        padding-top: 0;
        padding-bottom: 0;
    }

    .userRow{
        display: flex;
        justify-content: space-between;
        align-items: center;
        gap: 10pt;
        padding: 8pt 0;
        border-bottom: 1pt solid var(--border-color);
    }

    .userRow:last-child{
        border-bottom: none;
    }

    .userId{
        color: var(--text-color-secondary);
        font-size: var(--font-size-small);
    }

    .userName{
        font-family: var(--font-header);
        color: var(--accent-secondary);
    }

    .rootBadge{
        background-color: var(--rank-gold-surface);
        color: var(--rank-gold-text);
        font-size: var(--font-size-small);
        border-radius: calc(var(--border-radius-primary) - 2pt);
        padding: 1pt 5pt;
    }

    .userInfo{
        font-size: var(--font-size-small);
        color: var(--text-color-secondary);
    }

    .userActions{
        display: flex;
        gap: 5pt;
        flex-shrink: 0;
    }
</style>
