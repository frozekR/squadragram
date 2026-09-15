export type SkinsRarity = 'BLUE' |'RAGE' | 'INFINITE' | 'LEGENDARY'

export interface Character {
  ID: number;
  UUID: string; 
  Name: string;
  Description: string;
  Rarity: SkinsRarity;
}

export interface CharacterFilter {
  rarity?: SkinsRarity;
  name?: string;
}