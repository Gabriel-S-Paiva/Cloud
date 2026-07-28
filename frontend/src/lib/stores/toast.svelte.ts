import type { ToastItem } from "$lib/types";

class Toast {
    toastQueue = $state<ToastItem[]>([])

    get queue() {
		return this.toastQueue;
	}

	add(message: string, durationMs = 5000) {
		const id = crypto.randomUUID();
		this.toastQueue.push({ id, message, durationMs });
	}

	deleteById(id: string) {
		this.toastQueue = this.toastQueue.filter((t) => t.id !== id);
	}

	deleteIndex(index: number) {
		this.toastQueue.splice(index, 1);
	}
}

export const toast = new Toast()