<script lang="ts">
	// Define the possible views as a type
	type View = 'Users' | 'Notifications';

	// Svelte 5 state rune for navigation
	let currentView = $state<View>('Users');

	function setView(view: View) {
		currentView = view;
	}
</script>

<div class="dashboard">
	<h3>Admin Dashboard</h3>
	<div class="layout">
		
		<nav class="sidebar">
			<button class="transparent-btn"
				class:active={currentView === 'Users'} 
				onclick={() => setView('Users')}
			>
				Users
			</button>
			<button class="transparent-btn"
				class:active={currentView === 'Notifications'} 
				onclick={() => setView('Notifications')}
			>
				Notifications
			</button>
		</nav>

		<main>
			<h4>{currentView}</h4>
			
			{#if currentView === 'Users'}
				<p>User management interface goes here.</p>
			{:else if currentView === 'Notifications'}
				<p>Notification settings and logs go here.</p>
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
