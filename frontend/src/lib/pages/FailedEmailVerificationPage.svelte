<script lang="ts">
    interface Props {
        failureMessage?: string;
    }

    let { failureMessage = "" }: Props = $props();
</script>

<h2>&#x26A0; Invalid or Expired Verification Link</h2>
<p>
    Sorry, the link to verify your email address has either expired or is not
    valid.
</p>
<p>
    Error message: {failureMessage}
</p>

<style>
    h2 {
        color: goldenrod;
    }
</style>