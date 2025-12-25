<script lang="ts">
    interface Props {
        hasViewed: boolean;
    }

    let { hasViewed }: Props = $props();
</script>

{#if hasViewed} 
    <span>Viewed</span>
{:else}
    <span class="notification-badge"></span><span>&nbsp;New</span>
{/if}
