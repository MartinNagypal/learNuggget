<script>
    import '../style/LearningPage.css'
    import Card from '../lib/Card.svelte'

    let {name, id} = $props();

    let card = "Was ist Skefkenlnln<swjldna wkjndj kwand nwajdn wajkbdkwabd aortingmaorting aorting aorting"
    let antwort ="aortingaorting aorting aorting aorting aortingaorting aorting aorting aorting"
    let site = $state(0);
    let delCount = $state(0)

    function changeSite(){
        if(site === 0){
            site = 1;
        }
        else {
            site = 0;
        }
    }

    class CardElement {
        constructor(cardQuestion, cardAnswer, cardId) {
            this.question = cardQuestion;
            this.answer = cardAnswer;
            this.cardId = cardId;
        }
    }

    let cards = $state([]);

    async function fetchCards() {
        const params = new URLSearchParams({collectionID: id});
        const res = await fetch(`/api/cards?${params}`)
        console.log(res);
        const data = await res.json()
        cards = data.map(card => ({
            question: card.question,
            answer: card.answer,
            cardId: card.cardId
        }));
        console.log(cards)
    }

    addEventListener("keydown", (event) => {
        if(event.key === " "){
            changeSite();
            return
        }
        if(event.key === "ArrowRight"){
            nextCard();
            return;
        }
        if(event.key === "ArrowLeft"){
            previousCard();
            return;
        }
    })

    fetchCards();
    let currentCard = $state(0);
    let direction = $state(1);

    function nextCard() {
        direction = 1;
        site = 0;
        currentCard = (currentCard + 1) % cards.length;
        delCount = 0;
    }

    function previousCard() {
        direction = -1;
        site = 0;
        currentCard = (currentCard - 1 + cards.length) % cards.length;
        delCount = 0;
    }

    let startX;
    const threshold = 50;

    document.addEventListener("touchstart", (event) => {
        startX = event.touches[0].clientX;
    });

    document.addEventListener("touchend", (event) => {
        if (startX === undefined) return;

        const endX = event.changedTouches[0].clientX;
        const distance = endX - startX;

        if (distance < -threshold) {
            nextCard();
        } else if (distance > threshold) {
            previousCard();
        }

        startX = undefined;
    });

    async function deleteCard(){
        if(delCount === 0){
            delCount++
            return
        }
        delCount = 0;
        const res = await fetch(`/api/card/delete?id=${cards[currentCard].cardId}`, {
            method: 'DELETE'
        })
        if(res.ok){
            await fetchCards();
            nextCard();
        }
    }

</script>

<div class="main">
    <div class="learningPageContent">
        <div class="cardSlot">
            {#if cards.length != 0}
                {#key currentCard}
                    <Card collection={name} text={cards[currentCard].question} answer={cards[currentCard].answer} site={site} changeSite={changeSite} direction={direction} current={currentCard+1} len={cards.length}/>
                {/key}
            {:else}
                <span>Diese Kollektion enthält bisher noch keine Karten.</span>
            {/if}
        </div>

        {#if cards.length != 0}
            <div class="learningActions">
                <button class="dangerButton" onclick={deleteCard}><i class="fa-solid fa-trash"></i></button>
                <button class="standardButton" onclick={previousCard}><i class="fa-solid fa-angle-left fa-xl"></i></button>
                <button class="standardButton" onclick={changeSite}><i class="fa-solid fa-rotate fa-xl"></i></button>
                <button class="standardButton" onclick={nextCard}><i class="fa-solid fa-angle-right fa-xl"></i></button>
                <button class="accentButton"><i class="fa-regular fa-star"></i></button>
            </div>
        {/if}
    </div>
</div>

