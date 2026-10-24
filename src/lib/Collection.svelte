<script>
    import '../style/Collection.css';
    let {name, learningpage} = $props();
</script>

<div class="collection">
    <div class="name">{name}</div>
    <div class="stats">
        coming soon...
    </div>
    <div class="collectionActions">
        <div class="modifyButtons">
            <button class="standardButton"><i class="fa-solid fa-pen"></i></button>
            <button class="dangerButton"><i class="fa-solid fa-trash"></i></button>
        </div>
        <button class="standardButton" onclick={learningpage}>Lernen</button>
    </div>
</div>