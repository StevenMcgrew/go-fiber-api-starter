<script lang="ts">
    import { submitVerification } from "../../fetch";
    import { S } from "../../store.svelte";
    import type { User } from "../../types";

    let formData = {
        userId: S.user?.id,
        verificationCode: "",
    };

    let isLoading = false;
    let error = null;
    let data: any = null;

    function resetFormData() {
        formData.userId = S.user?.id;
        formData.verificationCode = "";
    }

    function setUser(data: any) {
        const user: User = {
            token: data.token,
            id: data.id,
            email: data.email,
            username: data.username,
            role: data.role,
            status: data.status,
            imageUrl: data.imageUrl,
        };
        S.user = user;
    }

    async function onsubmit(e: SubmitEvent) {
        e.preventDefault();
        isLoading = true;
        error = null;
        try {
            data = await submitVerification(formData);
        } catch (err) {
            error = err;
        } finally {
            isLoading = false;
            if (error === false) {
                resetFormData();
                setUser(data);
                S.showModal = null;
            }
        }
    }
</script>

<div class="form-wrapper">
    <form {onsubmit} class="auth-form">
        <h3>Verification</h3>
        <p>
            A verification code was sent to your email address. Please check
            your email and enter the code here:
        </p>
        <label for="verificationCode">Verification Code</label>
        <input
            bind:value={formData.verificationCode}
            id="verificationCode"
            name="verificationCode"
            type="text"
            required
            maxlength="6"
        />
        <div class="form-btn-box">
            <button type="submit">Verify</button>
        </div>
    </form>
</div>
