import type { ToastItem, ToastVariant } from "$lib/types";

class Toast {
  toastQueue = $state<ToastItem[]>([]);
  get queue() { return this.toastQueue; }

  private push(message: string, variant: ToastVariant, durationMs: number) {
    const id = crypto.randomUUID();
    this.toastQueue.push({ id, message, durationMs, variant });
  }

  info(message: string, durationMs = 5000)    { this.push(message, 'info', durationMs); }
  success(message: string, durationMs = 5000) { this.push(message, 'success', durationMs); }
  warning(message: string, durationMs = 6000) { this.push(message, 'warning', durationMs); }
  error(message: string, durationMs = 7000)   { this.push(message, 'error', durationMs); }

  deleteById(id: string) { this.toastQueue = this.toastQueue.filter((t) => t.id !== id); }
  deleteIndex(index: number) { this.toastQueue.splice(index, 1); }
}

export const toast = new Toast()