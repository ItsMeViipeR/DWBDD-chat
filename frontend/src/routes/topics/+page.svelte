<script lang="ts">
	import { onMount } from 'svelte';

	interface Topic {
		id: number;
		name: string;
		description: string;
		created_at: string;
	}

	let topics: Topic[] = $state([]);
	let loading = $state(true);

	async function fetchTopics() {
		try {
			const response = await fetch('http://localhost:8080/api/topics');
			const data = await response.json();
			topics = data.topics;
		} catch (err) {
			console.error('Erreur lors du chargement des sujets:', err);
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		fetchTopics();
	});
</script>

<div class="mx-auto max-w-4xl p-6">
	<h1 class="mb-6 text-2xl font-bold">Liste des Sujets</h1>

	{#if loading}
		<p class="text-gray-500">Chargement des sujets...</p>
	{:else}
		<div class="grid gap-4">
			{#each topics as topic (topic.id)}
				<a
					href={`http://localhost:5173/topics/${topic.id}`}
					class="block rounded-xl border border-gray-200 bg-white p-5 shadow-sm transition-all hover:border-indigo-300 hover:shadow-md"
				>
					<div class="flex items-center justify-between">
						<h2 class="text-lg font-semibold text-indigo-600">{topic.name}</h2>
						<span class="text-xs text-gray-400">ID: #{topic.id}</span>
					</div>
					<p class="mt-2 text-gray-600">{topic.description}</p>
				</a>
			{:else}
				<p class="text-gray-500 italic">Aucun sujet disponible pour le moment.</p>
			{/each}
		</div>
	{/if}
</div>
