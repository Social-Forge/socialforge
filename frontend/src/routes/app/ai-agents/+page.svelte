<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import * as Card from '$lib/components/ui/card';
	import * as Sheet from '$lib/components/ui/sheet';
	import * as Tabs from '$lib/components/ui/tabs';
	import * as NativeSelect from '$lib/components/ui/native-select';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { Input } from '$lib/components/ui/input';
	import { Textarea } from '$lib/components/ui/textarea';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch';
	import { Bot, Plus, Pencil, Trash2, Loader2 } from '@lucide/svelte';
	import type { AIAgent, AIKnowledge, AIPlaybook, AIAsset } from '$lib/server/ai-agent';

	let { data } = $props();
	const agents = $derived((data.agents ?? []) as AIAgent[]);

	let open = $state(false);
	let saving = $state(false);
	let error = $state('');
	let editingId = $state<string | null>(null);

	type Form = {
		name: string;
		provider: string;
		model: string;
		system_prompt: string;
		temperature: number;
		max_tokens: number;
		auto_reply_enabled: boolean;
		is_active: boolean;
		p_name: string;
		p_tone: string;
		p_gender: string;
		p_greeting: string;
		p_soul: string;
		g_rules: string;
		s_topics: string;
		s_handoff: boolean;
	};
	const empty: Form = {
		name: '', provider: 'claude', model: '', system_prompt: '',
		temperature: 0.7, max_tokens: 1024, auto_reply_enabled: true, is_active: true,
		p_name: '', p_tone: '', p_gender: '', p_greeting: '', p_soul: '',
		g_rules: '', s_topics: '', s_handoff: true
	};
	let form = $state<Form>({ ...empty });

	function lines(s: string): string[] {
		return s.split('\n').map((x) => x.trim()).filter(Boolean);
	}

	function openCreate() {
		editingId = null;
		form = { ...empty };
		error = '';
		open = true;
	}

	function openEdit(a: AIAgent) {
		editingId = a.id;
		error = '';
		form = {
			name: a.name,
			provider: a.provider,
			model: a.model ?? '',
			system_prompt: a.system_prompt ?? '',
			temperature: a.temperature ?? 0.7,
			max_tokens: a.max_tokens ?? 1024,
			auto_reply_enabled: a.auto_reply_enabled,
			is_active: a.is_active,
			p_name: a.persona?.agent_name ?? '',
			p_tone: a.persona?.tone ?? '',
			p_gender: a.persona?.gender ?? '',
			p_greeting: a.persona?.greeting ?? '',
			p_soul: a.persona?.soul ?? '',
			g_rules: (a.guardrails?.rules ?? []).join('\n'),
			s_topics: (a.safety?.sensitive_topics ?? []).join('\n'),
			s_handoff: a.safety?.handoff_to_human ?? true
		};
		open = true;
		loadResources(a.id);
	}

	function buildPayload() {
		return {
			name: form.name,
			provider: form.provider,
			model: form.model || undefined,
			system_prompt: form.system_prompt,
			temperature: Number(form.temperature) || undefined,
			max_tokens: Number(form.max_tokens) || undefined,
			auto_reply_enabled: form.auto_reply_enabled,
			is_active: form.is_active,
			persona: {
				agent_name: form.p_name,
				tone: form.p_tone,
				gender: form.p_gender,
				greeting: form.p_greeting,
				soul: form.p_soul
			},
			guardrails: { rules: lines(form.g_rules) },
			safety: { handoff_to_human: form.s_handoff, sensitive_topics: lines(form.s_topics) }
		};
	}

	async function save() {
		if (!form.name.trim() || !form.system_prompt.trim()) {
			error = 'Nama dan system prompt wajib diisi';
			return;
		}
		saving = true;
		error = '';
		try {
			const url = editingId ? `/api/ai-agents/${editingId}` : '/api/ai-agents';
			const res = await fetch(url, {
				method: editingId ? 'PUT' : 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(buildPayload())
			});
			const j = await res.json();
			if (!j.success) {
				error = j.message || 'Gagal menyimpan';
				return;
			}
			open = false;
			await invalidateAll();
		} finally {
			saving = false;
		}
	}

	async function removeAgent(a: AIAgent) {
		if (!confirm(`Hapus agent "${a.name}"?`)) return;
		await fetch(`/api/ai-agents/${a.id}`, { method: 'DELETE' });
		await invalidateAll();
	}

	// ---- nested resources ----
	let knowledge = $state<AIKnowledge[]>([]);
	let playbooks = $state<AIPlaybook[]>([]);
	let assets = $state<AIAsset[]>([]);

	async function loadResources(id: string) {
		const [k, p, a] = await Promise.all([
			fetch(`/api/ai-agents/${id}/knowledge`).then((r) => r.json()),
			fetch(`/api/ai-agents/${id}/playbooks`).then((r) => r.json()),
			fetch(`/api/ai-agents/${id}/assets`).then((r) => r.json())
		]);
		knowledge = k.data ?? [];
		playbooks = p.data ?? [];
		assets = a.data ?? [];
	}
	async function reloadResources() {
		if (editingId) await loadResources(editingId);
	}

	let kTitle = $state('');
	let kContent = $state('');
	async function addKnowledge() {
		if (!editingId || !kTitle.trim() || !kContent.trim()) return;
		await fetch(`/api/ai-agents/${editingId}/knowledge`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ title: kTitle, content: kContent })
		});
		kTitle = '';
		kContent = '';
		await reloadResources();
	}
	async function delKnowledge(id: string) {
		if (!editingId) return;
		await fetch(`/api/ai-agents/${editingId}/knowledge/${id}`, { method: 'DELETE' });
		await reloadResources();
	}

	let pbName = $state('');
	let pbKeywords = $state('');
	let pbInstruction = $state('');
	async function addPlaybook() {
		if (!editingId || !pbName.trim() || !pbInstruction.trim()) return;
		await fetch(`/api/ai-agents/${editingId}/playbooks`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({
				name: pbName,
				keywords: pbKeywords.split(',').map((x) => x.trim()).filter(Boolean),
				instruction: pbInstruction
			})
		});
		pbName = '';
		pbKeywords = '';
		pbInstruction = '';
		await reloadResources();
	}
	async function delPlaybook(id: string) {
		if (!editingId) return;
		await fetch(`/api/ai-agents/${editingId}/playbooks/${id}`, { method: 'DELETE' });
		await reloadResources();
	}

	let asName = $state('');
	let asType = $state('image');
	let asKey = $state('');
	let asDesc = $state('');
	async function addAsset() {
		if (!editingId || !asName.trim() || !asKey.trim()) return;
		await fetch(`/api/ai-agents/${editingId}/assets`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ name: asName, type: asType, storage_key: asKey, description: asDesc })
		});
		asName = '';
		asKey = '';
		asDesc = '';
		await reloadResources();
	}
	async function delAsset(id: string) {
		if (!editingId) return;
		await fetch(`/api/ai-agents/${editingId}/assets/${id}`, { method: 'DELETE' });
		await reloadResources();
	}

	function nstr(v: unknown): string {
		if (v && typeof v === 'object' && 'String' in v)
			return (v as any).Valid ? (v as any).String : '';
		return (v as string) ?? '';
	}
