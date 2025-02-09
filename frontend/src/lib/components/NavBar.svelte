<script lang="ts">
    import { link, location } from "svelte-spa-router";
</script>

<nav>
    <a use:link href="/" class={$location === "/" ? 'active' : ''}>Home</a>
    <a use:link href="/about" class={$location === "/about" ? 'active' : ''}>About</a>
    <a use:link href="/contact" class={$location === "/contact" ? 'active' : ''}>Contact</a>
</nav>

<style>
    nav {
        display: flex;
    }

    a {
        text-decoration: none;
        padding-right: 1rem;
        color: var(--nc-lk-1);
    }

    a.active {
        text-decoration: underline;
        text-decoration-color: rgb(197, 197, 197);
        text-underline-offset: 5px;
        text-decoration-thickness: 3px;
    }
</style>
