<script lang="ts">
    import { tblComp } from "../../types";
    import DeleteButton from "./tableComponents/DeleteButton.svelte";
    import NotificationStatus from "./tableComponents/NotificationStatus.svelte";

    type TableData = Record<string, any>;

    interface Props {
        data: TableData[];
        columns: {
            header: string;
            sortable: boolean;
            key: string;
            component?: string;
        }[];
    }

    let { data = [], columns }: Props = $props();

    // State runes
    let currentPage = $state(1);
    let rowsPerPage = $state(10);
    let sortKey = $state<string | null>(null);
    let sortOrder = $state<"asc" | "desc">("asc");

    // Derived sorted data
    const sortedData = $derived.by(() => {
        let result = [...data];
        if (sortKey) {
            result.sort((a, b) => {
                const valA = a[sortKey!];
                const valB = b[sortKey!];
                if (valA < valB) return sortOrder === "asc" ? -1 : 1;
                if (valA > valB) return sortOrder === "asc" ? 1 : -1;
                return 0;
            });
        }
        return result;
    });

    // Derived paginated data
    const paginatedData = $derived(
        sortedData.slice(
            (currentPage - 1) * rowsPerPage,
            currentPage * rowsPerPage,
        ),
    );

    const totalPages = $derived(Math.ceil(data.length / rowsPerPage));

    // Pagination Ellipsis Logic
    const getPageNumbers = $derived(() => {
        const pages = [];
        const maxVisible = 5;

        if (totalPages <= maxVisible) {
            for (let i = 1; i <= totalPages; i++) pages.push(i);
        } else {
            if (currentPage <= 3) {
                pages.push(1, 2, 3, 4, "...", totalPages);
            } else if (currentPage >= totalPages - 2) {
                pages.push(
                    1,
                    "...",
                    totalPages - 3,
                    totalPages - 2,
                    totalPages - 1,
                    totalPages,
                );
            } else {
                pages.push(
                    1,
                    "...",
                    currentPage - 1,
                    currentPage,
                    currentPage + 1,
                    "...",
                    totalPages,
                );
            }
        }
        return pages;
    });

    function toggleSort(key: string) {
        if (sortKey === key) {
            sortOrder = sortOrder === "asc" ? "desc" : "asc";
        } else {
            sortKey = key;
            sortOrder = "asc";
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
            <select bind:value={rowsPerPage} onchange={() => (currentPage = 1)}>
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
                        {#if sortKey === col.key}
                            {sortOrder === "asc" ? "↑" : "↓"}
                        {/if}
                    </th>
                {/each}
            </tr>
        </thead>
        <tbody>
            {#each paginatedData as row}
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
        <button
            class="chevron-btn"
            disabled={currentPage === 1}
            onclick={() => (currentPage = 1)}>&laquo;</button
        >
        <button
            class="chevron-btn"
            disabled={currentPage === 1}
            onclick={() => currentPage--}>&lsaquo;</button
        >

        {#each getPageNumbers() as page}
            {#if page === "..."}
                <span>...</span>
            {:else}
                <button
                    class:active={currentPage === page}
                    onclick={() => (currentPage = Number(page))}
                >
                    {page}
                </button>
            {/if}
        {/each}

        <button
            class="chevron-btn"
            disabled={currentPage === totalPages}
            onclick={() => currentPage++}>&rsaquo;</button
        >
        <button
            class="chevron-btn"
            disabled={currentPage === totalPages}
            onclick={() => (currentPage = totalPages)}>&raquo;</button
        >
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
    .chevron-btn {
        line-height: 0px;
        font-size: 23px;
    }
</style>
