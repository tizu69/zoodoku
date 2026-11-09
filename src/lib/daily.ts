import { dailyResult, type DailyResult } from './config.svelte.ts';

export const DAILY_SIZE = 8;
export type { DailyResult } from './config.svelte.ts';

function dateKey(date = new Date()): string {
	return [date.getFullYear(), date.getMonth() + 1, date.getDate()]
		.map((part, index) => (index === 0 ? String(part) : String(part).padStart(2, '0')))
		.join('-');
}

function dateNumber(key: string): number {
	const [year, month, day] = key.split('-').map(Number);
	return Date.UTC(year, month - 1, day) / 86_400_000;
}

export function dateDiff(date2: string, date1 = dateKey()) {
	return dateNumber(date1) - dateNumber(date2);
}

export function dailyDay(date = new Date()): number {
	const day = Date.UTC(date.getFullYear(), date.getMonth(), date.getDate());
	const start = Date.UTC(2026, 0, 1);
	return Math.floor((day - start) / 86_400_000);
}

export function dailySeed(date = new Date()): number {
	let value = (dailyDay(date) + 1) >>> 0;
	value = Math.imul(value ^ (value >>> 16), 0x45d9f3b);
	value = Math.imul(value ^ (value >>> 16), 0x45d9f3b);
	return (value ^ (value >>> 16)) >>> 0;
}

export function loadDailyResult(date = new Date()): DailyResult | undefined {
	const saved = dailyResult.current;
	return saved?.date === dateKey(date) ? saved : undefined;
}

export function saveDailyResult(
	result: Omit<NonNullable<DailyResult>, 'date' | 'streak'>,
	date = new Date()
): DailyResult {
	const today = dateKey(date);
	const current = loadDailyResult(date);
	const previous = dailyResult.current;
	const continuesStreak = previous?.result === 'win' && dateDiff(previous.date, today) === 1;
	const streak =
		result.result === 'win'
			? current?.result === 'win'
				? current.streak
				: continuesStreak
					? previous.streak + 1
					: 1
			: 0;
	const complete: DailyResult = { ...result, date: today, streak };

	dailyResult.current = complete;
	return complete;
}

export function formatDailyTime(seconds: number): string {
	const minutes = Math.floor(seconds / 60);
	return `${minutes}:${String(seconds % 60).padStart(2, '0')}`;
}
