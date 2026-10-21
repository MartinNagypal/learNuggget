<script>
    import '../style/Collection.css';
    let name = "Algorithmen"
    let count = 1;
    let progress = 67;
</script>

<div class="collection">
    <div class="name">{name}</div>
    <div class="stats">
        <div>{count} Karte</div>
        <div style="color: var(--success-color)">{count} Richtig</div>
        <div style="color: var(--failure-color)">{count} Falsch</div>
    </div>
    <div class="actions">
        <button class="standardButton">Lernen</button>
    </div>
</div>