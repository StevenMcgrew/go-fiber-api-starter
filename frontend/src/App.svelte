<script lang="ts">
    import Router from "svelte-spa-router";
    import { routes } from "./routes";
    import Header from "./lib/components/Header.svelte";
    import ModalDialog from "./lib/components/ModalDialog.svelte";
</script>

<Header />
<main>
    <Router {routes} />
</main>
<ModalDialog/>
