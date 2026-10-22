<script>
    import '../style/Collection.css';
    let {name} = $props();
</script>

<div class="collection">
    <div class="name">{name}</div>
    <div class="stats">
        coming soon...
    </div>
    <div class="actions">
        <div class="modifyButtons">
            <button class="standardButton"><i class="fa-solid fa-pen"></i></button>
            <button class="dangerButton"><i class="fa-solid fa-trash"></i></button>
        </div>
        <button class="standardButton">Lernen</button>
    </div>
</div>