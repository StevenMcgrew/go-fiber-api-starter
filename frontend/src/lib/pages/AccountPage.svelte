<script lang="ts">
    import { store } from "../../store.svelte";
    import { getEmailValidationWarnings, getUsernameValidationWarnings } from "../../validation";
    import { submitForm } from "../../fetch";
    import { modalComp } from "../../types";

    let usernameWarnings = "";
    let emailWarnings = "";
    let passwordWarnings = "";

    async function submit_username(e: SubmitEvent) {
        e.preventDefault();
        usernameWarnings = "";

        const form = e.currentTarget as HTMLFormElement;

        usernameWarnings = getUsernameValidationWarnings(form.username.value);
        if (usernameWarnings) {
            return;
        }

        const formData = new FormData(form);
        const url = `${$store.baseFetchUrl}/users/${$store.user.id}/username`;
        const method = "PATCH";
        const token = $store.user.token;
        let response: any = null;

        try {
            response = await submitForm(formData, url, method, token);
        } catch (err: any) {
            usernameWarnings = err.message;
        } finally {
            if (usernameWarnings === "") {
                form.reset();
                $store.user.username = response.data.username;
                $store.showToast = {
                    color: "green",
                    text: "Username updated!",
                };
            }
        }
    }

    async function submit_email(e: SubmitEvent) {
        e.preventDefault();
        emailWarnings = "";

        const form = e.currentTarget as HTMLFormElement;

        emailWarnings = getEmailValidationWarnings(form.email.value);
        if (emailWarnings) {
            return;
        }

        const formData = new FormData(form);
        const url = `${$store.baseFetchUrl}/users/${$store.user.id}/email`;
        const method = "PATCH";
        const token = $store.user.token;
        let response: any = null;

        try {
            response = await submitForm(formData, url, method, token);
        } catch (err: any) {
            emailWarnings = err.message;
        } finally {
            if (emailWarnings === "") {
                form.reset();
                $store.newEmailAddress = response.data.email;
                // TODO: implement separate verification form
                $store.showModal = modalComp.VerificationForm;
            }
        }
    }
</script>

<div class="user-account-page">
    <!-- Profile Section -->
    <div class="profile-section">
        <div>
            <img
                class="user-img profile-picture"
                src={$store.baseStorageUrl + $store.user.imageUrl}
                alt="user"
            />
            <button class="edit-picture-btn">Change Picture</button>
        </div>
        <h2>{$store.user.username}</h2>
    </div>

    <!-- Update Username Section -->
    <div class="update-section">
        <p>Username: {$store.user.username}</p>
        <form onsubmit={submit_username} class="update-form">
            <label for="username">New Username</label><br />
            <input
                type="text"
                id="username"
                name="username"
                autocomplete="off"
                spellcheck="false"
            /><br />
            <p class="error-text form-input-warning">{usernameWarnings}</p>
            <button type="submit">Update</button>
        </form>
    </div>

    <!-- Update Email Section -->
    <div class="update-section">
        <p>Email: {$store.user.email}</p>
        <form onsubmit={submit_email} class="update-form">
            <label for="email">New Email</label><br />
            <input
                type="email"
                id="email"
                name="email"
                autocomplete="off"
                spellcheck="false"
            /><br />
            <p class="error-text form-input-warning">{emailWarnings}</p>
            <button type="submit">Update</button>
        </form>
    </div>

    <!-- Update Password Section -->
    <div class="update-section">
        <p>Password</p>
        <form class="update-form">
            <label for="current-password">Current Password</label><br />
            <input
                type="password"
                id="current-password"
                name="current-password"
                autocomplete="off"
                spellcheck="false"
            /><br />
            <p class="error-text form-input-warning"></p>

            <label for="new-password">New Password</label><br />
            <input
                type="password"
                id="new-password"
                name="new-password"
                autocomplete="off"
                spellcheck="false"
            /><br />
            <p class="error-text form-input-warning">{passwordWarnings}</p>

            <label for="confirm-password">Confirm New Password</label><br />
            <input
                type="password"
                id="confirm-password"
                name="confirm-password"
                autocomplete="off"
                spellcheck="false"
            /><br />
            <p class="error-text form-input-warning"></p>

            <button type="submit">Update</button>
        </form>
    </div>
</div>

<style>
    .user-account-page {
        max-width: 600px;
        margin: 0 auto;
    }

    .profile-section {
        text-align: center;
        margin-bottom: 30px;
    }

    .profile-section img {
        width: 150px;
        height: 150px;
    }

    .edit-picture-btn {
        display: block;
        margin: 10px auto;
    }

    .update-section {
        margin-bottom: 35px;
    }

    .update-section > p {
        font-weight: bold;
        font-size: 18px;
        margin: 0;
    }

    .update-form input {
        width: 300px;
    }

    .update-form p {
        margin-bottom: 0;
    }

    .update-form label {
        font-size: 14px;
    }
</style>
