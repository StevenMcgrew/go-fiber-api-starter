<script lang="ts">
    import NavBar from "./NavBar.svelte";
    import Login from "./Login.svelte";
    import { S, orient } from "../../store.svelte";

    let headerWidth: number = $state(570);
    $effect(() => {
        if (headerWidth < 570) {
            S.orientLoginBtns = orient.vert
        } else {
            S.orientLoginBtns = orient.horiz
        }
    })
</script>

<header class="comp-header" bind:clientWidth={headerWidth}>
    <div class="upper-box">
        <a class="h1-anchor" href="/"><h1>Svelte SPA Starter</h1></a>
        {#if S.showLoginBtns}
            <Login/>
        {/if}
    </div>
    <NavBar/>
</header>

<style>
    .comp-header {
        container: header-cont / inline-size;
        padding-bottom: 1.5rem;
    }
    .upper-box {
        display: flex;
        justify-content: space-between;
        align-items: center;
    }
    .h1-anchor {
        color: var(--nc-tx-1);
        text-decoration: none;
    }
</style>
