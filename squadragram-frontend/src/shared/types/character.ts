export type CharacterRole = 'DAMAGE' | 'TANK' | 'TECHNIC'

export interface Character {
  ID: number;
  UUID: string; 
  Name: string;
  Description: string;
  Role: CharacterRole;
}

export interface CharacterFilter {
  role?: CharacterRole;
  name?: string;
}