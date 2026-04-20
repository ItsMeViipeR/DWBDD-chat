<script lang="ts">
	import { isTokenExpired } from '$lib/isTokenExpired';

	let name = $state('');
	let description = $state('');
	let message = $state('');
	let isError = $state(false);
	let token = localStorage.getItem('chat_token');

	if (!token || isTokenExpired(token)) {
		window.location.href = '/login';
	}

	async function createTopic(e: Event) {
		e.preventDefault();
		message = '';
		isError = false;

		try {
			const response = await fetch('http://localhost:8080/api/topics', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${token}`
				},
				body: JSON.stringify({ name, description: description !== '' ? description : '' })
			});

			const data = await response.json();

			if (response.ok) {
				message = 'Topic created successfully!';
				name = '';
				description = '';
			} else {
				message = data.error || 'Something went wrong';
				isError = true;
			}
		} catch (error) {
			message = error as string;
			isError = true;
		}
	}
</script>

<div class="flex min-h-screen items-center justify-center bg-gray-100">
	<form class="w-full max-w-md rounded bg-white p-8 shadow-md" onsubmit={createTopic}>
		<h2 class="mb-6 text-center text-2xl font-bold">Create Topic</h2>

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
			<label class="mb-1 block font-medium" for="description">Description</label>
			<input
				id="description"
				type="text"
				class="w-full rounded border px-3 py-2 focus:border-blue-300 focus:ring focus:outline-none"
				bind:value={description}
			/>
		</div>

		<button
			type="submit"
			class="w-full rounded bg-blue-600 py-2 text-white transition hover:bg-blue-700"
		>
			Create Topic
		</button>
	</form>
</div>
