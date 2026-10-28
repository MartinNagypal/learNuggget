<script>
    import '../style/login.css';
    import {fade, fly} from "svelte/transition";

    let registration = $state(0);

    function register(){
        if(registration === 0){
            registration++
        }
    }
</script>

<div class="login">
    <div class="loginContent">
        <div class="welcomeText">
            <span>Wilkommen</span>
            <span>bei</span>
            <img src="/logo.png" width="200px">
        </div>
        <span class="loginHeader">Benutzername</span>
        <input type="text" placeholder="Benutzername" class="mainInputStyle">
        <span class="loginHeader">Passwort</span>
        <input type="password" placeholder="Passwort" class="mainInputStyle">
        {#if registration === 1}
            <div class="register" in:fly={{ y: -20, duration: 500, opacity: 0 }} >
                <span class="loginHeader">Passwort Wiederholen</span>
                <input type="password" placeholder="Passwort" class="mainInputStyle">
            </div>
        {/if}
        <div class="loginAction">
            <button class="standardButton">Anmelden</button>
            <button class="accentButton" onclick={register}>Registrieren</button>
        </div>

    </div>

</div>