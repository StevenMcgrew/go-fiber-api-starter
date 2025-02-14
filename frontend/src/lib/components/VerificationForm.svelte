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
    let response: any = null;

    function resetFormData() {
        formData.userId = S.user?.id;
        formData.verificationCode = "";
    }

    function setUser(res: any) {
        const user: User = {
            token: res.data.token,
            id: res.data.id,
            email: res.data.email,
            username: res.data.username,
            role: res.data.role,
            status: res.data.status,
            imageUrl: res.data.imageUrl,
        };
        S.user = user;
    }

    async function onsubmit(e: SubmitEvent) {
        e.preventDefault();
        isLoading = true;
        error = null;
        try {
            response = await submitVerification(formData);
        } catch (err) {
            error = err;
        } finally {
            isLoading = false;
            if (error === false) {
                resetFormData();
                setUser(response);
                S.showModal = "";
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
