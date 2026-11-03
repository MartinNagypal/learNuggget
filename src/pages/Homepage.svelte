<script>
    import '../style/Homepage.css';
    import Group from "../lib/Group.svelte";
    import {getContext} from "svelte";

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

</script>

<main class="homepage">
    <div class="welcome">
        <span>Hallo <span style="color:var(--success-color)">{name}</span></span>
        <span class="welcomeOrange">Übe weiter</span>
        <span> da wo du aufgehört hast!</span>
    </div>
    <div class="collectionsHeader"><h4>Kollektionen </h4></div>
    <div class="collectionsContainer">
        {#each groups as group (group.id)}
            <Group name={group.name} id={group.id} create={0} learningPage={learningpage} reloadGroups={getGroups}/>
        {/each}
        <Group name="Standard" id={0} create={0} learningPage={learningpage} reloadGroups={getGroups}/>
        <Group name="Neue Gruppe" id={0} create={1} learningPage={learningpage} reloadGroups={getGroups}/>
    </div>
</main>