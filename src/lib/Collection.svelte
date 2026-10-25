<script>
    import '../style/Collection.css';
    let {name, learningpage} = $props();

    function printName(){
        console.log("this is " + name);
    }

</script>

<div class="collection">
    <div class="name">{name}</div>
    <div class="collectionActions">
        <div class="modifyButtons">
            <button class="standardButton" onclick={printName}><i class="fa-solid fa-pen"></i></button>
            <button class="dangerButton"><i class="fa-solid fa-trash"></i></button>
        </div>
        <button class="standardButton" onclick={learningpage}>Lernen</button>
    </div>
</div>