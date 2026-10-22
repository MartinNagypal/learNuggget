<script>
    import '../style/Card.css';
    let {collection, text, answer, site, nextCardRight, nextCardLeft} = $props();


</script>

<div class="card">
    {#if site === 0}
        <div class="content">
            <div class="cardTop">{collection}</div>
            <div class="cardContent">{text}</div>
            <div class="cardBottom">
                Tippe zum Umdrehen
            </div>
        </div>
    {:else}
        <div class="content">
            <div class="cardTop">{collection}</div>
            <div class="question">{text}</div>
            <div class="answer"> {answer}</div>
            <div class="actions">
                <button class="failedButton button">Falsch</button>
                <button class="successButton button">Richtig</button>
            </div>
        </div>
    {/if}

</div>