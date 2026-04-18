<script lang="ts">
	import { isTokenExpired } from '$lib/isTokenExpired';

	let email = $state('');
	let password = $state('');
	let errorMessage = $state('');
	let token = localStorage.getItem('chat_token');

	if (token && !isTokenExpired(token)) {
		window.location.href = '/';
	}

	async function loginUser(e: Event) {
		e.preventDefault();
		errorMessage = '';

		try {
			const response = await fetch('http://localhost:8080/api/login', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ email, password })
			});

			const data = await response.json();

			if (response.ok) {
				localStorage.setItem('chat_token', data.token);
				localStorage.setItem('chat_user', data.user);

				window.location.href = '/';
			} else {
				errorMessage = data.error || 'Identifiants invalides';
			}
		} catch (e) {
			console.log(e);
			errorMessage = 'Le serveur ne répond pas';
		}
	}
</script>

<div class="flex min-h-screen items-center justify-center bg-gray-100">
	<form class="w-full max-w-md rounded bg-white p-8 shadow-md" onsubmit={loginUser}>
		<h2 class="mb-6 text-center text-2xl font-bold">Login</h2>

		{#if errorMessage}
			<div class="mb-4 rounded bg-red-100 p-2 text-sm text-red-600">
				{errorMessage}
			</div>
		{/if}

		<div class="mb-4">
			<label class="mb-1 block font-medium" for="email">Email</label>
			<input
				id="email"
				type="email"
				class="w-full rounded border px-3 py-2 focus:border-blue-300 focus:ring focus:outline-none"
				bind:value={email}
				required
			/>
		</div>

		<div class="mb-4">
			<label class="mb-1 block font-medium" for="password">Password</label>
			<input
				id="password"
				type="password"
				class="w-full rounded border px-3 py-2 focus:border-blue-300 focus:ring focus:outline-none"
				bind:value={password}
				required
			/>
		</div>

		<button
			type="submit"
			class="w-full rounded bg-blue-600 py-2 text-white transition hover:bg-blue-700"
		>
			Sign Up
		</button>
	</form>
</div>
