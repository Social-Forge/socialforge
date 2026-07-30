import { writable } from 'svelte/store';

const toastStore = writable<ToastMessage[]>([]);

export const toast = {
	subscribe: toastStore.subscribe,
	add: (message: string, type: ToastMessage['type'] = 'info', duration = 3000) => {
		const id = Math.random().toString(36).substring(2, 9);

		toastStore.update((all) => [...all, { id, message, type, duration }]);
		if (duration > 0) {
			setTimeout(() => {
				toast.remove(id);
			}, duration);
		}

		return id;
	},
	remove: (id: string) => {
		toastStore.update((all) => all.filter((t) => t.id !== id));
	},
	success: (message: string, duration?: number) => toast.add(message, 'success', duration),
	error: (message: string, duration?: number) => toast.add(message, 'error', duration),
	warning: (message: string, duration?: number) => toast.add(message, 'warning', duration),
	info: (message: string, duration?: number) => toast.add(message, 'info', duration)
};
