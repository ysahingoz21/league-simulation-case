const LOGO_MAP: Record<string, string> = {
  'Türkiye': '/team-logos/turkey.png',
  'ABD': '/team-logos/usa.png',
  'Avustralya': '/team-logos/australia.png',
  'Paraguay': '/team-logos/paraguay.png',
}

export function getTeamLogo(teamName: string): string | null {
  return LOGO_MAP[teamName] ?? null
}
