<script>
    import '../style/LearningPage.css'
    import Topbar from '../lib/Topbar.svelte'
    import Card from '../lib/Card.svelte'

    let name = "Informatik"
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
        constructor(question, answer) {
            this.question = question;
            this.answer = answer;
        }
    }

    const cards = [
        new CardElement("Was sind algorithmen", "Rezepte für Programme"),
        new CardElement("Was ist Bogo Sort", "Effizienteste Sortieralgorithmus der Welt"),
        new CardElement("Was ist die beste Progammiersprache", "Python"),
    ]

    let currentCard = $state(0);

    function nextCard(){
        currentCard += 1;
    }

    function previousCard(){
        currentCard -= 1;
    }

</script>

<div class="main">
    <Topbar/>
    <div class="pageContent">
        <Card collection={name} text={cards[currentCard].question} answer={cards[currentCard].answer} site={site}/>

        <div class="actions">
            <button class="standardButton" onclick={previousCard}><i class="fa-solid fa-angle-left fa-xl"></i></button>
            <button class="standardButton" onclick={changeSite}><i class="fa-solid fa-rotate fa-xl"></i></button>
            <button class="standardButton" onclick={nextCard}><i class="fa-solid fa-angle-right fa-xl"></i></button>
        </div>
    </div>
</div>

