<script lang="ts">
	let name = $state('');
	let email = $state('');
	let password = $state('');
	let confirmPassword = $state('');

	// Fonction pour envoyer les données au backend Go
	async function registerUser(e: Event) {
		e.preventDefault();
		if (password !== confirmPassword) {
			alert('Les mots de passe ne correspondent pas !');
			return;
		}

		try {
			const response = await fetch('http://localhost:8080/api/register', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ name, email, password })
			});

			const result = await response.json();

			if (response.ok) {
				window.location.href = '/login';
			} else {
				alert('Erreur : ' + (result.error || "Échec de l'inscription"));
			}
		} catch (err) {
			console.error('Erreur réseau :', err);
			alert('Impossible de joindre le serveur Go.');
		}
	}
</script>

<div class="flex min-h-screen items-center justify-center bg-gray-100">
	<form class="w-full max-w-md rounded bg-white p-8 shadow-md" onsubmit={registerUser}>
		<h2 class="mb-6 text-center text-2xl font-bold">Sign Up</h2>

		<div class="mb-4">
			<label class="mb-1 block font-medium" for="name">Name</label>
			<input
				id="name"
				type="text"
				class="w-full rounded border px-3 py-2 focus:border-blue-300 focus:ring focus:outline-none"
				bind:value={name}
				required
			/>
		</div>

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

		<div class="mb-6">
			<label class="mb-1 block font-medium" for="confirmPassword">Confirm Password</label>
			<input
				id="confirmPassword"
				type="password"
				class="w-full rounded border px-3 py-2 focus:border-blue-300 focus:ring focus:outline-none"
				bind:value={confirmPassword}
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
