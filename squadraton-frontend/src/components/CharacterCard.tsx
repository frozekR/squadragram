const roleStyles: Record<string, {
  border: string;
  shadow: string;
}> = {
  ranged: {
    border: 'border-blue-500',
    shadow: 'shadow-[0_0_15px_rgba(59,130,246,0.5)]',
  },
  melee: {
    border: 'border-red-500',
    shadow: 'shadow-[0_0_15px_rgba(239,68,68,0.5)]',
  },
  default: {
    border: 'border-gray-500',
    shadow: 'shadow-md',
  }
};

interface CharacterCardProps {
  id: string | number;
  name: string;
  imageUrl: string;
  role?: string;
  isNew?: boolean;
  variant?: 'card' | 'banner';
}

export default function CharacterCard({ 
  id, 
  name, 
  imageUrl, 
  role = 'default', 
  isNew = false,
  variant = 'card' // 'card' (квадрат) или 'banner' (горизонтальный)
}: CharacterCardProps) {
  // Получаем стили в зависимости от роли (приводим к нижнему регистру на случай опечаток)
  const styles = roleStyles[role.toLowerCase()] || roleStyles.default;
  const linkHref = `/character/${id}`;

    return (
      <a 
        href={linkHref} 
        className={`
          relative block w-48 h-48 rounded-2xl overflow-hidden 
          border-4 transition-all duration-300 group
          ${styles.border} ${styles.shadow} ${styles.hoverShadow}
        `}
      >
        {/* Изображение персонажа */}
        <img 
          src={imageUrl} 
          alt={name} 
          className="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500 ease-out"
        />
        
        {/* Градиентная подложка снизу для читаемости текста */}
        <div className="absolute inset-0 bg-gradient-to-t from-black/90 via-black/20 to-transparent" />
        
        {/* Имя персонажа */}
        <div className="absolute bottom-3 left-0 right-0 text-center">
          <span className="text-white font-bold text-lg drop-shadow-md tracking-wide">
            {name}
          </span>
        </div>
      </a>
    );
}