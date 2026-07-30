import Root from '$lib/components/extras/rename/rename.svelte';
import Provider from '$lib/components/extras/rename/rename-provider.svelte';
import Cancel from '$lib/components/extras/rename/rename-cancel.svelte';
import Save from '$lib/components/extras/rename/rename-save.svelte';
import Edit from '$lib/components/extras/rename/rename-edit.svelte';

/**
 * This can be consumed one of two ways:
 * ```svelte
 * <script lang="ts">
 *  import { Rename } from '$lib/components/ui/rename';
 * </script>
 *
 * <Rename ... />
 * ```
 *
 * or
 *
 * ```svelte
 * <script lang="ts">
 *  import * as Rename from '$lib/components/ui/rename';
 * </script>
 *
 * <Rename.Provider>
 *  <Rename.Root ... />
 *  <Rename.Cancel/>
 *  <Rename.Save/>
 *  <Rename.Edit/>
 * </Rename.Provider>
 * ```
 */

export { Root as Rename, Provider, Root, Cancel, Save, Edit };
