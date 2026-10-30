<script>
    import '../style/LearningPage.css'
    import Card from '../lib/Card.svelte'
    import ProgressReport from "../lib/ProgressReport.svelte";

    let {name, id, gameMode} = $props();

    let site = $state(0);
    let delCount = $state(0)
    let cards = $state([]);
    let finished = $state(0);
    let known = $state(0);
    let wrong = $state(0);
    let currentCard = $state(0);
    let direction = $state(1);
    let progress = $state(0);
    let len = $state(0);

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

    async function fetchCards() {
        const params = new URLSearchParams({collectionID: id});
        const res = await fetch(`/api/cards?${params}`)
        console.log(res);
        const data = await res.json()
        cards = data.map((card, index) => ({
            question: card.question,
            answer: card.answer,
            cardId: card.cardId,
            originalIndex: index
        }));
        len = cards.length;
    }

    addEventListener("keydown", (event) => {
        if(gameMode === 2){
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
        }

    })

    fetchCards();

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

    function toggleKnowCard(){
        if(gameMode === 0){
            cards.splice(currentCard, 1);
            if (cards.length > 0) {
                currentCard = Math.min(currentCard, cards.length - 1);
            }
            site = 0;
            known++;
            progress = known;
            if(cards.length === 0){
                finished = 1;
            }
        }
        else if(gameMode === 1){
            cards.splice(currentCard, 1);
            site = 0;
            known++;
            progress++;
            if(cards.length === 0){
                finished = 1;
            }
        }

    }

    function toggleWrong() {
        if(gameMode === 0){
            wrong++;
            nextCard();
        }
        else if(gameMode === 1){
            cards.splice(currentCard, 1);
            site = 0;
            wrong++;
            progress++;
            if(cards.length === 0){
                finished = 1;
            }
        }

    }

    let startX;
    const threshold = 50;

    document.addEventListener("touchstart", (event) => {
        startX = event.touches[0].clientX;
    });

    document.addEventListener("touchend", (event) => {
        if(gameMode == 2){
            if (startX === undefined) return;

            const endX = event.changedTouches[0].clientX;
            const distance = endX - startX;

            if (distance < -threshold) {
                nextCard();
            } else if (distance > threshold) {
                previousCard();
            }

            startX = undefined;
        }

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
                {#key cards[currentCard].cardId}
                    <Card mode={gameMode} text={cards[currentCard].question} answer={cards[currentCard].answer} site={site} changeSite={changeSite} direction={direction} current={(cards[currentCard].originalIndex)+1} len={len} progress={progress} toggleKnownCard={toggleKnowCard} toggleWrong={toggleWrong}/>
                {/key}
            {:else if finished === 1}
                <ProgressReport mode={gameMode} len={len} wrongCnt={wrong} rightCnt={known} name={name} id={id}/>
            {:else}
                <span>Diese Kollektion enthält noch keine Karten.</span>
            {/if}
        </div>

        {#if cards.length != 0}
            <div class="learningActions">
                <button class="dangerButton" onclick={deleteCard}><i class="fa-solid fa-trash"></i></button>
                {#if gameMode === 2}
                    <button class="standardButton" onclick={previousCard}><i class="fa-solid fa-angle-left fa-xl"></i></button>
                {:else}
                    <button class="inactiveButton"><i class="fa-solid fa-angle-left fa-xl"></i></button>
                {/if}
                    <button class="standardButton" onclick={changeSite}><i class="fa-solid fa-rotate fa-xl"></i></button>
                {#if gameMode === 2}
                    <button class="standardButton" onclick={nextCard}><i class="fa-solid fa-angle-right fa-xl"></i></button>
                {:else}
                    <button class="inactiveButton"><i class="fa-solid fa-angle-right fa-xl"></i></button>
                {/if}
                <button class="accentButton"><i class="fa-regular fa-star"></i></button>
            </div>
        {/if}
    </div>
</div>

