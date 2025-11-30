<script lang="ts">
    import { store, clearStorageAndReload } from "../../store.svelte";
    import { modalComp, orient } from "../../types";
    import { clickOutside, getImgSrc } from "../../utils";

    let isMenuOpen = false;

    function toggleMenu() {
        isMenuOpen = !isMenuOpen;
    }

    function logOut() {
        $store.showModal = modalComp.LoggingOutMsg;
        clearStorageAndReload();
    }
</script>

<div
    class="btn-box {$store.orientLoginBtns === orient.vert ? orient.vert : ''}"
>
    {#if $store.user?.token}
        <p class="username-txt">{$store.user.username}</p>
        <!-- <a use:link href="/account">
            <img
                class="user-img"
                src={getImgSrc($store.baseStorageUrl, $store.user.imageUrl)}
                alt="user"
            />
        </a>
        <button onclick={logOut}>Log Out</button> -->
        <!-- Use an on:click outside event handler to close the menu automatically -->
        <div
            class="profile-dropdown-container"
            use:clickOutside={() => (isMenuOpen = false)}
        >
            <!-- Profile Picture Button/Icon -->
            <button
                class="profile-icon"
                id="profileButton"
                aria-label="User Profile Menu"
                aria-haspopup="true"
                aria-expanded={isMenuOpen}
                onclick={toggleMenu}
            >
                <!-- Placeholder for the user's profile image -->
                <img
                    src={getImgSrc($store.baseStorageUrl, $store.user.imageUrl)}
                    alt="User Avatar"
                />

                <!-- Red dot notification indicator for the icon itself (only visible if $hasUnreadNotifications is true) -->
                {#if $store.hasUnreadNotifications}
                    <span class="notification-badge notification-badge-icon"
                    ></span>
                {/if}
            </button>

            <!-- Dropdown Menu -->
            <div
                class="dropdown-menu panel"
                class:visible={isMenuOpen}
                id="dropdownMenu"
            >
                <ul>
                    <li><a onclick={() => (isMenuOpen = false)} href="#/account">Account</a></li>
                    <li>
                        <a onclick={() => (isMenuOpen = false)} href="#/notifications">Notifications</a>
                        <!-- Red dot notification indicator for the menu item -->
                        {#if $store.hasUnreadNotifications}
                            <span class="notification-badge"></span>
                        {/if}
                    </li>
                    <li>
                        <button class="logout-btn" onclick={logOut}
                            >Log Out</button
                        >
                    </li>
                </ul>
            </div>
        </div>
    {:else}
        <button
            class={$store.orientLoginBtns === orient.vert ? "btn-width" : ""}
            onclick={() => ($store.showModal = modalComp.LoginForm)}
        >
            Log In
        </button>
        <button
            class={$store.orientLoginBtns === orient.vert ? "btn-width" : ""}
            onclick={() => ($store.showModal = modalComp.SignUpForm)}
        >
            Sign Up
        </button>
    {/if}
</div>

<style>
    .btn-box {
        display: flex;
        gap: 0.5rem;
        align-items: center;
    }
    .vertical {
        flex-direction: column;
    }
    .btn-width {
        width: 89px;
    }
    .username-txt {
        margin: 0;
    }

    /* Container for positioning */
    .profile-dropdown-container {
        position: relative;
        display: inline-block;
    }

    .profile-icon {
        width: 50px;
        height: 50px;
        border-radius: 50%;
        padding: 0;
        border: 2px solid #ccc;
        cursor: pointer;
        background: none;
        position: relative; /* Needed for positioning the badge */
    }

    .profile-icon img {
        width: 100%;
        height: 100%;
        border-radius: 50%;
        display: block;
    }

    .dropdown-menu {
        position: absolute;
        right: 0;
        min-width: max-content;
        margin-top: 10px;
        display: none;
        opacity: 0;
        transition: opacity 0.2s ease-in-out;
    }

    .dropdown-menu.visible {
        display: block;
        opacity: 1;
    }

    .dropdown-menu ul {
        list-style-type: none;
        padding: 0;
        margin: 0;
    }

    .dropdown-menu li {
        padding: 16px 16px 0px 16px;
        cursor: pointer;
        margin: 0;
    }

    /* Notification Badge (Red Dot) Styling */
    .notification-badge {
        display: inline-block;
        height: 10px;
        width: 10px;
        background-color: #ff4136; /* Bright Red */
        border-radius: 50%;
        border: 1.6px solid white;
    }

    /* Specific positioning for the badge on the main profile icon */
    .profile-icon .notification-badge-icon {
        position: absolute;
        top: 0px;
        right: 0px;
    }

    li a {
        text-decoration: none;
    }

    .logout-btn {
        margin-bottom: 16px;
        margin-top: 4px;
    }
</style>
