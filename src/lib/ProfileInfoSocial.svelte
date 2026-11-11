<script>
    let {name} = $props();
    let profileScore = $state(0);

    async function getProfileSCore(){
        const res = await fetch("/api/user/score");
        const data = await res.json();
        profileScore = data.score;
    }

    getProfileSCore();

</script>

<div class="profileInfoSocial">
    <div class="profileInfoSocialHeader">Hello,</div>
    <div class="profileInfoSocialName">{name}</div>
    <div class="profileInfoSocialScoreDiv">
        <div class="profileInfoSocialScore">Score</div>
        <div class="profileInfoSocialScoreValue">{profileScore.toFixed(1)}%</div>
    </div>
</div>

<style>
    .profileInfoSocial{
        margin-top: 30pt;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 5pt;
    }

    .profileInfoSocialHeader{
        color: var(--text-color-secondary);
        font-family: var(--font-header);
        font-size: var(--font-size-middle);
    }

    .profileInfoSocialName{
        font-family: var(--font-header);
        color: var(--accent-secondary);
        font-size: var(--font-size-big);
    }

    .profileInfoSocialScoreDiv{
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
    }

    .profileInfoSocialScore{
        font-family: var(--font-header);
        color: var(--text-color-secondary);
    }
    .profileInfoSocialScoreValue{
        font-family: var(--font-header);
        font-size: var(--font-size-middle);
        color: var(--accent-primary);
    }

</style>