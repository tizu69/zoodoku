<script>
	import { dailyResult, level } from '#lib/config.svelte.ts';
	import { loadDailyResult } from '#lib/daily.ts';
	import { browser } from '$app/env';
	import { goto } from '$app/navigation';
	import ManageSaveOverlay from './ManageSaveOverlay.svelte';
	import API from '#api';

	let daily = $derived(loadDailyResult());
	let lastDaily = $derived(dailyResult.current);
	let managementVisible = $state(false);
</script>

<ManageSaveOverlay bind:managementVisible />

<main class="mt-1/3 mx-auto flex min-h-dvh max-w-lg flex-col items-center p-8">
	<img src="/icons/icon.svg" alt="logo" class="size-2/3" />

	<button
		onclick={async () => {
			const r = await API.greet({ name: prompt('whats ur name') ?? undefined });
			alert(r.greeting);
		}}
	>
		hi
	</button>

	{#if browser}
		<button
			class="mt-4 w-full clicky rounded-lg border border-depress bg-layer p-4 drop-shadow"
			onclick={() => goto('/daily')}
		>
			<p>Play Daily</p>
			{#if daily}
				<p class="text-xs">{daily.streak} day streak</p>
			{:else if lastDaily?.result === 'win'}
				<p class="text-xs">Don't lose your {lastDaily.streak} day streak!</p>
			{/if}
		</button>
		<button
			class="mt-4 w-full clicky rounded-lg border border-depress bg-layer p-4 drop-shadow"
			onclick={() => goto(`/L${level.current}`)}
		>
			Play: Level {level.current}
		</button>
		<div class="mt-6 grid w-full grid-cols-5 gap-1.5">
			{#each { length: level.current - 1 } as _, i}
				<button
					class="clicky rounded-lg border border-depress bg-layer py-4 text-center drop-shadow"
					onclick={() => goto(`/L${level.current - i - 1}`)}
				>
					{level.current - i - 1}
				</button>
			{/each}
		</div>
		<button class="mt-4 w-full py-2 underline" onclick={() => goto('/tutorial')}> Tutorial </button>
		<button class="w-full py-2 underline" onclick={() => (managementVisible = !managementVisible)}>
			Manage Progress
		</button>
	{:else}
		<p class="mt-4">Loading...</p>
	{/if}
	<i class="mt-4 text-sm">
		Copyright © 2026 Ruby H.
		<a href="https://g.tizu.dev/zoodoku" class="underline">Source</a>
	</i>
</main>
