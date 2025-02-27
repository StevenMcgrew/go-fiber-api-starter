<script lang="ts">
    import { store } from "../../store.svelte";
    import {
        getEmailValidationWarnings,
        getPasswordValidationWarnings,
        getUsernameValidationWarnings,
    } from "../../validation";
    import { submitForm } from "../../fetch";
    import { modalComp, toastColor } from "../../types";

    let usernameWarnings = "";
    let emailWarnings = "";
    let currentPwWarnings = "";
    let newPwWarnings = "";
    let repeatPwWarnings = "";
    let pwSubmissionWarnings = "";

    let isUsernameSubmitting = false;
    let isEmailSubmitting = false;
    let isPasswordSubmitting = false;

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
            isUsernameSubmitting = true;
            response = await submitForm(formData, url, method, token);
        } catch (err: any) {
            usernameWarnings = err.message;
        } finally {
            isUsernameSubmitting = false;
            if (usernameWarnings === "") {
                form.reset();
                $store.user.username = response.data.username;
                $store.showToast = {
                    color: toastColor.green,
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
        const url = `${$store.baseFetchUrl}/users/${$store.user.id}/change-email/request`;
        const method = "POST";
        const token = $store.user.token;
        let response: any = null;

        try {
            isEmailSubmitting = true;
            response = await submitForm(formData, url, method, token);
        } catch (err: any) {
            emailWarnings = err.message;
        } finally {
            isEmailSubmitting = false;
            if (emailWarnings === "") {
                $store.newEmailAddress = response.data;
                $store.showModal = modalComp.UpdateEmailVerificationForm;
            }
        }
    }

    async function submit_password(e: SubmitEvent) {
        e.preventDefault();
        currentPwWarnings = "";
        newPwWarnings = "";
        repeatPwWarnings = "";

        const form = e.currentTarget as HTMLFormElement;

        currentPwWarnings = getPasswordValidationWarnings(
            form.currentPassword.value,
        );
        newPwWarnings = getPasswordValidationWarnings(form.newPassword.value);
        if (form.newPassword.value !== form.repeatNewPassword.value) {
            repeatPwWarnings = "Must match New Password";
        }
        if (currentPwWarnings || newPwWarnings || repeatPwWarnings) {
            return;
        }

        const formData = new FormData(form);
        const url = `${$store.baseFetchUrl}/users/${$store.user.id}/password`;
        const method = "PATCH";
        const token = $store.user.token;
        let response: any = null;

        try {
            isPasswordSubmitting = true;
            response = await submitForm(formData, url, method, token);
        } catch (err: any) {
            pwSubmissionWarnings = err.message;
        } finally {
            isPasswordSubmitting = false;
            if (pwSubmissionWarnings === "") {
                form.reset();
                $store.showToast = {
                    color: toastColor.green,
                    text: "Password updated!",
                };
            }
        }
    }
</script>

{#if $store.user.token}
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
                <p
                    class="form-input-warning {usernameWarnings
                        ? 'error-text'
                        : ''}"
                >
                    {(() => {
                        if (usernameWarnings) {
                            return usernameWarnings;
                        }
                        if (isUsernameSubmitting) {
                            return "Submitting...";
                        }
                    })()}
                </p>
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
                    bind:value={$store.newEmailAddress}
                /><br />
                <p
                    class="form-input-warning {emailWarnings
                        ? 'error-text'
                        : ''}"
                >
                    {(() => {
                        if (emailWarnings) {
                            return emailWarnings;
                        }
                        if (isEmailSubmitting) {
                            return "Submitting...";
                        }
                    })()}
                </p>
                <button type="submit">Update</button>
            </form>
        </div>

        <!-- Update Password Section -->
        <div class="update-section">
            <p>Password</p>
            <form onsubmit={submit_password} class="update-form">
                <label for="current-password">Current Password</label><br />
                <input
                    type="password"
                    id="current-password"
                    name="currentPassword"
                    autocomplete="off"
                    spellcheck="false"
                /><br />
                <p class="error-text form-input-warning">{currentPwWarnings}</p>

                <label for="new-password">New Password</label><br />
                <input
                    type="password"
                    id="new-password"
                    name="newPassword"
                    autocomplete="off"
                    spellcheck="false"
                /><br />
                <p class="error-text form-input-warning">{newPwWarnings}</p>

                <label for="confirm-password">Confirm New Password</label><br />
                <input
                    type="password"
                    id="confirm-password"
                    name="repeatNewPassword"
                    autocomplete="off"
                    spellcheck="false"
                /><br />
                <p class="error-text form-input-warning">{repeatPwWarnings}</p>

                <button type="submit">Update</button>
                <p
                    class="form-input-warning {pwSubmissionWarnings
                        ? 'error-text'
                        : ''}"
                >
                    {(() => {
                        if (pwSubmissionWarnings) {
                            return pwSubmissionWarnings;
                        }
                        if (isPasswordSubmitting) {
                            return "Submitting...";
                        }
                    })()}
                </p>
            </form>
        </div>
    </div>
{:else}
    <p class="error-text">You must be logged in to view this page.</p>
{/if}

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
