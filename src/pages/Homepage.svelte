<script>
    import '../style/Homepage.css';
    import Collection from '../lib/Collection.svelte'
    import NewCollection from "../lib/NewCollection.svelte";
    import { fade } from 'svelte/transition';
    import { scale } from 'svelte/transition';
    import Group from "../lib/Group.svelte";
    import {get} from "svelte/store";

    let name = "Martin"
    let {learningpage} = $props();
    let groups = $state([])

    class GroupElement {
        constructor(groupName, groupId) {
            this.name = groupName;
            this.id = groupId;
        }
    }

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
        <Group name="Neue Gruppe" id="2" create={1} learningPage={learningpage} reloadGroups={getGroups}/>
    </div>
</main>