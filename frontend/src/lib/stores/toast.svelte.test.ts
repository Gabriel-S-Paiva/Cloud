import { describe, it, expect, beforeEach } from 'vitest';

import { toast } from './toast.svelte';
import { mockToastItem } from '$lib/mocks/testData';

describe('adding notifications', () => {
	beforeEach(() => {
		toast.toastQueue = [];
	});

	it('info uses the info variant with a 5000ms default duration', () => {
		toast.info('Test Info');
		expect(toast.queue).toContainEqual(
			expect.objectContaining({ message: 'Test Info', variant: 'info', durationMs: 5000 })
		);
	});

	it('success uses the success variant with a 5000ms default duration', () => {
		toast.success('Test Success');
		expect(toast.queue).toContainEqual(
			expect.objectContaining({ message: 'Test Success', variant: 'success', durationMs: 5000 })
		);
	});

	it('warning uses the warning variant with a 6000ms default duration', () => {
		toast.warning('Test warning');
		expect(toast.queue).toContainEqual(
			expect.objectContaining({ message: 'Test warning', variant: 'warning', durationMs: 6000 })
		);
	});

	it('error uses the error variant with a 7000ms default duration', () => {
		toast.error('Test error');
		expect(toast.queue).toContainEqual(
			expect.objectContaining({ message: 'Test error', variant: 'error', durationMs: 7000 })
		);
	});

	it('a custom duration overrides the variant default', () => {
		toast.error('Test duration', 8000);
		expect(toast.queue).toContainEqual(
			expect.objectContaining({ message: 'Test duration', variant: 'error', durationMs: 8000 })
		);
	});
});

describe('delete behaviour', () => {
	beforeEach(() => {
		toast.toastQueue = [];
	});

	it('deleteById removes only the matching toast', () => {
		toast.toastQueue = [
			mockToastItem({ id: 'Test Id' }),
			mockToastItem({ id: 'Np' }),
			mockToastItem({ id: 'Np' }),
			mockToastItem({ id: 'Np' })
		];

		toast.deleteById('Test Id');

		expect(toast.queue).not.toContainEqual(expect.objectContaining({ id: 'Test Id' }));
		expect(toast.queue).toHaveLength(3);
	});

	it('deleteById does nothing when the id is not found', () => {
		const toastMsm = [
			mockToastItem({ id: 'Test Id' }),
			mockToastItem({ id: 'Np' }),
			mockToastItem({ id: 'Np' }),
			mockToastItem({ id: 'Np' })
		];
		toast.toastQueue = toastMsm;

		toast.deleteById('Never seen Id');

		expect(toast.queue).toEqual(toastMsm);
	});

	it('deleteIndex removes only the item at that index', () => {
		const toastMsm = [
			mockToastItem({ id: 'Test Id' }),
			mockToastItem({ id: 'Np1' }),
			mockToastItem({ id: 'Epic Id' }),
			mockToastItem({ id: 'Np2' })
		];
		toast.toastQueue = [...toastMsm];
		const targetItem = toastMsm[2];

		toast.deleteIndex(2);

		expect(toast.queue).toHaveLength(3);
		expect(toast.queue).not.toContainEqual(targetItem);
	});
});
