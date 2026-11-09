<script>
    import {fade} from "svelte/transition";
    import {getContext} from "svelte";
    const setPage = getContext('setPage');
    let selected = $state(0);

    function select(id){
        selected = id;
        if(selected === 0){
            setPage(0);
        }
        else if(id === 2){
            setPage(2);
        }
        else if (id === 1){
            setPage(3);
        }
    }

</script>
<div class="navBar">
    {#key selected}
        <div class="navBarElements">
            <div class={selected === 0 ? "navBarSelected" : "navBarUnselected"} onclick={() => {select(0)}}> <i class="fa-solid fa-house"></i> {selected === 0 ? "Home" : ""}</div>
            <div class={selected === 1 ? "navBarSelected" : "navBarUnselected"} onclick={() => {select(1)}}><i class="fa-solid fa-graduation-cap fa-lg"></i> {selected === 1 ? "Learn" : ""}</div>
            <div class={selected === 2 ? "navBarSelected" : "navBarUnselected"} onclick={() => {select(2)}}><i class="fa-solid fa-people-group fa-lg"></i> {selected === 2 ? "Social" : ""}</div>
        </div>
    {/key}
</div>

<style>
    .navBar{
        display: flex;
        width: 100dvw;
        align-items: center;
        justify-content: center;
    }
    .navBarElements{
        display: grid;
        grid-template-columns: repeat(3, minmax(0, 1fr));
        align-items: center;
        text-align: center;
        box-sizing: border-box;
        width: 300pt;
        border: 1.5pt solid var(--border-color);
        border-radius: var(--border-radius-primary);
        padding: 7.5pt 10pt;
        margin-bottom: 10pt;
        margin-right: 10pt;
        margin-left: 10pt;
        backdrop-filter: blur(20px);
        user-select: none;
    }

    .navBarSelected{
        border-radius: calc(var(--border-radius-primary) - 3pt);
        background-color: var(--accent-secondary);
        color: var(--text-color-white);
        font-family: var(--font-header);
        padding: 5pt 5pt;
    }

    .navBarUnselected {
        border-radius: calc(var(--border-radius-primary) - 3pt);
        color: var(--accent-secondary);
    }

</style>