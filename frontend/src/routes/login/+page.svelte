<script lang="ts">
  import { auth } from '$lib/stores/auth.svelte';
  import { goto } from '$app/navigation';
  import Input from '$lib/components/UI/Input/Input.svelte';
  import { User, Lock, Eye, EyeOff } from '@lucide/svelte';
	import { toast } from '$lib/stores/toast.svelte';
	import Toast from '$lib/components/UI/Toast/Toast.svelte';

  let mode = $state<'login' | 'register'>('login');

  let loginUsername = $state('');
  let loginPassword = $state('');
  let PasswordVisible = $state(false);
  let loginError = $state<string | null>(null);

  let registerUsername = $state('');
  let registerPassword = $state('');

  async function handleLogin(e: Event) {
    e.preventDefault();
    loginError = null;
    try {
      await auth.login(loginUsername, loginPassword);
      goto('/home');
    } catch (err) {
      err instanceof Error ? toast.add(err.message) : toast.add('Login failed');
    }
  }

  async function handleRegister(e: Event) {
    e.preventDefault();
    try {
      await auth.register(registerUsername, registerPassword);
      toast.add('Account created - waiting for admin approval.')
    } catch (err) {
      err instanceof Error ? toast.add(err.message) : toast.add('Registration failed');
    }
  }
</script>

<div class="min-h-screen flex items-center justify-center bg-surface px-4">
  <div class="w-full max-w-sm">

    <h1 class="font-display text-2xl text-text text-center mb-6">Owned Cloud</h1>

    <div class="flex mb-6 border-b border-border">
      <button
        class="flex-1 pb-3 text-sm font-medium transition-colors {mode === 'login' ? 'text-accent border-b-2 border-accent' : 'text-text-muted'}"
        onclick={() => (mode = 'login')}
      >
        Login
      </button>
      <button
        class="flex-1 pb-3 text-sm font-medium transition-colors {mode === 'register' ? 'text-accent border-b-2 border-accent' : 'text-text-muted'}"
        onclick={() => (mode = 'register')}
      >
        Register
      </button>
    </div>

    {#if mode === 'login'}
      <form onsubmit={handleLogin} class="flex flex-col gap-3">
        <Input bind:value={loginUsername} placeholder="username">
          {#snippet left()}<User size={16} />{/snippet}
        </Input>
        <Input
          bind:value={loginPassword}
          placeholder="password"
          type={PasswordVisible ? 'text' : 'password'}
        >
          {#snippet left()}<Lock size={16} />{/snippet}
          {#snippet right()}
            <button type="button" onclick={() => (PasswordVisible = !PasswordVisible)}>
              {#if PasswordVisible}<EyeOff size={16} />{:else}<Eye size={16} />{/if}
            </button>
          {/snippet}
        </Input>
        <button type="submit" class="h-10 rounded-md bg-accent text-surface text-sm font-medium mt-1">
          Login
        </button>
        {#if loginError}<p class="text-danger text-sm">{loginError}</p>{/if}
      </form>
    {:else}
      <form onsubmit={handleRegister} class="flex flex-col gap-3">
        <Input bind:value={registerUsername} placeholder="username">
          {#snippet left()}<User size={16} />{/snippet}
        </Input>
        <Input
          bind:value={registerPassword}
          placeholder="password"
          type={PasswordVisible ? 'text' : 'password'}
        >
          {#snippet left()}<Lock size={16} />{/snippet}
          {#snippet right()}
            <button type="button" onclick={() => (PasswordVisible = !PasswordVisible)}>
              {#if PasswordVisible}<EyeOff size={16} />{:else}<Eye size={16} />{/if}
            </button>
          {/snippet}
        </Input>
        <button type="submit" class="h-10 rounded-md bg-accent text-surface text-sm font-medium mt-1">
          Register
        </button>
      </form>
    {/if}

  </div>
  <Toast/>
</div>