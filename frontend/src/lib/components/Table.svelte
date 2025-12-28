<script lang="ts">
    import type { Column, Pagination, Sorting, TableData } from "../../types";
    import { tblComp } from "../../types";
    import DeleteButton from "./tableComponents/DeleteButton.svelte";
    import NotificationStatus from "./tableComponents/NotificationStatus.svelte";

    interface Props {
        data: TableData[];
        columns: Column[];
        pagination: Pagination;
        sorting: Sorting;
    }

    let { data = [], columns, pagination, sorting }: Props = $props();

    // Pagination Ellipsis Logic
    const getPageNumbers = $derived(() => {
        const pages = [];
        const maxVisible = 5;

        if (pagination.totalPages <= maxVisible) {
            for (let i = 1; i <= pagination.totalPages; i++) pages.push(i);
        } else {
            if (pagination.page <= 3) {
                pages.push(1, 2, 3, 4, "...", pagination.totalPages);
            } else if (pagination.page >= pagination.totalPages - 2) {
                pages.push(
                    1,
                    "...",
                    pagination.totalPages - 3,
                    pagination.totalPages - 2,
                    pagination.totalPages - 1,
                    pagination.totalPages,
                );
            } else {
                pages.push(
                    1,
                    "...",
                    pagination.page - 1,
                    pagination.page,
                    pagination.page + 1,
                    "...",
                    pagination.totalPages,
                );
            }
        }
        return pages;
    });

    function toggleSort(key: string) {
        if (sorting.key === key) {
            sorting.direction = sorting.direction === "asc" ? "desc" : "asc";
        } else {
            sorting.key = key;
            sorting.direction = "asc";
        }
    }

    function handleEdit(id: number) {
        data = data.filter((r) => r.id !== id);
    }
</script>

<div class="table-container">
    <!-- Rows Per Page Select -->
    <div class="controls">
        <label>
            Rows per page:
            <select bind:value={pagination.perPage} onchange={() => (pagination.page = 1)}>
                <option value={10}>10</option>
                <option value={50}>50</option>
                <option value={100}>100</option>
            </select>
        </label>
    </div>

    <table>
        <thead>
            <tr>
                {#each columns as col}
                    <th
                        onclick={() => col.sortable && toggleSort(col.key)}
                        class:sortable={col.sortable}
                    >
                        {col.header}
                        {#if sorting.key === col.key}
                            {sorting.direction === "asc" ? "↑" : "↓"}
                        {/if}
                    </th>
                {/each}
            </tr>
        </thead>
        <tbody>
            {#each data as row}
                <tr>
                    {#each columns as col}
                        <td>
                            {#if col.component}
                                {#if col.component == tblComp.NotificationStatus}
                                    <NotificationStatus
                                        hasViewed={row[col.key]}
                                    />
                                {:else if col.component == tblComp.DeleteButton}
                                    <DeleteButton recordId={row[col.key]} onEdit={handleEdit}/>
                                    <!-- {:else if $store.showModal == modalComp.SignUpForm}
                                    <SignUpForm />
                                {:else if $store.showModal == modalComp.SignupVerificationForm}
                                    <SignupVerificationForm />
                                {:else if $store.showModal == modalComp.UpdateEmailVerificationForm}
                                    <UpdateEmailVerificationForm />
                                {:else if $store.showModal == modalComp.ForgotPasswordForm}
                                    <ForgotPasswordForm />
                                {:else if $store.showModal == modalComp.ResetPasswordForm}
                                    <ResetPasswordForm /> -->
                                {/if}
                            {:else}
                                {row[col.key]}
                            {/if}
                        </td>
                    {/each}
                </tr>
            {/each}
        </tbody>
    </table>

    <!-- Pagination Controls -->
    <div class="pagination">
        {#each getPageNumbers() as page}
            {#if page === "..."}
                <span>...</span>
            {:else}
                <button
                    class:active={pagination.page === page}
                    onclick={() => (pagination.page = Number(page))}
                >
                    {page}
                </button>
            {/if}
        {/each}
    </div>
</div>

<style>
    th.sortable {
        cursor: pointer;
    }
    .active {
        font-weight: bold;
        border: 2px solid white;
    }
    .pagination {
        margin-top: 1rem;
        display: flex;
        gap: 0.5rem;
        justify-content: center;
    }
    .pagination button {
        width: 36px;
    }
</style>
