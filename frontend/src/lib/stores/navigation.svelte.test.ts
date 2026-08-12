import { describe, it, expect, vi, beforeEach } from 'vitest';

import { navigation } from './navigation.svelte';
import type { PathSegment } from '$lib/types';

describe('inital storage state', () => {
	it('path starts as an empty array', async () => {
		expect(navigation.path).toEqual([]);
	});

	it('current Folder Should return null', async () => {
		expect(navigation.currentFolderId).toBeNull();
	});

	it('url path should return ""', async () => {
		expect(navigation.urlPath).toBe('');
	});
});

describe('navigation flow', () => {
	beforeEach(() => {
		navigation.reset();
	});
	it('pushing a single Segment', async () => {
		navigation.enter({ id: 1, displayName: 'Level1' });

		expect(navigation.path.length).toBe(1);
		expect(navigation.currentFolderId).toBe(1);
		expect(navigation.urlPath).toBe('Level1');
	});
	it('nested navigation', async () => {
		const nestedPath: PathSegment[] = [
			{ id: 1, displayName: 'Level1' },
			{ id: 2, displayName: 'Level2' },
			{ id: 3, displayName: 'Level3' },
			{ id: 4, displayName: 'Level4' }
		];
		navigation.setPath(nestedPath);

		expect(navigation.path.length).toBe(nestedPath.length);
		expect(navigation.currentFolderId).toBe(nestedPath.at(-1)!.id);
		expect(navigation.urlPath).toBe('Level1/Level2/Level3/Level4');
	});
});

describe('navigation using breadcrumbs', () => {
	beforeEach(() => {
		navigation.reset();
	});
	it('Triming Path', async () => {
		const nestedPath: PathSegment[] = [
			{ id: 1, displayName: 'Level1' },
			{ id: 2, displayName: 'Level2' },
			{ id: 3, displayName: 'Level3' },
			{ id: 4, displayName: 'Level4' }
		];
		navigation.setPath(nestedPath);

		navigation.goToDepth(1);

		expect(navigation.path.length).toBe(2);
		expect(navigation.currentFolderId).toBe(nestedPath.at(1)!.id);
		expect(navigation.urlPath).toBe('Level1/Level2');
	});

	it('Seleting root', async () => {
		const nestedPath: PathSegment[] = [
			{ id: 1, displayName: 'Level1' },
			{ id: 2, displayName: 'Level2' },
			{ id: 3, displayName: 'Level3' },
			{ id: 4, displayName: 'Level4' }
		];
		navigation.setPath(nestedPath);

		navigation.goToDepth(0);

		expect(navigation.path.length).toBe(1);
		expect(navigation.currentFolderId).toBe(nestedPath.at(0)!.id);
		expect(navigation.urlPath).toBe('Level1');
	});

	it('Selecting current depth', async () => {
		const nestedPath: PathSegment[] = [
			{ id: 1, displayName: 'Level1' },
			{ id: 2, displayName: 'Level2' },
			{ id: 3, displayName: 'Level3' },
			{ id: 4, displayName: 'Level4' }
		];
		navigation.setPath(nestedPath);

		navigation.goToDepth(nestedPath.length - 1);

		expect(navigation.path.length).toBe(nestedPath.length);
		expect(navigation.currentFolderId).toBe(nestedPath.at(nestedPath.length - 1)!.id);
		expect(navigation.urlPath).toBe('Level1/Level2/Level3/Level4');
	});

	it('Out of index', async () => {
		const nestedPath: PathSegment[] = [
			{ id: 1, displayName: 'Level1' },
			{ id: 2, displayName: 'Level2' },
			{ id: 3, displayName: 'Level3' },
			{ id: 4, displayName: 'Level4' }
		];
		navigation.setPath(nestedPath);

		navigation.goToDepth(-1);

		expect(navigation.path.length).toBe(0);
		expect(navigation.path).toEqual([]);
		expect(navigation.currentFolderId).toBe(null);
		expect(navigation.urlPath).toBe('');
	});
});

describe('Direct Operations', () => {
	beforeEach(() => {
		navigation.reset();
	});
	it('reset clears path', async () => {
		const nestedPath: PathSegment[] = [
			{ id: 1, displayName: 'Level1' },
			{ id: 2, displayName: 'Level2' },
			{ id: 3, displayName: 'Level3' },
			{ id: 4, displayName: 'Level4' }
		];
		navigation.setPath(nestedPath);

		navigation.reset();

		expect(navigation.path.length).toBe(0);
		expect(navigation.path).toEqual([]);
		expect(navigation.currentFolderId).toBe(null);
		expect(navigation.urlPath).toBe('');
	});

	it('setpath overwrites path', async () => {
		const nestedPath: PathSegment[] = [
			{ id: 1, displayName: 'Level1' },
			{ id: 2, displayName: 'Level2' }
		];
		const newPath: PathSegment[] = [
			{ id: 3, displayName: 'Level3' },
			{ id: 4, displayName: 'Level4' }
		];
		navigation.setPath(nestedPath);

		navigation.setPath(newPath);

		expect(navigation.path).toEqual(newPath);
		expect(navigation.currentFolderId).toBe(4);
		expect(navigation.urlPath).toBe('Level3/Level4');
	});
});
