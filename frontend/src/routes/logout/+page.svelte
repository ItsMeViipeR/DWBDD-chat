<script lang="ts">
	let message = $state('Vous êtes déconnecté.');
	let isError = $state(false);

	if (!localStorage.getItem('chat_token')) {
		message = "Vous n'êtes pas connecté.";
		isError = true;
	} else {
		try {
			localStorage.removeItem('chat_user');
			localStorage.removeItem('chat_token');
		} catch (e) {
			console.error(e);
			isError = true;
			message = 'Impossible de vous déconnecter.';
		}
	}
</script>

{#if isError}
	<div class="mb-4 rounded bg-red-100 p-2 text-sm text-red-600">
		{message}
	</div>
{:else}
	<div class="mb-4 rounded bg-green-100 p-2 text-sm text-green-600">
		{message}
	</div>
{/if}
