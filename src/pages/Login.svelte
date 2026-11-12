<script>
    import '../style/login.css';
    import {fade, fly} from "svelte/transition";
    import {getContext} from "svelte";

    const verifySession = getContext('verifySession');
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

        // Same rules as the server (NormalizeUsername in backend/functions/validation.go)
        username = username.trim();
        if(!/^[A-Za-z0-9_.-]{5,32}$/.test(username)){
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
            showInfo("Successfully registered", 0);
            username = "";
            password = "";
            await verifySession();
        }
        else{
            const data = await res.json();
            showInfo(data.message, 1);
        }
    }

    function checkPassLen(pass){
        if(pass.length < 8 || pass.length > 32){
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

    async function login(){
        const res = await fetch("/api/auth/login", {
            method: 'POST',
            credentials: 'include',
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                username: username,
                password: password
            })
        });
        const data = await res.json();
        if(!res.ok){
            showInfo(data.message, 1);
            password = "";
            return
        }
        showInfo(data.message, 0);
        username = "";
        password = "";
        await verifySession();
    }

</script>

<div class="login">
    <div class="loginContent">
        <div class="welcomeText">
            <span>Welcome</span>
            <span>to</span>
            <img src="/logo.png" width="200px">
        </div>
        <span class="loginHeader">Username</span>
        <input type="text" placeholder="Username" class="mainInputStyle" bind:value={username}>
        <div class="passReq">
            {#if usernameValid === -1}
                <span in:fly={{ y: -10, duration: 250, opacity: 0 }}>5 - 32 Character: Letters, Numbers, _ . -</span>
            {/if}
        </div>
        <span class="loginHeader">Password</span>
        <input type="password" placeholder="Password" class="mainInputStyle"  bind:value={password}>
        <div class="passReq">
            {#if passwordMinlen === -1}
                <span in:fly={{ y: -10, duration: 250, opacity: 0 }}>Length between 8 - 32 Character</span>
            {/if}
            {#if passwordNumber === -1}
                <span in:fly={{ y: -10, duration: 250, opacity: 0 }}>At least one Number</span>
            {/if}
            {#if passwordUpper === -1}
                <span in:fly={{ y: -10, duration: 250, opacity: 0 }}>At least one uppercase and one lowercase letter</span>
            {/if}
        </div>
        {#if registration === 1}
            <div class="register" in:fly={{ y: -20, duration: 500, opacity: 0 }} >
                <span class="loginHeader">Repeat password</span>
                <input type="password" placeholder="Password" class="mainInputStyle" bind:value={passwordAgain}>
                <div class="passReq">
                    {#if passwordMatch === -1}
                        <span in:fly={{ y: -10, duration: 250, opacity: 0 }}>Passwords do not match</span>
                    {/if}
                </div>
            </div>
        {/if}
        <div class="loginAction">
            <button class="standardButton" onclick={login}>Login</button>
            <button class="accentButton" onclick={register}>Register</button>
        </div>
    </div>
</div>