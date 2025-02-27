<script lang="ts">
    import NavBar from "./NavBar.svelte";
    import Account from "./Account.svelte";
    import { store } from "../../store.svelte";
    import { orient } from "../../types";

    let headerWidth: number = $state(570);
    $effect(() => {
        if (headerWidth < 570) {
            $store.orientLoginBtns = orient.vert
        } else {
            $store.orientLoginBtns = orient.horiz
        }
    })
</script>

<header class="comp-header" bind:clientWidth={headerWidth}>
    <div class="upper-box">
        <a class="h1-anchor" href="/"><h1>Svelte SPA Starter</h1></a>
        {#if $store.showLoginBtns}
            <Account/>
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
