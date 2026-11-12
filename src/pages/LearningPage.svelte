<script>
    import '../style/LearningPage.css'
    import Card from '../lib/Card.svelte'
    import ProgressReport from "../lib/ProgressReport.svelte";
    import { getContext } from "svelte";
    import { ALL_MARKED } from "../lib/markedCards.js";

    const showInfo = getContext('showInfo');

    let {name, id, gameMode, markedOnly = false} = $props();

    let loaded = $state(0);
    let site = $state(0);
    let delCount = $state(0)
    let delTimer;
    let cards = $state([]);
    let finished = $state(0);
    let known = $state(0);
    let wrong = $state(0);
    let currentCard = $state(0);
    let direction = $state(1);
    let progress = $state(0);
    let len = $state(0);
    let shuffled = $state(false);

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

    function shuffle(list) {
        for (let i = list.length - 1; i > 0; i--) {
            const j = Math.floor(Math.random() * (i + 1));
            [list[i], list[j]] = [list[j], list[i]];
        }
        return list;
    }

    async function fetchCards() {
        const params = new URLSearchParams({collectionID: id});
        const url = id === ALL_MARKED ? "/api/cards/marked" : `/api/cards?${params}`;
        const res = await fetch(url)
        console.log(res);
        const data = await res.json()
        const loadedCards = markedOnly ? data.filter(card => card.marked) : data;
        const newCards = loadedCards.map((card, index) => ({
            question: card.question,
            answer: card.answer,
            cardId: card.cardId,
            marked: card.marked,
            originalIndex: index
        }));
        cards = shuffled ? shuffle(newCards) : newCards;
        len = cards.length;
        loaded = 1;
    }
    fetchCards();

    function toggleOrder() {
        const currentId = cards[currentCard].cardId;
        shuffled = !shuffled;

        if (shuffled) {
            const current = cards[currentCard];
            const rest = cards.filter(card => card.cardId !== currentId);
            cards = [current, ...shuffle(rest)];
            currentCard = 0;
        } else {
            cards = [...cards].sort((a, b) => a.originalIndex - b.originalIndex);
            currentCard = cards.findIndex(card => card.cardId === currentId);
        }

        showInfo(shuffled ? "Random order" : "Original order", 0);
    }

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
            if (cards.length > 0) {
                currentCard = Math.min(currentCard, cards.length - 1);
            }
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
            if (cards.length > 0) {
                currentCard = Math.min(currentCard, cards.length - 1);
            }
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

    function handleKeydown(event) {
        if (gameMode === 2) {
            if (event.key === " ") { changeSite(); return; }
            if (event.key === "ArrowRight") { nextCard(); return; }
            if (event.key === "ArrowLeft") { previousCard(); return; }
        }
    }

    function handleTouchStart(event) {
        startX = event.touches[0].clientX;
    }

    function handleTouchEnd(event) {
        if (gameMode == 2) {
            if (startX === undefined) return;
            const distance = event.changedTouches[0].clientX - startX;
            if (distance < -threshold) nextCard();
            else if (distance > threshold) previousCard();
            startX = undefined;
        }
    }

    async function deleteCard(){
        if(gameMode === 2){
            if(delCount === 0){
                delCount++
                showInfo("Click again to delete this card", 1);
                clearTimeout(delTimer);
                delTimer = setTimeout(() => {
                    delCount = 0;
                }, 3000);
                return
            }
            clearTimeout(delTimer);
            delCount = 0;
            const res = await fetch(`/api/card/delete?id=${cards[currentCard].cardId}`, {
                method: 'DELETE'
            })
            if(res.ok){
                const deleted = cards[currentCard];
                cards.splice(currentCard, 1);
                for (const card of cards) {
                    if (card.originalIndex > deleted.originalIndex) {
                        card.originalIndex--;
                    }
                }
                len = cards.length;
                if (currentCard >= cards.length) {
                    currentCard = 0;
                }
                site = 0;
            }
            else {
                const data = await res.json();
                showInfo(data.message, 1);
            }
        }
        else{
            showInfo("Deleting cards is only available in Infinity mode", 1);
        }
    }

    async function markCard(cardId){
        const card = cards.find(c => c.cardId === cardId);
        if (!card) return;

        card.marked = !card.marked;

        const res = await fetch("/api/card/mark", {
            method: "PATCH",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ cardId: cardId, marked: card.marked })
        });
        if (!res.ok) {
            card.marked = !card.marked;
            showInfo("Failed to mark card", 1);
        }
    }

</script>



<svelte:window onkeydown={handleKeydown} />
<svelte:document ontouchstart={handleTouchStart} ontouchend={handleTouchEnd} />
<div class="main">
    {#if loaded === 1}
        <div class="learningPageContent">
            <div class="cardSlot">
                {#if cards.length != 0}
                    {#key cards[currentCard].cardId}
                        <Card mode={gameMode} text={cards[currentCard].question} answer={cards[currentCard].answer} site={site} changeSite={changeSite} direction={direction} current={(cards[currentCard].originalIndex)+1} len={len} progress={progress} toggleKnownCard={toggleKnowCard} toggleWrong={toggleWrong}/>
                    {/key}
                {:else if finished === 1}
                    <ProgressReport mode={gameMode} len={len} wrongCnt={wrong} rightCnt={known} name={name} id={id} markedOnly={markedOnly}/>
                {:else if markedOnly}
                    <span>No marked cards here.</span>
                {:else}
                    <span>This collection has no cards yet.</span>
                {/if}
            </div>

            {#if cards.length != 0}
                <div class="learningActions">
                    {#if gameMode === 2}
                        <button class="dangerButton" onclick={deleteCard}><i class="fa-solid fa-trash"></i></button>
                    {:else}
                        <button class="inactiveButton" onclick={deleteCard}><i class="fa-solid fa-trash"></i></button>
                    {/if}
                    {#if gameMode === 2}
                        <button class="standardButton" onclick={previousCard}><i class="fa-solid fa-angle-left fa-xl"></i></button>
                    {:else}
                        <button class="inactiveButton"><i class="fa-solid fa-angle-left fa-xl"></i></button>
                    {/if}
                        <button class="standardButton" onclick={toggleOrder} title={shuffled ? "Random order" : "In order"}><i class={shuffled ? "fa-solid fa-shuffle fa-xl" : "fa-solid fa-arrow-down-1-9 fa-xl"}></i></button>
                    {#if gameMode === 2}
                        <button class="standardButton" onclick={nextCard}><i class="fa-solid fa-angle-right fa-xl"></i></button>
                    {:else}
                        <button class="inactiveButton"><i class="fa-solid fa-angle-right fa-xl"></i></button>
                    {/if}
                    <button class="accentButton" onclick={() => markCard(cards[currentCard].cardId)}><i class={cards[currentCard].marked ? "fa-solid fa-star fa-lg" : "fa-regular fa-star fa-lg"}></i></button>
                </div>
            {/if}
        </div>
    {/if}
</div>

