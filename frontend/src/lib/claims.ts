function key(userId: string): string {
  return `claims:${userId}`;
}

export function getClaimIds(userId: string): string[] {
  const raw = localStorage.getItem(key(userId));
  if (!raw) return [];
  try {
    const parsed = JSON.parse(raw);
    return Array.isArray(parsed) ? (parsed as string[]) : [];
  } catch {
    return [];
  }
}

export function addClaimId(userId: string, claimId: string): void {
  const ids = getClaimIds(userId);
  if (ids.includes(claimId)) return;
  localStorage.setItem(key(userId), JSON.stringify([claimId, ...ids]));
}

export function removeClaimId(userId: string, claimId: string): void {
  const ids = getClaimIds(userId).filter((id) => id !== claimId);
  localStorage.setItem(key(userId), JSON.stringify(ids));
}
