<script>
    import '../style/LearningPage.css'
    import Card from '../lib/Card.svelte'

    let {name, id} = $props();

    let card = "Was ist Skefkenlnln<swjldna wkjndj kwand nwajdn wajkbdkwabd aortingmaorting aorting aorting"
    let antwort ="aortingaorting aorting aorting aorting aortingaorting aorting aorting aorting"
    let site = $state(0);

    function changeSite(){
        if(site === 0){
            site = 1;
        }
        else {
            site = 0;
        }
    }

    class CardElement {
        constructor(cardQuestion, cardAnswer) {
            this.question = cardQuestion;
            this.answer = cardAnswer;
        }
    }

    let cards = $state([]);

    async function fetchCards() {
        const params = new URLSearchParams({collectionID: id});
        const res = await fetch(`/api/cards?${params}`)
        console.log(res);
        const data = await res.json()
        cards = data.map(CardElement => ({
            question: CardElement.question,
            answer: CardElement.answer
        }));
    }

    fetchCards();
    let currentCard = $state(0);
    let direction = $state(1);

    function nextCard() {
        direction = 1;
        site = 0;
        currentCard = (currentCard + 1) % cards.length;
    }

    function previousCard() {
        direction = -1;
        site = 0;
        currentCard = (currentCard - 1 + cards.length) % cards.length;
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

</script>

<div class="main">
    <div class="learningPageContent">
        <div class="cardSlot">
            {#if cards.length != 0}
                {#key currentCard}
                    <Card collection={name} text={cards[currentCard].question} answer={cards[currentCard].answer} site={site} changeSite={changeSite} direction={direction}/>
                {/key}
            {:else}
                <span>Diese Kollektion enthält bisher noch keine Karten.</span>
            {/if}
        </div>

        {#if cards.length != 0}
            <div class="learningActions">
                <button class="standardButton" onclick={previousCard}><i class="fa-solid fa-angle-left fa-xl"></i></button>
                <button class="standardButton" onclick={changeSite}><i class="fa-solid fa-rotate fa-xl"></i></button>
                <button class="standardButton" onclick={nextCard}><i class="fa-solid fa-angle-right fa-xl"></i></button>
            </div>
        {/if}
    </div>
</div>

