import { describe, it, expect, beforeEach } from 'vitest';

import { navigation } from './navigation.svelte';
import { mockPathSegment } from '$lib/mocks/testData';
import type { PathSegment } from '$lib/types';

const nestedPath: PathSegment[] = [
	{ id: 1, displayName: 'Level1' },
	{ id: 2, displayName: 'Level2' },
	{ id: 3, displayName: 'Level3' },
	{ id: 4, displayName: 'Level4' }
];

describe('initial store state', () => {
	beforeEach(() => {
		navigation.reset();
	});

	it('path starts as an empty array', () => {
		expect(navigation.path).toEqual([]);
	});

	it('currentFolderId is null', () => {
		expect(navigation.currentFolderId).toBeNull();
	});

	it('urlPath is an empty string', () => {
		expect(navigation.urlPath).toBe('');
	});
});

describe('navigation flow', () => {
	beforeEach(() => {
		navigation.reset();
	});

	it('pushing a single segment', () => {
		navigation.enter(mockPathSegment({ id: 1, displayName: 'Level1' }));

		expect(navigation.path.length).toBe(1);
		expect(navigation.currentFolderId).toBe(1);
		expect(navigation.urlPath).toBe('Level1');
	});

	it('nested navigation', () => {
		navigation.setPath(nestedPath);

		expect(navigation.path.length).toBe(nestedPath.length);
		expect(navigation.currentFolderId).toBe(nestedPath.at(-1)!.id);
		expect(navigation.urlPath).toBe('Level1/Level2/Level3/Level4');
	});
});

describe('navigation using breadcrumbs', () => {
	beforeEach(() => {
		navigation.reset();
		navigation.setPath(nestedPath);
	});

	it('trims the path to the selected depth', () => {
		navigation.goToDepth(1);

		expect(navigation.path.length).toBe(2);
		expect(navigation.currentFolderId).toBe(nestedPath.at(1)!.id);
		expect(navigation.urlPath).toBe('Level1/Level2');
	});

	it('selecting root (depth 0) keeps only the first segment', () => {
		navigation.goToDepth(0);

		expect(navigation.path.length).toBe(1);
		expect(navigation.currentFolderId).toBe(nestedPath.at(0)!.id);
		expect(navigation.urlPath).toBe('Level1');
	});

	it('selecting the current depth keeps the full path', () => {
		navigation.goToDepth(nestedPath.length - 1);

		expect(navigation.path.length).toBe(nestedPath.length);
		expect(navigation.currentFolderId).toBe(nestedPath.at(-1)!.id);
		expect(navigation.urlPath).toBe('Level1/Level2/Level3/Level4');
	});

	it('an out-of-range depth clears the path entirely', () => {
		navigation.goToDepth(-1);

		expect(navigation.path).toEqual([]);
		expect(navigation.currentFolderId).toBeNull();
		expect(navigation.urlPath).toBe('');
	});
});

describe('direct operations', () => {
	beforeEach(() => {
		navigation.reset();
	});

	it('reset clears the path', () => {
		navigation.setPath(nestedPath);

		navigation.reset();

		expect(navigation.path).toEqual([]);
		expect(navigation.currentFolderId).toBeNull();
		expect(navigation.urlPath).toBe('');
	});

	it('setPath overwrites the existing path rather than merging', () => {
		const initialPath: PathSegment[] = [
			{ id: 1, displayName: 'Level1' },
			{ id: 2, displayName: 'Level2' }
		];
		const newPath: PathSegment[] = [
			{ id: 3, displayName: 'Level3' },
			{ id: 4, displayName: 'Level4' }
		];
		navigation.setPath(initialPath);

		navigation.setPath(newPath);

		expect(navigation.path).toEqual(newPath);
		expect(navigation.currentFolderId).toBe(4);
		expect(navigation.urlPath).toBe('Level3/Level4');
	});
});
