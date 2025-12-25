<script lang="ts">
    import { store } from "../../store.svelte";
    import { onMount } from "svelte";
    import DataTable from "../components/DataTable.svelte";
    import { tblComp, type Notification } from "../../types";
    import { customFetch } from "../../fetch";
    import Header from "../components/Header.svelte";

    // Reactive state for the API lifecycle
    let notifications = $state<Notification[]>([]);
    let isLoading = $state(true);
    let error = $state<string | null>(null);

    // Column definitions
    const columns = [
        { header: "Status", sortable: true, key: "hasViewed", component: tblComp.NotificationStatus },
        { header: "Timestamp", sortable: true, key: "createdAt" },
        { header: "Text", sortable: false, key: "textContent" },
        { header: "", sortable: false, key: "id", component: tblComp.DeleteButton }
    ];

    onMount(async () => {

        const url = `${$store.baseFetchUrl}/users/${$store.user.id}/notifications`;
        const token = $store.user.token;

        try {
            const response = await customFetch(url, token);
            const rawData: Notification[] = response.data;
            // Transform data
            notifications = rawData.map((n) => ({
                ...n,
                createdAt: new Intl.DateTimeFormat(navigator.language, {
                    dateStyle: "short",
                    timeStyle: "short",
                }).format(new Date(n.createdAt)),
            }));
        } catch (err: any) {
            error =
                err instanceof Error
                    ? err.message
                    : "An unknown error occurred";
        } finally {
            isLoading = false;
        }
    });
</script>

<main>
    <h3>Notifications</h3>
    {#if isLoading}
        <p>Loading...</p>
    {:else if error}
        <p class="error-text">Error: {error}</p>
    {:else}
        <DataTable data={notifications} {columns} />
    {/if}
</main>
