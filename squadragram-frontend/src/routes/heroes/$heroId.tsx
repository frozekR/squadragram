import { useState } from 'react';
import { createFileRoute, useParams } from '@tanstack/react-router';

type HeroData = {
  name?: string;
  bio?: string;
  skins?: unknown[];
  emotes?: unknown[];
  superAttacks?: unknown[];
};

type TabProps = {
  data?: HeroData;
};

const previewHero: HeroData = {
  name: 'Preview',
  bio: 'This local preview keeps the hero layout available while backend data is being prepared.',
  skins: ['Default skin', 'Battle skin'],
  emotes: ['Victory', 'Laugh'],
  superAttacks: ['Meteor Strike'],
};

export const Route = createFileRoute('/heroes/$heroId')({
  component: HeroPage,
});

// Заглушки для контента вкладок
const DetailsTab = ({ data }: TabProps) => <div className="text-gray-700">Детали: {data?.bio || 'Нет данных'}</div>;
const SkinsTab = ({ data }: TabProps) => <div className="text-gray-700">Скины: {data?.skins?.length || 0} шт.</div>;
const EmotesTab = ({ data }: TabProps) => <div className="text-gray-700">Эмоции: {data?.emotes?.length || 0} шт.</div>;
const SuperAttacksTab = ({ data }: TabProps) => <div className="text-gray-700">Супер атаки: {data?.superAttacks?.length || 0} шт.</div>;

export default function HeroPage() {
  const [activeTab, setActiveTab] = useState('details');
  const hero = { ...previewHero, name: `${previewHero.name}` };

  const tabs = [
    { id: 'details', label: 'Details' },
    { id: 'skins', label: 'Skins' },
    { id: 'emotes', label: 'Emotes' },
    { id: 'super-attacks', label: '(MAX) Super Attacks' },
  ];

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col items-center pt-10 px-4">
      
      {/* Заголовок страницы */}
      <h1 className="text-2xl font-bold text-gray-900 mb-8 uppercase tracking-wide">
        Hero {hero?.name} Details
      </h1>

      {/* Навигация по вкладкам */}
      <div className="flex flex-wrap justify-center gap-4 mb-10">
        {tabs.map((tab) => {
          const isActive = activeTab === tab.id;
          
          return (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`
                px-8 py-2.5 rounded-full font-semibold text-black transition-all duration-300 ease-in-out
                ${isActive 
                  ? 'bg-blue-400 shadow-[0_0_15px_rgba(34,211,238,0.6)] scale-105' 
                  : 'bg-gray-300 hover:bg-gray-400' 
                }
              `}
            >
              {tab.label}
            </button>
          );
        })}
      </div>

      {/* Область контента */}
      <div className="w-full max-w-4xl bg-white p-8 rounded-2xl shadow-sm border border-gray-100">
        {activeTab === 'details' && <DetailsTab data={hero} />}
        {activeTab === 'skins' && <SkinsTab data={hero} />}
        {activeTab === 'emotes' && <EmotesTab data={hero} />}
        {activeTab === 'super-attacks' && <SuperAttacksTab data={hero} />}
      </div>

    </div>
  );
}