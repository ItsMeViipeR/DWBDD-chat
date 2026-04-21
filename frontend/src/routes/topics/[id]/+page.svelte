<script lang="ts">
	import { page } from '$app/state';
	import { resolve } from '$app/paths';
	import { isTokenExpired } from '$lib/isTokenExpired';
	import { jwtDecode } from 'jwt-decode';
	import type { Token } from '$lib/token';

	const token = localStorage.getItem('chat_token');
	let userId: number | null = $state(null);

	if (token && !isTokenExpired(token)) {
		userId = jwtDecode<Token>(token).user_id;
	} else {
		window.location.href = '/login';
	}

	interface Message {
		id: number;
		user: {
			username: string;
			id: number;
		};
		content: string;
		created_at: string;
	}

	let topicId = $derived(page.params.id);
	let messages: Message[] = $state([]);
	let errorStatus: string | null = $state(null);
	let newMessage: string = $state('');
	let message: string = $state('');
	let isError: boolean = $state(false);

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

	async function sendMessage() {
		if (!newMessage.trim()) return;
		console.log('sending message', newMessage);

		if (!token || isTokenExpired(token)) {
			message = 'Veuillez vous connecter pour envoyer un message.';
			isError = true;
			return;
		}

		try {
			const response = await fetch(`http://localhost:8080/api/messages`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${token}`
				},
				body: JSON.stringify({ topic_id: Number(topicId), content: newMessage })
			});

			if (!response.ok) {
				message = 'Une erreur est survenue.';
				isError = true;
				return;
			}

			await fetchMessages(topicId!);

			newMessage = '';
			scrollTo(0, document.body.scrollHeight);
		} catch (error) {
			message =
				'Une erreur est survenue : ' + (error instanceof Error ? error.message : String(error));
			isError = true;
			return;
		}

		message = '';
	}

	let contextMenu = $state({ show: false, x: 0, y: 0, messageId: null as number | null });

	function handleContextMenu(event: MouseEvent, messageId: number) {
		event.preventDefault();

		contextMenu = {
			show: true,
			x: event.clientX,
			y: event.clientY,
			messageId: messageId
		};
	}

	function closeContextMenu() {
		contextMenu = {
			show: false,
			x: 0,
			y: 0,
			messageId: null
		};
	}

	async function deleteMessage(messageId: number) {
		if (!messageId) return;

		try {
			const response = await fetch(`http://localhost:8080/api/messages/${messageId}`, {
				method: 'DELETE',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${token}`
				}
			});

			if (!response.ok) {
				message = 'Une erreur est survenue.';
				isError = true;
				return;
			}

			await fetchMessages(topicId!);
			closeContextMenu();
		} catch (error) {
			message =
				'Une erreur est survenue : ' + (error instanceof Error ? error.message : String(error));
			isError = true;
		}

		closeContextMenu();
	}
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
				oncontextmenu={(e) => handleContextMenu(e, msg.id)}
				role="button"
				tabindex="0"
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

		<div class="h-36"></div>

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

		<div
			class="fixed right-0 bottom-0 left-0 bg-linear-to-t from-white via-white/95 to-transparent px-4 pt-12 pb-6"
		>
			<div class="mx-auto max-w-2xl">
				<div
					class="relative flex items-center gap-2 rounded-2xl border border-gray-200 bg-white p-2 shadow-[0_10px_40px_rgba(0,0,0,0.04)] transition-all duration-300 focus-within:border-indigo-400 focus-within:shadow-[0_10px_40px_rgba(79,70,229,0.08)]"
				>
					<textarea
						bind:value={newMessage}
						placeholder="Écrire une réponse..."
						rows="1"
						class="flex-1 resize-none border-none bg-transparent px-3 py-2.5 text-[15px] text-gray-800 placeholder-gray-400 focus:ring-0"
						onkeydown={(e) =>
							e.key === 'Enter' && !e.shiftKey && (e.preventDefault(), sendMessage())}
					></textarea>

					<button
						onclick={sendMessage}
						disabled={!newMessage.trim()}
						class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-indigo-600 text-white shadow-sm transition-all hover:bg-indigo-700 active:scale-90 disabled:bg-gray-100 disabled:text-gray-400"
						title="Envoyer"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							viewBox="0 0 24 24"
							fill="currentColor"
							class="h-5 w-5"
						>
							<path
								d="M3.478 2.405a.75.75 0 00-.926.94l2.432 7.905H13.5a.75.75 0 010 1.5H4.984l-2.432 7.905a.75.75 0 00.926.94 60.519 60.519 0 0018.445-8.986.75.75 0 000-1.218A60.517 60.517 0 003.478 2.405z"
							/>
						</svg>
					</button>
				</div>

				<div
					class="mt-3 flex items-center justify-center gap-2 text-[11px] font-medium text-gray-500"
				>
					<kbd
						class="rounded-md border border-gray-300 bg-white px-1.5 py-0.5 font-sans text-[10px] text-gray-600 shadow-sm"
						>Entrée</kbd
					>
					<span class="opacity-70">pour envoyer</span>
				</div>
			</div>
		</div>
	</div>
{/if}

<svelte:window onclick={closeContextMenu} onscroll={closeContextMenu} />

{#if contextMenu.show}
	{@const selectedMsg = messages.find((m) => m.id === contextMenu.messageId)}

	<div
		class="fixed z-50 min-w-37.5 rounded-lg border border-gray-200 bg-white py-1 shadow-xl"
		style="top: {contextMenu.y}px; left: {contextMenu.x}px;"
	>
		{#if selectedMsg && selectedMsg.user.id === userId}
			<button
				onclick={() => deleteMessage(contextMenu.messageId!)}
				class="flex w-full items-center gap-2 px-4 py-2 text-left text-sm text-red-600 hover:bg-red-50"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="h-4 w-4"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
					/>
				</svg>
				Supprimer le message
			</button>
		{:else}
			<div class="px-4 py-2 text-sm text-gray-400 italic">Aucune action disponible</div>
		{/if}
	</div>
{/if}
