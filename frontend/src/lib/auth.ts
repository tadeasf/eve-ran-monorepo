export function checkAdminAuth(): boolean {
  if (typeof window !== 'undefined') {
    // Client-side check
    return localStorage.getItem('admin_authenticated') === 'true';
  }
  // For server-side, we'll handle this differently in the middleware
  return false;
}

export function setAdminAuth() {
  if (typeof window !== 'undefined') {
    localStorage.setItem('admin_authenticated', 'true');
  }
}

export function clearAdminAuth() {
  if (typeof window !== 'undefined') {
    localStorage.removeItem('admin_authenticated');
  }
}

export function verifyCredentials(username: string, password: string): boolean {
  const adminUsername = process.env.ADMIN_USERNAME || 'admin';
  const adminPassword = process.env.ADMIN_PASSWORD || 'TnDrg2024SecCorp';

  return username === adminUsername && password === adminPassword;
}
