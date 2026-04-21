<script lang="ts">
	import { page } from '$app/state';
	import { resolve } from '$app/paths';

	interface Message {
		id: number;
		user: {
			username: string;
		};
		content: string;
		created_at: string;
	}

	let topicId = $derived(page.params.id);
	let messages: Message[] = $state([]);
	let errorStatus: string | null = $state(null);

	async function fetchMessages(id: string) {
		if (!id) return;
		errorStatus = null;

		const response = await fetch(`http://localhost:8080/api/messages?topic_id=${id}`);

		if (response.status === 404) {
			errorStatus = "Ce sujet n'existe pas.";
			messages = [];
			return;
		}

		if (!response.ok) {
			errorStatus = 'Une erreur est survenue.';
			return;
		}

		const data = await response.json();
		messages = data.messages;
	}

	const formatRelativeTime = (dateString: string) => {
		const date = new Date(dateString);
		const now = new Date();
		const diffInSeconds = Math.floor((date.getTime() - now.getTime()) / 1000);

		const units = [
			{ unit: 'year', seconds: 31536000 },
			{ unit: 'month', seconds: 2592000 },
			{ unit: 'day', seconds: 86400 },
			{ unit: 'hour', seconds: 3600 },
			{ unit: 'minute', seconds: 60 },
			{ unit: 'second', seconds: 1 }
		] as const;

		const rtf = new Intl.RelativeTimeFormat('fr', { numeric: 'auto' });

		for (const { unit, seconds } of units) {
			if (Math.abs(diffInSeconds) >= seconds || unit === 'second') {
				return rtf.format(Math.round(diffInSeconds / seconds), unit);
			}
		}
	};

	$effect(() => {
		fetchMessages(topicId!);
	});
</script>

{#if errorStatus}
	<div
		class="mx-auto my-8 flex max-w-md flex-col items-center justify-center rounded-2xl border border-red-200 bg-red-50 p-12 shadow-sm"
	>
		<div class="mb-4 text-red-500">
			<svg
				xmlns="http://www.w3.org/2000/svg"
				class="h-12 w-12"
				fill="none"
				viewBox="0 0 24 24"
				stroke="currentColor"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
				/>
			</svg>
		</div>
		<h2 class="mb-2 text-xl font-bold text-red-800">Oups !</h2>
		<p class="text-center text-red-600">{errorStatus}</p>
		<a
			href={resolve('/topics')}
			class="mt-6 rounded-lg bg-red-600 px-4 py-2 text-white transition-colors hover:bg-red-700"
		>
			Retour aux sujets
		</a>
	</div>
{:else}
	<header class="mx-auto mb-8 max-w-2xl border-b border-gray-100 pb-6">
		<h1 class="text-3xl font-extrabold tracking-tight text-gray-900">
			Sujet <span class="text-indigo-600">#{topicId}</span>
		</h1>
		<p class="mt-1 text-gray-500">Consultez et répondez aux messages de la communauté.</p>
	</header>

	<div class="mx-auto max-w-2xl space-y-4 px-4">
		{#each messages as msg (msg.id)}
			<div
				class="group rounded-2xl border border-gray-200 bg-white p-6 shadow-sm transition-shadow hover:shadow-md"
			>
				<div class="flex items-start gap-4">
					<div
						class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-indigo-100 font-bold text-indigo-700 uppercase"
					>
						{msg.user.username.substring(0, 2)}
					</div>

					<div class="min-w-0 flex-1">
						<div class="mb-1 flex items-center justify-between">
							<span class="truncate font-bold text-gray-900">
								{msg.user.username}
							</span>
							<span class="text-xs text-gray-400 italic">
								{formatRelativeTime(msg.created_at)}
							</span>
						</div>
						<p class="leading-relaxed wrap-break-word text-gray-700">
							{msg.content}
						</p>
					</div>
				</div>
			</div>
		{:else}
			<div class="text-center py-20 bg-gray-50 rounded-3xl border-2 border-dashed border-gray-200">
				<p class="text-gray-400 text-lg">Le silence est d'or... Aucun message ici.</p>
				<button class="mt-4 text-indigo-600 font-semibold hover:underline">
					Soyez le premier à écrire !
				</button>
			</div>
		{/each}
	</div>
{/if}
