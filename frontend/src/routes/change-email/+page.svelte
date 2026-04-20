<script lang="ts">
	import { jwtDecode } from 'jwt-decode';
	import { isTokenExpired } from '$lib/isTokenExpired';

	const chat_token = localStorage.getItem('chat_token');

	if (chat_token == null || isTokenExpired(chat_token)) {
		window.location.href = '/login';
	}

	interface Token {
		user_id: number;
		username: string;
		email: string;
		exp: number;
	}

	const decoded = jwtDecode<Token>(chat_token!);
	const email = decoded.email;

	let new_email = $state(email);
	let message = $state('');
	let isError = $state(false);

	async function changeEmail(e: Event) {
		e.preventDefault();
		message = '';

		try {
			const response = await fetch('http://localhost:8080/api/change_email', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${chat_token}`
				},
				body: JSON.stringify({ email: new_email })
			});

			const data = await response.json();

			if (response.ok) {
				localStorage.setItem('chat_token', data.token);
				message = 'Email modifié avec succès';
			} else {
				message = data.error || "Erreur lors de la modification de l'email";
				isError = true;
			}
		} catch (e) {
			console.log(e);
			message = 'Le serveur ne répond pas';
			isError = true;
		}
	}
</script>

<div class="flex min-h-screen items-center justify-center bg-gray-100">
	<form class="w-full max-w-md rounded bg-white p-8 shadow-md" onsubmit={changeEmail}>
		<h2 class="mb-6 text-center text-2xl font-bold">Login</h2>

		{#if message != ''}
			{#if isError}
				<div class="mb-4 rounded bg-red-100 p-2 text-sm text-red-600">
					{message}
				</div>
			{:else}
				<div class="mb-4 rounded bg-green-100 p-2 text-sm text-green-600">
					{message}
				</div>
			{/if}
		{/if}

		<div class="mb-4">
			<label class="mb-1 block font-medium" for="email">Email</label>
			<input
				id="email"
				type="email"
				class="w-full rounded border px-3 py-2 focus:border-blue-300 focus:ring focus:outline-none"
				bind:value={new_email}
				required
			/>
		</div>

		<button
			type="submit"
			class="w-full rounded bg-blue-600 py-2 text-white transition hover:bg-blue-700"
		>
			Sign In
		</button>
	</form>
</div>
