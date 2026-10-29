<script>
    import '../style/login.css';
    import {fade, fly} from "svelte/transition";
    import {getContext} from "svelte";

    const showInfo = getContext('showInfo');
    let registration = $state(0);
    let username = $state("")
    let password = $state("")
    let passwordAgain = $state("");
    let passwordMinlen = $state(0);
    let passwordUpper = $state(0);
    let passwordNumber = $state(0);
    let passwordMatch = $state(0);
    let usernameValid = $state(0);

    async function register(){
        let validPass = 1;
        passwordMinlen = 0;
        passwordNumber = 0;
        passwordUpper = 0;
        passwordMatch = 0;
        usernameValid = 0;

        if(username.length < 5){
            usernameValid = -1
        }

        if(!checkPassLen(password)){
            passwordMinlen = -1;
            validPass = 0;
        }

        if(!checkPassNum(password)){
            passwordNumber = -1;
            validPass = 0;
        }

        if(!checkPassUppercase(password)){
            passwordUpper = -1;
            validPass = 0;
        }

        if(validPass === 0 || usernameValid === -1){
            return
        }

        if(registration === 0){
            registration++
            return
        }

        if(password != passwordAgain){
            passwordMatch = -1;
            return
        }

        const res = await fetch("/api/auth/register", {
            method: 'POST',
            credentials: 'include',
            headers: {
                "Content-Type":"application/json"
            },
            body: JSON.stringify({
                username: username,
                password: password
            })
        });
        if(res.ok){
            registration = 0;
            showInfo("Succesfully registered", 0);
        }
        else{
            const data = await res.json()
            showInfo(data.message, 1);
        }
    }

    function checkPassLen(pass){
        if(pass.length < 8){
            return false;
        }
        return true;
    }

    function checkPassNum(pass){
        return /\d/.test(pass);
    }

    function checkPassUppercase(pass){
        return /[A-Z]/.test(password);
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
        <input type="text" placeholder="Benutzername" class="mainInputStyle" bind:value={username}>
        <div class="passReq">
            {#if usernameValid === -1}
                <span in:fly={{ y: -10, duration: 250, opacity: 0 }}>Mindestens 5 Character</span>
            {/if}
        </div>
        <span class="loginHeader">Passwort</span>
        <input type="password" placeholder="Passwort" class="mainInputStyle"  bind:value={password}>
        <div class="passReq">
            {#if passwordMinlen === -1}
                <span in:fly={{ y: -10, duration: 250, opacity: 0 }}>Mindestens 8 Character</span>
            {/if}
            {#if passwordNumber === -1}
                <span in:fly={{ y: -10, duration: 250, opacity: 0 }}>Mindestens eine Nummer</span>
            {/if}
            {#if passwordUpper === -1}
                <span in:fly={{ y: -10, duration: 250, opacity: 0 }}>Mindestens eine Großbuchstabe</span>
            {/if}
        </div>
        {#if registration === 1}
            <div class="register" in:fly={{ y: -20, duration: 500, opacity: 0 }} >
                <span class="loginHeader">Passwort Wiederholen</span>
                <input type="password" placeholder="Passwort" class="mainInputStyle" bind:value={passwordAgain}>
                <div class="passReq">
                    {#if passwordMatch === -1}
                        <span in:fly={{ y: -10, duration: 250, opacity: 0 }}>Passwörter sind nicht gleich</span>
                    {/if}
                </div>
            </div>
        {/if}
        <div class="loginAction">
            <button class="standardButton">Anmelden</button>
            <button class="accentButton" onclick={register}>Registrieren</button>
        </div>
    </div>
</div>