<script>
    import '../style/Homepage.css';
    import Group from "../lib/Group.svelte";
    import {getContext} from "svelte";
    import {fade} from "svelte/transition";

    const showInfo = getContext('showInfo');
    let name = $state("")
    let {learningpage} = $props();
    let groups = $state([])

    class GroupElement {
        constructor(groupName, groupId) {
            this.name = groupName;
            this.id = groupId;
        }
    }

    async function getUsername(){
        const res = await fetch("/api/auth/username");
        const data = await res.json();
        name = data.username;
    }

    getUsername();

    async function getGroups(){
        const res = await fetch("/api/group/fetchall");
        const data = await res.json();
        if (!res.ok) return;
        const newGroups = data.map(group => ({
            name: group.name,
            id: group.groupID
        }));
        if (JSON.stringify(groups) !== JSON.stringify(newGroups)) {
            groups = newGroups;
        }
    }

    getGroups();

    // Counted up after a group was created or deleted. Every Group watches it
    // and reloads its collections, so collections moved to "Default" when
    // their group is deleted show up right away.
    let collectionsVersion = $state(0);

    async function reloadGroupsAndCollections(){
        await getGroups();
        collectionsVersion++;
    }

</script>

<main class="homepage" in:fade={{duration: 200}}>
    <div class="welcome">
        <span>Welcome back</span>
        <span style="color:var(--success-color)">{name}</span>
        <span class="welcomeOrange">continue practicing!</span>
    </div>
    <div class="collectionsHeader"><h4>Collections</h4></div>
    <div class="collectionsContainer">
        {#each groups as group (group.id)}
            <Group name={group.name} id={group.id} create={0} learningPage={learningpage} reloadGroups={reloadGroupsAndCollections} collectionsVersion={collectionsVersion}/>
        {/each}
        <Group name="Default" id={0} create={0} learningPage={learningpage} reloadGroups={reloadGroupsAndCollections} collectionsVersion={collectionsVersion}/>
        <Group name="New Group" id={0} create={1} learningPage={learningpage} reloadGroups={reloadGroupsAndCollections}/>
    </div>
</main>