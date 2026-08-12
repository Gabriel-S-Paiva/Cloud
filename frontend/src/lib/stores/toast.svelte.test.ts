import { describe, it, expect, vi, beforeEach } from 'vitest';

import { toast } from './toast.svelte';
import type { ToastItem } from '../types';

describe('Adding notification', () => {
	it('info calls the correct variant', () => {
		toast.info('Test Info');
		expect(toast.queue).toContainEqual(
			expect.objectContaining({ message: 'Test Info', variant: 'info', durationMs: 5000 })
		);
	});

	it('sucess calls the correct variant', () => {
		toast.success('Test Sucess');
		expect(toast.queue).toContainEqual(
			expect.objectContaining({ message: 'Test Sucess', variant: 'success', durationMs: 5000 })
		);
	});

	it('warning calls the correct variant', () => {
		toast.warning('Test warning');
		expect(toast.queue).toContainEqual(
			expect.objectContaining({ message: 'Test warning', variant: 'warning', durationMs: 6000 })
		);
	});

	it('error calls the correct variant', () => {
		toast.error('Test error');
		expect(toast.queue).toContainEqual(
			expect.objectContaining({ message: 'Test error', variant: 'error', durationMs: 7000 })
		);
	});

	it('Custom duration overwrites default', () => {
		toast.error('Test duration', 8000);
		expect(toast.queue).toContainEqual(
			expect.objectContaining({ message: 'Test duration', variant: 'error', durationMs: 8000 })
		);
	});
});

describe('Delete Behaviour', () => {
	it('delete by id', () => {
		const toastMsm: ToastItem[] = [
			{ id: 'Test Id', message: 'Message', variant: 'success', durationMs: 6000 },
			{ id: 'Np', message: 'Message', variant: 'success', durationMs: 6000 },
			{ id: 'Np', message: 'Message', variant: 'success', durationMs: 6000 },
			{ id: 'Np', message: 'Message', variant: 'success', durationMs: 6000 }
		];
		toast.toastQueue = toastMsm;

		toast.deleteById('Test Id');

		expect(toast.queue).not.include(expect.objectContaining({ id: 'Test Id' }));
	});

	it('delete by id', () => {
		const toastMsm: ToastItem[] = [
			{ id: 'Test Id', message: 'Message', variant: 'success', durationMs: 6000 },
			{ id: 'Np', message: 'Message', variant: 'success', durationMs: 6000 },
			{ id: 'Np', message: 'Message', variant: 'success', durationMs: 6000 },
			{ id: 'Np', message: 'Message', variant: 'success', durationMs: 6000 }
		];
		toast.toastQueue = toastMsm;

		toast.deleteById('Never seen Id');

		expect(toast.queue).toEqual(toastMsm);
	});

	it('delete by index', () => {
		const toastMsm: ToastItem[] = [
			{ id: 'Test Id', message: 'Message', variant: 'success', durationMs: 6000 },
			{ id: 'Np1', message: 'Message', variant: 'success', durationMs: 6000 },
			{ id: 'Epic Id', message: 'Message', variant: 'success', durationMs: 6000 },
			{ id: 'Np2', message: 'Message', variant: 'success', durationMs: 6000 }
		];

		toast.toastQueue = [...toastMsm];

		const targetItem = toastMsm[2];

		toast.deleteIndex(2);

		expect(toast.queue).toHaveLength(3);
		expect(toast.queue).not.toContainEqual(targetItem);
	});
});
