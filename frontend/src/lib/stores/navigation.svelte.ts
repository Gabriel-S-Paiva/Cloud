import type { PathSegment } from "$lib/types";

class NavigationClass {
  path = $state<PathSegment[]>([]);

  get currentFolderId(): number | null {
    if (this.path.length === 0) return null;
    return this.path[this.path.length - 1].id;
  }

  get urlPath(): string {
    return this.path.map((segment) => segment.displayName).join('/');
  }

  enter(segment: PathSegment) {
    this.path.push(segment);
  }

  goToDepth(index: number) {
    this.path = this.path.slice(0, index + 1);
  }

  reset() {
    this.path = [];
  }

  setPath(segments: PathSegment[]) {
    this.path = segments;
  }
}

export const navigation = new NavigationClass();