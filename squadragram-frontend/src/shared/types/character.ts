export type CharacterRole = 'DAMAGE' | 'TANK' | 'TECHNICAL'

export interface Character {
  id: number
  uuid: string
  name: string
  description: string
  role: CharacterRole
}

export interface CharacterFilter {
  role?: CharacterRole;
  name?: string;
}