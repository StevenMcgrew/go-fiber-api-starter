<script lang="ts">
    import { customFetch } from "../../fetch";
    import {
        type User,
        type Notification,
        type TableData,
        type Column,
        type Pagination,
        type Sorting,
        tblComp,
        type AdminView,
    } from "../../types";
    import Table from "../components/Table.svelte";
    import { onMount } from "svelte";
    import { store } from "../../store.svelte";

    let tableData = $state<TableData[]>([]);
    let columns = $state<Column[]>([]);
    let pagination = $state<Pagination>({
        page: 0,
        perPage: 0,
        totalPages: 0,
        totalCount: 0,
    });
    let sorting = $state<Sorting>({
        key: "",
        direction: "asc",
    });

    onMount(async () => {
        setView($store.adminPageView);
    });

    function setView(view: AdminView) {
        $store.adminPageView = view;
        if ($store.adminPageView === "Users") {
            fetchUsers();
        } else if ($store.adminPageView === "Notifications") {
            fetchNotifications();
        }
    }

    function fetchUsers() {
        // TODO
    }

    async function fetchNotifications() {
        const url = `${$store.baseFetchUrl}/notifications/?page=3&per_page=10`;
        const token = $store.user.token;

        try {
            const response = await customFetch(url, token);
            const rawData: Notification[] = response.data;
            // Transform data
            tableData = rawData.map((n) => ({
                ...n,
                createdAt: new Intl.DateTimeFormat(navigator.language, {
                    dateStyle: "short",
                    timeStyle: "short",
                }).format(new Date(n.createdAt)),
            }));
            columns = [
                {
                    header: "Status",
                    sortable: true,
                    key: "hasViewed",
                    component: tblComp.NotificationStatus,
                },
                { header: "Timestamp", sortable: true, key: "createdAt" },
                { header: "Text", sortable: false, key: "textContent" },
                {
                    header: "",
                    sortable: false,
                    key: "id",
                    component: tblComp.DeleteButton,
                },
            ];
            pagination = {
                page: response.pagination.page,
                perPage: response.pagination.perPage,
                totalPages: response.pagination.totalPages,
                totalCount: response.pagination.totalCount,
            };
            sorting = {
                key: null,
                direction: "asc",
            }
        } catch (err: any) {
            console.log(err.message);
            // error =
            // err instanceof Error
            //     ? err.message
            //     : "An unknown error occurred";
        } finally {
            console.log("finally_finished");
            // isLoading = false;
        }
    }
</script>

<div class="dashboard">
    <h3>Admin Dashboard</h3>
    <div class="layout">
        <nav class="sidebar">
            <button
                class="transparent-btn"
                class:active={$store.adminPageView === "Users"}
                onclick={() => setView("Users")}
            >
                Users
            </button>
            <button
                class="transparent-btn"
                class:active={$store.adminPageView === "Notifications"}
                onclick={() => setView("Notifications")}
            >
                Notifications
            </button>
        </nav>

        <main>
            <h4>{$store.adminPageView}</h4>

            {#if $store.adminPageView === "Users"}
                <Table data={tableData} {columns} {pagination} {sorting}
                ></Table>
            {:else if $store.adminPageView === "Notifications"}
                <Table data={tableData} {columns} {pagination} {sorting}
                ></Table>
            {/if}
        </main>
    </div>
</div>

<style>
    .layout {
        padding-top: 1rem;
        display: flex;
        min-height: 100vh;
    }

    .sidebar {
        display: flex;
        flex-direction: column;
        border-right: 1px solid #ccc;
    }

    .sidebar button {
        text-align: left;
    }

    .sidebar button.active {
        text-decoration: underline;
        text-decoration-color: white;
        text-underline-offset: 3px;
    }

    main {
        flex: 1;
        padding: 1rem;
    }
</style>
