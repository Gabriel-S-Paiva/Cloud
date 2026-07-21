<script lang="ts">
  import { auth } from '$lib/stores/auth.svelte';
  import { goto } from '$app/navigation';

  let loginUsername = $state('');
  let loginPassword = $state('');
  let loginError = $state<string | null>(null);

  let registerUsername = $state('');
  let registerPassword = $state('');
  let registerError = $state<string | null>(null);
  let registerSuccess = $state(false);

  async function handleLogin(e: Event) {
    e.preventDefault();
    loginError = null;
    try {
      await auth.login(loginUsername, loginPassword);
      goto('/home');
    } catch (err) {
      loginError = err instanceof Error ? err.message : 'Login failed';
    }
  }

  async function handleRegister(e: Event) {
    e.preventDefault();
    registerError = null;
    try {
      await auth.register(registerUsername, registerPassword);
      registerSuccess = true;
    } catch (err) {
      registerError = err instanceof Error ? err.message : 'Registration failed';
    }
  }
</script>

<form onsubmit={handleLogin}>
  <input type="text" bind:value={loginUsername} placeholder="username" />
  <input type="password" bind:value={loginPassword} placeholder="password" />
  <button type="submit">Login</button>
  {#if loginError}<p>{loginError}</p>{/if}
</form>

<form onsubmit={handleRegister}>
  <input type="text" bind:value={registerUsername} placeholder="username" />
  <input type="password" bind:value={registerPassword} placeholder="password" />
  <button type="submit">Register</button>
  {#if registerError}<p>{registerError}</p>{/if}
  {#if registerSuccess}<p>Account created — waiting for admin approval.</p>{/if}
</form>