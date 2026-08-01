export interface ActionReturn<Parameter> {
	update?: (parameter: Parameter) => void;
	destroy?: () => void;
}

export interface Action<Parameter = void, Return = ActionReturn<Parameter>> {
	<Node extends HTMLElement>(node: Node, parameter: Parameter): Return;
}

export const longpress: Action<number> = (node, duration) => {
	let timer: number | undefined;
	let isLongPress = false;

	// ⭐ Fungsi untuk memulai timer
	function startTimer() {
		isLongPress = false;
		timer = window.setTimeout(() => {
			isLongPress = true;
			node.dispatchEvent(new CustomEvent('longpress'));
		}, duration);
	}

	// ⭐ Fungsi untuk membatalkan timer
	function clearTimer() {
		if (timer) {
			clearTimeout(timer);
			timer = undefined;
		}
	}

	// ⭐ Desktop events
	function handleMouseDown(event: MouseEvent) {
		// ⭐ Hanya untuk left click
		if (event.button === 0) {
			startTimer();
		}
	}

	function handleMouseUp(event: MouseEvent) {
		clearTimer();
		// ⭐ Jika bukan long press, dispatch shortpress
		if (!isLongPress) {
			node.dispatchEvent(new CustomEvent('shortpress'));
		}
	}

	function handleMouseLeave() {
		clearTimer();
	}

	// ⭐ Mobile events
	function handleTouchStart(event: TouchEvent) {
		// ⭐ Cegah default behavior (scroll, zoom)
		event.preventDefault();
		startTimer();
	}

	function handleTouchEnd(event: TouchEvent) {
		clearTimer();
		// ⭐ Jika bukan long press, dispatch shortpress
		if (!isLongPress) {
			node.dispatchEvent(new CustomEvent('shortpress'));
		}
	}

	function handleTouchMove(event: TouchEvent) {
		// ⭐ Jika user swipe, cancel long press
		clearTimer();
	}

	// ⭐ Register event listeners
	node.addEventListener('mousedown', handleMouseDown);
	node.addEventListener('mouseup', handleMouseUp);
	node.addEventListener('mouseleave', handleMouseLeave);

	node.addEventListener('touchstart', handleTouchStart, { passive: false });
	node.addEventListener('touchend', handleTouchEnd);
	node.addEventListener('touchmove', handleTouchMove, { passive: true });

	return {
		update(newDuration: number) {
			clearTimer();
			duration = newDuration;
		},
		destroy() {
			clearTimer();
			node.removeEventListener('mousedown', handleMouseDown);
			node.removeEventListener('mouseup', handleMouseUp);
			node.removeEventListener('mouseleave', handleMouseLeave);
			node.removeEventListener('touchstart', handleTouchStart);
			node.removeEventListener('touchend', handleTouchEnd);
			node.removeEventListener('touchmove', handleTouchMove);
		}
	};
};
