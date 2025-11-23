<script lang="ts">
	import { page } from '$app/stores';

	interface NavItem {
		label: string;
		href: string;
		icon?: string;
	}

	const navItems: NavItem[] = [
		{ label: 'Home', href: '/', icon: '🏠' },
		{ label: 'Logbook', href: '/logbook', icon: '📖' },
		{ label: 'Propagation', href: '/#propagation', icon: '📡' }
	];

	let mobileMenuOpen = $state(false);

	function toggleMobileMenu() {
		mobileMenuOpen = !mobileMenuOpen;
	}

	function isActive(href: string): boolean {
		if ($page?.url?.pathname === undefined) return false;
		if (href === '/') {
			return $page.url.pathname === '/';
		}
		return $page.url.pathname.startsWith(href);
	}
</script>

<nav class="bg-white shadow-sm">
	<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
		<div class="flex h-16 items-center justify-between">
			<!-- Logo -->
			<div class="flex items-center">
				<a href="/" class="flex items-center">
					<span class="text-2xl font-bold text-primary-600">Ham-Radio Cloud</span>
				</a>
			</div>

			<!-- Desktop Navigation -->
			<div class="hidden md:block">
				<div class="ml-10 flex items-baseline space-x-4">
					{#each navItems as item}
						<a
							href={item.href}
							class="rounded-md px-3 py-2 text-sm font-medium transition-colors {isActive(item.href)
								? 'bg-primary-100 text-primary-700'
								: 'text-gray-700 hover:bg-gray-100 hover:text-gray-900'}"
						>
							{#if item.icon}
								<span class="mr-1">{item.icon}</span>
							{/if}
							{item.label}
						</a>
					{/each}
				</div>
			</div>

			<!-- Auth Buttons (Desktop) -->
			<div class="hidden md:block">
				<div class="flex items-center gap-3">
					<a href="/login" class="text-sm font-medium text-gray-700 hover:text-gray-900">
						Sign In
					</a>
					<a href="/register" class="btn btn-primary btn-sm">Get Started</a>
				</div>
			</div>

			<!-- Mobile menu button -->
			<div class="md:hidden">
				<button
					onclick={toggleMobileMenu}
					type="button"
					class="inline-flex items-center justify-center rounded-md p-2 text-gray-700 hover:bg-gray-100 hover:text-gray-900 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-primary-500"
					aria-controls="mobile-menu"
					aria-expanded={mobileMenuOpen}
				>
					<span class="sr-only">Open main menu</span>
					{#if mobileMenuOpen}
						<!-- X icon -->
						<svg
							class="h-6 w-6"
							fill="none"
							viewBox="0 0 24 24"
							stroke-width="1.5"
							stroke="currentColor"
						>
							<path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
						</svg>
					{:else}
						<!-- Hamburger icon -->
						<svg
							class="h-6 w-6"
							fill="none"
							viewBox="0 0 24 24"
							stroke-width="1.5"
							stroke="currentColor"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"
							/>
						</svg>
					{/if}
				</button>
			</div>
		</div>
	</div>

	<!-- Mobile menu -->
	{#if mobileMenuOpen}
		<div class="md:hidden" id="mobile-menu">
			<div class="space-y-1 px-2 pb-3 pt-2">
				{#each navItems as item}
					<a
						href={item.href}
						class="block rounded-md px-3 py-2 text-base font-medium {isActive(item.href)
							? 'bg-primary-100 text-primary-700'
							: 'text-gray-700 hover:bg-gray-100 hover:text-gray-900'}"
						onclick={toggleMobileMenu}
					>
						{#if item.icon}
							<span class="mr-2">{item.icon}</span>
						{/if}
						{item.label}
					</a>
				{/each}
			</div>
			<div class="border-t border-gray-200 px-2 pb-3 pt-4">
				<div class="flex flex-col gap-2">
					<a
						href="/login"
						class="block rounded-md px-3 py-2 text-base font-medium text-gray-700 hover:bg-gray-100 hover:text-gray-900"
						onclick={toggleMobileMenu}
					>
						Sign In
					</a>
					<a
						href="/register"
						class="block rounded-md bg-primary-600 px-3 py-2 text-center text-base font-medium text-white hover:bg-primary-700"
						onclick={toggleMobileMenu}
					>
						Get Started
					</a>
				</div>
			</div>
		</div>
	{/if}
</nav>
