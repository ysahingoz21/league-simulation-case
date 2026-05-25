const LOGO_MAP: Record<string, string> = {
  Turkey: "/team-logos/turkey.png",
  USA: "/team-logos/usa.png",
  Australia: "/team-logos/australia.png",
  Paraguay: "/team-logos/paraguay.png",
};

export function getTeamLogo(teamName: string): string | null {
  return LOGO_MAP[teamName] ?? null;
}
