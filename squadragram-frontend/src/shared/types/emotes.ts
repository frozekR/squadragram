export type EmotesRarity = 'BLUE' | 'GOLD' | 'LEGENDARY'

export interface Character {
  ID: number;
  UUID: string; 
  Name: string;
  Description: string;
  Rarity: EmotesRarity;
}

export interface CharacterFilter {
  rarity?: EmotesRarity;
  name?: string;
}