<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { localizeHref } from '$lib/paraglide/runtime.js';
	import * as i18n from '$lib/paraglide/messages.js';

	let { data } = $props();
</script>

<div class="flex flex-col gap-6">
	<Card.Root>
		<Card.Header class="text-center">
			<Card.Title class="text-xl">{i18n.auth_verify_email()}</Card.Title>
			<Card.Description>
				{i18n.auth_verify_email_description()}
			</Card.Description>
		</Card.Header>
		<Card.CardFooter>
			<Button type="button" variant="default" href={localizeHref('/signin')}
				>{i18n.auth_button_back_to_login()}</Button
			>
		</Card.CardFooter>
	</Card.Root>
</div>
