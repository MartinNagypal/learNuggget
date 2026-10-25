<script>
    import '../style/LearningPage.css'
    import Card from '../lib/Card.svelte'

    let {name} = $props();

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

    const cards = [
        new CardElement("Was sind algorithmen", "Rezepte für Programme"),
        new CardElement("Was ist Bogo Sort", "Effizienteste Sortieralgorithmus der Welt"),
        new CardElement("Was ist die beste Progammiersprache", "Python"),
        new CardElement("Wovon hängt die Laufzeit ab", "Sie hängt von vielen Größen wie Hardware, Software, Implementierung, Anzahl der vergleiche ab oder auch Interne hardware wie der Scheduler der zwisxchendurch auch andere prozesse ablaufen lassen kann.Sie hängt von vielen Größen wie Hardware, Software, Implementierung, Anzahl der vergleiche ab oder auch Interne hardware wie der Scheduler der zwisxchendurch auch andere prozesse ablaufen lassen kann.\""),
    ]

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
            {#key currentCard}
                <Card collection={name} text={cards[currentCard].question} answer={cards[currentCard].answer} site={site} changeSite={changeSite} direction={direction}/>
            {/key}
        </div>

        <div class="learningActions">
            <button class="standardButton" onclick={previousCard}><i class="fa-solid fa-angle-left fa-xl"></i></button>
            <button class="standardButton" onclick={changeSite}><i class="fa-solid fa-rotate fa-xl"></i></button>
            <button class="standardButton" onclick={nextCard}><i class="fa-solid fa-angle-right fa-xl"></i></button>
        </div>
    </div>
</div>