</script>

<div class="mx-auto flex w-full max-w-6xl flex-col gap-6 p-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-xl font-semibold">AI Agents</h1>
			<p class="text-sm text-muted-foreground">
				Kelola agen AI, persona, knowledge, playbook & aset.
			</p>
		</div>
		<Button onclick={openCreate}><Plus class="mr-2 h-4 w-4" /> Agent Baru</Button>
	</div>

	{#if agents.length === 0}
		<Card.Root>
			<Card.Content class="flex flex-col items-center gap-3 py-16 text-center">
				<Bot class="h-10 w-10 text-muted-foreground" />
				<p class="text-muted-foreground">Belum ada AI agent. Buat yang pertama.</p>
				<Button onclick={openCreate}><Plus class="mr-2 h-4 w-4" /> Agent Baru</Button>
			</Card.Content>
		</Card.Root>
	{:else}
		<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
			{#each agents as a (a.id)}
				<Card.Root class="flex flex-col">
					<Card.Header>
						<div class="flex items-start justify-between gap-2">
							<div class="flex items-center gap-3">
								<div class="flex h-10 w-10 items-center justify-center rounded-full bg-primary/10">
									<Bot class="h-5 w-5 text-primary" />
								</div>
								<div>
									<Card.Title class="text-base">{a.name}</Card.Title>
									<div class="text-xs text-muted-foreground">{a.provider} · {a.model}</div>
								</div>
							</div>
							<Badge variant={a.is_active ? 'default' : 'secondary'}>
								{a.is_active ? 'Aktif' : 'Nonaktif'}
							</Badge>
						</div>
					</Card.Header>
					<Card.Content class="flex-1">
						<p class="line-clamp-3 text-sm text-muted-foreground">{a.system_prompt}</p>
						<div class="mt-3 flex flex-wrap gap-1">
							{#if a.persona?.agent_name}
								<Badge variant="outline">{a.persona.agent_name}</Badge>
							{/if}
							{#if a.auto_reply_enabled}<Badge variant="outline">Auto-reply</Badge>{/if}
						</div>
					</Card.Content>
					<Card.Footer class="gap-2">
						<Button variant="outline" size="sm" onclick={() => openEdit(a)}>
							<Pencil class="mr-1 h-3.5 w-3.5" /> Edit
						</Button>
						<Button variant="ghost" size="sm" class="text-red-600" onclick={() => removeAgent(a)}>
							<Trash2 class="mr-1 h-3.5 w-3.5" /> Hapus
						</Button>
					</Card.Footer>
				</Card.Root>
			{/each}
		</div>
	{/if}
</div>

<Sheet.Root bind:open>
	<Sheet.Content side="right" class="flex w-full flex-col gap-0 overflow-y-auto sm:max-w-xl">
		<Sheet.Header>
			<Sheet.Title>{editingId ? 'Edit Agent' : 'Agent Baru'}</Sheet.Title>
		</Sheet.Header>

		<Tabs.Root value="general" class="flex-1 px-4 pb-4">
			<Tabs.List class="flex w-full flex-wrap">
				<Tabs.Trigger value="general">Umum</Tabs.Trigger>
				<Tabs.Trigger value="persona">Persona</Tabs.Trigger>
				<Tabs.Trigger value="guardrails">Guardrails</Tabs.Trigger>
				{#if editingId}
					<Tabs.Trigger value="knowledge">Knowledge</Tabs.Trigger>
					<Tabs.Trigger value="playbook">Playbook</Tabs.Trigger>
					<Tabs.Trigger value="asset">Aset</Tabs.Trigger>
				{/if}
			</Tabs.List>

			<Tabs.Content value="general" class="flex flex-col gap-3 pt-4">
				<div class="flex flex-col gap-1.5">
					<Label>Nama</Label>
					<Input bind:value={form.name} placeholder="mis. Sari CS AI" />
				</div>
				<div class="grid grid-cols-2 gap-3">
					<div class="flex flex-col gap-1.5">
						<Label>Provider</Label>
						<NativeSelect.Root bind:value={form.provider}>
							<NativeSelect.Option value="claude">Claude</NativeSelect.Option>
							<NativeSelect.Option value="openai">OpenAI</NativeSelect.Option>
							<NativeSelect.Option value="google">Google</NativeSelect.Option>
						</NativeSelect.Root>
					</div>
					<div class="flex flex-col gap-1.5">
						<Label>Model (opsional)</Label>
						<Input bind:value={form.model} placeholder="default per provider" />
					</div>
				</div>
				<div class="flex flex-col gap-1.5">
					<Label>System Prompt</Label>
					<Textarea bind:value={form.system_prompt} rows={5} placeholder="Instruksi utama agen…" />
				</div>
				<div class="grid grid-cols-2 gap-3">
					<div class="flex flex-col gap-1.5">
						<Label>Temperature</Label>
						<Input type="number" step="0.1" min="0" max="2" bind:value={form.temperature} />
					</div>
					<div class="flex flex-col gap-1.5">
						<Label>Max Tokens</Label>
						<Input type="number" bind:value={form.max_tokens} />
					</div>
				</div>
				<div class="flex items-center justify-between rounded-md border p-3">
					<Label>Auto-reply</Label>
					<Switch bind:checked={form.auto_reply_enabled} />
				</div>
				<div class="flex items-center justify-between rounded-md border p-3">
					<Label>Aktif</Label>
					<Switch bind:checked={form.is_active} />
				</div>
			</Tabs.Content>

			<Tabs.Content value="persona" class="flex flex-col gap-3 pt-4">
				<div class="flex flex-col gap-1.5"><Label>Nama Persona</Label><Input bind:value={form.p_name} placeholder="mis. Sari" /></div>
				<div class="flex flex-col gap-1.5"><Label>Tone</Label><Input bind:value={form.p_tone} placeholder="ramah & ringkas" /></div>
				<div class="flex flex-col gap-1.5"><Label>Gender</Label><Input bind:value={form.p_gender} placeholder="perempuan" /></div>
				<div class="flex flex-col gap-1.5"><Label>Salam Pembuka</Label><Input bind:value={form.p_greeting} placeholder="Halo Kak" /></div>
				<div class="flex flex-col gap-1.5"><Label>Soul / Karakter</Label><Textarea bind:value={form.p_soul} rows={2} placeholder="Antusias, fokus closing" /></div>
			</Tabs.Content>

			<Tabs.Content value="guardrails" class="flex flex-col gap-3 pt-4">
				<div class="flex flex-col gap-1.5">
					<Label>Guardrails (satu aturan per baris)</Label>
					<Textarea bind:value={form.g_rules} rows={4} placeholder={'Jangan janji diskon tanpa konfirmasi\nJangan bahas SARA'} />
				</div>
				<div class="flex flex-col gap-1.5">
					<Label>Topik Sensitif → handoff (satu per baris)</Label>
					<Textarea bind:value={form.s_topics} rows={3} placeholder={'komplain hukum\nrefund besar'} />
				</div>
				<div class="flex items-center justify-between rounded-md border p-3">
					<Label>Handoff ke manusia</Label>
					<Switch bind:checked={form.s_handoff} />
				</div>
			</Tabs.Content>

			{#if editingId}
				<Tabs.Content value="knowledge" class="flex flex-col gap-3 pt-4">
					<div class="flex flex-col gap-2 rounded-md border p-3">
						<Input bind:value={kTitle} placeholder="Judul knowledge" />
						<Textarea bind:value={kContent} rows={3} placeholder="Isi pengetahuan…" />
						<Button size="sm" onclick={addKnowledge}>Tambah Knowledge</Button>
					</div>
					{#each knowledge as k (k.id)}
						<div class="flex items-start justify-between gap-2 rounded-md border p-3">
							<div>
								<div class="text-sm font-medium">{k.title}</div>
								<div class="line-clamp-2 text-xs text-muted-foreground">{k.content}</div>
							</div>
							<Button variant="ghost" size="icon" class="text-red-600" onclick={() => delKnowledge(k.id)}><Trash2 class="h-4 w-4" /></Button>
						</div>
					{/each}
				</Tabs.Content>

				<Tabs.Content value="playbook" class="flex flex-col gap-3 pt-4">
					<div class="flex flex-col gap-2 rounded-md border p-3">
						<Input bind:value={pbName} placeholder="Nama playbook" />
						<Input bind:value={pbKeywords} placeholder="keyword dipisah koma (harga, promo)" />
						<Textarea bind:value={pbInstruction} rows={2} placeholder="Instruksi saat keyword cocok…" />
						<Button size="sm" onclick={addPlaybook}>Tambah Playbook</Button>
					</div>
					{#each playbooks as p (p.id)}
						<div class="flex items-start justify-between gap-2 rounded-md border p-3">
							<div>
								<div class="text-sm font-medium">{p.name}</div>
								<div class="mt-1 flex flex-wrap gap-1">{#each p.keywords as kw}<Badge variant="outline">{kw}</Badge>{/each}</div>
								<div class="mt-1 line-clamp-2 text-xs text-muted-foreground">{p.instruction}</div>
							</div>
							<Button variant="ghost" size="icon" class="text-red-600" onclick={() => delPlaybook(p.id)}><Trash2 class="h-4 w-4" /></Button>
						</div>
					{/each}
				</Tabs.Content>

				<Tabs.Content value="asset" class="flex flex-col gap-3 pt-4">
					<div class="flex flex-col gap-2 rounded-md border p-3">
						<Input bind:value={asName} placeholder="Nama aset" />
						<NativeSelect.Root bind:value={asType}>
							<NativeSelect.Option value="image">Image</NativeSelect.Option>
							<NativeSelect.Option value="video">Video</NativeSelect.Option>
							<NativeSelect.Option value="document">Document</NativeSelect.Option>
						</NativeSelect.Root>
						<Input bind:value={asKey} placeholder="storage_key (MinIO object key)" />
						<Input bind:value={asDesc} placeholder="deskripsi (opsional)" />
						<Button size="sm" onclick={addAsset}>Tambah Aset</Button>
					</div>
					{#each assets as as (as.id)}
						<div class="flex items-start justify-between gap-2 rounded-md border p-3">
							<div>
								<div class="text-sm font-medium">{as.name} <Badge variant="outline">{as.type}</Badge></div>
								<div class="text-xs text-muted-foreground">{nstr(as.description)}</div>
							</div>
							<Button variant="ghost" size="icon" class="text-red-600" onclick={() => delAsset(as.id)}><Trash2 class="h-4 w-4" /></Button>
						</div>
					{/each}
				</Tabs.Content>
			{/if}
		</Tabs.Root>

		{#if error}
			<div class="mx-4 rounded-md bg-red-500/10 px-3 py-2 text-sm text-red-600">{error}</div>
		{/if}
		<Sheet.Footer class="flex-row gap-2 border-t p-4">
			<Button variant="outline" class="flex-1" onclick={() => (open = false)}>Batal</Button>
			<Button class="flex-1" onclick={save} disabled={saving}>
				{#if saving}<Loader2 class="mr-2 h-4 w-4 animate-spin" />{/if}
				{editingId ? 'Simpan' : 'Buat Agent'}
			</Button>
		</Sheet.Footer>
	</Sheet.Content>
</Sheet.Root>
