import { Action } from './Action';

export const ChevronButton = ({
  collapsed,
  onToggle,
  direction,
}: {
  collapsed: boolean;
  onToggle: () => void;
  direction: 'left' | 'right';
}) => {
  const alignment =
    direction === 'left'
      ? collapsed
        ? 'flex-end'
        : 'flex-start'
      : collapsed
        ? 'flex-start'
        : 'flex-end';

  return (
    <div
      style={{
        display: 'flex',
        justifyContent: alignment,
      }}
    >
      <Action
        icon={direction === 'left' ? 'ChevronLeft' : 'ChevronRight'}
        iconOff={direction === 'left' ? 'ChevronRight' : 'ChevronLeft'}
        isOn={!collapsed}
        onClick={onToggle}
        variant="subtle"
        size="sm"
        radius="md"
        style={{
          fontWeight: 'normal',
          paddingRight: direction === 'right' ? '0.5rem' : '',
          paddingLeft: direction !== 'right' ? '0.5rem' : '',
        }}
      />
    </div>
  );
};
